package kubernetes

import (
	"errors"
	"fmt"

	"github.com/go-openapi/spec"
)

// GVKExtension represents the OpenAPI extension x-kubernetes-group-version-kind
type GVKExtension struct {
	Group   APIGroup
	Version APIVersion
	Kind    APIKind
}

// getGVKExtension returns the GVK Kubernetes extension of a definition, if found
func getGVKExtension(extensions spec.Extensions) (*GVKExtension, bool, error) {
	extension, found := extensions["x-kubernetes-group-version-kind"]
	if !found {
		return nil, false, nil
	}

	var gvkMap map[string]interface{}

	gvks, ok := extension.([]interface{})
	if ok {
		if len(gvks) == 0 {
			return nil, false, nil
		}

		if len(gvks) > 1 {
			// TODO DeleteOptions in all groups
			return nil, false, nil
		}
		gvkMap, ok = (gvks[0]).(map[string]interface{})
		if !ok {
			return nil, false, fmt.Errorf("Error getting GVK")
		}
	} else {
		gvkMap, ok = extension.(map[string]interface{})
		if !ok {
			return nil, false, fmt.Errorf("x-kubernetes-group-version-kind is not an array nor a GVK structure")
		}
	}

	group, ok := gvkMap["group"].(string)
	if !ok {
		return nil, false, fmt.Errorf("Error getting GVK apigroup")
	}

	version, ok := gvkMap["version"].(string)
	if !ok {
		return nil, false, fmt.Errorf("Error getting GVK apiversion")
	}

	apiversion, err := NewAPIVersion(version)
	if err != nil {
		return nil, false, fmt.Errorf("Error creating APIVersion")
	}

	kind, ok := gvkMap["kind"].(string)
	if !ok {
		return nil, false, fmt.Errorf("Error getting GVK apikind")
	}
	return &GVKExtension{
		Group:   APIGroup(group),
		Version: *apiversion,
		Kind:    APIKind(kind),
	}, true, nil
}

// GetPatchStrategyExtension returns the PatchStrategy extension of a definition, or nil if not found
func GetPatchStrategyExtension(extensions spec.Extensions) (*string, error) {
	extension, found := extensions["x-kubernetes-patch-strategy"]
	if !found {
		return nil, nil
	}
	value, ok := extension.(string)
	if !ok {
		return nil, errors.New("x-bubernetes-patch-strategy is not a string")
	}
	return &value, nil
}

// GetPatchMergeKeyExtension returns the GetPatchMergeKey extension of a definition, or nil if not found
func GetPatchMergeKeyExtension(extensions spec.Extensions) (*string, error) {
	extension, found := extensions["x-kubernetes-patch-merge-key"]
	if !found {
		return nil, nil
	}
	value, ok := extension.(string)
	if !ok {
		return nil, errors.New("x-bubernetes-patch-merge-key is not a string")
	}
	return &value, nil
}

// GetListType returns the ListType extension of a definition, or nil if not found
func GetListType(definition spec.Schema) (*string, error) {
	extensions := definition.Extensions
	extension, found := extensions["x-kubernetes-list-type"]
	if !found {
		return nil, nil
	}
	value, ok := extension.(string)
	if !ok {
		return nil, errors.New("x-bubernetes-list-type is not a string")
	}
	return &value, nil
}

// GetListMapKeys returns the ListMapKeys extension of a definition, or nil if not found
func GetListMapKeys(definition spec.Schema) ([]string, error) {
	extensions := definition.Extensions
	extension, found := extensions["x-kubernetes-list-map-keys"]
	if !found {
		return nil, nil
	}
	value, ok := extension.([]interface{})
	if !ok {
		return nil, errors.New("x-bubernetes-list-map-keys is not an array")
	}
	var result []string
	for _, val := range value {
		v, ok := val.(string)
		if !ok {
			return nil, errors.New("x-bubernetes-list-map-keys value is not a string")
		}
		result = append(result, v)
	}
	return result, nil
}

// GetActionExtension returns the Action extension of an operation, or nil if not found
func getActionExtension(extensions spec.Extensions) (*ActionExtension, error) {
	extension, found := extensions["x-kubernetes-action"]
	if !found {
		return nil, nil
	}
	value, ok := extension.(string)
	if !ok {
		return nil, errors.New("x-bubernetes-action is not a string")
	}
	action := ActionExtension(value)
	return &action, nil
}
