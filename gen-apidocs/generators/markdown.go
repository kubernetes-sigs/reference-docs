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

	// HugoMode emits Hugo-flavored markdown (markdown tables with class
	// attributes that Hugo render hooks can target). Set by NewHugoMDWriter.
	HugoMode bool

	// linkMap maps kind name → page path for cross-reference resolution.
	linkMap map[string]linkInfo

	classifications map[string]defClassification
	inlinedByParent map[string][]*api.Definition

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
	Version  api.ApiVersion // used to pick the highest version on name collisions
}

// resourcePage is the view model resource.tmpl consumes.
type resourcePage struct {
	APIVersion  string
	Kind        string
	Import      string
	Title       string
	Weight      int
	Anchor      string
	Description string
	Sections    []fieldSection
	Operations  []templateOperation
}

type fieldSection struct {
	Title       string
	Anchor      string
	Description string
	Fields      []templateField
}

type templateField struct {
	Name          string
	Type          string
	TypeHref      string // relative path#anchor, empty for primitives and unknowns
	Description   string
	Required      bool
	ConstValue    string // non-empty for fields with a fixed value (apiVersion, kind)
	PatchStrategy string // x-kubernetes-patch-strategy
	PatchMergeKey string // x-kubernetes-patch-merge-key
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
	HugoMode    bool
}

type templateParam struct {
	Name        string
	Type        string
	TypeHref    string // relative path#anchor when Type is a known definition
	Description string
}

type templateResponse struct {
	Code        string
	Type        string
	TypeHref    string // relative path#anchor when Type is a known definition
	Description string
}

const hugoIndex = "_index.md"

const apimachineryPrefix = "io.k8s.apimachinery."

const (
	titleOverview    = "Overview"
	titleAPIGroups   = "API Groups"
	titleDefinitions = "Definitions"
	titleOperations  = "Operations"
	titleOldVersions = "Old API Versions"
)

// utilityStandalone pins core/v1 utility kinds the BFS can't distinguish
// from real sub-types.
var utilityStandalone = map[string]bool{
	"LocalObjectReference":      true,
	"NodeSelector":              true,
	"NodeSelectorTerm":          true,
	"TypedLocalObjectReference": true,
}

var _ DocWriter = (*MarkdownWriter)(nil)

var anchorRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

var (
	enumHeaderRegex = regexp.MustCompile(`\s+Possible enum values:`)
	enumBulletRegex = regexp.MustCompile(`\s+- ` + "`")
)

//go:embed templates/resource.tmpl
var resourceTemplateSrc string

// q quotes for YAML frontmatter; md escapes `<` for body text; hugoRef
// wraps a relative path in a {{< ref >}} shortcode.
var resourceTemplate = template.Must(template.New("resource").Funcs(template.FuncMap{
	"q":       strconv.Quote,
	"md":      escape,
	"hugoRef": hugoRef,
}).Parse(resourceTemplateSrc))

// hugoRef wraps a path in a {{< ref >}} shortcode resolved by Hugo at build time.
func hugoRef(path string) string {
	return `{{< ref "` + path + `" >}}`
}

func NewMarkdownWriter(config *api.Config, copyright, title string) DocWriter {
	outputDir := api.BuildDir
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "MarkdownWriter: failed to create output dir %s: %v\n", outputDir, err)
	}
	m := &MarkdownWriter{
		Config:    config,
		OutputDir: outputDir,
		linkMap:   make(map[string]linkInfo),
	}
	// Order matters: linkMap consults classifications to alias inlined types.
	m.classifications = m.classifyDefinitions()
	m.buildLinkMap(config)
	return m
}

func NewHugoMDWriter(config *api.Config, copyright, title string) DocWriter {
	w := NewMarkdownWriter(config, copyright, title).(*MarkdownWriter)
	w.HugoMode = true
	return w
}

func (m *MarkdownWriter) Extension() string {
	return ".md"
}

func (m *MarkdownWriter) DefaultStaticContent(title string) string {
	return "# " + title + "\n"
}

// Pipeline methods below follow the call order in writer.go's GenerateFiles().

// No-op: kubernetes-api/_index.md in k/website is hand-curated.
func (m *MarkdownWriter) WriteOverview() error {
	return nil
}

