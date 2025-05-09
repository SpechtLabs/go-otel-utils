package otelzap_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/sierrasoftworks/humane-errors-go"
	"github.com/spechtlabs/go-otel-utils/otelzap"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLogger() *bytes.Buffer {
	// Capture logs for later assertions
	enc := zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
	buf := &bytes.Buffer{}
	writer := zapcore.AddSync(buf) //zap.CombineWriteSyncers(zaptest.NewTestingWriter(t), )
	level := zap.NewAtomicLevelAt(zapcore.DebugLevel)

	otelZapLogger := otelzap.New(zap.New(zapcore.NewCore(enc, writer, level)))
	otelzap.ReplaceGlobals(otelZapLogger)

	return buf
}

func TestLogOnce(t *testing.T) {
	// Capture logs for later assertions
	buf := initLogger()
	buf.Reset()

	humaneErr := humane.New(
		"message",
		"advice",
	)

	// Write a log once
	otelzap.L().WithError(humaneErr).Error("Test Message", zap.String("foo", "bar"))
	assert.Contains(t, buf.String(), "error\tTest Message\t{\"foo\": \"bar\", \"error\": \"message\", \"error_advice\": [\"advice\"]}")
	buf.Reset()

	// And it shouldn't show up in the next one
	otelzap.L().Error("Test Message", zap.String("foo", "bar"))
	assert.Contains(t, buf.String(), "error\tTest Message\t{\"foo\": \"bar\"}")
}

func TestErrorContext(t *testing.T) {
	buf := initLogger()
	buf.Reset()
	ctx := context.Background()

	otelzap.L().ErrorContext(ctx, "Test Message", zap.String("foo", "bar"))
	assert.Contains(t, buf.String(), "error\tTest Message\t{\"foo\": \"bar\"}")

	buf.Reset()

	otelzap.L().Ctx(ctx).Error("Test Message", zap.String("foo", "bar"))
	assert.Contains(t, buf.String(), "error\tTest Message\t{\"foo\": \"bar\"}")
}

func TestErrorSugared(t *testing.T) {
	buf := initLogger()
	buf.Reset()
	otelzap.L().Sugar().Errorw("Test Message", "foo", "bar")
	assert.Contains(t, buf.String(), "error\tTest Message\t{\"foo\": \"bar\"}")
}

func TestErrorCtxSugared(t *testing.T) {
	buf := initLogger()
	buf.Reset()
	ctx := context.Background()
	otelzap.L().Ctx(ctx).Sugar().Errorw("Test Message", "foo", "bar")
	assert.Contains(t, buf.String(), "error\tTest Message\t{\"foo\": \"bar\"}")
}
