/*
Copyright 2017 The Kubernetes Authors.

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
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// MarkdownPostProcessing goes though the generated files
func MarkdownPostProcessing(cmd *cobra.Command, dir string, subdir string, processor func(string) string) error {
	for _, c := range cmd.Commands() {
		// if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() { // Qiming
		if !c.IsAvailableCommand() { // Qiming
			continue
		}

		parts := strings.Split(c.CommandPath(), " ")
		subdir := ""
		if len(parts) > 1 {
			subdir = parts[0] + "_" + parts[1]
		}

		if err := MarkdownPostProcessing(c, dir, subdir, processor); err != nil {
			return err
		}
	}

	// basename := strings.ReplaceAll(cmd.CommandPath(), " ", "_") + ".md"
	// filename := filepath.Join(dir, basename)

	fullname := strings.ReplaceAll(cmd.CommandPath(), " ", "_") + ".md"
	// indexFile := false
	if len(subdir) > 0 {
		parts := strings.Split(cmd.CommandPath(), " ")
		if len(parts) == 2 {
			// indexFile = true
			fullname = "_index.md"
		}
		fullname = filepath.Join(subdir, fullname)
	}
	filename := filepath.Join(dir, fullname)

	markdownBytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("filename: %s - %s", filename, err)
		return err
	}

	processedMarkDown := processor(string(markdownBytes))

	return os.WriteFile(filename, []byte(processedMarkDown), 0644)
}

// CleanupForInclude parts of markdown that will make difficult to use it as include in the website:
// - The title of the document (this allow more flexibility for include, e.g. include in tabs)
// - The sections see also, that assumes file will be used as a main page
func CleanupForInclude(md string) string {
	lines := strings.Split(md, "\n")

	cleanMd := ""
	for i, line := range lines {
		if i == 0 {
			continue
		}

		if strings.HasSuffix(strings.ToUpper(line), "SEE ALSO") {
			break
		}

		cleanMd += line
		if i < len(lines)-1 {
			cleanMd += "\n"
		}
	}

	return cleanMd
}
