package config

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/kubernetes-sigs/reference-docs/gen-resourcesdocs/pkg/kubernetes"
	"github.com/kubernetes-sigs/reference-docs/gen-resourcesdocs/pkg/outputs"
)

// OutputDocument outputs contents using output
func (o *TOC) OutputDocument(output outputs.Output) error {
	for p, tocPart := range o.Parts {
		if err := o.OutputPart(p, tocPart, output); err != nil {
			return err
		}
	}

	if err := o.OutputCommonParameters(len(o.Parts), output); err != nil {
		return err
	}

	return output.Terminate()
}

// OutputPart outputs a Part
func (o *TOC) OutputPart(i int, part *Part, output outputs.Output) error {
	outputPart, err := output.AddPart(i, part.Name)
	if err != nil {
		return err
	}

	for c, tocChapter := range part.Chapters {
		if err = o.OutputChapter(c, tocChapter, outputPart); err != nil {
			return err
		}
	}
	return nil
}

// OutputChapter outputs a chapter of the part
func (o *TOC) OutputChapter(i int, chapter *Chapter, outputPart outputs.Part) error {
	description := ""
	if len(chapter.Sections) > 0 {
		description = getEscapedFirstPhrase(chapter.Sections[0].Definition.Description)
	}
	gv := ""
	if chapter.Group != nil && chapter.Version != nil {
		gv = GetGV(*chapter.Group, *chapter.Version)
	}
	outputChapter, err := outputPart.AddChapter(i, chapter.Name, gv, chapter.Version, description, chapter.Key.GoImportPrefix())
	if err != nil {
		return err
	}

	if chapter.Group != nil && chapter.Version != nil {
		if err = outputChapter.SetAPIVersion(GetGV(*chapter.Group, *chapter.Version)); err != nil {
			return err
		}
	}
	if err = outputChapter.SetGoImport(chapter.Key.GoImportPrefix()); err != nil {
		return err
	}

	for s, tocSection := range chapter.Sections {
		if err = o.OutputSection(s, tocSection, outputChapter); err != nil {
			return err
		}
	}

	if chapter.Group != nil && chapter.Version != nil {
		gvkString := chapter.Group.String() + "." + chapter.Version.String() + "." + chapter.Name
		actions := o.Actions.Get(gvkString)
		if actions != nil {
			if err := o.OutputOperations(len(chapter.Sections), outputChapter, &actions); err != nil {
				return err
			}
		}
	}

	return outputChapter.Write()
}

// OutputSection outputs a section of the chapter
func (o *TOC) OutputSection(i int, section *Section, outputChapter outputs.Chapter) error {
	var apiVersion *string
	if section.Group != nil && section.Version != nil {
		a := GetGV(*section.Group, *section.Version)
		apiVersion = &a
	}
	outputSection, err := outputChapter.AddSection(i, section.Name, apiVersion)
	if err != nil {
		return err
	}

	if err = outputSection.AddDefinitionIndexEntry(section.Name); err != nil {
		return err
	}
	if err = outputSection.AddContent(section.Definition.Description); err != nil {
		return err
	}

	return o.OutputProperties(section.Name, section.Definition, outputSection, []string{}, section.Group, section.Version, section.Key)
}

