package diff

import (
	"context"
	"github.com/otiai10/copy"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	"keboola-as-code/src/api"
	"keboola-as-code/src/model"
	"keboola-as-code/src/utils"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDifferLoadState(t *testing.T) {
	defer utils.ResetEnv(t, os.Environ())
	api.SetStateOfTestProject(t, "minimal.json")

	// Same IDs in local and remote state
	utils.MustSetEnv("LOCAL_STATE_MAIN_BRANCH_ID", utils.MustGetEnv(`TEST_BRANCH_MAIN_ID`))
	utils.MustSetEnv("LOCAL_STATE_GENERIC_CONFIG_ID", utils.MustGetEnv(`TEST_BRANCH_ALL_CONFIG_EMPTY_ID`))
	projectDir, metadataDir := initLocalState(t, "minimal")
	differ := createDiffer(t, projectDir, metadataDir)

	assert.NoError(t, differ.LoadState())
	assert.Empty(t, differ.state.RemoteErrors())
	assert.Empty(t, differ.state.LocalErrors())
	assert.True(t, differ.stateLoaded)
	assert.Equal(t, []*model.BranchState{
		{
			Id: cast.ToInt(utils.MustGetEnv(`TEST_BRANCH_MAIN_ID`)),
			Remote: &model.Branch{
				Id:          cast.ToInt(utils.MustGetEnv(`TEST_BRANCH_MAIN_ID`)),
				Name:        "Main",
				Description: "Main branch",
				IsDefault:   true,
			},
			Local: &model.Branch{
				Id:          cast.ToInt(utils.MustGetEnv(`TEST_BRANCH_MAIN_ID`)),
				Name:        "Main",
				Description: "Main branch",
				IsDefault:   true,
			},
			Manifest: &model.BranchManifest{
				Id:           cast.ToInt(utils.MustGetEnv(`TEST_BRANCH_MAIN_ID`)),
				Path:         "main",
				RelativePath: "main",
				MetadataFile: "main/meta.json",
			},
		},
	}, differ.state.Branches())
	assert.Equal(t, []*model.ConfigState{
		{
			BranchId:    cast.ToInt(utils.MustGetEnv(`TEST_BRANCH_MAIN_ID`)),
			ComponentId: "ex-generic-v2",
			Id:          utils.MustGetEnv(`TEST_BRANCH_ALL_CONFIG_EMPTY_ID`),
			Remote: &model.Config{
				BranchId:          cast.ToInt(utils.MustGetEnv(`TEST_BRANCH_MAIN_ID`)),
				ComponentId:       "ex-generic-v2",
				Id:                utils.MustGetEnv(`TEST_BRANCH_ALL_CONFIG_EMPTY_ID`),
				Name:              "empty",
				Description:       "test fixture",
				ChangeDescription: "created by test",
				Config:            map[string]interface{}{},
				Rows:              []*model.ConfigRow{},
			},
			Local: &model.Config{
				BranchId:          cast.ToInt(utils.MustGetEnv(`TEST_BRANCH_MAIN_ID`)),
				ComponentId:       "ex-generic-v2",
				Id:                utils.MustGetEnv(`TEST_BRANCH_ALL_CONFIG_EMPTY_ID`),
				Name:              "todos",
				Description:       "todos config",
				ChangeDescription: "",
				Config: map[string]interface{}{
					"parameters": map[string]interface{}{
						"api": map[string]interface{}{
							"baseUrl": "https://jsonplaceholder.typicode.com",
						},
					},
				},
				Rows: []*model.ConfigRow{},
			},
			Manifest: &model.ConfigManifest{
				BranchId:     cast.ToInt(utils.MustGetEnv(`TEST_BRANCH_MAIN_ID`)),
				ComponentId:  "ex-generic-v2",
				Id:           utils.MustGetEnv(`TEST_BRANCH_ALL_CONFIG_EMPTY_ID`),
				Path:         "ex-generic-v2/456-todos",
				Rows:         []*model.ConfigRowManifest{},
				RelativePath: "main/ex-generic-v2/456-todos",
				MetadataFile: "main/ex-generic-v2/456-todos/meta.json",
				ConfigFile:   "main/ex-generic-v2/456-todos/config.json",
			},
		},
	}, differ.state.Configs())
}

func initLocalState(t *testing.T, localState string) (string, string) {
	_, testFile, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(testFile)
	localStateDir := filepath.Join(testDir, "..", "fixtures", "local", localState)
	projectDir := t.TempDir()
	metadataDir := filepath.Join(projectDir, ".keboola")
	err := copy.Copy(localStateDir, projectDir)
	if err != nil {
		t.Fatalf("Copy error: %s", err)
	}
	utils.ReplaceEnvsDir(projectDir)
	return projectDir, metadataDir
}

func createDiffer(t *testing.T, projectDir, metadataDir string) *Differ {
	a, _ := api.TestStorageApiWithToken(t)
	logger, _ := utils.NewDebugLogger()
	return NewDiffer(projectDir, metadataDir, context.Background(), a, logger)
}
