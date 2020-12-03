package kubernetes

// APIKind represents the Kind of a Kubernetes resource
type APIKind string

func (o APIKind) String() string {
	return string(o)
}
