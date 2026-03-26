package config

import "testing"

func TestGetEscapedFirstPhrase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single sentence with period",
			input:    "Foo is a bar.",
			expected: "Foo is a bar.",
		},
		{
			name:     "X.509 abbreviation not split",
			input:    "ClusterTrustBundle is a cluster-scoped container for X.509 trust anchors (root certificates).",
			expected: "ClusterTrustBundle is a cluster-scoped container for X.509 trust anchors (root certificates).",
		},
		{
			name:     "paragraph boundary takes first sentence",
			input:    "ClusterTrustBundle is a cluster-scoped container for X.509 trust anchors (root certificates).\n\nClusterTrustBundle objects are considered readable by any authenticated user.",
			expected: "ClusterTrustBundle is a cluster-scoped container for X.509 trust anchors (root certificates).",
		},
		{
			name:     "e.g. abbreviation not split",
			input:    "ServiceCIDR defines a range of IP addresses using CIDR format (e.g. 192.168.0.0/24 or 2001:db2::/64).",
			expected: "ServiceCIDR defines a range of IP addresses using CIDR format (e.g. 192.168.0.0/24 or 2001:db2::/64).",
		},
		{
			name:     "two sentences split on uppercase",
			input:    "First sentence. Second sentence.",
			expected: "First sentence.",
		},
		{
			name:     "double space between sentences",
			input:    "RoleBinding references a role, but does not contain it.  It can reference a Role in the same namespace.",
			expected: "RoleBinding references a role, but does not contain it.",
		},
		{
			name:     "sentence boundary before list continuation",
			input:    "StatefulSet represents a set of pods with consistent identities. Identities are defined as:\n  - Network: A single stable DNS and hostname.\n  - Storage: As many VolumeClaims as requested.",
			expected: "StatefulSet represents a set of pods with consistent identities.",
		},
		{
			name:     "multi-paragraph with numbered list",
			input:    "CertificateSigningRequest objects provide a mechanism to obtain x509 certificates.\n\nKubelets use this API to obtain:\n 1. client certificates",
			expected: "CertificateSigningRequest objects provide a mechanism to obtain x509 certificates.",
		},
		{
			name:     "no period adds one",
			input:    "ResourceQuota sets aggregate quota restrictions enforced per namespace",
			expected: "ResourceQuota sets aggregate quota restrictions enforced per namespace.",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "description with quotes",
			input:    `She said "hello". Then left.`,
			expected: `She said \"hello\".`,
		},
		{
			name:     "no sentence boundary returns whole string",
			input:    "ServiceAccount binds together: * a name * a principal * a set of secrets",
			expected: "ServiceAccount binds together: * a name * a principal * a set of secrets.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getEscapedFirstPhrase(tt.input)
			if got != tt.expected {
				t.Errorf("getEscapedFirstPhrase(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
