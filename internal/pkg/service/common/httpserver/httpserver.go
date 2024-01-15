package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/httpserver/middleware"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/servicectx"
	"github.com/keboola/keboola-as-code/internal/pkg/telemetry"
)

const (
	requestTimeout          = 30 * time.Second
	readHeaderTimeout       = 10 * time.Second
	gracefulShutdownTimeout = 30 * time.Second
)

type dependencies interface {
	Logger() log.Logger
	Process() *servicectx.Process
	Telemetry() telemetry.Telemetry
}

// Start HTTP server.
func Start(ctx context.Context, d dependencies, cfg Config) error {
	logger, tel := d.Logger(), d.Telemetry()

	// Create server components
	com := newComponents(cfg, logger)

	// Register middlewares
	middlewareCfg := middleware.NewConfig(cfg.MiddlewareOptions...)
	com.Muxer.Use(middleware.OpenTelemetryExtractRoute())
	handler := middleware.Wrap(
		com.Muxer,
		middleware.ContextTimout(requestTimeout),
		middleware.RequestInfo(),
		middleware.Filter(middlewareCfg),
		middleware.Logger(logger),
		middleware.OpenTelemetry(tel.TracerProvider(), tel.MeterProvider(), middlewareCfg),
	)
	// Mount endpoints
	cfg.Mount(com)
	logger.InfofCtx(ctx, "mounted HTTP endpoints")

	// Start HTTP server
	srv := &http.Server{Addr: cfg.ListenAddress, Handler: handler, ReadHeaderTimeout: readHeaderTimeout}
	proc := d.Process()
	proc.Add(func(shutdown servicectx.ShutdownFn) {
		// Start HTTP server in a separate goroutine.
		logger.InfofCtx(ctx, "HTTP server listening on %q", cfg.ListenAddress)
		serverErr := srv.ListenAndServe()         // ListenAndServe blocks while the server is running
		shutdown(context.Background(), serverErr) // nolint: contextcheck // intentionally creating new context for the shutdown operation
	})

	// Register graceful shutdown
	proc.OnShutdown(func(ctx context.Context) {
		// Shutdown gracefully with a timeout.
		ctx, cancel := context.WithTimeout(ctx, gracefulShutdownTimeout)
		defer cancel()

		logger.InfofCtx(ctx, "shutting down HTTP server at %q", cfg.ListenAddress)

		if err := srv.Shutdown(ctx); err != nil {
			logger.ErrorfCtx(ctx, `HTTP server shutdown error: %s`, err)
		}
		logger.Info(ctx, "HTTP server shutdown finished")
	})

	return nil
}
