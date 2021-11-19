package local

import (
	"context"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/keboola/keboola-as-code/internal/pkg/env"
	"github.com/keboola/keboola-as-code/internal/pkg/filesystem"
	"github.com/keboola/keboola-as-code/internal/pkg/filesystem/aferofs"
	"github.com/keboola/keboola-as-code/internal/pkg/json"
	"github.com/keboola/keboola-as-code/internal/pkg/manifest"
	"github.com/keboola/keboola-as-code/internal/pkg/model"
	"github.com/keboola/keboola-as-code/internal/pkg/testhelper"
	"github.com/keboola/keboola-as-code/internal/pkg/utils"
)

type testMapper struct {
	onLoadCalls []string
}

func (*testMapper) MapBeforeLocalSave(recipe *model.LocalSaveRecipe) error {
	if _, ok := recipe.Object.(*model.Config); ok {
		recipe.Configuration.Content.Set("key", "overwritten")
		recipe.Configuration.Content.Set("new", "value")
	}
	return nil
}

func (*testMapper) MapAfterLocalLoad(recipe *model.LocalLoadRecipe) error {
	if _, ok := recipe.Object.(*model.Config); ok {
		recipe.Configuration.Content.Set("parameters", "overwritten")
		recipe.Configuration.Content.Set("new", "value")
	}
	return nil
}

func (t *testMapper) OnObjectsLoad(event model.OnObjectsLoadEvent) error {
	for _, object := range event.NewObjects {
		t.onLoadCalls = append(t.onLoadCalls, object.Desc())
	}
	return nil
}

func TestLocalUnitOfWork_workersFor(t *testing.T) {
	t.Parallel()
	manager, _ := newTestLocalManager(t)
	uow := manager.NewUnitOfWork(context.Background())

	lock := &sync.Mutex{}
	var order []int

	for _, level := range []int{3, 2, 4, 1} {
		level := level
		uow.workersFor(level).AddWorker(func() error {
			lock.Lock()
			defer lock.Unlock()
			order = append(order, level)
			return nil
		})
		uow.workersFor(level).AddWorker(func() error {
			lock.Lock()
			defer lock.Unlock()
			order = append(order, level)
			return nil
		})
	}

	// Not started
	time.Sleep(10 * time.Millisecond)
	assert.Empty(t, order)

	// Invoke
	assert.NoError(t, uow.Invoke())
	assert.Equal(t, []int{1, 1, 2, 2, 3, 3, 4, 4}, order)

	// Cannot be reused
	assert.PanicsWithError(t, `invoked local.UnitOfWork cannot be reused`, func() {
		uow.Invoke()
	})
}

func TestBeforeLocalSaveMapper(t *testing.T) {
	t.Parallel()
	manager, mapper := newTestLocalManager(t)
	fs := manager.Fs()
	uow := manager.NewUnitOfWork(context.Background())

	// Add test mapper
	testMapperInst := &testMapper{}
	mapper.AddMapper(testMapperInst)

	// Test object
	configKey := model.ConfigKey{BranchId: 123, ComponentId: `foo.bar`, Id: `456`}
	configState := &model.ConfigState{
		ConfigManifest: &model.ConfigManifest{
			ConfigKey: configKey,
			Paths: model.Paths{
				PathInProject: model.NewPathInProject(`branch`, `config`),
			},
		},
		Remote: &model.Config{
			ConfigKey: configKey,
			Name:      "name",
			Content: utils.PairsToOrderedMap([]utils.Pair{
				{Key: "key", Value: "internal value"},
			}),
		},
	}

	// Save object
	uow.SaveObject(configState, configState.Remote, model.ChangedFields{})
	assert.NoError(t, uow.Invoke())

	// File content has been mapped
	configFile, err := fs.ReadFile(filesystem.Join(`branch`, `config`, model.ConfigFile), `config file`)
	assert.NoError(t, err)
	assert.Equal(t, "{\n  \"key\": \"overwritten\",\n  \"new\": \"value\"\n}", strings.TrimSpace(configFile.Content))
}

