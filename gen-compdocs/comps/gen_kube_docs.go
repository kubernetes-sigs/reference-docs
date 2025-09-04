/*
Copyright 2014 The Kubernetes Authors.

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

package comps

import (
	"context"
	goflag "flag"
	"fmt"
	"os"

	"github.com/kubernetes-sigs/reference-docs/gen-compdocs/generators"
	"github.com/spf13/pflag"
	cliflag "k8s.io/component-base/cli/flag"
	kubectlcmd "k8s.io/kubectl/pkg/cmd"
	"k8s.io/kubernetes/cmd/genutils"
	apiservapp "k8s.io/kubernetes/cmd/kube-apiserver/app"
	cmapp "k8s.io/kubernetes/cmd/kube-controller-manager/app"
	schapp "k8s.io/kubernetes/cmd/kube-scheduler/app"
	kubeadmapp "k8s.io/kubernetes/cmd/kubeadm/app/cmd"
	kubeletapp "k8s.io/kubernetes/cmd/kubelet/app"
)

func GenerateFiles(path, module string) error {
	outDir, err := genutils.OutDir(path)
	if err != nil {
		return fmt.Errorf("failed to get output directory: %w", err)
	}

	switch module {
	case "kube-apiserver":
		pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
		apiserver := apiservapp.NewAPIServerCommand()
		return generators.GenMarkdownTree(apiserver, outDir, true)

	case "kube-controller-manager":
		pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
		controllermanager := cmapp.NewControllerManagerCommand()
		return generators.GenMarkdownTree(controllermanager, outDir, true)

	case "kube-scheduler":
		pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
		scheduler := schapp.NewSchedulerCommand()
		return generators.GenMarkdownTree(scheduler, outDir, true)

	case "kubelet":
		kubelet := kubeletapp.NewKubeletCommand(context.TODO())
		return generators.GenMarkdownTree(kubelet, outDir, true)

	case "kubeadm":
		pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
		pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
		// Ignore irrelevant flags
		// pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)

		pflag.Set("logtostderr", "true")
		// We do not want these flags to show up in --help
		// These MarkHidden calls must be after the lines above
		pflag.CommandLine.MarkHidden("version")
		pflag.CommandLine.MarkHidden("log-flush-frequency")
		pflag.CommandLine.MarkHidden("alsologtostderr")
		pflag.CommandLine.MarkHidden("log-backtrace-at")
		pflag.CommandLine.MarkHidden("log-dir")
		pflag.CommandLine.MarkHidden("logtostderr")
		pflag.CommandLine.MarkHidden("stderrthreshold")
		pflag.CommandLine.MarkHidden("vmodule")

		// generate docs for kubeadm
		kubeadm := kubeadmapp.NewKubeadmCommand(os.Stdin, os.Stdout, os.Stderr)
		if err := generators.GenMarkdownTree(kubeadm, outDir, false); err != nil {
			return fmt.Errorf("failed to generate markdown tree: %w", err)
		}

		// cleanup generated code for usage as include in the website
		return generators.MarkdownPostProcessing(kubeadm, outDir, "", generators.CleanupForInclude)

	case "kubectl":
		kubectl := kubectlcmd.NewDefaultKubectlCommand()
		pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
		pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)

		return generators.GenMarkdownTree(kubectl, outDir, true)

	default:
		return fmt.Errorf("module %s is not supported", module)
	}
}
