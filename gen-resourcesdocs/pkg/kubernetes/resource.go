package kubernetes

import (
	"fmt"
	"sort"

	"github.com/go-openapi/spec"
)

// Resource represent a Kubernetes API resource
type Resource struct {
	Key Key
	GVKExtension
	Definition spec.Schema

	// Replaced indicates if this version is replaced by another one
	ReplacedBy *Key
	// Documented indicates if this resource was included in the TOC
	Documented bool
}

// LessThan returns true if 'o' is a newer version than 'p'
func (o *Resource) LessThan(p *Resource) bool {
	return o.Group.Replaces(p.Group) || (o.Group == p.Group && p.Version.LessThan(&o.Version))
}

// Replaces returns true if 'o' replaces 'p'
func (o *Resource) Replaces(p *Resource) bool {
	return o.Group.Replaces(p.Group) || o.Version.Replaces(&p.Version)
}

// Equals returns true if a resource is referenced by group/version/kind
func (o *Resource) Equals(group APIGroup, version APIVersion, kind APIKind) bool {
	return o.Group == group && o.Version.Equals(&version) && o.Kind == kind
}

// GetGV returns the group/version of a resource (used for apiVersion:)
func (o *Resource) GetGV() string {
	if o.Group == "" {
		return o.Version.String()
	}
	return fmt.Sprintf("%s/%s", o.Group, o.Version.String())
}

// ResourceList is the list of resources for a given Kind
type ResourceList []*Resource

func (a ResourceList) Len() int           { return len(a) }
func (a ResourceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ResourceList) Less(i, j int) bool { return a[i].LessThan(a[j]) }

// ResourceMap contains a map of resources, classified by Kind
type ResourceMap map[APIKind]ResourceList

// Add a resource to the resource list
func (o *ResourceMap) Add(resource *Resource) {
	list, ok := (*o)[resource.Kind]
	if ok {
		for _, otherResource := range list {
			if resource.Replaces(otherResource) {
				otherResource.ReplacedBy = &resource.Key
			} else if otherResource.Replaces(resource) {
				resource.ReplacedBy = &otherResource.Key
			}
		}
		list = append(list, resource)
	} else {
		list = []*Resource{resource}
	}
	sort.Sort(list)
	(*o)[resource.Kind] = list
}
