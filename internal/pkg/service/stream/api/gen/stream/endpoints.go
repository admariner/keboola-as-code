// Code generated by goa v3.20.1, DO NOT EDIT.
//
// stream endpoints
//
// Command:
// $ goa gen github.com/keboola/keboola-as-code/api/stream --output
// ./internal/pkg/service/stream/api

package stream

import (
	"context"
	"io"

	dependencies "github.com/keboola/keboola-as-code/internal/pkg/service/stream/dependencies"
	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Endpoints wraps the "stream" service endpoints.
type Endpoints struct {
	APIRootIndex          goa.Endpoint
	APIVersionIndex       goa.Endpoint
	HealthCheck           goa.Endpoint
	CreateSource          goa.Endpoint
	UpdateSource          goa.Endpoint
	ListSources           goa.Endpoint
	ListDeletedSources    goa.Endpoint
	GetSource             goa.Endpoint
	DeleteSource          goa.Endpoint
	GetSourceSettings     goa.Endpoint
	UpdateSourceSettings  goa.Endpoint
	TestSource            goa.Endpoint
	SourceStatisticsClear goa.Endpoint
	DisableSource         goa.Endpoint
	EnableSource          goa.Endpoint
	UndeleteSource        goa.Endpoint
	ListSourceVersions    goa.Endpoint
	SourceVersionDetail   goa.Endpoint
	RollbackSourceVersion goa.Endpoint
	CreateSink            goa.Endpoint
	GetSink               goa.Endpoint
	GetSinkSettings       goa.Endpoint
	UpdateSinkSettings    goa.Endpoint
	ListSinks             goa.Endpoint
	ListDeletedSinks      goa.Endpoint
	UpdateSink            goa.Endpoint
	DeleteSink            goa.Endpoint
	SinkStatisticsTotal   goa.Endpoint
	SinkStatisticsFiles   goa.Endpoint
	SinkStatisticsClear   goa.Endpoint
	DisableSink           goa.Endpoint
	EnableSink            goa.Endpoint
	UndeleteSink          goa.Endpoint
	ListSinkVersions      goa.Endpoint
	SinkVersionDetail     goa.Endpoint
	RollbackSinkVersion   goa.Endpoint
	GetTask               goa.Endpoint
	AggregationSources    goa.Endpoint
}

// TestSourceRequestData holds both the payload and the HTTP request body
// reader of the "TestSource" method.
type TestSourceRequestData struct {
	// Payload is the method payload.
	Payload *TestSourcePayload
	// Body streams the HTTP request body.
	Body io.ReadCloser
}

// NewEndpoints wraps the methods of the "stream" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		APIRootIndex:          NewAPIRootIndexEndpoint(s),
		APIVersionIndex:       NewAPIVersionIndexEndpoint(s),
		HealthCheck:           NewHealthCheckEndpoint(s),
		CreateSource:          NewCreateSourceEndpoint(s, a.APIKeyAuth),
		UpdateSource:          NewUpdateSourceEndpoint(s, a.APIKeyAuth),
		ListSources:           NewListSourcesEndpoint(s, a.APIKeyAuth),
		ListDeletedSources:    NewListDeletedSourcesEndpoint(s, a.APIKeyAuth),
		GetSource:             NewGetSourceEndpoint(s, a.APIKeyAuth),
		DeleteSource:          NewDeleteSourceEndpoint(s, a.APIKeyAuth),
		GetSourceSettings:     NewGetSourceSettingsEndpoint(s, a.APIKeyAuth),
		UpdateSourceSettings:  NewUpdateSourceSettingsEndpoint(s, a.APIKeyAuth),
		TestSource:            NewTestSourceEndpoint(s, a.APIKeyAuth),
		SourceStatisticsClear: NewSourceStatisticsClearEndpoint(s, a.APIKeyAuth),
		DisableSource:         NewDisableSourceEndpoint(s, a.APIKeyAuth),
		EnableSource:          NewEnableSourceEndpoint(s, a.APIKeyAuth),
		UndeleteSource:        NewUndeleteSourceEndpoint(s, a.APIKeyAuth),
		ListSourceVersions:    NewListSourceVersionsEndpoint(s, a.APIKeyAuth),
		SourceVersionDetail:   NewSourceVersionDetailEndpoint(s, a.APIKeyAuth),
		RollbackSourceVersion: NewRollbackSourceVersionEndpoint(s, a.APIKeyAuth),
		CreateSink:            NewCreateSinkEndpoint(s, a.APIKeyAuth),
		GetSink:               NewGetSinkEndpoint(s, a.APIKeyAuth),
		GetSinkSettings:       NewGetSinkSettingsEndpoint(s, a.APIKeyAuth),
		UpdateSinkSettings:    NewUpdateSinkSettingsEndpoint(s, a.APIKeyAuth),
		ListSinks:             NewListSinksEndpoint(s, a.APIKeyAuth),
		ListDeletedSinks:      NewListDeletedSinksEndpoint(s, a.APIKeyAuth),
		UpdateSink:            NewUpdateSinkEndpoint(s, a.APIKeyAuth),
		DeleteSink:            NewDeleteSinkEndpoint(s, a.APIKeyAuth),
		SinkStatisticsTotal:   NewSinkStatisticsTotalEndpoint(s, a.APIKeyAuth),
		SinkStatisticsFiles:   NewSinkStatisticsFilesEndpoint(s, a.APIKeyAuth),
		SinkStatisticsClear:   NewSinkStatisticsClearEndpoint(s, a.APIKeyAuth),
		DisableSink:           NewDisableSinkEndpoint(s, a.APIKeyAuth),
		EnableSink:            NewEnableSinkEndpoint(s, a.APIKeyAuth),
		UndeleteSink:          NewUndeleteSinkEndpoint(s, a.APIKeyAuth),
		ListSinkVersions:      NewListSinkVersionsEndpoint(s, a.APIKeyAuth),
		SinkVersionDetail:     NewSinkVersionDetailEndpoint(s, a.APIKeyAuth),
		RollbackSinkVersion:   NewRollbackSinkVersionEndpoint(s, a.APIKeyAuth),
		GetTask:               NewGetTaskEndpoint(s, a.APIKeyAuth),
		AggregationSources:    NewAggregationSourcesEndpoint(s, a.APIKeyAuth),
	}
}

