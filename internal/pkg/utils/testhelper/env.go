// nolint forbidigo
package testhelper

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/acarl005/stripansi"
	"github.com/spf13/cast"

	"github.com/keboola/keboola-as-code/internal/pkg/filesystem"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

type EnvProvider interface {
	GetOrErr(key string) (string, error)
	MustGet(key string) string
}

const envPlaceholderTemplate = `%s[a-zA-Z0-9][a-zA-Z0-9\-_]*[a-zA-Z0-9]%s`

// MustReplaceEnvsStringWithSeparator replaces ENVs in given string with chosen separator.
func MustReplaceEnvsStringWithSeparator(str string, provider EnvProvider, envSeparator string) string {
	return regexp.
		MustCompile(fmt.Sprintf(envPlaceholderTemplate, envSeparator, envSeparator)).
		ReplaceAllStringFunc(str, func(s string) string {
			return provider.MustGet(strings.Trim(s, envSeparator))
		})
}

// ReplaceEnvsStringWithSeparator replaces ENVs in given string with chosen separator.
func ReplaceEnvsStringWithSeparator(str string, provider EnvProvider, envSeparator string) (string, error) {
	errs := errors.NewMultiError()
	res := regexp.
		MustCompile(fmt.Sprintf(envPlaceholderTemplate, envSeparator, envSeparator)).
		ReplaceAllStringFunc(str, func(s string) string {
			res, err := provider.GetOrErr(strings.Trim(s, envSeparator))
			if err != nil {
				errs.Append(err)
				return s
			}
			return res
		})
	return res, errs.ErrorOrNil()
}

// MustReplaceEnvsString replaces ENVs in given string.
func MustReplaceEnvsString(str string, provider EnvProvider) string {
	return MustReplaceEnvsStringWithSeparator(str, provider, "%%")
}

// MustReplaceEnvsFileWithSeparator replaces ENVs in given file with chosen separator.
func MustReplaceEnvsFileWithSeparator(fs filesystem.Fs, path string, provider EnvProvider, envSeparator string) {
	file, err := fs.ReadFile(filesystem.NewFileDef(path))
	if err != nil {
		panic(err)
	}
	file.Content = MustReplaceEnvsStringWithSeparator(file.Content, provider, envSeparator)
	if err := fs.WriteFile(file); err != nil {
		panic(errors.Errorf("cannot write to file \"%s\": %w", path, err))
	}
}

// MustReplaceEnvsFile replaces ENVs in given file.
func MustReplaceEnvsFile(fs filesystem.Fs, path string, provider EnvProvider) {
	MustReplaceEnvsFileWithSeparator(fs, path, provider, "%%")
}

// MustReplaceEnvsDirWithSeparator replaces ENVs in all files in root directory and sub-directories with chosen separator.
func MustReplaceEnvsDirWithSeparator(fs filesystem.Fs, root string, provider EnvProvider, envSeparator string) {
	// Iterate over directory structure
	err := fs.Walk(root, func(p string, info filesystem.FileInfo, err error) error {
		// Stop on error
		if err != nil {
			return err
		}

		// Ignore hidden files, except .env*, .gitignore
		if IsIgnoredFile(p, info) {
			return nil
		}

		if path.Ext(p) == ".sql" || path.Ext(p) == ".py" {
			return nil
		}

		// Process file
		if !info.IsDir() {
			MustReplaceEnvsFileWithSeparator(fs, p, provider, envSeparator)
		}

		return nil
	})
	if err != nil {
		panic(errors.Errorf("cannot walk over dir \"%s\": %w", root, err))
	}
}

// MustReplaceEnvsDir replaces ENVs in all files in root directory and sub-directories.
func MustReplaceEnvsDir(fs filesystem.Fs, root string, provider EnvProvider) {
	MustReplaceEnvsDirWithSeparator(fs, root, provider, "%%")
}

// stripAnsiWriter strips ANSI characters from
type stripAnsiWriter struct {
	lock   *sync.Mutex
	buf    *bytes.Buffer
	writer io.Writer
}

func newStripAnsiWriter(writer io.Writer) *stripAnsiWriter {
	return &stripAnsiWriter{
		lock:   &sync.Mutex{},
		buf:    &bytes.Buffer{},
		writer: writer,
	}
}

func (w *stripAnsiWriter) writeBuffer() error {
	if _, err := w.writer.Write([]byte(stripansi.Strip(w.buf.String()))); err != nil {
		return err
	}
	w.buf.Reset()
	return nil
}

func (w *stripAnsiWriter) Write(p []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	// Append to the buffer
	n, err := w.buf.Write(p)

	// We can only remove an ANSI escape seq if the whole expression is present.
	// ... so if buffer contains new line -> flush
	if bytes.Contains(w.buf.Bytes(), []byte("\n")) {
		if err := w.writeBuffer(); err != nil {
			return 0, err
		}
	}

	return n, err
}

func (w *stripAnsiWriter) Close() error {
	w.lock.Lock()
	defer w.lock.Unlock()

	if err := w.writeBuffer(); err != nil {
		return err
	}
	return nil
}

type nopCloser struct {
	io.Writer
}

func (n *nopCloser) Close() error {
	return nil
}

func TestIsVerbose() bool {
	value := os.Getenv("TEST_VERBOSE")
	if value == "" {
		value = "false"
	}
	return cast.ToBool(value)
}

func VerboseStdout() io.WriteCloser {
	if TestIsVerbose() {
		return newStripAnsiWriter(os.Stdout)
	}

	return &nopCloser{io.Discard}
}

func VerboseStderr() io.WriteCloser {
	if TestIsVerbose() {
		return newStripAnsiWriter(os.Stderr)
	}

	return &nopCloser{io.Discard}
}
