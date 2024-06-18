// Code generated by goa v3.16.2, DO NOT EDIT.
//
// HTTP request path constructors for the stream service.
//
// Command:
// $ goa gen github.com/keboola/keboola-as-code/api/stream --output
// ./internal/pkg/service/stream/api

package server

import (
	"fmt"
)

// APIRootIndexStreamPath returns the URL path to the stream service ApiRootIndex HTTP endpoint.
func APIRootIndexStreamPath() string {
	return "/"
}

// APIVersionIndexStreamPath returns the URL path to the stream service ApiVersionIndex HTTP endpoint.
func APIVersionIndexStreamPath() string {
	return "/v1"
}

// HealthCheckStreamPath returns the URL path to the stream service HealthCheck HTTP endpoint.
func HealthCheckStreamPath() string {
	return "/health-check"
}

// CreateSourceStreamPath returns the URL path to the stream service CreateSource HTTP endpoint.
func CreateSourceStreamPath(branchID string) string {
	return fmt.Sprintf("/v1/branches/%v/sources", branchID)
}

// UpdateSourceStreamPath returns the URL path to the stream service UpdateSource HTTP endpoint.
func UpdateSourceStreamPath(branchID string, sourceID string) string {
	return fmt.Sprintf("/v1/branches/%v/sources/%v", branchID, sourceID)
}

// ListSourcesStreamPath returns the URL path to the stream service ListSources HTTP endpoint.
func ListSourcesStreamPath(branchID string) string {
	return fmt.Sprintf("/v1/branches/%v/sources", branchID)
}

// GetSourceStreamPath returns the URL path to the stream service GetSource HTTP endpoint.
func GetSourceStreamPath(branchID string, sourceID string) string {
	return fmt.Sprintf("/v1/branches/%v/sources/%v", branchID, sourceID)
}

// DeleteSourceStreamPath returns the URL path to the stream service DeleteSource HTTP endpoint.
func DeleteSourceStreamPath(branchID string, sourceID string) string {
	return fmt.Sprintf("/v1/branches/%v/sources/%v", branchID, sourceID)
}

// GetSourceSettingsStreamPath returns the URL path to the stream service GetSourceSettings HTTP endpoint.
func GetSourceSettingsStreamPath(branchID string, sourceID string) string {
	return fmt.Sprintf("/v1/branches/%v/sources/%v/settings", branchID, sourceID)
}

// UpdateSourceSettingsStreamPath returns the URL path to the stream service UpdateSourceSettings HTTP endpoint.
func UpdateSourceSettingsStreamPath(branchID string, sourceID string) string {
	return fmt.Sprintf("/v1/branches/%v/sources/%v/settings", branchID, sourceID)
}

// TestSourceStreamPath returns the URL path to the stream service TestSource HTTP endpoint.
func TestSourceStreamPath(branchID string, sourceID string) string {
	return fmt.Sprintf("/v1/branches/%v/sources/%v/test", branchID, sourceID)
}

// CreateSinkStreamPath returns the URL path to the stream service CreateSink HTTP endpoint.
func CreateSinkStreamPath(branchID string, sourceID string) string {
	return fmt.Sprintf("/v1/branches/%v/sources/%v/sinks", branchID, sourceID)
}

// GetSinkStreamPath returns the URL path to the stream service GetSink HTTP endpoint.
func GetSinkStreamPath(branchID string, sourceID string, sinkID string) string {
	return fmt.Sprintf("/v1/branches/%v/sources/%v/sinks/%v", branchID, sourceID, sinkID)
}

// GetSinkSettingsStreamPath returns the URL path to the stream service GetSinkSettings HTTP endpoint.
func GetSinkSettingsStreamPath(branchID string, sourceID string, sinkID string) string {
	return fmt.Sprintf("/v1/branches/%v/sources/%v/sinks/%v/settings", branchID, sourceID, sinkID)
}

// UpdateSinkSettingsStreamPath returns the URL path to the stream service UpdateSinkSettings HTTP endpoint.
func UpdateSinkSettingsStreamPath(branchID string, sourceID string, sinkID string) string {
	return fmt.Sprintf("/v1/branches/%v/sources/%v/sinks/%v/settings", branchID, sourceID, sinkID)
}

// ListSinksStreamPath returns the URL path to the stream service ListSinks HTTP endpoint.
func ListSinksStreamPath(branchID string, sourceID string) string {
	return fmt.Sprintf("/v1/branches/%v/sources/%v/sinks", branchID, sourceID)
}

// UpdateSinkStreamPath returns the URL path to the stream service UpdateSink HTTP endpoint.
func UpdateSinkStreamPath(branchID string, sourceID string, sinkID string) string {
	return fmt.Sprintf("/v1/branches/%v/sources/%v/sinks/%v", branchID, sourceID, sinkID)
}

// DeleteSinkStreamPath returns the URL path to the stream service DeleteSink HTTP endpoint.
func DeleteSinkStreamPath(branchID string, sourceID string, sinkID string) string {
	return fmt.Sprintf("/v1/branches/%v/sources/%v/sinks/%v", branchID, sourceID, sinkID)
}

// SinkStatisticsTotalStreamPath returns the URL path to the stream service SinkStatisticsTotal HTTP endpoint.
func SinkStatisticsTotalStreamPath(branchID string, sourceID string, sinkID string) string {
	return fmt.Sprintf("/v1/branches/%v/sources/%v/sinks/%v/statistics/total", branchID, sourceID, sinkID)
}

// SinkStatisticsFilesStreamPath returns the URL path to the stream service SinkStatisticsFiles HTTP endpoint.
func SinkStatisticsFilesStreamPath(branchID string, sourceID string, sinkID string) string {
	return fmt.Sprintf("/v1/branches/%v/sources/%v/sinks/%v/statistics/files", branchID, sourceID, sinkID)
}

// GetTaskStreamPath returns the URL path to the stream service GetTask HTTP endpoint.
func GetTaskStreamPath(taskID string) string {
	return fmt.Sprintf("/v1/tasks/%v", taskID)
}

// AggregateSourcesStreamPath returns the URL path to the stream service AggregateSources HTTP endpoint.
func AggregateSourcesStreamPath(branchID string) string {
	return fmt.Sprintf("/v1/branches/%v/aggregation/sources", branchID)
}
