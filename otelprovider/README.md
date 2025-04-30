# OTelProvider

An opinionated helper library for easily setting up OpenTelemetry log and trace providers in Go applications.

## Overview

OTelProvider simplifies the process of configuring OpenTelemetry (OTel) in your Go applications by providing ready-to-use functions for setting up tracing and logging with sensible defaults. This library helps reduce boilerplate code and ensures consistent telemetry configuration across your services.

## Features

- Simple API for initializing OpenTelemetry trace providers
- Structured logging with OpenTelemetry integration
- Resource detection for common environment information
- Support for multiple exporters (OTLP, Console)
- Environment-based configuration

## Installation

``` bash
go get github.com/spechtlabs/go-otel-utils/otelprovider
```

## Usage

### Basic Setup

``` go
package main

import (
    "context"
    "log"
    
    "github.com/spechtlabs/go-otel-utils/otelprovider"
)

func main() {
    ctx := context.Background()
    
    // Initialize the trace provider
    tp, err := otelprovider.NewTraceProvider(ctx, "my-service")
    if err != nil {
        log.Fatalf("Failed to initialize trace provider: %v", err)
    }
    defer tp.Shutdown(ctx)
    
    // Initialize the logger
    logger, err := otelprovider.NewLogger("my-service", "development")
    if err != nil {
        log.Fatalf("Failed to initialize logger: %v", err)
    }
    
    // Your application code here
    logger.Info("Application started successfully")
    
    // Create spans and add logs as needed
    ctx, span := tp.Tracer("component-name").Start(ctx, "operation-name")
    defer span.End()
    
    // ...
}
```

### Complete Example

See the [example/](example/) directory for a complete working example of how to use the otelprovider in your application.

### Configuration Options

The library offers various configuration options through environment variables:

- `OTEL_EXPORTER_OTLP_ENDPOINT`: Endpoint for the OTLP exporter
- `OTEL_SERVICE_NAME`: Default service name if not specified
- `OTEL_ENVIRONMENT`: Environment (development, staging, production)
- `LOG_LEVEL`: Logging level (debug, info, warn, error)

## Opinionated Decisions

The otelprovider library makes several opinionated choices to simplify telemetry setup:

1. **Sensible Defaults**: Provides reasonable default configurations that work out of the box
2. **Resource Detection**: Automatically detects and attaches environment information to telemetry data
3. **Structured Logging**: Uses structured logging for better searchability and context
4. **Correlation**: Ensures logs are correlated with traces for better observability
5. **Environment-Aware**: Configures exporters and sampling based on the runtime environment

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
