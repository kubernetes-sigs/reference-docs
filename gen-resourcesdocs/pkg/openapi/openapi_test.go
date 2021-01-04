package openapi_test

import (
	"testing"

	"github.com/kubernetes-sigs/reference-docs/gen-resourcesdocs/pkg/openapi"
)

func TestLoadOpenAPISpecV119(t *testing.T) {
	spec, err := openapi.LoadOpenAPISpec("../../api/v1.19/swagger.json")
	if err != nil {
		t.Errorf("Failed to load spec")
	}
	if len(spec.Definitions) != 617 {
		t.Errorf("Spec should contain %d definition but contains %d", 617, len(spec.Definitions))
	}
}
