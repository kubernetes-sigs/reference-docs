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
	"regexp"
	"strings"

	"github.com/go-openapi/spec"
)

type ApiGroup string

func (g ApiGroup) String() string {
	return string(g)
}

func (g ApiGroup) LessThan(other ApiGroup) bool {
	// "apps" group APIs are newer than "extensions" group APIs
	if g.String() == "apps" && other.String() == "extensions" {
		return true
	}
	if other.String() == "apps" && g.String() == "extensions" {
		return false
	}

	// "policy" group APIs are newer than "extensions" group APIs
	if g == "policy" && other.String() == "extensions" {
		return true
	}
	if other.String() == "policy" && g.String() == "extensions" {
		return false
	}

	// "networking" group APIs are newer than "extensions" group APIs
	if g.String() == "networking" && other.String() == "extensions" {
		return true
	}
	if other.String() == "networking" && g.String() == "extensions" {
		return false
	}

	return strings.Compare(g.String(), other.String()) < 0
}

type ApiGroups []ApiGroup

func (a ApiGroups) Len() int      { return len(a) }
func (a ApiGroups) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ApiGroups) Less(i, j int) bool {
	return a[i].LessThan(a[j])
}

type ApiKind string

func (k ApiKind) String() string {
	return string(k)
}

type ApiVersion string

func (a ApiVersion) String() string {
	return string(a)
}

func (this ApiVersion) LessThan(that ApiVersion) bool {
	re := regexp.MustCompile("(v\\d+)(alpha|beta|)(\\d*)")
	thisMatches := re.FindStringSubmatch(string(this))
	thatMatches := re.FindStringSubmatch(string(that))

	a := []string{thisMatches[1]}
	if len(thisMatches) > 2 {
		a = []string{thisMatches[2], thisMatches[1], thisMatches[0]}
	}

	b := []string{thatMatches[1]}
	if len(thatMatches) > 2 {
		b = []string{thatMatches[2], thatMatches[1], thatMatches[0]}
	}

	for i := 0; i < len(a) && i < len(b); i++ {
		v1 := ""
		v2 := ""
		if i < len(a) {
			v1 = a[i]
		}
		if i < len(b) {
			v2 = b[i]
		}
		// If the "beta" or "alpha" is missing, then it is ga (empty string comes before non-empty string)
		if len(v1) == 0 || len(v2) == 0 {
			return v1 < v2
		}
		// The string with the higher number comes first (or in the case of alpha/beta, beta comes first)
		if v1 != v2 {
			return v1 > v2
		}
	}

	return false
}

type ApiVersions []ApiVersion

func (a ApiVersions) Len() int      { return len(a) }
func (a ApiVersions) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ApiVersions) Less(i, j int) bool {
	return a[i].LessThan(a[j])
}

type SortDefinitionsByName []*Definition

func (a SortDefinitionsByName) Len() int      { return len(a) }
func (a SortDefinitionsByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortDefinitionsByName) Less(i, j int) bool {
	if a[i].Name == a[j].Name {
		if a[i].Version.String() == a[j].Version.String() {
			return a[i].Group.String() < a[j].Group.String()
		}
		return a[i].Version.LessThan(a[j].Version)
	}
	return a[i].Name < a[j].Name
}

type SortDefinitionsByVersion []*Definition

func (a SortDefinitionsByVersion) Len() int      { return len(a) }
func (a SortDefinitionsByVersion) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortDefinitionsByVersion) Less(i, j int) bool {
	if a[i].Group.String() == a[j].Group.String() {
		return a[i].Version.LessThan(a[j].Version)
	}

	return a[i].Group.LessThan(a[j].Group)
}

type Definition struct {
	// open-api schema for the definition
	schema spec.Schema
	// Display name of the definition (e.g. Deployment)
	Name      string
	Group     ApiGroup
	ShowGroup bool

	// Api version of the definition (e.g. v1beta1)
	Version                 ApiVersion
	Kind                    ApiKind
	DescriptionWithEntities string
	GroupFullName           string

	// InToc is true if this definition should appear in the table of contents
	InToc        bool
	IsInlined    bool
	IsOldVersion bool

	FoundInField     bool
	FoundInOperation bool

	// Inline is a list of definitions that should appear inlined with this one in the documentations
	Inline SortDefinitionsByName

	// AppearsIn is a list of definition that this one appears in - e.g. PodSpec in Pod
	AppearsIn SortDefinitionsByName

	OperationCategories []*OperationCategory

	// Fields is a list of fields in this definition
	Fields Fields

	OtherVersions SortDefinitionsByName
	NewerVersions SortDefinitionsByName

	Sample SampleConfig

	FullName string
	Resource string
}

type GroupVersions map[string]ApiVersions

// Definitions indexes open-api definitions
type Definitions struct {
	All    map[string]*Definition
	ByKind map[string]SortDefinitionsByVersion

	// Available API groups and their versions
	GroupVersions GroupVersions
}

