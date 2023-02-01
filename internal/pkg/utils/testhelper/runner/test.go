package runner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/shlex"
	"github.com/keboola/go-utils/pkg/orderedmap"
	"github.com/keboola/go-utils/pkg/wildcards"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"

	"github.com/keboola/keboola-as-code/internal/pkg/encoding/json"
	"github.com/keboola/keboola-as-code/internal/pkg/env"
	"github.com/keboola/keboola-as-code/internal/pkg/filesystem"
	"github.com/keboola/keboola-as-code/internal/pkg/filesystem/aferofs"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/testhelper"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/testproject"
	"github.com/keboola/keboola-as-code/internal/pkg/validator"
)

const (
	dumpDirCtxKey        = ctxKey("dumpDir")
	envFileName          = "env"
	expectedStatePath    = "expected-state.json"
	expectedStdoutPath   = "expected-server-stdout"
	expectedStderrPath   = "expected-server-stderr"
	inDirName            = `in`
	initialStateFileName = "initial-state.json"
)

type ctxKey string

type Options func(c *runConfig)

type runConfig struct {
	addEnvVarsFromFile bool
	assertDirContent   bool
	assertEtcdState    bool
	assertProjectState bool
	cliBinaryPath      string
	copyInToWorkingDir bool
	initProjectState   bool
	apiServerConfig    apiServerConfig
}

func WithAddEnvVarsFromFile() Options {
	return func(c *runConfig) {
		c.addEnvVarsFromFile = true
	}
}

func WithAssertDirContent() Options {
	return func(c *runConfig) {
		c.assertDirContent = true
	}
}

func WithAssertEtcdState() Options {
	return func(c *runConfig) {
		c.assertEtcdState = true
	}
}

func WithAssertProjectState() Options {
	return func(c *runConfig) {
		c.assertProjectState = true
	}
}

func WithCopyInToWorkingDir() Options {
	return func(c *runConfig) {
		c.copyInToWorkingDir = true
	}
}

func WithInitProjectState() Options {
	return func(c *runConfig) {
		c.initProjectState = true
	}
}

func WithRunAPIServerAndRequests(
	path string,
	args []string,
	envs *env.Map,
	requestDecoratorFn func(request *APIRequest),
) Options {
	return func(c *runConfig) {
		c.apiServerConfig = apiServerConfig{
			path:               path,
			args:               args,
			envs:               envs,
			requestDecoratorFn: requestDecoratorFn,
		}
	}
}

func WithRunCLIBinary(path string) Options {
	return func(c *runConfig) {
		c.cliBinaryPath = path
	}
}

type Test struct {
	Runner
	ctx          context.Context
	env          *env.Map
	envProvider  testhelper.EnvProvider
	project      *testproject.Project
	t            *testing.T
	testDir      string
	testDirFS    filesystem.Fs
	workingDir   string
	workingDirFS filesystem.Fs
}

func (t *Test) EnvProvider() testhelper.EnvProvider {
	return t.envProvider
}

func (t *Test) T() *testing.T {
	return t.t
}

func (t *Test) TestDirFS() filesystem.Fs {
	return t.testDirFS
}

func (t *Test) WorkingDirFS() filesystem.Fs {
	return t.workingDirFS
}

func (t *Test) Run(opts ...Options) {
	t.t.Helper()

	c := runConfig{}
	for _, o := range opts {
		o(&c)
	}

	if c.copyInToWorkingDir {
		// Copy .in to the working dir of the current test.
		t.copyInToWorkingDir()
	}

	if c.initProjectState {
		// Set initial project state from the test file.
		t.initProjectState()
	}

	if c.addEnvVarsFromFile {
		// Load additional env vars from the test file.
		t.addEnvVarsFromFile()
	}

	// Replace all %%ENV_VAR%% in all files of the working directory.
	testhelper.MustReplaceEnvsDir(t.workingDirFS, `/`, t.envProvider)

	if c.cliBinaryPath != "" {
		// Run a CLI binary
		t.runCLIBinary(c.cliBinaryPath)
	}

	if c.apiServerConfig.path != "" {
		// Run an API server binary
		t.runAPIServer(
			c.apiServerConfig.path,
			c.apiServerConfig.args,
			c.apiServerConfig.envs,
			c.apiServerConfig.requestDecoratorFn,
		)
	}

	if c.assertDirContent {
		t.assertDirContent()
	}

	if c.assertProjectState {
		t.assertProjectState()
	}
}

func (t *Test) copyInToWorkingDir() {
	if !t.testDirFS.IsDir(inDirName) {
		t.t.Fatalf(`Missing directory "%s" in "%s".`, inDirName, t.testDir)
	}
	assert.NoError(t.t, aferofs.CopyFs2Fs(t.testDirFS, inDirName, t.workingDirFS, `/`))
}

