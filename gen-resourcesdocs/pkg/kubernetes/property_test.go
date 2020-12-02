package kubernetes

import (
	"reflect"
	"testing"

	"github.com/go-openapi/jsonreference"
	"github.com/go-openapi/spec"
)

func TestIsRequired(t *testing.T) {
	tests := []struct {
		Name     string
		Required []string
		Expected bool
	}{
		{
			Name:     "found",
			Required: []string{"a", "found", "c"},
			Expected: true,
		},
		{
			Name:     "notfound",
			Required: []string{"a", "b", "c"},
		},
	}

	for _, test := range tests {
		result := isRequired(test.Name, test.Required)
		if result != test.Expected {
			t.Errorf("Should be %v but is %v", test.Expected, result)
		}
	}
}

func TestGetTypeNameAndKey(t *testing.T) {
	tests := []struct {
		Definition   spec.Schema
		ExpectedName string
		ExpectedKey  *Key
	}{
		{
			Definition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{
						"boolean",
					},
				},
			},
			ExpectedName: "boolean",
			ExpectedKey:  nil,
		},
		{
			Definition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{
						"array",
					},
					Items: &spec.SchemaOrArray{
						Schema: &spec.Schema{
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{
									"string",
								},
							},
						},
					},
				},
			},
			ExpectedName: "[]string",
			ExpectedKey:  nil,
		},
		{
			Definition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Ref: spec.Ref{
						Ref: jsonreference.MustCreateRef("#/definitions/io.k8s.api.core.v1.PodSpec"),
					},
				},
			},
			ExpectedName: "PodSpec",
			ExpectedKey:  func(k Key) *Key { return &k }("io.k8s.api.core.v1.PodSpec"),
		},
	}

	for _, test := range tests {
		resultName, resultKey := GetTypeNameAndKey(test.Definition)
		if resultName != test.ExpectedName {
			t.Errorf("Name should be %s but is %s", test.ExpectedName, resultName)
		}
		if !reflect.DeepEqual(resultKey, test.ExpectedKey) {
			t.Errorf("Key should be %v but is %v", test.ExpectedKey, resultKey)
		}
	}
}

func TestNewProperty(t *testing.T) {
	tests := []struct {
		Name     string
		Details  spec.Schema
		Required []string
		Expected *Property
	}{

		{
			Name: "onefield",
			Details: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{
						"boolean",
					},
					Description: "a description",
				},
			},
			Required: []string{"a", "b"},
			Expected: &Property{
				Name:        "onefield",
				Type:        "boolean",
				TypeKey:     nil,
				Description: "a description",
			},
		},
	}

	for _, test := range tests {
		result, _ := NewProperty(test.Name, test.Details, test.Required)
		if !reflect.DeepEqual(test.Expected, result) {
			t.Errorf("Should be %#v but is %#v", test.Expected, result)
		}
	}
}
