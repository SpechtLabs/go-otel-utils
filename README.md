# go-otel-utils

A collection of Go utilities for simplifying OpenTelemetry integration in Go applications.

## Overview

This project offers a suite of tools designed to reduce the boilerplate code required to properly implement OpenTelemetry tracing and logging in Go applications. The utilities are designed with sensible defaults while allowing the flexibility to customize as needed.

## Components

The project consists of the following packages:

### otelprovider

An opinionated helper library for easily setting up OpenTelemetry trace and log providers with batteries included. The package handles the common initialization patterns and provides a clean API for your applications.
[See full documentation in otelprovider/README.md](otelprovider/README.md)

#### Features

- Simple initialization of OpenTelemetry trace providers
- Structured logging with OpenTelemetry integration
- Automatic resource detection for environment information
- Support for multiple exporters (OTLP, Console)
- Environment-aware configuration

### otelzap

A wrapper for the popular [zap](https://github.com/uber-go/zap) logging library that integrates seamlessly with OpenTelemetry. This package enables correlation between logs and traces for better observability.
[See full documentation in otelzap/README.md](otelzap/README.md)

## Installation

Install the packages you need:

``` bash
# For just the OpenTelemetry provider utilities
go get github.com/spechtlabs/go-otel-utils/otelprovider

# For zap integration with OpenTelemetry
go get github.com/spechtlabs/go-otel-utils/otelzap
```

## Getting Started

### Basic Setup with otelprovider

``` go
package main

import (
    "context"
    "log"

    "github.com/spechtlabs/go-otel-utils/otelprovider"
)

func main() {
    ctx := context.Background()

    // Initialize trace provider
    tp, err := otelprovider.NewTraceProvider(ctx, "my-service")
    if err != nil {
        log.Fatalf("Failed to initialize trace provider: %v", err)
    }
    defer tp.Shutdown(ctx)

    // Initialize logger
    logger, err := otelprovider.NewLogger("my-service", "development")
    if err != nil {
        log.Fatalf("Failed to initialize logger: %v", err)
    }

    // Your application code here
    logger.Info("Application started successfully")
}
```

### Complete Example

Explore the [example/](example/) directory for a complete working demonstration of these utilities in action.

## Configuration

The utilities support configuration through environment variables:

- `OTEL_EXPORTER_OTLP_ENDPOINT`: Endpoint for the OTLP exporter
- `OTEL_SERVICE_NAME`: Default service name if not specified
- `OTEL_ENVIRONMENT`: Environment (development, staging, production)
- `LOG_LEVEL`: Logging level (debug, info, warn, error)

## Design Philosophy

The go-otel-utils project follows these principles:

1. **Simplicity First**: Provide straightforward APIs that hide complexity
2. **Sensible Defaults**: Work well out-of-the-box with minimal configuration
3. **Production Ready**: Built with real-world production use cases in mind
4. **Extensibility**: Allow customization when the defaults don't fit
5. **Standards Compliance**: Follow OpenTelemetry best practices and standards

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.
