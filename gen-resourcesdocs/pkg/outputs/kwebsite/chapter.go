package kwebsite

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/kubernetes-sigs/reference-docs/gen-resourcesdocs/pkg/outputs"
)

// Chapter of a KWebsite output
// implements the outputs.Chapter interface
type Chapter struct {
	kwebsite *KWebsite
	part     *Part
	name     string
	data     *ChapterData
	domain   string
}

type ChapterData struct {
	ApiVersion  string
	Version     string
	Import      string
	Kind        string
	Metadata    ChapterMetadata
	ChapterName string
	Sections    []SectionData
}

type ChapterMetadata struct {
	Description string
	Title       string
	Weight      int
}

type SectionData struct {
	Name            string
	Description     string
	Fields          []FieldData
	FieldCategories []FieldCategoryData
	Operations      []OperationData
}

type FieldCategoryData struct {
	Name   string
	Fields []FieldData
}

type FieldData struct {
	Name           string
	Value          string
	Description    string
	Type           string
	TypeDefinition string
	Indent         int
}

type OperationData struct {
	Verb          string
	Title         string
	RequestMethod string
	RequestPath   string
	Parameters    []ParameterData
	Responses     []ResponseData
}

type ParameterData struct {
	Title       string
	Description string
}

type ResponseData struct {
	Code        int
	Type        string
	Description string
}

// SetAPIVersion writes the APIVersion for a chapter
func (o Chapter) SetAPIVersion(s string) error {
	return nil
}

// SetGoImport writes the Go import for a chapter
func (o Chapter) SetGoImport(s string) error {
	return nil
}

// AddSection adds a section to the chapter
func (o Chapter) AddSection(i int, name string, apiVersion *string) (outputs.Section, error) {
	o.data.Sections = append(o.data.Sections, SectionData{
		Name: name,
	})

	return Section{
		kwebsite: o.kwebsite,
		part:     o.part,
		chapter:  &o,
	}, nil
}

func (o Chapter) Write() error {
	chaptername := escapeName(o.data.ChapterName, o.data.Version)
	filename := filepath.Join(o.kwebsite.Directory, o.part.name, chaptername) + ".md"
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	templateFile := "chapter.tmpl"
	if len(o.data.Sections) == 1 {
		templateFile = "chapter-single-definition.tmpl"
	}
	t := template.Must(template.New(templateFile).Funcs(sprig.TxtFuncMap()).ParseFiles(filepath.Join(o.kwebsite.TemplatesDir, templateFile)))
	return t.Execute(f, o.data)
}
