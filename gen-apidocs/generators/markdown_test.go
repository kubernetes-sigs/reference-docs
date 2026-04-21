/*
Copyright 2024 The Kubernetes Authors.

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

package generators

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kubernetes-sigs/reference-docs/gen-apidocs/generators/api"
)

// update controls golden-file regeneration. Run `go test -update` to
// rewrite any *.golden files from current output.
var update = flag.Bool("update", false, "rewrite *.golden files from test output")

func TestAnchor(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Deployment", "Deployment"},
		{"Pod Spec", "Pod-Spec"},
		{"v1.Container", "v1-Container"},
		{"foo__bar  baz", "foo-bar-baz"},
		{"---leading-and-trailing---", "leading-and-trailing"},
	}
	for _, c := range cases {
		if got := anchor(c.in); got != c.want {
			t.Errorf("anchor(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestEscape(t *testing.T) {
	// We only escape `<` — `>` and other characters pass through to match
	// gen-resourcesdocs chapter.tmpl behaviour.
	if got := escape("a <b> c"); got != `a \<b> c` {
		t.Errorf("escape: got %q", got)
	}
	if got := escape("no change"); got != "no change" {
		t.Errorf("escape: got %q", got)
	}
}

func TestKebabCase(t *testing.T) {
	cases := map[string]string{
		"Workloads APIs":       "workloads-apis",
		"Service Discovery":    "service-discovery",
		"Cluster - Admin":      "cluster-admin",
		"  Leading trailing  ": "leading-trailing",
	}
	for in, want := range cases {
		if got := kebabCase(in); got != want {
			t.Errorf("kebabCase(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestGroupVersionString(t *testing.T) {
	cases := []struct {
		group   string
		version api.ApiVersion
		want    string
	}{
		{"", "v1", "v1"},
		{"core", "v1", "v1"},
		{"apps", "v1", "apps/v1"},
		{"batch", "v1beta1", "batch/v1beta1"},
		{"apiextensions.k8s.io", "v1", "apiextensions.k8s.io/v1"},
	}
	for _, c := range cases {
		if got := groupVersionString(c.group, c.version); got != c.want {
			t.Errorf("groupVersionString(%q, %q) = %q, want %q", c.group, c.version, got, c.want)
		}
	}
}

func TestWritePipeTable(t *testing.T) {
	var buf bytes.Buffer
	writePipeTable(&buf, []string{"Group", "Versions"}, func(row func(cells ...string)) {
		row("`apps`", "`v1`")
		row("`batch`", "`v1, v1beta1`")
	})
	want := "| Group | Versions |\n" +
		"| --- | --- |\n" +
		"| `apps` | `v1` |\n" +
		"| `batch` | `v1, v1beta1` |\n"
	if got := buf.String(); got != want {
		t.Errorf("writePipeTable mismatch:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

// TestWriteResource_Golden exercises WriteResource end-to-end once it's
// implemented. It's skipped in PR 1 so CI stays green while the stub is
// in place. Unskip when WriteResource lands.
func TestWriteResource_Golden(t *testing.T) {
	t.Skip("pending WriteResource implementation — remove Skip when the body lands")

	// Reference scaffold for the implementor:
	//
	// tmp := t.TempDir()
	// api.BuildDir = tmp
	// m := &MarkdownWriter{
	//     Config:    minimalConfig(),
	//     OutputDir: filepath.Join(tmp, "markdown"),
	//     linkMap:   make(map[string]linkInfo),
	// }
	// os.MkdirAll(filepath.Join(m.OutputDir, "workloads-apis"), 0755)
	// m.currentCategory = mdCategory{name: "Workloads APIs", slug: "workloads-apis"}
	//
	// r := fabricateDeploymentResource()
	// if err := m.WriteResource(r); err != nil {
	//     t.Fatal(err)
	// }
	// compareWithGolden(t, filepath.Join(m.OutputDir, "workloads-apis", "deployment-v1.md"),
	//     "testdata/deployment-v1.golden.md")
}

// --- helpers ---

func mustContainInOrder(t *testing.T, s string, substrs ...string) {
	t.Helper()
	i := 0
	for _, sub := range substrs {
		idx := strings.Index(s[i:], sub)
		if idx < 0 {
			t.Fatalf("expected substring %q after position %d, got:\n%s", sub, i, s)
		}
		i += idx + len(sub)
	}
}

// compareWithGolden is used by the WriteResource test above once the
// implementation lands. Supports `go test -update` for easy regeneration.
func compareWithGolden(t *testing.T, gotPath, goldenPath string) { //nolint:unused // scaffold for future tests
	t.Helper()
	got, err := os.ReadFile(gotPath)
	if err != nil {
		t.Fatalf("read got: %v", err)
	}
	if *update {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(goldenPath, got, 0644); err != nil {
			t.Fatalf("write golden: %v", err)
		}
		return
	}
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden: %v (run `go test -update` to create it)", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("golden mismatch for %s\n--- got ---\n%s\n--- want ---\n%s",
			filepath.Base(gotPath), got, want)
	}
}
