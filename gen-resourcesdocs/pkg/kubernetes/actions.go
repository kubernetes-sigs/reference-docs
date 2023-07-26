package kubernetes

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/go-openapi/spec"
)

// ActionExtension represents the OpenAPI extension x-bubernetes-action
type ActionExtension string

const (
	// ActionConnect is the "connect" Kubernetes action
	// DELETE, GET, HEAD, OPTIONS, PATCH, POST, PUT
	ActionConnect = "connect"
	// ActionDelete is the "delete" Kubernetes action
	// DELETE
	ActionDelete = "delete"
	// ActionDeleteCollection is the "deletecollection" Kubernetes action
	// DELETE
	ActionDeleteCollection = "deletecollection"
	// ActionGet is the "get" Kubernetes action
	// GET
	ActionGet = "get"
	// ActionList is the "list" Kubernetes action
	// GET
	ActionList = "list"
	// ActionPatch is the "patch" Kubernetes action
	// PATCH
	ActionPatch = "patch"
	// ActionCreate is the "post" Kubernetes action
	// POST
	ActionCreate = "post"
	// ActionUpdate is the "put" Kubernetes action
	// PUT
	ActionUpdate = "put"
	// ActionWatch is the "watch" Kubernetes action
	// GET
	ActionWatch = "watch"
	// ActionWatchList is the "watchlist" Kubernetes action
	// GET
	ActionWatchList = "watchlist"
)

var (
	// ActionsOrder indicates the natural order of actions
	actionsOrder = map[ActionExtension]int8{
		ActionGet:              0,
		ActionWatch:            1,
		ActionList:             2,
		ActionWatchList:        3,
		ActionCreate:           4,
		ActionUpdate:           5,
		ActionPatch:            6,
		ActionDelete:           7,
		ActionDeleteCollection: 8,
		ActionConnect:          9,
	}

	actionsVerb = map[ActionExtension]string{
		ActionGet:              "get",
		ActionWatch:            "watch",
		ActionList:             "list",
		ActionWatchList:        "watchlist",
		ActionCreate:           "create",
		ActionUpdate:           "update",
		ActionPatch:            "patch",
		ActionDelete:           "delete",
		ActionDeleteCollection: "deletecollection",
		ActionConnect:          "connect",
	}
)

// String returns the string representation of an ActionExtension
func (o ActionExtension) String() string {
	return string(o)
}

// LessThan returns true if o appears before p in natural order
func (o ActionExtension) LessThan(p ActionExtension) bool {
	return actionsOrder[o] < actionsOrder[p]
}

// Verb returns the verb associated with the action
func (o ActionExtension) Verb() string {
	return actionsVerb[o]
}

// ActionPath represents the path of an action
type ActionPath string

func (o ActionPath) String() string {
	return string(o)
}

func (o ActionPath) isNamespaced() bool {
	return strings.Contains(o.String(), "/namespaces/{namespace}")
}

// LessThan returns true if o appears before p in natural order
func (o ActionPath) LessThan(p ActionPath) bool {
	if o.isNamespaced() && !p.isNamespaced() {
		return true
	}
	if !o.isNamespaced() && p.isNamespaced() {
		return false
	}
	return o.String() < p.String()
}

// ActionInfo contains information about a specific endpoint
type ActionInfo struct {
	// Path of the endpoint
	Path ActionPath
	// Kubernetes action mapped to the endpoint
	Action ActionExtension
	// Definition of the action
	Operation *spec.Operation
	// HTTP Method
	HTTPMethod string
	// Parameters of the actions at path level plus operation level
	Parameters ParametersList
}

// ActionInfoList represents a list of actions info
type ActionInfoList []ActionInfo

func (a ActionInfoList) Len() int      { return len(a) }
func (a ActionInfoList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ActionInfoList) Less(i, j int) bool {
	if a[i].Action.LessThan(a[j].Action) {
		return true
	}
	if a[j].Action.LessThan(a[i].Action) {
		return false
	}
	return a[i].Path.LessThan(a[j].Path)
}

// Actions represents a map of ActionInfo, mapped by GVK
type Actions map[string]ActionInfoList

// Add an action to the collection of actions
func (o Actions) Add(specParameters map[string]spec.Parameter, key string, operation *spec.Operation, httpMethod string, pathParameters []spec.Parameter) {

	desc := operation.Description
	if strings.Contains(strings.ToLower(desc), "deprecated") {
		return
	}

	action, err := getActionExtension(operation.Extensions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting extension")
		return
	}
	if action != nil {

		gvk, found, err := getGVKExtension(operation.Extensions)
		if err != nil {
			//fmt.Fprintf(os.Stderr, "Error getting GVK extension for %s|%s: %s\n", key, httpMethod, err)
		} else if !found {
			//fmt.Fprintf(os.Stderr, "GVK extension not found for %s|%s\n", key, httpMethod)
		} else {
			gvkString := gvk.Group.String() + "." + gvk.Version.String() + "." + gvk.Kind.String()

			list := new(ParametersList)
			for _, pathParam := range pathParameters {
				list.Add(specParameters, pathParam)
			}
			for _, opParam := range operation.Parameters {
				list.Add(specParameters, opParam)
			}
			sort.Sort(list)

			newActionInfo := ActionInfo{
				Path:       ActionPath(key),
				Action:     *action,
				Operation:  operation,
				HTTPMethod: httpMethod,
				Parameters: *list,
			}
			if o[gvkString] != nil {
				o[gvkString] = append(o[gvkString], newActionInfo)
			} else {
				o[gvkString] = []ActionInfo{newActionInfo}
			}
		}
	} else {
		//fmt.Fprintf(os.Stderr, "No action for %s|%s\n", key, httpMethod)
	}
}

// Get the actions for a specific GVK
func (o Actions) Get(gvk string) ActionInfoList {
	return o[gvk]
}

// Sort sorts the list of actions for each GVK
func (o Actions) Sort() {
	for k := range o {
		sort.Sort(o[k])
	}
}

func (o Actions) findCommonParameters() {
	for _, actionList := range o {
		for _, action := range actionList {
			ResourcesDescriptions.addActionParameters(&action.Parameters)
		}
	}
	for k, list := range ResourcesDescriptions {
		if len(list) == 1 && list[0].count > 10 {
			ParametersAnnex[k] = struct{}{}
		} else if len(list) == 2 && k == "fieldManager" {
			ParametersAnnex[k] = struct{}{}
			if len(list[0].Description) > len(list[1].Description) {
				list = []descriptionInfo{list[0]}
			} else {
				list = []descriptionInfo{list[1]}
			}
		}
	}
}