func (t *Test) initProjectState() {
	if t.testDirFS.IsFile(initialStateFileName) {
		err := t.project.SetState(filesystem.Join(t.testDir, initialStateFileName))
		assert.NoError(t.t, err)
	}
}

func (t *Test) addEnvVarsFromFile() {
	if t.testDirFS.Exists(envFileName) {
		envFile, err := t.testDirFS.ReadFile(filesystem.NewFileDef(envFileName))
		if err != nil {
			t.t.Fatalf(`Cannot load "%s" file %s`, envFileName, err)
		}

		// Replace all %%ENV_VAR%% in "env" file
		envFileContent := testhelper.MustReplaceEnvsString(envFile.Content, t.envProvider)

		// Parse "env" file
		envsFromFile, err := env.LoadEnvString(envFileContent)
		if err != nil {
			t.t.Fatalf(`Cannot load "%s" file: %s`, envFileName, err)
		}

		// Merge
		t.env.Merge(envsFromFile, true)
	}
}

func (t *Test) runCLIBinary(path string) {
	// Load command arguments from file
	argsFile, err := t.TestDirFS().ReadFile(filesystem.NewFileDef("args"))
	if err != nil {
		t.T().Fatalf(`cannot open "%s" test file %s`, "args", err)
	}

	// Load and parse command arguments
	argsStr := strings.TrimSpace(argsFile.Content)
	argsStr = testhelper.MustReplaceEnvsString(argsStr, t.EnvProvider())
	args, err := shlex.Split(argsStr)
	if err != nil {
		t.T().Fatalf(`Cannot parse args "%s": %s`, argsStr, err)
	}

	// Prepare command
	cmd := exec.CommandContext(t.ctx, path, args...) // nolint:gosec
	cmd.Env = t.env.ToSlice()
	cmd.Dir = t.workingDir

	// Setup command input/output
	cmdInOut, err := setupCmdInOut(t.t, t.envProvider, t.testDirFS, cmd)
	if err != nil {
		t.t.Fatal(err.Error())
	}

	// Start command
	if err := cmd.Start(); err != nil {
		t.t.Fatalf("Cannot start command: %s", err)
	}

	// Always terminate the command
	defer func() {
		_ = cmd.Process.Kill()
	}()

	// Error handler for errors in interaction
	interactionErrHandler := func(err error) {
		if err != nil {
			t.t.Fatal(err)
		}
	}

	// Wait for command
	exitCode := 0
	err = cmdInOut.InteractAndWait(t.ctx, cmd, interactionErrHandler)
	if err != nil {
		t.t.Logf(`cli command failed: %s`, err.Error())
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			exitCode = exitError.ExitCode()
		} else {
			t.t.Fatalf("Command failed: %s", err)
		}
	}

	// Get outputs
	stdout := cmdInOut.StdoutString()
	stderr := cmdInOut.StderrString()

	expectedCode := cast.ToInt(t.ReadFileFromTestDir("expected-code"))
	assert.Equal(
		t.t,
		expectedCode,
		exitCode,
		"Unexpected exit code.\nSTDOUT:\n%s\n\nSTDERR:\n%s\n\n",
		stdout,
		stderr,
	)

	expectedStdout := t.ReadFileFromTestDir("expected-stdout")
	wildcards.Assert(t.t, expectedStdout, stdout, "Unexpected STDOUT.")

	expectedStderr := t.ReadFileFromTestDir("expected-stderr")
	wildcards.Assert(t.t, expectedStderr, stderr, "Unexpected STDERR.")
}

type apiServerConfig struct {
	path               string
	args               []string
	envs               *env.Map
	requestDecoratorFn func(request *APIRequest)
}

