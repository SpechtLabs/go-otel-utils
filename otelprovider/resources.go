package otelprovider

import (
	"os"
	"path/filepath"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func newOtelResources() *resource.Resource {
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = filepath.Base(os.Args[0])
	}

	serviceVersion := os.Getenv("OTEL_SERVICE_VERSION")
	if serviceVersion == "" {
		serviceVersion = "0.0.0-unset"
	}

	// Build our attributes as a schemaless resource (no schema URL) rather than
	// pinning a specific semconv schema. resource.Merge rejects two resources
	// with different, non-empty schema URLs — and resource.Default()'s schema
	// advances with every OTel SDK release. A schemaless resource carries no
	// schema URL, so the merge always succeeds and inherits Default()'s schema,
	// regardless of which SDK version the consumer builds against.
	res, err := resource.Merge(resource.Default(),
		resource.NewSchemaless(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		))
	if err != nil {
		// Merge only errors on conflicting schema URLs, which a schemaless
		// resource cannot trigger; treat any future error as non-fatal and
		// fall back to the default resource instead of crashing init.
		return resource.Default()
	}

	return res
}
