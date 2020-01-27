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
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type TOCItem struct {
	Level       int
	Title       string
	Link        string
	Buffer      string
	SubSections []*TOCItem
}

type TOC struct {
	Title     string
	Copyright string
	Sections  []*TOCItem
}

type HTMLWriter struct {
	TOC            TOC
	CurrentSection *TOCItem
}

func NewHTMLWriter(copyright, title string) DocWriter {
	writer := HTMLWriter{
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

func (h *HTMLWriter) WriteCommands(toc ToC, params KubectlSpec) {

	/* map the top level command name to command structure */
	m := map[string]TopLevelCommand{}
	for _, g := range params.TopLevelCommandGroups {
		for _, tlc := range g.Commands {
			m[tlc.MainCommand.Name] = tlc
		}
	}

	for _, c := range toc.Categories {

		fmt.Printf("Category name:%s", c.Name)

		// write a category page
		h.writeCategoryPage(c.Name)

		// Write each of the commands in this category
		for _, cm := range c.Commands {
			if tlc, found := m[cm]; !found {
				fmt.Printf("Could not find top level command %s\n", cm)
				os.Exit(1)
			} else {

				// write pages to buffer
				h.writeCommandPage(tlc)

				delete(m, cm)
			}
		}
	}
	if len(m) > 0 {
		for k := range m {
			fmt.Printf("Kubectl command %s missing from toc.yaml\n", k)
		}
		os.Exit(1)
	}

}

// write category or group Name
func (h *HTMLWriter) writeCategorySection(groupName string) string {
	var category = strings.ToLower(strings.Replace(groupName, " ", "-", -1))
	return fmt.Sprintf("<DIV class=\"category-name\"><H2 id=\"-strong-%s-strong-\"><STRONG>%s</STRONG></H2>\n</DIV>", category, groupName)
}

// write command Name
func (h *HTMLWriter) writeName(cmd *Command) string {
	return fmt.Sprintf("<DIV><H2 id=\"-strong-%s-strong-\"><STRONG>%s</STRONG></H2>\n</DIV>", cmd.Name, cmd.Name)
}

// write subcommand Path
func (h *HTMLWriter) writePath(cmd *Command, mainCmdName string) string {
	fmt.Printf("main cmd name:%s,path:%s\n", mainCmdName, cmd.Path)
		if len(cmd.Path) > 0 {
			return fmt.Sprintf("<DIV><H3 id=\"-strong-%s-%s-strong-\"><STRONG>%s <em>%s</em></STRONG></H3>\n</DIV>", mainCmdName, cmd.Path, mainCmdName, cmd.Path)
		}
		return ""
}

// write command Description
func (h *HTMLWriter) writeDescription(cmd *Command) string {
	if len(cmd.Description) > 0 {
		return fmt.Sprintf("<DIV><p>%s</p>\n</DIV>", template.HTMLEscapeString(cmd.Description))
	}
	return ""
}

// write command Usage
func (h *HTMLWriter) writeUsage(cmd *Command) string {
	if len(cmd.Usage) > 0 {
		return fmt.Sprintf("<DIV><p><H3 id=\"usage\">Usage</H3><p><code>%s</code></p></DIV>", cmd.Usage)
	}
	return ""
}

// write command Example
// REVISIT: this is an html fomatted string from markdown
func (h *HTMLWriter) writeExamples(cmd *Command) string {
	if len(cmd.Example) > 0 {
		return fmt.Sprintf("<DIV><H3>Examples</H3><p>%s</p>\n</DIV>", cmd.Example)
	}
	return ""
}

// write a command
func (h *HTMLWriter) writeCmdSection(cmd *Command) string {
	var buf string
	buf = h.writeName(cmd)
	buf += h.writeDescription(cmd)
	buf += h.writeUsage(cmd)
	buf += h.writeOptions(cmd)
	buf += h.writeExamples(cmd)
	return buf
}

// write a sub command
func (h *HTMLWriter) writeSubCmdSection(cmd *Command, mainCmdName string) string {
	var buf string
	buf += h.writePath(cmd, mainCmdName)
	buf += h.writeDescription(cmd)
	buf += h.writeUsage(cmd)
	buf += h.writeOptions(cmd)
	buf += h.writeExamples(cmd)
	return buf
}


// write command Options
func (h *HTMLWriter) writeOptions(cmd *Command) string {
	if len(cmd.Options) > 0 {
		var buf string
		buf = "<DIV><p><H3 id=\"flags\">Flags</H3></p><br>"
		buf += "<div class=\"table-responsive table-options\"><table class=\"table table-bordered\">" +
			"<thead class=\"thead-light\"><tr>" +
			"<th>Name</th>" +
			"<th>Shorthand</th>" +
			"<th>Default</th>" +
			"<th>Usage</th>" +
			"</tr>" +
			"</thead>" +
			"<tbody>"

		for _, option := range cmd.Options {
			buf += "<tr><td><code>" + option.Name + "</code></td><td>" + option.Shorthand + "</td><td>" +
				option.DefaultValue + "</td><td>" + option.Usage + "</td></tr>"
		}
		buf += "</tbody></table></DIV>"
		return buf
	}
	return ""
}

// write category
func (h *HTMLWriter) writeCategoryPage(groupName string) {

	var buf string
	buf = h.writeCategorySection(groupName)

	var category = strings.ToLower(strings.Replace(groupName, " ", "-", -1))

	item := TOCItem{
		Level:  1,
		Title:  groupName,
		Link:   "-strong-" + category + "-strong-",
		Buffer: buf,
	}
	h.TOC.Sections = append(h.TOC.Sections, &item)
	h.CurrentSection = &item
}

// write command
func (h *HTMLWriter) writeCommandPage(params TopLevelCommand) {

	var buf string
	buf = h.writeCmdSection(params.MainCommand)

	item := TOCItem{
		Level:  1,
		Title:  params.MainCommand.Name,
		Link:   "-strong-" + params.MainCommand.Name + "-strong-",
		Buffer: buf,
	}
	h.TOC.Sections = append(h.TOC.Sections, &item)
	h.CurrentSection = &item

	// for each sub command, create cmd and link item in TOC
	for _, c := range params.SubCommands {
		h.writeSubCommandPage(c, params.MainCommand.Name)
	}
}

// write sub command
func (h *HTMLWriter) writeSubCommandPage(subCmd *Command, mainCmdName string) {
	var buf string
	buf = h.writeSubCmdSection(subCmd, mainCmdName)

	item := TOCItem{
		Level:  2,
		//Title:  mainCmdName + " " + subCmd.Path,
		Title:  subCmd.Path,
		Link:   "-strong-" + mainCmdName + "-" + subCmd.Path + "-strong-",
		Buffer: buf,
	}
	h.CurrentSection.SubSections = append(h.CurrentSection.SubSections, &item)
}

func (h *HTMLWriter) generateNavContent() string {
	nav := ""
	for _, sec := range h.TOC.Sections {
		// class for level-1 navigation item
		nav += "<UL>\n"
		if strings.Contains(sec.Link, "strong") {
			nav += fmt.Sprintf(" <LI class=\"nav-level-1 strong-nav\"><A href=\"#%s\" class=\"nav-item\"><STRONG>%s</STRONG></A></LI>\n", sec.Link, sec.Title)
		} else {
			nav += fmt.Sprintf(" <LI class=\"nav-level-1\"><A href=\"#%s\" class=\"nav-item\">%s</A></LI>\n",
				sec.Link, sec.Title)
		}

		// close H1 items which have no subsections or strong navs
		if len(sec.SubSections) == 0 || (sec.Level == 1 && strings.Contains(sec.Link, "strong")) {
			nav += "</UL>\n"
		}

		// short circuit to next if no sub-sections
		if len(sec.SubSections) == 0 {
			continue
		}

		// wrapper1
		nav += fmt.Sprintf(" <UL id=\"%s-nav\" style=\"display: none;\">\n", sec.Link)
		for _, sub := range sec.SubSections {
			nav += "  <UL>\n"
			if strings.Contains(sub.Link, "strong") {
				nav += fmt.Sprintf("   <LI class=\"nav-level-%d strong-nav\"><A href=\"#%s\" class=\"nav-item\"><STRONG>%s</STRONG></A></LI>\n",
					sub.Level, sub.Link, sub.Title)
			} else {
				nav += fmt.Sprintf("   <LI class=\"nav-level-%d\"><A href=\"#%s\" class=\"nav-item\"><em>%s</em></A></LI>\n",
					sub.Level, sub.Link, sub.Title)
			}
			// close this H1/H2 if possible
			if len(sub.SubSections) == 0 {
				nav += " </UL>\n"
				continue
			}
			// end wrapper2
			nav += "   </UL>\n"
			nav += "  </UL>\n"
		}
		// end wrapper1
		nav += " </UL>\n"
		// end top UL
		nav += "</UL>\n"
	}

	return nav
}

func (h *HTMLWriter) generateNavJS() {
	// generate NavData
	var tmp string
	flatToc := []string{}
	os.MkdirAll(BuildDir, os.ModePerm)

	navjs, err := os.Create(filepath.Join(BuildDir, "navData.js"))
	defer navjs.Close()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%v", err))
		os.Exit(1)
	}

	s1 := []string{}
	for _, sec := range h.TOC.Sections {
		flatToc = append([]string{"\"" + sec.Link + "\""}, flatToc...)
		s2 := []string{}
		for _, sub := range sec.SubSections {
			flatToc = append([]string{"\"" + sub.Link + "\""}, flatToc...)
		}

		if strings.Contains(sec.Link, "strong") {
			tmp = "{\"section\":\"" + sec.Link + "\",\"subsections\":[]}"
			s2res := strings.Join(s2, ",")
			if len(s2res) > 0 {
				tmp = s2res + "," + tmp
			}
		} else {
			tmp = "{\"section\":\"" + sec.Link + "\",\"subsections\":[" + strings.Join(s2, ",") + "]}"
		}
		s1 = append([]string{tmp}, s1...)
	}
	fmt.Fprintf(navjs, "(function(){navData={\"toc\":["+strings.Join(s1, ",")+"],\"flatToc\":["+strings.Join(flatToc, ",")+"]};})();")
}

func (h *HTMLWriter) generateHTML(navContent string) {
	html, err := os.Create(filepath.Join(BuildDir, "kubectl-commands.html"))
	defer html.Close()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%v", err))
		os.Exit(1)
	}

	fmt.Fprintf(html, "<!DOCTYPE html>\n<HTML>\n<HEAD>\n<META charset=\"UTF-8\">\n")
	fmt.Fprintf(html, "<TITLE>%s</TITLE>\n", h.TOC.Title)
	fmt.Fprintf(html, "<LINK rel=\"shortcut icon\" href=\"favicon.ico\" type=\"image/vnd.microsoft.icon\">\n")
	fmt.Fprintf(html, "<LINK rel=\"stylesheet\" href=\"css/bootstrap.min.css\">\n")
	fmt.Fprintf(html, "<LINK rel=\"stylesheet\" href=\"css/font-awesome.min.css\" type=\"text/css\">\n")
	fmt.Fprintf(html, "<LINK rel=\"stylesheet\" href=\"css/stylesheet.css\" type=\"text/css\">\n")
	fmt.Fprintf(html, "</HEAD>\n<BODY>\n")
	fmt.Fprintf(html, "<DIV id=\"wrapper\">\n")
	fmt.Fprintf(html, "<DIV id=\"sidebar-wrapper\" class=\"side-nav side-bar-nav\">\n")

	// html buffer
	buf := "<DIV class=\"row\">\n  <DIV class=\"col-md-6 copyright\">\n " + h.TOC.Copyright + "\n  </DIV>\n"
	buf += "  <DIV class=\"col-md-6 text-right\">\n"
	buf += fmt.Sprintf("    <DIV>Generated at: %s</DIV>\n", time.Now().Format("2006-01-02 15:04:05 (MST)"))
	buf += "</DIV></DIV>\n"
	buf += fmt.Sprintf("<DIV class=\"col-md-12 text-left kubectl-title\"><H1>" + h.TOC.Title + " v" + *KubernetesRelease)
	buf += "  </H1></DIV>\n"

	const OK = "\033[32mOK\033[0m"
	const NOT_FOUND = "\033[31mNot found\033[0m"
	for _, sec := range h.TOC.Sections {
		content := sec.Buffer

		if len(content) > 0 {
			buf += string(content)
			fmt.Printf(sec.Title, OK)
		} else {
			fmt.Printf(NOT_FOUND)
		}

		for _, sub := range sec.SubSections {
			subdata := sub.Buffer
			if len(subdata) > 0 {
				buf += string(subdata)
				fmt.Println(sub.Title, OK)
			} else {
				fmt.Println(NOT_FOUND)
			}
		}
	}

	fmt.Fprintf(html, "%s</DIV>\n", navContent)
	fmt.Fprintf(html, "<DIV id=\"page-content-wrapper\" class=\"body-content container\">\n")
	fmt.Fprintf(html, "%s", string(buf))
	fmt.Fprintf(html, "</DIV>\n</DIV>\n")
	fmt.Fprintf(html, "<SCRIPT src=\"/js/jquery-3.2.1.min.js\"></SCRIPT>\n")
	fmt.Fprintf(html, "<SCRIPT src=\"js/jquery.scrollTo.min.js\"></SCRIPT>\n")
	fmt.Fprintf(html, "<SCRIPT src=\"/js/bootstrap-4.3.1.min.js\"></SCRIPT>\n")
	fmt.Fprintf(html, "<SCRIPT src=\"js/navData.js\"></SCRIPT>\n")
	fmt.Fprintf(html, "<SCRIPT src=\"js/scroll.js\"></SCRIPT>\n")
	fmt.Fprintf(html, "</BODY>\n</HTML>\n")
}

func (h *HTMLWriter) Finalize() {
	os.MkdirAll(BuildDir, os.ModePerm)
	h.generateNavJS()
	navContent := h.generateNavContent()
	h.generateHTML(navContent)
}
