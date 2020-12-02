package kubernetes

import (
	"strings"
)

// Key of the resource in the OpenAPI definition (e.g. io.k8s.api.core.v1.Pod)
type Key string

// GoImportPrefix returns the path to use for this group in go import
func (o Key) GoImportPrefix() string {
	parts := strings.Split(o.String(), ".")
	return parts[1] + "." + parts[0] + "/" + strings.Join(parts[2:len(parts)-1], "/")
}

// RemoveResourceName removes the last part of the key corresponding to the resource name
func (o Key) RemoveResourceName() string {
	parts := strings.Split(o.String(), ".")
	return strings.Join(parts[:len(parts)-1], ".")
}

// ResourceName returns the resource name part of a key
func (o Key) ResourceName() string {
	parts := strings.Split(o.String(), ".")
	return parts[len(parts)-1]
}

// String returns a string representation of the Key
func (o Key) String() string {
	return string(o)
}
