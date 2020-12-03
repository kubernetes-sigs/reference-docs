package kubernetes_test

import (
	"testing"

	"github.com/feloy/kubernetes-api-reference/pkg/kubernetes"
)

func TestAPIGroupReplaces(t *testing.T) {
	tests := []struct {
		Group1   kubernetes.APIGroup
		Group2   kubernetes.APIGroup
		Expected bool
	}{
		{"networking", "extensions", true},
		{"events.k8s.io", "", true},
	}

	for _, test := range tests {
		result := test.Group1.Replaces(test.Group2)
		if result != test.Expected {
			t.Errorf("%s replaces %s: expected %v but got %v", test.Group1, test.Group2, test.Expected, result)
		}
	}
}
