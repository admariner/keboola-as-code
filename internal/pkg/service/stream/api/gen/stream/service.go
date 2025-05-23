// Code generated by goa v3.20.1, DO NOT EDIT.
//
// stream service
//
// Command:
// $ goa gen github.com/keboola/keboola-as-code/api/stream --output
// ./internal/pkg/service/stream/api

package stream

import (
	"context"
	"io"

	"github.com/keboola/keboola-as-code/internal/pkg/service/common/task"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/definition"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/definition/key"
	dependencies "github.com/keboola/keboola-as-code/internal/pkg/service/stream/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/mapping/table/column"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/model"
	"github.com/keboola/keboola-sdk-go/v2/pkg/keboola"
	"goa.design/goa/v3/security"
)

// A service for continuously importing data to the Keboola platform.
type Service interface {
	// Redirect to /v1.
	APIRootIndex(context.Context, dependencies.PublicRequestScope) (err error)
	// List API name and link to documentation.
	APIVersionIndex(context.Context, dependencies.PublicRequestScope) (res *ServiceDetail, err error)
	// HealthCheck implements HealthCheck.
	HealthCheck(context.Context, dependencies.PublicRequestScope) (res string, err error)
	// Create a new source in the branch.
	CreateSource(context.Context, dependencies.BranchRequestScope, *CreateSourcePayload) (res *Task, err error)
	// Update the source.
	UpdateSource(context.Context, dependencies.SourceRequestScope, *UpdateSourcePayload) (res *Task, err error)
	// List all sources in the branch.
	ListSources(context.Context, dependencies.BranchRequestScope, *ListSourcesPayload) (res *SourcesList, err error)
	// List all deleted sources in the branch.
	ListDeletedSources(context.Context, dependencies.BranchRequestScope, *ListDeletedSourcesPayload) (res *SourcesList, err error)
	// Get the source definition.
	GetSource(context.Context, dependencies.SourceRequestScope, *GetSourcePayload) (res *Source, err error)
	// Delete the source.
	DeleteSource(context.Context, dependencies.SourceRequestScope, *DeleteSourcePayload) (res *Task, err error)
	// Get source settings.
	GetSourceSettings(context.Context, dependencies.SourceRequestScope, *GetSourceSettingsPayload) (res *SettingsResult, err error)
	// Update source settings.
	UpdateSourceSettings(context.Context, dependencies.SourceRequestScope, *UpdateSourceSettingsPayload) (res *Task, err error)
	// Tests configured mapping of the source and its sinks.
	TestSource(context.Context, dependencies.SourceRequestScope, *TestSourcePayload, io.ReadCloser) (res *TestResult, err error)
	// Clears all statistics of the source.
	SourceStatisticsClear(context.Context, dependencies.SourceRequestScope, *SourceStatisticsClearPayload) (err error)
	// Disables the source.
	DisableSource(context.Context, dependencies.SourceRequestScope, *DisableSourcePayload) (res *Task, err error)
	// Enables the source.
	EnableSource(context.Context, dependencies.SourceRequestScope, *EnableSourcePayload) (res *Task, err error)
	// Undelete the source.
	UndeleteSource(context.Context, dependencies.SourceRequestScope, *UndeleteSourcePayload) (res *Task, err error)
	// List all source versions.
	ListSourceVersions(context.Context, dependencies.SourceRequestScope, *ListSourceVersionsPayload) (res *EntityVersions, err error)
	// Source version detail.
	SourceVersionDetail(context.Context, dependencies.SourceRequestScope, *SourceVersionDetailPayload) (res *Version, err error)
	// Rollback source version.
	RollbackSourceVersion(context.Context, dependencies.SourceRequestScope, *RollbackSourceVersionPayload) (res *Task, err error)
	// Create a new sink in the source.
	CreateSink(context.Context, dependencies.SourceRequestScope, *CreateSinkPayload) (res *Task, err error)
	// Get the sink definition.
	GetSink(context.Context, dependencies.SinkRequestScope, *GetSinkPayload) (res *Sink, err error)
	// Get the sink settings.
	GetSinkSettings(context.Context, dependencies.SinkRequestScope, *GetSinkSettingsPayload) (res *SettingsResult, err error)
	// Update sink settings.
	UpdateSinkSettings(context.Context, dependencies.SinkRequestScope, *UpdateSinkSettingsPayload) (res *Task, err error)
	// List all sinks in the source.
	ListSinks(context.Context, dependencies.SourceRequestScope, *ListSinksPayload) (res *SinksList, err error)
	// List all deleted sinks in the source.
	ListDeletedSinks(context.Context, dependencies.SourceRequestScope, *ListDeletedSinksPayload) (res *SinksList, err error)
	// Update the sink.
	UpdateSink(context.Context, dependencies.SinkRequestScope, *UpdateSinkPayload) (res *Task, err error)
	// Delete the sink.
	DeleteSink(context.Context, dependencies.SinkRequestScope, *DeleteSinkPayload) (res *Task, err error)
	// Get total statistics of the sink.
	SinkStatisticsTotal(context.Context, dependencies.SinkRequestScope, *SinkStatisticsTotalPayload) (res *SinkStatisticsTotalResult, err error)
	// Get files statistics of the sink.
	SinkStatisticsFiles(context.Context, dependencies.SinkRequestScope, *SinkStatisticsFilesPayload) (res *SinkStatisticsFilesResult, err error)
	// Clears all statistics of the sink.
	SinkStatisticsClear(context.Context, dependencies.SinkRequestScope, *SinkStatisticsClearPayload) (err error)
	// Disables the sink.
	DisableSink(context.Context, dependencies.SinkRequestScope, *DisableSinkPayload) (res *Task, err error)
	// Enables the sink.
	EnableSink(context.Context, dependencies.SinkRequestScope, *EnableSinkPayload) (res *Task, err error)
	// Undelete the sink.
	UndeleteSink(context.Context, dependencies.SinkRequestScope, *UndeleteSinkPayload) (res *Task, err error)
	// List all sink versions.
	ListSinkVersions(context.Context, dependencies.SinkRequestScope, *ListSinkVersionsPayload) (res *EntityVersions, err error)
	// Sink version detail.
	SinkVersionDetail(context.Context, dependencies.SinkRequestScope, *SinkVersionDetailPayload) (res *Version, err error)
	// Rollback sink version.
	RollbackSinkVersion(context.Context, dependencies.SinkRequestScope, *RollbackSinkVersionPayload) (res *Task, err error)
	// Get details of a task.
	GetTask(context.Context, dependencies.ProjectRequestScope, *GetTaskPayload) (res *Task, err error)
	// Details about sources for the UI.
	AggregationSources(context.Context, dependencies.BranchRequestScope, *AggregationSourcesPayload) (res *AggregatedSourcesResult, err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// APIKeyAuth implements the authorization logic for the APIKey security scheme.
	APIKeyAuth(ctx context.Context, key string, schema *security.APIKeyScheme) (context.Context, error)
}

// APIName is the name of the API as defined in the design.
const APIName = "stream"

// APIVersion is the version of the API as defined in the design.
const APIVersion = "1.0"

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "stream"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [38]string{"ApiRootIndex", "ApiVersionIndex", "HealthCheck", "CreateSource", "UpdateSource", "ListSources", "ListDeletedSources", "GetSource", "DeleteSource", "GetSourceSettings", "UpdateSourceSettings", "TestSource", "SourceStatisticsClear", "DisableSource", "EnableSource", "UndeleteSource", "ListSourceVersions", "SourceVersionDetail", "RollbackSourceVersion", "CreateSink", "GetSink", "GetSinkSettings", "UpdateSinkSettings", "ListSinks", "ListDeletedSinks", "UpdateSink", "DeleteSink", "SinkStatisticsTotal", "SinkStatisticsFiles", "SinkStatisticsClear", "DisableSink", "EnableSink", "UndeleteSink", "ListSinkVersions", "SinkVersionDetail", "RollbackSinkVersion", "GetTask", "AggregationSources"}

// A mapping from imported data to a destination table.
type AggregatedSink struct {
	ProjectID ProjectID
	BranchID  BranchID
	SourceID  SourceID
	SinkID    SinkID
	Type      SinkType
	// Human readable name of the sink.
	Name string
	// Description of the source.
	Description string
	Table       *TableSink
	Version     *Version
	Created     *CreatedEntity
	Deleted     *DeletedEntity
	Disabled    *DisabledEntity
	Statistics  *AggregatedStatistics
}

type AggregatedSinks []*AggregatedSink

// Source of data for further processing, start of the stream, max 100 sources
// per a branch.
type AggregatedSource struct {
	ProjectID ProjectID
	BranchID  BranchID
	SourceID  SourceID
	Type      SourceType
	// Human readable name of the source.
	Name string
	// Description of the source.
	Description string
	// HTTP source details for "type" = "http".
	HTTP     *HTTPSource
	Version  *Version
	Created  *CreatedEntity
	Deleted  *DeletedEntity
	Disabled *DisabledEntity
	Sinks    AggregatedSinks
}

type AggregatedSources []*AggregatedSource

// AggregatedSourcesResult is the result type of the stream service
// AggregationSources method.
type AggregatedSourcesResult struct {
	ProjectID ProjectID
	BranchID  BranchID
	Page      *PaginatedResponse
	Sources   AggregatedSources
}

type AggregatedStatistics struct {
	Total  *Level
	Levels *Levels
	Files  SinkFiles
}

// AggregationSourcesPayload is the payload type of the stream service
// AggregationSources method.
type AggregationSourcesPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	// Request records after the ID.
	AfterID string
	// Maximum number of returned records.
	Limit int
}

