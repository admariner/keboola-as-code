// Code generated by goa v3.19.1, DO NOT EDIT.
//
// stream client
//
// Command:
// $ goa gen github.com/keboola/keboola-as-code/api/stream --output
// ./internal/pkg/service/stream/api

package stream

import (
	"context"
	"io"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "stream" service client.
type Client struct {
	APIRootIndexEndpoint          goa.Endpoint
	APIVersionIndexEndpoint       goa.Endpoint
	HealthCheckEndpoint           goa.Endpoint
	CreateSourceEndpoint          goa.Endpoint
	UpdateSourceEndpoint          goa.Endpoint
	ListSourcesEndpoint           goa.Endpoint
	GetSourceEndpoint             goa.Endpoint
	DeleteSourceEndpoint          goa.Endpoint
	GetSourceSettingsEndpoint     goa.Endpoint
	UpdateSourceSettingsEndpoint  goa.Endpoint
	TestSourceEndpoint            goa.Endpoint
	SourceStatisticsClearEndpoint goa.Endpoint
	DisableSourceEndpoint         goa.Endpoint
	EnableSourceEndpoint          goa.Endpoint
	ListSourceVersionsEndpoint    goa.Endpoint
	CreateSinkEndpoint            goa.Endpoint
	GetSinkEndpoint               goa.Endpoint
	GetSinkSettingsEndpoint       goa.Endpoint
	UpdateSinkSettingsEndpoint    goa.Endpoint
	ListSinksEndpoint             goa.Endpoint
	UpdateSinkEndpoint            goa.Endpoint
	DeleteSinkEndpoint            goa.Endpoint
	SinkStatisticsTotalEndpoint   goa.Endpoint
	SinkStatisticsFilesEndpoint   goa.Endpoint
	SinkStatisticsClearEndpoint   goa.Endpoint
	DisableSinkEndpoint           goa.Endpoint
	EnableSinkEndpoint            goa.Endpoint
	ListSinkVersionsEndpoint      goa.Endpoint
	GetTaskEndpoint               goa.Endpoint
	AggregationSourcesEndpoint    goa.Endpoint
}

// NewClient initializes a "stream" service client given the endpoints.
func NewClient(aPIRootIndex, aPIVersionIndex, healthCheck, createSource, updateSource, listSources, getSource, deleteSource, getSourceSettings, updateSourceSettings, testSource, sourceStatisticsClear, disableSource, enableSource, listSourceVersions, createSink, getSink, getSinkSettings, updateSinkSettings, listSinks, updateSink, deleteSink, sinkStatisticsTotal, sinkStatisticsFiles, sinkStatisticsClear, disableSink, enableSink, listSinkVersions, getTask, aggregationSources goa.Endpoint) *Client {
	return &Client{
		APIRootIndexEndpoint:          aPIRootIndex,
		APIVersionIndexEndpoint:       aPIVersionIndex,
		HealthCheckEndpoint:           healthCheck,
		CreateSourceEndpoint:          createSource,
		UpdateSourceEndpoint:          updateSource,
		ListSourcesEndpoint:           listSources,
		GetSourceEndpoint:             getSource,
		DeleteSourceEndpoint:          deleteSource,
		GetSourceSettingsEndpoint:     getSourceSettings,
		UpdateSourceSettingsEndpoint:  updateSourceSettings,
		TestSourceEndpoint:            testSource,
		SourceStatisticsClearEndpoint: sourceStatisticsClear,
		DisableSourceEndpoint:         disableSource,
		EnableSourceEndpoint:          enableSource,
		ListSourceVersionsEndpoint:    listSourceVersions,
		CreateSinkEndpoint:            createSink,
		GetSinkEndpoint:               getSink,
		GetSinkSettingsEndpoint:       getSinkSettings,
		UpdateSinkSettingsEndpoint:    updateSinkSettings,
		ListSinksEndpoint:             listSinks,
		UpdateSinkEndpoint:            updateSink,
		DeleteSinkEndpoint:            deleteSink,
		SinkStatisticsTotalEndpoint:   sinkStatisticsTotal,
		SinkStatisticsFilesEndpoint:   sinkStatisticsFiles,
		SinkStatisticsClearEndpoint:   sinkStatisticsClear,
		DisableSinkEndpoint:           disableSink,
		EnableSinkEndpoint:            enableSink,
		ListSinkVersionsEndpoint:      listSinkVersions,
		GetTaskEndpoint:               getTask,
		AggregationSourcesEndpoint:    aggregationSources,
	}
}

// APIRootIndex calls the "ApiRootIndex" endpoint of the "stream" service.
func (c *Client) APIRootIndex(ctx context.Context) (err error) {
	_, err = c.APIRootIndexEndpoint(ctx, nil)
	return
}

// APIVersionIndex calls the "ApiVersionIndex" endpoint of the "stream" service.
func (c *Client) APIVersionIndex(ctx context.Context) (res *ServiceDetail, err error) {
	var ires any
	ires, err = c.APIVersionIndexEndpoint(ctx, nil)
	if err != nil {
		return
	}
	return ires.(*ServiceDetail), nil
}

// HealthCheck calls the "HealthCheck" endpoint of the "stream" service.
func (c *Client) HealthCheck(ctx context.Context) (res string, err error) {
	var ires any
	ires, err = c.HealthCheckEndpoint(ctx, nil)
	if err != nil {
		return
	}
	return ires.(string), nil
}

// CreateSource calls the "CreateSource" endpoint of the "stream" service.
// CreateSource may return the following errors:
//   - "stream.api.sourceAlreadyExists" (type *GenericError): Source already exists in the branch.
//   - "stream.api.resourceLimitReached" (type *GenericError): Resource limit reached.
//   - error: internal error
func (c *Client) CreateSource(ctx context.Context, p *CreateSourcePayload) (res *Task, err error) {
	var ires any
	ires, err = c.CreateSourceEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Task), nil
}

