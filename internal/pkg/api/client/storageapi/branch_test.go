package storageapi_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/keboola/keboola-as-code/internal/pkg/api/client/storageapi"
	"github.com/keboola/keboola-as-code/internal/pkg/env"
	"github.com/keboola/keboola-as-code/internal/pkg/http/client"
	"github.com/keboola/keboola-as-code/internal/pkg/json"
	"github.com/keboola/keboola-as-code/internal/pkg/model"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/testhelper"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/testproject"
)

func TestBranchApiCalls(t *testing.T) {
	t.Parallel()
	project := testproject.GetTestProject(t, env.Empty())
	project.SetState("empty.json")
	api := project.StorageApi()

	var job1 *model.Job
	var job2 *model.Job
	var job3 *model.Job
	var err error

	// Get default branch
	defaultBranch, err := api.GetDefaultBranch()
	assert.NoError(t, err)
	assert.NotNil(t, defaultBranch)
	assert.Equal(t, "Main", defaultBranch.Name)
	assert.True(t, defaultBranch.IsDefault)

	// Default branch cannot be created
	assert.PanicsWithError(t, "default branch cannot be created", func() {
		api.CreateBranch(&model.Branch{
			Name:        "Foo",
			Description: "Foo branch",
			IsDefault:   true,
		})
	})

	// Create branch, wait for successful job status
	branchFoo := &model.Branch{
		Name:        "Foo",
		Description: "Foo branch",
		IsDefault:   false,
	}
	job2, err = api.CreateBranch(branchFoo)
	assert.NoError(t, err)
	assert.NotNil(t, job2)
	assert.Equal(t, "success", job2.Status)
	assert.NotEmpty(t, branchFoo.Id)

	// Create branch with callback
	branchBar := &model.Branch{
		Name:        "Bar",
		Description: "Bar branch",
		IsDefault:   false,
	}
	onSuccessCalled := false
	request := api.CreateBranchRequest(branchBar).
		OnSuccess(func(response *client.Response) {
			// OnSuccess callback called when job is in successful state
			job := response.Result().(*model.Job)
			assert.NoError(t, response.Err())
			assert.NotNil(t, job)
			assert.Equal(t, "success", job.Status)
			assert.NotEmpty(t, branchBar.Id)
			onSuccessCalled = true
		}).
		Send()
	assert.NoError(t, request.Err())
	assert.True(t, request.IsSent())
	assert.True(t, request.IsDone())
	assert.True(t, onSuccessCalled)

	// Create branch, already exists
	branchBarDuplicate := &model.Branch{
		Name:        "Bar",
		Description: "Bar branch",
		IsDefault:   false,
	}
	job1, err = api.CreateBranch(branchBarDuplicate)
	assert.Nil(t, job1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "There already is a branch with name \"Bar\"")

	// Update branch
	branchFoo.Name = "Foo modified"
	branchFoo.Description = "Foo description modified"
	_, err = api.UpdateBranch(branchFoo, model.NewChangedFields("name", "description"))
	assert.NoError(t, err)

	// Update main branch description
	defaultBranch.Description = "Default branch"
	_, err = api.UpdateBranch(defaultBranch, model.NewChangedFields("description"))
	assert.NoError(t, err)

	// Cannot update default branch name
	defaultBranch.Name = "Not Allowed"
	assert.PanicsWithError(t, `changed field "name" not found in API values`, func() {
		api.UpdateBranch(defaultBranch, model.NewChangedFields("name", "description"))
	})

	// List branches
	var branches []*model.Branch
	branches, err = api.ListBranches()
	assert.NotNil(t, branches)
	assert.NoError(t, err)
	var encoded string
	testhelper.AssertWildcards(t, expectedBranchesAll(), json.MustEncodeString(branches, true), "Unexpected branches state")

	// Update metadata
	branchFoo.Metadata = map[string]string{"KBC.KaC.meta1": fmt.Sprintf("%d", rand.Intn(100))}
	assert.NoError(t, api.AppendBranchMetadataRequest(branchFoo).Send().Err())

	// List metadata
	req := api.ListBranchMetadataRequest(branchFoo.Id).Send()
	assert.NoError(t, req.Err())
	assert.NotNil(t, req.Result)
	branchMetadata := *req.Result.(*[]storageapi.Metadata)
	assert.Len(t, branchMetadata, 1)
	assert.Equal(t, "KBC.KaC.meta1", branchMetadata[0].Key)
	assert.Equal(t, branchFoo.Metadata["KBC.KaC.meta1"], branchMetadata[0].Value)

	// Delete metadata
	assert.NoError(t, api.DeleteBranchMetadataRequest(branchFoo.Id, branchMetadata[0].Id).Send().Err())

	// Check that metadata is deleted
	req = api.ListBranchMetadataRequest(branchFoo.Id).Send()
	assert.NoError(t, req.Err())
	assert.NotNil(t, req.Result)
	branchMetadata = *req.Result.(*[]storageapi.Metadata)
	assert.Len(t, branchMetadata, 0)

	// Delete branch
	job3, err = api.DeleteBranch(branchFoo.BranchKey)
	assert.NoError(t, err)
	assert.NotNil(t, job3)
	assert.Equal(t, "success", job3.Status)

	// Delete branch with callback
	onSuccessCalled = false
	request = api.DeleteBranchRequest(branchBar.BranchKey).
		OnSuccess(func(response *client.Response) {
			// OnSuccess callback called when job is in successful state
			job := response.Result().(*model.Job)
			assert.NoError(t, response.Err())
			assert.NotNil(t, job)
			assert.Equal(t, "success", job.Status)
			onSuccessCalled = true
		}).
		Send()
	assert.NoError(t, request.Err())
	assert.True(t, request.IsSent())
	assert.True(t, request.IsDone())
	assert.True(t, onSuccessCalled)

	// List branches
	branches, err = api.ListBranches()
	assert.NotNil(t, branches)
	assert.NoError(t, err)
	encoded, err = json.EncodeString(branches, true)
	assert.NoError(t, err)
	testhelper.AssertWildcards(t, expectedBranchesMain(), encoded, "Unexpected branches state")
}

func expectedBranchesAll() string {
	return `[
  {
    "id": %s,
    "name": "Foo modified",
    "description": "Foo description modified",
    "isDefault": false
  },
  {
    "id": %s,
    "name": "Bar",
    "description": "Bar branch",
    "isDefault": false
  },
  {
    "id": %s,
    "name": "Main",
    "description": "Default branch",
    "isDefault": true
  }
]`
}

func expectedBranchesMain() string {
	return `[
  {
    "id": %s,
    "name": "Main",
    "description": "Default branch",
    "isDefault": true
  }
]`
}