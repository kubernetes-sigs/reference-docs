package kwebsite

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kubernetes-sigs/reference-docs/gen-resourcesdocs/pkg/outputs"
)

// KWebsite output
// implements the Output interface
type KWebsite struct {
	Directory    string
	TemplatesDir string
}

// NewKWebsite returns a new KWebsite
func NewKWebsite(dir string, templatesDir string) *KWebsite {
	return &KWebsite{Directory: dir, TemplatesDir: templatesDir}
}

// NewPart creates a new part for the output
func (o *KWebsite) NewPart(i int, name string) (outputs.Part, error) {
	partname := escapeName(name)
	dirname := filepath.Join(o.Directory, partname)
	if err := os.Mkdir(dirname, 0755); err != nil {
		return nil, err
	}
	return Part{
		kwebsite: o,
		name:     partname,
	}, nil
}

// AddPart adds a part to the output
func (o *KWebsite) AddPart(i int, name string) (outputs.Part, error) {
	partname := escapeName(name)
	if err := o.addPartIndex(partname, name, i+1); err != nil {
		return Part{}, fmt.Errorf("Error writing index file for part %s: %s", name, err)
	}
	return Part{
		kwebsite: o,
		name:     partname,
	}, nil
}

// Terminate kwebsite document
func (o *KWebsite) Terminate() error {
	return nil
}
