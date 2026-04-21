/*
Copyright 2016 The Kubernetes Authors.

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
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/kubernetes-sigs/reference-docs/gen-apidocs/generators/api"
)

// MarkdownWriter emits Hugo-compatible markdown under:
//
//	_index.md                              — top-level part listing
//	<category-slug>/_index.md              — per-category index
//	<category-slug>/<resource>-<ver>.md    — per-resource page
//	definitions/<name>-<ver>-<group>.md    — standalone definitions
//	operations/<op-id>.md                  — orphaned operations
type MarkdownWriter struct {
	Config          *api.Config
	OutputDir       string
	currentCategory mdCategory
	resourceWeight  int
	categoryWeight  int

	// linkMap is populated during render for the PR 2 cross-reference pass.
	linkMap map[string]linkInfo

	toc []*mdTOCItem

	// finalized guards against Finalize being called twice by GenerateFiles.
	finalized bool
}

type mdCategory struct {
	name string
	slug string
}

type mdTOCItem struct {
	title    string
	path     string
	weight   int
	children []*mdTOCItem
}

type linkInfo struct {
	Category string
	Filename string
	Anchor   string
}

const hugoIndex = "_index.md"

// Header page titles — kept as named constants so the TOC-ordering
// switch and the page emitters never drift.
const (
	titleOverview    = "Overview"
	titleAPIGroups   = "API Groups"
	titleDefinitions = "Definitions"
	titleOperations  = "Operations"
	titleOldVersions = "Old API Versions"
)

var _ DocWriter = (*MarkdownWriter)(nil)

// anchorRegex must stay in sync with the Sprig regex in
// gen-resourcesdocs/templates/chapter.tmpl; external k/website links rely
// on the exact anchor format.
var anchorRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

//go:embed templates/resource.tmpl
var resourceTemplateSrc string

// Template funcs: q = YAML-safe quote (frontmatter strings);
// md = markdown `<` escape (body text). Descriptions are passed raw
// so the template picks the right escape per site.
var resourceTemplate = template.Must(template.New("resource").Funcs(template.FuncMap{
	"q":  strconv.Quote,
	"md": escape,
}).Parse(resourceTemplateSrc))

type resourcePage struct {
	APIVersion  string
	Kind        string
	Import      string
	Title       string
	Weight      int
	Anchor      string
	Description string
	Fields      []templateField
	Operations  []templateOperation
}

type templateField struct {
	Name          string
	Type          string
	Description   string
	Required      bool
	ConstValue    string // non-empty for fields with a fixed value (apiVersion, kind)
	PatchStrategy string // x-kubernetes-patch-strategy — e.g. "merge", "retainKeys"
	PatchMergeKey string // x-kubernetes-patch-merge-key — e.g. "name"
}

type templateOperation struct {
	Verb        string
	Title       string
	Method      string
	Path        string
	PathParams  []templateParam
	QueryParams []templateParam
	BodyParams  []templateParam
	Responses   []templateResponse
}

type templateParam struct {
	Name        string
	Type        string
	Description string
}

type templateResponse struct {
	Code        string
	Type        string
	Description string
}

func NewMarkdownWriter(config *api.Config, copyright, title string) DocWriter {
	outputDir := filepath.Join(api.BuildDir, "markdown")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "MarkdownWriter: failed to create output dir %s: %v\n", outputDir, err)
	}
	return &MarkdownWriter{
		Config:    config,
		OutputDir: outputDir,
		linkMap:   make(map[string]linkInfo),
	}
}

func (m *MarkdownWriter) Extension() string {
	return ".md"
}

func (m *MarkdownWriter) DefaultStaticContent(title string) string {
	return "# " + title + "\n"
}

func (m *MarkdownWriter) WriteOverview() error {
	if err := m.writeSection("_overview.md", "API Overview"); err != nil {
		return fmt.Errorf("markdown: overview: %w", err)
	}
	m.toc = append(m.toc, &mdTOCItem{
		title:  titleOverview,
		path:   "_overview.md",
		weight: m.nextCategoryWeight(),
	})
	return nil
}

func (m *MarkdownWriter) WriteAPIGroupVersions(gvs api.GroupVersions) error {
	path := filepath.Join(m.OutputDir, "_group_versions.md")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("markdown: group versions: %w", err)
	}
	defer f.Close()

	fmt.Fprintln(f, "# API Groups")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "The API Groups and their versions are summarized in the following table.")
	fmt.Fprintln(f)

	groups := make(api.ApiGroups, 0, len(gvs))
	for g := range gvs {
		groups = append(groups, api.ApiGroup(g))
	}
	sort.Sort(groups)

	writePipeTable(f, []string{"Group", "Versions"}, func(row func(cells ...string)) {
		for _, g := range groups {
			versions := gvs[g.String()]
			sort.Sort(versions)
			vs := make([]string, 0, len(versions))
			for _, v := range versions {
				vs = append(vs, v.String())
			}
			row("`"+g.String()+"`", "`"+strings.Join(vs, ", ")+"`")
		}
	})

	m.toc = append(m.toc, &mdTOCItem{
		title:  titleAPIGroups,
		path:   "_group_versions.md",
		weight: m.nextCategoryWeight(),
	})
	return nil
}

func (m *MarkdownWriter) WriteResourceCategory(name, file string) error {
	slug := kebabCase(name)
	dir := filepath.Join(m.OutputDir, slug)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("markdown: category dir: %w", err)
	}

	weight := m.nextCategoryWeight()
	indexPath := filepath.Join(dir, hugoIndex)
	f, err := os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("markdown: category index: %w", err)
	}
	defer f.Close()

	writeSectionFrontmatter(f, name, "", weight)

	if body := readOptionalSection(file + ".md"); body != "" {
		fmt.Fprintln(f, body)
	} else {
		fmt.Fprintf(f, "# %s\n", name)
	}

	m.currentCategory = mdCategory{name: name, slug: slug}
	m.resourceWeight = 0
	m.toc = append(m.toc, &mdTOCItem{
		title:  name,
		path:   filepath.Join(slug, hugoIndex),
		weight: weight,
	})
	return nil
}

func (m *MarkdownWriter) WriteDefinitionsOverview() error {
	if err := m.writeSection("_definitions.md", titleDefinitions); err != nil {
		return fmt.Errorf("markdown: definitions overview: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(m.OutputDir, "definitions"), 0755); err != nil {
		return fmt.Errorf("markdown: definitions dir: %w", err)
	}
	m.toc = append(m.toc, &mdTOCItem{
		title:  titleDefinitions,
		path:   "_definitions.md",
		weight: m.nextCategoryWeight(),
	})
	return nil
}

func (m *MarkdownWriter) WriteOrphanedOperationsOverview() error {
	if err := m.writeSection("_operations.md", titleOperations); err != nil {
		return fmt.Errorf("markdown: operations overview: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(m.OutputDir, "operations"), 0755); err != nil {
		return fmt.Errorf("markdown: operations dir: %w", err)
	}
	m.toc = append(m.toc, &mdTOCItem{
		title:  titleOperations,
		path:   "_operations.md",
		weight: m.nextCategoryWeight(),
	})
	return nil
}

func (m *MarkdownWriter) WriteOldVersionsOverview() error {
	if err := m.writeSection("_oldversions.md", titleOldVersions); err != nil {
		return fmt.Errorf("markdown: old versions overview: %w", err)
	}
	m.toc = append(m.toc, &mdTOCItem{
		title:  titleOldVersions,
		path:   "_oldversions.md",
		weight: m.nextCategoryWeight(),
	})
	return nil
}

func (m *MarkdownWriter) WriteResource(r *api.Resource) error {
	filename := fmt.Sprintf("%s-%s.md", strings.ToLower(r.Name), r.Definition.Version)
	path := filepath.Join(m.OutputDir, m.currentCategory.slug, filename)

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("markdown: resource %s: %w", r.Name, err)
	}
	defer f.Close()

	if err := resourceTemplate.Execute(f, m.buildResourcePage(r)); err != nil {
		return fmt.Errorf("markdown: resource %s body: %w", r.Name, err)
	}

	return nil
}

func (m *MarkdownWriter) buildResourcePage(r *api.Resource) resourcePage {
	page := m.buildDefinitionPage(r.Definition)
	for _, oc := range r.Definition.OperationCategories {
		for _, o := range oc.Operations {
			page.Operations = append(page.Operations, buildTemplateOperation(o))
		}
	}
	return page
}

// buildDefinitionPage fills in the resource-page fields common to both
// resource pages (which then add operations) and standalone definition
// pages (which don't).
func (m *MarkdownWriter) buildDefinitionPage(d *api.Definition) resourcePage {
	page := resourcePage{
		APIVersion:  groupVersionString(d.GroupFullName, d.Version),
		Kind:        d.Name,
		Import:      d.GoImportPath(),
		Title:       d.Name,
		Weight:      m.nextResourceWeight(),
		Anchor:      anchor(d.Name),
		Description: d.DescriptionWithEntities,
	}

	required := map[string]bool{}
	for _, name := range d.RequiredFields() {
		required[name] = true
	}

	for _, fld := range d.Fields {
		page.Fields = append(page.Fields, templateField{
			Name:          fld.Name,
			Type:          fld.Type,
			Description:   fld.Description,
			Required:      required[fld.Name],
			ConstValue:    constValueFor(fld.Name, page.APIVersion, page.Kind),
			PatchStrategy: fld.PatchStrategy,
			PatchMergeKey: fld.PatchMergeKey,
		})
	}

	return page
}

// constValueFor returns the fixed value of a definition field when one exists.
// Today only apiVersion and kind are constants on resource objects — the
// swagger doesn't tag them explicitly, but Kubernetes manifests always carry
// the same literal values for a given GVK, so we surface them as hints.
func constValueFor(fieldName, apiVersion, kind string) string {
	switch fieldName {
	case "apiVersion":
		return apiVersion
	case "kind":
		return kind
	}
	return ""
}

func buildTemplateOperation(o *api.Operation) templateOperation {
	op := templateOperation{
		Verb:   strings.ToLower(o.HttpMethod),
		Title:  o.Type.Name,
		Method: o.HttpMethod,
		Path:   o.Path,
	}

	convert := func(params api.Fields) []templateParam {
		out := make([]templateParam, 0, len(params))
		for _, p := range params {
			out = append(out, templateParam{
				Name:        p.Name,
				Type:        p.Type,
				Description: p.Description,
			})
		}
		return out
	}
	op.PathParams = convert(o.PathParams)
	op.QueryParams = convert(o.QueryParams)
	op.BodyParams = convert(o.BodyParams)

	responses := append(api.HttpResponses(nil), o.HttpResponses...)
	sort.Slice(responses, func(i, j int) bool {
		return responses[i].Code < responses[j].Code
	})
	for _, rsp := range responses {
		op.Responses = append(op.Responses, templateResponse{
			Code:        rsp.Code,
			Type:        rsp.Field.Type,
			Description: rsp.Field.Description,
		})
	}

	return op
}

// WriteDefinition emits a standalone definition page under definitions/.
// Uses the same template as resources but with no operations section.
func (m *MarkdownWriter) WriteDefinition(d *api.Definition) error {
	filename := strings.ToLower(d.Name) + "-" + string(d.Version)
	if d.Group != "" && d.Group != "core" {
		filename += "-" + string(d.Group)
	}
	filename += ".md"
	path := filepath.Join(m.OutputDir, "definitions", filename)

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("markdown: definition %s: %w", d.Name, err)
	}
	defer f.Close()

	if err := resourceTemplate.Execute(f, m.buildDefinitionPage(d)); err != nil {
		return fmt.Errorf("markdown: definition %s body: %w", d.Name, err)
	}
	return nil
}

// WriteOperation emits a single orphaned operation as its own page
// under operations/. Uses the shared "operation" define from the
// resource template so the shape matches operations that render inline
// on resource pages.
func (m *MarkdownWriter) WriteOperation(o *api.Operation) error {
	filename := operationSlug(o.ID) + ".md"
	path := filepath.Join(m.OutputDir, "operations", filename)

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("markdown: operation %s: %w", o.ID, err)
	}
	defer f.Close()

	writeSectionFrontmatter(f, o.Type.Name, o.Description(), m.nextResourceWeight())
	if err := resourceTemplate.ExecuteTemplate(f, "operation", buildTemplateOperation(o)); err != nil {
		return fmt.Errorf("markdown: operation %s body: %w", o.ID, err)
	}
	return nil
}

// operationSlug sanitizes an operation ID into a filesystem-safe name:
// lowercase, non-alphanumerics collapsed to '-'.
func operationSlug(id string) string {
	return strings.Trim(anchorRegex.ReplaceAllString(strings.ToLower(id), "-"), "-")
}

func (m *MarkdownWriter) Finalize() error {
	if m.finalized {
		return nil
	}
	m.finalized = true

	path := filepath.Join(m.OutputDir, hugoIndex)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("markdown: finalize: %w", err)
	}
	defer f.Close()

	writeSectionFrontmatter(f,
		m.Config.SpecTitle,
		fmt.Sprintf("Kubernetes API reference, version %s.", m.Config.SpecVersion),
		0)

	fmt.Fprintf(f, "# %s\n\n", m.Config.SpecTitle)
	fmt.Fprintf(f, "_Version: %s_\n\n", m.Config.SpecVersion)

	// Categories are added to m.toc in the order GenerateFiles encounters
	// them, which in --auto-detect mode depends on map iteration and is
	// therefore non-deterministic. Sort here so top-level _index.md diffs
	// stay quiet across runs. Header pages (Overview, API Groups, etc.)
	// stay at their fixed positions via tocSortRank.
	sort.SliceStable(m.toc, func(i, j int) bool {
		ri, rj := tocSortRank(m.toc[i].title), tocSortRank(m.toc[j].title)
		if ri != rj {
			return ri < rj
		}
		return m.toc[i].title < m.toc[j].title
	})

	for _, item := range m.toc {
		fmt.Fprintf(f, "- [%s](./%s)\n", item.title, item.path)
	}
	return nil
}

// tocSortRank pins well-known header pages at fixed positions ahead of
// the variable-length list of resource categories. Unknown titles fall
// into the categories bucket and sort alphabetically among themselves.
func tocSortRank(title string) int {
	switch title {
	case titleOverview:
		return 0
	case titleAPIGroups:
		return 1
	case titleDefinitions:
		return 3
	case titleOperations:
		return 4
	case titleOldVersions:
		return 5
	default:
		return 2
	}
}

// writeSectionFrontmatter emits the minimal frontmatter block used by
// non-resource pages (category _index.md, top-level _index.md). Resource
// pages have richer frontmatter and go through resource.tmpl instead.
// Description is optional; empty omits the line.
func writeSectionFrontmatter(w io.Writer, title, description string, weight int) {
	fmt.Fprintln(w, "---")
	fmt.Fprintln(w, `content_type: "api_reference"`)
	if description != "" {
		fmt.Fprintf(w, "description: %q\n", description)
	}
	fmt.Fprintf(w, "title: %q\n", title)
	fmt.Fprintf(w, "weight: %d\n", weight)
	fmt.Fprintln(w, "auto_generated: true")
	fmt.Fprintln(w, "---")
	fmt.Fprintln(w)
}

func anchor(s string) string {
	return strings.Trim(anchorRegex.ReplaceAllString(s, "-"), "-")
}

// escape is the minimal markdown escape needed so OpenAPI description text
// containing `<foo>` renders as literal rather than HTML.
func escape(s string) string {
	return strings.ReplaceAll(s, "<", `\<`)
}

func kebabCase(s string) string {
	return strings.Trim(anchorRegex.ReplaceAllString(strings.ToLower(s), "-"), "-")
}

func groupVersionString(group string, version api.ApiVersion) string {
	if group == "" || group == "core" {
		return version.String()
	}
	return fmt.Sprintf("%s/%s", group, version.String())
}

func writePipeTable(w io.Writer, headers []string, rowFn func(row func(cells ...string))) {
	fmt.Fprintln(w, "| "+strings.Join(headers, " | ")+" |")
	sep := make([]string, len(headers))
	for i := range sep {
		sep[i] = "---"
	}
	fmt.Fprintln(w, "| "+strings.Join(sep, " | ")+" |")
	rowFn(func(cells ...string) {
		fmt.Fprintln(w, "| "+strings.Join(cells, " | ")+" |")
	})
}

func (m *MarkdownWriter) writeSection(filename, title string) error {
	content := readOptionalSection(filename)
	if content == "" {
		content = "# " + title + "\n"
	}
	dst := filepath.Join(m.OutputDir, filename)
	return os.WriteFile(dst, []byte(content), 0644)
}

// readOptionalSection returns config/sections/<name> if present, "" otherwise.
// Read errors are swallowed to match the HTML writer's best-effort behaviour.
func readOptionalSection(name string) string {
	src := filepath.Join(api.SectionsDir, name)
	data, err := os.ReadFile(src)
	if err != nil {
		return ""
	}
	return string(data)
}

func (m *MarkdownWriter) nextCategoryWeight() int {
	m.categoryWeight += 10
	return m.categoryWeight
}

func (m *MarkdownWriter) nextResourceWeight() int {
	m.resourceWeight += 10
	return m.resourceWeight
}

