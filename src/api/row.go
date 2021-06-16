package api

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"keboola-as-code/src/client"
	"keboola-as-code/src/json"
	"keboola-as-code/src/model"
	"strconv"
)

func (a *StorageApi) GetConfigRow(branchId int, componentId string, configId string, rowId string) (*model.ConfigRow, error) {
	response := a.GetConfigRowRequest(branchId, componentId, configId, rowId).Send().Response()
	if response.HasResult() {
		return response.Result().(*model.ConfigRow), nil
	}
	return nil, response.Error()
}

func (a *StorageApi) CreateConfigRow(row *model.ConfigRow) (*model.ConfigRow, error) {
	request, err := a.CreateConfigRowRequest(row)
	if err != nil {
		return nil, err
	}

	response := request.Send().Response()
	if response.HasResult() {
		return response.Result().(*model.ConfigRow), nil
	}
	return nil, response.Error()
}

func (a *StorageApi) UpdateConfigRow(row *model.ConfigRow) (*model.ConfigRow, error) {
	request, err := a.UpdateConfigRowRequest(row)
	if err != nil {
		return nil, err
	}

	response := request.Send().Response()
	if response.HasResult() {
		return response.Result().(*model.ConfigRow), nil
	}
	return nil, response.Error()
}

// DeleteConfigRow - only config row in main branch can be deleted!
func (a *StorageApi) DeleteConfigRow(componentId string, configId string, rowId string) *client.Response {
	return a.DeleteConfigRowRequest(componentId, configId, rowId).Send().Response()
}

// GetConfigRowRequest https://keboola.docs.apiary.io/#reference/components-and-configurations/manage-configuration-rows/row-detail
func (a *StorageApi) GetConfigRowRequest(branchId int, componentId string, configId string, rowId string) *client.Request {
	return a.
		NewRequest(resty.MethodGet, fmt.Sprintf("branch/%d/components/%s/configs/%s/rows/%s", branchId, componentId, configId, rowId)).
		SetResult(&model.ConfigRow{
			BranchId:    branchId,
			ComponentId: componentId,
			ConfigId:    configId,
		})
}

// CreateConfigRowRequest https://keboola.docs.apiary.io/#reference/components-and-configurations/create-or-list-configuration-rows/create-development-branch-configuration-row
func (a *StorageApi) CreateConfigRowRequest(row *model.ConfigRow) (*client.Request, error) {
	// Id is autogenerated
	if row.Id != "" {
		panic("config id is set but it should be auto-generated")
	}

	// Encode config to JSON
	configJson, err := json.Encode(row.Config, false)
	if err != nil {
		panic(fmt.Errorf(`cannot JSON encode config row configuration: %s`, err))
	}

	// Create request
	request := a.
		NewRequest(resty.MethodPost, fmt.Sprintf("branch/%d/components/%s/configs/%s/rows", row.BranchId, row.ComponentId, row.ConfigId)).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetMultipartFormData(map[string]string{
			"name":              row.Name,
			"description":       row.Description,
			"changeDescription": row.ChangeDescription,
			"isDisabled":        strconv.FormatBool(row.IsDisabled),
			"configuration":     string(configJson),
		}).
		SetResult(row)

	return request, nil
}

// UpdateConfigRowRequest https://keboola.docs.apiary.io/#reference/components-and-configurations/manage-configuration-rows/update-row-for-development-branch
func (a *StorageApi) UpdateConfigRowRequest(row *model.ConfigRow) (*client.Request, error) {
	// Id is required
	if row.Id == "" {
		panic("config row id must be set")
	}

	// Encode config to JSON
	configJson, err := json.Encode(row.Config, false)
	if err != nil {
		panic(fmt.Errorf(`cannot JSON encode config row configuration: %s`, err))
	}

	// Create request
	request := a.
		NewRequest(resty.MethodPut, fmt.Sprintf("branch/%d/components/%s/configs/%s/rows/%s", row.BranchId, row.ComponentId, row.ConfigId, row.Id)).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetMultipartFormData(map[string]string{
			"name":              row.Name,
			"description":       row.Description,
			"changeDescription": row.ChangeDescription,
			"isDisabled":        strconv.FormatBool(row.IsDisabled),
			"configuration":     string(configJson),
		}).
		SetResult(&model.ConfigRow{})

	return request, nil
}

// DeleteConfigRowRequest https://keboola.docs.apiary.io/#reference/components-and-configurations/manage-configuration-rows/update-row
// Only config in main branch can be removed!
func (a *StorageApi) DeleteConfigRowRequest(componentId string, configId string, rowId string) *client.Request {
	return a.NewRequest(resty.MethodDelete, fmt.Sprintf("components/%s/configs/%s/rows/%s", componentId, configId, rowId))
}