// ID of the branch.
type BranchID = keboola.BranchID

// ID of the branch or "default".
type BranchIDOrDefault = key.BranchIDOrDefault

// Information about the operation actor.
type By struct {
	// Date and time of deletion.
	Type string
	// ID of the token.
	TokenID *string
	// Description of the token.
	TokenDesc *string
	// ID of the user.
	UserID *string
	// Name of the user.
	UserName *string
}

// CreateSinkPayload is the payload type of the stream service CreateSink
// method.
type CreateSinkPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
	// Optional ID, if not filled in, it will be generated from name. Cannot be
	// changed later.
	SinkID *SinkID
	Type   SinkType
	// Human readable name of the sink.
	Name string
	// Description of the source.
	Description *string
	Table       *TableSinkCreate
}

// CreateSourcePayload is the payload type of the stream service CreateSource
// method.
type CreateSourcePayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	// Optional ID, if not filled in, it will be generated from name. Cannot be
	// changed later.
	SourceID *SourceID
	Type     SourceType
	// Human readable name of the source.
	Name string
	// Description of the source.
	Description *string
}

// Information about the entity creation.
type CreatedEntity struct {
	// Date and time of deletion.
	At string
	// Who created the entity.
	By *By
}

// DeleteSinkPayload is the payload type of the stream service DeleteSink
// method.
type DeleteSinkPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
	SinkID          SinkID
}

