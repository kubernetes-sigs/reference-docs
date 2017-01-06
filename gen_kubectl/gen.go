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

package gen_kubectl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/kubernetes-incubator/reference-docs/lib"
	"gopkg.in/yaml.v2"
	"path/filepath"
	"flag"
)

var GenKubectlDir = flag.String("gen-kubectl-dir", "gen_kubectl/", "Directory containing kubectl files")

func getYamlFile() string {
	return filepath.Join(*GenKubectlDir, *lib.KubernetesVersion, "kubectl.yaml")
}

func getTocFile() string {
	return filepath.Join(*GenKubectlDir, *lib.KubernetesVersion, "toc.yaml")
}

func getTemplateFile(name string) string {
	return filepath.Join(*GenKubectlDir, name)
}

func getStaticIncludesDir() string {
	return filepath.Join(*GenKubectlDir, "static_includes")
}

func GenerateFiles() {
	spec := KubectlSpec{}

	if len(getYamlFile()) < 1 {
		fmt.Printf("Must specify --yaml-file.\n")
		os.Exit(2)
	}

	contents, err := ioutil.ReadFile(getYamlFile())
	if err != nil {
		fmt.Printf("Failed to read yaml file %s: %v", getYamlFile(), err)
	}

	err = yaml.Unmarshal(contents, &spec)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	toc := ToC{}
	if len(getTocFile()) < 1 {
		fmt.Printf("Must specify --toc-file.\n")
		os.Exit(2)
	}

	contents, err = ioutil.ReadFile(getTocFile())
	if err != nil {
		fmt.Printf("Failed to read yaml file %s: %v", getTocFile(), err)
	}

	err = yaml.Unmarshal(contents, &toc)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	manifest := &Manifest{}
	manifest.Title = "Kubectl Reference Docs"
	manifest.Copyright = "<a href=\"https://github.com/kubernetes/kubernetes\">Copyright 2016 The Kubernetes Authors.</a>"

	NormalizeSpec(&spec)

	if _, err := os.Stat(*GenKubectlDir + "/includes"); os.IsNotExist(err) {
		os.Mkdir(*GenKubectlDir + "/includes", os.FileMode(0700))
	}

	WriteCommandFiles(manifest, toc, spec)
	WriteManifest(manifest)
}

func NormalizeSpec(spec *KubectlSpec) {
	for _, g  := range spec.TopLevelCommandGroups {
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
}

func FormatDescription(input string) string {
	return strings.Replace(input, "   *", "*", 10000)
}

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
			if len(strings.TrimSpace(strings.Replace(line, "#", ">bdocs-tab:example", 1))) < 1 {
				continue
			}
		}

		// Format comments as code blocks
		if strings.HasPrefix(line, "#") {
			if last == "command" {
				// Close command if it is open
				result += "\n```\n\n"
			}

			if last == "comment" {
				// Add to the previous code block
				result += " " + line
			} else {
				// Start a new code block
				result += strings.Replace(line, "#", ">bdocs-tab:example", 1)
			}
			last = "comment"
		} else {
			if last != "command" {
				// Open a new code section
				result += "\n\n```bdocs-tab:example_shell"
			}
			result += "\n" + line
			last = "command"
		}
	}

	// Close the final command if needed
	if last == "command" {
		result += "\n```\n"
	}
	return result
}

func WriteManifest(manifest *Manifest) {
	jsonbytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		fmt.Printf("Could not Marshal manfiest %+v due to error: %v.\n", manifest, err)
	} else {
		jsonfile, err := os.Create(*GenKubectlDir + "/" + lib.JsonOutputFile)
		if err != nil {
			fmt.Printf("Could not create file %s due to error: %v.\n", lib.JsonOutputFile, err)
		} else {
			defer jsonfile.Close()
			_, err := jsonfile.Write(jsonbytes)
			if err != nil {
				fmt.Printf("Failed to write bytes %s to file %s: %v.\n", jsonbytes, lib.JsonOutputFile, err)
			}
		}
	}

}

func WriteCommandFiles(manifest *Manifest, toc ToC,  params KubectlSpec) {
	t, err := template.New("command.template").ParseFiles(getTemplateFile("command.template"))
	if err != nil {
		fmt.Printf("Failed to parse template: %v", err)
		os.Exit(1)
	}


	m := map[string]TopLevelCommand{}
	for _, g := range params.TopLevelCommandGroups {
		for _, tlc := range g.Commands {
			m[tlc.MainCommand.Name] = tlc
		}
	}

	for _, c := range toc.Categories {
		// Write the category include
		if len(c.Include) > 0 {

		}

		fn := strings.Replace(c.Name, " ", "_", -1)
		manifest.Docs = append(manifest.Docs, Doc{strings.ToLower(fmt.Sprintf("_generated_category_%s.md", fn))})
		WriteCategoryFile(c)

		// Write each of the commands in this category
		for _, cm := range c.Commands {
			if tlc, found := m[cm]; !found {
				fmt.Printf("Could not find top level command %s\n", cm)
				os.Exit(1)
			} else {
				WriteCommandFile(manifest, t, tlc)
				delete(m, cm)
			}
		}
	}
	if len(m) > 0 {
		for k := range m {
			fmt.Printf("Kubectl command %s missing from table of contents\n", k)
		}
		os.Exit(1)
	}
}

func WriteCategoryFile(c Category) {
	ct, err := template.New("category.template").ParseFiles(getTemplateFile("category.template"))
	if err != nil {
		fmt.Printf("Failed to parse template: %v", err)
		os.Exit(1)
	}

	fn := strings.Replace(c.Name, " ", "_", -1)
	f, err := os.Create(*GenKubectlDir + "/includes/_generated_category_" + strings.ToLower(fmt.Sprintf("%s.md", fn)))
	if err != nil {
		fmt.Printf("Failed to open index: %v", err)
		os.Exit(1)
	}
	defer f.Close()
	err = ct.Execute(f, c)
	if err != nil {
		fmt.Printf("Failed to execute template: %v", err)
		os.Exit(1)
	}
}

func WriteCommandFile(manifest *Manifest, t *template.Template, params TopLevelCommand) {
	params.MainCommand.Description = strings.Replace(params.MainCommand.Description, "|", "&#124;", -1)
	for _, o := range params.MainCommand.Options {
		o.Usage = strings.Replace(o.Usage, "|", "&#124;", -1)
	}
	for _, sc := range params.SubCommands {
		for _, o := range sc.Options {
			o.Usage = strings.Replace(o.Usage, "|", "&#124;", -1)
		}
	}
	f, err := os.Create(*GenKubectlDir + "/includes/_generated_" + strings.ToLower(params.MainCommand.Name) + ".md")
	if err != nil {
		fmt.Printf("Failed to open index: %v", err)
		os.Exit(1)
	}
	defer f.Close()

	err = t.Execute(f, params)
	if err != nil {
		fmt.Printf("Failed to execute template: %v", err)
		os.Exit(1)
	}
	manifest.Docs = append(manifest.Docs, Doc{"_generated_" + strings.ToLower(params.MainCommand.Name) + ".md"})
}
