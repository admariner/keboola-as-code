// Code generated by goa v3.5.5, DO NOT EDIT.
//
// templates HTTP server
//
// Command:
// $ goa gen github.com/keboola/keboola-as-code/design --output
// ./internal/pkg/template/api

package server

import (
	"context"
	"net/http"

	templates "github.com/keboola/keboola-as-code/internal/pkg/template/api/gen/templates"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Server lists the templates service endpoint HTTP handlers.
type Server struct {
	Mounts             []*MountPoint
	IndexEndpoint      http.Handler
	HealthCheck        http.Handler
	GenHTTPOpenapiJSON http.Handler
}

// ErrorNamer is an interface implemented by generated error structs that
// exposes the name of the error as defined in the design.
type ErrorNamer interface {
	ErrorName() string
}

// MountPoint holds information about the mounted endpoints.
type MountPoint struct {
	// Method is the name of the service method served by the mounted HTTP handler.
	Method string
	// Verb is the HTTP method used to match requests to the mounted handler.
	Verb string
	// Pattern is the HTTP request path pattern used to match requests to the
	// mounted handler.
	Pattern string
}

// New instantiates HTTP handlers for all the templates service endpoints using
// the provided encoder and decoder. The handlers are mounted on the given mux
// using the HTTP verb and path defined in the design. errhandler is called
// whenever a response fails to be encoded. formatter is used to format errors
// returned by the service methods prior to encoding. Both errhandler and
// formatter are optional and can be nil.
func New(
	e *templates.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
	fileSystemGenHTTPOpenapiJSON http.FileSystem,
) *Server {
	if fileSystemGenHTTPOpenapiJSON == nil {
		fileSystemGenHTTPOpenapiJSON = http.Dir(".")
	}
	return &Server{
		Mounts: []*MountPoint{
			{"IndexEndpoint", "GET", "/"},
			{"HealthCheck", "GET", "/health-check"},
			{"./gen/http/openapi.json", "GET", "/openapi.json"},
		},
		IndexEndpoint:      NewIndexEndpointHandler(e.IndexEndpoint, mux, decoder, encoder, errhandler, formatter),
		HealthCheck:        NewHealthCheckHandler(e.HealthCheck, mux, decoder, encoder, errhandler, formatter),
		GenHTTPOpenapiJSON: http.FileServer(fileSystemGenHTTPOpenapiJSON),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "templates" }

// Use wraps the server handlers with the given middleware.
func (s *Server) Use(m func(http.Handler) http.Handler) {
	s.IndexEndpoint = m(s.IndexEndpoint)
	s.HealthCheck = m(s.HealthCheck)
}

// Mount configures the mux to serve the templates endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountIndexEndpointHandler(mux, h.IndexEndpoint)
	MountHealthCheckHandler(mux, h.HealthCheck)
	MountGenHTTPOpenapiJSON(mux, goahttp.Replace("", "/./gen/http/openapi.json", h.GenHTTPOpenapiJSON))
}

// Mount configures the mux to serve the templates endpoints.
func (s *Server) Mount(mux goahttp.Muxer) {
	Mount(mux, s)
}

// MountIndexEndpointHandler configures the mux to serve the "templates"
// service "index" endpoint.
func MountIndexEndpointHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/", f)
}

// NewIndexEndpointHandler creates a HTTP handler which loads the HTTP request
// and calls the "templates" service "index" endpoint.
func NewIndexEndpointHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		encodeResponse = EncodeIndexEndpointResponse(encoder)
		encodeError    = goahttp.ErrorEncoder(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "index")
		ctx = context.WithValue(ctx, goa.ServiceKey, "templates")
		var err error
		res, err := endpoint(ctx, nil)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountHealthCheckHandler configures the mux to serve the "templates" service
// "health-check" endpoint.
func MountHealthCheckHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/health-check", f)
}

// NewHealthCheckHandler creates a HTTP handler which loads the HTTP request
// and calls the "templates" service "health-check" endpoint.
func NewHealthCheckHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		encodeResponse = EncodeHealthCheckResponse(encoder)
		encodeError    = goahttp.ErrorEncoder(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "health-check")
		ctx = context.WithValue(ctx, goa.ServiceKey, "templates")
		var err error
		res, err := endpoint(ctx, nil)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountGenHTTPOpenapiJSON configures the mux to serve GET request made to
// "/openapi.json".
func MountGenHTTPOpenapiJSON(mux goahttp.Muxer, h http.Handler) {
	mux.Handle("GET", "/openapi.json", h.ServeHTTP)
}
