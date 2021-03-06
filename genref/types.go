package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	texttemplate "text/template"
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting"
	"k8s.io/gengo/types"
)

// apiPackage is a collection of Go packages where API type definitions are found.
type apiPackage struct {
	apiGroup   string
	apiVersion string

	// The Go packages related to this API package. There can be more than one
	// Go package related to the same API package.
	GoPackages []*types.Package

	// List of Types defined. Note that multiple 'types.Package's can define
	// Types for the same apiVersion.
	Types []*apiType

	// Title is set from config
	Title string
}

// DisplayName returns the full name of the API package
func (p *apiPackage) DisplayName() string {
	return fmt.Sprintf("%s/%s", p.apiGroup, p.apiVersion)
}

// GroupName returns the API group the package contains.
func (p *apiPackage) GroupName() string {
	return p.apiGroup
}

// Anchor generates a valid anchor ID for an API package based on its name.
func (p *apiPackage) Anchor() string {
	s := strings.Replace(p.DisplayName(), " ", "", -1)
	s = strings.Replace(s, "/", "-", -1)
	return strings.Replace(s, ".", "-", -1)
}

// VisibleTypes enumerates all visible types contained in a package.
func (p *apiPackage) VisibleTypes() []*apiType {
	var result []*apiType
	for _, t := range sortTypes(p.Types) {
		if !t.isHidden() {
			result = append(result, t)
		}
	}
	return result
}

// GetComment returns the rendered HTML format of the package comment.
func (p *apiPackage) GetComment() template.HTML {
	comments := p.GoPackages[0].DocComments
	return renderComments(comments)
}

// apiMember is a wrapper of types.Member
type apiMember struct {
	types.Member
}

// IsOptional tests if the apiMember is an optional one.
func (m *apiMember) IsOptional() bool {
	tags := types.ExtractCommentTags("+", m.CommentLines)
	_, ok := tags["optional"]
	return ok
}

// FieldName returns the member name when used in serialized format.
func (m *apiMember) FieldName() string {
	v := reflect.StructTag(m.Tags).Get("json")
	v = strings.TrimSuffix(v, ",omitempty")
	v = strings.TrimSuffix(v, ",inline")
	if v != "" {
		return v
	}
	return m.Name
}

// GetType translates the Type field of an apiMember to an apiType reference
func (m *apiMember) GetType() *apiType {
	return &apiType{*m.Type}
}

// Test if a field is an inline one
func (m *apiMember) IsInline() bool {
	return strings.Contains(reflect.StructTag(m.Tags).Get("json"), ",inline")
}

// Test if a member is supposed to be hidden.
func (m *apiMember) Hidden() bool {
	for _, v := range config.HiddenMemberFields {
		if m.Name == v {
			return true
		}
	}
	return false
}

// GetComment returns the rendered HTML output from the field comment.
func (m *apiMember) GetComment() template.HTML {
	return renderComments(m.CommentLines)
}

// apiType is a wrapper of type.Type
type apiType struct {
	types.Type
}

// isLocal tests if the type should be treated as a local definition
func (t *apiType) isLocal() bool {
	t = t.deref()
	if t.Kind == types.Builtin {
		return false
	}
	_, ok := typePkgMap[t.String()]
	return ok
}

// isHidden tests if a type is supposed to be hidden.
func (t *apiType) isHidden() bool {
	for _, pattern := range config.HideTypePatterns {
		if regexp.MustCompile(pattern).MatchString(t.Name.String()) {
			return true
		}
	}
	if !t.IsExported() && unicode.IsLower(rune(t.Name.Name[0])) {
		// types that start with lowercase
		return true
	}
	return false
}

// typeId returns the type Identifier in the format of PackagePath.Name
func (t *apiType) typeId() string {
	t = t.deref()
	return t.Name.String()
}

// deref returns the underlying type when t is a pointer, map, or slice.
func (t *apiType) deref() *apiType {
	if t.Elem != nil {
		return &apiType{*t.Elem}
	}
	return t
}

// GetMembers returns a list of apiMembers each of which is from Type.Members
func (t *apiType) GetMembers() []*apiMember {
	var result []*apiMember
	for _, m := range t.Members {
		member := &apiMember{m}
		result = append(result, member)
	}
	return result
}

// IsExported tests if a type is exported
func (t *apiType) IsExported() bool {
	comments := strings.Join(t.SecondClosestCommentLines, "\n")
	if strings.Contains(comments, "+genclient") {
		return true
	}

	if strings.Contains(comments, "+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object") {
		return true
	}
	return false
}

// Referenced tests if the API type is referenced anywhere in the package
func (t *apiType) Referenced() bool {
	typeName := t.String()
	_, found := references[typeName]
	return found
}

// APIGroup looks up API group for the given type
func (t *apiType) APIGroup() string {
	t = t.deref()

	p := typePkgMap[t.String()]
	if p == nil {
		pwarning("Cannot read apiVersion for %s from type=>pkg map", t.Name.String())
		return "<UNKNOWN_API_GROUP>"
	}

	return p.DisplayName()
}

// Anchor returns the #anchor string for the local type
func (t *apiType) Anchor() string {
	var s string
	group := t.APIGroup()
	if group[0] == '/' {
		s = fmt.Sprintf("%s", t.Name.Name)
	} else {
		s = fmt.Sprintf("%s.%s", group, t.Name.Name)
	}
	s = strings.Replace(s, "/", "-", -1)
	return strings.Replace(s, ".", "-", -1)
}

