package kubernetes

import (
	"strings"

	"github.com/go-openapi/spec"
)

// ParameterIn represenets the position of a parameter of an operation
type ParameterIn string

var paramInOrder = map[ParameterIn]int8{
	"path":  0,
	"body":  1,
	"query": 2,
}

// ParametersAnnex indicates the common parameters
// that are displayed in an annex
var ParametersAnnex = map[string]struct{}{}

// LessThan returns true if o appears before p in the natural order
func (o ParameterIn) LessThan(p ParameterIn) bool {
	return paramInOrder[o] < paramInOrder[p]
}

// ParametersList is a list of parameters
type ParametersList []spec.Parameter

func (a ParametersList) Len() int      { return len(a) }
func (a ParametersList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ParametersList) Less(i, j int) bool {
	if a[i].In == a[j].In {
		return a[i].Name < a[j].Name
	}
	return ParameterIn(a[i].In).LessThan(ParameterIn(a[j].In))
}

// Add a parameter to the list
func (a *ParametersList) Add(specParameters map[string]spec.Parameter, parameter spec.Parameter) {
	if !parameter.Ref.GetPointer().IsEmpty() {
		key := Key(strings.TrimPrefix(parameter.Ref.GetPointer().String(), "/parameters/"))
		parameter = specParameters[key.String()]
	}
	desc := parameter.Description
	if strings.Contains(strings.ToLower(desc), "deprecated") {
		return
	}
	*a = append(*a, parameter)
}

// ParameterInAnnex returns true if param is displayed in annex
func ParameterInAnnex(param spec.Parameter) bool {
	_, found := ParametersAnnex[param.Name]
	return found
}

type descriptionInfo struct {
	Description string
	count       int
}

type ResourcesMap map[string][]descriptionInfo

var ResourcesDescriptions = ResourcesMap{}

func (o *ResourcesMap) add(param spec.Parameter) {
	if _, ok := (*o)[param.Name]; !ok {
		(*o)[param.Name] = []descriptionInfo{
			{
				Description: param.Description,
				count:       1,
			},
		}
	} else {
		list := (*o)[param.Name]
		for k, descInfo := range list {
			if descInfo.Description == param.Description {
				(*o)[param.Name][k].count++
				return
			}
		}
		(*o)[param.Name] = append((*o)[param.Name], descriptionInfo{
			Description: param.Description,
			count:       1,
		})
	}
}

func (o ResourcesMap) addActionParameters(params *ParametersList) {
	for _, param := range *params {
		o.add(param)
	}
}
