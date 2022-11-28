package http

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"time"

	goaHTTP "goa.design/goa/v3/http"
	httpMiddleware "goa.design/goa/v3/http/middleware"
	dataDog "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"

	"github.com/keboola/keboola-as-code/internal/pkg/service/buffer/api/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/buffer/api/gen/buffer"
	bufferSvr "github.com/keboola/keboola-as-code/internal/pkg/service/buffer/api/gen/http/buffer/server"
	"github.com/keboola/keboola-as-code/internal/pkg/service/buffer/api/openapi"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/httpserver"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/muxer"
	swaggerui "github.com/keboola/keboola-as-code/third_party"
)

const (
	ErrorNamePrefix   = "buffer."
	ExceptionIdPrefix = "keboola-buffer-"
)

// HandleHTTPServer starts configures and starts a HTTP server on the given
// URL. It shuts down the server if any error is received in the error channel.
func HandleHTTPServer(ctx context.Context, d dependencies.ForServer, u *url.URL, endpoints *buffer.Endpoints, errCh chan error, debug bool) {
	logger := d.Logger()

	// Trace endpoint start, finish and error
	endpoints.Use(TraceEndpointsMiddleware(d))

	// Build the service HTTP request multiplexer and configure it to serve
	// HTTP requests to the service endpoints.
	errorWr := httpserver.NewErrorWriter(logger, ErrorNamePrefix, ExceptionIdPrefix)
	errorFmt := httpserver.FormatError
	encoder := httpserver.NewEncoder(logger, errorWr)
	decoder := httpserver.NewDecoder()
	mux := muxer.New(errorWr)

	// Wrap the endpoints with the transport specific layers. The generated
	// server packages contains code generated from the design which maps
	// the service input and output data structures to HTTP requests and
	// responses.
	docsFs := http.FS(openapi.Fs)
	swaggerUiFs := http.FS(swaggerui.SwaggerFS)
	server := bufferSvr.New(endpoints, mux, decoder, encoder, errorWr.Write, errorFmt, docsFs, docsFs, docsFs, docsFs, swaggerUiFs)
	if debug {
		servers := goaHTTP.Servers{server}
		servers.Use(httpMiddleware.Debug(mux, os.Stdout))
	}

	// Configure the mux.
	bufferSvr.Mount(mux, server)

	// Wrap the multiplexer with additional middlewares. Middlewares mounted
	// here apply to all the service endpoints.
	var handler http.Handler = mux
	handler = LogMiddleware(d, handler)
	handler = ContextMiddleware(d, handler)
	handler = dataDog.WrapHandler(handler, "buffer-api", "", dataDog.WithIgnoreRequest(func(r *http.Request) bool {
		// Trace all requests except health check
		return r.URL.Path == "/health-check"
	}))

	// Start HTTP server using default configuration, change the code to
	// configure the server as required by your service.
	srv := &http.Server{
		Addr:              u.Host,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}
	for _, m := range server.Mounts {
		logger.Infof("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}

	wg := d.ServerWaitGroup()
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Start HTTP server in a separate goroutine.
		go func() {
			logger.Infof("HTTP server listening on %q", u.Host)
			errCh <- srv.ListenAndServe()
		}()

		<-ctx.Done()
		logger.Infof("shutting down HTTP server at %q", u.Host)

		// Shutdown gracefully with a 30s timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()
}