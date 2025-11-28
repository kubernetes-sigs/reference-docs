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

package api

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

// inlineDefinition is a definition that should be inlined when displaying a Concept
// instead of appearing the in "Definitions"
type inlineDefinition struct {
	Name  string
	Match string
}

var _INLINE_DEFINITIONS = []inlineDefinition{
	{Name: "Spec", Match: "${resource}Spec"},
	{Name: "Status", Match: "${resource}Status"},
	{Name: "List", Match: "${resource}List"},
	{Name: "Strategy", Match: "${resource}Strategy"},
	{Name: "Rollback", Match: "${resource}Rollback"},
	{Name: "RollingUpdate", Match: "RollingUpdate${resource}"},
	{Name: "EventSource", Match: "${resource}EventSource"},
}

func NewDefinitions(config *Config, specs []*loads.Document) (*Definitions, error) {
	s := &Definitions{
		All:           map[string]*Definition{},
		ByKind:        map[string]SortDefinitionsByVersion{},
		GroupVersions: map[string]ApiVersions{},
	}

	if err := LoadDefinitions(config, specs, s); err != nil {
		return nil, fmt.Errorf("failed to load definitions from OpenAPI specs: %w", err)
	}

	s.initialize()
	return s, nil
}

func (s *Definitions) initialize() {
	// initialize fields for all definitions
	for _, d := range s.All {
		s.InitializeFields(d)
	}

	for _, d := range s.All {
		s.ByKind[d.Name] = append(s.ByKind[d.Name], d)
	}

	// If there are multiple versions for an object.  Mark all by the newest as old
	// Sort the ByKind index in by version with newer versions coming before older versions.
	for k, l := range s.ByKind {
		if len(l) <= 1 {
			continue
		}
		sort.Sort(l)
		// Mark all version as old
		for i, d := range l {
			if len(l) > 1 {
				if i == 0 {
					fmt.Printf("Current Version: %s.%s.%s", d.Group, d.Version, k)
					if len(l) > i-1 {
						fmt.Printf(" Old Versions: [")
					}
				} else {
					fmt.Printf("%s.%s.%s", d.Group, d.Version, k)
					if len(l) > i-1 {
						fmt.Printf(",")
					}
					d.IsOldVersion = true
				}
			}
		}
		if len(l) > 1 {
			fmt.Printf("]\n")
		}
	}

	// Initialize OtherVersions
	for _, d := range s.All {
		defs := s.ByKind[d.Name]
		others := []*Definition{}
		for _, def := range defs {
			if def.Version != d.Version {
				others = append(others, def)
			}
		}
		d.OtherVersions = others
	}

	// Initialize AppearsIn and FoundInField
	for _, d := range s.All {
		for _, r := range s.getReferences(d) {
			r.AppearsIn = append(r.AppearsIn, d)
			r.FoundInField = true
		}
	}

	// Initialize Inline, IsInlined
	// Note: examples of inline definitions are "Spec", "Status", "List", etc
	for _, d := range s.All {
		for _, name := range s.getInlineDefinitionNames(d.Name) {
			if cr, ok := s.GetByVersionKind(string(d.Group), string(d.Version), name); ok {
				d.Inline = append(d.Inline, cr)
				cr.IsInlined = true
				cr.FoundInField = true
			}
		}
	}
}

func (s *Definitions) getInlineDefinitionNames(parent string) []string {
	names := []string{}
	for _, id := range _INLINE_DEFINITIONS {
		name := strings.ReplaceAll(id.Match, "${resource}", parent)
		names = append(names, name)
	}
	return names
}

