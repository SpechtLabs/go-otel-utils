package otelprovider

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/spechtlabs/go-otel-utils/otelzap"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/zap"
)

type Logger struct {
	providerOptions []log.LoggerProviderOption
	insecure        bool
	resources       *resource.Resource
	register        bool
}

func NewLogger(opts ...LoggerOption) *log.LoggerProvider {
	l := &Logger{
		insecure:        false,
		providerOptions: []log.LoggerProviderOption{},
		resources:       newOtelResources(),
		register:        true,
	}

	for _, opt := range opts {
		opt(l)
	}

	l.providerOptions = append(l.providerOptions, log.WithResource(l.resources))
	logProvider := log.NewLoggerProvider(l.providerOptions...)

	// Register the Provider globally
	if l.register {
		global.SetLoggerProvider(logProvider)
	}

	return logProvider
}

// TracerOption applies a configuration to the given config.
type LoggerOption func(t *Logger)

func WithLogInsecure() LoggerOption {
	return func(t *Logger) {
		t.insecure = true
	}
}

func WithGrpcLogEndpoint(otelGrpcEndpoint string) LoggerOption {
	return func(t *Logger) {
		grpcExporterOptions := []otlploggrpc.Option{
			otlploggrpc.WithEndpoint(otelGrpcEndpoint),
		}

		if t.insecure {
			grpcExporterOptions = append(grpcExporterOptions, otlploggrpc.WithInsecure())
		}

		grpcExporter, err := otlploggrpc.New(context.Background(), grpcExporterOptions...)
		if err != nil {
			otelzap.L().Sugar().Fatalw("Failed to create OTLP gRPC logs exporter", zap.Error(err))
		}

		batcher := log.NewBatchProcessor(grpcExporter,
			log.WithMaxQueueSize(10_000),
			log.WithExportMaxBatchSize(10_000),
			log.WithExportInterval(10*time.Second),
			log.WithExportTimeout(10*time.Second),
		)

		t.providerOptions = append(t.providerOptions, log.WithProcessor(batcher))
	}
}

func WithHttpLogEndpoint(otelHttpEndpoint string) LoggerOption {
	return func(t *Logger) {
		httpExporterOptions := []otlploghttp.Option{
			otlploghttp.WithEndpoint(otelHttpEndpoint),
		}

		if t.insecure {
			httpExporterOptions = append(httpExporterOptions, otlploghttp.WithInsecure())
		}

		httpExporter, err := otlploghttp.New(context.Background(), httpExporterOptions...)
		if err != nil {
			otelzap.L().Sugar().Fatalw("Failed to create OTLP HTTP logs exporter", zap.Error(err))
		}

		batcher := log.NewBatchProcessor(httpExporter,
			log.WithMaxQueueSize(10_000),
			log.WithExportMaxBatchSize(10_000),
			log.WithExportInterval(10*time.Second),
			log.WithExportTimeout(10*time.Second),
		)

		t.providerOptions = append(t.providerOptions, log.WithProcessor(batcher))
	}
}

func WithLogAutomaticEnv() LoggerOption {
	return func(t *Logger) {
		otelEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
		if otelEndpoint == "" {
			return // if no endpoint is set, do not configure the exporter
		}

		otelInsecure := os.Getenv("OTEL_EXPORTER_OTLP_INSECURE") == "true"

		if otelInsecure {
			WithLogInsecure()(t)
		}

		if strings.Contains(otelEndpoint, "4317") {
			WithGrpcLogEndpoint(otelEndpoint)(t)
		} else if strings.Contains(otelEndpoint, "4318") {
			WithHttpLogEndpoint(otelEndpoint)(t)
		}
	}
}

func WithLogResources(res *resource.Resource) LoggerOption {
	return func(t *Logger) {
		t.resources = res
	}
}

func WithoutRegisterLogProvider() LoggerOption {
	return func(t *Logger) {
		t.register = false
	}
}
