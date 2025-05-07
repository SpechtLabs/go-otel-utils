package otelzap_test

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sierrasoftworks/humane-errors-go"
	"github.com/spechtlabs/go-otel-utils/otelprovider"
	"github.com/spechtlabs/go-otel-utils/otelzap"
	"go.uber.org/zap"
)

func ExampleNew() {
	// Initialize zapLogger as normal
	zapLogger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("failed to initialize logger: %v", err)
		os.Exit(1)
	}

	// Initialize otel logging provider
	logProvider := otelprovider.NewLogger(
		otelprovider.WithLogAutomaticEnv(),
	)

	// Create otelZap Logger
	otelZapLogger := otelzap.New(zapLogger,
		otelzap.WithMinLevel(zap.InfoLevel),          // sets the minimal zap logging level on which the log message is recorded on the span.
		otelzap.WithAnnotateLevel(zap.WarnLevel),     // sets the minimal zap logging level on which spans will be annotated with the log fields as metadata.
		otelzap.WithErrorStatusLevel(zap.ErrorLevel), // sets the minimal zap logging level on which the span status is set to codes.Error.
		otelzap.WithStackTrace(false),                // configures the logger to capture logs with a stack trace.
		otelzap.WithLoggerProvider(logProvider),      // configures the logger to send logs via OTLP
	)

	humane_err := humane.Wrap(
		ioutil.WriteFile("demo.txt", []byte("This is an example"), os.ModePerm),
		"We couldn't write the demo.txt file to the current directory.",
		"Ensure you have write permissions to the current directory.",
		"Make sure you have free space on your disk.",
	)

	otelZapLogger.WithError(humane_err).Info("hello from zap",
		zap.String("foo", "bar"),
	)
}
