package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	goaHTTP "goa.design/goa/v3/http"
	httpMdlwr "goa.design/goa/v3/http/middleware"
	"goa.design/goa/v3/middleware"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"

	openapiapi "github.com/keboola/keboola-as-code/api/templates"
	templatesSvr "github.com/keboola/keboola-as-code/internal/pkg/template/api/gen/http/templates/server"
	"github.com/keboola/keboola-as-code/internal/pkg/template/api/gen/templates"
)

// handleHTTPServer starts configures and starts a HTTP server on the given
// URL. It shuts down the server if any error is received in the error channel.
func handleHTTPServer(ctx context.Context, u *url.URL, templatesEndpoints *templates.Endpoints, wg *sync.WaitGroup, errc chan error, logger *log.Logger, debug bool) {
	// Setup goa log adapter.
	var (
		adapter middleware.Logger
	)
	{
		adapter = middleware.NewLogger(logger)
	}

	// Provide the transport specific request decoder and response encoder.
	// The goa http package has built-in support for JSON, XML and gob.
	// Other encodings can be used by providing the corresponding functions,
	// see goa.design/implement/encoding.
	var (
		dec = goaHTTP.RequestDecoder
		enc = goaHTTP.ResponseEncoder
	)

	// Build the service HTTP request multiplexer and configure it to serve
	// HTTP requests to the service endpoints.
	var mux goaHTTP.Muxer
	{
		mux = goaHTTP.NewMuxer()
	}

	// Wrap the endpoints with the transport specific layers. The generated
	// server packages contains code generated from the design which maps
	// the service input and output data structures to HTTP requests and
	// responses.
	var (
		templatesServer *templatesSvr.Server
	)
	{
		eh := errorHandler(logger)
		docsFS := http.FS(openapiapi.ApiDocsFS)
		templatesServer = templatesSvr.New(templatesEndpoints, mux, dec, enc, eh, nil, docsFS, docsFS, docsFS, docsFS)
		if debug {
			servers := goaHTTP.Servers{
				templatesServer,
			}
			servers.Use(httpMdlwr.Debug(mux, os.Stdout))
		}
	}
	// Configure the mux.
	templatesSvr.Mount(mux, templatesServer)

	// Wrap the multiplexer with additional middlewares. Middlewares mounted
	// here apply to all the service endpoints.
	var handler http.Handler = mux
	{
		handler = httpMdlwr.Log(adapter)(handler)
		handler = httpMdlwr.RequestID()(handler)
		handler = httptrace.WrapHandler(handler, "templates-api", "")
	}

	// Start HTTP server using default configuration, change the code to
	// configure the server as required by your service.
	srv := &http.Server{Addr: u.Host, Handler: handler}
	for _, m := range templatesServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}

	(*wg).Add(1)
	go func() {
		defer (*wg).Done()

		// Start HTTP server in a separate goroutine.
		go func() {
			logger.Printf("HTTP server listening on %q", u.Host)
			errc <- srv.ListenAndServe()
		}()

		<-ctx.Done()
		logger.Printf("shutting down HTTP server at %q", u.Host)

		// Shutdown gracefully with a 30s timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()
}

// errorHandler returns a function that writes and logs the given error.
// The function also writes and logs the error unique ID so that it's possible
// to correlate.
func errorHandler(logger *log.Logger) func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		id := ctx.Value(middleware.RequestIDKey).(string)
		_, _ = w.Write([]byte("[" + id + "] encoding: " + err.Error()))
		logger.Printf("[%s] ERROR: %s", id, err.Error())
	}
}