// DeleteSourcePayload is the payload type of the stream service DeleteSource
// method.
type DeleteSourcePayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
}

// Information about the deleted entity.
type DeletedEntity struct {
	// Date and time of deletion.
	At string
	// Who deleted the entity, for example "system", "user", ...
	By *By
}

// DisableSinkPayload is the payload type of the stream service DisableSink
// method.
type DisableSinkPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
	SinkID          SinkID
}

// DisableSourcePayload is the payload type of the stream service DisableSource
// method.
type DisableSourcePayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
}

// Information about the disabled entity.
type DisabledEntity struct {
	// Date and time of disabling.
	At string
	// Who disabled the entity, for example "system", "user", ...
	By *By
	// Why was the entity disabled?
	Reason string
}

// EnableSinkPayload is the payload type of the stream service EnableSink
// method.
type EnableSinkPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
	SinkID          SinkID
}

// EnableSourcePayload is the payload type of the stream service EnableSource
// method.
type EnableSourcePayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
}

// EntityVersions is the result type of the stream service ListSourceVersions
// method.
type EntityVersions struct {
	Versions []*Version
	Page     *PaginatedResponse
}

type FileState = model.FileState

// Generic error.
type GenericError struct {
	// HTTP status code.
	StatusCode int
	// Name of error.
	Name string
	// Error message.
	Message string
}

