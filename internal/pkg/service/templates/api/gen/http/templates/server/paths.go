// Code generated by goa v3.14.6, DO NOT EDIT.
//
// HTTP request path constructors for the templates service.
//
// Command:
// $ goa gen github.com/keboola/keboola-as-code/api/templates --output
// ./internal/pkg/service/templates/api

package server

import (
	"fmt"
)

// APIRootIndexTemplatesPath returns the URL path to the templates service ApiRootIndex HTTP endpoint.
func APIRootIndexTemplatesPath() string {
	return "/"
}

// APIVersionIndexTemplatesPath returns the URL path to the templates service ApiVersionIndex HTTP endpoint.
func APIVersionIndexTemplatesPath() string {
	return "/v1"
}

// HealthCheckTemplatesPath returns the URL path to the templates service HealthCheck HTTP endpoint.
func HealthCheckTemplatesPath() string {
	return "/health-check"
}

// RepositoriesIndexTemplatesPath returns the URL path to the templates service RepositoriesIndex HTTP endpoint.
func RepositoriesIndexTemplatesPath() string {
	return "/v1/repositories"
}

// RepositoryIndexTemplatesPath returns the URL path to the templates service RepositoryIndex HTTP endpoint.
func RepositoryIndexTemplatesPath(repository string) string {
	return fmt.Sprintf("/v1/repositories/%v", repository)
}

// TemplatesIndexTemplatesPath returns the URL path to the templates service TemplatesIndex HTTP endpoint.
func TemplatesIndexTemplatesPath(repository string) string {
	return fmt.Sprintf("/v1/repositories/%v/templates", repository)
}

// TemplateIndexTemplatesPath returns the URL path to the templates service TemplateIndex HTTP endpoint.
func TemplateIndexTemplatesPath(repository string, template string) string {
	return fmt.Sprintf("/v1/repositories/%v/templates/%v", repository, template)
}

// VersionIndexTemplatesPath returns the URL path to the templates service VersionIndex HTTP endpoint.
func VersionIndexTemplatesPath(repository string, template string, version string) string {
	return fmt.Sprintf("/v1/repositories/%v/templates/%v/%v", repository, template, version)
}

// InputsIndexTemplatesPath returns the URL path to the templates service InputsIndex HTTP endpoint.
func InputsIndexTemplatesPath(repository string, template string, version string) string {
	return fmt.Sprintf("/v1/repositories/%v/templates/%v/%v/inputs", repository, template, version)
}

// ValidateInputsTemplatesPath returns the URL path to the templates service ValidateInputs HTTP endpoint.
func ValidateInputsTemplatesPath(repository string, template string, version string) string {
	return fmt.Sprintf("/v1/repositories/%v/templates/%v/%v/validate", repository, template, version)
}

// UseTemplateVersionTemplatesPath returns the URL path to the templates service UseTemplateVersion HTTP endpoint.
func UseTemplateVersionTemplatesPath(repository string, template string, version string) string {
	return fmt.Sprintf("/v1/repositories/%v/templates/%v/%v/use", repository, template, version)
}

// InstancesIndexTemplatesPath returns the URL path to the templates service InstancesIndex HTTP endpoint.
func InstancesIndexTemplatesPath(branch string) string {
	return fmt.Sprintf("/v1/project/%v/instances", branch)
}

// InstanceIndexTemplatesPath returns the URL path to the templates service InstanceIndex HTTP endpoint.
func InstanceIndexTemplatesPath(branch string, instanceID string) string {
	return fmt.Sprintf("/v1/project/%v/instances/%v", branch, instanceID)
}

// UpdateInstanceTemplatesPath returns the URL path to the templates service UpdateInstance HTTP endpoint.
func UpdateInstanceTemplatesPath(branch string, instanceID string) string {
	return fmt.Sprintf("/v1/project/%v/instances/%v", branch, instanceID)
}

// DeleteInstanceTemplatesPath returns the URL path to the templates service DeleteInstance HTTP endpoint.
func DeleteInstanceTemplatesPath(branch string, instanceID string) string {
	return fmt.Sprintf("/v1/project/%v/instances/%v", branch, instanceID)
}

// UpgradeInstanceTemplatesPath returns the URL path to the templates service UpgradeInstance HTTP endpoint.
func UpgradeInstanceTemplatesPath(branch string, instanceID string, version string) string {
	return fmt.Sprintf("/v1/project/%v/instances/%v/upgrade/%v", branch, instanceID, version)
}

// UpgradeInstanceInputsIndexTemplatesPath returns the URL path to the templates service UpgradeInstanceInputsIndex HTTP endpoint.
func UpgradeInstanceInputsIndexTemplatesPath(branch string, instanceID string, version string) string {
	return fmt.Sprintf("/v1/project/%v/instances/%v/upgrade/%v/inputs", branch, instanceID, version)
}

// UpgradeInstanceValidateInputsTemplatesPath returns the URL path to the templates service UpgradeInstanceValidateInputs HTTP endpoint.
func UpgradeInstanceValidateInputsTemplatesPath(branch string, instanceID string, version string) string {
	return fmt.Sprintf("/v1/project/%v/instances/%v/upgrade/%v/inputs", branch, instanceID, version)
}

// GetTaskTemplatesPath returns the URL path to the templates service GetTask HTTP endpoint.
func GetTaskTemplatesPath(taskID string) string {
	return fmt.Sprintf("/v1/tasks/%v", taskID)
}
