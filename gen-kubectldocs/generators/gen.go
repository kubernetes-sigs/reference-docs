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
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/sets"
)

const JsonOutputFile = "manifest.json"

var (
	KubernetesVersion = flag.String("kubernetes-version", "", "Version of Kubernetes to generate docs for.")
	GenKubectlDir     = flag.String("gen-kubectl-dir", "generators", "Directory containing kubectl files")
)

func getTocFile() string {
	return filepath.Join(*GenKubectlDir, *KubernetesVersion, "toc.yaml")
}

func getStaticIncludesDir() string {
	return filepath.Join(*GenKubectlDir, *KubernetesVersion, "static_includes")
}

func GenerateFiles() error {
	if *KubernetesVersion == "" {
		return errors.New("must specify --kubernetes-version")
	}

	spec := GetSpec()

	contents, err := os.ReadFile(getTocFile())
	if err != nil {
		return fmt.Errorf("failed to read TOC file %s: %w", getTocFile(), err)
	}

	toc := ToC{}
	if err = yaml.Unmarshal(contents, &toc); err != nil {
		return fmt.Errorf("failed to unmarshal %s: %w", getTocFile(), err)
	}

	manifest := &Manifest{}
	manifest.Title = "Kubectl Reference Docs"
	manifest.Copyright = "<a href=\"https://github.com/kubernetes/kubernetes\">Copyright 2020 The Kubernetes Authors.</a>"

	NormalizeSpec(&spec)

	if err := os.MkdirAll(*GenKubectlDir+"/includes", os.FileMode(0700)); err != nil {
		return err
	}

	if err := WriteCommandFiles(manifest, toc, spec); err != nil {
		return fmt.Errorf("failed to write command files: %w", err)
	}

	if err := WriteManifest(manifest); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	return nil
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
}

func FormatDescription(input string) string {
	/* This fixes an error when the description is a string followed by a
	   new line and another string that is indented >= four spaces. The marked.js parser
	   throws a parsing error. Error found in generated file: build/_generated_rollout.md */
	input = strings.Replace(input, "\n   ", "\n ", 10)
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

func WriteManifest(manifest *Manifest) error {
	jsonbytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest %+v: %w", manifest, err)
	}

	jsonfile, err := os.Create(*GenKubectlDir + "/" + JsonOutputFile)
	if err != nil {
		return fmt.Errorf("failed to create manifest file %s: %w", JsonOutputFile, err)
	}
	defer jsonfile.Close()

	_, err = jsonfile.Write(jsonbytes)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", JsonOutputFile, err)
	}

	return nil
}

func WriteCommandFiles(manifest *Manifest, toc ToC, params KubectlSpec) error {
	t, err := template.New("command.template").Parse(CommandTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	missingCommands := map[string]TopLevelCommand{}
	for _, g := range params.TopLevelCommandGroups {
		for _, tlc := range g.Commands {
			missingCommands[tlc.MainCommand.Name] = tlc
		}
	}

	err = filepath.Walk(getStaticIncludesDir(), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			to := filepath.Join(*GenKubectlDir, "includes", filepath.Base(path))
			return os.Link(path, to)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to copy includes %w", err)
	}

	for _, c := range toc.Categories {
		if len(c.Include) > 0 {
			// Use the static category include
			manifest.Docs = append(manifest.Docs, Doc{strings.ToLower(c.Include)})
		} else {
			// Write a general category include
			fn := strings.ReplaceAll(c.Name, " ", "_")
			manifest.Docs = append(manifest.Docs, Doc{strings.ToLower(fmt.Sprintf("_generated_category_%s.md", fn))})
			if err := WriteCategoryFile(c); err != nil {
				return fmt.Errorf("failed to write category file: %w", err)
			}
		}

		// Write each of the commands in this category
		for _, cm := range c.Commands {
			if tlc, found := missingCommands[cm]; !found {
				return fmt.Errorf("could not find top level command %s", cm)
			} else {
				if err := WriteCommandFile(manifest, t, tlc); err != nil {
					return fmt.Errorf("failed to write command file: %w", err)
				}
				delete(missingCommands, cm)
			}
		}
	}

	if len(missingCommands) > 0 {
		return fmt.Errorf("kubectl commands %v are missing from table of contents", sets.List(sets.KeySet(missingCommands)))
	}

	return nil
}

func WriteCategoryFile(c Category) error {
	ct, err := template.New("category.template").Parse(CategoryTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	fn := strings.ReplaceAll(c.Name, " ", "_")
	f, err := os.Create(*GenKubectlDir + "/includes/_generated_category_" + strings.ToLower(fmt.Sprintf("%s.md", fn)))
	if err != nil {
		return fmt.Errorf("failed to open index: %w", err)
	}
	defer f.Close()

	err = ct.Execute(f, c)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func WriteCommandFile(manifest *Manifest, t *template.Template, params TopLevelCommand) error {
	replacer := strings.NewReplacer(
		"|", "&#124;",
		"<", "&lt;",
		">", "&gt;",
		"[", "<span>[</span>",
		"]", "<span>]</span>",
		"\n", "<br>",
	)

	params.MainCommand.Description = replacer.Replace(params.MainCommand.Description)
	for _, o := range params.MainCommand.Options {
		o.Usage = replacer.Replace(o.Usage)
	}
	for _, sc := range params.SubCommands {
		for _, o := range sc.Options {
			o.Usage = replacer.Replace(o.Usage)
		}
	}
	f, err := os.Create(*GenKubectlDir + "/includes/_generated_" + strings.ToLower(params.MainCommand.Name) + ".md")
	if err != nil {
		return fmt.Errorf("failed to open index: %w", err)
	}
	defer f.Close()

	err = t.Execute(f, params)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	manifest.Docs = append(manifest.Docs, Doc{"_generated_" + strings.ToLower(params.MainCommand.Name) + ".md"})

	return nil
}
