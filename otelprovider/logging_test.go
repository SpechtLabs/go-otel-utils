package otelprovider

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
)

// otelzap picks up the log provider through otel.GetLoggerProvider, so
// NewLogger has to register it there. A provider set anywhere else, such as
// the deprecated otel/log/global package, would leave otelzap on the no-op
// provider and silently drop every log record.
func TestNewLoggerRegistersGlobalProvider(t *testing.T) {
	previous := otel.GetLoggerProvider()
	t.Cleanup(func() { otel.SetLoggerProvider(previous) })

	lp := NewLogger()
	t.Cleanup(func() { _ = lp.Shutdown(context.Background()) })

	if otel.GetLoggerProvider() != lp {
		t.Fatal("expected NewLogger to register its provider with otel.SetLoggerProvider")
	}
}

func TestNewLoggerWithoutRegisterLeavesGlobalProvider(t *testing.T) {
	previous := otel.GetLoggerProvider()
	t.Cleanup(func() { otel.SetLoggerProvider(previous) })

	lp := NewLogger(WithoutRegisterLogProvider())
	t.Cleanup(func() { _ = lp.Shutdown(context.Background()) })

	if otel.GetLoggerProvider() == lp {
		t.Fatal("expected WithoutRegisterLogProvider to keep the global provider unchanged")
	}
}
