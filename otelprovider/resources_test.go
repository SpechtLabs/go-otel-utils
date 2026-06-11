package otelprovider

import "testing"

// newOtelResources merges our service attributes with resource.Default(), whose
// schema URL advances with every OTel SDK release. Building the resource must
// not depend on that schema matching ours — otherwise each SDK bump panics the
// consumer at init (the recurring "conflicting Schema URL" crash). This test
// guards against reintroducing a hard-coded, conflicting schema URL.
func TestNewOtelResourcesDoesNotConflictWithSDKSchema(t *testing.T) {
	res := newOtelResources()
	if res == nil {
		t.Fatal("expected a resource, got nil")
	}

	// Sanity check: the service attributes we set are present.
	var sawServiceName bool
	for _, attr := range res.Attributes() {
		if attr.Key == "service.name" {
			sawServiceName = true
		}
	}
	if !sawServiceName {
		t.Fatal("expected service.name to be set on the resource")
	}
}
