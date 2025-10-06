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
	"flag"
	"fmt"
	"html"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"gopkg.in/yaml.v2"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

var AllowErrors = flag.Bool("allow-errors", false, "If true, don't fail on errors.")
var WorkDir = flag.String("work-dir", "", "Working directory for the generator.")
var UseTags = flag.Bool("use-tags", false, "If true, use the openapi tags instead of the config yaml.")
var KubernetesRelease = flag.String("kubernetes-release", "", "Kubernetes release version.")

// titleCase converts a string to title case as a replacement for deprecated strings.Title
func titleCase(s string) string {
	if s == "" {
		return s
	}
	words := strings.Fields(s)
	for i, word := range words {
		if word != "" {
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}
	return strings.Join(words, " ")
}

// Directory for output files
var BuildDir string

// Directory for configuration and data files
var ConfigDir string

// Directory for static sections
var SectionsDir string

// Directory for temporary files that will eventually get merged into the HTML output file.
var IncludesDir string

// Directory for versioned configuration file and swagger.json
var VersionedConfigDir string

func NewConfig() (*Config, error) {
	// Initialize global directories
	BuildDir = filepath.Join(*WorkDir, "build")
	ConfigDir = filepath.Join(*WorkDir, "config")
	IncludesDir = filepath.Join(BuildDir, "includes")
	SectionsDir = filepath.Join(ConfigDir, "sections")

	k8sRelease := fmt.Sprintf("v%s", strings.ReplaceAll(*KubernetesRelease, ".", "_"))
	VersionedConfigDir = filepath.Join(ConfigDir, k8sRelease)

	config, err := loadAndInitializeConfig()
	if err != nil {
		return nil, err
	}

	if err := processDefinitionsAndOperations(config); err != nil {
		return nil, err
	}

	return config, nil
}

// loadAndInitializeConfig loads configuration and specs, then initializes basic config
func loadAndInitializeConfig() (*Config, error) {
	config, err := LoadConfigFromYAML()
	if err != nil {
		return nil, fmt.Errorf("failed to load config yaml: %w", err)
	}

	specs, err := LoadOpenApiSpec()
	if err != nil {
		return nil, fmt.Errorf("failed to load openapi spec: %w", err)
	}

	// Parse spec version
	ParseSpecInfo(specs, config)

	// Set the spec version
	config.SpecVersion = fmt.Sprintf("v%s.%s", *KubernetesRelease, "0")

	// Initialize all of the operations
	defs, err := NewDefinitions(config, specs)
	if err != nil {
		return nil, fmt.Errorf("failed to init definitions: %w", err)
	}
	config.Definitions = *defs

	return config, nil
}

// processDefinitionsAndOperations handles the main processing logic
func processDefinitionsAndOperations(config *Config) error {
	specs, err := LoadOpenApiSpec()
	if err != nil {
		return fmt.Errorf("failed to load openapi spec: %w", err)
	}

	if *UseTags {
		// Initialize the config and ToC from the tags on definitions
		if err := config.genConfigFromTags(specs); err != nil {
			return fmt.Errorf("failed to generate config from tags: %w", err)
		}
	} else {
		// Initialization for ToC resources only
		if err := config.visitResourcesInToc(); err != nil {
			return fmt.Errorf("failed to visit resources in TOC: %w", err)
		}
	}

	if err := config.initOperations(specs); err != nil {
		return fmt.Errorf("failed to init operations: %w", err)
	}

	// replace unicode escape sequences with HTML entities.
	config.escapeDescriptions()

	config.CleanUp()

	// Prune anything that shouldn't be in the ToC
	if *UseTags {
		config.pruneResourceCategories()
	}

	return nil
}

// pruneResourceCategories removes resources that shouldn't be in the ToC
func (c *Config) pruneResourceCategories() {
	categories := []ResourceCategory{}
	for _, cat := range c.ResourceCategories {
		resources := Resources{}
		for _, r := range cat.Resources {
			if d, f := c.Definitions.GetByVersionKind(r.Group, r.Version, r.Name); f {
				if d.InToc {
					resources = append(resources, r)
				}
			}
		}
		cat.Resources = resources
		if len(resources) > 0 {
			categories = append(categories, cat)
		}
	}
	c.ResourceCategories = categories
}

func (c *Config) genConfigFromTags(specs []*loads.Document) error {
	log.Printf("Using OpenAPI extension tags to configure.")

	// build the apis from the observed groups
	groupsMap := map[ApiGroup]DefinitionList{}
	for _, d := range c.Definitions.All {
		if strings.HasSuffix(d.Name, "List") {
			continue
		}
		if strings.HasSuffix(d.Name, "Status") {
			continue
		}
		if strings.HasPrefix(d.Description(), "Deprecated. Please use") {
			// Don't look at deprecated types
			continue
		}
		if err := d.initExample(c); err != nil {
			return fmt.Errorf("failed to init example: %w", err)
		}
		g := d.Group
		groupsMap[g] = append(groupsMap[g], d)
	}

	groupsList := ApiGroups{}
	for g := range groupsMap {
		groupsList = append(groupsList, g)
	}

	sort.Sort(groupsList)

	for _, g := range groupsList {
		groupName := titleCase(string(g))
		c.ApiGroups = append(c.ApiGroups, ApiGroup(groupName))
		rc := ResourceCategory{
			Include: string(g),
			Name:    groupName,
		}
		defList := groupsMap[g]
		sort.Sort(defList)
		for _, d := range defList {
			r := &Resource{
				Name:       d.Name,
				Group:      string(d.Group),
				Version:    string(d.Version),
				Definition: d,
			}
			rc.Resources = append(rc.Resources, r)
		}
		c.ResourceCategories = append(c.ResourceCategories, rc)
	}

	return nil
}

func (config *Config) initOperationsFromTags(specs []*loads.Document) error {
	if *UseTags {
		ops := map[string]map[string][]*Operation{}
		defs := map[string]*Definition{}
		for _, d := range config.Definitions.All {
			name := fmt.Sprintf("%s.%s.%s", d.Group, d.Version, d.GetResourceName())
			defs[name] = d
		}

		VisitOperations(specs, func(operation Operation) {
			if o, found := config.Operations[operation.ID]; found && o.Definition != nil {
				return
			}
			op := operation
			o := &op
			config.Operations[operation.ID] = o
			group, version, kind, sub := o.GetGroupVersionKindSub()
			if sub == "status" {
				return
			}
			if len(group) == 0 {
				return
			}
			key := fmt.Sprintf("%s.%s.%s", group, version, kind)
			o.Definition = defs[key]

			// Index by group and subresource
			if _, f := ops[key]; !f {
				ops[key] = map[string][]*Operation{}
			}
			ops[key][sub] = append(ops[key][sub], o)
		})

		for key, subMap := range ops {
			def := defs[key]
			if def == nil {
				return fmt.Errorf("unable to locate resource %s in resource map: %v", key, defs)
			}
			subs := []string{}
			for s := range subMap {
				subs = append(subs, s)
			}
			sort.Strings(subs)
			for _, s := range subs {
				cat := &OperationCategory{}
				cat.Name = titleCase(s) + " Operations"
				for _, op := range subMap[s] {
					ot := OperationType{}
					ot.Name = op.GetMethod() + " " + titleCase(s)
					op.Type = ot
					cat.Operations = append(cat.Operations, op)
				}
				def.OperationCategories = append(def.OperationCategories, cat)
			}
		}
	}

	return nil
}

// initOperations returns all Operations found in the Documents
func (c *Config) initOperations(specs []*loads.Document) error {
	c.Operations = Operations{}
	VisitOperations(specs, func(op Operation) {
		c.Operations[op.ID] = &op
	})

	if err := c.mapOperationsToDefinitions(); err != nil {
		return err
	}

	if err := c.initOperationsFromTags(specs); err != nil {
		return err
	}

	VisitOperations(specs, func(target Operation) {
		if op, ok := c.Operations[target.ID]; !ok || op.Definition == nil {
			if !c.OpExcluded(op.ID) {
				log.Printf("No definition found for operation %s [%s]", op.ID, op.Path)
			} else {
				log.Printf("Operation excluded: %s", op.ID)
			}
		}
	})

	if err := c.initOperationParameters(specs); err != nil {
		return err
	}

	// Clear the operations.  We still have to calculate the operations because that is how we determine
	// the API Group for each definition.
	if !*BuildOps {
		c.Operations = Operations{}
		c.OperationCategories = []OperationCategory{}
		for _, d := range c.Definitions.All {
			d.OperationCategories = []*OperationCategory{}
		}
	}

	return nil
}

func (c *Config) OpExcluded(op string) bool {
	for _, pattern := range c.ExcludedOperations {
		if strings.Contains(op, pattern) {
			return true
		}
	}
	return false
}

// CleanUp sorts and dedups fields
func (c *Config) CleanUp() {
	for _, d := range c.Definitions.All {
		sort.Sort(d.AppearsIn)
		sort.Sort(d.Fields)
		dedup := SortDefinitionsByName{}
		var last *Definition
		for _, i := range d.AppearsIn {
			if last != nil &&
				i.Name == last.Name &&
				i.Group.String() == last.Group.String() &&
				i.Version.String() == last.Version.String() {
				continue
			}
			last = i
			dedup = append(dedup, i)
		}
		d.AppearsIn = dedup
	}
}

// LoadConfigFromYAML reads the config yaml file into a struct
func LoadConfigFromYAML() (*Config, error) {
	config := &Config{}

	f := filepath.Join(VersionedConfigDir, "config.yaml")
	contents, err := os.ReadFile(f)
	if err != nil {
		if !*UseTags {
			return nil, fmt.Errorf("failed to read yaml file %s: %w", f, err)
		}
	} else if err = yaml.Unmarshal(contents, config); err != nil {
		return nil, err
	}

	writeCategory := OperationCategory{
		Name: "Write Operations",
		OperationTypes: []OperationType{
			{
				Name:  "Create",
				Match: "create${group}${version}(Namespaced)?${resource}",
			},
			{
				Name:  "Create Eviction",
				Match: "create${group}${version}(Namespaced)?${resource}Eviction",
			},
			{
				Name:  "Patch",
				Match: "patch${group}${version}(Namespaced)?${resource}",
			},
			{
				Name:  "Replace",
				Match: "replace${group}${version}(Namespaced)?${resource}",
			},
			{
				Name:  "Delete",
				Match: "delete${group}${version}(Namespaced)?${resource}",
			},
			{
				Name:  "Delete Collection",
				Match: "delete${group}${version}Collection(Namespaced)?${resource}",
			},
		},
	}

	readCategory := OperationCategory{
		Name: "Read Operations",
		OperationTypes: []OperationType{
			{
				Name:  "Read",
				Match: "read${group}${version}(Namespaced)?${resource}",
			},
			{
				Name:  "List",
				Match: "list${group}${version}(Namespaced)?${resource}",
			},
			{
				Name:  "List All Namespaces",
				Match: "list${group}${version}(Namespaced)?${resource}ForAllNamespaces",
			},
			{
				Name:  "Watch",
				Match: "watch${group}${version}(Namespaced)?${resource}",
			},
			{
				Name:  "Watch List",
				Match: "watch${group}${version}(Namespaced)?${resource}List",
			},
			{
				Name:  "Watch List All Namespaces",
				Match: "watch${group}${version}(Namespaced)?${resource}ListForAllNamespaces",
			},
		},
	}

	statusCategory := OperationCategory{
		Name: "Status Operations",
		OperationTypes: []OperationType{
			{
				Name:  "Patch Status",
				Match: "patch${group}${version}(Namespaced)?${resource}Status",
			},
			{
				Name:  "Read Status",
				Match: "read${group}${version}(Namespaced)?${resource}Status",
			},
			{
				Name:  "Replace Status",
				Match: "replace${group}${version}(Namespaced)?${resource}Status",
			},
		},
	}

	resizeCategory := OperationCategory{
		Name: "Resize Operations",
		OperationTypes: []OperationType{
			{
				Name:  "Read Resize",
				Match: "read${group}${version}(Namespaced)?${resource}Resize",
			},
			{
				Name:  "Patch Resize",
				Match: "patch${group}${version}(Namespaced)?${resource}Resize",
			},
			{
				Name:  "Replace Resize",
				Match: "replace${group}${version}(Namespaced)?${resource}Resize",
			},
		},
	}

	ephemeralCategory := OperationCategory{
		Name: "EphemeralContainers Operations",
		OperationTypes: []OperationType{
			{
				Name:  "Patch EphemeralContainers",
				Match: "patch${group}${version}(Namespaced)?${resource}Ephemeralcontainers",
			},
			{
				Name:  "Read EphemeralContainers",
				Match: "read${group}${version}(Namespaced)?${resource}Ephemeralcontainers",
			},
			{
				Name:  "Replace EphemeralContainers",
				Match: "replace${group}${version}(Namespaced)?${resource}Ephemeralcontainers",
			},
		},
	}

	config.OperationCategories = append(
		[]OperationCategory{
			writeCategory,
			readCategory,
			statusCategory,
			resizeCategory,
			ephemeralCategory,
		},
		config.OperationCategories...,
	)

	return config, nil
}

func (c *Config) initOperationParameters(specs []*loads.Document) error {
	s := c.Definitions
	for _, op := range c.Operations {
		pathItem := op.item

		location := ""
		var param spec.Parameter
		var found bool
		// Path parameters
		for _, p := range pathItem.Parameters {
			if p.In == "" {
				paramID := strings.Split(p.Ref.String(), "/")[2]
				swagger := specs[0].Spec()
				if param, found = swagger.Parameters[paramID]; found {
					location = param.In
				}
			} else {
				location = p.In
				param = p
			}

			switch location {
			case "path":
				op.PathParams = append(op.PathParams, s.parameterToField(param))
			case "query":
				op.QueryParams = append(op.QueryParams, s.parameterToField(param))
			case "body":
				op.BodyParams = append(op.BodyParams, s.parameterToField(param))
			default:
				return fmt.Errorf("unknown location %q", location)
			}
		}

		// Query parameters
		location = ""
		for _, p := range op.op.Parameters {
			if p.In == "" {
				paramID := strings.Split(p.Ref.String(), "/")[2]
				swagger := specs[0].Spec()
				if param, found = swagger.Parameters[paramID]; found {
					location = param.In
				}
			} else {
				location = p.In
				param = p
			}

			switch location {
			case "path":
				op.PathParams = append(op.PathParams, s.parameterToField(param))
			case "query":
				op.QueryParams = append(op.QueryParams, s.parameterToField(param))
			case "body":
				op.BodyParams = append(op.BodyParams, s.parameterToField(param))
			default:
				return fmt.Errorf("unknown location %q", location)
			}
		}

		for code, response := range op.op.Responses.StatusCodeResponses {
			if response.Schema == nil {
				continue
			}
			r := &HttpResponse{
				Field: Field{
					Description: strings.ReplaceAll(response.Description, "\n", " "),
					Type:        GetTypeName(*response.Schema),
					Name:        fmt.Sprintf("%d", code),
				},
				Code: fmt.Sprintf("%d", code),
			}
			if IsComplex(*response.Schema) {
				r.Definition, _ = s.GetForSchema(*response.Schema)
				if r.Definition != nil {
					r.Definition.FoundInOperation = true
				}
			}
			op.HttpResponses = append(op.HttpResponses, r)
		}
	}

	return nil
}

func (c *Config) getOperationGroupName(group string) string {
	for k, v := range c.OperationGroupMap {
		if strings.ToLower(group) == k {
			return v
		}
	}
	return titleCase(group)
}

func (c *Config) getOperationId(match string, group string, version ApiVersion, kind string) string {
	ver := []rune(string(version))
	ver[0] = unicode.ToUpper(ver[0])

	match = strings.ReplaceAll(match, "${group}", group)
	match = strings.ReplaceAll(match, "${version}", string(ver))
	match = strings.ReplaceAll(match, "${resource}", kind)
	return match
}

func (c *Config) setOperation(match, namespace string, ot *OperationType, oc *OperationCategory, d *Definition) error {
	key := strings.ReplaceAll(match, "(Namespaced)?", namespace)
	if o, ok := c.Operations[key]; ok {
		// Each operation should have exactly 1 definition
		if o.Definition != nil {
			return fmt.Errorf(
				"found multiple matching definitions [%s/%s/%s, %s/%s/%s] for operation key: %s",
				d.Group, d.Version, d.Name, o.Definition.Group, o.Definition.Version, o.Definition.Name, key)
		}
		o.Type = *ot
		o.Definition = d
		if err := o.initExample(c); err != nil {
			return fmt.Errorf("failed to init example: %w", err)
		}
		oc.Operations = append(oc.Operations, o)

		// When using tags for the configuration, everything with an operation goes in the ToC
		if *UseTags && !o.Definition.IsOldVersion {
			o.Definition.InToc = true
		}
	}

	return nil
}

// mapOperationsToDefinitions adds operations to the definitions they operate
func (c *Config) mapOperationsToDefinitions() error {
	for _, d := range c.Definitions.All {
		if d.IsInlined {
			continue
		}

		// XXX: The TokenRequest definition has operation defined as "createCoreV1NamespacedServiceAccountToken"!
		if d.Name == "TokenRequest" && d.Group.String() == "authentication" && d.Version == "v1" {
			operationId := "createCoreV1NamespacedServiceAccountToken"
			if o, ok := c.Operations[operationId]; ok {
				ot := OperationType{
					Name:  "Create",
					Match: "createCoreV1NamespacedServiceAccountToken",
				}
				oc := OperationCategory{
					Name:           "Write Operations",
					OperationTypes: []OperationType{ot},
				}

				o.Definition = d
				o.Definition.InToc = true
				if err := o.initExample(c); err != nil {
					return fmt.Errorf("failed to init example: %w", err)
				}
				oc.Operations = append(oc.Operations, o)
			}
			continue
		}

		for i := range c.OperationCategories {
			oc := c.OperationCategories[i]
			for j := range oc.OperationTypes {
				ot := oc.OperationTypes[j]
				groupName := c.getOperationGroupName(d.Group.String())
				operationId := c.getOperationId(ot.Match, groupName, d.Version, d.Name)
				if err := c.setOperation(operationId, "Namespaced", &ot, &oc, d); err != nil {
					return err
				}
				if err := c.setOperation(operationId, "", &ot, &oc, d); err != nil {
					return err
				}
			}

			if len(oc.Operations) > 0 {
				d.OperationCategories = append(d.OperationCategories, &oc)
			}
		}
	}

	return nil
}

// The OpenAPI spec has escape sequences like \u003c. When the spec is unmarshaled,
// the escape sequences get converted to ordinary characters. For example,
// \u003c gets converted to a regular < character. But we can't use  regular <
// and > characters in our HTML document. This function replaces these regular
// characters with HTML entities: <, >, &, ', and ".
func (c *Config) escapeDescriptions() {
	for _, d := range c.Definitions.All {
		d.DescriptionWithEntities = html.EscapeString(d.Description())

		for _, f := range d.Fields {
			f.DescriptionWithEntities = html.EscapeString(f.Description)
		}
	}

	for _, op := range c.Operations {
		for _, p := range op.BodyParams {
			p.DescriptionWithEntities = html.EscapeString(p.Description)
		}
		for _, p := range op.QueryParams {
			p.DescriptionWithEntities = html.EscapeString(p.Description)
		}
		for _, p := range op.PathParams {
			p.DescriptionWithEntities = html.EscapeString(p.Description)
		}
		for _, r := range op.HttpResponses {
			r.DescriptionWithEntities = html.EscapeString(r.Description)
		}
	}
}

// For each resource in the ToC, look up its definition and visit it.
func (c *Config) visitResourcesInToc() error {
	missing := false
	for _, cat := range c.ResourceCategories {
		for _, r := range cat.Resources {
			if d, ok := c.Definitions.GetByVersionKind(r.Group, r.Version, r.Name); ok {
				d.InToc = true // Mark as in Toc
				if err := d.initExample(c); err != nil {
					return fmt.Errorf("failed to init example: %w", err)
				}
				r.Definition = d
			} else {
				fmt.Printf("\033[31mCould not find definition for resource in TOC: %s %s %s.\033[0m\n", r.Group, r.Version, r.Name)
				missing = true
			}
		}
	}
	if missing {
		fmt.Printf("\033[36mAll known definitions: %v\033[0m\n", c.Definitions.All)
	}

	return nil
}
