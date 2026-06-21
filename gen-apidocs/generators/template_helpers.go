/*
Copyright 2026 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package generators

import (
	_ "embed"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

//go:embed templates/resource.tmpl
var resourceTemplateSrc string

// Template helpers registered on resourceTemplate:
//
//	q             quotes for YAML frontmatter
//	md            escapes `<` for body text
//	mdCell        escapes pipes and newlines for safe use inside markdown table cells
//	hugoRef       wraps a relative path in a {{< ref >}} shortcode
//	hugoShortcode emits a Hugo shortcode with positional string args
//	dict          builds a map for passing named args to named templates
var resourceTemplate = template.Must(template.New("resource").Funcs(template.FuncMap{
	"q":             strconv.Quote,
	"md":            escape,
	"mdCell":        mdCell,
	"hugoRef":       hugoRef,
	"hugoShortcode": hugoShortcode,
	"dict":          dict,
}).Parse(resourceTemplateSrc))

var (
	enumHeaderRegex = regexp.MustCompile(`\s+Possible enum values:`)
	enumBulletRegex = regexp.MustCompile(`\s+- ` + "`")
)

// escape covers the only markdown-breaking character in OpenAPI descriptions:
// raw `<` that would otherwise be read as HTML.
func escape(s string) string {
	s = strings.ReplaceAll(s, "<", `\<`)
	s = enumHeaderRegex.ReplaceAllString(s, "<br/><br/>Possible enum values:")
	s = enumBulletRegex.ReplaceAllString(s, "<br/> - `")
	return s
}

// mdCell escapes a description for safe use inside a markdown table cell:
// pipes break the row, and raw newlines split the cell across rows.
func mdCell(s string) string {
	s = escape(s)
	s = strings.ReplaceAll(s, "|", `\|`)
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}

// hugoRef wraps a path in a {{< ref >}} shortcode resolved by Hugo at build time.
func hugoRef(path string) string {
	return `{{< ref "` + path + `" >}}`
}

// hugoShortcode emits a Hugo shortcode invocation: `{{< name "arg1" "arg2" >}}`.
// String args are quoted to preserve embedded whitespace.
func hugoShortcode(name string, args ...string) string {
	if len(args) == 0 {
		return `{{< ` + name + ` >}}`
	}
	quoted := make([]string, len(args))
	for i, a := range args {
		quoted[i] = strconv.Quote(a)
	}
	return `{{< ` + name + ` ` + strings.Join(quoted, " ") + ` >}}`
}

// dict builds a map[string]any from alternating key/value pairs, used to pass
// named arguments to {{template "name" ...}} calls.
func dict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("dict: odd argument count %d", len(values))
	}
	m := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		k, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict: non-string key at index %d (%T)", i, values[i])
		}
		m[k] = values[i+1]
	}
	return m, nil
}