// Use applies the given middleware to all the "stream" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.APIRootIndex = m(e.APIRootIndex)
	e.APIVersionIndex = m(e.APIVersionIndex)
	e.HealthCheck = m(e.HealthCheck)
	e.CreateSource = m(e.CreateSource)
	e.UpdateSource = m(e.UpdateSource)
	e.ListSources = m(e.ListSources)
	e.ListDeletedSources = m(e.ListDeletedSources)
	e.GetSource = m(e.GetSource)
	e.DeleteSource = m(e.DeleteSource)
	e.GetSourceSettings = m(e.GetSourceSettings)
	e.UpdateSourceSettings = m(e.UpdateSourceSettings)
	e.TestSource = m(e.TestSource)
	e.SourceStatisticsClear = m(e.SourceStatisticsClear)
	e.DisableSource = m(e.DisableSource)
	e.EnableSource = m(e.EnableSource)
	e.UndeleteSource = m(e.UndeleteSource)
	e.ListSourceVersions = m(e.ListSourceVersions)
	e.SourceVersionDetail = m(e.SourceVersionDetail)
	e.RollbackSourceVersion = m(e.RollbackSourceVersion)
	e.CreateSink = m(e.CreateSink)
	e.GetSink = m(e.GetSink)
	e.GetSinkSettings = m(e.GetSinkSettings)
	e.UpdateSinkSettings = m(e.UpdateSinkSettings)
	e.ListSinks = m(e.ListSinks)
	e.ListDeletedSinks = m(e.ListDeletedSinks)
	e.UpdateSink = m(e.UpdateSink)
	e.DeleteSink = m(e.DeleteSink)
	e.SinkStatisticsTotal = m(e.SinkStatisticsTotal)
	e.SinkStatisticsFiles = m(e.SinkStatisticsFiles)
	e.SinkStatisticsClear = m(e.SinkStatisticsClear)
	e.DisableSink = m(e.DisableSink)
	e.EnableSink = m(e.EnableSink)
	e.UndeleteSink = m(e.UndeleteSink)
	e.ListSinkVersions = m(e.ListSinkVersions)
	e.SinkVersionDetail = m(e.SinkVersionDetail)
	e.RollbackSinkVersion = m(e.RollbackSinkVersion)
	e.GetTask = m(e.GetTask)
	e.AggregationSources = m(e.AggregationSources)
}

