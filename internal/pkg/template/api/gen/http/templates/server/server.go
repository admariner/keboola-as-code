// Code generated by goa v3.5.5, DO NOT EDIT.
//
// templates HTTP server
//
// Command:
// $ goa gen github.com/keboola/keboola-as-code/api/templates --output
// ./internal/pkg/template/api

package server

import (
	"context"
	"net/http"

	templates "github.com/keboola/keboola-as-code/internal/pkg/template/api/gen/templates"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
	"goa.design/plugins/v3/cors"
)

// Server lists the templates service endpoint HTTP handlers.
type Server struct {
	Mounts          []*MountPoint
	IndexRoot       http.Handler
	HealthCheck     http.Handler
	IndexEndpoint   http.Handler
	Foo             http.Handler
	CORS            http.Handler
	GenOpenapiJSON  http.Handler
	GenOpenapiYaml  http.Handler
	GenOpenapi3JSON http.Handler
	GenOpenapi3Yaml http.Handler
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
	fileSystemGenOpenapiJSON http.FileSystem,
	fileSystemGenOpenapiYaml http.FileSystem,
	fileSystemGenOpenapi3JSON http.FileSystem,
	fileSystemGenOpenapi3Yaml http.FileSystem,
) *Server {
	if fileSystemGenOpenapiJSON == nil {
		fileSystemGenOpenapiJSON = http.Dir(".")
	}
	if fileSystemGenOpenapiYaml == nil {
		fileSystemGenOpenapiYaml = http.Dir(".")
	}
	if fileSystemGenOpenapi3JSON == nil {
		fileSystemGenOpenapi3JSON = http.Dir(".")
	}
	if fileSystemGenOpenapi3Yaml == nil {
		fileSystemGenOpenapi3Yaml = http.Dir(".")
	}
	return &Server{
		Mounts: []*MountPoint{
			{"IndexRoot", "GET", "/"},
			{"HealthCheck", "GET", "/health-check"},
			{"IndexEndpoint", "GET", "/v1"},
			{"Foo", "GET", "/v1/foo"},
			{"CORS", "OPTIONS", "/"},
			{"CORS", "OPTIONS", "/health-check"},
			{"CORS", "OPTIONS", "/v1"},
			{"CORS", "OPTIONS", "/v1/foo"},
			{"CORS", "OPTIONS", "/v1/documentation/openapi.json"},
			{"CORS", "OPTIONS", "/v1/documentation/openapi.yaml"},
			{"CORS", "OPTIONS", "/v1/documentation/openapi3.json"},
			{"CORS", "OPTIONS", "/v1/documentation/openapi3.yaml"},
			{"gen/openapi.json", "GET", "/v1/documentation/openapi.json"},
			{"gen/openapi.yaml", "GET", "/v1/documentation/openapi.yaml"},
			{"gen/openapi3.json", "GET", "/v1/documentation/openapi3.json"},
			{"gen/openapi3.yaml", "GET", "/v1/documentation/openapi3.yaml"},
		},
		IndexRoot:       NewIndexRootHandler(e.IndexRoot, mux, decoder, encoder, errhandler, formatter),
		HealthCheck:     NewHealthCheckHandler(e.HealthCheck, mux, decoder, encoder, errhandler, formatter),
		IndexEndpoint:   NewIndexEndpointHandler(e.IndexEndpoint, mux, decoder, encoder, errhandler, formatter),
		Foo:             NewFooHandler(e.Foo, mux, decoder, encoder, errhandler, formatter),
		CORS:            NewCORSHandler(),
		GenOpenapiJSON:  http.FileServer(fileSystemGenOpenapiJSON),
		GenOpenapiYaml:  http.FileServer(fileSystemGenOpenapiYaml),
		GenOpenapi3JSON: http.FileServer(fileSystemGenOpenapi3JSON),
		GenOpenapi3Yaml: http.FileServer(fileSystemGenOpenapi3Yaml),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "templates" }

// Use wraps the server handlers with the given middleware.
func (s *Server) Use(m func(http.Handler) http.Handler) {
	s.IndexRoot = m(s.IndexRoot)
	s.HealthCheck = m(s.HealthCheck)
	s.IndexEndpoint = m(s.IndexEndpoint)
	s.Foo = m(s.Foo)
	s.CORS = m(s.CORS)
}

// Mount configures the mux to serve the templates endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountIndexRootHandler(mux, h.IndexRoot)
	MountHealthCheckHandler(mux, h.HealthCheck)
	MountIndexEndpointHandler(mux, h.IndexEndpoint)
	MountFooHandler(mux, h.Foo)
	MountCORSHandler(mux, h.CORS)
	MountGenOpenapiJSON(mux, goahttp.Replace("", "/gen/openapi.json", h.GenOpenapiJSON))
	MountGenOpenapiYaml(mux, goahttp.Replace("", "/gen/openapi.yaml", h.GenOpenapiYaml))
	MountGenOpenapi3JSON(mux, goahttp.Replace("", "/gen/openapi3.json", h.GenOpenapi3JSON))
	MountGenOpenapi3Yaml(mux, goahttp.Replace("", "/gen/openapi3.yaml", h.GenOpenapi3Yaml))
}

