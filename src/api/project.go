package api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"keboola-as-code/src/client"
	"keboola-as-code/src/fixtures"
	"keboola-as-code/src/model"
	"keboola-as-code/src/utils"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// testProject manages cleanup and setup of the test project
type testProject struct {
	t             *testing.T
	lock          *sync.Mutex
	testDir       string
	stateFile     *fixtures.StateFile
	api           *StorageApi
	defaultBranch *model.Branch
	envs          []string
}

func SetStateOfTestProject(t *testing.T, projectStateFilePath string) {
	p := newTestProject(t, projectStateFilePath)
	p.Clear()
	p.InitState()
}

// newTestProject creates testProject and loads state from the stateFilePath
func newTestProject(t *testing.T, stateFilePath string) *testProject {
	_, testFile, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(testFile)

	// Load state file
	stateFile, err := fixtures.LoadStateFile(stateFilePath)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	// Create API
	a, _ := TestStorageApiWithToken(t)
	if utils.TestProjectId() != a.ProjectId() {
		assert.FailNow(t, "TEST_PROJECT_ID and token project id are different.")
	}

	// Load default branch
	defaultBranch, err := a.GetDefaultBranch()
	if err != nil {
		assert.FailNow(t, "cannot get default branch")
	}

	// Create
	p := &testProject{t, &sync.Mutex{}, testDir, stateFile, a, defaultBranch, nil}
	p.log("Initializing test project \"%s\", id: \"%d\".", p.api.ProjectName(), p.api.ProjectId())
	return p
}

// Clear deletes all project branches (except default) and all configurations
func (p *testProject) Clear() {
	p.log("Clearing project ...")
	startTime := time.Now()

	// Delete all configs in default branch, it cannot be deleted
	pool := p.api.NewPool()
	pool.Request(p.api.DeleteConfigsInBranchRequest(p.defaultBranch.Id)).Send()
	if err := pool.StartAndWait(); err != nil {
		assert.FailNow(p.t, fmt.Sprintf("cannot delete branches: %s", err))
	}

	// Load branches
	branches, err := p.api.ListBranches()
	if err != nil {
		assert.FailNow(p.t, fmt.Sprintf("cannot load branches: %s", err))
	}

	// Delete all dev-branches sequentially, parallel requests don't work with this endpoint
	for _, branch := range *branches {
		if !branch.IsDefault {
			p.api.DeleteBranchRequest(branch.Id).Send()
		}
	}

	p.log("Test project cleared | %s", time.Since(startTime))
}

// InitState creates branches and configurations according stateFile
func (p *testProject) InitState() {
	startTime := time.Now()
	p.log("Setting project state ...")

	// Create configs in default branch, they will be auto-copied to dev-branches
	pool := p.api.NewPool()
	p.CreateConfigsInBranch(pool, p.stateFile.AllBranchesConfigs, p.defaultBranch, "TEST_BRANCH_ALL_CONFIG")
	if err := pool.StartAndWait(); err != nil {
		assert.FailNow(p.t, fmt.Sprintf("cannot create configs in default branch: %s", err))
	}

	// Create branches sequentially, parallel requests don't work with this endpoint
	branchesByName := make(map[string]*model.Branch)
	for _, fixture := range p.stateFile.Branches {
		branch := fixture.Branch.ToModel(p.defaultBranch)
		branchesByName[branch.Name] = branch
		if branch.IsDefault {
			p.defaultBranch.Description = fixture.Branch.Description
			if _, err := p.api.UpdateBranch(p.defaultBranch, []string{"description"}); err != nil {
				assert.FailNow(p.t, fmt.Sprintf("cannot set default branch description: %s", err))
			}
			p.setEnv(fmt.Sprintf("TEST_BRANCH_%s_ID", branch.Name), strconv.Itoa(branch.Id))
		} else {
			p.api.
				CreateBranchRequest(branch).
				OnSuccess(func(response *client.Response) *client.Response {
					p.log(`crated branch "%s", id: "%d"`, branch.Name, branch.Id)
					p.setEnv(fmt.Sprintf("TEST_BRANCH_%s_ID", branch.Name), strconv.Itoa(branch.Id))
					return response
				}).
				Send()
		}
	}

	// Create configs in dev-branches
	pool = p.api.NewPool()
	for _, branch := range p.stateFile.Branches {
		modelBranch := branchesByName[branch.Branch.Name]
		envPrefix := fmt.Sprintf("TEST_BRANCH_%s_CONFIG", modelBranch.Name)
		p.CreateConfigsInBranch(pool, branch.Configs, modelBranch, envPrefix)
	}
	if err := pool.StartAndWait(); err != nil {
		assert.FailNow(p.t, fmt.Sprintf("cannot create configs: %s", err))
	}

	// Log ENVs
	for _, env := range p.envs {
		p.log(fmt.Sprintf(`Set ENV "%s"`, env))
	}

	// Done
	p.log("Project state set | %s", time.Since(startTime))
}

// CreateConfigsInBranch loads configs from JSON files and creates them in the test project
func (p *testProject) CreateConfigsInBranch(pool *client.Pool, names []string, branch *model.Branch, envPrefix string) {
	for _, name := range names {
		config := fixtures.LoadConfig(p.t, name)
		config.BranchId = branch.Id
		if request, err := p.api.CreateConfigRequest(config); err == nil {
			p.log("creating config \"%s/%s/%s\"", branch.Name, config.ComponentId, config.Name)
			pool.
				Request(request).
				OnSuccess(func(response *client.Response) *client.Response {
					p.setEnv(fmt.Sprintf("%s_%s_ID", envPrefix, config.Name), config.Id)
					return response
				}).
				Send()
		} else {
			assert.FailNow(p.t, fmt.Sprintf("cannot create create config request: %s", err))
		}
	}
}

// setEnv set ENV variable, all ENVs are logged at the end of InitState method
func (p *testProject) setEnv(key string, value string) {
	// Normalize key
	key = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(key, "_")
	key = strings.ToUpper(key)
	key = strings.Trim(key, "_")

	// Set
	utils.MustSetEnv(key, value)

	// Log
	p.lock.Lock()
	defer p.lock.Unlock()
	p.envs = append(p.envs, fmt.Sprintf("%s=%s", key, value))
}

func (p *testProject) log(format string, a ...interface{}) {
	if utils.TestIsVerbose() {
		a = append([]interface{}{p.t.Name()}, a...)
		fmt.Println(fmt.Sprintf("Fixtures[%s]: "+format, a...))
	}
}