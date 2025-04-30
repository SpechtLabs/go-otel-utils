package otelzap

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerWithCtx is a wrapper for Logger that also carries a context.Context.
type LoggerWithCtx struct {
	ctx context.Context
	l   *Logger
}

// Context returns logger's context.
func (l LoggerWithCtx) Context() context.Context {
	return l.ctx
}

// Logger returns the underlying logger.
func (l LoggerWithCtx) Logger() *Logger {
	return l.l
}

// ZapLogger returns the underlying zap logger.
func (l LoggerWithCtx) ZapLogger() *zap.Logger {
	return l.l.Logger
}

// Sugar returns a sugared logger with the context.
func (l LoggerWithCtx) Sugar() SugaredLoggerWithCtx {
	return SugaredLoggerWithCtx{
		ctx: l.ctx,
		s:   l.l.Sugar(),
	}
}

// WithOptions clones the current Logger, applies the supplied Options,
// and returns the resulting Logger. It's safe to use concurrently.
func (l LoggerWithCtx) WithOptions(opts ...zap.Option) LoggerWithCtx {
	return LoggerWithCtx{
		ctx: l.ctx,
		l:   l.l.WithOptions(opts...),
	}
}

// Clone clones the current logger applying the supplied options.
func (l LoggerWithCtx) Clone(opts ...Option) LoggerWithCtx {
	return LoggerWithCtx{
		ctx: l.ctx,
		l:   l.l.Clone(opts...),
	}
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l LoggerWithCtx) Debug(msg string, fields ...zapcore.Field) {
	fields = l.l.logFields(l.ctx, zap.DebugLevel, msg, fields)
	l.l.skipCaller.Debug(msg, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l LoggerWithCtx) Info(msg string, fields ...zapcore.Field) {
	fields = l.l.logFields(l.ctx, zap.InfoLevel, msg, fields)
	l.l.skipCaller.Info(msg, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l LoggerWithCtx) Warn(msg string, fields ...zapcore.Field) {
	fields = l.l.logFields(l.ctx, zap.WarnLevel, msg, fields)
	l.l.skipCaller.Warn(msg, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l LoggerWithCtx) Error(msg string, fields ...zapcore.Field) {
	fields = l.l.logFields(l.ctx, zap.ErrorLevel, msg, fields)
	l.l.skipCaller.Error(msg, fields...)
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func (l LoggerWithCtx) DPanic(msg string, fields ...zapcore.Field) {
	fields = l.l.logFields(l.ctx, zap.DPanicLevel, msg, fields)
	l.l.skipCaller.DPanic(msg, fields...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func (l LoggerWithCtx) Panic(msg string, fields ...zapcore.Field) {
	fields = l.l.logFields(l.ctx, zap.PanicLevel, msg, fields)
	l.l.skipCaller.Panic(msg, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func (l LoggerWithCtx) Fatal(msg string, fields ...zapcore.Field) {
	fields = l.l.logFields(l.ctx, zap.FatalLevel, msg, fields)
	l.l.skipCaller.Fatal(msg, fields...)
}