// UpdateSource calls the "UpdateSource" endpoint of the "stream" service.
// UpdateSource may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - error: internal error
func (c *Client) UpdateSource(ctx context.Context, p *UpdateSourcePayload) (res *Task, err error) {
	var ires any
	ires, err = c.UpdateSourceEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Task), nil
}

// ListSources calls the "ListSources" endpoint of the "stream" service.
func (c *Client) ListSources(ctx context.Context, p *ListSourcesPayload) (res *SourcesList, err error) {
	var ires any
	ires, err = c.ListSourcesEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*SourcesList), nil
}

// GetSource calls the "GetSource" endpoint of the "stream" service.
// GetSource may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - error: internal error
func (c *Client) GetSource(ctx context.Context, p *GetSourcePayload) (res *Source, err error) {
	var ires any
	ires, err = c.GetSourceEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Source), nil
}

// DeleteSource calls the "DeleteSource" endpoint of the "stream" service.
// DeleteSource may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - error: internal error
func (c *Client) DeleteSource(ctx context.Context, p *DeleteSourcePayload) (res *Task, err error) {
	var ires any
	ires, err = c.DeleteSourceEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Task), nil
}

// GetSourceSettings calls the "GetSourceSettings" endpoint of the "stream"
// service.
// GetSourceSettings may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - error: internal error
func (c *Client) GetSourceSettings(ctx context.Context, p *GetSourceSettingsPayload) (res *SettingsResult, err error) {
	var ires any
	ires, err = c.GetSourceSettingsEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*SettingsResult), nil
}

// UpdateSourceSettings calls the "UpdateSourceSettings" endpoint of the
// "stream" service.
// UpdateSourceSettings may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - "stream.api.forbidden" (type *GenericError): Modification of protected settings is forbidden.
//   - error: internal error
func (c *Client) UpdateSourceSettings(ctx context.Context, p *UpdateSourceSettingsPayload) (res *Task, err error) {
	var ires any
	ires, err = c.UpdateSourceSettingsEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Task), nil
}

