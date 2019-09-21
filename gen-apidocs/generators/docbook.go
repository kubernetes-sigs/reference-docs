/*
Copyright 2019 The Kubernetes Authors.

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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kubernetes-incubator/reference-docs/gen-apidocs/generators/api"
)

type DocbookTOCItem struct {
	Level       int
	Title       string
	Link        string
	File        string
	FileClose   string
	SubSections []*DocbookTOCItem
}

type DocbookTOC struct {
	Title     string
	Copyright string
	Sections  []*DocbookTOCItem
}

type DocbookWriter struct {
	Config         *api.Config
	TOC            DocbookTOC
	CurrentSection *DocbookTOCItem
}

func NewDocbookWriter(config *api.Config, copyright, title string) DocWriter {
	writer := DocbookWriter{
		Config: config,
		TOC: DocbookTOC{
			Copyright: copyright,
			Title:     title,
			Sections:  []*DocbookTOCItem{},
		},
	}
	return &writer
}

func (h *DocbookWriter) Extension() string {
	return ".xml"
}

func (h *DocbookWriter) DefaultStaticContent(title string) string {
	titleID := strings.ToLower(strings.Replace(title, " ", "-", -1))
	return fmt.Sprintf("<part id=\"strong-%s-strong\"><title>%s</title>\n", titleID, title)
}

func (h *DocbookWriter) WriteOverview() {
	fn := "_overview.xml"
	writeStaticFile("Overview", fn, h.DefaultStaticContent("Overview"))
	item := DocbookTOCItem{
		Level: 1,
		Title: "Overview",
		Link:  "strong-api-overview-strong",
		File:  fn,
	}
	h.TOC.Sections = append(h.TOC.Sections, &item)
	h.CurrentSection = &item
}

func (h *DocbookWriter) WriteResourceCategory(name, file string) {
	writeStaticFile(name, "_"+file+".xml", h.DefaultStaticContent(name))
	link := strings.Replace(strings.ToLower(name), " ", "-", -1)
	item := DocbookTOCItem{
		Level:     1,
		Title:     strings.ToUpper(name),
		Link:      "strong-" + link + "-strong",
		File:      "_" + file + ".xml",
		FileClose: "_" + file + "_close.xml",
	}
	h.TOC.Sections = append(h.TOC.Sections, &item)
	h.CurrentSection = &item

	writeStaticFile(name, "_"+file+"_close.xml", h.staticContentClose())
}

func (h *DocbookWriter) WriteResource(r *api.Resource) {
	fn := "_" + conceptFileName(r.Definition) + ".xml"
	fnClose := "_" + conceptFileName(r.Definition) + "_close.xml"

	pathClose := *api.ConfigDir + "/includes/" + fnClose
	wClose, err := os.Create(pathClose)
	defer wClose.Close()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%v", err))
		os.Exit(1)
	}
	fmt.Fprintf(wClose, "</chapter>\n")

	path := *api.ConfigDir + "/includes/" + fn
	w, err := os.Create(path)
	defer w.Close()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%v", err))
		os.Exit(1)
	}

	dvg := fmt.Sprintf("%s %s %s", r.Name, r.Definition.Version, r.Definition.GroupDisplayName())
	linkID := getLink(dvg)
	fmt.Fprintf(w, "<chapter id=\"%s\"><title>%s</title>\n", linkID, dvg)
	h.writeSample(w, r.Definition)

	// GVK
	fmt.Fprintf(w, "<informaltable>\n<tgroup cols=\"3\"><thead><row><entry>Group</entry><entry>Version</entry><entry>Kind</entry></row></thead>\n<tbody>\n")
	fmt.Fprintf(w, "<row><entry><systemitem>%s</systemitem></entry><entry><systemitem>%s</systemitem></entry><entry><systemitem>%s</systemitem></entry></row>\n",
		r.Definition.GroupDisplayName(), r.Definition.Version.String(), r.Name)
	fmt.Fprintf(w, "</tbody></tgroup>\n</informaltable>\n")

	if r.DescriptionWarning != "" {
		fmt.Fprintf(w, "<warning><para>%s</para></warning>\n", a2link(r.DescriptionWarning))
	}
	if r.DescriptionNote != "" {
		fmt.Fprintf(w, "<note><para>%s</para></note>\n", a2link(r.DescriptionNote))
	}

	h.writeOtherVersions(w, r.Definition)
	h.writeAppearsIn(w, r.Definition)
	h.writeFields(w, r.Definition)

	// Inline
	if r.Definition.Inline.Len() > 0 {
		for _, d := range r.Definition.Inline {
			fmt.Fprintf(w, "<sect1 id=\"%s\"><title>%s %s %s</title>\n", d.LinkID(), d.Name, d.Version, d.Group)
			h.writeAppearsIn(w, d)
			h.writeFields(w, d)
			fmt.Fprintf(w, "</sect1>")
		}
	}

	item := DocbookTOCItem{
		Level:     1,
		Title:     dvg,
		Link:      linkID,
		File:      fn,
		FileClose: fnClose,
	}
	h.CurrentSection.SubSections = append(h.CurrentSection.SubSections, &item)

	// Operations
	if len(r.Definition.OperationCategories) == 0 {
		return
	}

	for i, c := range r.Definition.OperationCategories {
		if len(c.Operations) == 0 {
			continue
		}
		catID := strings.Replace(strings.ToLower(c.Name), " ", "-", -1) + "-" + r.Definition.LinkID()
		catID = "strong-" + catID + "-strong"
		if i > 0 {
			fmt.Fprintf(w, "</sect1>\n")
		}
		fmt.Fprintf(w, "<sect1 id=\"%s\"><title>%s</title>\n", catID, c.Name)
		OCItem := DocbookTOCItem{
			Level: 2,
			Title: c.Name,
			Link:  catID,
		}
		h.CurrentSection.SubSections = append(h.CurrentSection.SubSections, &OCItem)

		for j, o := range c.Operations {
			opID := strings.Replace(strings.ToLower(o.Type.Name), " ", "-", -1) + "-" + r.Definition.LinkID()
			if j > 0 {
				fmt.Fprintf(w, "</simplesect>\n")
			}
			fmt.Fprintf(w, "<simplesect id=\"%s\"><title>%s</title>\n", opID, o.Type.Name)
			OPItem := DocbookTOCItem{
				Level: 2,
				Title: o.Type.Name,
				Link:  opID,
			}
			OCItem.SubSections = append(OCItem.SubSections, &OPItem)

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

			fmt.Fprintf(w, "<bridgehead renderas=\"sect4\">HTTP Request</bridgehead>\n")
			fmt.Fprintf(w, "<para>%s</para>\n", o.Description())
			fmt.Fprintf(w, "<synopsis>%s</synopsis>\n", o.GetDisplayHttp())

			h.writeRequestParams(w, o)
			h.writeResponseParams(w, o)
		}
		fmt.Fprintf(w, "</simplesect>\n")
	}
	fmt.Fprintf(w, "</sect1>\n")
}

func (h *DocbookWriter) WriteDefinitionsOverview() {
	writeStaticFile("Definitions", "_definitions.xml", h.DefaultStaticContent("Definitions"))
	item := DocbookTOCItem{
		Level:     1,
		Title:     "DEFINITIONS",
		Link:      "strong-definitions-strong",
		File:      "_definitions.xml",
		FileClose: "_definitions_close.xml",
	}
	h.TOC.Sections = append(h.TOC.Sections, &item)
	h.CurrentSection = &item

	writeStaticFile("Definitions", "_definitions_close.xml", h.staticContentClose())
}

func (h *DocbookWriter) WriteDefinition(d *api.Definition) {
	fn := "_" + definitionFileName(d) + ".xml"
	path := *api.ConfigDir + "/includes/" + fn
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%v", err))
		os.Exit(1)
	}
	nvg := fmt.Sprintf("%s %s %s", d.Name, d.Version, d.GroupDisplayName())
	linkID := getLink(nvg)
	fmt.Fprintf(f, "<sect1 id=\"%s\"><title>%s</title>\n", linkID, nvg)
	fmt.Fprintf(f, "<informaltable>\n<tgroup cols=\"3\"><thead><row><entry>Group</entry><entry>Version</entry><entry>Kind</entry></row></thead>\n<tbody>\n")
	fmt.Fprintf(f, "<row><entry><systemitem>%s</systemitem></entry><entry><systemitem>%s</systemitem></entry><entry><systemitem>%s</systemitem></entry></row>\n",
		d.GroupDisplayName(), d.Version, d.Name)
	fmt.Fprintf(f, "</tbody></tgroup>\n</informaltable>\n")

	fmt.Fprintf(f, "<para>%s</para>\n", d.DescriptionWithEntities)
	h.writeOtherVersions(f, d)
	h.writeAppearsIn(f, d)
	h.writeFields(f, d)

	fmt.Fprintf(f, "</sect1>\n")
	item := DocbookTOCItem{
		Level: 2,
		Title: nvg,
		Link:  linkID,
		File:  fn,
	}
	h.CurrentSection.SubSections = append(h.CurrentSection.SubSections, &item)
}

func (h *DocbookWriter) WriteOldVersionsOverview() {
	writeStaticFile("Old Versions", "_oldversions.xml", h.DefaultStaticContent("Old Versions"))
	item := DocbookTOCItem{
		Level:     1,
		Title:     "OLD API VERSIONS",
		Link:      "strong-old-api-versions-strong",
		File:      "_oldversions.xml",
		FileClose: "_oldversions_close.xml",
	}
	h.TOC.Sections = append(h.TOC.Sections, &item)
	h.CurrentSection = &item
	writeStaticFile("Old Versions", "_oldversions_close.xml", h.staticContentClose())
}

func (h *DocbookWriter) staticContentClose() string {
	return fmt.Sprintf("</part>\n")
}

func (h *DocbookWriter) writeOtherVersions(w io.Writer, d *api.Definition) {
	if d.OtherVersions.Len() == 0 {
		return
	}

	fmt.Fprint(w, "<caution><para>Other API versions of this object exist: ")
	for _, v := range d.OtherVersions {
		fmt.Fprintf(w, "%s\n", a2link(v.VersionLink()))
	}
	fmt.Fprintf(w, "</para></caution>\n")
}

func (h *DocbookWriter) writeAppearsIn(w io.Writer, d *api.Definition) {
	if d.AppearsIn.Len() != 0 {
		fmt.Fprintf(w, "<note><para> Appears In: <itemizedlist>\n")
		for _, a := range d.AppearsIn {
			fmt.Fprintf(w, "  <listitem><para>%s</para></listitem>\n", a2link(a.FullHrefLink()))
		}
		fmt.Fprintf(w, " </itemizedlist>\n</para></note>\n")
	}
}

func (h *DocbookWriter) writeFields(w io.Writer, d *api.Definition) {
	if len(d.Fields) == 0 {
		return
	}
	fmt.Fprintf(w, "<variablelist>\n")

	for _, field := range d.Fields {
		fmt.Fprintf(w, "<varlistentry><term>%s", field.Name)

		if field.Link() != "" {
			fmt.Fprintf(w, " (<emphasis>%s</emphasis>)", a2link(field.FullLink()))
		}
		fmt.Fprintf(w, "</term><listitem>")

		if field.PatchStrategy != "" {
			fmt.Fprintf(w, "<para>patch strategy: %s</para>", field.PatchStrategy)
		}
		if field.PatchMergeKey != "" {
			fmt.Fprintf(w, "<para>patch merge key: %s</para>", field.PatchMergeKey)
		}
		fmt.Fprintf(w, "<para>%s</para></listitem></varlistentry>\n", field.DescriptionWithEntities)
	}
	fmt.Fprintf(w, "</variablelist>\n")
}

func (h *DocbookWriter) Finalize() {
	// generate NavData
	os.MkdirAll(*api.ConfigDir+"/build", os.ModePerm)
	h.generateBook()
}

// END OF INTERFACE IMPLEMENTATION

func (h *DocbookWriter) writeSample(w io.Writer, d *api.Definition) {
	if d.Sample.Sample == "" {
		return
	}

	note := d.Sample.Note
	for _, s := range d.GetSamples() {
		fmt.Fprintf(w, "<example><title>%s</title><programlisting>%s</programlisting></example>\n\n", note, html.EscapeString(s.Text))
	}
}

func (h *DocbookWriter) writeOperationSample(w io.Writer, req bool, op string, examples []api.ExampleText) {

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
			msg = "<command>curl</command> command (<emphasis>requires <command>kubectl proxy</command> to be running</emphasis>)"
		} else if eType == "kubectl" && strings.Contains(msg, "Command") { // `kubectl` command
			msg = "<command>kubectl</command> command"
		}
		fmt.Fprintf(w, "<example id=\"%s\"><title>%s</title><programlisting>%s</programlisting></example>", sampleID, msg, e.Text)
	}
}

func (h *DocbookWriter) writeParams(w io.Writer, title string, params api.Fields) {
	fmt.Fprintf(w, "<bridgehead renderas=\"sect4\">%s</bridgehead>\n", title)
	fmt.Fprintf(w, "<variablelist>\n")

	for _, p := range params {
		fmt.Fprintf(w, "<varlistentry><term>%s", p.Name)

		if p.Link() != "" {
			fmt.Fprintf(w, " (<emphasis>%s</emphasis>)", a2link(p.FullLink()))
		}
		fmt.Fprintf(w, "</term><listitem>")

		fmt.Fprintf(w, "<para>%s</para></listitem></varlistentry>\n", p.Description)
	}
	fmt.Fprintf(w, "</variablelist>\n")
}

func (h *DocbookWriter) writeRequestParams(w io.Writer, o *api.Operation) {
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

func (h *DocbookWriter) writeResponseParams(w io.Writer, o *api.Operation) {
	if o.HttpResponses.Len() == 0 {
		return
	}

	fmt.Fprintf(w, "<bridgehead renderas=\"sect4\">Response</bridgehead><informaltable>\n<tgroup cols=\"2\"><thead><row><entry>Code</entry><entry>Description</entry></row></thead>\n<tbody>\n")
	for _, p := range o.HttpResponses {
		fmt.Fprintf(w, "<row><entry>%s", p.Name)
		if p.Field.Link() != "" {
			fmt.Fprintf(w, " (<emphasis>%s</emphasis>)", a2link(p.Field.FullLink()))
		}
		fmt.Fprintf(w, "</entry><entry>%s</entry></row>\n", p.Field.Description)
	}
	fmt.Fprintf(w, "</tbody></tgroup>\n</informaltable>\n")
}

func (h *DocbookWriter) generateBook() {
	html, err := os.Create(*api.ConfigDir + "/build/index.xml")
	defer html.Close()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%v", err))
		os.Exit(1)
	}

	fmt.Fprintf(html, "<?xml version=\"1.0\"?>\n<!DOCTYPE book PUBLIC \"-//OASIS//DTD DocBook XML V4.5//EN\" \"http://www.oasis-open.org/docbook/xml/4.5/docbookx.dtd\">\n")
	fmt.Fprintf(html, "<book>\n")
	fmt.Fprintf(html, "<bookinfo>\n")
	fmt.Fprintf(html, "<releaseinfo>%s</releaseinfo>\n", a2ulink(h.TOC.Copyright))
	fmt.Fprintf(html, "<releaseinfo>Generated at %s</releaseinfo>\n", time.Now().Format("2006-01-02 15:04:05 (MST)"))

	pos := strings.LastIndex(h.Config.SpecVersion, ".")
	release := fmt.Sprintf("release-%s", h.Config.SpecVersion[1:pos])
	specLink := "https://github.com/kubernetes/kubernetes/blob/" + release + "/api/openapi-spec/swagger.json"
	fmt.Fprintf(html, "<releaseinfo>API Version: <ulink url=\"%s\">%s</ulink></releaseinfo>\n", specLink, h.Config.SpecVersion)
	fmt.Fprintf(html, "<title>%s</title><subtitle>%s</subtitle>\n", h.TOC.Title, h.Config.SpecVersion)
	fmt.Fprintf(html, "</bookinfo>\n")

	buf := ""
	for _, sec := range h.TOC.Sections {
		fmt.Printf("Collecting %s ... ", sec.File)
		content, err := ioutil.ReadFile(filepath.Join(*api.ConfigDir, "includes", sec.File))
		if err == nil {
			buf += string(content)
		}

		for _, sub := range sec.SubSections {
			if len(sub.File) > 0 {
				subdata, err := ioutil.ReadFile(filepath.Join(*api.ConfigDir, "includes", sub.File))
				fmt.Printf("Collecting %s ... ", sub.File)
				if err == nil {
					buf += string(subdata)
					fmt.Printf("OK\n")
				}
			}

			if len(sub.FileClose) > 0 {
				subdata, err := ioutil.ReadFile(filepath.Join(*api.ConfigDir, "includes", sub.FileClose))
				fmt.Printf("Collecting %s ... ", sub.FileClose)
				if err == nil {
					buf += string(subdata)
					fmt.Printf("OK\n")
				}
			}
		}
		content, err = ioutil.ReadFile(filepath.Join(*api.ConfigDir, "includes", sec.FileClose))
		if err == nil {
			buf += string(content)
		}
		fmt.Printf("OK\n")
	}
	fmt.Fprintf(html, "%s", string(buf))

	fmt.Fprintf(html, "</book>\n")
}

func a2link(str string) string {
	result := strings.Replace(str, "<a href=\"#", "<link linkend=\"", -1)
	return strings.Replace(result, "</a>", "</link>", -1)
}

func a2ulink(str string) string {
	result := strings.Replace(str, "<a href=", "<ulink url=", -1)
	return strings.Replace(result, "</a>", "</ulink>", -1)
}
