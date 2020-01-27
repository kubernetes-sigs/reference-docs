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
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

var KubernetesRelease = flag.String("kubernetes-release", "", "Version of Kubernetes to generate docs for.")

//var MyPages = Pages{}

var GenKubectlDir = flag.String("gen-kubectl-dir", "gen-kubectldocs", "Directory containing kubectl build files")

// Directory for output files
var BuildDir string

// Directory for configuration and data files
var ConfigDir string

// Versioned directory for configuration file
var VersionedConfigDir string

func getTocFile() string {
	return filepath.Join(VersionedConfigDir, "toc.yaml")
}

func getStaticIncludesDir() string {
	return filepath.Join(VersionedConfigDir, "static_includes")
}

type Doc struct {
	Filename string `json:"filename,omitempty"`
}

type DocWriter interface {
	Extension() string
	WriteCommands(toc ToC, params KubectlSpec)
	Finalize()
}

func GenerateFiles() {

	BuildDir = filepath.Join(*GenKubectlDir, "build")
	ConfigDir = filepath.Join(*GenKubectlDir, "config")

	var versionChar = "v"
	var k8sRelease = fmt.Sprintf("%s%s", versionChar, strings.ReplaceAll(*KubernetesRelease, ".", "_"))
	VersionedConfigDir = filepath.Join(ConfigDir, k8sRelease)

	// get the kubectl command specification
	spec := GetSpec()

	// categories from toc.yaml
	toc := ToC{}

	// REVISIT
	if len(getTocFile()) < 1 {
		fmt.Printf("Must have toc.yaml file, %s", getTocFile())
		os.Exit(2)
	}

	ensureBuildDirs()

	contents, err := ioutil.ReadFile(getTocFile())
	if err != nil {
		fmt.Printf("Failed to read yaml file %s: %v", getTocFile(), err)
	}

	err = yaml.Unmarshal(contents, &toc)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	NormalizeSpec(&spec)

	copyright_tmpl := "<a href=\"https://github.com/kubernetes/kubernetes\">&#xa9;Copyright 2016-%s The Kubernetes Authors.</a>"
	now := time.Now().Format("2006")
	copyright := fmt.Sprintf(copyright_tmpl, now)

	var title = "Kubectl Reference Docs"

	writer := NewHTMLWriter(copyright, title)
	writer.WriteCommands(toc, spec)
	writer.Finalize()
}

func NormalizeSpec(spec *KubectlSpec) {
	for _, g := range spec.TopLevelCommandGroups {
		for _, c := range g.Commands {
			FormatCommand(c.MainCommand)
			for _, s := range c.SubCommands {
				FormatCommand(s)
			}
		}
	}
}

func FormatCommand(c *Command) {
	c.Example = FormatExample(c.Example)
	c.Description = FormatDescription(c.Description)
	c.Options = FormatOptions(c.Options)
}

func FormatDescription(input string) string {
	/* This fixes an error when the description is a string followed by a
	   new line and another string that is indented >= four spaces. The marked.js parser
	   throws a parsing error. Error found in generated file: build/_generated_rollout.md */
	input = strings.Replace(input, "\n   ", "\n ", 10)

	input = strings.Replace(input, "|", "&#124;", -1)

	return strings.Replace(input, "   *", "*", 10000)
}

// REVISIT, simplify the md parsing and html formatting of examples
func FormatExample(input string) string {
	last := ""
	result := ""

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if len(line) < 1 {
			continue
		}

		// Skip empty lines
		if strings.HasPrefix(line, "#") {
			if len(strings.TrimSpace(strings.Replace(line, "#", "", 1))) < 1 {
				continue
			}
		}

		if strings.HasPrefix(line, "#") {

			if last == "command" {
				// Close command if it is open
				result += "\n</code></pre></DIV>"
			}

			// have text describing example
			result += "<DIV class=\"cmd-example-text\"><p><h4>" + strings.Replace(line, "#", "", 1) + "</h4></p></DIV>"

			last = "comment"
		} else {
			// codeblock
			if last != "command" {
				// Open a new code section
				result += "<DIV class=\"cmd-example-code\"><pre><code>"
			}
			result += "\n" + line

			// added the first command, more?
			last = "command"
		}
	}

	// Close the final command
	if last == "command" {
		result += "\n</code></pre></DIV>\n"
	}
	return result
}

func FormatOptions(options Options) Options {
	for _, o := range options {
		o.Usage = strings.Replace(o.Usage, "|", "&#124;", -1)
	}
	return options
}

func ensureBuildDirs() {
	if _, err := os.Stat(*GenKubectlDir); os.IsNotExist(err) {
		os.Mkdir(*GenKubectlDir, os.FileMode(0700))
	}
	if _, err := os.Stat(BuildDir); os.IsNotExist(err) {
		os.Mkdir(BuildDir, os.FileMode(0700))
	}
	if _, err := os.Stat(ConfigDir); os.IsNotExist(err) {
		os.Mkdir(ConfigDir, os.FileMode(0700))
	}
	if _, err := os.Stat(VersionedConfigDir); os.IsNotExist(err) {
		os.Mkdir(VersionedConfigDir, os.FileMode(0700))
	}
}