// Link returns an anchor to the type if it can be generated. returns
// empty string if it is not a local type or unrecognized external type.
func (t *apiType) Link() string {
	t = t.deref() // dereference kind=Pointer

	if t.Kind == types.Builtin {
		return ""
	}

	if t.isLocal() {
		return "#" + t.Anchor()
	}

	var arrIndex = func(a []string, i int) string {
		return a[(len(a)+i)%len(a)]
	}

	// types like k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta,
	// k8s.io/api/core/v1.Container, k8s.io/api/autoscaling/v1.CrossVersionObjectReference,
	// github.com/knative/build/pkg/apis/build/v1alpha1.BuildSpec
	if t.Kind == types.Struct || t.Kind == types.Pointer || t.Kind == types.Interface || t.Kind == types.Alias {
		// gives {{ ImportPath.Identifier }} for type
		id := t.typeId()
		// to parse [meta, v1] from "k8s.io/apimachinery/pkg/apis/meta/v1"
		segments := strings.Split(t.Name.Package, "/")

		for _, v := range config.ExternalPackages {
			r, err := regexp.Compile(v.Match)
			if err != nil {
				perror("Pattern %q failed to compile: %+v", v.Match, err)
				return ""
			}
			// The type identifier is identified as a type from an "external" package
			if r.MatchString(id) {
				tpl, err := texttemplate.New("").Funcs(map[string]interface{}{
					"lower":    strings.ToLower,
					"arrIndex": arrIndex,
				}).Parse(v.Target)
				if err != nil {
					perror("Failed to parse the 'target': %s", v.Target)
					return ""
				}

				var b bytes.Buffer
				err = tpl.Execute(&b, map[string]interface{}{
					"TypeIdentifier":  t.Name.Name,
					"PackagePath":     t.Name.Package,
					"PackageSegments": segments,
				})
				if err != nil {
					perror("Failed to execute template: %+v", err)
					return ""
				}
				return b.String()
			}
		}

		// We are here if the type identifier for the type is not listed as an
		// external one. This means we have to parse it.
		perror("External link source for '%s.%s' is not found.", t.Name.Package, t.Name.Name)
	}
	return ""
}

// DisplayName deterimines how a type is displayed in the docs.
func (t *apiType) DisplayName() string {
	s := t.typeId()
	if t.isLocal() {
		s = t.deref().Name.Name
	}
	if t.Kind == types.Pointer {
		s = strings.TrimLeft(s, "*")
	}

	switch t.Kind {
	case types.Struct,
		types.Interface,
		types.Alias,
		types.Pointer,
		types.Slice,
		types.Builtin:
		// noop
	case types.Map: // return original name
		return t.Name.Name
	default:
		pwarning("Type '%s' has kind='%v' which is unhandled", t.Name, t.Kind)
	}

	// strip prefix if desired
	for _, prefix := range config.StripPrefix {
		if strings.HasPrefix(s, prefix) {
			s = strings.Replace(s, prefix, "", 1)
		}
	}

	if t.Kind == types.Slice {
		s = "[]" + s
	}

	return s
}

// GetComment returns the rendered comment doc for the type.
func (t *apiType) GetComment() template.HTML {
	return renderComments(t.CommentLines)
}

// References returns a list of types where the current type is referenced.
func (t *apiType) References() []*apiType {
	var out []*apiType
	m := make(map[*apiType]struct{})
	for _, ref := range references[t.String()] {
		if !ref.isHidden() {
			m[ref] = struct{}{}
		}
	}
	for k := range m {
		out = append(out, k)
	}
	sortTypes(out)
	return out
}

// groupName extracts the "//+groupName" meta-comment from the specified
// package's comments, or returns empty string if it cannot be found.
func groupName(gopkg *types.Package) string {
	p := gopkg.Constants["GroupName"]
	if p != nil {
		return *p.ConstValue
	}
	m := types.ExtractCommentTags("+", gopkg.Comments)
	v := m["groupName"]
	if len(v) == 1 {
		return v[0]
	}
	return ""
}

// isVendorPackage determines if package is coming from vendor/ dir.
func isVendorPackage(gopkg *types.Package) bool {
	vendorPattern := string(os.PathSeparator) + "vendor" + string(os.PathSeparator)
	return strings.Contains(gopkg.SourcePath, vendorPattern)
}

// sortTypes is a utility function for sorting types in alphabetic order
func sortTypes(typs []*apiType) []*apiType {
	sort.Slice(typs, func(i, j int) bool {
		t1, t2 := typs[i], typs[j]
		if t1.IsExported() && !t2.IsExported() {
			return true
		} else if !t1.IsExported() && t2.IsExported() {
			return false
		}
		return t1.Name.Name < t2.Name.Name
	})
	return typs
}

// renderComments is a utility function for processing a list of strings into
// safe and valid HTML snippets.
func renderComments(comments []string) template.HTML {
	var res string
	// filter out tags in comments
	var list []string
	for _, v := range comments {
		if !strings.HasPrefix(strings.TrimSpace(v), "+") {
			list = append(list, v)
		}
	}
	doc := strings.Join(list, "\n")
	if !config.MarkdownDisabled {
		// This is for blackfriday
		// res = string(blackfriday.Run([]byte(doc)))
		var buf bytes.Buffer
		md := goldmark.New(
			goldmark.WithExtensions(
				highlighting.Highlighting,
			),
		)
		if err := md.Convert([]byte(doc), &buf); err != nil {
			perror("Bad doc detected: %+v", err)
			res = doc
		} else {
			res = buf.String()
		}
		// replace '*' by '&lowast;'
		res = strings.Replace(doc, "*", "&lowast;", -1)
	} else {
		res = strings.Replace(doc, "\n\n", string(template.HTML("<br/><br/>")), -1)
	}
	return template.HTML(res)
}
