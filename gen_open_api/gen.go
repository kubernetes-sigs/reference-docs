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

package gen_open_api

import (
	"flag"
	"github.com/pwittrock/kubedocs/lib"
	"github.com/pwittrock/kubedocs/gen_open_api/api"
)

var openApiDir = flag.String("open-api-dir", "", "Directory containing open-api specs.")

func GenerateSlateFiles() {
	// Load the yaml config
	config := api.NewConfig(*lib.YamlFile, *openApiDir)

	PrintInfo(config)
	WriteTemplates(config)
}
