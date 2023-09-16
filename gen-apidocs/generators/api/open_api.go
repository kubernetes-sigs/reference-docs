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
	"strings"

	"github.com/go-openapi/loads"
)

const (
	patchStrategyKey = "x-kubernetes-patch-strategy"
	patchMergeKeyKey = "x-kubernetes-patch-merge-key"
	resourceNameKey  = "x-kubernetes-resource"
	typeKey          = "x-kubernetes-group-version-kind"
)

// Loads all of the open-api documents
func LoadOpenApiSpec() ([]*loads.Document, error) {
	docs := []*loads.Document{}
	err := filepath.Walk(VersionedConfigDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		ext := filepath.Ext(path)
		if ext != ".json" {
			return nil
		}
		var d *loads.Document
		d, err = loads.JSONSpec(path)
		if err != nil {
			return fmt.Errorf("Could not load json file %s as api-spec: %v\n", path, err)
		}
		docs = append(docs, d)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return docs, nil
}

func LoadDefinitions(config *Config, specs []*loads.Document, s *Definitions) error {
	var versionList ApiVersions

	for _, spec := range specs {
		for name, spec := range spec.Spec().Definitions {
			resource := ""
			if r, ok := spec.Extensions.GetString(resourceNameKey); ok {
				resource = r
			}

			// This actually skips the following groups, i.e. old definitions
			//  'io.k8s.kubernetes.pkg.api.*'
			//  'io.k8s.kubernetes.pkg.apis.*'
			if strings.HasPrefix(spec.Description, "Deprecated. Please use") {
				continue
			}

			// NOTE:
			if strings.Contains(name, "JSONSchemaPropsOrStringArray") {
				continue
			}

			group, version, kind := GuessGVK(name)
			if group == "" {
				continue
			} else if group == "error" {
				return fmt.Errorf("could not locate group for %s", name)
			}

			full_group, found := config.GroupFullNames[group]
			if !found {
				// fall back to group name if no mapping found
				fmt.Printf("\033[31mWarning: full name for '%s' not provided, guessing...\033[0m\n", group)
				full_group = group
			}

			d := &Definition{
				schema:        spec,
				Name:          kind,
				Version:       ApiVersion(version),
				Kind:          ApiKind(kind),
				Group:         ApiGroup(group),
				GroupFullName: full_group,
				ShowGroup:     true,
				Resource:      resource,
			}

			s.All[d.Key()] = d

			// skip "io.k8s.apimachinery.pkg.api.resource.*"
			// skip "meta" group also
			if version == "resource" || group == "meta" {
				continue
			}

			versionList, found = s.GroupVersions[full_group]
			if !found {
				s.GroupVersions[full_group] = ApiVersions{ApiVersion(version)}
			} else {
				found = false
				for _, v := range versionList {
					if v.String() == version {
						found = true
					}
				}
				if !found {
					versionList = append(versionList, ApiVersion(version))
					s.GroupVersions[full_group] = versionList
				}
			}
		}
	}

	return nil
}

func ParseSpecInfo(specs []*loads.Document, cfg *Config) {
	// The following loop can be optimized, there is now only one spec for analysis
	for _, spec := range specs {
		cfg.SpecTitle = spec.Spec().Info.InfoProps.Title
	}
}