// NewAPIRootIndexEndpoint returns an endpoint function that calls the method
// "ApiRootIndex" of service "stream".
func NewAPIRootIndexEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		deps := ctx.Value(dependencies.PublicRequestScopeCtxKey).(dependencies.PublicRequestScope)
		return nil, s.APIRootIndex(ctx, deps)
	}
}

// NewAPIVersionIndexEndpoint returns an endpoint function that calls the
// method "ApiVersionIndex" of service "stream".
func NewAPIVersionIndexEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		deps := ctx.Value(dependencies.PublicRequestScopeCtxKey).(dependencies.PublicRequestScope)
		return s.APIVersionIndex(ctx, deps)
	}
}

// NewHealthCheckEndpoint returns an endpoint function that calls the method
// "HealthCheck" of service "stream".
func NewHealthCheckEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		deps := ctx.Value(dependencies.PublicRequestScopeCtxKey).(dependencies.PublicRequestScope)
		return s.HealthCheck(ctx, deps)
	}
}

// NewCreateSourceEndpoint returns an endpoint function that calls the method
// "CreateSource" of service "stream".
func NewCreateSourceEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*CreateSourcePayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.BranchRequestScopeCtxKey).(dependencies.BranchRequestScope)
		return s.CreateSource(ctx, deps, p)
	}
}

// NewUpdateSourceEndpoint returns an endpoint function that calls the method
// "UpdateSource" of service "stream".
func NewUpdateSourceEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*UpdateSourcePayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SourceRequestScopeCtxKey).(dependencies.SourceRequestScope)
		return s.UpdateSource(ctx, deps, p)
	}
}

// NewListSourcesEndpoint returns an endpoint function that calls the method
// "ListSources" of service "stream".
func NewListSourcesEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*ListSourcesPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.BranchRequestScopeCtxKey).(dependencies.BranchRequestScope)
		return s.ListSources(ctx, deps, p)
	}
}

// NewListDeletedSourcesEndpoint returns an endpoint function that calls the
// method "ListDeletedSources" of service "stream".
func NewListDeletedSourcesEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*ListDeletedSourcesPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.BranchRequestScopeCtxKey).(dependencies.BranchRequestScope)
		return s.ListDeletedSources(ctx, deps, p)
	}
}

// NewGetSourceEndpoint returns an endpoint function that calls the method
// "GetSource" of service "stream".
func NewGetSourceEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*GetSourcePayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SourceRequestScopeCtxKey).(dependencies.SourceRequestScope)
		return s.GetSource(ctx, deps, p)
	}
}

// NewDeleteSourceEndpoint returns an endpoint function that calls the method
// "DeleteSource" of service "stream".
func NewDeleteSourceEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*DeleteSourcePayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SourceRequestScopeCtxKey).(dependencies.SourceRequestScope)
		return s.DeleteSource(ctx, deps, p)
	}
}

// NewGetSourceSettingsEndpoint returns an endpoint function that calls the
// method "GetSourceSettings" of service "stream".
func NewGetSourceSettingsEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*GetSourceSettingsPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SourceRequestScopeCtxKey).(dependencies.SourceRequestScope)
		return s.GetSourceSettings(ctx, deps, p)
	}
}

// NewUpdateSourceSettingsEndpoint returns an endpoint function that calls the
// method "UpdateSourceSettings" of service "stream".
func NewUpdateSourceSettingsEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*UpdateSourceSettingsPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SourceRequestScopeCtxKey).(dependencies.SourceRequestScope)
		return s.UpdateSourceSettings(ctx, deps, p)
	}
}

// NewTestSourceEndpoint returns an endpoint function that calls the method
// "TestSource" of service "stream".
func NewTestSourceEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		ep := req.(*TestSourceRequestData)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, ep.Payload.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SourceRequestScopeCtxKey).(dependencies.SourceRequestScope)
		return s.TestSource(ctx, deps, ep.Payload, ep.Body)
	}
}

