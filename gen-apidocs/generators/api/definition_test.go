/*
Copyright 2026 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"strings"
	"testing"

	"github.com/go-openapi/spec"
)

// endpointSubsetDescription mirrors the real doc comment on
// k8s.io/api/core/v1's EndpointSubset: a summary sentence, a blank line, an
// indented struct-literal example, another blank line, then more prose. This
// is the shape that produced the garbled front-matter ".lead" reported in
// kubernetes/website#56418.
const endpointSubsetDescription = "Endpoints is a collection of endpoints that implement the actual service. Example:\n" +
	"\n" +
	"\t Name: \"mysvc\",\n" +
	"\t Subsets: [\n" +
	"\t   {\n" +
	"\t     Addresses: [{\"ip\": \"10.10.1.1\"}, {\"ip\": \"10.10.2.2\"}],\n" +
	"\t     Ports: [{\"name\": \"a\", \"port\": 8675}, {\"name\": \"b\", \"port\": 309}]\n" +
	"\t   },\n" +
	"\t ]\n" +
	"\n" +
	"Endpoints is a legacy API and does not contain information about all Service features."

func TestDefinitionSummaryVsDescription(t *testing.T) {
	d := &Definition{schema: spec.Schema{
		SchemaProps: spec.SchemaProps{Description: endpointSubsetDescription},
	}}

	wantSummary := "Endpoints is a collection of endpoints that implement the actual service. Example:"
	if got := d.Summary(); got != wantSummary {
		t.Errorf("Summary() = %q, want %q", got, wantSummary)
	}
	if strings.Contains(d.Summary(), "\n") || strings.Contains(d.Summary(), "\t") {
		t.Errorf("Summary() must not contain embedded newlines/tabs, got %q", d.Summary())
	}
	if strings.Contains(d.Summary(), "mysvc") {
		t.Errorf("Summary() must not include the embedded example, got %q", d.Summary())
	}

	// Description() (used for the page body) must still carry the full text,
	// including the example — only the front-matter Summary() is trimmed.
	if got := d.Description(); got != endpointSubsetDescription {
		t.Errorf("Description() = %q, want the untouched original description", got)
	}
	if !strings.Contains(d.Description(), "mysvc") {
		t.Errorf("Description() must still include the embedded example")
	}
}

func TestDefinitionSummarySingleParagraph(t *testing.T) {
	// A description with no embedded example (the common case) should pass
	// through Summary() unchanged aside from trimming.
	d := &Definition{schema: spec.Schema{
		SchemaProps: spec.SchemaProps{Description: "PodSpec describes how the pod will look."},
	}}
	want := "PodSpec describes how the pod will look."
	if got := d.Summary(); got != want {
		t.Errorf("Summary() = %q, want %q", got, want)
	}
}
