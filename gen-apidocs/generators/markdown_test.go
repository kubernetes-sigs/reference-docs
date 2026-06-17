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

const (
	// Match a real API group so fixture mirrors what auto-detect produces.
	testCategoryName = "Apps"
	testCategorySlug = "apps"
)

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
	// We only escape `<`; `>` and other characters pass through unchanged.
	if got := escape("a <b> c"); got != `a \<b> c` {
		t.Errorf("escape: got %q", got)
	}
	if got := escape("no change"); got != "no change" {
		t.Errorf("escape: got %q", got)
	}
	if got := escape("Allowed values.  Possible enum values:  - `\"A\"` first.  - `\"B\"` second."); got != "Allowed values.<br/><br/>Possible enum values:<br/> - `\"A\"` first.<br/> - `\"B\"` second." {
		t.Errorf("escape enum list: got %q", got)
	}
}

func TestKebabCase(t *testing.T) {
	cases := map[string]string{
		testCategoryName:       testCategorySlug,
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

func TestOperationSlug(t *testing.T) {
	cases := map[string]string{
		"listCoreV1Pod":                             "listcorev1pod",
		"readAppsV1NamespacedDeployment":            "readappsv1namespaceddeployment",
		"watchCore.V1.Pod":                          "watchcore-v1-pod",
		"Some/Weird ID":                             "some-weird-id",
	}
	for in, want := range cases {
		if got := operationSlug(in); got != want {
			t.Errorf("operationSlug(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestWriteResourceGolden(t *testing.T) {
	m, cleanup := newTestWriter(t)
	defer cleanup()

	// The category dir mimicks what WriteResourceCategory would have created.
	if err := os.MkdirAll(filepath.Join(m.OutputDir, testCategorySlug), 0755); err != nil {
		t.Fatal(err)
	}
	m.currentCategory = mdCategory{name: testCategoryName, slug: testCategorySlug}

	r := fabricateDeploymentResource()
	if err := m.WriteResource(r); err != nil {
		t.Fatalf("WriteResource: %v", err)
	}

	compareWithGolden(t,
		filepath.Join(m.OutputDir, testCategorySlug, "deployment-v1.md"),
		"testdata/deployment-v1.golden.md")
}

func TestWriteOperationGolden(t *testing.T) {
	m, cleanup := newTestWriter(t)
	defer cleanup()

	if err := os.MkdirAll(filepath.Join(m.OutputDir, "operations"), 0755); err != nil {
		t.Fatal(err)
	}

	o := fabricateOperation()
	if err := m.WriteOperation(o); err != nil {
		t.Fatalf("WriteOperation: %v", err)
	}

	compareWithGolden(t,
		filepath.Join(m.OutputDir, "operations", "listcorev1pod.md"),
		"testdata/operation-list.golden.md")
}

func TestWriteOperationGoldenHugoMode(t *testing.T) {
	m, cleanup := newTestWriter(t)
	defer cleanup()
	m.HugoMode = true

	if err := os.MkdirAll(filepath.Join(m.OutputDir, "operations"), 0755); err != nil {
		t.Fatal(err)
	}

	o := fabricateOperation()
	if err := m.WriteOperation(o); err != nil {
		t.Fatalf("WriteOperation: %v", err)
	}

	compareWithGolden(t,
		filepath.Join(m.OutputDir, "operations", "listcorev1pod.md"),
		"testdata/operation-list-hugo.golden.md")
}

func TestResolveType(t *testing.T) {
	m := &MarkdownWriter{linkMap: map[string]linkInfo{}}
	m.linkResources([]api.ResourceCategory{
		{
			Name: testCategoryName,
			Resources: api.Resources{
				{Name: "Deployment", Definition: &api.Definition{Name: "Deployment", Version: api.ApiVersion("v1")}},
			},
		},
	})
	m.linkDefinitions(map[string]*api.Definition{
		"io.k8s.api.apps.v1.DeploymentSpec": {
			Name: "DeploymentSpec", Group: api.ApiGroup("apps"), Version: api.ApiVersion("v1"),
		},
		"io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta": {
			Name: "ObjectMeta", Group: api.ApiGroup("meta"), Version: api.ApiVersion("v1"),
		},
	})

	cases := []struct {
		typeName, currentCategory, want string
	}{
		// From the Deployment resource page, DeploymentSpec lives in definitions/
		{"DeploymentSpec", testCategorySlug, "../definitions/deployment-spec-v1-apps#DeploymentSpec"},
		// ObjectMeta same — different dir
		{"ObjectMeta", testCategorySlug, "../definitions/object-meta-v1-meta#ObjectMeta"},
		// Self-reference (inside the definitions dir)
		{"ObjectMeta", "definitions", "object-meta-v1-meta#ObjectMeta"},
		// Primitive / unknown → empty string
		{"string", testCategorySlug, ""},
		{"UnknownThing", testCategorySlug, ""},
		// Resource path: Deployment lives under apps/
		{"Deployment", "definitions", "../apps/deployment-v1#Deployment"},
	}
	for _, c := range cases {
		if got := m.resolveType(c.typeName, c.currentCategory); got != c.want {
			t.Errorf("resolveType(%q, %q) = %q, want %q", c.typeName, c.currentCategory, got, c.want)
		}
	}
}

func TestRecordLinkPrecedence(t *testing.T) {
	m := &MarkdownWriter{linkMap: map[string]linkInfo{}}

	// Pretend Deployment is first seen as a standalone definition...
	m.recordLink("Deployment", "definitions", "deployment-v1-apps", api.ApiVersion("v1"))
	if got := m.linkMap["Deployment"].Category; got != "definitions" {
		t.Fatalf("initial category = %q, want definitions", got)
	}

	// ...then as a resource. Resource should win.
	m.recordLink("Deployment", "apps", "deployment-v1", api.ApiVersion("v1"))
	if got := m.linkMap["Deployment"].Category; got != "apps" {
		t.Errorf("after resource record, category = %q, want apps", got)
	}

	// A later standalone-definition record for the same name must not clobber.
	m.recordLink("Deployment", "definitions", "deployment-v1-apps", api.ApiVersion("v1"))
	if got := m.linkMap["Deployment"].Category; got != "apps" {
		t.Errorf("after late definitions record, category = %q, want apps", got)
	}

	// Version bump within same bucket wins.
	m.recordLink("HPA", "autoscaling", "horizontalpodautoscaler-v1", api.ApiVersion("v1"))
	m.recordLink("HPA", "autoscaling", "horizontalpodautoscaler-v2", api.ApiVersion("v2"))
	if got := m.linkMap["HPA"].Version; string(got) != "v2" {
		t.Errorf("HPA version = %q, want v2", got)
	}

	// Older version must not clobber newer.
	m.recordLink("HPA", "autoscaling", "horizontalpodautoscaler-v1", api.ApiVersion("v1"))
	if got := m.linkMap["HPA"].Version; string(got) != "v2" {
		t.Errorf("after old version, HPA version = %q, want v2", got)
	}
}

// TestClassifyDefinitions verifies the BFS home-finder on a synthetic
// reference graph: types with a unique closest top-level inline there,
// types with ties at the minimum distance stay standalone.
func TestClassifyDefinitions(t *testing.T) {
	mkDef := func(name string, inToc bool) *api.Definition {
		return &api.Definition{
			Name:    name,
			Group:   api.ApiGroup("test"),
			Version: api.ApiVersion("v1"),
			Kind:    api.ApiKind(name),
			InToc:   inToc,
		}
	}

	// Reference graph (parent --refs--> child; AppearsIn is reverse).
	//
	//   Pod (InToc) --> PodSpec --> Container
	//                          --> Volume --> AzureDiskVolumeSource
	//                          <-- PodTemplateSpec (shared)
	//   Deployment (InToc) --> DeploymentSpec --> PodTemplateSpec
	//   PodTemplate (InToc) --> PodTemplateSpec
	//   ObjectMeta -- referenced directly by Pod, Deployment, PodTemplate
	pod := mkDef("Pod", true)
	deployment := mkDef("Deployment", true)
	podTemplate := mkDef("PodTemplate", true)
	podSpec := mkDef("PodSpec", false)
	deploymentSpec := mkDef("DeploymentSpec", false)
	podTemplateSpec := mkDef("PodTemplateSpec", false)
	container := mkDef("Container", false)
	volume := mkDef("Volume", false)
	azureDisk := mkDef("AzureDiskVolumeSource", false)
	objectMeta := mkDef("ObjectMeta", false)

	// Populate AppearsIn (= "who references me?").
	podSpec.AppearsIn = api.SortDefinitionsByName{pod, podTemplateSpec}
	deploymentSpec.AppearsIn = api.SortDefinitionsByName{deployment}
	podTemplateSpec.AppearsIn = api.SortDefinitionsByName{podTemplate, deploymentSpec}
	container.AppearsIn = api.SortDefinitionsByName{podSpec}
	volume.AppearsIn = api.SortDefinitionsByName{podSpec}
	azureDisk.AppearsIn = api.SortDefinitionsByName{volume}
	objectMeta.AppearsIn = api.SortDefinitionsByName{pod, deployment, podTemplate}

	all := map[string]*api.Definition{}
	for _, d := range []*api.Definition{
		pod, deployment, podTemplate,
		podSpec, deploymentSpec, podTemplateSpec,
		container, volume, azureDisk, objectMeta,
	} {
		all[d.Key()] = d
	}

	m := &MarkdownWriter{Config: &api.Config{
		Definitions: api.Definitions{All: all},
	}}
	got := m.classifyDefinitions()

	cases := []struct {
		name       string
		def        *api.Definition
		wantMode   classifyMode
		wantParent *api.Definition // only checked when wantMode == classifyInline
	}{
		// Container's only chain hits Pod at dist 2 (Container -> PodSpec -> Pod)
		// while Deployment / PodTemplate sit at dist 5+. Pod wins uniquely.
		{"Container inlines into Pod", container, classifyInline, pod},
		// AzureDisk is one hop deeper than Container; same winner.
		{"AzureDiskVolumeSource inlines into Pod", azureDisk, classifyInline, pod},
		// Volume is shared via PodSpec only; same winner as Container.
		{"Volume inlines into Pod", volume, classifyInline, pod},
		// ObjectMeta is at distance 1 from three InToc resources → tie → standalone.
		{"ObjectMeta is standalone", objectMeta, classifyStandalone, nil},
		// Pod itself is InToc → not classified (own resource page).
		{"Pod is skip", pod, classifySkip, nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cls := got[c.def.Key()]
			if cls.Mode != c.wantMode {
				t.Fatalf("Mode = %v, want %v", cls.Mode, c.wantMode)
			}
			if c.wantMode == classifyInline && cls.InlineInto != c.wantParent {
				t.Fatalf("InlineInto = %v, want %v", cls.InlineInto, c.wantParent)
			}
		})
	}
}

// TestLinkMapAliasesInlinedChildren verifies that a definition classified
// as inline gets a linkMap entry that points at its parent resource page,
// not at a standalone definitions/ file (which won't exist after Stage 2).
func TestLinkMapAliasesInlinedChildren(t *testing.T) {
	pod := &api.Definition{
		Name: "Pod", Group: api.ApiGroup("core"), Version: api.ApiVersion("v1"),
		Kind: api.ApiKind("Pod"), InToc: true,
	}
	container := &api.Definition{
		Name: "Container", Group: api.ApiGroup("core"), Version: api.ApiVersion("v1"),
		Kind: api.ApiKind("Container"),
	}

	m := &MarkdownWriter{
		linkMap: map[string]linkInfo{},
		classifications: map[string]defClassification{
			container.Key(): {Mode: classifyInline, InlineInto: pod},
		},
		inlinedByParent: map[string][]*api.Definition{
			pod.Key(): {container},
		},
	}

	m.linkResources([]api.ResourceCategory{
		{
			Name:      "Workloads",
			Resources: api.Resources{{Name: "Pod", Definition: pod}},
		},
	})

	// Inlined child must alias to the parent's slug + filename.
	got, ok := m.linkMap["Container"]
	if !ok {
		t.Fatal("linkMap missing entry for inlined Container")
	}
	if got.Category != "workloads" || got.Filename != "pod-v1" {
		t.Errorf("Container linkInfo = {%q,%q}, want {workloads,pod-v1}", got.Category, got.Filename)
	}
}

// TestLinkDefinitionsSkipsInlined verifies that an inlined type does not get
// a stale definitions/ entry — only the parent-aliased entry from linkResources.
func TestLinkDefinitionsSkipsInlined(t *testing.T) {
	azureDisk := &api.Definition{
		Name: "AzureDiskVolumeSource", Group: api.ApiGroup("core"),
		Version: api.ApiVersion("v1"), Kind: api.ApiKind("AzureDiskVolumeSource"),
	}
	objectMeta := &api.Definition{
		Name: "ObjectMeta", Group: api.ApiGroup("meta"),
		Version: api.ApiVersion("v1"), Kind: api.ApiKind("ObjectMeta"),
	}

	m := &MarkdownWriter{
		linkMap: map[string]linkInfo{},
		classifications: map[string]defClassification{
			azureDisk.Key():   {Mode: classifyInline}, // InlineInto irrelevant for this test
			objectMeta.Key():  {Mode: classifyStandalone},
		},
	}

	m.linkDefinitions(map[string]*api.Definition{
		azureDisk.Key():  azureDisk,
		objectMeta.Key(): objectMeta,
	})

	if _, exists := m.linkMap["AzureDiskVolumeSource"]; exists {
		t.Error("inlined AzureDiskVolumeSource got a definitions/ entry; should have been skipped")
	}
	if got, ok := m.linkMap["ObjectMeta"]; !ok || got.Category != "definitions" {
		t.Errorf("ObjectMeta should keep its standalone definitions/ entry; got = %+v ok=%v", got, ok)
	}
}

// --- fixture helpers ---

func newTestWriter(t *testing.T) (*MarkdownWriter, func()) {
	t.Helper()
	tmp := t.TempDir()
	prevBuildDir := api.BuildDir
	api.BuildDir = tmp
	m := &MarkdownWriter{
		Config:    &api.Config{SpecTitle: "Test", SpecVersion: "v1.0.0"},
		OutputDir: filepath.Join(tmp, "markdown"),
		linkMap:   make(map[string]linkInfo),
	}
	if err := os.MkdirAll(m.OutputDir, 0755); err != nil {
		t.Fatal(err)
	}
	return m, func() { api.BuildDir = prevBuildDir }
}

func fabricateDeploymentResource() *api.Resource {
	return &api.Resource{
		Name: "Deployment",
		Definition: &api.Definition{
			Name:                    "Deployment",
			Group:                   api.ApiGroup("apps"),
			GroupFullName:           "apps",
			Version:                 api.ApiVersion("v1"),
			Kind:                    api.ApiKind("Deployment"),
			DescriptionWithEntities: "Deployment enables declarative updates for Pods and ReplicaSets.",
			SwaggerKey:              "io.k8s.api.apps.v1.Deployment",
			Fields: api.Fields{
				{Name: "apiVersion", Type: "string", Description: "APIVersion defines the versioned schema of this representation of an object."},
				{Name: "kind", Type: "string", Description: "Kind is a string value representing the REST resource."},
				{Name: "metadata", Type: "ObjectMeta", Description: "Standard object's metadata."},
				{Name: "spec", Type: "DeploymentSpec", Description: "Specification of the desired behavior of the Deployment."},
				{Name: "status", Type: "DeploymentStatus", Description: "Most recently observed status of the Deployment."},
			},
		},
	}
}

func fabricateOperation() *api.Operation {
	return &api.Operation{
		ID:         "listCoreV1Pod",
		Type:       api.OperationType{Name: "List Pods"},
		Path:       "/api/v1/namespaces/{namespace}/pods",
		HttpMethod: "GET",
		PathParams: api.Fields{
			{Name: "namespace", Type: "string", Description: "object name and auth scope, such as for teams and users"},
		},
		QueryParams: api.Fields{
			{Name: "watch", Type: "boolean", Description: "Watch for changes to the described resources."},
		},
		HttpResponses: api.HttpResponses{
			{Code: "200", Field: api.Field{Type: "PodList", Description: "OK"}},
			{Code: "401", Field: api.Field{Description: "Unauthorized"}},
		},
	}
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
