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
	"k8s.io/klog/v2"
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
	WriteOrphanedOperationsOverview() error
	WriteDefinition(d *api.Definition) error
	WriteOperation(o *api.Operation) error
	WriteOldVersionsOverview() error
	Finalize() error
}

func GenerateFiles() error {
	// load the yaml config
	config, err := api.NewConfig()
	if err != nil {

		return fmt.Errorf("failed to load config: %w", err)
	}

	PrintInfo(config)

	if err := ensureDirectories(); err != nil {

		return fmt.Errorf("failed to ensure directories: %w", err)
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

		return fmt.Errorf("failed to write overview: %w", err)
	}

	// write API groups
	if err := writer.WriteAPIGroupVersions(config.Definitions.GroupVersions); err != nil {
		return fmt.Errorf("failed to write API group versions: %w", err)
	}

	// write resource definitions
	for _, c := range config.ResourceCategories {
		if err := writer.WriteResourceCategory(c.Name, c.Include); err != nil {
			return fmt.Errorf("failed to write resource category '%s': %w", c.Name, err)
		}

		for _, r := range c.Resources {
			if r.Definition == nil {
				// Use klog for structured logging instead of fmt.Printf
				klog.Warningf("Missing definition for item in TOC %s", r.Name)
				continue
			}
			if err := writer.WriteResource(r); err != nil {

				return fmt.Errorf("failed to write resource '%s': %w", r.Name, err)
			}
		}
	}

	// write orphaned operation endpoints
	orphanedIDs := make([]string, 0)
	for id, o := range config.Operations {
		if o.Definition == nil && !config.OpExcluded(o.ID) {
			orphanedIDs = append(orphanedIDs, id)
		}
	}

	if len(orphanedIDs) > 0 {
		if err := writer.WriteOrphanedOperationsOverview(); err != nil {

			return fmt.Errorf("failed to write orphaned operations overview: %w", err)
		}

		sort.Strings(orphanedIDs)

		for _, opKey := range orphanedIDs {
			if err := writer.WriteOperation(config.Operations[opKey]); err != nil {

				return fmt.Errorf("failed to write orphaned operation '%s': %w", opKey, err)
			}
		}
	}

	if err := writer.WriteDefinitionsOverview(); err != nil {

		return fmt.Errorf("failed to write definitions overview: %w", err)
	}

	// add other definition imports
	definitions := api.SortDefinitionsByName{}
	for _, d := range config.Definitions.All {
		// don't add definitions for top level resources in the toc or inlined resources
		if d.InToc || d.IsInlined || d.IsOldVersion {
			continue
		}
		definitions = append(definitions, d)
	}
	sort.Sort(definitions)
	for _, d := range definitions {
		if err := writer.WriteDefinition(d); err != nil {

			return fmt.Errorf("failed to write definition '%s': %w", d.Name, err)
		}
	}

	if err := writer.WriteOldVersionsOverview(); err != nil {

		return fmt.Errorf("failed to write old versions overview: %w", err)
	}

	oldversions := api.SortDefinitionsByName{}
	for _, d := range config.Definitions.All {
		// don't add definitions for top level resources in the toc or inlined resources
		if d.IsOldVersion {
			oldversions = append(oldversions, d)
		}
	}
	sort.Sort(oldversions)
	for _, d := range oldversions {
		// skip Inlined definitions
		if d.IsInlined {
			continue
		}
		r := &api.Resource{Definition: d, Name: d.Name}
		if err := writer.WriteResource(r); err != nil {

			return fmt.Errorf("failed to write old version resource '%s': %w", d.Name, err)
		}
	}

	if err := writer.Finalize(); err != nil {
		// add context to finalize errors
		return fmt.Errorf("failed to finalize writer: %w", err)
	}

	return nil
}

func ensureDirectories() error {
	if err := os.MkdirAll(api.BuildDir, os.FileMode(0700)); err != nil {

		return fmt.Errorf("failed to create build dir '%s': %w", api.BuildDir, err)
	}
	if err := os.MkdirAll(api.IncludesDir, os.FileMode(0700)); err != nil {

		return fmt.Errorf("failed to create includes dir '%s': %w", api.IncludesDir, err)
	}
	return nil
}

func definitionFileName(d *api.Definition) string {
	name := "generated_" + strings.ToLower(strings.ReplaceAll(d.Name, ".", "_"))
	return fmt.Sprintf("%s_%s_%s_definition", name, d.Version, d.Group)
}

func operationFileName(o *api.Operation) string {
	name := "generated_" + strings.ToLower(strings.ReplaceAll(o.ID, ".", "_"))
	return fmt.Sprintf("%s_operation", name)
}

func conceptFileName(d *api.Definition) string {
	name := "generated_" + strings.ToLower(strings.ReplaceAll(d.Name, ".", "_"))
	return fmt.Sprintf("%s_%s_%s_concept", name, d.Version, d.Group)
}

func getLink(s string) string {
	tmp := strings.ReplaceAll(s, ".", "-")
	return strings.ToLower(strings.ReplaceAll(tmp, " ", "-"))
}

func writeStaticFile(filename, defaultContent string) error {
	src := filepath.Join(api.SectionsDir, filename)
	dst := filepath.Join(api.IncludesDir, filename)

	// only try to read the source file, handle error if it doesn't exist (removes double syscall)
	content, readErr := os.ReadFile(src)
	if readErr == nil {
		// if file exists and is readable, use its content
		defaultContent = string(content)
	} else if !os.IsNotExist(readErr) {
		return fmt.Errorf("failed to read source file %s: %w", src, readErr)
	}

	// structured logging using klog instead of fmt.Printf for consistency
	klog.Info("Creating file ", dst)

	if err := os.WriteFile(dst, []byte(defaultContent), 0644); err != nil {
		return fmt.Errorf("failed to write static file '%s': %w", dst, err)
	}
	return nil
}
