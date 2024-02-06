// Code generated by goa v3.14.6, DO NOT EDIT.
//
// stream endpoints
//
// Command:
// $ goa gen github.com/keboola/keboola-as-code/api/buffer --output
// ./internal/pkg/service/buffer/api

package stream

import (
	"context"

	dependencies "github.com/keboola/keboola-as-code/internal/pkg/service/buffer/dependencies"
	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Endpoints wraps the "stream" service endpoints.
type Endpoints struct {
	APIRootIndex         goa.Endpoint
	APIVersionIndex      goa.Endpoint
	HealthCheck          goa.Endpoint
	CreateSource         goa.Endpoint
	UpdateSource         goa.Endpoint
	ListSources          goa.Endpoint
	GetSource            goa.Endpoint
	DeleteSource         goa.Endpoint
	GetSourceSettings    goa.Endpoint
	UpdateSourceSettings goa.Endpoint
	RefreshSourceTokens  goa.Endpoint
	CreateSink           goa.Endpoint
	GetSink              goa.Endpoint
	GetSinkSettings      goa.Endpoint
	UpdateSinkSettings   goa.Endpoint
	ListSinks            goa.Endpoint
	UpdateSink           goa.Endpoint
	DeleteSink           goa.Endpoint
	GetTask              goa.Endpoint
}

// NewEndpoints wraps the methods of the "stream" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		APIRootIndex:         NewAPIRootIndexEndpoint(s),
		APIVersionIndex:      NewAPIVersionIndexEndpoint(s),
		HealthCheck:          NewHealthCheckEndpoint(s),
		CreateSource:         NewCreateSourceEndpoint(s, a.APIKeyAuth),
		UpdateSource:         NewUpdateSourceEndpoint(s, a.APIKeyAuth),
		ListSources:          NewListSourcesEndpoint(s, a.APIKeyAuth),
		GetSource:            NewGetSourceEndpoint(s, a.APIKeyAuth),
		DeleteSource:         NewDeleteSourceEndpoint(s, a.APIKeyAuth),
		GetSourceSettings:    NewGetSourceSettingsEndpoint(s, a.APIKeyAuth),
		UpdateSourceSettings: NewUpdateSourceSettingsEndpoint(s, a.APIKeyAuth),
		RefreshSourceTokens:  NewRefreshSourceTokensEndpoint(s, a.APIKeyAuth),
		CreateSink:           NewCreateSinkEndpoint(s, a.APIKeyAuth),
		GetSink:              NewGetSinkEndpoint(s, a.APIKeyAuth),
		GetSinkSettings:      NewGetSinkSettingsEndpoint(s, a.APIKeyAuth),
		UpdateSinkSettings:   NewUpdateSinkSettingsEndpoint(s, a.APIKeyAuth),
		ListSinks:            NewListSinksEndpoint(s, a.APIKeyAuth),
		UpdateSink:           NewUpdateSinkEndpoint(s, a.APIKeyAuth),
		DeleteSink:           NewDeleteSinkEndpoint(s, a.APIKeyAuth),
		GetTask:              NewGetTaskEndpoint(s, a.APIKeyAuth),
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
	e.GetSource = m(e.GetSource)
	e.DeleteSource = m(e.DeleteSource)
	e.GetSourceSettings = m(e.GetSourceSettings)
	e.UpdateSourceSettings = m(e.UpdateSourceSettings)
	e.RefreshSourceTokens = m(e.RefreshSourceTokens)
	e.CreateSink = m(e.CreateSink)
	e.GetSink = m(e.GetSink)
	e.GetSinkSettings = m(e.GetSinkSettings)
	e.UpdateSinkSettings = m(e.UpdateSinkSettings)
	e.ListSinks = m(e.ListSinks)
	e.UpdateSink = m(e.UpdateSink)
	e.DeleteSink = m(e.DeleteSink)
	e.GetTask = m(e.GetTask)
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
		deps := ctx.Value(dependencies.ProjectRequestScopeCtxKey).(dependencies.ProjectRequestScope)
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
		deps := ctx.Value(dependencies.ProjectRequestScopeCtxKey).(dependencies.ProjectRequestScope)
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
		deps := ctx.Value(dependencies.ProjectRequestScopeCtxKey).(dependencies.ProjectRequestScope)
		return s.ListSources(ctx, deps, p)
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
		deps := ctx.Value(dependencies.ProjectRequestScopeCtxKey).(dependencies.ProjectRequestScope)
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
		deps := ctx.Value(dependencies.ProjectRequestScopeCtxKey).(dependencies.ProjectRequestScope)
		return nil, s.DeleteSource(ctx, deps, p)
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
		deps := ctx.Value(dependencies.ProjectRequestScopeCtxKey).(dependencies.ProjectRequestScope)
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
		deps := ctx.Value(dependencies.ProjectRequestScopeCtxKey).(dependencies.ProjectRequestScope)
		return s.UpdateSourceSettings(ctx, deps, p)
	}
}

// NewRefreshSourceTokensEndpoint returns an endpoint function that calls the
// method "RefreshSourceTokens" of service "stream".
func NewRefreshSourceTokensEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*RefreshSourceTokensPayload)
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
		return s.RefreshSourceTokens(ctx, deps, p)
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
		deps := ctx.Value(dependencies.ProjectRequestScopeCtxKey).(dependencies.ProjectRequestScope)
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
		deps := ctx.Value(dependencies.ProjectRequestScopeCtxKey).(dependencies.ProjectRequestScope)
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
		deps := ctx.Value(dependencies.ProjectRequestScopeCtxKey).(dependencies.ProjectRequestScope)
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
		deps := ctx.Value(dependencies.ProjectRequestScopeCtxKey).(dependencies.ProjectRequestScope)
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
		deps := ctx.Value(dependencies.ProjectRequestScopeCtxKey).(dependencies.ProjectRequestScope)
		return s.ListSinks(ctx, deps, p)
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
		deps := ctx.Value(dependencies.ProjectRequestScopeCtxKey).(dependencies.ProjectRequestScope)
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
		deps := ctx.Value(dependencies.ProjectRequestScopeCtxKey).(dependencies.ProjectRequestScope)
		return nil, s.DeleteSink(ctx, deps, p)
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