// GetSinkPayload is the payload type of the stream service GetSink method.
type GetSinkPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
	SinkID          SinkID
}

// GetSinkSettingsPayload is the payload type of the stream service
// GetSinkSettings method.
type GetSinkSettingsPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
	SinkID          SinkID
}

// GetSourcePayload is the payload type of the stream service GetSource method.
type GetSourcePayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
}

// GetSourceSettingsPayload is the payload type of the stream service
// GetSourceSettings method.
type GetSourceSettingsPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
}

// GetTaskPayload is the payload type of the stream service GetTask method.
type GetTaskPayload struct {
	StorageAPIToken string
	TaskID          TaskID
}

// HTTP source details for "type" = "http".
type HTTPSource struct {
	// URL of the HTTP source. Contains secret used for authentication.
	URL string
}

type Level struct {
	// Timestamp of the first received record.
	FirstRecordAt *string
	// Timestamp of the last received record.
	LastRecordAt *string
	RecordsCount uint64
	// Compressed size of data in bytes.
	CompressedSize uint64
	// Uncompressed size of data in bytes.
	UncompressedSize uint64
}

type Levels struct {
	Local   *Level
	Staging *Level
	Target  *Level
}

// ListDeletedSinksPayload is the payload type of the stream service
// ListDeletedSinks method.
type ListDeletedSinksPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
	// Request records after the ID.
	AfterID string
	// Maximum number of returned records.
	Limit int
}

// ListDeletedSourcesPayload is the payload type of the stream service
// ListDeletedSources method.
type ListDeletedSourcesPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	// Request records after the ID.
	AfterID string
	// Maximum number of returned records.
	Limit int
}

// ListSinkVersionsPayload is the payload type of the stream service
// ListSinkVersions method.
type ListSinkVersionsPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
	SinkID          SinkID
	// Request records after the ID.
	AfterID string
	// Maximum number of returned records.
	Limit int
}

// ListSinksPayload is the payload type of the stream service ListSinks method.
type ListSinksPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
	// Request records after the ID.
	AfterID string
	// Maximum number of returned records.
	Limit int
}

// ListSourceVersionsPayload is the payload type of the stream service
// ListSourceVersions method.
type ListSourceVersionsPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
	// Request records after the ID.
	AfterID string
	// Maximum number of returned records.
	Limit int
}

// ListSourcesPayload is the payload type of the stream service ListSources
// method.
type ListSourcesPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	// Request records after the ID.
	AfterID string
	// Maximum number of returned records.
	Limit int
}

type PaginatedResponse struct {
	// Current limit.
	Limit int
	// Total count of all records.
	TotalCount int
	// Current offset.
	AfterID string
	// ID of the last record in the response.
	LastID string
}

// ID of the project.
type ProjectID = keboola.ProjectID

// RollbackSinkVersionPayload is the payload type of the stream service
// RollbackSinkVersion method.
type RollbackSinkVersionPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
	SinkID          SinkID
	// Version number counted from 1.
	VersionNumber definition.VersionNumber
}

// RollbackSourceVersionPayload is the payload type of the stream service
// RollbackSourceVersion method.
type RollbackSourceVersionPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
	// Version number counted from 1.
	VersionNumber definition.VersionNumber
}

// ServiceDetail is the result type of the stream service ApiVersionIndex
// method.
type ServiceDetail struct {
	// Name of the API
	API string
	// URL of the API documentation.
	Documentation string
}

// One setting key-value pair.
type SettingPatch struct {
	// Key path.
	Key string
	// A new key value. Use null to reset the value to the default value.
	Value any
}

// One setting key-value pair.
type SettingResult struct {
	// Key path.
	Key string
	// Value type.
	Type string
	// Key description.
	Description string
	// Actual value.
	Value any
	// Default value.
	DefaultValue any
	// True, if the default value is locally overwritten.
	Overwritten bool
	// True, if only a super admin can modify the key.
	Protected bool
	// Validation rules as a string definition.
	Validation *string
}

type SettingsPatch []*SettingPatch

// SettingsResult is the result type of the stream service GetSourceSettings
// method.
type SettingsResult struct {
	Settings []*SettingResult
}