// NewSourceStatisticsClearEndpoint returns an endpoint function that calls the
// method "SourceStatisticsClear" of service "stream".
func NewSourceStatisticsClearEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*SourceStatisticsClearPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SourceRequestScopeCtxKey).(dependencies.SourceRequestScope)
		return nil, s.SourceStatisticsClear(ctx, deps, p)
	}
}

// NewDisableSourceEndpoint returns an endpoint function that calls the method
// "DisableSource" of service "stream".
func NewDisableSourceEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*DisableSourcePayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SourceRequestScopeCtxKey).(dependencies.SourceRequestScope)
		return s.DisableSource(ctx, deps, p)
	}
}

// NewEnableSourceEndpoint returns an endpoint function that calls the method
// "EnableSource" of service "stream".
func NewEnableSourceEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*EnableSourcePayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SourceRequestScopeCtxKey).(dependencies.SourceRequestScope)
		return s.EnableSource(ctx, deps, p)
	}
}

// NewUndeleteSourceEndpoint returns an endpoint function that calls the method
// "UndeleteSource" of service "stream".
func NewUndeleteSourceEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*UndeleteSourcePayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SourceRequestScopeCtxKey).(dependencies.SourceRequestScope)
		return s.UndeleteSource(ctx, deps, p)
	}
}

// NewListSourceVersionsEndpoint returns an endpoint function that calls the
// method "ListSourceVersions" of service "stream".
func NewListSourceVersionsEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*ListSourceVersionsPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SourceRequestScopeCtxKey).(dependencies.SourceRequestScope)
		return s.ListSourceVersions(ctx, deps, p)
	}
}

// NewSourceVersionDetailEndpoint returns an endpoint function that calls the
// method "SourceVersionDetail" of service "stream".
func NewSourceVersionDetailEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*SourceVersionDetailPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SourceRequestScopeCtxKey).(dependencies.SourceRequestScope)
		return s.SourceVersionDetail(ctx, deps, p)
	}
}

// NewRollbackSourceVersionEndpoint returns an endpoint function that calls the
// method "RollbackSourceVersion" of service "stream".
func NewRollbackSourceVersionEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*RollbackSourceVersionPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SourceRequestScopeCtxKey).(dependencies.SourceRequestScope)
		return s.RollbackSourceVersion(ctx, deps, p)
	}
}

// NewCreateSinkEndpoint returns an endpoint function that calls the method
// "CreateSink" of service "stream".
func NewCreateSinkEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*CreateSinkPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SourceRequestScopeCtxKey).(dependencies.SourceRequestScope)
		return s.CreateSink(ctx, deps, p)
	}
}

// NewGetSinkEndpoint returns an endpoint function that calls the method
// "GetSink" of service "stream".
func NewGetSinkEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*GetSinkPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SinkRequestScopeCtxKey).(dependencies.SinkRequestScope)
		return s.GetSink(ctx, deps, p)
	}
}

// NewGetSinkSettingsEndpoint returns an endpoint function that calls the
// method "GetSinkSettings" of service "stream".
func NewGetSinkSettingsEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*GetSinkSettingsPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SinkRequestScopeCtxKey).(dependencies.SinkRequestScope)
		return s.GetSinkSettings(ctx, deps, p)
	}
}

// NewUpdateSinkSettingsEndpoint returns an endpoint function that calls the
// method "UpdateSinkSettings" of service "stream".
func NewUpdateSinkSettingsEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*UpdateSinkSettingsPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SinkRequestScopeCtxKey).(dependencies.SinkRequestScope)
		return s.UpdateSinkSettings(ctx, deps, p)
	}
}

// NewListSinksEndpoint returns an endpoint function that calls the method
// "ListSinks" of service "stream".
func NewListSinksEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*ListSinksPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SourceRequestScopeCtxKey).(dependencies.SourceRequestScope)
		return s.ListSinks(ctx, deps, p)
	}
}

// NewListDeletedSinksEndpoint returns an endpoint function that calls the
// method "ListDeletedSinks" of service "stream".
func NewListDeletedSinksEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*ListDeletedSinksPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SourceRequestScopeCtxKey).(dependencies.SourceRequestScope)
		return s.ListDeletedSinks(ctx, deps, p)
	}
}

// NewUpdateSinkEndpoint returns an endpoint function that calls the method
// "UpdateSink" of service "stream".
func NewUpdateSinkEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*UpdateSinkPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SinkRequestScopeCtxKey).(dependencies.SinkRequestScope)
		return s.UpdateSink(ctx, deps, p)
	}
}

