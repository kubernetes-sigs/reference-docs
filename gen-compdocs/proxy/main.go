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
		fmt.Fprintf(os.Stderr, "usage: %s [output-dir]\n", os.Args[0])
		os.Exit(1)
	}

	GenKubeProxy(path)
}

func GenKubeProxy(path string) {
	outDir, err := genutils.OutDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get output directory: %v\n", err)
		os.Exit(1)
	}

	proxy := proxyapp.NewProxyCommand()
	generators.GenMarkdownTree(proxy, outDir, true)
}