func (s *Definitions) getReferences(d *Definition) []*Definition {
	refs := []*Definition{}
	// Find all of the definitions referenced by this definition
	for _, p := range d.schema.Properties {
		if p.AdditionalProperties != nil && p.AdditionalProperties.Schema != nil {
			additionalProperty := p.AdditionalProperties.Schema.Ref
			if len(additionalProperty.String()) > 0 {
				group, version, kind := GetDefinitionVersionKindFromString(additionalProperty.String())
				definition, ok := s.GetByVersionKind(group, version, kind)
				if ok {
					refs = append(refs, definition)
				}
			}
		}
		if !IsComplex(p) {
			// Skip primitive types and collections of primitive types
			continue
		}
		// Look up the definition for the referenced definitions
		if schema, ok := s.GetForSchema(p); ok {
			refs = append(refs, schema)
		} else {
			g, v, k := GetDefinitionVersionKind(p)
			fmt.Printf("Could not locate referenced property of %s: %s (%s/%s).\n", d.Name, g, k, v)
		}
	}

	// Additionally probe the entire schema (including nested additionalProperties, items, allOf/oneOf/anyOf)
	// for referenced definitions using a recursive walker. This covers deeply nested cases not found
	// by the top-level property inspection above.
	nested := s.findReferencedDefinitions(&d.schema, nil)
	// Merge nested results into refs, avoiding duplicates
	present := map[string]bool{}
	for _, r := range refs {
		if r != nil {
			present[r.Key()] = true
		}
	}
	for _, r := range nested {
		if r == nil {
			continue
		}
		if !present[r.Key()] {
			refs = append(refs, r)
			present[r.Key()] = true
		}
	}

	return refs
}

// findReferencedDefinitions walks a schema recursively and returns all Definitions referenced from it.
// It deduplicates results and prevents infinite loops using a visited set keyed by schema.Ref.String()
// or address for inline schemas.
func (s *Definitions) findReferencedDefinitions(schema *spec.Schema, visited map[string]bool) []*Definition {
	if schema == nil {
		return nil
	}
	if visited == nil {
		visited = map[string]bool{}
	}
	// produce a stable key for visited detection
	key := schema.Ref.String()
	if key == "" {
		key = fmt.Sprintf("addr:%p", schema)
	}
	if visited[key] {
		return nil
	}
	visited[key] = true

	results := []*Definition{}

	// If this schema itself maps to a Definition (usually via $ref), record it and stop descending.
	if def, ok := s.GetForSchema(*schema); ok {
		return []*Definition{def}
	}

	// Recurse into AdditionalProperties (map value type)
	if schema.AdditionalProperties != nil && schema.AdditionalProperties.Schema != nil {
		results = append(results, s.findReferencedDefinitions(schema.AdditionalProperties.Schema, visited)...)
	}

	// Recurse into items (arrays)
	if schema.Items != nil {
		if schema.Items.Schema != nil {
			results = append(results, s.findReferencedDefinitions(schema.Items.Schema, visited)...)
		}
		for i := range schema.Items.Schemas {
			results = append(results, s.findReferencedDefinitions(&schema.Items.Schemas[i], visited)...)
		}
	}

	// Recurse into properties
	for _, prop := range schema.Properties {
		// prop is a spec.Schema (value) so take its address
		results = append(results, s.findReferencedDefinitions(&prop, visited)...)
	}

	// Recurse into composition keywords
	for i := range schema.AllOf {
		results = append(results, s.findReferencedDefinitions(&schema.AllOf[i], visited)...)
	}
	for i := range schema.OneOf {
		results = append(results, s.findReferencedDefinitions(&schema.OneOf[i], visited)...)
	}
	for i := range schema.AnyOf {
		results = append(results, s.findReferencedDefinitions(&schema.AnyOf[i], visited)...)
	}

	// Deduplicate and filter out nil definitions
	seen := make(map[string]bool, len(results))
	out := make([]*Definition, 0, len(results))

	for _, d := range results {
		if d == nil {
			continue
		}
		key := d.Key()
		if _, exists := seen[key]; !exists {
			seen[key] = true
			out = append(out, d)
		}
	}

	return out
}

func (s *Definitions) parameterToField(param spec.Parameter) *Field {
	f := &Field{
		Name:        param.Name,
		Description: strings.ReplaceAll(param.Description, "\n", " "),
	}
	if param.Schema != nil {
		f.Type = GetTypeName(*param.Schema)
		if fieldType, ok := s.GetForSchema(*param.Schema); ok {
			f.Definition = fieldType
		}
	}
	return f
}

// GetByVersionKind looks up a definition using its primary key (version,kind)
func (s *Definitions) GetByVersionKind(group, version, kind string) (*Definition, bool) {
	key := &Definition{Group: ApiGroup(group), Version: ApiVersion(version), Kind: ApiKind(kind)}
	r, f := s.All[key.Key()]
	return r, f
}

func (s *Definitions) GetForSchema(schema spec.Schema) (*Definition, bool) {
	g, v, k := GetDefinitionVersionKind(schema)
	if len(k) <= 0 {
		return nil, false
	}
	return s.GetByVersionKind(g, v, k)
}