// Mount configures the mux to serve the templates endpoints.
func (s *Server) Mount(mux goahttp.Muxer) {
	Mount(mux, s)
}

// MountIndexRootHandler configures the mux to serve the "templates" service
// "index-root" endpoint.
func MountIndexRootHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := HandleTemplatesOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/", f)
}

// NewIndexRootHandler creates a HTTP handler which loads the HTTP request and
// calls the "templates" service "index-root" endpoint.
func NewIndexRootHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "index-root")
		ctx = context.WithValue(ctx, goa.ServiceKey, "templates")
		http.Redirect(w, r, "/v1", http.StatusMovedPermanently)
	})
}

// MountHealthCheckHandler configures the mux to serve the "templates" service
// "health-check" endpoint.
func MountHealthCheckHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := HandleTemplatesOrigin(h).(http.HandlerFunc)
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

// MountIndexEndpointHandler configures the mux to serve the "templates"
// service "index" endpoint.
func MountIndexEndpointHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := HandleTemplatesOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/v1", f)
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

// MountFooHandler configures the mux to serve the "templates" service "foo"
// endpoint.
func MountFooHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := HandleTemplatesOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/v1/foo", f)
}

// NewFooHandler creates a HTTP handler which loads the HTTP request and calls
// the "templates" service "foo" endpoint.
func NewFooHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeFooRequest(mux, decoder)
		encodeResponse = EncodeFooResponse(encoder)
		encodeError    = goahttp.ErrorEncoder(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "foo")
		ctx = context.WithValue(ctx, goa.ServiceKey, "templates")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
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

// MountGenOpenapiJSON configures the mux to serve GET request made to
// "/v1/documentation/openapi.json".
func MountGenOpenapiJSON(mux goahttp.Muxer, h http.Handler) {
	mux.Handle("GET", "/v1/documentation/openapi.json", HandleTemplatesOrigin(h).ServeHTTP)
}

// MountGenOpenapiYaml configures the mux to serve GET request made to
// "/v1/documentation/openapi.yaml".
func MountGenOpenapiYaml(mux goahttp.Muxer, h http.Handler) {
	mux.Handle("GET", "/v1/documentation/openapi.yaml", HandleTemplatesOrigin(h).ServeHTTP)
}

// MountGenOpenapi3JSON configures the mux to serve GET request made to
// "/v1/documentation/openapi3.json".
func MountGenOpenapi3JSON(mux goahttp.Muxer, h http.Handler) {
	mux.Handle("GET", "/v1/documentation/openapi3.json", HandleTemplatesOrigin(h).ServeHTTP)
}

// MountGenOpenapi3Yaml configures the mux to serve GET request made to
// "/v1/documentation/openapi3.yaml".
func MountGenOpenapi3Yaml(mux goahttp.Muxer, h http.Handler) {
	mux.Handle("GET", "/v1/documentation/openapi3.yaml", HandleTemplatesOrigin(h).ServeHTTP)
}

// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service templates.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = HandleTemplatesOrigin(h)
	mux.Handle("OPTIONS", "/", h.ServeHTTP)
	mux.Handle("OPTIONS", "/health-check", h.ServeHTTP)
	mux.Handle("OPTIONS", "/v1", h.ServeHTTP)
	mux.Handle("OPTIONS", "/v1/foo", h.ServeHTTP)
	mux.Handle("OPTIONS", "/v1/documentation/openapi.json", h.ServeHTTP)
	mux.Handle("OPTIONS", "/v1/documentation/openapi.yaml", h.ServeHTTP)
	mux.Handle("OPTIONS", "/v1/documentation/openapi3.json", h.ServeHTTP)
	mux.Handle("OPTIONS", "/v1/documentation/openapi3.yaml", h.ServeHTTP)
}

// NewCORSHandler creates a HTTP handler which returns a simple 200 response.
func NewCORSHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
}

// HandleTemplatesOrigin applies the CORS response headers corresponding to the
// origin for the service templates.
func HandleTemplatesOrigin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			h.ServeHTTP(w, r)
			return
		}
		if cors.MatchOrigin(origin, "*") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				w.Header().Set("Access-Control-Allow-Headers", "X-StorageApi-Token")
			}
			h.ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, r)
		return
	})
}
