package otelzap_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"sync"
	"testing"

	"github.com/spechtlabs/go-otel-utils/otelzap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// lockedBuffer is a bytes.Buffer that several goroutines can log to.
type lockedBuffer struct {
	buf bytes.Buffer
	mu  sync.Mutex
}

func (b *lockedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *lockedBuffer) Sync() error { return nil }

// entries decodes every JSON log line written so far.
func (b *lockedBuffer) entries(t *testing.T) []map[string]any {
	t.Helper()
	b.mu.Lock()
	defer b.mu.Unlock()
	var out []map[string]any
	for line := range strings.SplitSeq(strings.TrimSpace(b.buf.String()), "\n") {
		entry := map[string]any{}
		require.NoError(t, json.Unmarshal([]byte(line), &entry), line)
		out = append(out, entry)
	}
	return out
}

func newJSONLogger(opts ...otelzap.Option) (*otelzap.Logger, *lockedBuffer) {
	buf := &lockedBuffer{}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), buf, zapcore.DebugLevel)
	return otelzap.New(zap.New(core), opts...), buf
}

func TestWithLeavesTheLoggerUnchanged(t *testing.T) {
	logger, buf := newJSONLogger()

	logger.With(zap.String("request_id", "a")).Info("first")
	_ = logger.With(zap.String("request_id", "b"))
	logger.Info("second")

	entries := buf.entries(t)
	require.Len(t, entries, 2)
	assert.Equal(t, "a", entries[0]["request_id"])
	assert.NotContains(t, entries[1], "request_id", "With must not add fields to the logger it was called on")
}

func TestWithOnASharedLoggerConcurrently(t *testing.T) {
	logger, buf := newJSONLogger()

	const n = 64
	var wg sync.WaitGroup
	for i := range n {
		wg.Go(func() {
			logger.With(zap.Int("request", i)).Info("handled", zap.Int("handler", i))
		})
	}
	wg.Wait()

	entries := buf.entries(t)
	require.Len(t, entries, n)
	for _, entry := range entries {
		assert.Equal(t, entry["handler"], entry["request"], "a request's field landed on another request's line: %v", entry)
	}
}

func TestWithOptionsCopiesDoNotShareFields(t *testing.T) {
	// Three extra fields, then a fourth through WithOptions, leave the
	// field slice with spare capacity: the case where appending to a
	// copied slice writes into the original's backing array.
	base, buf := newJSONLogger(otelzap.WithExtraFields(
		zap.String("f1", "1"), zap.String("f2", "2"), zap.String("f3", "3"),
	))
	base = base.WithOptions(zap.Fields(zap.String("f4", "4")))

	a := base.WithOptions(zap.Fields(zap.String("copy", "a")))
	b := base.WithOptions(zap.Fields(zap.String("copy", "b")))
	a.Info("from a")
	b.Info("from b")

	for _, entry := range buf.entries(t) {
		want := map[string]string{"from a": "a", "from b": "b"}[entry["msg"].(string)]
		assert.Equal(t, want, entry["copy"], "%v", entry)
	}
}
