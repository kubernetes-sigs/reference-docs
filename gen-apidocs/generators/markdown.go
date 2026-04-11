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
	"log"
	"os"
	"path/filepath"

	"github.com/kubernetes-sigs/reference-docs/gen-apidocs/generators/api"
)

type MarkdownWriter struct {
	Config          *api.Config
	OutputDir       string // build/markdown
	currentCategory string
	resourceWeight  int
	categoryWeight  int
	linkMap         map[string]linkInfo
}

type linkInfo struct {
	Category string
	Filename string
	Anchor   string
}

var _ DocWriter = (*MarkdownWriter)(nil)

func NewMarkdownWriter(config *api.Config, copyright, title string) DocWriter {
	outputDir := filepath.Join(api.BuildDir, "markdown")

	writer := MarkdownWriter{
		Config:    config,
		OutputDir: outputDir,
		linkMap:   make(map[string]linkInfo),
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Printf("MarkdownWriter: failed to create output dir %s: %v", outputDir, err)
	}

	return &writer
}

func (m *MarkdownWriter) Extension() string {
	return ".md"
}

func (m *MarkdownWriter) DefaultStaticContent(title string) string {
	return "# " + title + "\n"
}

func (m *MarkdownWriter) WriteOverview() error {
	fmt.Println("MarkdownWriter.WriteOverview")
	return nil
}

func (m *MarkdownWriter) WriteResource(r *api.Resource) error {
	fmt.Printf("MarkdownWriter.WriteResource: %s\n", r.Name)
	return nil
}

func (m *MarkdownWriter) WriteAPIGroupVersions(gvs api.GroupVersions) error {
	fmt.Println("MarkdownWriter.WriteAPIGroupVersions")
	return nil
}

func (m *MarkdownWriter) WriteResourceCategory(name, file string) error {
	fmt.Printf("MarkdownWriter.WriteResourceCategory: %s\n", name)
	return nil
}

func (m *MarkdownWriter) WriteDefinitionsOverview() error {
	fmt.Println("MarkdownWriter.WriteDefinitionsOverview")
	return nil
}

func (m *MarkdownWriter) WriteOrphanedOperationsOverview() error {
	fmt.Println("MarkdownWriter.WriteOrphanedOperationsOverview")
	return nil
}

func (m *MarkdownWriter) WriteDefinition(d *api.Definition) error {
	fmt.Printf("MarkdownWriter.WriteDefinition: %s\n", d.Name)
	return nil
}

func (m *MarkdownWriter) WriteOperation(o *api.Operation) error {
	fmt.Printf("MarkdownWriter.WriteOperation: %s\n", o.ID)
	return nil
}

func (m *MarkdownWriter) WriteOldVersionsOverview() error {
	fmt.Println("MarkdownWriter.WriteOldVersionsOverview")
	return nil
}

func (m *MarkdownWriter) Finalize() error {
	fmt.Println("MarkdownWriter.Finalize")
	return nil
}
