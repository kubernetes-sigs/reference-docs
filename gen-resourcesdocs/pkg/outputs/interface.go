package outputs

import "github.com/kubernetes-sigs/reference-docs/gen-resourcesdocs/pkg/kubernetes"

// Output is an interface for output formats
type Output interface {
	AddPart(i int, name string) (Part, error)
	NewPart(i int, name string) (Part, error)
	Terminate() error
}

// Part is an interface to a part of an output
type Part interface {
	AddChapter(i int, name string, gv string, version *kubernetes.APIVersion, description string, importPrefix string, domain string) (Chapter, error)
}

// Chapter is an interface to a chapter of an output
type Chapter interface {
	SetAPIVersion(s string) error
	SetGoImport(s string) error
	AddSection(i int, name string, apiVersion *string) (Section, error)
	Write() error
}

// Section is an interface to a section of an output
type Section interface {
	AddContent(s string) error
	AddTypeDefinition(typ string, description string) error
	StartPropertyList() error
	AddFieldCategory(name string) error
	AddProperty(name string, property *kubernetes.Property, linkend []string, indent int, defname string, shortName string) error
	EndProperty() error
	EndPropertyList() error
	AddOperation(operation *kubernetes.ActionInfo, linkends kubernetes.LinkEnds) error
	AddDefinitionIndexEntry(d string) error
}
