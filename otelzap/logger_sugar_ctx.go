package otelzap

import (
	"context"

	"go.uber.org/zap"
)

type SugaredLoggerWithCtx struct {
	ctx context.Context
	s   *SugaredLogger
}

// Desugar unwraps a SugaredLogger, exposing the original Logger. Desugaring
// is quite inexpensive, so it's reasonable for a single application to use
// both Loggers and SugaredLoggers, converting between them on the boundaries
// of performance-sensitive code.
func (s SugaredLoggerWithCtx) Desugar() LoggerWithCtx {
	return LoggerWithCtx{
		ctx: s.ctx,
		l:   s.s.Desugar(),
	}
}

// Debugf uses fmt.Sprintf to log a templated message.
func (s SugaredLoggerWithCtx) Debugf(template string, args ...interface{}) {
	s.s.logArgs(s.ctx, zap.DebugLevel, template, args)
	s.s.skipCaller.Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (s SugaredLoggerWithCtx) Infof(template string, args ...interface{}) {
	s.s.logArgs(s.ctx, zap.InfoLevel, template, args)
	s.s.skipCaller.Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (s SugaredLoggerWithCtx) Warnf(template string, args ...interface{}) {
	s.s.logArgs(s.ctx, zap.WarnLevel, template, args)
	s.s.skipCaller.Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (s SugaredLoggerWithCtx) Errorf(template string, args ...interface{}) {
	s.s.logArgs(s.ctx, zap.ErrorLevel, template, args)
	s.s.skipCaller.Errorf(template, args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func (s SugaredLoggerWithCtx) DPanicf(template string, args ...interface{}) {
	s.s.logArgs(s.ctx, zap.DPanicLevel, template, args)
	s.s.skipCaller.DPanicf(template, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func (s SugaredLoggerWithCtx) Panicf(template string, args ...interface{}) {
	s.s.logArgs(s.ctx, zap.PanicLevel, template, args)
	s.s.skipCaller.Panicf(template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (s SugaredLoggerWithCtx) Fatalf(template string, args ...interface{}) {
	s.s.logArgs(s.ctx, zap.FatalLevel, template, args)
	s.s.skipCaller.Fatalf(template, args...)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//
//	s.With(keysAndValues).Debug(msg)
func (s SugaredLoggerWithCtx) Debugw(msg string, keysAndValues ...interface{}) {
	s.s.logKVs(s.ctx, zap.DebugLevel, msg, keysAndValues)
	s.s.skipCaller.Debugw(msg, keysAndValues...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (s SugaredLoggerWithCtx) Infow(msg string, keysAndValues ...interface{}) {
	s.s.logKVs(s.ctx, zap.InfoLevel, msg, keysAndValues)
	s.s.skipCaller.Infow(msg, keysAndValues...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (s SugaredLoggerWithCtx) Warnw(msg string, keysAndValues ...interface{}) {
	s.s.logKVs(s.ctx, zap.WarnLevel, msg, keysAndValues)
	s.s.skipCaller.Warnw(msg, keysAndValues...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (s SugaredLoggerWithCtx) Errorw(msg string, keysAndValues ...interface{}) {
	s.s.logKVs(s.ctx, zap.ErrorLevel, msg, keysAndValues)
	s.s.skipCaller.Errorw(msg, keysAndValues...)
}

// DPanicw logs a message with some additional context. In development, the
// logger then panics. (See DPanicLevel for details.) The variadic key-value
// pairs are treated as they are in With.
func (s SugaredLoggerWithCtx) DPanicw(msg string, keysAndValues ...interface{}) {
	s.s.logKVs(s.ctx, zap.DPanicLevel, msg, keysAndValues)
	s.s.skipCaller.DPanicw(msg, keysAndValues...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func (s SugaredLoggerWithCtx) Panicw(msg string, keysAndValues ...interface{}) {
	s.s.logKVs(s.ctx, zap.PanicLevel, msg, keysAndValues)
	s.s.skipCaller.Panicw(msg, keysAndValues...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func (s SugaredLoggerWithCtx) Fatalw(msg string, keysAndValues ...interface{}) {
	s.s.logKVs(s.ctx, zap.FatalLevel, msg, keysAndValues)
	s.s.skipCaller.Fatalw(msg, keysAndValues...)
}
