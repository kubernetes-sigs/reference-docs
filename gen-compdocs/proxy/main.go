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

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kubernetes-sigs/reference-docs/gen-compdocs/generators"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/cmd/genutils"
	proxyapp "k8s.io/kubernetes/cmd/kube-proxy/app"
)

func init() {
	klog.InitFlags(nil)
}

func main() {
	// use os.Args instead of "flags" because "flags" will mess up the man pages!
	path := ""
	if len(os.Args) == 2 {
		path = os.Args[1]
	} else {
		log.Fatalf("usage: %s [output-dir]", os.Args[0])
	}

	if err := GenKubeProxy(path); err != nil {
		log.Fatalf("failure: %v", err)
	}
}

func GenKubeProxy(path string) error {
	outDir, err := genutils.OutDir(path)
	if err != nil {
		return fmt.Errorf("failed to get output directory: %w", err)
	}

	proxy := proxyapp.NewProxyCommand()

	return generators.GenMarkdownTree(proxy, outDir, true)
}
