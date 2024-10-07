// Code generated by goa v3.19.1, DO NOT EDIT.
//
// apps-proxy HTTP server types
//
// Command:
// $ goa gen github.com/keboola/keboola-as-code/api/appsproxy --output
// ./internal/pkg/service/appsproxy/api

package server

import (
	appsproxy "github.com/keboola/keboola-as-code/internal/pkg/service/appsproxy/api/gen/apps_proxy"
)

// APIVersionIndexResponseBody is the type of the "apps-proxy" service
// "ApiVersionIndex" endpoint HTTP response body.
type APIVersionIndexResponseBody struct {
	// Name of the API
	API string `form:"api" json:"api" xml:"api"`
	// URL of the API documentation.
	Documentation string `form:"documentation" json:"documentation" xml:"documentation"`
}

// ValidateResponseBody is the type of the "apps-proxy" service "Validate"
// endpoint HTTP response body.
type ValidateResponseBody struct {
	// All authorization providers.
	Configuration []*ConfigurationResponseBody `form:"configuration,omitempty" json:"configuration,omitempty" xml:"configuration,omitempty"`
}

// ConfigurationResponseBody is used to define fields on response body types.
type ConfigurationResponseBody struct {
	// Unique ID of provider.
	ID string `form:"id" json:"id" xml:"id"`
	// Client ID of provider.
	ClientID string `form:"clientID" json:"clientID" xml:"clientID"`
	// Client secret provided by OIDC provider.
	ClientSecret string `form:"clientSecret" json:"clientSecret" xml:"clientSecret"`
}

// NewAPIVersionIndexResponseBody builds the HTTP response body from the result
// of the "ApiVersionIndex" endpoint of the "apps-proxy" service.
func NewAPIVersionIndexResponseBody(res *appsproxy.ServiceDetail) *APIVersionIndexResponseBody {
	body := &APIVersionIndexResponseBody{
		API:           res.API,
		Documentation: res.Documentation,
	}
	return body
}

// NewValidateResponseBody builds the HTTP response body from the result of the
// "Validate" endpoint of the "apps-proxy" service.
func NewValidateResponseBody(res *appsproxy.Validations) *ValidateResponseBody {
	body := &ValidateResponseBody{}
	if res.Configuration != nil {
		body.Configuration = make([]*ConfigurationResponseBody, len(res.Configuration))
		for i, val := range res.Configuration {
			body.Configuration[i] = marshalAppsproxyConfigurationToConfigurationResponseBody(val)
		}
	}
	return body
}

// NewValidatePayload builds a apps-proxy service Validate endpoint payload.
func NewValidatePayload(storageAPIToken string) *appsproxy.ValidatePayload {
	v := &appsproxy.ValidatePayload{}
	v.StorageAPIToken = storageAPIToken

	return v
}
