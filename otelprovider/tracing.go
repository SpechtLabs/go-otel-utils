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

// Tracer collects the configuration NewTracer builds a trace provider from.
type Tracer struct {
	resources       *resource.Resource
	providerOptions []trace.TracerProviderOption
	insecure        bool
	register        bool
}

// NewTracer builds an OpenTelemetry trace provider with the given options,
// describing this process with the resource from WithTraceResources (by
// default, one detected from the environment). Unless
// WithoutRegisterTraceProvider is given, it also becomes the global trace
// provider.
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

// TracerOption configures the trace provider NewTracer builds.
type TracerOption func(t *Tracer)

// WithTraceInsecure turns off TLS for the exporters the endpoint options after
// it set up.
func WithTraceInsecure() TracerOption {
	return func(t *Tracer) {
		t.insecure = true
	}
}

// WithGrpcTraceEndpoint exports spans over OTLP/gRPC to the given endpoint, in
// batches. It exits the process when the exporter can't be created.
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

// WithHttpTraceEndpoint exports spans over OTLP/HTTP to the given endpoint, in
// batches. It exits the process when the exporter can't be created.
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

// WithTraceAutomaticEnv configures the exporter from
// OTEL_EXPORTER_OTLP_ENDPOINT: OTLP/gRPC when the endpoint contains port 4317,
// OTLP/HTTP when it contains 4318. OTEL_EXPORTER_OTLP_INSECURE=true turns off
// TLS. Without an endpoint, nothing is exported.
func WithTraceAutomaticEnv() TracerOption {
	return func(t *Tracer) {
		otelEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
		if otelEndpoint == "" {
			return // if no endpoint is set, do not configure the exporter
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

// WithTraceResources replaces the resource that describes this process in
// every span.
func WithTraceResources(res *resource.Resource) TracerOption {
	return func(t *Tracer) {
		t.resources = res
	}
}

// WithoutRegisterTraceProvider leaves the global trace provider alone;
// NewTracer only returns the provider it built.
func WithoutRegisterTraceProvider() TracerOption {
	return func(t *Tracer) {
		t.register = false
	}
}
