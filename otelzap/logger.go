package otelzap

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/smithy-go/logging"
	"github.com/sierrasoftworks/humane-errors-go"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
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
	// extraFieldsOnce contains a number of zap.Fields that are added to only the next log entry
	extraFieldsOnce []zap.Field
	callerDepth     int
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

// WithError adds a humane.Error to the logging context.
//
// For example,
//
//		 sugaredLogger.WithError(
//	    humane.New("foo", "bar")
//		)
func (l *Logger) WithError(err error) *Logger {
	zapFields := make([]zap.Field, 0)
	zapFields = append(zapFields, zap.Error(err))

	advice := make([]string, 0)
	causes := make([]error, 0)
	for err != nil {
		var herr humane.Error
		if ok := errors.As(err, &herr); ok {
			causes = append(causes, err)
			advice = append(advice, herr.Advice()...)
		}

		err = errors.Unwrap(err)
	}

	if len(advice) > 0 {
		zapFields = append(zapFields, zap.Strings("error_advice", advice))
	}

	if len(causes) > 1 {
		zapFields = append(zapFields, zap.Errors("error_causes", causes[1:]))
	}

	return l.With(zapFields...)
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	l.extraFieldsOnce = append(l.extraFieldsOnce, fields...)
	return l
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

// Log logs a message at the specified level. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
// Any Fields that require  evaluation (such as Objects) are evaluated upon
// invocation of Log.
func (l *Logger) Log(lvl zapcore.Level, msg string, fields ...zapcore.Field) {
	fields = l.logFields(fields)
	l.skipCaller.Log(lvl, msg, fields...)
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Debug(msg string, fields ...zapcore.Field) {
	fields = l.logFields(fields)
	l.skipCaller.Debug(msg, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Info(msg string, fields ...zapcore.Field) {
	fields = l.logFields(fields)
	l.skipCaller.Info(msg, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Warn(msg string, fields ...zapcore.Field) {
	fields = l.logFields(fields)
	l.skipCaller.Warn(msg, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Error(msg string, fields ...zapcore.Field) {
	fields = l.logFields(fields)
	l.skipCaller.Error(msg, fields...)
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func (l *Logger) DPanic(msg string, fields ...zapcore.Field) {
	fields = l.logFields(fields)
	l.skipCaller.DPanic(msg, fields...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func (l *Logger) Panic(msg string, fields ...zapcore.Field) {
	fields = l.logFields(fields)
	l.skipCaller.Panic(msg, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func (l *Logger) Fatal(msg string, fields ...zapcore.Field) {
	fields = l.logFields(fields)
	l.skipCaller.Fatal(msg, fields...)
}

func (l *Logger) LogContext(ctx context.Context, lvl zapcore.Level, msg string, fields ...zapcore.Field) {
	fields = l.logFields(fields)
	l.Ctx(ctx).l.skipCaller.Log(lvl, msg, fields...)
}

func (l *Logger) DebugContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.Ctx(ctx).l.skipCaller.Debug(msg, fields...)
}

func (l *Logger) InfoContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.Ctx(ctx).l.skipCaller.Info(msg, fields...)
}

func (l *Logger) WarnContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.Ctx(ctx).l.skipCaller.Warn(msg, fields...)
}

func (l *Logger) ErrorContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.Ctx(ctx).l.skipCaller.Error(msg, fields...)
}

func (l *Logger) DPanicContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.Ctx(ctx).l.skipCaller.DPanic(msg, fields...)
}

func (l *Logger) PanicContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.Ctx(ctx).l.skipCaller.Panic(msg, fields...)
}

func (l *Logger) FatalContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.Ctx(ctx).l.skipCaller.Fatal(msg, fields...)
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

func (l *Logger) logFields(fields []zapcore.Field) []zapcore.Field {
	if len(l.extraFields) > 0 {
		fields = append(fields, l.extraFields...)
	}

	if len(l.extraFieldsOnce) > 0 {
		fields = append(fields, l.extraFieldsOnce...)
		l.extraFieldsOnce = make([]zap.Field, 0)
	}

	return fields
}
