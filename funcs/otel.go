package funcs

import (
	"github.com/tsingsun/woocoo/pkg/conf"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func ZipkinTracer(cnf *conf.Configuration) (opts []sdktrace.TracerProviderOption) {
	exporter, err := zipkin.New(cnf.String("traceExporterEndpoint"))
	if err != nil {
		panic(err)
	}
	return append(opts, sdktrace.WithBatcher(exporter))
}

func PrometheusProvider(cnf *conf.Configuration) (opts []metric.Option) {
	exporter, err := prometheus.New()
	if err != nil {
		panic(err)
	}
	return []metric.Option{
		metric.WithReader(exporter),
	}
}
