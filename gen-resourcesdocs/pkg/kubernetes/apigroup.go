package kubernetes

// APIGroup represents the group of a Kubernetes API
type APIGroup string

func (o APIGroup) String() string {
	return string(o)
}

// Replaces returns true if 'o' group is replaced by 'p' group
func (o APIGroup) Replaces(p APIGroup) bool {
	// * replaces extensions
	if o.String() != "extensions" && p.String() == "extensions" {
		return true
	}

	// events replaces core
	if o.String() == "events.k8s.io" && p.String() == "" {
		return true
	}

	// autoscaling replaces apps
	if o.String() == "autoscaling" && p.String() == "apps" {
		return true
	}

	return false
}