// Sink is the result type of the stream service GetSink method.
type Sink struct {
	ProjectID ProjectID
	BranchID  BranchID
	SourceID  SourceID
	SinkID    SinkID
	Type      SinkType
	// Human readable name of the sink.
	Name string
	// Description of the source.
	Description string
	Table       *TableSink
	Version     *Version
	Created     *CreatedEntity
	Deleted     *DeletedEntity
	Disabled    *DisabledEntity
}

type SinkFile struct {
	State       FileState
	OpenedAt    string
	ClosingAt   *string
	ImportingAt *string
	ImportedAt  *string
	// Number of failed attempts.
	RetryAttempt *int
	// Reason of the last failed attempt.
	RetryReason *string
	// Next attempt time.
	RetryAfter *string
	Statistics *SinkFileStatistics
}

type SinkFileStatistics struct {
	Total  *Level
	Levels *Levels
}

// List of recent sink files.
type SinkFiles []*SinkFile

// Unique ID of the sink.
type SinkID = key.SinkID

// SinkStatisticsClearPayload is the payload type of the stream service
// SinkStatisticsClear method.
type SinkStatisticsClearPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
	SinkID          SinkID
}

// SinkStatisticsFilesPayload is the payload type of the stream service
// SinkStatisticsFiles method.
type SinkStatisticsFilesPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
	SinkID          SinkID
	// Filter for not imported files. If set to true, only not imported files will
	// be included.
	FailedFiles bool
}

// SinkStatisticsFilesResult is the result type of the stream service
// SinkStatisticsFiles method.
type SinkStatisticsFilesResult struct {
	Files SinkFiles
}

// SinkStatisticsTotalPayload is the payload type of the stream service
// SinkStatisticsTotal method.
type SinkStatisticsTotalPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
	SinkID          SinkID
}

// SinkStatisticsTotalResult is the result type of the stream service
// SinkStatisticsTotal method.
type SinkStatisticsTotalResult struct {
	Total  *Level
	Levels *Levels
}

type SinkType = definition.SinkType

// SinkVersionDetailPayload is the payload type of the stream service
// SinkVersionDetail method.
type SinkVersionDetailPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
	SinkID          SinkID
	// Version number counted from 1.
	VersionNumber definition.VersionNumber
}

// List of sinks, max 100 sinks per a source.
type Sinks []*Sink

// SinksList is the result type of the stream service ListSinks method.
type SinksList struct {
	ProjectID ProjectID
	BranchID  BranchID
	SourceID  SourceID
	Page      *PaginatedResponse
	Sinks     Sinks
}

// Source is the result type of the stream service GetSource method.
type Source struct {
	ProjectID ProjectID
	BranchID  BranchID
	SourceID  SourceID
	Type      SourceType
	// Human readable name of the source.
	Name string
	// Description of the source.
	Description string
	// HTTP source details for "type" = "http".
	HTTP     *HTTPSource
	Version  *Version
	Created  *CreatedEntity
	Deleted  *DeletedEntity
	Disabled *DisabledEntity
}

// Unique ID of the source.
type SourceID = key.SourceID

// SourceStatisticsClearPayload is the payload type of the stream service
// SourceStatisticsClear method.
type SourceStatisticsClearPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
}

type SourceType = definition.SourceType

// SourceVersionDetailPayload is the payload type of the stream service
// SourceVersionDetail method.
type SourceVersionDetailPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
	// Version number counted from 1.
	VersionNumber definition.VersionNumber
}

// List of sources, max 100 sources per a branch.
type Sources []*Source

// SourcesList is the result type of the stream service ListSources method.
type SourcesList struct {
	ProjectID ProjectID
	BranchID  BranchID
	Page      *PaginatedResponse
	Sources   Sources
}

// An output mapping defined by a template.
type TableColumn struct {
	// Column mapping type. This represents a static mapping (e.g. `body` or
	// `headers`), or a custom mapping using a template language (`template`).
	Type column.Type
	// Column name.
	Name string
	// Path to the value.
	Path *string
	// Fallback value if path doesn't exist.
	DefaultValue *string
	// Set to true if path value should use raw string instead of json-encoded
	// value.
	RawString *bool
	// Template mapping details. Only for "type" = "template".
	Template *TableColumnTemplate
}

