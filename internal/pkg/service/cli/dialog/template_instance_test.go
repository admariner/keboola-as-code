package dialog_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/keboola/keboola-as-code/internal/pkg/fixtures"
	"github.com/keboola/keboola-as-code/internal/pkg/model"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/dependencies"
	loadState "github.com/keboola/keboola-as-code/pkg/lib/operation/state/load"
)

func TestAskTemplateInstance_Interactive(t *testing.T) {
	t.Parallel()

	// Test dependencies
	dialog, _, console := createDialogs(t, true)
	d := dependencies.NewMocked(t)
	projectState, err := d.MockedProject(fixtures.MinimalProjectFs(t)).LoadState(loadState.Options{LoadLocalState: true}, d)
	assert.NoError(t, err)
	branch, _ := projectState.LocalObjects().Get(model.BranchKey{ID: 123})

	instanceID := "inst1"
	templateID := "tmpl1"
	version := "1.0.1"
	instanceName := "Instance 1"
	repositoryName := "repo"
	tokenID := "1234"
	assert.NoError(t, branch.(*model.Branch).Metadata.UpsertTemplateInstance(time.Now(), instanceID, instanceName, templateID, repositoryName, version, tokenID, nil))

	// Interaction
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		assert.NoError(t, console.ExpectString("Select branch:"))

		assert.NoError(t, console.SendEnter()) // enter - Main

		assert.NoError(t, console.ExpectString("Select template instance:"))

		assert.NoError(t, console.SendEnter()) // enter - tmpl1

		assert.NoError(t, console.ExpectEOF())
	}()

	// Run
	branchKey, instance, err := dialog.AskTemplateInstance(projectState)
	assert.NoError(t, err)
	assert.NoError(t, console.Tty().Close())
	wg.Wait()
	assert.NoError(t, console.Close())

	assert.Equal(t, model.BranchKey{ID: 123}, branchKey)
	assert.Equal(t, instanceID, instance.InstanceID)
}

func TestAskTemplateInstance_Noninteractive_InvalidInstance(t *testing.T) {
	t.Parallel()

	// Test dependencies
	dialog, o, _ := createDialogs(t, true)
	d := dependencies.NewMocked(t)
	projectState, err := d.MockedProject(fixtures.MinimalProjectFs(t)).LoadState(loadState.Options{LoadLocalState: true}, d)
	assert.NoError(t, err)
	branch, _ := projectState.LocalObjects().Get(model.BranchKey{ID: 123})

	instanceID := "inst1"
	templateID := "tmpl1"
	version := "1.0.1"
	instanceName := "Instance 1"
	repositoryName := "repo"
	tokenID := "1234"
	assert.NoError(t, branch.(*model.Branch).Metadata.UpsertTemplateInstance(time.Now(), instanceID, instanceName, templateID, repositoryName, version, tokenID, nil))

	o.Set("branch", 123)
	o.Set("instance", "inst2")
	_, _, err = dialog.AskTemplateInstance(projectState)
	assert.Error(t, err)
	assert.Equal(t, `template instance "inst2" was not found in branch "Main"`, err.Error())
}

func TestAskTemplateInstance_Noninteractive(t *testing.T) {
	t.Parallel()

	// Test dependencies
	dialog, o, _ := createDialogs(t, true)
	d := dependencies.NewMocked(t)
	projectState, err := d.MockedProject(fixtures.MinimalProjectFs(t)).LoadState(loadState.Options{LoadLocalState: true}, d)
	assert.NoError(t, err)
	branch, _ := projectState.LocalObjects().Get(model.BranchKey{ID: 123})

	instanceID := "inst1"
	templateID := "tmpl1"
	version := "1.0.1"
	instanceName := "Instance 1"
	repositoryName := "repo"
	tokenID := "1234"
	assert.NoError(t, branch.(*model.Branch).Metadata.UpsertTemplateInstance(time.Now(), instanceID, instanceName, templateID, repositoryName, version, tokenID, nil))

	o.Set("branch", 123)
	o.Set("instance", "inst1")
	_, _, err = dialog.AskTemplateInstance(projectState)
	assert.NoError(t, err)
}