type DefinitionList []*Definition

type Config struct {
	ApiGroups           []ApiGroup          `yaml:"api_groups,omitempty"`
	ExampleLocation     string              `yaml:"example_location,omitempty"`
	OperationCategories []OperationCategory `yaml:"operation_categories,omitempty"`
	ResourceCategories  []ResourceCategory  `yaml:"resource_categories,omitempty"`
	ExcludedOperations  []string            `yaml:"excluded_operations,omitempty"`

	// Used to map the group as the resource sees it to the group as the operation sees it
	GroupMap map[string]string

	GroupFullNames map[string]string `yaml:"group_full_names,omitempty"`

	Definitions Definitions
	Operations  Operations
	SpecTitle   string
	SpecVersion string
}

type Field struct {
	Name                    string
	Type                    string
	Description             string
	DescriptionWithEntities string

	Definition *Definition // Optional Definition for complex types

	PatchStrategy string
	PatchMergeKey string
}

type Fields []*Field

func (a Fields) Len() int           { return len(a) }
func (a Fields) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Fields) Less(i, j int) bool { return a[i].Name < a[j].Name }

func (f Field) Link() string {
	if f.Definition != nil {
		return strings.Replace(f.Type, f.Definition.Name, f.Definition.MdLink(), -1)
	} else {
		return f.Type
	}
}

func (f Field) FullLink() string {
	if f.Definition != nil {
		return strings.Replace(f.Type, f.Definition.Name, f.Definition.HrefLink(), -1)
	} else {
		return f.Type
	}
}

// Operation defines a highlevel operation type such as Read, Replace, Patch
type OperationType struct {
	// Name is the display name of this operation
	Name string `yaml:",omitempty"`
	// Match is the regular expression of operation IDs that match this group where '${resource}' matches the resource name.
	Match string `yaml:",omitempty"`
}

type ExampleText struct {
	Tab  string
	Type string
	Text string
	Msg  string
}

type HttpResponse struct {
	Field
	Code string
}

type HttpResponses []*HttpResponse

type Operation struct {
	item          spec.PathItem
	op            *spec.Operation
	ID            string
	Type          OperationType
	Path          string
	HttpMethod    string
	Definition    *Definition
	BodyParams    Fields
	QueryParams   Fields
	PathParams    Fields
	HttpResponses HttpResponses

	ExampleConfig ExampleConfig
}

type Operations map[string]*Operation

// OperationCategory defines a group of related operations
type OperationCategory struct {
	// Name is the display name of this group
	Name string `yaml:",omitempty"`
	// Operations are the collection of Operations in this group
	OperationTypes []OperationType `yaml:"operation_types,omitempty"`
	// Default is true if this is the default operation group for operations that do not match any other groups
	Default bool `yaml:",omitempty"`

	Operations []*Operation
}

type ExampleProvider interface {
	GetTab() string
	GetRequestMessage() string
	GetResponseMessage() string
	GetRequestType() string
	GetResponseType() string
	GetSampleType() string
	GetSample(d *Definition) string
	GetRequest(o *Operation) string
	GetResponse(o *Operation) string
}

type EmptyExample struct{}
type CurlExample struct{}
type KubectlExample struct{}

type Resource struct {
	// Name is the display name of this Resource
	Name    string `yaml:",omitempty"`
	Version string `yaml:",omitempty"`
	Group   string `yaml:",omitempty"`

	// DescriptionWarning is a warning message to show along side this resource when displaying it
	DescriptionWarning string `yaml:"description_warning,omitempty"`
	// DescriptionNote is a note message to show along side this resource when displaying it
	DescriptionNote string `yaml:"description_note,omitempty"`
	// ConceptGuide is a link to the concept guide for this resource if it exists
	ConceptGuide string `yaml:"concept_guide,omitempty"`
	// RelatedTasks is as list of tasks related to this concept
	RelatedTasks []string `yaml:"related_tasks,omitempty"`

	// Definition of the object
	Definition *Definition
}

type Resources []*Resource

// ResourceCategory defines a category of Concepts
type ResourceCategory struct {
	// Name is the display name of this group
	Name string `yaml:",omitempty"`
	// Include is the name of the _<resource_category>.html file to include in the index.html
	Include string `yaml:",omitempty"`
	// Resources are the collection of Resources in this group
	Resources Resources `yaml:",omitempty"`
}

type ExampleConfig struct {
	Name         string `yaml:",omitempty"`
	Namespace    string `yaml:",omitempty"`
	Request      string `yaml:",omitempty"`
	Response     string `yaml:",omitempty"`
	RequestNote  string `yaml:",omitempty"`
	ResponseNote string `yaml:",omitempty"`
}

type SampleConfig struct {
	Note   string `yaml:",omitempty"`
	Sample string `yaml:",omitempty"`
}

type ResourceVisitor func(resource *Resource, d *Definition)