func TestAfterLocalLoadMapper(t *testing.T) {
	t.Parallel()
	manager, mapper := newTestLocalManager(t)
	fs := manager.Fs()
	uow := manager.NewUnitOfWork(context.Background())

	// Add test mapper
	testMapperInst := &testMapper{}
	mapper.AddMapper(testMapperInst)

	// Init dir
	_, testFile, _, _ := runtime.Caller(0)
	testDir := filesystem.Dir(testFile)
	inputDir := filesystem.Join(testDir, `..`, `fixtures`, `local`, `minimal`)
	assert.NoError(t, aferofs.CopyFs2Fs(nil, inputDir, fs, ``))

	// Replace placeholders in files
	envs := env.Empty()
	envs.Set("TEST_KBC_STORAGE_API_HOST", "foo.bar")
	envs.Set("LOCAL_STATE_MAIN_BRANCH_ID", "111")
	envs.Set("LOCAL_STATE_GENERIC_CONFIG_ID", "456")
	testhelper.ReplaceEnvsDir(fs, `/`, envs)

	// Load objects
	m, err := manifest.LoadManifest(fs)
	assert.NoError(t, err)
	uow.LoadAll(m.Content)
	assert.NoError(t, uow.Invoke())

	// Internal state has been mapped
	configState := manager.state.MustGet(model.ConfigKey{BranchId: 111, ComponentId: `ex-generic-v2`, Id: `456`}).(*model.ConfigState)
	assert.Equal(t, `{"parameters":"overwritten","new":"value"}`, json.MustEncodeString(configState.Local.Content, false))

	// OnObjectsLoad event has been called
	assert.Equal(t, []string{`branch "111"`, `config "branch:111/component:ex-generic-v2/config:456"`}, testMapperInst.onLoadCalls)
}

func TestSkipChildrenLoadIfParentIsInvalid(t *testing.T) {
	t.Parallel()
	manager, _ := newTestLocalManager(t)
	fs := manager.Fs()
	uow := manager.NewUnitOfWork(context.Background())

	// Init dir
	_, testFile, _, _ := runtime.Caller(0)
	testDir := filesystem.Dir(testFile)
	inputDir := filesystem.Join(testDir, `..`, `fixtures`, `local`, `branch-invalid-meta-json`)
	assert.NoError(t, aferofs.CopyFs2Fs(nil, inputDir, fs, ``))

	// Setup manifest
	branchManifest := &model.BranchManifest{
		BranchKey: model.BranchKey{Id: 123},
		Paths: model.Paths{
			PathInProject: model.NewPathInProject(``, `main`),
		},
	}
	configManifest := &model.ConfigManifestWithRows{
		ConfigManifest: &model.ConfigManifest{
			ConfigKey: model.ConfigKey{BranchId: 123, ComponentId: `foo.bar`, Id: `456`},
			Paths: model.Paths{
				PathInProject: model.NewPathInProject(`main`, `config`),
			},
		},
	}
	manager.manifest.Content.Branches = []*model.BranchManifest{branchManifest}
	manager.manifest.Content.Configs = []*model.ConfigManifestWithRows{configManifest}
	assert.False(t, branchManifest.State().IsInvalid())
	assert.False(t, configManifest.State().IsInvalid())
	assert.False(t, branchManifest.State().IsNotFound())
	assert.False(t, configManifest.State().IsNotFound())

	// Load all
	uow.LoadAll(manager.manifest.Content)

	// Invoke and check error
	// Branch is invalid, so config does not read at all (no error that it does not exist).
	err := uow.Invoke()
	expectedErr := `
branch metadata file "main/meta.json" is invalid:
  - invalid character 'f' looking for beginning of object key string, offset: 3
`
	assert.Error(t, err)
	assert.Equal(t, strings.Trim(expectedErr, "\n"), err.Error())

	// Check manifest records
	assert.True(t, branchManifest.State().IsInvalid())
	assert.True(t, configManifest.State().IsInvalid())
	assert.False(t, branchManifest.State().IsNotFound())
	assert.False(t, configManifest.State().IsNotFound())
}
