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

//$ go run main.go
//$ docker run -v $(pwd)/source:/slate/source -v $(pwd)/build:/slate/build pwittrock/slate
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kubernetes-incubator/reference-docs/gen_kubectl"
	"github.com/kubernetes-incubator/reference-docs/gen_open_api"
)

var docType = flag.String("doc-type", "open-api", "Type of docs to generate.")

func main() {
	flag.Parse()

	switch *docType {
	case "open-api":
		gen_open_api.GenerateSlateFiles()
	case "kubectl":
		gen_kubectl.GenerateSlateFiles()
	//case "open-api-old":
	//	gen_open_api_old.GenerateSlateFiles()
	default:
		fmt.Printf("Must provide type as either open-api or kubectl, was %s\n", *docType)
		os.Exit(2)
	}
}
