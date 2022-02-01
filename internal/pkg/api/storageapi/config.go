package storageapi

import (
	"fmt"

	"github.com/go-resty/resty/v2"

	"github.com/keboola/keboola-as-code/internal/pkg/http/client"
	"github.com/keboola/keboola-as-code/internal/pkg/model"
)

func (a *Api) ListComponents(branchId model.BranchId) ([]*model.ComponentWithConfigs, error) {
	response := a.ListComponentsRequest(branchId).Send().Response
	if response.HasResult() {
		return *response.Result().(*[]*model.ComponentWithConfigs), nil
	}
	return nil, response.Err()
}

func (a *Api) GetConfig(branchId model.BranchId, componentId model.ComponentId, configId model.ConfigId) (*model.Config, error) {
	response := a.GetConfigRequest(branchId, componentId, configId).Send().Response
	if response.HasResult() {
		return response.Result().(*model.Config), nil
	}
	return nil, response.Err()
}

func (a *Api) CreateConfig(config *model.ConfigWithRows) (*model.ConfigWithRows, error) {
	request, err := a.CreateConfigRequest(config)
	if err != nil {
		return nil, err
	}

	response := request.Send().Response
	if response.HasResult() {
		return response.Result().(*model.ConfigWithRows), nil
	}
	return nil, response.Err()
}

func (a *Api) UpdateConfig(config *model.Config, changed model.ChangedFields) (*model.Config, error) {
	request, err := a.UpdateConfigRequest(config, changed)
	if err != nil {
		return nil, err
	}

	response := request.Send().Response
	if response.HasResult() {
		return response.Result().(*model.Config), nil
	}
	return nil, response.Err()
}

func (a *Api) DeleteConfig(key model.ConfigKey) error {
	return a.DeleteConfigRequest(key).Send().Err()
}

func (a *Api) ListComponentsRequest(branchId model.BranchId) *client.Request {
	components := make([]*model.ComponentWithConfigs, 0)
	return a.
		NewRequest(resty.MethodGet, "branch/{branchId}/components").
		SetPathParam("branchId", branchId.String()).
		SetQueryParam("include", "configuration,rows").
		SetResult(&components).
		OnSuccess(func(response *client.Response) {
			if response.Result() != nil {
				// Add missing values
				for _, component := range components {
					component.BranchId = branchId

					// Cache
					a.Components().Set(component.Component)

					// Set config IDs
					for _, config := range component.Configs {
						config.BranchId = branchId
						config.ComponentId = component.Id
						config.SortRows()

						// Set rows IDs
						for _, row := range config.Rows {
							row.BranchId = branchId
							row.ComponentId = component.Id
							row.ConfigId = config.Id
						}
					}
				}
			}
		})
}

// ListConfigMetadataRequest https://keboola.docs.apiary.io/#reference/components-and-configurations/copy-configuration-rows/search-component-configurations
func (a *Api) ListConfigMetadataRequest(branchId model.BranchId) *client.Request {
	metadata := ConfigMetadataResponse{}
	return a.
		NewRequest(resty.MethodGet, "branch/{branchId}/search/component-configurations").
		SetPathParam("branchId", branchId.String()).
		SetQueryParam("include", "filteredMetadata").
		SetResult(metadata)
}

type (
	ConfigMetadataResponse     []ConfigMetadataResponseItem
	ConfigMetadataResponseItem struct {
		ComponentId model.ComponentId `json:"idComponent"`
		ConfigId    model.ConfigId    `json:"configurationId"`
		Metadata    []ConfigMetadata  `json:"metadata"`
	}
	ConfigMetadata struct {
		Id        string `json:"id"`
		Key       string `json:"key"`
		Value     string `json:"value"`
		Timestamp string `json:"timestamp"`
	}
)

func (c ConfigMetadataResponse) MetadataMap(branchId model.BranchId) map[model.ConfigKey][]ConfigMetadata {
	result := make(map[model.ConfigKey][]ConfigMetadata)
	for _, metadataWrapper := range c {
		configKey := model.ConfigKey{BranchId: branchId, ComponentId: metadataWrapper.ComponentId, Id: metadataWrapper.ConfigId}
		result[configKey] = metadataWrapper.Metadata
	}
	return result
}