// OutputProperties outputs the properties of a definition
func (o *TOC) OutputProperties(defname string, definition spec.Schema, outputSection outputs.Section, prefix []string, group *kubernetes.APIGroup, version *kubernetes.APIVersion, key *kubernetes.Key) error {
	requiredProperties := definition.Required

	var apiVersion *string
	if group != nil && version != nil {
		a := GetGV(*group, *version)
		apiVersion = &a
	}

	// Search configured categories
	var fieldCategories []FieldCategory
	if key != nil {
		fieldCategories = o.Categories.Find(*key)

		if fieldCategories != nil {
			if err := checkAllFieldsPresent(fieldCategories, definition.Properties); err != nil {
				return fmt.Errorf("error on fields configuration: %s", err)
			}
		}
	}

	if fieldCategories == nil {
		// Categories config not found, create a default one
		ordered := orderedPropertyKeys(requiredProperties, definition.Properties, apiVersion != nil)
		fieldCategories = []FieldCategory{
			{
				Name:   "",
				Fields: ordered,
			},
		}
	}

	for _, fieldCategory := range fieldCategories {

		if len(prefix) == 0 {
			// NOTE: category names are not displayed for sub-fields (that would be a hell of a mess...)
			if fieldCategory.Name != "" {
				if err := outputSection.AddFieldCategory(fieldCategory.Name); err != nil {
					return err
				}
			}

			if err := outputSection.StartPropertyList(); err != nil {
				return err
			}
		}

		for _, name := range fieldCategory.Fields {
			if apiVersion != nil && (name == "apiVersion" || name == "kind") {
				var property *kubernetes.Property
				if name == "apiVersion" {
					property = kubernetes.NewHardCodedValueProperty(name, *apiVersion)
				} else if name == "kind" {
					property = kubernetes.NewHardCodedValueProperty(name, defname)
				}
				if err := outputSection.AddProperty(name, property, []string{}, 0, defname, name); err != nil {
					return err
				}
				continue
			}

			details := definition.Properties[name]
			property, err := kubernetes.NewProperty(name, details, requiredProperties)
			if err != nil {
				return err
			}
			var linkend []string
			if property.TypeKey != nil {
				linkend = o.LinkEnds[*property.TypeKey]
			}
			completeName := prefix
			completeName = append(completeName, name)
			if err = outputSection.AddProperty(strings.Join(completeName, "."), property, linkend, len(prefix), defname, name); err != nil {
				return err
			}

			if property.TypeKey != nil && len(linkend) == 0 {
				// The type is documented inline
				if target, found := (*o.Definitions)[property.TypeKey.String()]; found {
					o.setDocumentedDefinition(property.TypeKey, defname+"/"+strings.Join(completeName, "."))

					indexedType := property.Type
					indexedType = strings.TrimPrefix(indexedType, "[]")
					indexedType = strings.TrimPrefix(indexedType, "map[string]")

					if err = outputSection.AddDefinitionIndexEntry(indexedType); err != nil {
						return err
					}

					if err = outputSection.AddTypeDefinition(property.TypeKey.ResourceName(), target.Description); err != nil {
						return err
					}

					sublist := false
					if len(prefix) == 0 {
						sublist = true
						if err := outputSection.StartPropertyList(); err != nil {
							return err
						}
					} else if err = outputSection.EndProperty(); err != nil {
						return err
					}

					if err := o.OutputProperties(defname, target, outputSection, append(prefix, name), nil, nil, property.TypeKey); err != nil {
						return err
					}

					if sublist {
						if err := outputSection.EndPropertyList(); err != nil {
							return err
						}
					}
				}
			}
			if err = outputSection.EndProperty(); err != nil {
				return err
			}
		}
		if len(prefix) == 0 {
			if err := outputSection.EndPropertyList(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *TOC) setDocumentedDefinition(key *kubernetes.Key, from string) {
	o.DocumentedDefinitions[*key] = append(o.DocumentedDefinitions[*key], from)
}

// OutputOperations outputs the Operations chapter
func (o *TOC) OutputOperations(i int, outputChapter outputs.Chapter, operations *kubernetes.ActionInfoList) error {
	operationsSection, err := outputChapter.AddSection(i, "Operations", nil)
	if err != nil {
		return err
	}
	for i, operation := range *operations {
		if err := o.OutputOperation(i, operationsSection, &operation); err != nil {
			return err
		}
	}
	return nil
}

// OutputOperation outputs details of an Operation
func (o *TOC) OutputOperation(i int, outputSection outputs.Section, operation *kubernetes.ActionInfo) error {
	return outputSection.AddOperation(operation, o.LinkEnds)
}

// OutputCommonParameters outputs the parameters in common
func (o *TOC) OutputCommonParameters(i int, output outputs.Output) error {
	outputPart, err := output.NewPart(i, "Common Parameters")
	if err != nil {
		return err
	}

	outputChapter, err := outputPart.AddChapter(i, "Common Parameters", "", nil, "", "")
	if err != nil {
		return err
	}

	params := make([]string, len(kubernetes.ParametersAnnex))
	j := 0
	for k := range kubernetes.ParametersAnnex {
		params[j] = k
		j++
	}
	sort.Strings(params)
	for i, param := range params {
		if len(kubernetes.ResourcesDescriptions[param][0].Description) == 0 {
			continue
		}
		outputSection, err := outputChapter.AddSection(i, param, nil)
		if err != nil {
			return err
		}
		if err = outputSection.AddContent(kubernetes.ResourcesDescriptions[param][0].Description); err != nil {
			return err
		}
	}

	return outputChapter.Write()
}

// orderedPropertyKeys returns the keys of m alphabetically ordered
// keys in required will be placed first
func orderedPropertyKeys(required []string, m map[string]spec.Schema, isResource bool) []string {
	sort.Strings(required)

	if isResource {
		mkeys := make(map[string]struct{})
		for k := range m {
			mkeys[k] = struct{}{}
		}
		for _, special := range []string{"metadata", "kind", "apiVersion"} {
			if !isRequired(special, required) {
				if _, ok := mkeys[special]; ok {
					required = append([]string{special}, required...)
				}
			}
		}
	}

	keys := make([]string, len(m)-len(required))
	i := 0
	for k := range m {
		if !isRequired(k, required) {
			keys[i] = k
			i++
		}
	}
	sort.Strings(keys)
	return append(required, keys...)
}

// isRequired returns true if k is in the required array
func isRequired(k string, required []string) bool {
	for _, r := range required {
		if r == k {
			return true
		}
	}
	return false
}

func checkAllFieldsPresent(configuredFields []FieldCategory, definedFields map[string]spec.Schema) error {
	already := map[string]struct{}{}

	count := 0
	for _, category := range configuredFields {
		for _, field := range category.Fields {
			if _, found := already[field]; found {
				return fmt.Errorf("field %s found twice", field)
			}
			already[field] = struct{}{}
			if _, found := definedFields[field]; !found {
				return fmt.Errorf("field %s not defined in Spec", field)
			}
			count++
		}
	}
	if len(definedFields) != count {
		forgotten := []string{}
		for defined := range definedFields {
			if _, found := already[defined]; !found {
				forgotten = append(forgotten, defined)
			}
		}
		return fmt.Errorf("%d fields configured but %d fields in Spec, missing %v", count, len(definedFields), forgotten)
	}
	return nil
}
