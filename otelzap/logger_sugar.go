package otelzap

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// A SugaredLogger wraps the base Logger functionality in a slower, but less
// verbose, API. Any Logger can be converted to a SugaredLogger with its Sugar
// method.
//
// Unlike the Logger, the SugaredLogger doesn't insist on structured logging.
// For each log level, it exposes three methods: one for loosely-typed
// structured logging, one for println-style formatting, and one for
// printf-style formatting. For example, SugaredLoggers can produce InfoLevel
// output with Infow ("info with" structured context), Info, or Infof.
type SugaredLogger struct {
	*zap.SugaredLogger
	skipCaller *zap.SugaredLogger

	l *Logger
}

// Desugar unwraps a SugaredLogger, exposing the original Logger. Desugaring
// is quite inexpensive, so it's reasonable for a single application to use
// both Loggers and SugaredLoggers, converting between them on the boundaries
// of performance-sensitive code.
func (s *SugaredLogger) Desugar() *Logger {
	return s.l
}

// With adds a variadic number of fields to the logging context. It accepts a
// mix of strongly-typed Field objects and loosely-typed key-value pairs. When
// processing pairs, the first element of the pair is used as the field key
// and the second as the field value.
//
// For example,
//
//	 sugaredLogger.With(
//	   "hello", "world",
//	   "failure", errors.New("oh no"),
//	   Stack(),
//	   "count", 42,
//	   "user", User{Name: "alice"},
//	)
//
// is the equivalent of
//
//	unsugared.With(
//	  String("hello", "world"),
//	  String("failure", "oh no"),
//	  Stack(),
//	  Int("count", 42),
//	  Object("user", User{Name: "alice"}),
//	)
//
// Note that the keys in key-value pairs should be strings. In development,
// passing a non-string key panics. In production, the logger is more
// forgiving: a separate error is logged, but the key-value pair is skipped
// and execution continues. Passing an orphaned key triggers similar behavior:
// panics in development and errors in production.
func (s *SugaredLogger) With(args ...interface{}) *SugaredLogger {
	return &SugaredLogger{
		SugaredLogger: s.SugaredLogger.With(args...),
		skipCaller:    s.skipCaller,
		l:             s.l,
	}
}

// Ctx returns a new sugared logger with the context.
func (s *SugaredLogger) Ctx(ctx context.Context) SugaredLoggerWithCtx {
	return SugaredLoggerWithCtx{
		ctx: ctx,
		s:   s,
	}
}

// Debugf uses fmt.Sprintf to log a templated message.
func (s *SugaredLogger) DebugfContext(ctx context.Context, template string, args ...interface{}) {
	s.logArgs(ctx, zap.DebugLevel, template, args)
	s.Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (s *SugaredLogger) InfofContext(ctx context.Context, template string, args ...interface{}) {
	s.logArgs(ctx, zap.InfoLevel, template, args)
	s.Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (s *SugaredLogger) WarnfContext(ctx context.Context, template string, args ...interface{}) {
	s.logArgs(ctx, zap.WarnLevel, template, args)
	s.Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (s *SugaredLogger) ErrorfContext(ctx context.Context, template string, args ...interface{}) {
	s.logArgs(ctx, zap.ErrorLevel, template, args)
	s.Errorf(template, args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func (s *SugaredLogger) DPanicfContext(ctx context.Context, template string, args ...interface{}) {
	s.logArgs(ctx, zap.DPanicLevel, template, args)
	s.DPanicf(template, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func (s *SugaredLogger) PanicfContext(ctx context.Context, template string, args ...interface{}) {
	s.logArgs(ctx, zap.PanicLevel, template, args)
	s.Panicf(template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (s *SugaredLogger) FatalfContext(ctx context.Context, template string, args ...interface{}) {
	s.logArgs(ctx, zap.FatalLevel, template, args)
	s.Fatalf(template, args...)
}

func (s *SugaredLogger) logArgs(
	ctx context.Context, lvl zapcore.Level, template string, args []interface{},
) {
	if lvl < s.l.minLevel {
		return
	}

	kvs := make([]zapcore.Field, 0, 1+numExtraAttr)
	kvs = append(kvs, zap.String("log.template", template))
	s.l.LogContext(ctx, lvl, fmt.Sprintf(template, args...), kvs...)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (s *SugaredLogger) DebugwContext(
	ctx context.Context, msg string, keysAndValues ...interface{},
) {
	s.logKVs(ctx, zap.DebugLevel, msg, keysAndValues)
	s.Debugw(msg, keysAndValues...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (s *SugaredLogger) InfowContext(
	ctx context.Context, msg string, keysAndValues ...interface{},
) {
	s.logKVs(ctx, zap.InfoLevel, msg, keysAndValues)
	s.Infow(msg, keysAndValues...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (s *SugaredLogger) WarnwContext(
	ctx context.Context, msg string, keysAndValues ...interface{},
) {
	s.logKVs(ctx, zap.WarnLevel, msg, keysAndValues)
	s.Warnw(msg, keysAndValues...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (s *SugaredLogger) ErrorwContext(
	ctx context.Context, msg string, keysAndValues ...interface{},
) {
	s.logKVs(ctx, zap.ErrorLevel, msg, keysAndValues)
	s.Errorw(msg, keysAndValues...)
}

// DPanicw logs a message with some additional context. In development, the
// logger then panics. (See DPanicLevel for details.) The variadic key-value
// pairs are treated as they are in With.
func (s *SugaredLogger) DPanicwContext(
	ctx context.Context, msg string, keysAndValues ...interface{},
) {
	s.logKVs(ctx, zap.DPanicLevel, msg, keysAndValues)
	s.DPanicw(msg, keysAndValues...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func (s *SugaredLogger) PanicwContext(
	ctx context.Context, msg string, keysAndValues ...interface{},
) {
	s.logKVs(ctx, zap.PanicLevel, msg, keysAndValues)
	s.Panicw(msg, keysAndValues...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func (s *SugaredLogger) FatalwContext(
	ctx context.Context, msg string, keysAndValues ...interface{},
) {
	s.logKVs(ctx, zap.FatalLevel, msg, keysAndValues)
	s.Fatalw(msg, keysAndValues...)
}

func (s *SugaredLogger) logKVs(
	ctx context.Context, lvl zapcore.Level, msg string, args []interface{},
) {
	if lvl < s.l.minLevel {
		return
	}

	kvs := make([]zapcore.Field, 0, len(args)/2)

	for i := 0; i < len(args); i++ {
		field := args[i]

		switch field := field.(type) {

		// in case it's a zapcore.Field we know that key and value are encoded in the zapcore.Field
		case zapcore.Field:
			kvs = append(kvs, field)

		// in case it's a string, we assume it's key + value separate
		case string:
			kvs = append(kvs, zap.Any(field, args[i+1]))

			// Also increment i because we just read args[i+1]
			i += 1
		}
	}

	s.l.LogContext(ctx, lvl, msg, kvs...)
}
