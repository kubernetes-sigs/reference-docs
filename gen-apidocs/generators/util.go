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
	"fmt"
	"regexp"
	"strings"

	"github.com/kubernetes-sigs/reference-docs/gen-apidocs/generators/api"
)

var (
	anchorRegex    = regexp.MustCompile(`[^a-zA-Z0-9]+`)
	kebabBoundary1 = regexp.MustCompile(`([a-z0-9])([A-Z])`)
	kebabBoundary2 = regexp.MustCompile(`([A-Z])([A-Z][a-z])`)
)

// anchor produces a stable in-page anchor from a kind or section name.
func anchor(s string) string {
	return strings.Trim(anchorRegex.ReplaceAllString(s, "-"), "-")
}

// kebabCase lowercases and slugifies any string into a kebab-case identifier
// used for directory and file names.
func kebabCase(s string) string {
	return strings.Trim(anchorRegex.ReplaceAllString(strings.ToLower(s), "-"), "-")
}

// kebabName converts a CamelCase API kind (e.g. "PodTemplate") into its
// kebab-cased file form (e.g. "pod-template"), preserving acronym boundaries.
func kebabName(s string) string {
	s = kebabBoundary2.ReplaceAllString(s, "$1-$2")
	s = kebabBoundary1.ReplaceAllString(s, "$1-$2")
	return strings.ToLower(s)
}

// groupVersionString formats the canonical "group/version" string used in
// API references, collapsing the core group to a bare version.
func groupVersionString(group string, version api.ApiVersion) string {
	if group == "" || group == "core" {
		return version.String()
	}
	return fmt.Sprintf("%s/%s", group, version.String())
}

// operationSlug derives a filename-safe slug from an operation ID.
func operationSlug(id string) string {
	return strings.Trim(anchorRegex.ReplaceAllString(strings.ToLower(id), "-"), "-")
}

// constValueFor hard-codes the two fields Kubernetes manifests always carry
// with fixed values (apiVersion and kind). Swagger does not tag them as const
// so we derive them from the GVK.
func constValueFor(fieldName, apiVersion, kind string) string {
	switch fieldName {
	case "apiVersion":
		return apiVersion
	case "kind":
		return kind
	}
	return ""
}

func PrintInfo(config *api.Config) {
	definitions := config.Definitions

	// collect orphaned and missing TOC results in slices instead of multiple flags/loops
	var orphaned []string
	var missingToc []struct {
		Name       string
		Operations []string
	}

	// ignored names kept in a slice for clarity and easier maintenance
	ignored := []string{"meta.v1.APIVersions", "meta.v1.Patch"}

	// single loop over definitions.All for efficiency
	for name, d := range definitions.All {
		// orphaned check inlined
		ignore := false
		if !d.FoundInField && !d.FoundInOperation {
			for _, ig := range ignored {
				if strings.Contains(name, ig) {
					ignore = true
					break
				}
			}
			if !ignore {
				orphaned = append(orphaned, name)
			}
		}

		// missing TOC check inlined
		if !d.InToc && len(d.OperationCategories) > 0 && !d.IsOldVersion && !d.IsInlined {
			var ops []string
			for _, oc := range d.OperationCategories {
				for _, o := range oc.Operations {
					ops = append(ops, o.ID)
				}
			}
			missingToc = append(missingToc, struct {
				Name       string
				Operations []string
			}{Name: name, Operations: ops})
		}
	}

	// print orphaned results
	if len(orphaned) > 0 {
		fmt.Println("----------------------------------")
		fmt.Println("Orphaned Definitions:")
		for _, name := range orphaned {
			fmt.Printf("[%s]\n", name)
		}
		if !*api.AllowErrors {
			fmt.Println("Possible orphaned definitions found.")
		}
	}

	// print missing TOC results
	if len(missingToc) > 0 {
		fmt.Println("----------------------------------")
		fmt.Println("Definitions with Operations Missing from Toc (Excluding old version):")
		for _, item := range missingToc {
			fmt.Printf("[%s]\n", item.Name)
			for _, op := range item.Operations {
				fmt.Printf("\t [%s]\n", op)
			}
		}
	}
}
