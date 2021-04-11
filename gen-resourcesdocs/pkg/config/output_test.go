package config_test

import (
	"testing"

	"github.com/kubernetes-sigs/reference-docs/gen-resourcesdocs/pkg/config"
	"github.com/kubernetes-sigs/reference-docs/gen-resourcesdocs/pkg/kubernetes"
	"github.com/kubernetes-sigs/reference-docs/gen-resourcesdocs/pkg/outputs"
)

type FakeOutput struct{}

func (o FakeOutput) Prepare() error                                   { return nil }
func (o FakeOutput) NewPart(i int, name string) (outputs.Part, error) { return FakePart{}, nil }
func (o FakeOutput) AddPart(i int, name string) (outputs.Part, error) { return FakePart{}, nil }
func (o FakeOutput) Terminate() error                                 { return nil }

type FakePart struct{}

func (o FakePart) AddChapter(i int, name string, gv string, version *kubernetes.APIVersion, description string, importPrefix string, domain string) (outputs.Chapter, error) {
	return FakeChapter{}, nil
}

type FakeChapter struct{}

func (o FakeChapter) SetAPIVersion(s string) error { return nil }
func (o FakeChapter) SetGoImport(s string) error   { return nil }
func (o FakeChapter) AddSection(i int, name string, apiVersion *string) (outputs.Section, error) {
	return FakeSection{}, nil
}
func (o FakeChapter) Write() error { return nil }

type FakeSection struct{}

func (o FakeSection) AddContent(s string) error                              { return nil }
func (o FakeSection) AddTypeDefinition(typ string, description string) error { return nil }
func (o FakeSection) AddFieldCategory(name string) error                     { return nil }

func (o FakeSection) AddProperty(name string, property *kubernetes.Property, linkend []string, indent int, defname string, shortName string) error {
	return nil
}
func (o FakeSection) EndProperty() error       { return nil }
func (o FakeSection) StartPropertyList() error { return nil }
func (o FakeSection) EndPropertyList() error   { return nil }
func (o FakeSection) AddOperation(operation *kubernetes.ActionInfo, linkends kubernetes.LinkEnds) error {
	return nil
}
func (o FakeSection) AddDefinitionIndexEntry(d string) error { return nil }

func TestOutputDocumentV119(t *testing.T) {
	outputDocumentVersion(t, "v1.19")
}

func TestOutputDocumentV120(t *testing.T) {
	outputDocumentVersion(t, "v1.20")
}

func outputDocumentVersion(t *testing.T, version string) {
	spec, err := kubernetes.NewSpec("../../api/" + version + "/swagger.json")
	if err != nil {
		t.Errorf("Error loding swagger file")
	}

	toc, err := config.LoadTOC("../../config/" + version + "/toc.yaml")
	if err != nil {
		t.Errorf("LoadTOC should not fail")
	}

	err = toc.PopulateAssociates(spec)
	if err != nil {
		t.Errorf("%s", err)
	}
	toc.Definitions = &spec.Swagger.Definitions
	toc.OutputDocument(FakeOutput{})

}
