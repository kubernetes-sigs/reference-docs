package kubernetes_test

import (
	"testing"

	"github.com/feloy/kubernetes-api-reference/pkg/kubernetes"
)

func TestResourceLessThan(t *testing.T) {
	v1 := newAPIVersionAssert(t, "v1")
	v1beta1 := newAPIVersionAssert(t, "v1beta1")
	v2alpha1 := newAPIVersionAssert(t, "v2alpha1")

	tests := []struct {
		R1       kubernetes.Resource
		R2       kubernetes.Resource
		Expected bool
	}{
		// General case
		{
			R1: kubernetes.Resource{
				Key: "key1",
				GVKExtension: kubernetes.GVKExtension{
					Group:   kubernetes.APIGroup("apps"),
					Version: *v1,
				},
			},
			R2: kubernetes.Resource{
				Key: "key2",
				GVKExtension: kubernetes.GVKExtension{
					Group:   kubernetes.APIGroup("apps"),
					Version: *v1beta1,
				},
			},
			Expected: true,
		},
		// Cronjob resource in v1.18
		{
			R1: kubernetes.Resource{
				Key: "key1",
				GVKExtension: kubernetes.GVKExtension{
					Group:   kubernetes.APIGroup("batch"),
					Version: *v2alpha1,
				},
			},
			R2: kubernetes.Resource{
				Key: "key2",
				GVKExtension: kubernetes.GVKExtension{
					Group:   kubernetes.APIGroup("batch"),
					Version: *v1beta1,
				},
			},
			Expected: true,
		},
		// Event resource in v1.18
		{
			R1: kubernetes.Resource{
				Key: "key1",
				GVKExtension: kubernetes.GVKExtension{
					Group:   kubernetes.APIGroup(""),
					Version: *v1,
				},
			},
			R2: kubernetes.Resource{
				Key: "key2",
				GVKExtension: kubernetes.GVKExtension{
					Group:   kubernetes.APIGroup("events.k8s.io"),
					Version: *v1,
				},
			},
			Expected: false,
		},
		// Ingress resource in v1.18
		{
			R1: kubernetes.Resource{
				Key: "key1",
				GVKExtension: kubernetes.GVKExtension{
					Group:   kubernetes.APIGroup("networking.k8s.io"),
					Version: *v1beta1,
				},
			},
			R2: kubernetes.Resource{
				Key: "key2",
				GVKExtension: kubernetes.GVKExtension{
					Group:   kubernetes.APIGroup("extensions"),
					Version: *v1beta1,
				},
			},
			Expected: true,
		},
	}

	for _, test := range tests {
		result := test.R1.LessThan(&test.R2)
		if result != test.Expected {
			t.Errorf("%s < %s: expected %v but got %v", test.R1.GetGV(), test.R2.GetGV(), test.Expected, result)
		}
	}
}

func TestResourceGetGV(t *testing.T) {
	v1 := newAPIVersionAssert(t, "v1")

	tests := []struct {
		Input    kubernetes.Resource
		Expected string
	}{
		{
			Input: kubernetes.Resource{
				Key: "key1",
				GVKExtension: kubernetes.GVKExtension{
					Group:   kubernetes.APIGroup("apps"),
					Version: *v1,
				},
			},
			Expected: "apps/v1",
		},
		{
			Input: kubernetes.Resource{
				Key: "key1",
				GVKExtension: kubernetes.GVKExtension{
					Group:   kubernetes.APIGroup(""),
					Version: *v1,
				},
			},
			Expected: "v1",
		},
		{
			Input: kubernetes.Resource{
				Key: "key1",
				GVKExtension: kubernetes.GVKExtension{
					Group:   kubernetes.APIGroup("storage.k8s.io"),
					Version: *v1,
				},
			},
			Expected: "storage.k8s.io/v1",
		},
	}

	for _, test := range tests {
		result := test.Input.GetGV()
		if result != test.Expected {
			t.Errorf("%#v: Expected %s but got %s", test.Input, test.Expected, result)
		}
	}
}

func TestResourceAdd(t *testing.T) {
	v1 := newAPIVersionAssert(t, "v1")

	resources := kubernetes.ResourceMap{}
	resources.Add(&kubernetes.Resource{
		Key: "key1",
		GVKExtension: kubernetes.GVKExtension{
			Group:   kubernetes.APIGroup("extensions"),
			Version: *v1,
			Kind:    kubernetes.APIKind("Kind1"),
		},
	})
	resources.Add(&kubernetes.Resource{
		Key: "key1",
		GVKExtension: kubernetes.GVKExtension{
			Group:   kubernetes.APIGroup("apps"),
			Version: *v1,
			Kind:    kubernetes.APIKind("Kind1"),
		},
	})
	resources.Add(&kubernetes.Resource{
		Key: "key1",
		GVKExtension: kubernetes.GVKExtension{
			Group:   kubernetes.APIGroup("apps"),
			Version: *v1,
			Kind:    kubernetes.APIKind("Kind2"),
		},
	})
	if len(resources) != 2 {
		t.Errorf("Len of resources should be %d but is %d", 2, len(resources))
	}
	if _, ok := resources["Kind1"]; !ok {
		t.Errorf("Key Kind1 should exist")
	}
	if _, ok := resources["Kind2"]; !ok {
		t.Errorf("Key Kind2 should exist")
	}
	if len(resources["Kind1"]) != 2 {
		t.Errorf("Len of versions for Kind1 should be %d but is %d", 2, len(resources["Kind1"]))
	}
	if len(resources["Kind2"]) != 1 {
		t.Errorf("Len of versions for Kind2 should be %d but is %d", 1, len(resources["Kind2"]))
	}
	if resources["Kind1"][0].Group != "apps" {
		t.Errorf("Recent version for Kind1 should be %s but is %s", "apps", resources["Kind1"][0].Group)
	}
	if resources["Kind1"][1].Group != "extensions" {
		t.Errorf("Previous version for Kind1 should be %s but is %s", "extensions", resources["Kind1"][1].Group)
	}
}

func TestGoImportPrefix(t *testing.T) {
	tests := []struct {
		Key      kubernetes.Key
		Expected string
	}{
		{
			Key:      "io.k8s.api.core.Pod",
			Expected: "k8s.io/api/core",
		},
	}

	for _, test := range tests {
		result := test.Key.GoImportPrefix()
		if result != test.Expected {
			t.Errorf("%s: Expected %s but got %s", test.Key, test.Expected, result)
		}
	}
}

func TestRemoveResourceName(t *testing.T) {
	tests := []struct {
		Key      kubernetes.Key
		Expected string
	}{
		{
			Key:      "io.k8s.api.core.PodSpec",
			Expected: "io.k8s.api.core",
		},
	}

	for _, test := range tests {
		result := test.Key.RemoveResourceName()
		if result != test.Expected {
			t.Errorf("%s: Expected %s but got %s", test.Key, test.Expected, result)
		}
	}
}
