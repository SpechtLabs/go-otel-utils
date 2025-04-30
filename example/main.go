package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"github.com/SpechtLabs/go-otel-utils/otelprovider"
	"github.com/SpechtLabs/go-otel-utils/otelzap"
)

func main() {
	logProvider := otelprovider.NewLogger(
		otelprovider.WithLogAutomaticEnv(),
	)

	traceProvider := otelprovider.NewTracer(
		otelprovider.WithTraceAutomaticEnv(),
	)

	// Initialize Logging
	debug := os.Getenv("DEBUG") == "true"
	var zapLogger *zap.Logger
	var err error
	if debug {
		zapLogger, err = zap.NewDevelopment()
	} else {
		zapLogger, err = zap.NewProduction()
	}
	if err != nil {
		fmt.Printf("failed to initialize logger: %v", err)
		os.Exit(1)
	}

	// Replace zap global
	undoZapGlobals := zap.ReplaceGlobals(zapLogger)

	// Redirect stdlib log to zap
	undoStdLogRedirect := zap.RedirectStdLog(zapLogger)

	// Create otelLogger
	otelZapLogger := otelzap.New(zapLogger,
		otelzap.WithCaller(true),
		otelzap.WithMinLevel(zap.InfoLevel),
		otelzap.WithAnnotateLevel(zap.WarnLevel),
		otelzap.WithErrorStatusLevel(zap.ErrorLevel),
		otelzap.WithStackTrace(false),
		otelzap.WithLoggerProvider(logProvider),
	)

	// Replace global otelZap logger
	undoOtelZapGlobals := otelzap.ReplaceGlobals(otelZapLogger)

	defer func() {
		if err := traceProvider.ForceFlush(context.Background()); err != nil {
			otelzap.L().Warn("failed to flush traces")
		}

		if err := logProvider.ForceFlush(context.Background()); err != nil {
			otelzap.L().Warn("failed to flush logs")
		}

		if err := traceProvider.Shutdown(context.Background()); err != nil {
			panic(err)
		}

		if err := logProvider.Shutdown(context.Background()); err != nil {
			panic(err)
		}

		undoStdLogRedirect()
		undoOtelZapGlobals()
		undoZapGlobals()
	}()

	tracer := otel.Tracer("tracer")

	ctx := context.Background()
	baseCtx, baseSpan := tracer.Start(ctx, "root")
	defer baseSpan.End()

	debugCtx, debugSpan := tracer.Start(baseCtx, "debug")
	otelzap.L().Sugar().Ctx(debugCtx).Debugw("hello from zap",
		zap.Error(errors.New("hello world")),
		zap.String("foo", "bar"))
	debugSpan.End()

	infoCtx, infoSpan := tracer.Start(baseCtx, "info")
	otelzap.L().Sugar().Ctx(infoCtx).Infow("hello from zap",
		zap.Error(errors.New("hello world")),
		zap.String("foo", "bar"))
	infoSpan.End()

	warnCtx, warnSpan := tracer.Start(baseCtx, "warn")
	otelzap.L().Sugar().Ctx(warnCtx).Warnw("hello from zap",
		zap.Error(errors.New("hello world")),
		zap.String("foo", "bar"))
	warnSpan.End()

	errorCtx, errorSpan := tracer.Start(baseCtx, "error")
	otelzap.L().Sugar().Ctx(errorCtx).Errorw("hello from zap",
		zap.Error(errors.New("hello world")),
		zap.String("foo", "bar"))
	errorSpan.End()
}