func (t *Test) runAPIServer(
	path string,
	addArgs []string,
	addEnvs *env.Map,
	requestDecoratorFn func(request *APIRequest),
) {
	// Get a free port
	port, err := getFreePort()
	if err != nil {
		t.t.Fatalf("Could not receive a free port: %s", err)
	}
	apiURL := fmt.Sprintf("http://localhost:%d", port)

	args := append([]string{fmt.Sprintf("--http-port=%d", port)}, addArgs...)

	// Envs
	envs := env.Empty()
	envs.Set("PATH", os.Getenv("PATH")) // nolint:forbidigo
	envs.Set("KBC_STORAGE_API_HOST", t.project.StorageAPIHost())
	envs.Set("DATADOG_ENABLED", "false")
	envs.Merge(addEnvs, false)

	// Start API server
	stdout := newCmdOut()
	stderr := newCmdOut()
	cmd := exec.Command(path, args...)
	cmd.Env = envs.ToSlice()
	cmd.Stdout = io.MultiWriter(stdout, testhelper.VerboseStdout())
	cmd.Stderr = io.MultiWriter(stderr, testhelper.VerboseStderr())
	if err := cmd.Start(); err != nil {
		t.t.Fatalf("Server failed to start: %s", err)
	}

	cmdWaitCh := make(chan error, 1)
	go func() {
		cmdWaitCh <- cmd.Wait()
		close(cmdWaitCh)
	}()

	// Kill API server after test
	t.t.Cleanup(func() {
		_ = cmd.Process.Kill()
	})

	// Wait for API server
	if err = waitForAPI(cmdWaitCh, apiURL); err != nil {
		t.t.Fatalf(
			"Unexpected error while waiting for API: %s\n\nServer STDERR:%s\n\nServer STDOUT:%s\n",
			err,
			stderr.String(),
			stdout.String(),
		)
	}

	// Run the requests
	requestsOk := t.runRequests(apiURL, requestDecoratorFn)

	// Shutdown API server
	_ = cmd.Process.Signal(syscall.SIGTERM)
	select {
	case <-time.After(10 * time.Second):
		t.t.Fatalf("timeout while waiting for server shutdown")
	case <-cmdWaitCh:
		// continue
	}

	// Dump process stdout/stderr
	assert.NoError(t.t, t.workingDirFS.WriteFile(filesystem.NewRawFile("process-stdout.txt", stdout.String())))
	assert.NoError(t.t, t.workingDirFS.WriteFile(filesystem.NewRawFile("process-stderr.txt", stderr.String())))

	// Check API server stdout/stderr
	if requestsOk {
		if t.testDirFS.IsFile(expectedStdoutPath) {
			expected := t.ReadFileFromTestDir(expectedStdoutPath)
			wildcards.Assert(t.t, expected, stdout.String(), "Unexpected STDOUT.")
		}
		if t.testDirFS.IsFile(expectedStderrPath) {
			expected := t.ReadFileFromTestDir(expectedStderrPath)
			wildcards.Assert(t.t, expected, stderr.String(), "Unexpected STDERR.")
		}
	}
}

type APIRequest struct {
	Path    string            `json:"path" validate:"required"`
	Method  string            `json:"method" validate:"required,oneof=DELETE GET PATCH POST PUT"`
	Body    interface{}       `json:"body"`
	Headers map[string]string `json:"headers"`
}

func (t *Test) runRequests(apiURL string, requestDecoratorFn func(*APIRequest)) bool {
	client := resty.New()
	client.SetBaseURL(apiURL)

	// Dump raw HTTP request
	client.SetPreRequestHook(func(client *resty.Client, request *http.Request) error {
		if dumpDir, ok := request.Context().Value(dumpDirCtxKey).(string); ok {
			reqDump, err := httputil.DumpRequest(request, true)
			assert.NoError(t.t, err)
			assert.NoError(t.t, t.workingDirFS.WriteFile(filesystem.NewRawFile(filesystem.Join(dumpDir, "request.txt"), string(reqDump))))
		}
		return nil
	})

	// Request folders should be named e.g. 001-request1, 002-request2
	dirs, err := t.testDirFS.Glob("[0-9][0-9][0-9]-*")
	assert.NoError(t.t, err)
	for _, requestDir := range dirs {
		// Read the request file
		requestFileStr := t.ReadFileFromTestDir(filesystem.Join(requestDir, "request.json"))
		assert.NoError(t.t, err)

		request := &APIRequest{}
		err = json.DecodeString(requestFileStr, request)
		assert.NoError(t.t, err)
		err = validator.New().Validate(context.Background(), request)
		assert.NoError(t.t, err)

		// Send the request
		r := client.R()
		if request.Body != nil {
			if v, ok := request.Body.(string); ok {
				r.SetBody(v)
			} else if v, ok := request.Body.(map[string]any); ok && resty.IsJSONType(request.Headers["Content-Type"]) {
				r.SetBody(v)
			} else {
				assert.FailNow(t.t, fmt.Sprintf("request.json for request %s is malformed, body must be JSON for proper JSON content type or string otherwise", requestDir))
			}
		}
		r.SetHeaders(request.Headers)

		// Decorate the request
		requestDecoratorFn(request)

		// Send request
		r.SetContext(context.WithValue(r.Context(), dumpDirCtxKey, requestDir))
		resp, err := r.Execute(request.Method, request.Path)
		assert.NoError(t.t, err)

		// Dump raw HTTP response
		if err == nil {
			respDump, err := httputil.DumpResponse(resp.RawResponse, false)
			assert.NoError(t.t, err)
			assert.NoError(t.t, t.workingDirFS.WriteFile(filesystem.NewRawFile(filesystem.Join(requestDir, "response.txt"), string(respDump)+string(resp.Body()))))
		}

		// Compare response body
		expectedRespBody := t.ReadFileFromTestDir(filesystem.Join(requestDir, "expected-response.json"))

		// Decode && encode json to unite indentation of the response with expected-response.json
		respMap := orderedmap.New()
		if resp.String() != "" {
			err = json.DecodeString(resp.String(), &respMap)
		}
		assert.NoError(t.t, err)
		respBody, err := json.EncodeString(respMap, true)
		assert.NoError(t.t, err)

		// Compare response status code
		expectedCode := cast.ToInt(t.ReadFileFromTestDir(filesystem.Join(requestDir, "expected-http-code")))
		ok1 := assert.Equal(
			t.t,
			expectedCode,
			resp.StatusCode(),
			"Unexpected status code for request \"%s\".\nRESPONSE:\n%s\n\n",
			requestDir,
			resp.String(),
		)

		// Assert response body
		ok2 := wildcards.Assert(t.t, expectedRespBody, respBody, fmt.Sprintf("Unexpected response for request %s.", requestDir))

		// If the request failed, skip other requests
		if !ok1 || !ok2 {
			t.t.Errorf(`request "%s" failed, skipping the other requests`, requestDir)
			return false
		}
	}

	return true
}

