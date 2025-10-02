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
	"strconv"
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

	// "events" group APIs are newer than "core" group APIs
	if g.String() == "events" && other.String() == "core" {
		return true
	}
	if other.String() == "events" && g.String() == "core" {
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
type versionType int

const (
	// Bigger the version type number, higher priority it is
	versionTypeAlpha versionType = iota
	versionTypeBeta
	versionTypeGA
)

var kubeVersionRegex = regexp.MustCompile(`^v([\d]+)(?:(alpha|beta)([\d]+))?$`)

func (a ApiVersion) String() string {
	return string(a)
}

// logic copied from k8s.apimachinery
func parseKubeVersion(v string) (majorVersion int, vType versionType, minorVersion int, ok bool) {
	var err error
	submatches := kubeVersionRegex.FindStringSubmatch(v)
	if len(submatches) != 4 {
		return 0, 0, 0, false
	}
	switch submatches[2] {
	case "alpha":
		vType = versionTypeAlpha
	case "beta":
		vType = versionTypeBeta
	case "":
		vType = versionTypeGA
	default:
		return 0, 0, 0, false
	}
	if majorVersion, err = strconv.Atoi(submatches[1]); err != nil {
		return 0, 0, 0, false
	}
	if vType != versionTypeGA {
		if minorVersion, err = strconv.Atoi(submatches[3]); err != nil {
			return 0, 0, 0, false
		}
	}
	return majorVersion, vType, minorVersion, true
}

// logic copied from k8s.apimachinery
func compareVersionStrings(v1, v2 string) int {
	if v1 == v2 {
		return 0
	}
	v1major, v1type, v1minor, ok1 := parseKubeVersion(v1)
	v2major, v2type, v2minor, ok2 := parseKubeVersion(v2)
	switch {
	case !ok1 && !ok2:
		return strings.Compare(v2, v1)
	case !ok1 && ok2:
		return -1
	case ok1 && !ok2:
		return 1
	}
	if v1type != v2type {
		return int(v1type) - int(v2type)
	}
	if v1major != v2major {
		return v1major - v2major
	}
	return v1minor - v2minor
}

func (ver ApiVersion) LessThan(that ApiVersion) bool {
	res := compareVersionStrings(string(ver), string(that))
	return res > 0
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
	OperationGroupMap map[string]string `yaml:"operation_group_map,omitempty"`

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
		return strings.ReplaceAll(f.Type, f.Definition.Name, f.Definition.MdLink())
	} else {
		return f.Type
	}
}

func (f Field) FullLink() string {
	if f.Definition != nil {
		return strings.ReplaceAll(f.Type, f.Definition.Name, f.Definition.HrefLink())
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
