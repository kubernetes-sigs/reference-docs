package kwebsite

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/stoewer/go-strcase"
)

type PartIndex struct {
	Title  string
	Weight int
}

func (o *KWebsite) addPartIndex(subdir string, name string, weight int) error {
	dirname := filepath.Join(o.Directory, subdir)
	err := os.Mkdir(dirname, 0755)
	if err != nil {
		return err
	}

	t := template.Must(template.ParseFiles(filepath.Join(o.TemplatesDir, "part-index.tmpl")))
	data := PartIndex{
		Title:  name,
		Weight: weight,
	}

	filename := filepath.Join(o.Directory, subdir, "_index.md")
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.Execute(f, data)
}

// escapeName returns a name usable as file name
func escapeName(parts ...string) string {
	result := []string{}
	for _, s := range parts {
		if s != "" {
			result = append(result, strcase.KebabCase(s))
		}
	}
	return strings.Join(result, "-")
}

// headingID returns the ID built by hugo for a given header
func headingID(s string) string {
	result := strings.ReplaceAll(s, " ", "-")
	return result
}
