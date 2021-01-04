package kubernetes_test

import (
	"testing"

	"github.com/kubernetes-sigs/reference-docs/gen-resourcesdocs/pkg/kubernetes"
)

func TestSpecV119(t *testing.T) {
	spec, err := kubernetes.NewSpec("../../api/v1.19/swagger.json")
	if err != nil {
		t.Errorf("NewSpec should not return an errors but returns %s", err)
	}
	if len(*spec.Resources) != 112 {
		t.Errorf("Spec should contain %d resources but contains %d", 112, len(*spec.Resources))
	}
}

func TestGetResourceV119(t *testing.T) {
	spec, err := kubernetes.NewSpec("../../api/v1.19/swagger.json")
	if err != nil {
		t.Errorf("NewSpec should not return an errors but returns %s", err)
	}
	v1 := newAPIVersionAssert(t, "v1")
	_, res := spec.GetResource("", *v1, "Pod", false)
	if res.Description != "Pod is a collection of containers that can run on a host. This resource is created by clients and scheduled onto hosts." {
		t.Error("Error getting definition of Pod")
	}
}
