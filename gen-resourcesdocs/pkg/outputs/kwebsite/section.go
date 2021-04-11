package kwebsite

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kubernetes-sigs/reference-docs/gen-resourcesdocs/pkg/kubernetes"
	"github.com/leonelquinteros/gotext"
)

// Section of a Hugo output
// implements the outputs.Section interface
type Section struct {
	kwebsite *KWebsite
	part     *Part
	chapter  *Chapter
}

// AddContent adds content to a section
func (o Section) AddContent(s string) error {
	i := len(o.chapter.data.Sections)
	gotext.SetDomain(o.chapter.domain)
	o.chapter.data.Sections[i-1].Description = gotext.Get(removeEOL(s))
	return nil
}

// AddTypeDefinition adds the definition of a type to a section
func (o Section) AddTypeDefinition(typ string, description string) error {
	i := len(o.chapter.data.Sections)
	cats := o.chapter.data.Sections[i-1].FieldCategories
	var fields *[]FieldData
	if len(cats) == 0 {
		fields = &o.chapter.data.Sections[i-1].Fields
	} else {
		fields = &cats[len(cats)-1].Fields
	}
	j := len(*fields)
	(*fields)[j-1].Type = typ
	gotext.SetDomain(o.chapter.domain)
	(*fields)[j-1].TypeDefinition = "*" + gotext.Get(removeEOL(description)) + "*"
	return nil
}

// StartPropertyList starts the list of properties
func (o Section) StartPropertyList() error {
	return nil
}

func (o Section) AddFieldCategory(name string) error {
	i := len(o.chapter.data.Sections)
	o.chapter.data.Sections[i-1].FieldCategories = append(o.chapter.data.Sections[i-1].FieldCategories, FieldCategoryData{
		Name: name,
	})
	return nil
}

// AddProperty adds a property to the section
func (o Section) AddProperty(name string, property *kubernetes.Property, linkend []string, indent int, defname string, shortName string) error {
	if property.HardCodedValue != nil {
		i := len(o.chapter.data.Sections)
		cats := o.chapter.data.Sections[i-1].FieldCategories
		var fields *[]FieldData
		if len(cats) == 0 {
			fields = &o.chapter.data.Sections[i-1].Fields
		} else {
			fields = &cats[len(cats)-1].Fields
		}
		*fields = append(*fields, FieldData{
			Name:   "**" + name + "**",
			Value:  *property.HardCodedValue,
			Indent: 0,
		})
		return nil
	}

	gotext.SetDomain(o.chapter.domain)
	description := gotext.Get(removeEOL(property.Description))

	gotext.SetDomain("gen-resourcesdocs")

	required := ""
	if property.Required {
		required = gotext.Get(", required")
	}

	typ := property.Type
	if property.TypeKey != nil && len(linkend) > 0 {
		typ = o.kwebsite.LinkEnd(linkend, property.Type)
	}
	title := fmt.Sprintf("**%s** (%s)%s", name, typ, required)

	listType := ""
	if property.ListType != nil {
		if *property.ListType == "atomic" {
			listType = gotext.Get("Atomic: will be replaced during a merge")
		} else if *property.ListType == "set" {
			listType = gotext.Get("Set: unique values will be kept during a merge")
		} else if *property.ListType == "map" {
			if len(property.ListMapKeys) == 1 {
				listType = fmt.Sprintf(gotext.Get("Map: unique values on key %s will be kept during a merge"), property.ListMapKeys[0])
			} else {
				listType = fmt.Sprintf(gotext.Get("Map: unique values on keys `%s` will be kept during a merge"), strings.Join(property.ListMapKeys, ", "))
			}
		}
	}
	if len(listType) > 0 {
		description = "*" + listType + "*\n\n" + description
	}

	var patches string
	if property.MergeStrategyKey != nil && property.RetainKeysStrategy {
		patches = fmt.Sprintf(gotext.Get("Patch strategies: retainKeys, merge on key `%s`"), *property.MergeStrategyKey)
	} else if property.MergeStrategyKey != nil {
		patches = fmt.Sprintf(gotext.Get("Patch strategy: merge on key `%s`"), *property.MergeStrategyKey)
	} else if property.RetainKeysStrategy {
		patches = gotext.Get("Patch strategy: retainKeys")
	}

	if len(patches) > 0 {
		description = "*" + patches + "*\n\n" + description
	}

	i := len(o.chapter.data.Sections)
	cats := o.chapter.data.Sections[i-1].FieldCategories
	var fields *[]FieldData
	if len(cats) == 0 {
		fields = &o.chapter.data.Sections[i-1].Fields
	} else {
		fields = &cats[len(cats)-1].Fields
	}
	*fields = append(*fields, FieldData{
		Name:        title,
		Description: description,
		Indent:      indent,
	})
	if indent > 3 {
		fmt.Printf("Warning: %s/%s indentation is %d\n", o.chapter.name, name, indent)
	}
	return nil
}

