// Package prometheus provides HTTP metrics endpoint with OpenTelemetry metrics for Prometheus scraper.
package prometheus

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/bridge/opencensus"
	export "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/servicectx"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/netutils"
)

const (
	Endpoint        = "metrics"
	startTimeout    = 10 * time.Second
	shutdownTimeout = 30 * time.Second
)

type errLogger struct {
	logger log.Logger
}

func (l *errLogger) Println(v ...any) {
	l.logger.Error(v...)
}

// ServeMetrics starts HTTP server for Prometheus metrics and return OpenTelemetry metrics provider.
// Inspired by: https://github.com/open-telemetry/opentelemetry-go/blob/main/example/prometheus/main.go
func ServeMetrics(ctx context.Context, serviceName, listenAddr string, logger log.Logger, proc *servicectx.Process) (*metric.MeterProvider, error) {
	logger = logger.AddPrefix("[metrics]")

	// Create resource
	res, err := resource.New(
		ctx,
		resource.WithAttributes(attribute.String("service_name", serviceName)),
		// resource.WithFromEnv(), // unused
		// resource.WithTelemetrySDK(), // unused
	)
	if err != nil {
		return nil, err
	}

	// Create metrics registry and exporter
	registry := prometheus.NewRegistry()
	exporter, err := export.New(export.WithRegisterer(registry), export.WithoutScopeInfo(), export.WithoutUnits())
	if err != nil {
		return nil, err
	}

	// Register legacy OpenCensus metrics, for go-cloud (https://github.com/google/go-cloud/issues/2877)
	exporter.RegisterProducer(opencensus.NewMetricProducer())

	// Create HTTP metrics server
	opts := promhttp.HandlerOpts{ErrorLog: &errLogger{logger: logger}}
	handler := http.NewServeMux()
	handler.Handle("/"+Endpoint, promhttp.HandlerFor(registry, opts))
	srv := &http.Server{Addr: listenAddr, Handler: handler, ReadHeaderTimeout: 10 * time.Second}
	proc.Add(func(ctx context.Context, shutdown servicectx.ShutdownFn) {
		logger.Infof(`HTTP server listening on "%s/%s"`, listenAddr, Endpoint)
		shutdown(srv.ListenAndServe())
	})
	proc.OnShutdown(func() {
		logger.Infof(`shutting down HTTP server at "%s"`, listenAddr)

		// Shutdown gracefully with a 30s timeout.
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Errorf(`HTTP server shutdown error: %s`, err)
		}
		logger.Info("HTTP server shutdown finished")
	})

	// Wait for server
	if err := netutils.WaitForTCP(srv.Addr, startTimeout); err != nil {
		return nil, errors.Errorf(`metrics server did not start: %w`, err)
	}

	// Create OpenTelemetry metrics provider
	provider := metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithResource(res),
		metric.WithView(View()),
	)
	return provider, nil
}

func View() metric.View {
	ignoreAttrs := metric.NewView(
		metric.Instrument{Name: "*"},
		metric.Stream{AttributeFilter: func(value attribute.KeyValue) bool {
			switch value.Key {
			// Remove invalid otelhttp metric attributes with high cardinality.
			// https://github.com/open-telemetry/opentelemetry-go-contrib/issues/3765
			case "net.sock.peer.addr",
				"net.sock.peer.port",
				"http.user_agent",
				"http.client_ip",
				"http.request_content_length",
				"http.response_content_length":
				return false
			// Remove unused attributes.
			case "http.flavor":
				return false
			}
			return true
		}},
	)
	rename := func(inst metric.Instrument) metric.Instrument {
		if strings.HasPrefix(inst.Name, "http.server") {
			inst.Name = "keboola.go." + inst.Name
		}
		return inst
	}
	return func(inst metric.Instrument) (metric.Stream, bool) {
		inst = rename(inst)
		return ignoreAttrs(inst)
	}
}
