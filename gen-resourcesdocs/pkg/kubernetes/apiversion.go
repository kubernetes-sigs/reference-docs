package kubernetes

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
)

// VersionStage represents the stage of the version: alpha, beta, ga
type VersionStage int

const (
	// StageAlpha is the first stage during development
	StageAlpha VersionStage = iota
	// StageBeta is the next stage
	StageBeta
	// StageGA is the latest stage
	StageGA
)

var stageName = map[VersionStage]string{
	StageAlpha: "alpha",
	StageBeta:  "beta",
	StageGA:    "",
}

// APIVersion represents the version of a Kubernetes API (v1aplha1, v1beta2, v1, v2, etc)
type APIVersion struct {
	// Version is the *1* in v1alpha2
	Version int
	// Stage is *alpha* in v1alpha2
	Stage VersionStage
	// StageVersion is *2* in v1alpha2, or nil in v1
	StageVersion *int
}

// NewAPIVersion creates a new APIVersion struct from the literal version (for example v1alpha1)
func NewAPIVersion(literal string) (apiversion *APIVersion, err error) {
	re := regexp.MustCompile("^v(\\d+)((alpha|beta|)(\\d))?$")
	parts := re.FindStringSubmatch(literal)
	if parts == nil || len(parts) != 5 {
		return nil, fmt.Errorf("Error parsing %s", literal)
	}

	apiversion = &APIVersion{}
	apiversion.Version, err = strconv.Atoi(parts[1])
	if err != nil {
		log.Printf("Error parsing integer '%s'", parts[1])
		return nil, err
	}
	apiversion.Stage, err = getVersionStage(parts[3])
	if err != nil {
		return nil, err
	}
	if apiversion.Stage != StageGA {
		var res int
		res, err = strconv.Atoi(parts[4])
		if err != nil {
			return nil, err
		}
		apiversion.StageVersion = &res
	}
	return
}

func getVersionStage(stage string) (VersionStage, error) {
	switch stage {
	case stageName[StageAlpha]:
		return StageAlpha, nil
	case stageName[StageBeta]:
		return StageBeta, nil
	case stageName[StageGA]:
		return StageGA, nil
	default:
		return 0, fmt.Errorf("unknown stage %s", stage)
	}
}

// String returns the literal representation of an APIVersion
func (o *APIVersion) String() string {
	if o == nil {
		return ""
	}
	if o.Stage == StageGA {
		return fmt.Sprintf("v%d", o.Version)
	}
	return fmt.Sprintf("v%d%s%d", o.Version, stageName[o.Stage], *o.StageVersion)
}

// Equals returns true if 'o' and 'p' represent the same version
func (o *APIVersion) Equals(p *APIVersion) bool {
	return o.String() == p.String()
}

// LessThan returns true if 'o' version comes before 'p' version
func (o *APIVersion) LessThan(p *APIVersion) bool {
	if o.Version != p.Version {
		return o.Version < p.Version
	}
	if o.Stage != StageGA {
		if o.Stage != p.Stage {
			return o.Stage < p.Stage
		}
		return *o.StageVersion < *p.StageVersion
	}
	return false
}

// Replaces returns true if 'o' version replaces 'p' version
func (o *APIVersion) Replaces(p *APIVersion) bool {
	return o.Version == p.Version && p.LessThan(o)
}

// UnmarshalYAML helps unmarshal APIVersion values from YAML
func (o *APIVersion) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	err := unmarshal(&str)
	if err != nil {
		return err
	}
	v, err := NewAPIVersion(str)
	if err != nil {
		return err
	}
	*o = *v
	return nil
}