// NewDeleteSinkEndpoint returns an endpoint function that calls the method
// "DeleteSink" of service "stream".
func NewDeleteSinkEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*DeleteSinkPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SinkRequestScopeCtxKey).(dependencies.SinkRequestScope)
		return s.DeleteSink(ctx, deps, p)
	}
}

// NewSinkStatisticsTotalEndpoint returns an endpoint function that calls the
// method "SinkStatisticsTotal" of service "stream".
func NewSinkStatisticsTotalEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*SinkStatisticsTotalPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SinkRequestScopeCtxKey).(dependencies.SinkRequestScope)
		return s.SinkStatisticsTotal(ctx, deps, p)
	}
}

// NewSinkStatisticsFilesEndpoint returns an endpoint function that calls the
// method "SinkStatisticsFiles" of service "stream".
func NewSinkStatisticsFilesEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*SinkStatisticsFilesPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SinkRequestScopeCtxKey).(dependencies.SinkRequestScope)
		return s.SinkStatisticsFiles(ctx, deps, p)
	}
}

// NewSinkStatisticsClearEndpoint returns an endpoint function that calls the
// method "SinkStatisticsClear" of service "stream".
func NewSinkStatisticsClearEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*SinkStatisticsClearPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SinkRequestScopeCtxKey).(dependencies.SinkRequestScope)
		return nil, s.SinkStatisticsClear(ctx, deps, p)
	}
}

// NewDisableSinkEndpoint returns an endpoint function that calls the method
// "DisableSink" of service "stream".
func NewDisableSinkEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*DisableSinkPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SinkRequestScopeCtxKey).(dependencies.SinkRequestScope)
		return s.DisableSink(ctx, deps, p)
	}
}

// NewEnableSinkEndpoint returns an endpoint function that calls the method
// "EnableSink" of service "stream".
func NewEnableSinkEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*EnableSinkPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SinkRequestScopeCtxKey).(dependencies.SinkRequestScope)
		return s.EnableSink(ctx, deps, p)
	}
}

// NewUndeleteSinkEndpoint returns an endpoint function that calls the method
// "UndeleteSink" of service "stream".
func NewUndeleteSinkEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*UndeleteSinkPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SinkRequestScopeCtxKey).(dependencies.SinkRequestScope)
		return s.UndeleteSink(ctx, deps, p)
	}
}

// NewListSinkVersionsEndpoint returns an endpoint function that calls the
// method "ListSinkVersions" of service "stream".
func NewListSinkVersionsEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*ListSinkVersionsPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SinkRequestScopeCtxKey).(dependencies.SinkRequestScope)
		return s.ListSinkVersions(ctx, deps, p)
	}
}

// NewSinkVersionDetailEndpoint returns an endpoint function that calls the
// method "SinkVersionDetail" of service "stream".
func NewSinkVersionDetailEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*SinkVersionDetailPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SinkRequestScopeCtxKey).(dependencies.SinkRequestScope)
		return s.SinkVersionDetail(ctx, deps, p)
	}
}

// NewRollbackSinkVersionEndpoint returns an endpoint function that calls the
// method "RollbackSinkVersion" of service "stream".
func NewRollbackSinkVersionEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*RollbackSinkVersionPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.SinkRequestScopeCtxKey).(dependencies.SinkRequestScope)
		return s.RollbackSinkVersion(ctx, deps, p)
	}
}

// NewGetTaskEndpoint returns an endpoint function that calls the method
// "GetTask" of service "stream".
func NewGetTaskEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*GetTaskPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.ProjectRequestScopeCtxKey).(dependencies.ProjectRequestScope)
		return s.GetTask(ctx, deps, p)
	}
}

// NewAggregationSourcesEndpoint returns an endpoint function that calls the
// method "AggregationSources" of service "stream".
func NewAggregationSourcesEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*AggregationSourcesPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "storage-api-token",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.StorageAPIToken, &sc)
		if err != nil {
			return nil, err
		}
		deps := ctx.Value(dependencies.BranchRequestScopeCtxKey).(dependencies.BranchRequestScope)
		return s.AggregationSources(ctx, deps, p)
	}
}