func (t *Test) ReadFileFromTestDir(path string) string {
	file, err := t.testDirFS.ReadFile(filesystem.NewFileDef(path))
	assert.NoError(t.t, err)
	return testhelper.MustReplaceEnvsString(strings.TrimSpace(file.Content), t.envProvider)
}

func (t *Test) assertDirContent() {
	// Expected state dir
	expectedDir := "out"
	if !t.testDirFS.IsDir(expectedDir) {
		t.t.Fatalf(`Missing directory "%s" in "%s".`, expectedDir, t.testDirFS.BasePath())
	}

	// Copy expected state and replace ENVs
	expectedDirFS := aferofs.NewMemoryFsFrom(filesystem.Join(t.testDirFS.BasePath(), expectedDir))
	testhelper.MustReplaceEnvsDir(expectedDirFS, `/`, t.envProvider)

	// Compare actual and expected dirs
	testhelper.AssertDirectoryContentsSame(t.t, expectedDirFS, `/`, t.workingDirFS, `/`)
}

func (t *Test) assertProjectState() {
	if t.testDirFS.IsFile(expectedStatePath) {
		expectedState := t.ReadFileFromTestDir(expectedStatePath)

		// Load actual state
		actualState, err := t.project.NewSnapshot()
		assert.NoError(t.t, err)

		// Write actual state
		err = t.workingDirFS.WriteFile(filesystem.NewRawFile("actual-state.json", json.MustEncodeString(actualState, true)))
		assert.NoError(t.t, err)

		// Compare expected and actual state
		wildcards.Assert(
			t.t,
			testhelper.MustReplaceEnvsString(expectedState, t.envProvider),
			json.MustEncodeString(actualState, true),
			`unexpected project state, compare "expected-state.json" from test and "actual-state.json" from ".out" dir.`,
		)
	}
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func waitForAPI(cmdErrCh <-chan error, apiURL string) error {
	client := resty.New()

	serverStartTimeout := 45 * time.Second
	timeout := time.After(serverStartTimeout)
	tick := time.Tick(200 * time.Millisecond) // nolint:staticcheck
	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		// Handle timeout
		case <-timeout:
			return errors.Errorf("server didn't start within %s", serverStartTimeout)
		// Handle server termination
		case err := <-cmdErrCh:
			if err == nil {
				return errors.New("the server was terminated unexpectedly")
			} else {
				return errors.Errorf("the server was terminated unexpectedly with error: %w", err)
			}
		// Periodically test health check endpoint
		case <-tick:
			resp, err := client.R().Get(fmt.Sprintf("%s/health-check", apiURL))
			if err != nil && !strings.Contains(err.Error(), "connection refused") {
				return err
			}
			if resp.StatusCode() == 200 {
				return nil
			}
		}
	}
}

// cmdOut is used to prevent race conditions, see https://hackmysql.com/post/reading-os-exec-cmd-output-without-race-conditions/
type cmdOut struct {
	buf  *bytes.Buffer
	lock *sync.Mutex
}

func newCmdOut() *cmdOut {
	return &cmdOut{buf: &bytes.Buffer{}, lock: &sync.Mutex{}}
}

func (o *cmdOut) Write(p []byte) (int, error) {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.buf.Write(p)
}

func (o *cmdOut) String() string {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.buf.String()
}