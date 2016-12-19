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
	"sort"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/loads"
)


// Definitions indexes open-api definitions
type Definitions struct {
	ByVersionKind map[string]*Definition
	ByKind map[string]SortDefinitionsByVersion
}

func (d *Definitions) GetAllDefinitions() map[string]*Definition {
	return d.ByVersionKind
}

func (d *Definition) GroupDisplayName() string {
	if len(d.Group) <= 0 {
		return "Core"
	}
	return string(d.Group)
}

func (d *Definitions) GetOtherVersions(this *Definition) []*Definition {
	defs := d.ByKind[this.Name]
	others := []*Definition{}
	for _, def := range defs {
		if def.Version != this.Version {
			others = append(others, def)
		}
	}
	return others
}

// GetByVersionKind looks up a definition using its primary key (version,kind)
func (d *Definitions) GetByVersionKind(version, kind string) (*Definition, bool) {
	key := &Definition{Version: ApiVersion(version), Kind: ApiKind(kind)}
	r, f := d.ByVersionKind[key.Key()]
	return r, f
}

// GetByKey looks up a definition from its key (version.kind)
func (d *Definitions) GetByKey(key string) (*Definition, bool) {
	r, f := d.ByVersionKind[key]
	return r, f
}

// IsComplex returns true if the schema is for a complex (non-primitive) defintions
func (d *Definitions) IsComplex(s spec.Schema) bool {
	_, k := GetDefinitionVersionKind(s)
	return len(k) > 0
}


func (d *Definitions) GetForSchema(s spec.Schema) (*Definition, bool) {
	v, k := GetDefinitionVersionKind(s)
	if len(k) <= 0 {
		return nil, false
	}
	return d.GetByVersionKind(v, k)
}

func (d *Definitions) Put(defintion *Definition) {
	d.ByVersionKind[defintion.Key()] = defintion
}

// Initializes the fields for all definitions
func (d *Definitions) InitializeFieldsForAll() {
	for _, definition := range d.GetAllDefinitions() {
		d.InitializeFields(definition)
	}
}

// Initializes the fields for a definition
func (d *Definitions) InitializeFields(definition *Definition) {
	for fieldName, property := range definition.schema.Properties {
		def := strings.Replace(property.Description, "\n", " ", -1)
		field := &Field{
			Name:        fieldName,
			Type:        GetTypeName(property),
			Description: def,
		}
		if fieldDefinition, found := d.GetForSchema(property); found {
			field.Definition = fieldDefinition
		}
		definition.Fields = append(definition.Fields, field)
	}
}

func (d *Definitions) InitializeOtherVersions() {
	for _, definition := range d.GetAllDefinitions() {
		definition.OtherVersions = d.GetOtherVersions(definition)
	}
}

type Definition struct {
	// open-api schema for the definition
	schema              spec.Schema
	// Display name of the definition (e.g. Deployment)
	Name                string
	Group               ApiGroup
	// Api version of the definition (e.g. v1beta1)
	Version             ApiVersion
	Kind                ApiKind

	// InToc is true if this definition should appear in the table of contents
	InToc               bool
	IsInlined           bool
	IsOldVersion        bool

	FoundInField        bool
	FoundInOperation    bool

	// Inline is a list of definitions that should appear inlined with this one in the documentations
	Inline              SortDefinitionsByName

	// AppearsIn is a list of definition that this one appears in - e.g. PodSpec in Pod
	AppearsIn           SortDefinitionsByName

	OperationCategories []*OperationCategory

	// Fields is a list of fields in this definition
	Fields              Fields

	OtherVersions       SortDefinitionsByName
	NewerVersions       SortDefinitionsByName

	Sample              SampleConfig
}

func (d *Definition) Key() string {
	return fmt.Sprintf("%s.%s", d.Version, d.Kind)
}

func (d *Definition) MdLink() string {
	return fmt.Sprintf("[%s](#%s-%s)", d.Name, strings.ToLower(d.Name), d.Version)
}

func (d *Definition) HrefLink() string {
	return fmt.Sprintf("<a href=\"#%s-%s\">%s</a>", strings.ToLower(d.Name), d.Version, d.Name)
}

func (d *Definition) VersionLink() string {
	return fmt.Sprintf("<a href=\"#%s-%s\">%s</a>", strings.ToLower(d.Name), d.Version, d.Version)
}

func (d Definition) Description() string {
	return d.schema.Description
}

func VisitDefinitions(specs []*loads.Document, fn func(definition *Definition)) {
	for _, spec := range specs {
		for name, spec := range spec.Spec().Definitions {
			parts := strings.Split(name, ".")
			if len(parts) < 2 {
				fmt.Printf("Error: Could not find version and type for definition %s.\n", name)
				continue
			}
			fn(&Definition{
				schema:  spec,
				Name:    parts[1],
				Version: ApiVersion(parts[0]),
				Kind:    ApiKind(parts[1]),
			})
		}
	}
}

func (d *Definition) GetSamples() []ExampleText {
	r := []ExampleText{}
	for _, p := range GetExampleProviders() {
		r = append(r, ExampleText{
			Tab: p.GetTab(),
			Type: p.GetSampleType(),
			Text: p.GetSample(d),
		})
	}
	return r
}

func GetDefinitions(specs []*loads.Document) Definitions {
	d := Definitions{
		ByVersionKind: map[string]*Definition{},
		ByKind: map[string]SortDefinitionsByVersion{},
	}
	VisitDefinitions(specs, func(definition *Definition) {
		d.Put(definition)
	})
	d.InitializeFieldsForAll()
	for _, def := range d.GetAllDefinitions() {
		d.ByKind[def.Name] = append(d.ByKind[def.Name], def)
	}

	// If there are multiple versions for an object.  Mark all by the newest as old
	// Sort the ByKind index in by version with newer versions coming before older versions.
	for _, l := range d.ByKind {
		if len(l) <= 1 {
			continue
		}
		sort.Sort(l)
		// Mark all version as old
		for i, d := range l {
			if i > 0 {
				d.IsOldVersion = true
			}
		}
	}
	d.InitializeOtherVersions()
	d.initAppearsIn()
	d.initInlinedDefinitions()
	return d
}
