// Code generated by goa v3.11.1, DO NOT EDIT.
//
// templates service
//
// Command:
// $ goa gen github.com/keboola/keboola-as-code/api/templates --output
// ./internal/pkg/service/templates/api

package templates

import (
	"context"

	"github.com/keboola/keboola-as-code/internal/pkg/service/common/task"
	dependencies "github.com/keboola/keboola-as-code/internal/pkg/service/templates/api/dependencies"
	"goa.design/goa/v3/security"
)

// Service for applying templates to Keboola projects.
type Service interface {
	// Redirect to /v1.
	APIRootIndex(dependencies.ForPublicRequest) (err error)
	// List API name and link to documentation.
	APIVersionIndex(dependencies.ForPublicRequest) (res *ServiceDetail, err error)
	// HealthCheck implements HealthCheck.
	HealthCheck(dependencies.ForPublicRequest) (res string, err error)
	// List all template repositories defined in the project.
	RepositoriesIndex(dependencies.ForProjectRequest, *RepositoriesIndexPayload) (res *Repositories, err error)
	// Get details of specified repository. Use "keboola" for default Keboola
	// repository.
	RepositoryIndex(dependencies.ForProjectRequest, *RepositoryIndexPayload) (res *Repository, err error)
	// List all templates  defined in the repository.
	TemplatesIndex(dependencies.ForProjectRequest, *TemplatesIndexPayload) (res *Templates, err error)
	// Get detail and versions of specified template.
	TemplateIndex(dependencies.ForProjectRequest, *TemplateIndexPayload) (res *TemplateDetail, err error)
	// Get details of specified template version.
	VersionIndex(dependencies.ForProjectRequest, *VersionIndexPayload) (res *VersionDetailExtended, err error)
	// Get inputs for the "use" API call.
	InputsIndex(dependencies.ForProjectRequest, *InputsIndexPayload) (res *Inputs, err error)
	// Validate inputs for the "use" API call.
	// Only configured steps should be send.
	ValidateInputs(dependencies.ForProjectRequest, *ValidateInputsPayload) (res *ValidationResult, err error)
	// Validate inputs and use template in the branch.
	// Only configured steps should be send.
	UseTemplateVersion(dependencies.ForProjectRequest, *UseTemplateVersionPayload) (res *Task, err error)
	// InstancesIndex implements InstancesIndex.
	InstancesIndex(dependencies.ForProjectRequest, *InstancesIndexPayload) (res *Instances, err error)
	// InstanceIndex implements InstanceIndex.
	InstanceIndex(dependencies.ForProjectRequest, *InstanceIndexPayload) (res *InstanceDetail, err error)
	// UpdateInstance implements UpdateInstance.
	UpdateInstance(dependencies.ForProjectRequest, *UpdateInstancePayload) (res *InstanceDetail, err error)
	// DeleteInstance implements DeleteInstance.
	DeleteInstance(dependencies.ForProjectRequest, *DeleteInstancePayload) (err error)
	// UpgradeInstance implements UpgradeInstance.
	UpgradeInstance(dependencies.ForProjectRequest, *UpgradeInstancePayload) (res *Task, err error)
	// UpgradeInstanceInputsIndex implements UpgradeInstanceInputsIndex.
	UpgradeInstanceInputsIndex(dependencies.ForProjectRequest, *UpgradeInstanceInputsIndexPayload) (res *Inputs, err error)
	// UpgradeInstanceValidateInputs implements UpgradeInstanceValidateInputs.
	UpgradeInstanceValidateInputs(dependencies.ForProjectRequest, *UpgradeInstanceValidateInputsPayload) (res *ValidationResult, err error)
	// Get details of a task.
	GetTask(dependencies.ForProjectRequest, *GetTaskPayload) (res *Task, err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// APIKeyAuth implements the authorization logic for the APIKey security scheme.
	APIKeyAuth(ctx context.Context, key string, schema *security.APIKeyScheme) (context.Context, error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "templates"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [19]string{"ApiRootIndex", "ApiVersionIndex", "HealthCheck", "RepositoriesIndex", "RepositoryIndex", "TemplatesIndex", "TemplateIndex", "VersionIndex", "InputsIndex", "ValidateInputs", "UseTemplateVersion", "InstancesIndex", "InstanceIndex", "UpdateInstance", "DeleteInstance", "UpgradeInstance", "UpgradeInstanceInputsIndex", "UpgradeInstanceValidateInputs", "GetTask"}

// Author of template or repository.
type Author struct {
	// Name of the author.
	Name string
	// Link to the author website.
	URL string
}

// Date of change and who made it.
type ChangeInfo struct {
	// Date and time of the change.
	Date string
	// The token by which the change was made.
	TokenID string
}

// The configuration that is part of the template instance.
type Config struct {
	// Component ID.
	ComponentID string
	// Configuration ID.
	ConfigID string
	// Name of the configuration.
	Name string
}

// DeleteInstancePayload is the payload type of the templates service
// DeleteInstance method.
type DeleteInstancePayload struct {
	StorageAPIToken string
	// ID of the template instance.
	InstanceID string
	// ID of the branch. Use "default" for default branch.
	Branch string
}

// Generic error
type GenericError struct {
	// HTTP status code.
	StatusCode int
	// Name of error.
	Name string
	// Error message.
	Message string
}

// GetTaskPayload is the payload type of the templates service GetTask method.
type GetTaskPayload struct {
	StorageAPIToken string
	TaskID          TaskID
}

// User input.
type Input struct {
	// Unique ID of the input.
	ID string
	// Name of the input.
	Name string
	// Description of the input.
	Description string
	// Type of the input.
	Type string
	// Kind of the input.
	Kind string
	// Default value, match defined type.
	Default interface{}
	// Input options for type = select OR multiselect.
	Options []*InputOption
	// Component id for "oauth" kind inputs.
	ComponentID *string
	// OAuth input id for "oauthAccounts" kind inputs.
	OauthInputID *string
}

// Input option for type = select OR multiselect.
type InputOption struct {
	// Visible label of the option.
	Label string
	// Value of the option.
	Value string
}

// Validation Detail of the input.
type InputValidationResult struct {
	// Input ID.
	ID string
	// If false, the input should be hidden to user.
	Visible bool
	// Error message.
	Error *string
}

// Input value filled in by user.
type InputValue struct {
	// Unique ID of the input.
	ID string
	// Input value filled in by user in the required type.
	Value interface{}
}

// Inputs is the result type of the templates service InputsIndex method.
type Inputs struct {
	// List of the step groups.
	StepGroups []*StepGroup
	// Initial state - same structure as the validation result.
	InitialState *ValidationResult
}

// InputsIndexPayload is the payload type of the templates service InputsIndex
// method.
type InputsIndexPayload struct {
	StorageAPIToken string
	// Semantic version of the template. Use "default" for default version.
	Version string
	// ID of the template.
	Template string
	// Name of the template repository. Use "keboola" for default Keboola
	// repository.
	Repository string
}

// ID of the template.
type Instance struct {
	// ID of the template.
	TemplateID string
	// ID of the template instance.
	InstanceID string
	// ID of the branch.
	Branch string
	// Name of the template repository.
	RepositoryName string
	// Semantic version of the template.
	Version string
	// Name of the instance.
	Name string
	// Instance creation date and token.
	Created *ChangeInfo
	// Instance update date and token.
	Updated    *ChangeInfo
	MainConfig *MainConfig
}

// InstanceDetail is the result type of the templates service InstanceIndex
// method.
type InstanceDetail struct {
	// Information about the template version. Can be null if the repository or
	// template no longer exists. If the exact version is not found, the nearest
	// one is used.
	VersionDetail *VersionDetail
	// All configurations from the instance.
	Configurations []*Config
	// ID of the template.
	TemplateID string
	// ID of the template instance.
	InstanceID string
	// ID of the branch.
	Branch string
	// Name of the template repository.
	RepositoryName string
	// Semantic version of the template.
	Version string
	// Name of the instance.
	Name string
	// Instance creation date and token.
	Created *ChangeInfo
	// Instance update date and token.
	Updated    *ChangeInfo
	MainConfig *MainConfig
}

// InstanceIndexPayload is the payload type of the templates service
// InstanceIndex method.
type InstanceIndexPayload struct {
	StorageAPIToken string
	// ID of the template instance.
	InstanceID string
	// ID of the branch. Use "default" for default branch.
	Branch string
}

// Instances is the result type of the templates service InstancesIndex method.
type Instances struct {
	// All instances found in branch.
	Instances []*Instance
}

// InstancesIndexPayload is the payload type of the templates service
// InstancesIndex method.
type InstancesIndexPayload struct {
	StorageAPIToken string
	// ID of the branch. Use "default" for default branch.
	Branch string
}

// Main config of the instance, usually an orchestration. Optional.
type MainConfig struct {
	// Component ID.
	ComponentID string
	// Configuration ID.
	ConfigID string
}

// Project locked error
type ProjectLockedError struct {
	// HTTP status code.
	StatusCode int
	// Name of error.
	Name string
	// Error message.
	Message string
	// Indicates how long the user agent should wait before making a follow-up
	// request.
	RetryAfter string
}

// Repositories is the result type of the templates service RepositoriesIndex
// method.
type Repositories struct {
	// All template repositories defined in the project.
	Repositories []*Repository
}

// RepositoriesIndexPayload is the payload type of the templates service
// RepositoriesIndex method.
type RepositoriesIndexPayload struct {
	StorageAPIToken string
}

// Repository is the result type of the templates service RepositoryIndex
// method.
type Repository struct {
	// Template repository name. Use "keboola" for default Keboola repository.
	Name string
	// Git URL to the repository.
	URL string
	// Git branch or tag.
	Ref    string
	Author *Author
}

// RepositoryIndexPayload is the payload type of the templates service
// RepositoryIndex method.
type RepositoryIndexPayload struct {
	StorageAPIToken string
	// Name of the template repository. Use "keboola" for default Keboola
	// repository.
	Repository string
}

// ServiceDetail is the result type of the templates service ApiVersionIndex
// method.
type ServiceDetail struct {
	// Name of the API
	API string
	// URL of the API documentation.
	Documentation string
}

// Step is a container for inputs.
type Step struct {
	// Unique ID of the step.
	ID string
	// Icon for UI. Component icon if it starts with "component:...", or a common
	// icon if it starts with "common:...".
	Icon string
	// Name of the step.
	Name string
	// Description of the step.
	Description string
	// Name of the dialog with inputs.
	DialogName string
	// Description of the dialog with inputs.
	DialogDescription string
	// Inputs in the step.
	Inputs []*Input
}

// Step group is a container for steps.
type StepGroup struct {
	// Unique ID of the step group.
	ID string
	// Description of the step group, a tooltip explaining what needs to be
	// configured.
	Description string
	// The number of steps that must be configured.
	Required string
	// Steps in the group.
	Steps []*Step
}

// Validation Detail of the step group.
type StepGroupValidationResult struct {
	// Step group ID.
	ID string
	// True if the required number of steps is configured and all inputs are valid.
	Valid bool
	// Are all inputs valid?
	Error *string
	// List of Details for the steps.
	Steps []*StepValidationResult
}

// Step with input values filled in by user.
type StepPayload struct {
	// Unique ID of the step.
	ID string
	// Input values.
	Inputs []*InputValue
}

// Validation Detail of the step.
type StepValidationResult struct {
	// Step ID.
	ID string
	// True if the step was part of the sent payload.
	Configured bool
	// True if all inputs in the step are valid.
	Valid bool
	// List of Details for the inputs.
	Inputs []*InputValidationResult
}

// Task is the result type of the templates service UseTemplateVersion method.
type Task struct {
	ID TaskID
	// URL of the task.
	URL string
	// Task status, one of: processing, success, error
	Status string
	// Shortcut for status != "processing".
	IsFinished bool
	// Date and time of the task creation.
	CreatedAt string
	// Date and time of the task end.
	FinishedAt *string
	// Duration of the task in milliseconds.
	Duration *int64
	Result   *string
	Error    *string
}

// Unique ID of the task.
type TaskID = task.ID

// Template.
type Template struct {
	// Template ID.
	ID string
	// Template name.
	Name string
	// List of categories the template belongs to.
	Categories []string
	// List of components used in the template.
	Components []string
	Author     *Author
	// Short description of the template.
	Description string
	// Recommended version of the template.
	DefaultVersion string
	// All available versions of the template.
	Versions []*Version
}

// TemplateDetail is the result type of the templates service TemplateIndex
// method.
type TemplateDetail struct {
	// Information about the repository.
	Repository *Repository
	// Template ID.
	ID string
	// Template name.
	Name string
	// List of categories the template belongs to.
	Categories []string
	// List of components used in the template.
	Components []string
	Author     *Author
	// Short description of the template.
	Description string
	// Recommended version of the template.
	DefaultVersion string
	// All available versions of the template.
	Versions []*Version
}

// TemplateIndexPayload is the payload type of the templates service
// TemplateIndex method.
type TemplateIndexPayload struct {
	StorageAPIToken string
	// ID of the template.
	Template string
	// Name of the template repository. Use "keboola" for default Keboola
	// repository.
	Repository string
}

// Templates is the result type of the templates service TemplatesIndex method.
type Templates struct {
	// Information about the repository.
	Repository *Repository
	// All template defined in the repository.
	Templates []*Template
}

// TemplatesIndexPayload is the payload type of the templates service
// TemplatesIndex method.
type TemplatesIndexPayload struct {
	StorageAPIToken string
	// Name of the template repository. Use "keboola" for default Keboola
	// repository.
	Repository string
}

// UpdateInstancePayload is the payload type of the templates service
// UpdateInstance method.
type UpdateInstancePayload struct {
	StorageAPIToken string
	// New name of the instance.
	Name string
	// ID of the template instance.
	InstanceID string
	// ID of the branch. Use "default" for default branch.
	Branch string
}

// UpgradeInstanceInputsIndexPayload is the payload type of the templates
// service UpgradeInstanceInputsIndex method.
type UpgradeInstanceInputsIndexPayload struct {
	StorageAPIToken string
	// Semantic version of the template. Use "default" for default version.
	Version string
	// ID of the template instance.
	InstanceID string
	// ID of the branch. Use "default" for default branch.
	Branch string
}

// UpgradeInstancePayload is the payload type of the templates service
// UpgradeInstance method.
type UpgradeInstancePayload struct {
	StorageAPIToken string
	// Semantic version of the template. Use "default" for default version.
	Version string
	// ID of the template instance.
	InstanceID string
	// ID of the branch. Use "default" for default branch.
	Branch string
	// Steps with input values filled in by user.
	Steps []*StepPayload
}

// UpgradeInstanceValidateInputsPayload is the payload type of the templates
// service UpgradeInstanceValidateInputs method.
type UpgradeInstanceValidateInputsPayload struct {
	StorageAPIToken string
	// Semantic version of the template. Use "default" for default version.
	Version string
	// ID of the template instance.
	InstanceID string
	// ID of the branch. Use "default" for default branch.
	Branch string
	// Steps with input values filled in by user.
	Steps []*StepPayload
}

// UseTemplateVersionPayload is the payload type of the templates service
// UseTemplateVersion method.
type UseTemplateVersionPayload struct {
	StorageAPIToken string
	// Name of the new template instance.
	Name string
	// ID of the branch. Use "default" for default branch.
	Branch string
	// Semantic version of the template. Use "default" for default version.
	Version string
	// ID of the template.
	Template string
	// Name of the template repository. Use "keboola" for default Keboola
	// repository.
	Repository string
	// Steps with input values filled in by user.
	Steps []*StepPayload
}

// ValidateInputsPayload is the payload type of the templates service
// ValidateInputs method.
type ValidateInputsPayload struct {
	StorageAPIToken string
	// Semantic version of the template. Use "default" for default version.
	Version string
	// ID of the template.
	Template string
	// Name of the template repository. Use "keboola" for default Keboola
	// repository.
	Repository string
	// Steps with input values filled in by user.
	Steps []*StepPayload
}

type ValidationError struct {
	// Name of error.
	Name string
	// Error message.
	Message          string
	ValidationResult *ValidationResult
}

// ValidationResult is the result type of the templates service ValidateInputs
// method.
type ValidationResult struct {
	// True if all groups and inputs are valid.
	Valid bool
	// List of Details for the step groups.
	StepGroups []*StepGroupValidationResult
}

// Template version.
type Version struct {
	// Semantic version.
	Version string
	// If true, then the version is ready for production use.
	Stable bool
	// Optional short description of the version. Can be empty.
	Description string
}

type VersionDetail struct {
	// List of components used in the template version.
	Components []string
	// Extended description of the template in Markdown format.
	LongDescription string
	// Readme of the template version in Markdown format.
	Readme string
	// Semantic version.
	Version string
	// If true, then the version is ready for production use.
	Stable bool
	// Optional short description of the version. Can be empty.
	Description string
}

// VersionDetailExtended is the result type of the templates service
// VersionIndex method.
type VersionDetailExtended struct {
	// Information about the repository.
	Repository *Repository
	// Information about the template.
	Template *Template
	// List of components used in the template version.
	Components []string
	// Extended description of the template in Markdown format.
	LongDescription string
	// Readme of the template version in Markdown format.
	Readme string
	// Semantic version.
	Version string
	// If true, then the version is ready for production use.
	Stable bool
	// Optional short description of the version. Can be empty.
	Description string
}

// VersionIndexPayload is the payload type of the templates service
// VersionIndex method.
type VersionIndexPayload struct {
	StorageAPIToken string
	// Semantic version of the template. Use "default" for default version.
	Version string
	// ID of the template.
	Template string
	// Name of the template repository. Use "keboola" for default Keboola
	// repository.
	Repository string
}

// Error returns an error description.
func (e *GenericError) Error() string {
	return "Generic error"
}

// ErrorName returns "GenericError".
//
// Deprecated: Use GoaErrorName - https://github.com/goadesign/goa/issues/3105
func (e *GenericError) ErrorName() string {
	return e.GoaErrorName()
}

// GoaErrorName returns "GenericError".
func (e *GenericError) GoaErrorName() string {
	return e.Name
}

// Error returns an error description.
func (e *ProjectLockedError) Error() string {
	return "Project locked error"
}

// ErrorName returns "ProjectLockedError".
//
// Deprecated: Use GoaErrorName - https://github.com/goadesign/goa/issues/3105
func (e *ProjectLockedError) ErrorName() string {
	return e.GoaErrorName()
}

// GoaErrorName returns "ProjectLockedError".
func (e *ProjectLockedError) GoaErrorName() string {
	return e.Name
}

// Error returns an error description.
func (e *ValidationError) Error() string {
	return ""
}

// ErrorName returns "ValidationError".
//
// Deprecated: Use GoaErrorName - https://github.com/goadesign/goa/issues/3105
func (e *ValidationError) ErrorName() string {
	return e.GoaErrorName()
}

// GoaErrorName returns "ValidationError".
func (e *ValidationError) GoaErrorName() string {
	return e.Name
}