// UpdateConfigMetadataRequest https://keboola.docs.apiary.io/#reference/metadata/components-configurations-metadata/create-or-update
func (a *Api) UpdateConfigMetadataRequest(config *model.Config) *client.Request {
	formBody := make(map[string]string)
	i := 0
	for k, v := range config.Metadata {
		formBody[fmt.Sprintf("metadata[%d][key]", i)] = k
		formBody[fmt.Sprintf("metadata[%d][value]", i)] = v
		i++
	}
	return a.
		NewRequest(resty.MethodPost, "branch/{branchId}/components/{componentId}/configs/{configId}/metadata").
		SetPathParam("branchId", config.BranchId.String()).
		SetPathParam("componentId", config.ComponentId.String()).
		SetPathParam("configId", config.Id.String()).
		SetFormBody(formBody)
}

// GetConfigRequest https://keboola.docs.apiary.io/#reference/components-and-configurations/manage-configurations/development-branch-configuration-detail
func (a *Api) GetConfigRequest(branchId model.BranchId, componentId model.ComponentId, configId model.ConfigId) *client.Request {
	config := &model.Config{}
	config.BranchId = branchId
	config.ComponentId = componentId
	return a.
		NewRequest(resty.MethodGet, "branch/{branchId}/components/{componentId}/configs/{configId}").
		SetPathParam("branchId", branchId.String()).
		SetPathParam("componentId", componentId.String()).
		SetPathParam("configId", configId.String()).
		SetResult(config)
}

// CreateConfigRequest https://keboola.docs.apiary.io/#reference/components-and-configurations/component-configurations/create-development-branch-configuration
func (a *Api) CreateConfigRequest(config *model.ConfigWithRows) (*client.Request, error) {
	// Data
	values, err := config.ToApiValues()
	if err != nil {
		return nil, err
	}

	// Create config with the defined ID
	if config.Id != "" {
		values["configurationId"] = config.Id.String()
	}

	// Create config
	var configRequest *client.Request
	configRequest = a.
		NewRequest(resty.MethodPost, "branch/{branchId}/components/{componentId}/configs").
		SetPathParam("branchId", config.BranchId.String()).
		SetPathParam("componentId", config.ComponentId.String()).
		SetFormBody(values).
		SetResult(config).
		// Create config rows
		OnSuccess(func(response *client.Response) {
			for _, row := range config.Rows {
				row.BranchId = config.BranchId
				row.ComponentId = config.ComponentId
				row.ConfigId = config.Id
				rowRequest, err := a.CreateConfigRowRequest(row)
				if err != nil {
					response.SetErr(err)
				}
				configRequest.WaitFor(rowRequest)
				response.Sender().Request(rowRequest).Send()
			}
		})

	return configRequest, nil
}

// UpdateConfigRequest https://keboola.docs.apiary.io/#reference/components-and-configurations/manage-configurations/update-development-branch-configuration
func (a *Api) UpdateConfigRequest(config *model.Config, changed model.ChangedFields) (*client.Request, error) {
	// Id is required
	if config.Id == "" {
		panic("config id must be set")
	}

	// Data
	values, err := config.ToApiValues()
	if err != nil {
		return nil, err
	}

	// Update config
	request := a.
		NewRequest(resty.MethodPut, "branch/{branchId}/components/{componentId}/configs/{configId}").
		SetPathParam("branchId", config.BranchId.String()).
		SetPathParam("componentId", config.ComponentId.String()).
		SetPathParam("configId", config.Id.String()).
		SetFormBody(getChangedValues(values, changed)).
		SetResult(config)

	return request, nil
}

// DeleteConfigRequest https://keboola.docs.apiary.io/#reference/components-and-configurations/manage-configurations/delete-configuration
func (a *Api) DeleteConfigRequest(key model.ConfigKey) *client.Request {
	return a.NewRequest(resty.MethodDelete, "branch/{branchId}/components/{componentId}/configs/{configId}").
		SetPathParam("branchId", key.BranchId.String()).
		SetPathParam("componentId", key.ComponentId.String()).
		SetPathParam("configId", key.Id.String())
}

func (a *Api) DeleteConfigsInBranchRequest(key model.BranchKey) *client.Request {
	return a.ListComponentsRequest(key.Id).
		OnSuccess(func(response *client.Response) {
			for _, component := range *response.Result().(*[]*model.ComponentWithConfigs) {
				for _, config := range component.Configs {
					response.
						Sender().
						Request(a.DeleteConfigRequest(config.ConfigKey)).
						Send()
				}
			}
		})
}