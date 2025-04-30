package otelprovider

import (
	"context"
	"os"
	"strings"

	"github.com/spechtlabs/go-otel-utils/otelzap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

type Tracer struct {
	providerOptions []trace.TracerProviderOption
	insecure        bool
	resources       *resource.Resource
	register        bool
}

func NewTracer(opts ...TracerOption) *trace.TracerProvider {
	t := &Tracer{
		insecure:        false,
		providerOptions: []trace.TracerProviderOption{},
		resources:       newOtelResources(),
		register:        true,
	}

	for _, opt := range opts {
		opt(t)
	}

	t.providerOptions = append(t.providerOptions, trace.WithResource(t.resources))
	traceProvider := trace.NewTracerProvider(t.providerOptions...)

	// Register the Provider globally
	if t.register {
		otel.SetTracerProvider(traceProvider)
	}

	return traceProvider
}

// TracerOption applies a configuration to the given config.
type TracerOption func(t *Tracer)

func WithTraceInsecure() TracerOption {
	return func(t *Tracer) {
		t.insecure = true
	}
}

func WithGrpcTraceEndpoint(otelGrpcEndpoint string) TracerOption {
	return func(t *Tracer) {
		grpcExporterOptions := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(otelGrpcEndpoint),
		}

		if t.insecure {
			grpcExporterOptions = append(grpcExporterOptions, otlptracegrpc.WithInsecure())
		}

		grpcExporter, err := otlptrace.New(context.Background(), otlptracegrpc.NewClient(grpcExporterOptions...))
		if err != nil {
			otelzap.L().Sugar().Fatalw("Failed to create OTLP gRPC trace exporter", zap.Error(err))
		}

		batcher := trace.NewBatchSpanProcessor(grpcExporter)

		t.providerOptions = append(t.providerOptions, trace.WithSpanProcessor(batcher))
	}
}

func WithHttpTraceEndpoint(otelHttpEndpoint string) TracerOption {
	return func(t *Tracer) {
		httpExporterOptions := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(otelHttpEndpoint),
		}

		if t.insecure {
			httpExporterOptions = append(httpExporterOptions, otlptracehttp.WithInsecure())
		}

		httpExporter, err := otlptrace.New(context.Background(), otlptracehttp.NewClient(httpExporterOptions...))
		if err != nil {
			otelzap.L().Sugar().Fatalw("Failed to create OTLP gRPC trace exporter", zap.Error(err))
		}

		batcher := trace.NewBatchSpanProcessor(httpExporter)

		t.providerOptions = append(t.providerOptions, trace.WithSpanProcessor(batcher))
	}
}

func WithTraceAutomaticEnv() TracerOption {
	return func(t *Tracer) {
		otelEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
		if otelEndpoint == "" {
			otelEndpoint = "localhost:4317"
		}

		otelInsecure := os.Getenv("OTEL_EXPORTER_OTLP_INSECURE") == "true"

		if otelInsecure {
			WithTraceInsecure()(t)
		}

		if strings.Contains(otelEndpoint, "4317") {
			WithGrpcTraceEndpoint(otelEndpoint)(t)
		} else if strings.Contains(otelEndpoint, "4318") {
			WithHttpTraceEndpoint(otelEndpoint)(t)
		}
	}
}

func WithTraceResources(res *resource.Resource) TracerOption {
	return func(t *Tracer) {
		t.resources = res
	}
}

func WithoutRegisterTraceProvider() TracerOption {
	return func(t *Tracer) {
		t.register = false
	}
}
