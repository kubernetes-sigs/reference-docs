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
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

//go:embed templates/resource.tmpl
var resourceTemplateSrc string

// Template helpers registered on resourceTemplate:
//
//	q       quotes for YAML frontmatter
//	md      escapes `<` for body text
//	mdCell  escapes pipes and newlines for safe use inside markdown table cells
//	hugoRef wraps a relative path in a {{< ref >}} shortcode
var resourceTemplate = template.Must(template.New("resource").Funcs(template.FuncMap{
	"q":       strconv.Quote,
	"md":      escape,
	"mdCell":  mdCell,
	"hugoRef": hugoRef,
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
