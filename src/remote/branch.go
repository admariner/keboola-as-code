package remote

import (
	"fmt"

	"keboola-as-code/src/client"
	"keboola-as-code/src/model"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
)

func (a *StorageApi) GetDefaultBranch() (*model.Branch, error) {
	branches, err := a.ListBranches()
	if err != nil {
		return nil, err
	}

	for _, branch := range branches {
		if branch.IsDefault {
			return branch, nil
		}
	}

	return nil, fmt.Errorf("default branch not found")
}

func (a *StorageApi) GetBranch(branchId int) (*model.Branch, error) {
	response := a.GetBranchRequest(branchId).Send().Response
	if response.HasResult() {
		return response.Result().(*model.Branch), nil
	}
	return nil, response.Err()
}

func (a *StorageApi) CreateBranch(branch *model.Branch) (*model.Job, error) {
	response := a.CreateBranchRequest(branch).Send().Response
	if response.HasResult() {
		return response.Result().(*model.Job), nil
	}
	return nil, response.Err()
}

func (a *StorageApi) UpdateBranch(branch *model.Branch, changed []string) (*model.Branch, error) {
	response := a.UpdateBranchRequest(branch, changed).Send().Response
	if response.HasResult() {
		return response.Result().(*model.Branch), nil
	}
	return nil, response.Err()
}

func (a *StorageApi) ListBranches() ([]*model.Branch, error) {
	response := a.ListBranchesRequest().Send().Response
	if response.HasResult() {
		return *response.Result().(*[]*model.Branch), nil
	}
	return nil, response.Err()
}

func (a *StorageApi) DeleteBranch(branchId int) (*model.Job, error) {
	response := a.DeleteBranchRequest(branchId).Send().Response
	if response.HasResult() {
		return response.Result().(*model.Job), nil
	}
	return nil, response.Err()
}

// GetBranchRequest https://keboola.docs.apiary.io/#reference/development-branches/branch-manipulation/branch-detail
func (a *StorageApi) GetBranchRequest(branchId int) *client.Request {
	branch := &model.Branch{}
	return a.
		NewRequest(resty.MethodGet, "dev-branches/{branchId}").
		SetPathParam("branchId", cast.ToString(branchId)).
		SetResult(branch)
}

// CreateBranchRequest https://keboola.docs.apiary.io/#reference/development-branches/branches/create-branch
func (a *StorageApi) CreateBranchRequest(branch *model.Branch) *client.Request {
	job := &model.Job{}
	// Id is autogenerated
	if branch.Id != 0 {
		panic(fmt.Errorf("branch id is set but it should be auto-generated"))
	}

	// Default branch cannot be created
	if branch.IsDefault {
		panic(fmt.Errorf("default branch cannot be created"))
	}

	// Create request
	request := a.
		NewRequest(resty.MethodPost, "dev-branches").
		SetBody(map[string]string{
			"name":        branch.Name,
			"description": branch.Description,
		}).
		SetResult(job)

	request.OnSuccess(waitForJob(a, request, job, func(response *client.Response) {
		// Set branch id from the job results
		branch.Id = cast.ToInt(job.Results["id"])
	}))

	return request
}

// UpdateBranchRequest https://keboola.docs.apiary.io/#reference/development-branches/branches/update-branch
func (a *StorageApi) UpdateBranchRequest(branch *model.Branch, changed []string) *client.Request {
	// Id is required
	if branch.Id == 0 {
		panic("branch id must be set")
	}

	// Data
	all := map[string]string{
		"description": branch.Description,
	}

	// Name of the default branch cannot be changed
	if !branch.IsDefault {
		all["name"] = branch.Name
	}

	// Create request
	request := a.
		NewRequest(resty.MethodPut, "dev-branches/{branchId}").
		SetPathParam("branchId", cast.ToString(branch.Id)).
		SetBody(getChangedValues(all, changed)).
		SetResult(branch)

	return request
}

// ListBranchesRequest https://keboola.docs.apiary.io/#reference/development-branches/branches/list-branches
func (a *StorageApi) ListBranchesRequest() *client.Request {
	branches := make([]*model.Branch, 0)
	return a.
		NewRequest(resty.MethodGet, "dev-branches").
		SetResult(&branches)

}

// DeleteBranchRequest https://keboola.docs.apiary.io/#reference/development-branches/branch-manipulation/delete-branch
func (a *StorageApi) DeleteBranchRequest(branchId int) *client.Request {
	job := &model.Job{}
	request := a.
		NewRequest(resty.MethodDelete, "dev-branches/{branchId}").
		SetPathParam("branchId", cast.ToString(branchId)).
		SetResult(job)
	request.OnSuccess(waitForJob(a, request, job, nil))
	return request
}