// EndProperty ends a property
func (o Section) EndProperty() error {
	return nil
}

// EndPropertyList ends the list of properties
func (o Section) EndPropertyList() error {
	return nil
}

// AddOperation adds an operation
func (o Section) AddOperation(operation *kubernetes.ActionInfo, linkends kubernetes.LinkEnds) error {
	sentences := strings.Split(operation.Operation.Description, ".")

	if len(sentences) > 1 {
		fmt.Printf("SHOULD NOT HAPPEN, sentences: %d\n", len(sentences))
	}

	dataParams := []ParameterData{}
	for _, param := range operation.Parameters {

		required := ""
		if param.Required {
			required = ", required"
		}

		typ := param.Type
		if param.Schema != nil {
			t, typeKey := kubernetes.GetTypeNameAndKey(*param.Schema)
			linkend, found := linkends[*typeKey]
			if found {
				typ = o.kwebsite.LinkEnd(linkend, t)
			} else {
				typ = t
				fmt.Printf("SHOULD NOT HAPPEN: %s\n", typeKey)
			}
		}

		desc := param.Description
		if len(desc) > 0 && kubernetes.ParameterInAnnex(param) {
			desc = o.kwebsite.LinkEnd([]string{"common-parameters", "common-parameters"}, param.Name)
		}

		dataParams = append(dataParams, ParameterData{
			Title:       paramName(param.Name, param.In) + ": " + typ + required,
			Description: desc,
		})
	}

	codes := make([]int, len(operation.Operation.Responses.StatusCodeResponses))
	i := 0
	for code := range operation.Operation.Responses.StatusCodeResponses {
		codes[i] = code
		i++
	}
	sort.Ints(codes)
	responsesData := []ResponseData{}
	for _, code := range codes {
		response := operation.Operation.Responses.StatusCodeResponses[code]

		typ := ""
		if response.Schema != nil {
			t, typeKey := kubernetes.GetTypeNameAndKey(*response.Schema)
			if typeKey != nil {
				linkend, found := linkends[*typeKey]
				if found {
					typ = o.kwebsite.LinkEnd(linkend, t)
				} else {
					typ = t
					fmt.Printf("SHOULD NOT HAPPEN: %s\n", typeKey)
				}
			} else {
				typ = t
			}
		}

		responsesData = append(responsesData, ResponseData{
			Code:        code,
			Type:        typ,
			Description: response.Description,
		})
	}

	i = len(o.chapter.data.Sections)
	ops := &o.chapter.data.Sections[i-1].Operations
	*ops = append(*ops, OperationData{
		Verb:          operation.Action.Verb(),
		Title:         sentences[0],
		RequestMethod: operation.HTTPMethod,
		RequestPath:   operation.Path.String(),
		Parameters:    dataParams,
		Responses:     responsesData,
	})

	return nil
}

func (o Section) AddDefinitionIndexEntry(d string) error {
	return nil
}

func paramName(s string, in string) string {
	switch in {
	case "path":
		return fmt.Sprintf("**%s** (*in path*)", s)
	case "query":
		return fmt.Sprintf("**%s** (*in query*)", s)
	default:
		return fmt.Sprintf("**%s**", s)
	}
}

func removeEOL(s string) string {
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	return s
}
