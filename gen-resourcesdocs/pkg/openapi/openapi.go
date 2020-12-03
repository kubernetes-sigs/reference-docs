package openapi

import (
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

// LoadOpenAPISpec loads the open-api document
func LoadOpenAPISpec(filename string) (*spec.Swagger, error) {
	d, err := loads.JSONSpec(filename)
	if err != nil {
		return nil, err
	}
	return d.Spec(), nil
}