func (m *MarkdownWriter) WriteAPIGroupVersions(gvs api.GroupVersions) error {
	path := filepath.Join(m.OutputDir, "group-versions.md")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("markdown: group versions: %w", err)
	}
	defer f.Close()

	weight := m.nextCategoryWeight()
	writeSectionFrontmatter(f, titleAPIGroups, "Kubernetes API groups and their served versions.", weight)
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
		path:   "group-versions.md",
		weight: weight,
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

func (m *MarkdownWriter) WriteResource(r *api.Resource) error {
	slug := m.currentCategory.slug
	if r.Definition != nil && r.Definition.IsOldVersion {
		return nil // markdown backend omits old-version pages; current version is canonical
	}

	filename := fmt.Sprintf("%s-%s.md", kebabName(r.Name), r.Definition.Version)
	path := filepath.Join(m.OutputDir, slug, filename)

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("markdown: resource %s: %w", r.Name, err)
	}
	defer f.Close()

	if err := resourceTemplate.Execute(f, m.buildResourcePage(r, slug)); err != nil {
		return fmt.Errorf("markdown: resource %s body: %w", r.Name, err)
	}

	return nil
}

// definitions/_index.md is required for Hugo to nest definition pages
// under the section; without it children flatten up to the parent level.
func (m *MarkdownWriter) WriteDefinitionsOverview() error {
	dir := filepath.Join(m.OutputDir, "definitions")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("markdown: definitions dir: %w", err)
	}
	indexPath := filepath.Join(dir, hugoIndex)
	f, err := os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("markdown: definitions _index.md: %w", err)
	}
	defer f.Close()
	weight := m.nextCategoryWeight()
	writeSectionFrontmatter(f, titleDefinitions, "", weight, hideFromNav())

	m.toc = append(m.toc, &mdTOCItem{
		title:  titleDefinitions,
		path:   filepath.Join("definitions", hugoIndex),
		weight: weight,
	})
	return nil
}

func (m *MarkdownWriter) WriteDefinition(d *api.Definition) error {
	if m.classifications[d.Key()].Mode == classifyInline {
		return nil
	}
	filename := kebabName(d.Name) + "-" + string(d.Version)
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

	if err := resourceTemplate.Execute(f, m.buildDefinitionPage(d, "definitions")); err != nil {
		return fmt.Errorf("markdown: definition %s body: %w", d.Name, err)
	}
	return nil
}

func (m *MarkdownWriter) WriteOrphanedOperationsOverview() error {
	dir := filepath.Join(m.OutputDir, "operations")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("markdown: operations dir: %w", err)
	}
	indexPath := filepath.Join(dir, hugoIndex)
	f, err := os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("markdown: operations _index.md: %w", err)
	}
	defer f.Close()
	weight := m.nextCategoryWeight()
	writeSectionFrontmatter(f, titleOperations, "", weight)

	m.toc = append(m.toc, &mdTOCItem{
		title:  titleOperations,
		path:   filepath.Join("operations", hugoIndex),
		weight: weight,
	})
	return nil
}

func (m *MarkdownWriter) WriteOperation(o *api.Operation) error {
	filename := operationSlug(o.ID) + ".md"
	path := filepath.Join(m.OutputDir, "operations", filename)

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("markdown: operation %s: %w", o.ID, err)
	}
	defer f.Close()

	writeSectionFrontmatter(f, o.Type.Name, o.Description(), m.nextResourceWeight())
	if err := resourceTemplate.ExecuteTemplate(f, "operation", m.buildTemplateOperation(o, "operations")); err != nil {
		return fmt.Errorf("markdown: operation %s body: %w", o.ID, err)
	}
	return nil
}