// TestSource calls the "TestSource" endpoint of the "stream" service.
// TestSource may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - error: internal error
func (c *Client) TestSource(ctx context.Context, p *TestSourcePayload, req io.ReadCloser) (res *TestResult, err error) {
	var ires any
	ires, err = c.TestSourceEndpoint(ctx, &TestSourceRequestData{Payload: p, Body: req})
	if err != nil {
		return
	}
	return ires.(*TestResult), nil
}

// SourceStatisticsClear calls the "SourceStatisticsClear" endpoint of the
// "stream" service.
// SourceStatisticsClear may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - error: internal error
func (c *Client) SourceStatisticsClear(ctx context.Context, p *SourceStatisticsClearPayload) (err error) {
	_, err = c.SourceStatisticsClearEndpoint(ctx, p)
	return
}

// DisableSource calls the "DisableSource" endpoint of the "stream" service.
// DisableSource may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - error: internal error
func (c *Client) DisableSource(ctx context.Context, p *DisableSourcePayload) (res *Task, err error) {
	var ires any
	ires, err = c.DisableSourceEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Task), nil
}

// EnableSource calls the "EnableSource" endpoint of the "stream" service.
// EnableSource may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - error: internal error
func (c *Client) EnableSource(ctx context.Context, p *EnableSourcePayload) (res *Task, err error) {
	var ires any
	ires, err = c.EnableSourceEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Task), nil
}

// ListSourceVersions calls the "ListSourceVersions" endpoint of the "stream"
// service.
// ListSourceVersions may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - error: internal error
func (c *Client) ListSourceVersions(ctx context.Context, p *ListSourceVersionsPayload) (res *EntityVersions, err error) {
	var ires any
	ires, err = c.ListSourceVersionsEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*EntityVersions), nil
}

// CreateSink calls the "CreateSink" endpoint of the "stream" service.
// CreateSink may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - "stream.api.sinkAlreadyExists" (type *GenericError): Sink already exists in the source.
//   - "stream.api.resourceLimitReached" (type *GenericError): Resource limit reached.
//   - error: internal error
func (c *Client) CreateSink(ctx context.Context, p *CreateSinkPayload) (res *Task, err error) {
	var ires any
	ires, err = c.CreateSinkEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Task), nil
}

// GetSink calls the "GetSink" endpoint of the "stream" service.
// GetSink may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - "stream.api.sinkNotFound" (type *GenericError): Sink not found error.
//   - error: internal error
func (c *Client) GetSink(ctx context.Context, p *GetSinkPayload) (res *Sink, err error) {
	var ires any
	ires, err = c.GetSinkEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Sink), nil
}

// GetSinkSettings calls the "GetSinkSettings" endpoint of the "stream" service.
// GetSinkSettings may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - "stream.api.sinkNotFound" (type *GenericError): Sink not found error.
//   - error: internal error
func (c *Client) GetSinkSettings(ctx context.Context, p *GetSinkSettingsPayload) (res *SettingsResult, err error) {
	var ires any
	ires, err = c.GetSinkSettingsEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*SettingsResult), nil
}

// UpdateSinkSettings calls the "UpdateSinkSettings" endpoint of the "stream"
// service.
// UpdateSinkSettings may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - "stream.api.sinkNotFound" (type *GenericError): Sink not found error.
//   - "stream.api.forbidden" (type *GenericError): Modification of protected settings is forbidden.
//   - error: internal error
func (c *Client) UpdateSinkSettings(ctx context.Context, p *UpdateSinkSettingsPayload) (res *Task, err error) {
	var ires any
	ires, err = c.UpdateSinkSettingsEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Task), nil
}

// ListSinks calls the "ListSinks" endpoint of the "stream" service.
// ListSinks may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - error: internal error
func (c *Client) ListSinks(ctx context.Context, p *ListSinksPayload) (res *SinksList, err error) {
	var ires any
	ires, err = c.ListSinksEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*SinksList), nil
}

// UpdateSink calls the "UpdateSink" endpoint of the "stream" service.
// UpdateSink may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - "stream.api.sinkNotFound" (type *GenericError): Sink not found error.
//   - error: internal error
func (c *Client) UpdateSink(ctx context.Context, p *UpdateSinkPayload) (res *Task, err error) {
	var ires any
	ires, err = c.UpdateSinkEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Task), nil
}