// Template column definition, for "type" = "template".
type TableColumnTemplate struct {
	Language string
	Content  string
}

// List of export column mappings. An export may have a maximum of 100 columns.
type TableColumns []*TableColumn

type TableID string

// Table mapping definition.
type TableMapping struct {
	Columns TableColumns
}

// Table sink configuration for "type" = "table".
type TableSink struct {
	Type    TableType
	TableID TableID
	Mapping *TableMapping
}

// Table sink configuration for "type" = "table".
type TableSinkCreate struct {
	Type    TableType
	TableID TableID
	Mapping *TableMapping
}

// Table sink configuration for "type" = "table".
type TableSinkUpdate struct {
	Type    *TableType
	TableID *TableID
	Mapping *TableMapping
}

type TableType = definition.TableType

// Task is the result type of the stream service CreateSource method.
type Task struct {
	TaskID TaskID
	// Task type.
	Type string
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
	Outputs  *TaskOutputs
}

// Unique ID of the task.
type TaskID = task.ID

// Outputs generated by the task.
type TaskOutputs struct {
	// Absolute URL of the entity.
	URL *string
	// ID of the parent project.
	ProjectID *ProjectID
	// ID of the parent branch.
	BranchID *BranchID
	// ID of the created/updated source.
	SourceID *SourceID
	// ID of the created/updated sink.
	SinkID *SinkID
}

// TestResult is the result type of the stream service TestSource method.
type TestResult struct {
	ProjectID ProjectID
	BranchID  BranchID
	SourceID  SourceID
	// Table for each configured sink.
	Tables []*TestResultTable
}

// Generated table column value, part of the test result.
type TestResultColumn struct {
	// Column name.
	Name string
	// Column value.
	Value string
}

// Generated table row, part of the test result.
type TestResultRow struct {
	// Generated columns.
	Columns []*TestResultColumn
}

// Generated table rows, part of the test result.
type TestResultTable struct {
	SinkID  SinkID
	TableID TableID
	// Generated rows.
	Rows []*TestResultRow
}

// TestSourcePayload is the payload type of the stream service TestSource
// method.
type TestSourcePayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
}

// UndeleteSinkPayload is the payload type of the stream service UndeleteSink
// method.
type UndeleteSinkPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
	SinkID          SinkID
}

// UndeleteSourcePayload is the payload type of the stream service
// UndeleteSource method.
type UndeleteSourcePayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
}

// UpdateSinkPayload is the payload type of the stream service UpdateSink
// method.
type UpdateSinkPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
	SinkID          SinkID
	// Description of the modification, description of the version.
	ChangeDescription *string
	Type              *SinkType
	// Human readable name of the sink.
	Name *string
	// Description of the source.
	Description *string
	Table       *TableSinkUpdate
}

// UpdateSinkSettingsPayload is the payload type of the stream service
// UpdateSinkSettings method.
type UpdateSinkSettingsPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
	SinkID          SinkID
	// Description of the modification, description of the version.
	ChangeDescription *string
	Settings          SettingsPatch
}

// UpdateSourcePayload is the payload type of the stream service UpdateSource
// method.
type UpdateSourcePayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
	// Description of the modification, description of the version.
	ChangeDescription *string
	Type              *SourceType
	// Human readable name of the source.
	Name *string
	// Description of the source.
	Description *string
}

// UpdateSourceSettingsPayload is the payload type of the stream service
// UpdateSourceSettings method.
type UpdateSourceSettingsPayload struct {
	StorageAPIToken string
	BranchID        BranchIDOrDefault
	SourceID        SourceID
	// Description of the modification, description of the version.
	ChangeDescription *string
	Settings          SettingsPatch
}

// Version is the result type of the stream service SourceVersionDetail method.
type Version struct {
	// Version number counted from 1.
	Number definition.VersionNumber
	// Hash of the entity state.
	Hash string
	// Description of the change.
	Description string
	// Date and time of the modification.
	At string
	// Who modified the entity.
	By *By
}

// Error returns an error description.
func (e *GenericError) Error() string {
	return "Generic error."
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
