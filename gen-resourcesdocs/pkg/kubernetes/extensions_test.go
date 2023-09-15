package kubernetes

import (
	"testing"

	"github.com/go-openapi/spec"
)

func TestGetGVKExtension(t *testing.T) {
	definition := spec.Schema{
		VendorExtensible: spec.VendorExtensible{
			Extensions: spec.Extensions{
				"x-kubernetes-group-version-kind": []interface{}{
					map[string]interface{}{
						"group":   "apps",
						"version": "v1",
						"kind":    "Deployment",
					},
				},
			},
		},
	}
	extension, found, err := getGVKExtension(definition.Extensions)
	if !found {
		t.Fatal("Extension should be found")
	}
	if err != nil {
		t.Fatalf("Extension should be found without error: %v", err)
	}
	if extension.Group != "apps" {
		t.Errorf("Group should be %s but is %s", "apps", extension.Group)
	}
	if extension.Version.String() != "v1" {
		t.Errorf("Version should be %s but is %s", "v1", extension.Version.String())
	}
	if extension.Kind != "Deployment" {
		t.Errorf("Kind should be %s but is %s", "Deployment", extension.Kind)
	}
}