// DeleteSink calls the "DeleteSink" endpoint of the "stream" service.
// DeleteSink may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - "stream.api.sinkNotFound" (type *GenericError): Sink not found error.
//   - error: internal error
func (c *Client) DeleteSink(ctx context.Context, p *DeleteSinkPayload) (res *Task, err error) {
	var ires any
	ires, err = c.DeleteSinkEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Task), nil
}

// SinkStatisticsTotal calls the "SinkStatisticsTotal" endpoint of the "stream"
// service.
// SinkStatisticsTotal may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - "stream.api.sinkNotFound" (type *GenericError): Sink not found error.
//   - error: internal error
func (c *Client) SinkStatisticsTotal(ctx context.Context, p *SinkStatisticsTotalPayload) (res *SinkStatisticsTotalResult, err error) {
	var ires any
	ires, err = c.SinkStatisticsTotalEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*SinkStatisticsTotalResult), nil
}

// SinkStatisticsFiles calls the "SinkStatisticsFiles" endpoint of the "stream"
// service.
// SinkStatisticsFiles may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - "stream.api.sinkNotFound" (type *GenericError): Sink not found error.
//   - error: internal error
func (c *Client) SinkStatisticsFiles(ctx context.Context, p *SinkStatisticsFilesPayload) (res *SinkStatisticsFilesResult, err error) {
	var ires any
	ires, err = c.SinkStatisticsFilesEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*SinkStatisticsFilesResult), nil
}

// SinkStatisticsClear calls the "SinkStatisticsClear" endpoint of the "stream"
// service.
// SinkStatisticsClear may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - "stream.api.sinkNotFound" (type *GenericError): Sink not found error.
//   - error: internal error
func (c *Client) SinkStatisticsClear(ctx context.Context, p *SinkStatisticsClearPayload) (err error) {
	_, err = c.SinkStatisticsClearEndpoint(ctx, p)
	return
}

// DisableSink calls the "DisableSink" endpoint of the "stream" service.
// DisableSink may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - "stream.api.sinkNotFound" (type *GenericError): Sink not found error.
//   - error: internal error
func (c *Client) DisableSink(ctx context.Context, p *DisableSinkPayload) (res *Task, err error) {
	var ires any
	ires, err = c.DisableSinkEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Task), nil
}

// EnableSink calls the "EnableSink" endpoint of the "stream" service.
// EnableSink may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - "stream.api.sinkNotFound" (type *GenericError): Sink not found error.
//   - error: internal error
func (c *Client) EnableSink(ctx context.Context, p *EnableSinkPayload) (res *Task, err error) {
	var ires any
	ires, err = c.EnableSinkEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Task), nil
}

// ListSinkVersions calls the "ListSinkVersions" endpoint of the "stream"
// service.
// ListSinkVersions may return the following errors:
//   - "stream.api.sourceNotFound" (type *GenericError): Source not found error.
//   - "stream.api.sinkNotFound" (type *GenericError): Sink not found error.
//   - error: internal error
func (c *Client) ListSinkVersions(ctx context.Context, p *ListSinkVersionsPayload) (res *EntityVersions, err error) {
	var ires any
	ires, err = c.ListSinkVersionsEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*EntityVersions), nil
}

// GetTask calls the "GetTask" endpoint of the "stream" service.
// GetTask may return the following errors:
//   - "stream.api.taskNotFound" (type *GenericError): Task not found error.
//   - error: internal error
func (c *Client) GetTask(ctx context.Context, p *GetTaskPayload) (res *Task, err error) {
	var ires any
	ires, err = c.GetTaskEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Task), nil
}

// AggregationSources calls the "AggregationSources" endpoint of the "stream"
// service.
func (c *Client) AggregationSources(ctx context.Context, p *AggregationSourcesPayload) (res *AggregatedSourcesResult, err error) {
	var ires any
	ires, err = c.AggregationSourcesEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*AggregatedSourcesResult), nil
}