// Initializes the fields for a definition
func (s *Definitions) InitializeFields(d *Definition) {
	for fieldName, property := range d.schema.Properties {
		des := strings.ReplaceAll(property.Description, "\n", " ")
		f := &Field{
			Name:        fieldName,
			Type:        GetTypeName(property),
			Description: EscapeAsterisks(des),
		}
		if len(property.Extensions) > 0 {
			if ps, ok := property.Extensions.GetString(patchStrategyKey); ok {
				f.PatchStrategy = ps
			}
			if pmk, ok := property.Extensions.GetString(patchMergeKeyKey); ok {
				f.PatchMergeKey = pmk
			}
		}

		if fd, ok := s.GetForSchema(property); ok {
			f.Definition = fd
		}
		d.Fields = append(d.Fields, f)
	}
}

func (d *Definition) GroupDisplayName() string {
	if len(d.GroupFullName) > 0 {
		return d.GroupFullName
	}
	if len(d.Group) <= 0 || d.Group == "core" {
		return "Core"
	}
	return string(d.Group)
}

func (d *Definition) Key() string {
	return fmt.Sprintf("%s.%s.%s", d.Group, d.Version, d.Kind)
}

func (d *Definition) LinkID() string {
	groupName := strings.ReplaceAll(strings.ToLower(d.GroupFullName), ".", "-")
	link := fmt.Sprintf("%s-%s-%s", d.Name, d.Version, groupName)
	return strings.ToLower(link)
}

func (d *Definition) MdLink() string {
	groupName := strings.ReplaceAll(strings.ToLower(d.GroupFullName), ".", "-")
	return fmt.Sprintf("[%s](#%s-%s-%s)", d.Name, strings.ToLower(d.Name), d.Version, groupName)
}

func (d *Definition) HrefLink() string {
	groupName := strings.ReplaceAll(strings.ToLower(d.GroupFullName), ".", "-")
	return fmt.Sprintf("<a href=\"#%s-%s-%s\">%s</a>", strings.ToLower(d.Name), d.Version, groupName, d.Name)
}

func (d *Definition) FullHrefLink() string {
	groupName := strings.ReplaceAll(strings.ToLower(d.GroupFullName), ".", "-")
	return fmt.Sprintf("<a href=\"#%s-%s-%s\">%s [%s/%s]</a>", strings.ToLower(d.Name),
		d.Version, groupName, d.Name, d.Group, d.Version)
}

func (d *Definition) VersionLink() string {
	groupName := strings.ReplaceAll(strings.ToLower(d.GroupFullName), ".", "-")
	return fmt.Sprintf("<a href=\"#%s-%s-%s\">%s</a>", strings.ToLower(d.Name), d.Version, groupName, d.Version)
}

func (d *Definition) Description() string {
	return EscapeAsterisks(d.schema.Description)
}

func (d *Definition) GetResourceName() string {
	if len(d.Resource) > 0 {
		return d.Resource
	}
	resource := strings.ToLower(d.Name)
	if strings.HasSuffix(resource, "y") {
		return strings.TrimSuffix(resource, "y") + "ies"
	}
	return resource + "s"
}

func (d *Definition) initExample(config *Config) error {
	path := filepath.Join(ConfigDir, config.ExampleLocation, d.Name, d.Name+".yaml")
	file := strings.ReplaceAll(strings.ToLower(path), " ", "_")

	// missing files are okay
	if _, err := os.Stat(file); err != nil {
		return nil
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read example file %s for definition %s: %w", file, d.Name, err)
	}

	if err = yaml.Unmarshal(content, &d.Sample); err != nil {
		return fmt.Errorf("failed to unmarshal YAML content from file %s for definition %s: %w", file, d.Name, err)
	}

	return nil
}

func (d *Definition) GetSamples() []ExampleText {
	r := []ExampleText{}
	for _, p := range EmptyExampleProviders {
		r = append(r, ExampleText{
			Tab:  p.GetTab(),
			Type: p.GetSampleType(),
			Text: p.GetSample(d),
		})
	}
	return r
}

func (a DefinitionList) Len() int      { return len(a) }
func (a DefinitionList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a DefinitionList) Less(i, j int) bool {
	return strings.Compare(a[i].Name, a[j].Name) < 0
}
