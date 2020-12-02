package kubernetes

import (
	"fmt"
	"strings"

	"github.com/feloy/kubernetes-api-reference/pkg/openapi"
	"github.com/go-openapi/spec"
)

// Spec represents the Kubernetes API Specification
type Spec struct {
	// Swagger is the openAPI representation of the k8s spec
	// populated by calling getSwagger
	Swagger *spec.Swagger

	// Resources is the list of K8s resources
	// populated by calling getResources
	Resources *ResourceMap

	// Actions is the list of endpoints defined in the API
	Actions Actions

	// GVToKey maps beetween Kubernetes Group/Version and Swagger definition key
	GVToKey GVToKeyMap
}

// GVToKeyMap maps Kubernetes resource Group/Version with Spec Definition key (without Kind)
// e.g. GVToKey["v1"]: "io.k8s.api.core.v1"
type GVToKeyMap map[string][]string

// Add adds a new match between key and resource GV
func (o GVToKeyMap) Add(key string, resource *Resource) {
	parts := strings.Split(key, ".")
	if len(parts) == 0 {
		return
	}
	subkey := strings.Join(parts[0:len(parts)-1], ".")
	gv := resource.GetGV()
	if _, found := o[gv]; !found {
		o[gv] = []string{subkey}
	} else {
		found := false
		for _, k := range o[gv] {
			if k == subkey {
				found = true
			}
		}
		if !found {
			o[gv] = append(o[gv], subkey)
		}
	}
}

// NewSpec creates a new Spec from a K8s spec file
func NewSpec(filename string) (*Spec, error) {
	spec := &Spec{}
	err := spec.getSwagger(filename)
	if err != nil {
		return nil, err
	}
	err = spec.getResources()
	if err != nil {
		return nil, err
	}

	err = spec.getActions()
	if err != nil {
		return nil, err
	}

	return spec, nil
}

// GetSwagger populates the swagger representation of the Spec in file `filename`
func (o *Spec) getSwagger(filename string) error {
	var err error
	o.Swagger, err = openapi.LoadOpenAPISpec(filename)
	return err
}

// GetResources populates the resources defined in the spec
// and maps definitions keys to Resources GVs
func (o *Spec) getResources() error {
	o.Resources = &ResourceMap{}
	o.GVToKey = GVToKeyMap{}

	for key, definition := range o.Swagger.Definitions {
		gvk, found, err := getGVKExtension(definition.Extensions)
		if err != nil {
			return fmt.Errorf("%s: %f", key, err)
		}
		if !found {
			continue
		}
		resource := &Resource{
			Key:          Key(key),
			GVKExtension: *gvk,
			Definition:   definition,
		}
		o.Resources.Add(resource)
		o.GVToKey.Add(key, resource)
	}
	return nil
}

// GetResource returns the resource referenced by group/version/kind, or nil if not found
func (o *Spec) GetResource(group APIGroup, version APIVersion, kind APIKind, markAsDocumented bool) (Key, *spec.Schema) {
	if resources, ok := (*o.Resources)[kind]; ok {
		for r, resource := range resources {
			if resource.Equals(group, version, kind) {
				if markAsDocumented {
					(*o.Resources)[kind][r].Documented = true
				}
				return resource.Key, &resource.Definition
			}
		}
	}
	return "", nil
}

// GetDefinition returns the definition referenced by key
func (o *Spec) GetDefinition(key Key) *spec.Schema {
	if s, found := o.Swagger.Definitions[key.String()]; found {
		return &s
	}
	return nil
}

func (o *Spec) getActions() error {
	o.Actions = make(Actions)

	paths := o.Swagger.Paths.Paths
	for key, path := range paths {
		if path.Get != nil {
			o.Actions.Add(key, path.Get, "GET", path.Parameters)
		}
		if path.Put != nil {
			o.Actions.Add(key, path.Put, "PUT", path.Parameters)
		}
		if path.Post != nil {
			o.Actions.Add(key, path.Post, "POST", path.Parameters)
		}
		if path.Delete != nil {
			o.Actions.Add(key, path.Delete, "DELETE", path.Parameters)
		}
		if path.Options != nil {
			o.Actions.Add(key, path.Options, "OPTIONS", path.Parameters)
		}
		if path.Head != nil {
			o.Actions.Add(key, path.Head, "HEAD", path.Parameters)
		}
		if path.Patch != nil {
			o.Actions.Add(key, path.Patch, "PATCH", path.Parameters)
		}
	}

	o.Actions.findCommonParameters()
	return nil
}
