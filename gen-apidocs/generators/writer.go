/*
Copyright 2018 The Kubernetes Authors.

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
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kubernetes-sigs/reference-docs/gen-apidocs/generators/api"
)

type Doc struct {
	Filename string `json:"filename,omitempty"`
}

type DocWriter interface {
	Extension() string
	DefaultStaticContent(title string) string
	WriteOverview() error
	WriteAPIGroupVersions(gvs api.GroupVersions) error
	WriteResourceCategory(name, file string) error
	WriteResource(r *api.Resource) error
	WriteDefinitionsOverview() error
	WriteDefinition(d *api.Definition) error
	WriteOldVersionsOverview() error
	Finalize() error
}

func GenerateFiles() error {
	// Load the yaml config
	config, err := api.NewConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	PrintInfo(config)

	if err := ensureDirectories(); err != nil {
		return err
	}

	copyright_tmpl := "<a href=\"https://github.com/kubernetes/kubernetes\">Copyright 2016-%s The Kubernetes Authors.</a>"
	now := time.Now().Format("2006")
	copyright := fmt.Sprintf(copyright_tmpl, now)
	var title string
	if !*api.BuildOps {
		title = "Kubernetes Resource Reference Docs"
	} else {
		title = "Kubernetes API Reference Docs"
	}

	writer := NewHTMLWriter(config, copyright, title)
	if err := writer.WriteOverview(); err != nil {
		return err
	}

	// Write API groups
	if err := writer.WriteAPIGroupVersions(config.Definitions.GroupVersions); err != nil {
		return err
	}

	// Write resource definitions
	for _, c := range config.ResourceCategories {
		if err := writer.WriteResourceCategory(c.Name, c.Include); err != nil {
			return err
		}
		for _, r := range c.Resources {
			if r.Definition == nil {
				fmt.Printf("Warning: Missing definition for item in TOC %s\n", r.Name)
				continue
			}
			if err := writer.WriteResource(r); err != nil {
				return err
			}
		}
	}

	if err := writer.WriteDefinitionsOverview(); err != nil {
		return err
	}

	// Add other definition imports
	definitions := api.SortDefinitionsByName{}
	for _, d := range config.Definitions.All {
		// Don't add definitions for top level resources in the toc or inlined resources
		if d.InToc || d.IsInlined || d.IsOldVersion {
			continue
		}
		definitions = append(definitions, d)
	}
	sort.Sort(definitions)
	for _, d := range definitions {
		if err := writer.WriteDefinition(d); err != nil {
			return err
		}
	}

	if err := writer.WriteOldVersionsOverview(); err != nil {
		return err
	}

	oldversions := api.SortDefinitionsByName{}
	for _, d := range config.Definitions.All {
		// Don't add definitions for top level resources in the toc or inlined resources
		if d.IsOldVersion {
			oldversions = append(oldversions, d)
		}
	}
	sort.Sort(oldversions)
	for _, d := range oldversions {
		// Skip Inlined definitions
		if d.IsInlined {
			continue
		}
		r := &api.Resource{Definition: d, Name: d.Name}
		if err := writer.WriteResource(r); err != nil {
			return err
		}
	}

	if err := writer.Finalize(); err != nil {
		return err
	}

	return nil
}

func ensureDirectories() error {
	if err := os.MkdirAll(api.BuildDir, os.FileMode(0700)); err != nil {
		return err
	}
	if err := os.MkdirAll(api.IncludesDir, os.FileMode(0700)); err != nil {
		return err
	}

	return nil
}

func definitionFileName(d *api.Definition) string {
	name := "generated_" + strings.ToLower(strings.ReplaceAll(d.Name, ".", "_"))
	return fmt.Sprintf("%s_%s_%s_definition", name, d.Version, d.Group)
}

func conceptFileName(d *api.Definition) string {
	name := "generated_" + strings.ToLower(strings.ReplaceAll(d.Name, ".", "_"))
	return fmt.Sprintf("%s_%s_%s_concept", name, d.Version, d.Group)
}

func getLink(s string) string {
	tmp := strings.ReplaceAll(s, ".", "-")
	return strings.ToLower(strings.ReplaceAll(tmp, " ", "-"))
}

func writeStaticFile(title, location, defaultContent string) error {
	fn := filepath.Join(api.SectionsDir, location)
	to := filepath.Join(api.IncludesDir, location)
	_, err := os.Stat(fn)
	if err == nil {
		// copy the file if it exists
		return os.Link(fn, to)
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat file %s: %w", fn, err)
	}

	fmt.Printf("Creating file %s\n", to)

	return os.WriteFile(to, []byte(defaultContent), 0644)
}