// No-op: old versions render as resource pages routed to their group folder.
func (m *MarkdownWriter) WriteOldVersionsOverview() error {
	return nil
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

	// Sort because --auto-detect populates categories via map iteration,
	// which is non-deterministic. tocSortRank pins header pages.
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

func (m *MarkdownWriter) buildResourcePage(r *api.Resource, currentCategory string) resourcePage {
	page := m.buildDefinitionPage(r.Definition, currentCategory)
	for _, oc := range r.Definition.OperationCategories {
		for _, o := range oc.Operations {
			page.Operations = append(page.Operations, m.buildTemplateOperation(o, currentCategory))
		}
	}
	return page
}

// One H2 section per top-level inlined child (Spec, Status, List); deeper
// inlines stay flattened with dot-notation in their containing section.
func (m *MarkdownWriter) buildDefinitionPage(d *api.Definition, currentCategory string) resourcePage {
	page := resourcePage{
		APIVersion:  groupVersionString(d.GroupFullName, d.Version),
		Kind:        d.Name,
		Import:      d.GoImportPath(),
		Title:       d.Name,
		Weight:      m.nextResourceWeight(),
		Anchor:      anchor(d.Name),
		Description: d.DescriptionWithEntities,
	}

	// Inline closure rooted at d gates which types may flatten inline;
	// anything outside renders as a cross-page link.
	allowInline := map[string]bool{}
	var collect func(x *api.Definition)
	collect = func(x *api.Definition) {
		for _, c := range x.Inline {
			if allowInline[c.Key()] {
				continue
			}
			allowInline[c.Key()] = true
			collect(c)
		}
	}
	collect(d)

	// Section-worthy = directly referenced field type, or d's List type.
	sectionTypes := map[string]bool{}
	sectionOrder := []*api.Definition{}
	addSection := func(def *api.Definition) {
		if def == nil || sectionTypes[def.Key()] {
			return
		}
		sectionTypes[def.Key()] = true
		sectionOrder = append(sectionOrder, def)
	}
	for _, fld := range d.Fields {
		if fld.Definition != nil && allowInline[fld.Definition.Key()] {
			addSection(fld.Definition)
		}
	}
	for _, c := range d.Inline {
		if c.Name == d.Name+"List" {
			addSection(c)
		}
	}

	if siblings := m.inlinedByParent[d.Key()]; len(siblings) > 0 {
		sorted := append([]*api.Definition(nil), siblings...)
		sort.Slice(sorted, func(i, j int) bool { return sorted[i].Name < sorted[j].Name })
		for _, c := range sorted {
			addSection(c)
		}
	}

	visited := map[string]bool{d.Key(): true}
	root := fieldSection{
		Title:       d.Name,
		Anchor:      anchor(d.Name),
		Description: d.DescriptionWithEntities,
	}
	m.appendFields(&root, d, "", currentCategory, allowInline, sectionTypes, visited)
	page.Sections = append(page.Sections, root)

	for _, s := range sectionOrder {
		section := fieldSection{
			Title:       s.Name,
			Anchor:      anchor(s.Name),
			Description: s.DescriptionWithEntities,
		}
		svisited := map[string]bool{d.Key(): true, s.Key(): true}
		m.appendFields(&section, s, "", currentCategory, allowInline, sectionTypes, svisited)
		page.Sections = append(page.Sections, section)
	}
	return page
}

// visited guards against pathological cycles.
func (m *MarkdownWriter) appendFields(section *fieldSection, d *api.Definition, prefix, currentCategory string, allowInline, sectionTypes, visited map[string]bool) {
	required := map[string]bool{}
	for _, name := range d.RequiredFields() {
		required[name] = true
	}

	for _, fld := range d.Fields {
		fullName := fld.Name
		if prefix != "" {
			fullName = prefix + "." + fld.Name
		}

		typeHref := m.resolveType(fld.Type, currentCategory)
		if fld.Definition != nil && sectionTypes[fld.Definition.Key()] {
			typeHref = "#" + anchor(fld.Definition.Name)
		}

		section.Fields = append(section.Fields, templateField{
			Name:          fullName,
			Type:          fld.Type,
			TypeHref:      typeHref,
			Description:   fld.Description,
			Required:      required[fld.Name],
			ConstValue:    constValueFor(fullName, "", ""),
			PatchStrategy: fld.PatchStrategy,
			PatchMergeKey: fld.PatchMergeKey,
		})

		if fld.Definition == nil {
			continue
		}
		key := fld.Definition.Key()
		if !allowInline[key] || sectionTypes[key] || visited[key] {
			continue
		}
		visited[key] = true
		m.appendFields(section, fld.Definition, fullName, currentCategory, allowInline, sectionTypes, visited)
	}
}

func (m *MarkdownWriter) buildTemplateOperation(o *api.Operation, currentCategory string) templateOperation {
	op := templateOperation{
		Verb:     strings.ToLower(o.HttpMethod),
		Title:    o.Type.Name,
		Method:   o.HttpMethod,
		Path:     o.Path,
		HugoMode: m.HugoMode,
	}

	convert := func(params api.Fields) []templateParam {
		out := make([]templateParam, 0, len(params))
		for _, p := range params {
			out = append(out, templateParam{
				Name:        p.Name,
				Type:        p.Type,
				TypeHref:    m.resolveType(p.Type, currentCategory),
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
			TypeHref:    m.resolveType(rsp.Field.Type, currentCategory),
			Description: rsp.Field.Description,
		})
	}

	return op
}

// On collisions: resource entries beat definition entries; higher version wins.
func (m *MarkdownWriter) buildLinkMap(config *api.Config) {
	m.linkResources(config.ResourceCategories)
	m.linkDefinitions(config.Definitions.All)
}

type classifyMode int

const (
	classifySkip classifyMode = iota
	classifyInline
	classifyStandalone
)

type defClassification struct {
	Mode       classifyMode
	InlineInto *api.Definition // non-nil iff Mode == classifyInline
}

// classifyDefinitions returns each definition's emit mode.
func (m *MarkdownWriter) classifyDefinitions() map[string]defClassification {
	out := make(map[string]defClassification, len(m.Config.Definitions.All))
	m.inlinedByParent = map[string][]*api.Definition{}
	for _, d := range m.Config.Definitions.All {
		if d.IsOldVersion || d.IsInlined || d.InToc {
			out[d.Key()] = defClassification{Mode: classifySkip}
			continue
		}
		if strings.HasPrefix(d.SwaggerKey, apimachineryPrefix) || utilityStandalone[d.Name] {
			out[d.Key()] = defClassification{Mode: classifyStandalone}
			continue
		}
		if home := m.closestTopLevelHome(d); home != nil {
			out[d.Key()] = defClassification{Mode: classifyInline, InlineInto: home}
			m.inlinedByParent[home.Key()] = append(m.inlinedByParent[home.Key()], d)
			continue
		}
		out[d.Key()] = defClassification{Mode: classifyStandalone}
	}
	return out
}

// closestTopLevelHome returns the unique closest InToc ancestor of d via
// AppearsIn, or nil on a tie or no reachable top-level.
func (m *MarkdownWriter) closestTopLevelHome(d *api.Definition) *api.Definition {
	type qItem struct {
		def  *api.Definition
		dist int
	}
	seen := map[string]bool{d.Key(): true}
	queue := []qItem{{def: d, dist: 0}}

	minDist := -1
	winners := []*api.Definition{}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if minDist >= 0 && cur.dist > minDist {
			break
		}

		for _, parent := range cur.def.AppearsIn {
			if parent == nil || seen[parent.Key()] {
				continue
			}
			seen[parent.Key()] = true
			nextDist := cur.dist + 1
			if parent.InToc {
				if minDist < 0 {
					minDist = nextDist
				}
				if nextDist == minDist {
					winners = append(winners, parent)
				}
				continue
			}
			queue = append(queue, qItem{def: parent, dist: nextDist})
		}
	}

	if len(winners) == 1 {
		return winners[0]
	}
	return nil
}

func (m *MarkdownWriter) linkResources(categories []api.ResourceCategory) {
	for _, c := range categories {
		slug := kebabCase(c.Name)
		for _, r := range c.Resources {
			if r.Definition == nil {
				continue
			}
			filename := kebabName(r.Name) + "-" + string(r.Definition.Version)
			m.recordLink(r.Definition.Name, slug, filename, r.Definition.Version)
			// Types inlined into this page (pattern-matched and refcount-derived)
			// share the parent's URL; the anchor on the parent is what selects them.
			for _, child := range r.Definition.Inline {
				m.recordLink(child.Name, slug, filename, child.Version)
			}
			for _, child := range m.inlinedByParent[r.Definition.Key()] {
				m.recordLink(child.Name, slug, filename, child.Version)
			}
		}
	}
}

func (m *MarkdownWriter) linkDefinitions(all map[string]*api.Definition) {
	for _, d := range all {
		if d.InToc || d.IsInlined || d.IsOldVersion {
			continue
		}
		if m.classifications[d.Key()].Mode == classifyInline {
			continue
		}
		filename := kebabName(d.Name) + "-" + string(d.Version)
		if d.Group != "" && d.Group != "core" {
			filename += "-" + string(d.Group)
		}
		m.recordLink(d.Name, "definitions", filename, d.Version)
	}
}

func (m *MarkdownWriter) recordLink(name, category, filename string, version api.ApiVersion) {
	if existing, ok := m.linkMap[name]; ok {
		if existing.Category != "definitions" && category == "definitions" {
			return
		}
		// LessThan returns true when the receiver is the higher version.
		if existing.Category == category && existing.Version.LessThan(version) {
			return
		}
	}
	m.linkMap[name] = linkInfo{
		Category: category,
		Filename: filename,
		Anchor:   anchor(name),
		Version:  version,
	}
}

// resolveType returns a relative path#anchor for typeName, or "" if unknown.
// currentCategory is the slug containing the calling page; "MutatingWebhook
// array" is normalised to "MutatingWebhook" before lookup.
func (m *MarkdownWriter) resolveType(typeName, currentCategory string) string {
	info, ok := m.linkMap[typeName]
	if !ok {
		if bare, stripped := strings.CutSuffix(typeName, " array"); stripped {
			info, ok = m.linkMap[bare]
		}
	}
	if !ok {
		return ""
	}
	var path string
	if info.Category == currentCategory {
		path = info.Filename
	} else {
		path = "../" + info.Category + "/" + info.Filename
	}
	return path + "#" + info.Anchor
}

// readOptionalSection swallows read errors on purpose; missing or
// unreadable section files fall back to generated content.
func readOptionalSection(name string) string {
	src := filepath.Join(api.SectionsDir, name)
	data, err := os.ReadFile(src)
	if err != nil {
		return ""
	}
	return string(data)
}

type sectionFrontmatterOpt func(w io.Writer)

// hideFromNav drops the section and its children from the Docsy sidebar.
// Cross-ref-only sections (definitions) would otherwise add hundreds of
// nav entries and slow every page render across k/website.
func hideFromNav() sectionFrontmatterOpt {
	return func(w io.Writer) {
		fmt.Fprintln(w, "_build:")
		fmt.Fprintln(w, "  list: never")
		fmt.Fprintln(w, "toc_hide: true")
	}
}

// writeSectionFrontmatter emits the minimal frontmatter for non-resource
// pages. Resource pages go through resource.tmpl instead.
func writeSectionFrontmatter(w io.Writer, title, description string, weight int, opts ...sectionFrontmatterOpt) {
	fmt.Fprintln(w, "---")
	fmt.Fprintln(w, `content_type: "api_reference"`)
	if description != "" {
		fmt.Fprintf(w, "description: %q\n", description)
	}
	fmt.Fprintf(w, "title: %q\n", title)
	fmt.Fprintf(w, "weight: %d\n", weight)
	fmt.Fprintln(w, "auto_generated: true")
	for _, o := range opts {
		o(w)
	}
	fmt.Fprintln(w, "---")
	fmt.Fprintln(w)
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

// tocSortRank pins header pages at fixed positions; categories share rank 2
// and sort alphabetically among themselves.
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

func anchor(s string) string {
	return strings.Trim(anchorRegex.ReplaceAllString(s, "-"), "-")
}

// escape covers the only markdown-breaking character in OpenAPI descriptions:
// raw `<` that would otherwise be read as HTML.
func escape(s string) string {
	s = strings.ReplaceAll(s, "<", `\<`)
	s = enumHeaderRegex.ReplaceAllString(s, "<br/><br/>Possible enum values:")
	s = enumBulletRegex.ReplaceAllString(s, "<br/> - `")
	return s
}

func kebabCase(s string) string {
	return strings.Trim(anchorRegex.ReplaceAllString(strings.ToLower(s), "-"), "-")
}

var (
	kebabBoundary1 = regexp.MustCompile(`([a-z0-9])([A-Z])`)
	kebabBoundary2 = regexp.MustCompile(`([A-Z])([A-Z][a-z])`)
)

func kebabName(s string) string {
	s = kebabBoundary2.ReplaceAllString(s, "$1-$2")
	s = kebabBoundary1.ReplaceAllString(s, "$1-$2")
	return strings.ToLower(s)
}

func groupVersionString(group string, version api.ApiVersion) string {
	if group == "" || group == "core" {
		return version.String()
	}
	return fmt.Sprintf("%s/%s", group, version.String())
}

func operationSlug(id string) string {
	return strings.Trim(anchorRegex.ReplaceAllString(strings.ToLower(id), "-"), "-")
}

// constValueFor hard-codes the two fields Kubernetes manifests always
// carry with fixed values (apiVersion and kind). Swagger doesn't tag
// them as const so we derive them from the GVK.
func constValueFor(fieldName, apiVersion, kind string) string {
	switch fieldName {
	case "apiVersion":
		return apiVersion
	case "kind":
		return kind
	}
	return ""
}

func (m *MarkdownWriter) nextCategoryWeight() int {
	m.categoryWeight += 10
	return m.categoryWeight
}

func (m *MarkdownWriter) nextResourceWeight() int {
	m.resourceWeight += 10
	return m.resourceWeight
}
