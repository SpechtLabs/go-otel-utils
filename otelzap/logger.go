package otelzap

import (
	"context"
	"fmt"
	"runtime"

	"github.com/aws/smithy-go/logging"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a thin wrapper for zap.Logger that adds Ctx method.
type Logger struct {
	*zap.Logger
	skipCaller *zap.Logger

	provider   log.LoggerProvider
	version    string
	schemaURL  string
	otelLogger log.Logger

	minLevel         zapcore.Level
	errorStatusLevel zapcore.Level
	minAnnotateLevel zapcore.Level

	caller     bool
	stackTrace bool

	// extraFields contains a number of zap.Fields that are added to every log entry
	extraFields []zap.Field
	callerDepth int
}

// New creates a new Logger instance with specified options and returns it along
// with an undo function used for cleanup.
func New(logger *zap.Logger, opts ...Option) *Logger {
	l := &Logger{
		Logger:     logger,
		skipCaller: logger.WithOptions(zap.AddCallerSkip(1)),

		provider: global.GetLoggerProvider(),

		minLevel:         zap.InfoLevel,
		errorStatusLevel: zap.ErrorLevel,
		minAnnotateLevel: zap.WarnLevel,
		caller:           true,
		callerDepth:      0,
	}
	for _, opt := range opts {
		opt(l)
	}
	l.otelLogger = l.newOtelLogger(logger.Name())

	return l
}

func (l *Logger) newOtelLogger(name string) log.Logger {
	var opts []log.LoggerOption
	if l.version != "" {
		opts = append(opts, log.WithInstrumentationVersion(l.version))
	}
	if l.schemaURL != "" {
		opts = append(opts, log.WithSchemaURL(l.schemaURL))
	}
	return l.provider.Logger(name, opts...)
}

// WithOptions clones the current Logger, applies the supplied Options,
// and returns the resulting Logger. It's safe to use concurrently.
func (l *Logger) WithOptions(opts ...zap.Option) *Logger {
	extraFields := []zap.Field{}
	// zap.New side effect is extracting fields from .WithOptions(zap.Fields(...))
	zap.New(&fieldExtractorCore{extraFields: &extraFields}, opts...)
	clone := *l
	clone.Logger = l.Logger.WithOptions(opts...)
	clone.skipCaller = l.skipCaller.WithOptions(opts...)
	clone.extraFields = append(clone.extraFields, extraFields...)
	return &clone
}

// Sugar wraps the Logger to provide a more ergonomic, but slightly slower,
// API. Sugaring a Logger is quite inexpensive, so it's reasonable for a
// single application to use both Loggers and SugaredLoggers, converting
// between them on the boundaries of performance-sensitive code.
func (l *Logger) Sugar() *SugaredLogger {
	return &SugaredLogger{
		SugaredLogger: l.Logger.Sugar(),
		skipCaller:    l.skipCaller.Sugar(),
		l:             l,
	}
}

// Clone clones the current logger applying the supplied options.
func (l *Logger) Clone(opts ...Option) *Logger {
	clone := *l
	for _, opt := range opts {
		opt(&clone)
	}
	return &clone
}

// Ctx returns a new logger with the context.
func (l *Logger) Ctx(ctx context.Context) LoggerWithCtx {
	return LoggerWithCtx{
		ctx: ctx,
		l:   l,
	}
}

func (l *Logger) DebugContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	fields = l.logFields(ctx, zap.DebugLevel, msg, fields)
	l.skipCaller.Debug(msg, fields...)
}

func (l *Logger) InfoContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	fields = l.logFields(ctx, zap.InfoLevel, msg, fields)
	l.skipCaller.Info(msg, fields...)
}

func (l *Logger) WarnContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	fields = l.logFields(ctx, zap.WarnLevel, msg, fields)
	l.skipCaller.Warn(msg, fields...)
}

func (l *Logger) ErrorContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	fields = l.logFields(ctx, zap.ErrorLevel, msg, fields)
	l.skipCaller.Error(msg, fields...)
}

func (l *Logger) DPanicContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	fields = l.logFields(ctx, zap.DPanicLevel, msg, fields)
	l.skipCaller.DPanic(msg, fields...)
}

func (l *Logger) PanicContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	fields = l.logFields(ctx, zap.PanicLevel, msg, fields)
	l.skipCaller.Panic(msg, fields...)
}

func (l *Logger) FatalContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	fields = l.logFields(ctx, zap.FatalLevel, msg, fields)
	l.skipCaller.Fatal(msg, fields...)
}

func (l *Logger) Logf(classification logging.Classification, format string, fields ...interface{}) {
	msg := fmt.Sprintf(format, fields...)

	switch classification {
	case logging.Warn:
		l.skipCaller.Warn(msg)

	case logging.Debug:
		l.skipCaller.Debug(msg)

	default:
		l.skipCaller.Info(msg)
	}
}

func (l *Logger) logFields(
	ctx context.Context, lvl zapcore.Level, msg string, fields []zapcore.Field,
) []zapcore.Field {
	if len(l.extraFields) > 0 {
		fields = append(fields, l.extraFields...)
	}

	if lvl >= l.minLevel {
		l.log(ctx, lvl, msg, convertFields(fields))
	}
	return fields
}

func (l *Logger) log(
	ctx context.Context, lvl zapcore.Level, msg string, kvs []log.KeyValue,
) {
	if lvl >= l.minAnnotateLevel || lvl >= l.errorStatusLevel {
		if span := trace.SpanFromContext(ctx); span.IsRecording() {
			if lvl >= l.minAnnotateLevel {
				for _, kv := range kvs {
					span.SetAttributes(Attribute(kv.Key, kv.Value))
				}
			}

			if lvl >= l.errorStatusLevel {
				span.SetStatus(codes.Error, msg)
				span.RecordError(fmt.Errorf("%s", msg))
			}
		}
	}

	record := log.Record{}
	record.SetBody(log.StringValue(msg))
	record.SetSeverity(convertLevel(lvl))

	if l.caller {
		if fn, file, line, ok := runtimeCaller(4 + l.callerDepth); ok {
			if fn != "" {
				kvs = append(kvs, log.String("code.function", fn))
			}
			if file != "" {
				kvs = append(kvs, log.String("code.filepath", file))
				kvs = append(kvs, log.Int("code.lineno", line))
			}
		}
	}

	if l.stackTrace {
		stackTrace := make([]byte, 2048)
		n := runtime.Stack(stackTrace, false)
		kvs = append(kvs, log.String("exception.stacktrace", string(stackTrace[:n])))
	}

	if len(kvs) > 0 {
		record.AddAttributes(kvs...)
	}

	l.otelLogger.Emit(ctx, record)
}
