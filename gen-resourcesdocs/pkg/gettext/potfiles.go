package gettext

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/kubernetes-sigs/reference-docs/gen-resourcesdocs/pkg/kubernetes"
)

type Definition struct {
	Name   string
	Schema spec.Schema
}

type PotFiles struct {
	Path  string
	Files map[string][]Definition
}

func NewPotFiles(path string) PotFiles {
	return PotFiles{
		Path:  path,
		Files: map[string][]Definition{},
	}
}

func (o PotFiles) Add(key kubernetes.Key, definition spec.Schema) {
	domain := key.RemoveResourceName()
	if _, ok := o.Files[domain]; !ok {
		o.Files[domain] = []Definition{
			{
				Name:   key.ResourceName(),
				Schema: definition,
			},
		}
	} else {
		o.Files[domain] = append(o.Files[domain], Definition{
			Name:   key.ResourceName(),
			Schema: definition,
		})
	}
}

func (o PotFiles) CreateFiles() error {
	for k, definitions := range o.Files {
		err := createPoFile(o.Path, k, definitions)
		if err != nil {
			return err
		}
	}
	return nil
}

func createPoFile(path string, domain string, definitions []Definition) error {
	f, err := os.Create(filepath.Join(path, domain+".pot"))
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, `# Translations for Kubernetes API specifications.
#
#, fuzzy
msgid ""
msgstr ""
"Project-Id-Version: %s\n"
"PO-Revision-Date: YEAR-MO-DA HO:MI +ZONE\n"
"Last-Translator: FULL NAME <EMAIL@ADDRESS>\n"
"Language-Team: LANGUAGE <LL@li.org>\n"
"MIME-Version: 1.0\n"
"Content-Type: text/plain; charset=UTF-8\n"
"Content-Transfer-Encoding: 8bit\n"
`, domain)

	createdEntries := &PoEntries{}

	sort.Slice(definitions, func(i, j int) bool {
		return definitions[i].Name < definitions[j].Name
	})
	for _, definition := range definitions {
		writeDefinition(f, definition, createdEntries)
	}
	return nil
}

func writeDefinition(f io.Writer, definition Definition, createdEntries *PoEntries) {

	description := definition.Schema.Description
	addEntry(f, createdEntries, definition.Name, description)

	properties := definition.Schema.Properties
	keys := make([]string, 0, len(properties))
	for k := range properties {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if len(properties[k].Description) == 0 {
			continue
		}
		entry := fmt.Sprintf("%s.%s", definition.Name, k)
		addEntry(f, createdEntries, entry, properties[k].Description)
	}
}

func addEntry(f io.Writer, createdEntries *PoEntries, entry string, description string) {

	ok, old := createdEntries.Add(description, entry)
	if ok {
		lines := strings.Split(description, "\n")
		fmt.Fprintf(f, "\n#: %s\n", entry)
		if len(lines) == 1 {
			fmt.Fprintf(f, "msgid \"%s\"\n", escapeMsg(description))
		} else {
			fmt.Fprintf(f, "msgid \"\"\n")
			for _, line := range lines {
				fmt.Fprintf(f, "\"%s\"\n", escapeMsg(line))
			}
		}
		fmt.Fprint(f, "msgstr \"\"\n")
	} else {
		fmt.Fprintf(f, "\n# %s\n", entry)
		fmt.Fprintf(f, "# same as %s\n", *old)
	}
}

func escapeMsg(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}
