/*
Copyright 2016 The Kubernetes Authors.

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
	"fmt"
	"html"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kubernetes-sigs/reference-docs/gen-apidocs/generators/api"
)

type TOCItem struct {
	Level       int
	Title       string
	Link        string
	File        string
	SubSections []*TOCItem
}

func (ti *TOCItem) ToHTML() string {
	nav := ""
	nav += fmt.Sprintf("<LI class=\"nav-level-%d\" data-level=\"%d\">\n", ti.Level, ti.Level)
	nav += fmt.Sprintf("  <A href=\"#%s\" class=\"nav-item\">%s</A>", ti.Link, ti.Title)

	if len(ti.SubSections) > 0 {
		nav += "\n"
		nav += fmt.Sprintf("  <UL id=\"%s-nav\">\n", ti.Link)

		for _, subItem := range ti.SubSections {
			nav += subItem.ToHTML()
			nav += "\n"
		}

		nav += "  </UL>"
	}

	nav += "\n"
	nav += "</LI>"

	return nav
}

type TOC struct {
	Title     string
	Copyright string
	Sections  []*TOCItem
}

type HTMLWriter struct {
	Config *api.Config
	TOC    TOC

	// currentTOCItem is used to remember the current item between
	// calls to e.g. WriteResourceCategory() followed by WriteResource().
	currentTOCItem *TOCItem
}

func NewHTMLWriter(config *api.Config, copyright, title string) DocWriter {
	writer := HTMLWriter{
		Config: config,
		TOC: TOC{
			Copyright: copyright,
			Title:     title,
			Sections:  []*TOCItem{},
		},
	}
	return &writer
}

func (h *HTMLWriter) Extension() string {
	return ".html"
}

func (h *HTMLWriter) WriteOverview() error {
	filename := "_overview.html"
	if err := writeStaticFile(filename, h.SectionHeading("API Overview")); err != nil {
		return err
	}

	item := TOCItem{
		Level: 1,
		Title: "Overview",
		Link:  "api-overview",
		File:  filename,
	}
	h.TOC.Sections = append(h.TOC.Sections, &item)
	h.currentTOCItem = &item

	return nil
}

func (h *HTMLWriter) WriteAPIGroupVersions(gvs api.GroupVersions) error {
	fn := "_group_versions.html"
	path := filepath.Join(api.IncludesDir, fn)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprint(f, "<DIV id=\"api-groups\">\n")
	fmt.Fprint(f, h.SectionHeading("API Groups")+"\n")
	fmt.Fprint(f, "<P>The API Groups and their versions are summarized in the following table.</P>\n")
	fmt.Fprint(f, "<TABLE class=\"col-md-8\">\n<THEAD><TR><TH>Group</TH><TH>Versions</TH></TR></THEAD>\n<TBODY>\n")

	groups := api.ApiGroups{}
	for group := range gvs {
		groups = append(groups, api.ApiGroup(group))
	}
	sort.Sort(groups)

	for _, group := range groups {
		versionList := gvs[group.String()]
		sort.Sort(versionList)
		var versions []string
		for _, v := range versionList {
			versions = append(versions, v.String())
		}

		fmt.Fprintf(f, "<TR><TD><CODE>%s</CODE></TD><TD><CODE>%s</CODE></TD></TR>\n",
			group, strings.Join(versions, ", "))
	}
	fmt.Fprint(f, "</TBODY>\n</TABLE>\n")
	fmt.Fprint(f, "</DIV>\n")

	item := TOCItem{
		Level: 1,
		Title: "API Groups",
		Link:  "api-groups",
		File:  fn,
	}
	h.TOC.Sections = append(h.TOC.Sections, &item)
	h.currentTOCItem = &item

	return nil
}

func (h *HTMLWriter) WriteResourceCategory(name, file string) error {
	if err := writeStaticFile("_"+file+".html", h.ResourceCategoryHeading(name)); err != nil {
		return err
	}

	link := strings.ReplaceAll(strings.ToLower(name), " ", "-")
	item := TOCItem{
		Level: 1,
		Title: strings.ToUpper(name),
		Link:  link,
		File:  "_" + file + ".html",
	}
	h.TOC.Sections = append(h.TOC.Sections, &item)
	h.currentTOCItem = &item

	return nil
}

func (h *HTMLWriter) ResourceCategoryHeading(title string) string {
	sectionID := strings.ToLower(strings.ReplaceAll(title, " ", "-"))
	return fmt.Sprintf(`<H1 class="toc-item resource-category" id="%s">%s</H1>`, sectionID, title)
}

func (h *HTMLWriter) SectionHeading(title string) string {
	return fmt.Sprintf(`<H1 class="toc-item section">%s</H1>`, title)
}

func (h *HTMLWriter) DefaultStaticContent(title string) string {
	titleID := strings.ToLower(strings.ReplaceAll(title, " ", "-"))
	return fmt.Sprintf("<H1 class=\"strong\" id=\"%s\">%s</H1>\n", titleID, title)
}

func (h *HTMLWriter) writeOtherVersions(w io.Writer, d *api.Definition) {
	if d.OtherVersions.Len() == 0 {
		return
	}

	fmt.Fprint(w, "<DIV class=\"alert alert-success col-md-8\"><I class=\"fa fa-toggle-right\"></I> Other API versions of this object exist:\n")
	for _, v := range d.OtherVersions {
		fmt.Fprintf(w, "%s\n", v.VersionLink())
	}
	fmt.Fprintf(w, "</DIV>\n")
}

func (h *HTMLWriter) writeAppearsIn(w io.Writer, d *api.Definition) {
	if d.AppearsIn.Len() != 0 {
		fmt.Fprintf(w, "<DIV class=\"alert alert-info col-md-8\"><I class=\"fa fa-info-circle\"></I> Appears In:\n <UL>\n")
		for _, a := range d.AppearsIn {
			fmt.Fprintf(w, "  <LI>%s</LI>\n", a.FullHrefLink())
		}
		fmt.Fprintf(w, " </UL>\n</DIV>\n")
	}
}

func (h *HTMLWriter) writeFields(w io.Writer, d *api.Definition) {
	fmt.Fprintf(w, "<TABLE>\n<THEAD><TR><TH>Field</TH><TH>Description</TH></TR></THEAD>\n<TBODY>\n")

	for _, field := range d.Fields {
		fmt.Fprintf(w, "<TR><TD><CODE>%s</CODE>", field.Name)
		if field.Link() != "" {
			fmt.Fprintf(w, "<BR /><I>%s</I>", field.FullLink())
		}
		if field.PatchStrategy != "" {
			fmt.Fprintf(w, "<BR /><B>patch strategy</B>: <I>%s</I>", field.PatchStrategy)
		}
		if field.PatchMergeKey != "" {
			fmt.Fprintf(w, "<BR /><B>patch merge key</B>: <I>%s</I>", field.PatchMergeKey)
		}
		fmt.Fprintf(w, "</TD><TD>%s</TD></TR>\n", field.DescriptionWithEntities)
	}
	fmt.Fprintf(w, "</TBODY>\n</TABLE>\n")
}

func (h *HTMLWriter) WriteDefinitionsOverview() error {
	if err := writeStaticFile("_definitions.html", h.SectionHeading("Definitions")); err != nil {
		return err
	}

	item := TOCItem{
		Level: 1,
		Title: "DEFINITIONS",
		Link:  "definitions",
		File:  "_definitions.html",
	}
	h.TOC.Sections = append(h.TOC.Sections, &item)
	h.currentTOCItem = &item

	return nil
}

func (h *HTMLWriter) WriteDefinition(d *api.Definition) error {
	fn := "_" + definitionFileName(d) + ".html"
	path := filepath.Join(api.IncludesDir, fn)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	nvg := fmt.Sprintf("%s %s %s", d.Name, d.Version, d.GroupDisplayName())
	linkID := getLink(nvg)

	fmt.Fprintf(f, "<DIV class=\"definition-container\" id=\"%s\">\n", linkID)
	defer fmt.Fprint(f, "</DIV>\n")

	fmt.Fprintf(f, "<H2 class=\"definition\">%s</H2>\n", nvg)

	// GVK
	fmt.Fprintf(f, "<TABLE class=\"col-md-8\">\n<THEAD><TR><TH>Group</TH><TH>Version</TH><TH>Kind</TH></TR></THEAD>\n<TBODY>\n")
	fmt.Fprintf(f, "<TR><TD><CODE>%s</CODE></TD><TD><CODE>%s</CODE></TD><TD><CODE>%s</CODE></TD></TR>\n",
		d.GroupDisplayName(), d.Version, d.Name)
	fmt.Fprintf(f, "</TBODY>\n</TABLE>\n")

	fmt.Fprintf(f, "<P>%s</P>\n", d.DescriptionWithEntities)
	h.writeOtherVersions(f, d)
	h.writeAppearsIn(f, d)
	h.writeFields(f, d)

	// Definitions are added to the TOC to enable the generator to later collect
	// all the individual definition files, but definitions will not show up
	// in the nav treet because it would take up too much screen estate.
	item := TOCItem{
		Level: 2,
		Title: nvg,
		Link:  linkID,
		File:  fn,
	}
	h.currentTOCItem.SubSections = append(h.currentTOCItem.SubSections, &item)

	return nil
}

func (h *HTMLWriter) writeSamples(w io.Writer, d *api.Definition) {
	if d.Sample.Sample == "" {
		return
	}

	fmt.Fprintf(w, "<DIV class=\"samples-container\">\n")

	note := d.Sample.Note
	for _, s := range d.GetSamples() {
		sType := strings.Split(s.Tab, ":")[1]
		linkID := sType + "-" + d.LinkID()
		fmt.Fprintf(w, "<BUTTON class=\"btn btn-info\" type=\"button\" data-toggle=\"collapse\"\n")
		fmt.Fprintf(w, "  data-target=\"#%s\" aria-controls=\"%s\"\n", linkID, linkID)
		fmt.Fprintf(w, "  aria-expanded=\"false\">show %s</BUTTON>\n", sType)
	}

	for _, s := range d.GetSamples() {
		sType := strings.Split(s.Tab, ":")[1]
		linkID := sType + "-" + d.LinkID()
		lType := strings.Split(s.Type, ":")[1]
		lang := strings.Split(lType, "_")[1]
		fmt.Fprintf(w, "<DIV class=\"collapse\" id=\"%s\">\n", linkID)
		fmt.Fprintf(w, "  <DIV class=\"panel panel-default\">\n<DIV class=\"panel-heading\">%s</DIV>\n", note)
		fmt.Fprintf(w, "  <DIV class=\"panel-body\">\n<PRE class=\"%s\">", sType)
		fmt.Fprintf(w, "<CODE class=\"lang-%s\">\n", lang)
		// TODO: Add language highlight
		fmt.Fprintf(w, "%s\n</CODE></PRE></DIV></DIV></DIV>\n", html.EscapeString(s.Text))
	}

	fmt.Fprint(w, "</DIV>\n")
}

func (h *HTMLWriter) writeOperationSample(w io.Writer, req bool, op string, examples []api.ExampleText) {
	// e.Tab bdocs-tab:kubectl  | bdocs-tab:curl
	// e.Msg `kubectl` Command  | Output | Response Body | `curl` Command (*requires `kubectl proxy` to be running*)
	// e.Type bdocs-tab:kubectl_shell
	// e.Text <actual command>

	for _, e := range examples {
		eType := strings.Split(e.Tab, ":")[1]
		var sampleID string
		var btnText string
		if req {
			sampleID = "req-" + eType + "-" + op
			btnText = eType + " request"
		} else {
			sampleID = "res-" + eType + "-" + op
			btnText = eType + " response"
		}
		fmt.Fprintf(w, "<BUTTON class=\"btn btn-info\" type=\"button\" data-toggle=\"collapse\"\n")
		fmt.Fprintf(w, "  data-target=\"#%s\" aria-controls=\"%s\"\n", sampleID, sampleID)
		fmt.Fprintf(w, "  aria-expanded=\"false\">%s example</BUTTON>\n", btnText)
	}

	for _, e := range examples {
		eType := strings.Split(e.Tab, ":")[1]
		var sampleID string
		if req {
			sampleID = "req-" + eType + "-" + op
		} else {
			sampleID = "res-" + eType + "-" + op
		}
		msg := e.Msg
		if eType == "curl" && strings.Contains(msg, "proxy") {
			msg = "<CODE>curl</CODE> command (<I>requires <code>kubectl proxy</code> to be running</I>)"
		} else if eType == "kubectl" && strings.Contains(msg, "Command") { // `kubectl` command
			msg = "<CODE>kubectl</CODE> command"
		}
		lType := strings.Split(e.Type, ":")[1]
		lang := strings.Split(lType, "_")[1]
		fmt.Fprintf(w, "<DIV class=\"collapse\" id=\"%s\">\n", sampleID)
		fmt.Fprintf(w, "  <DIV class=\"panel panel-default\">\n<DIV class=\"panel-heading\">%s</DIV>\n", msg)
		fmt.Fprintf(w, "  <DIV class=\"panel-body\">\n<PRE class=\"%s\">", eType)
		fmt.Fprintf(w, "<CODE class=\"lang-%s\">\n", lang)
		// TODO: Add language highlight
		fmt.Fprintf(w, "%s\n</CODE></PRE></DIV></DIV></DIV>\n", e.Text)
	}
}

func (h *HTMLWriter) writeParams(w io.Writer, title string, params api.Fields) {
	fmt.Fprintf(w, "<H3>%s</H3>\n", title)
	fmt.Fprintf(w, "<TABLE>\n<THEAD><TR><TH>Parameter</TH><TH>Description</TH></TR></THEAD>\n<TBODY>\n")
	for _, p := range params {
		fmt.Fprintf(w, "<TR><TD><CODE>%s</CODE>", p.Name)
		if p.Link() != "" {
			fmt.Fprintf(w, "<br /><I>%s</I>", p.FullLink())
		}
		fmt.Fprintf(w, "</TD><TD>%s</TD></TR>\n", p.Description)
	}
	fmt.Fprintf(w, "</TBODY>\n</TABLE>\n")
}

func (h *HTMLWriter) writeRequestParams(w io.Writer, o *api.Operation) {
	// Operation path params
	if o.PathParams.Len() > 0 {
		h.writeParams(w, "Path Parameters", o.PathParams)
	}

	// operation query params
	if o.QueryParams.Len() > 0 {
		h.writeParams(w, "Query Parameters", o.QueryParams)
	}

	// operation body params
	if o.BodyParams.Len() > 0 {
		h.writeParams(w, "Body Parameters", o.BodyParams)
	}
}

func (h *HTMLWriter) writeResponseParams(w io.Writer, o *api.Operation) {
	if o.HttpResponses.Len() == 0 {
		return
	}

	fmt.Fprintf(w, "<H3>Response</H3>\n")
	fmt.Fprintf(w, "<TABLE>\n<THEAD><TR><TH>Code</TH><TH>Description</TH></TR></THEAD>\n<TBODY>\n")
	responses := o.HttpResponses
	sort.Slice(responses, func(i, j int) bool {
		return strings.Compare(responses[i].Name, responses[j].Name) < 0
	})
	for _, p := range responses {
		fmt.Fprintf(w, "<TR><TD>%s", p.Name)
		if p.Field.Link() != "" {
			fmt.Fprintf(w, "<br /><I>%s</I>", p.Field.FullLink())
		}
		fmt.Fprintf(w, "</TD><TD>%s</TD></TR>\n", p.Field.Description)
	}
	fmt.Fprintf(w, "</TBODY>\n</TABLE>\n")
}

func (h *HTMLWriter) WriteResource(r *api.Resource) error {
	fn := "_" + conceptFileName(r.Definition) + ".html"
	path := filepath.Join(api.IncludesDir, fn)

	w, err := os.Create(path)
	if err != nil {
		return err
	}
	defer w.Close()

	dvg := fmt.Sprintf("%s %s %s", r.Name, r.Definition.Version, r.Definition.GroupDisplayName())
	linkID := getLink(dvg)

	fmt.Fprintf(w, "<DIV class=\"resource-container\" id=\"%s\">\n", linkID)
	defer fmt.Fprint(w, "</DIV>\n")

	fmt.Fprintf(w, "<H1 class=\"toc-item resource\">%s</H1>\n", dvg)

	h.writeSamples(w, r.Definition)

	// GVK
	fmt.Fprintf(w, "<TABLE class=\"col-md-8\">\n<THEAD><TR><TH>Group</TH><TH>Version</TH><TH>Kind</TH></TR></THEAD>\n<TBODY>\n")
	fmt.Fprintf(w, "<TR><TD><CODE>%s</CODE></TD><TD><CODE>%s</CODE></TD><TD><CODE>%s</CODE></TD></TR>\n",
		r.Definition.GroupDisplayName(), r.Definition.Version, r.Name)
	fmt.Fprintf(w, "</TBODY>\n</TABLE>\n")

	if r.DescriptionWarning != "" {
		fmt.Fprintf(w, "<DIV class=\"alert alert-warning col-md-8\"><P><I class=\"fa fa-warning\"></I> <B>Warning:</B></P><P>%s</P></DIV>\n", r.DescriptionWarning)
	}
	if r.DescriptionNote != "" {
		fmt.Fprintf(w, "<DIV class=\"alert alert-info col-md-8\"><I class=\"fa fa-bullhorn\"></I> %s</DIV>\n", r.DescriptionNote)
	}

	h.writeOtherVersions(w, r.Definition)
	h.writeAppearsIn(w, r.Definition)
	h.writeFields(w, r.Definition)

	// Inline
	if r.Definition.Inline.Len() > 0 {
		fmt.Fprintf(w, "<DIV class=\"inline-definitions-container\">\n")
		for _, d := range r.Definition.Inline {
			fmt.Fprintf(w, "<H3 class=\"inline-definition\" id=\"%s\">%s %s %s</H3>\n", d.LinkID(), d.Name, d.Version, d.Group)
			h.writeAppearsIn(w, d)
			h.writeFields(w, d)
		}
		fmt.Fprint(w, "</DIV>\n")
	}

	resourceItem := TOCItem{
		Level: 2,
		Title: dvg,
		Link:  linkID,
		File:  fn,
	}
	h.currentTOCItem.SubSections = append(h.currentTOCItem.SubSections, &resourceItem)

	// Operations
	if len(r.Definition.OperationCategories) == 0 {
		return nil
	}

	for _, c := range r.Definition.OperationCategories {
		if len(c.Operations) == 0 {
			continue
		}

		catID := strings.ReplaceAll(strings.ToLower(c.Name), " ", "-") + "-" + r.Definition.LinkID()
		fmt.Fprintf(w, "<DIV class=\"operation-category-container\" id=\"%s\">\n", catID)
		fmt.Fprintf(w, "<H2 class=\"toc-item operation-category\">%s</H2>\n", c.Name)

		ocItem := TOCItem{
			Level: 3,
			Title: c.Name,
			Link:  catID,
		}
		resourceItem.SubSections = append(resourceItem.SubSections, &ocItem)

		for _, o := range c.Operations {
			opID := strings.ReplaceAll(strings.ToLower(o.Type.Name), " ", "-") + "-" + r.Definition.LinkID()
			fmt.Fprintf(w, "<DIV class=\"operation-container\" id=\"%s\">\n", opID)
			fmt.Fprintf(w, "<H2 class=\"toc-item operation\">%s</H2>\n", o.Type.Name)

			OPItem := TOCItem{
				Level: 4,
				Title: o.Type.Name,
				Link:  opID,
			}
			ocItem.SubSections = append(ocItem.SubSections, &OPItem)

			// Example requests
			requests := o.GetExampleRequests()
			if len(requests) > 0 {
				h.writeOperationSample(w, true, opID, requests)
			}
			// Example responses
			responses := o.GetExampleResponses()
			if len(responses) > 0 {
				h.writeOperationSample(w, false, opID, responses)
			}

			fmt.Fprintf(w, "<P>%s</P>\n", o.Description())
			fmt.Fprintf(w, "<H3>HTTP Request</H3>\n")
			fmt.Fprintf(w, "<CODE>%s</CODE>\n", o.GetDisplayHttp())

			h.writeRequestParams(w, o)
			h.writeResponseParams(w, o)

			fmt.Fprint(w, "</DIV>\n")
		}

		fmt.Fprint(w, "</DIV>\n")
	}

	return nil
}

func (h *HTMLWriter) WriteOldVersionsOverview() error {
	if err := writeStaticFile("_oldversions.html", h.SectionHeading("Old API Versions")); err != nil {
		return err
	}

	item := TOCItem{
		Level: 1,
		Title: "OLD API VERSIONS",
		Link:  "old-api-versions",
		File:  "_oldversions.html",
	}
	h.TOC.Sections = append(h.TOC.Sections, &item)
	h.currentTOCItem = &item

	return nil
}

func (h *HTMLWriter) generateNavContent() string {
	nav := "<UL id=\"navigation\">\n"

	for _, sec := range h.TOC.Sections {
		nav += sec.ToHTML()
		nav += "\n"
	}

	nav += "</UL>\n"

	return nav
}

func (h *HTMLWriter) generateIndex(navContent string) error {
	html, err := os.Create(filepath.Join(api.BuildDir, "index.html"))
	if err != nil {
		return err
	}
	defer html.Close()

	/* Make sure the following stylesheets exist in kubernetes/website repo:
	   kubernetes/website/static/css/bootstrap-4.3.1.min.css
	   kubernetes/website/static/css/fontawesome-4.7.0.min.css
	   kubernetes/website/static/css/style_apiref.css
	*/
	fmt.Fprintf(html, "<!DOCTYPE html>\n<HTML lang=\"en\">\n<HEAD>\n<META charset=\"UTF-8\">\n")
	fmt.Fprintf(html, "<TITLE>%s</TITLE>\n", h.TOC.Title)
	fmt.Fprintf(html, "<LINK rel=\"shortcut icon\" href=\"favicon.ico\" type=\"image/vnd.microsoft.icon\">\n")
	fmt.Fprintf(html, "<LINK rel=\"stylesheet\" href=\"/css/bootstrap-4.3.1.min.css\" type=\"text/css\">\n")
	fmt.Fprintf(html, "<LINK rel=\"stylesheet\" href=\"/css/fontawesome-4.7.0.min.css\" type=\"text/css\">\n")
	fmt.Fprintf(html, "<LINK rel=\"stylesheet\" href=\"/css/style_apiref.css\" type=\"text/css\">\n")
	fmt.Fprintf(html, "</HEAD>\n<BODY>\n")
	fmt.Fprintf(html, "<DIV id=\"wrapper\" class=\"container-fluid\">\n")
	fmt.Fprintf(html, "<DIV class=\"row\">\n")
	fmt.Fprintf(html, "<DIV id=\"sidebar-wrapper\" class=\"col-xs-4 col-sm-3 col-md-2 side-nav side-bar-nav\">\n")

	// html buffer
	buf := "<DIV class=\"row\">\n  <DIV class=\"col-md-6 copyright\">\n " + h.TOC.Copyright + "\n  </DIV>\n"
	buf += "  <DIV class=\"col-md-6 text-right\">\n"
	buf += fmt.Sprintf("    <DIV>Generated at: %s</DIV>\n", time.Now().Format("2006-01-02 15:04:05 (MST)"))
	pos := strings.LastIndex(h.Config.SpecVersion, ".")
	release := fmt.Sprintf("release-%s", h.Config.SpecVersion[1:pos])
	spec_link := "https://github.com/kubernetes/kubernetes/blob/" + release + "/api/openapi-spec/swagger.json"
	buf += fmt.Sprintf("    <DIV>API Version: <a href=\"%s\">%s</a></DIV>\n", spec_link, h.Config.SpecVersion)
	buf += "  </DIV>\n</DIV>"
	const OK = "\033[32mOK\033[0m"
	const NOT_FOUND = "\033[31mNot found\033[0m"
	for _, sec := range h.TOC.Sections {
		fmt.Printf("Collecting %s ... ", sec.File)
		content, err := os.ReadFile(filepath.Join(api.IncludesDir, sec.File))
		if err == nil {
			buf += string(content)
			fmt.Println(OK)
		} else {
			fmt.Println(NOT_FOUND)
		}

		for _, sub := range sec.SubSections {
			if len(sub.File) > 0 {
				subdata, err := os.ReadFile(filepath.Join(api.IncludesDir, sub.File))
				fmt.Printf("Collecting %s ... ", sub.File)
				if err == nil {
					buf += string(subdata)
					fmt.Println(OK)
				} else {
					fmt.Println(NOT_FOUND)
				}
			}

			for _, subsub := range sub.SubSections {
				if len(subsub.File) > 0 {
					subsubdata, err := os.ReadFile(filepath.Join(api.IncludesDir, subsub.File))
					fmt.Printf("Collecting %s ...", subsub.File)
					if err == nil {
						buf += string(subsubdata)
						fmt.Println(OK)
					} else {
						fmt.Println(NOT_FOUND)
					}
				}
			}
		}
	}

	/*
		Make sure the following scripts exist in kubernetes/website repo:
		kubernetes/website/static/js/jquery-3.6.0.min.js
		kubernetes/website/static/js/jquery.scrollTo-2.1.3.min.js
		kubernetes/website/static/js/bootstrap-4.6.1.min.js
		kubernetes/website/static/js/scroll-apiref.js
	*/
	fmt.Fprintf(html, "%s</DIV>\n", navContent)
	fmt.Fprintf(html, "<DIV id=\"page-content-wrapper\" class=\"col-xs-8 offset-xs-4 col-sm-9 offset-sm-3 col-md-10 offset-md-2 body-content\">\n")
	fmt.Fprintf(html, "%s", string(buf))
	fmt.Fprintf(html, "\n</DIV>\n</DIV>\n</DIV>\n")
	fmt.Fprintf(html, "<SCRIPT src=\"/js/jquery-3.6.0.min.js\"></SCRIPT>\n")
	fmt.Fprintf(html, "<SCRIPT src=\"/js/jquery.scrollTo-2.1.3.min.js\"></SCRIPT>\n")
	fmt.Fprintf(html, "<SCRIPT src=\"/js/bootstrap-4.6.1.min.js\"></SCRIPT>\n")
	fmt.Fprintf(html, "<SCRIPT src=\"/js/scroll-apiref.js\"></SCRIPT>\n")
	fmt.Fprintf(html, "</BODY>\n</HTML>\n")

	return nil
}

func (h *HTMLWriter) Finalize() error {
	if err := os.MkdirAll(api.BuildDir, os.ModePerm); err != nil {
		return err
	}

	navContent := h.generateNavContent()

	if err := h.generateIndex(navContent); err != nil {
		return err
	}

	return nil
}
