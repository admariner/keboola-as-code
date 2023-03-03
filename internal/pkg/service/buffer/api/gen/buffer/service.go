// Code generated by goa v3.11.1, DO NOT EDIT.
//
// buffer service
//
// Command:
// $ goa gen github.com/keboola/keboola-as-code/api/buffer --output
// ./internal/pkg/service/buffer/api

package buffer

import (
	"context"
	"io"

	dependencies "github.com/keboola/keboola-as-code/internal/pkg/service/buffer/api/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/buffer/store/key"
	"github.com/keboola/keboola-as-code/internal/pkg/service/buffer/store/model/column"
	"goa.design/goa/v3/security"
)

// A service for continuously importing data to Keboola storage.
type Service interface {
	// Redirect to /v1.
	APIRootIndex(dependencies.ForPublicRequest) (err error)
	// List API name and link to documentation.
	APIVersionIndex(dependencies.ForPublicRequest) (res *ServiceDetail, err error)
	// HealthCheck implements HealthCheck.
	HealthCheck(dependencies.ForPublicRequest) (res string, err error)
	// Create a new receiver for a given project
	CreateReceiver(dependencies.ForProjectRequest, *CreateReceiverPayload) (res *Task, err error)
	// Update a receiver export.
	UpdateReceiver(dependencies.ForProjectRequest, *UpdateReceiverPayload) (res *Receiver, err error)
	// List all receivers for a given project.
	ListReceivers(dependencies.ForProjectRequest, *ListReceiversPayload) (res *ReceiversList, err error)
	// Get the configuration of a receiver.
	GetReceiver(dependencies.ForProjectRequest, *GetReceiverPayload) (res *Receiver, err error)
	// Delete a receiver.
	DeleteReceiver(dependencies.ForProjectRequest, *DeleteReceiverPayload) (err error)
	// Each export uses its own token scoped to the target bucket, this endpoint
	// refreshes all of those tokens.
	RefreshReceiverTokens(dependencies.ForProjectRequest, *RefreshReceiverTokensPayload) (res *Receiver, err error)
	// Create a new export for an existing receiver.
	CreateExport(dependencies.ForProjectRequest, *CreateExportPayload) (res *Task, err error)
	// Get the configuration of an export.
	GetExport(dependencies.ForProjectRequest, *GetExportPayload) (res *Export, err error)
	// List all exports for a given receiver.
	ListExports(dependencies.ForProjectRequest, *ListExportsPayload) (res *ExportsList, err error)
	// Update a receiver export.
	UpdateExport(dependencies.ForProjectRequest, *UpdateExportPayload) (res *Task, err error)
	// Delete a receiver export.
	DeleteExport(dependencies.ForProjectRequest, *DeleteExportPayload) (err error)
	// Upload data into the receiver.
	Import(dependencies.ForPublicRequest, *ImportPayload, io.ReadCloser) (err error)
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
const ServiceName = "buffer"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [16]string{"ApiRootIndex", "ApiVersionIndex", "HealthCheck", "CreateReceiver", "UpdateReceiver", "ListReceivers", "GetReceiver", "DeleteReceiver", "RefreshReceiverTokens", "CreateExport", "GetExport", "ListExports", "UpdateExport", "DeleteExport", "Import", "GetTask"}

// An output mapping defined by a template.
type Column struct {
	// Sets this column as a part of the primary key of the destination table.
	PrimaryKey bool
	// Column mapping type. This represents a static mapping (e.g. `body` or
	// `headers`), or a custom mapping using a template language (`template`).
	Type column.Type
	// Column name.
	Name string
	// Template mapping details.
	Template *Template
}

// Table import triggers.
type Conditions struct {
	// Maximum import buffer size in number of records.
	Count int
	// Maximum import buffer size in bytes. Units: B, KB, MB.
	Size string
	// Minimum import interval. Units: [s]econd,[m]inute,[h]our.
	Time string
}

type CreateExportData struct {
	// Optional ID, if not filled in, it will be generated from name. Cannot be
	// changed later.
	ID *ExportID
	// Human readable name of the export.
	Name string
	// Export column mapping.
	Mapping *Mapping
	// Table import conditions.
	Conditions *Conditions
}

// CreateExportPayload is the payload type of the buffer service CreateExport
// method.
type CreateExportPayload struct {
	StorageAPIToken string
	ReceiverID      ReceiverID
	// Optional ID, if not filled in, it will be generated from name. Cannot be
	// changed later.
	ID *ExportID
	// Human readable name of the export.
	Name string
	// Export column mapping.
	Mapping *Mapping
	// Table import conditions.
	Conditions *Conditions
}

// CreateReceiverPayload is the payload type of the buffer service
// CreateReceiver method.
type CreateReceiverPayload struct {
	StorageAPIToken string
	// Optional ID, if not filled in, it will be generated from name. Cannot be
	// changed later.
	ID *ReceiverID
	// Human readable name of the receiver.
	Name string
	// List of exports, max 20 exports per a receiver.
	Exports []*CreateExportData
}

// DeleteExportPayload is the payload type of the buffer service DeleteExport
// method.
type DeleteExportPayload struct {
	StorageAPIToken string
	ReceiverID      ReceiverID
	ExportID        ExportID
}

// DeleteReceiverPayload is the payload type of the buffer service
// DeleteReceiver method.
type DeleteReceiverPayload struct {
	StorageAPIToken string
	ReceiverID      ReceiverID
}

// Export is the result type of the buffer service GetExport method.
type Export struct {
	ID         ExportID
	ReceiverID ReceiverID
	// Human readable name of the export.
	Name string
	// Export column mapping.
	Mapping *Mapping
	// Table import conditions.
	Conditions *Conditions
}

// Unique ID of the export.
type ExportID = key.ExportID

// ExportsList is the result type of the buffer service ListExports method.
type ExportsList struct {
	Exports []*Export
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

// GetExportPayload is the payload type of the buffer service GetExport method.
type GetExportPayload struct {
	StorageAPIToken string
	ReceiverID      ReceiverID
	ExportID        ExportID
}

// GetReceiverPayload is the payload type of the buffer service GetReceiver
// method.
type GetReceiverPayload struct {
	StorageAPIToken string
	ReceiverID      ReceiverID
}

// GetTaskPayload is the payload type of the buffer service GetTask method.
type GetTaskPayload struct {
	StorageAPIToken string
	ReceiverID      ReceiverID
	Type            string
	TaskID          TaskID
}

// ImportPayload is the payload type of the buffer service Import method.
type ImportPayload struct {
	ProjectID  ProjectID
	ReceiverID ReceiverID
	// Secret used for authentication.
	Secret      string
	ContentType string
}

// ListExportsPayload is the payload type of the buffer service ListExports
// method.
type ListExportsPayload struct {
	StorageAPIToken string
	ReceiverID      ReceiverID
}

// ListReceiversPayload is the payload type of the buffer service ListReceivers
// method.
type ListReceiversPayload struct {
	StorageAPIToken string
}

// Export column mapping.
type Mapping struct {
	// Destination table ID.
	TableID string
	// Enables incremental loading to the table.
	Incremental *bool
	// List of export column mappings. An export may have a maximum of 50 columns.
	Columns []*Column
}

// ID of the project
type ProjectID = key.ProjectID

// Receiver is the result type of the buffer service UpdateReceiver method.
type Receiver struct {
	ID ReceiverID
	// URL of the receiver. Contains secret used for authentication.
	URL string
	// Human readable name of the receiver.
	Name string
	// List of exports, max 20 exports per a receiver.
	Exports []*Export
}

// Unique ID of the receiver.
type ReceiverID = key.ReceiverID

// ReceiversList is the result type of the buffer service ListReceivers method.
type ReceiversList struct {
	Receivers []*Receiver
}

// RefreshReceiverTokensPayload is the payload type of the buffer service
// RefreshReceiverTokens method.
type RefreshReceiverTokensPayload struct {
	StorageAPIToken string
	ReceiverID      ReceiverID
}

// ServiceDetail is the result type of the buffer service ApiVersionIndex
// method.
type ServiceDetail struct {
	// Name of the API
	API string
	// URL of the API documentation.
	Documentation string
}

// Task is the result type of the buffer service CreateReceiver method.
type Task struct {
	ID         TaskID
	ReceiverID ReceiverID
	// URL of the task.
	URL  string
	Type string
	// Date and time of the task creation.
	CreatedAt string
	// Date and time of the task end.
	FinishedAt *string
	IsFinished bool
	// Duration of the task in milliseconds.
	Duration *int64
	Result   *string
	Error    *string
}

// Unique ID of the task.
type TaskID = key.TaskID

type Template struct {
	Language string
	Content  string
}

// UpdateExportPayload is the payload type of the buffer service UpdateExport
// method.
type UpdateExportPayload struct {
	StorageAPIToken string
	// Human readable name of the export.
	Name *string
	// Export column mapping.
	Mapping *Mapping
	// Table import conditions.
	Conditions *Conditions
	ReceiverID ReceiverID
	ExportID   ExportID
}

// UpdateReceiverPayload is the payload type of the buffer service
// UpdateReceiver method.
type UpdateReceiverPayload struct {
	StorageAPIToken string
	// Human readable name of the receiver.
	Name       *string
	ReceiverID ReceiverID
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
