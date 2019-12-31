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
	"github.com/kubernetes-sigs/reference-docs/gen-apidocs/generators/api"
	"strings"
)

func PrintInfo(config *api.Config) {
	definitions := config.Definitions

	hasOrphaned := false
	for name, d := range definitions.All {
		if !d.FoundInField && !d.FoundInOperation {
			if !strings.Contains(name, "meta.v1.APIVersions") && !strings.Contains(name, "meta.v1.Patch") {
				hasOrphaned = true
			}
		}
	}
	if hasOrphaned {
		fmt.Printf("----------------------------------\n")
		fmt.Printf("Orphaned Definitions:\n")
		for name, d := range definitions.All {
			if !d.FoundInField && !d.FoundInOperation {
				if !strings.Contains(name, "meta.v1.APIVersions") && !strings.Contains(name, "meta.v1.Patch") {
					fmt.Printf("[%s]\n", name)
				}
			}
		}
		if !*api.AllowErrors {
			panic("Orphaned definitions found.")
		}
	}

	missingFromToc := false
	for _, d := range definitions.All {
		if !d.InToc && len(d.OperationCategories) > 0 && !d.IsOldVersion && !d.IsInlined {
			missingFromToc = true
		}
	}

	if missingFromToc {
		fmt.Printf("----------------------------------\n")
		fmt.Printf("Definitions with Operations Missing from Toc (Excluding old version):\n")
		for name, d := range definitions.All {
			if !d.InToc && len(d.OperationCategories) > 0 && !d.IsOldVersion && !d.IsInlined {
				fmt.Printf("[%s]\n", name)
				for _, oc := range d.OperationCategories {
					for _, o := range oc.Operations {
						fmt.Printf("\t [%s]\n", o.ID)
					}
				}
			}
		}
	}

	//fmt.Printf("Old definitions:\n")
	//for name, d := range definitions.All {
	//	if !d.InToc && len(d.OperationCategories) > 0 && d.IsOldVersion && !d.IsInlined {
	//		fmt.Printf("[%s]\n", name)
	//		for _, oc := range d.OperationCategories {
	//			for _, o := range oc.Operations {
	//				fmt.Printf("\t [%s]\n", o.ID)
	//			}
	//		}
	//	}
	//}
}
