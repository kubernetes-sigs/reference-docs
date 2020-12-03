package config

import (
	"io/ioutil"

	"github.com/feloy/kubernetes-api-reference/pkg/kubernetes"
	"gopkg.in/yaml.v2"
)

// FieldCategory is a list of fields regrouped in the same category
type FieldCategory struct {
	Name   string   `yaml:"name"`
	Fields []string `yaml:"fields"`
}

// Category is the list of fields categories for a specific definition
type Category struct {
	Definition      kubernetes.Key  `yaml:"definition"`
	FieldCategories []FieldCategory `yaml:"field_categories"`
}

// Categories is the list of fields categories for all definitions
type Categories []Category

// LoadCategories from a configuration file
func LoadCategories(filenames []string) (Categories, error) {
	var result Categories
	for _, filename := range filenames {
		var fileCats Categories
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			return result, err
		}
		err = yaml.Unmarshal(content, &fileCats)
		if err != nil {
			return result, err
		}
		result = append(result, fileCats...)
	}
	return result, nil
}

func (o Categories) Find(key kubernetes.Key) []FieldCategory {
	for _, category := range o {
		if category.Definition == key {
			return category.FieldCategories
		}
	}
	return nil
}
