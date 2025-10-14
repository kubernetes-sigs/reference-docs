/*
Copyright 2018 The Kubernetes Authors.

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
	"bytes"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
)

type byName []*cobra.Command

func (s byName) Len() int           { return len(s) }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }

const generated_warning = `
<!--
The file is auto-generated from the Go source code of the component using a generic
[generator](https://github.com/kubernetes-sigs/reference-docs/). To learn how
to generate the reference documentation, please read
[Contributing to the reference documentation](/docs/contribute/generate-ref-docs/).
To update the reference content, please follow the
[Contributing upstream](/docs/contribute/generate-ref-docs/contribute-upstream/)
guide. You can file document formatting bugs against the
[reference-docs](https://github.com/kubernetes-sigs/reference-docs/) project.
-->
`

func GenMarkdownTree(cmd *cobra.Command, dir string, withTitle bool) error {
	identity := func(s string) string { return s }
	emptyStr := func(s string) string { return "" }
	return GenMarkdownTreeCustom(cmd, dir, "", emptyStr, identity, withTitle)
}

func GenMarkdownTreeCustom(cmd *cobra.Command, dir string, subdir string, filePrepender, linkHandler func(string) string, withTitle bool) error {
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		// CommandPath example: 'kubectl top pod'
		parts := strings.Split(c.CommandPath(), " ")
		subdir := ""
		if len(parts) > 1 {
			subdir = parts[0] + "_" + parts[1]
			fullPath := filepath.Join(dir, subdir)
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				os.Mkdir(fullPath, 0770)
			}
		}
		if err := GenMarkdownTreeCustom(c, dir, subdir, filePrepender, linkHandler, withTitle); err != nil {
			return err
		}
	}

	fullname := strings.ReplaceAll(cmd.CommandPath(), " ", "_") + ".md"
	indexFile := false
	if len(subdir) > 0 {
		parts := strings.Split(cmd.CommandPath(), " ")
		if len(parts) == 2 {
			indexFile = true
			fullname = "_index.md"
		}
		fullname = filepath.Join(subdir, fullname)
	}
	filename := filepath.Join(dir, fullname)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.WriteString(f, filePrepender(filename)); err != nil {
		return err
	}
	if err := GenReference(cmd, f, linkHandler, withTitle, indexFile); err != nil {
		return err
	}
	return nil
}

func GenReference(cmd *cobra.Command, w io.Writer, linkHandler func(string) string, withTitle bool, indexFile bool) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()
	name := cmd.CommandPath()
	short := cmd.Short
	long := cmd.Long

	if len(long) == 0 {
		long = short
	}

	// Writes a Markdown file with content_type as tool-reference.
	// Adds headings (heading shortcode).
	// Note: Files generated for kubeadm tool are snippets of Markdown without a title.
	// These snippets are included in the corresponding kubeadm pages.
	if withTitle {
		if _, err := fmt.Fprintf(w, "---\ntitle: %s\ncontent_type: tool-reference\nweight: 30\nauto_generated: true\n", name); err != nil {
			return err
		}
		if indexFile {
			if _, err := fmt.Fprintf(w, "no_list: true\n"); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(w, "---\n\n"); err != nil {
			return err
		}

		// Print the "generated" warning
		fmt.Fprintf(w, "%s\n\n", generated_warning)

		if _, err := fmt.Fprintf(w, "%s\n\n", `## {{% heading "synopsis" %}}`); err != nil {
			return err
		}

		// Escape any '<', '>' characters found in the long description
		long = strings.ReplaceAll(long, "<", "&lt;")
		long = strings.ReplaceAll(long, ">", "&gt;")
		if _, err := fmt.Fprintf(w, "\n%s\n\n", long); err != nil {
			return err
		}

		if cmd.Runnable() {
			if _, err := fmt.Fprintf(w, "```\n%s\n```\n\n", cmd.UseLine()); err != nil {
				return err
			}
		}
	} else {
		// Print the "generated" warning
		fmt.Fprintf(w, "%s\n\n", generated_warning)

		trimmedShort := strings.TrimSpace(strings.TrimSuffix(short, "."))
		trimmedLong := strings.TrimSpace(strings.TrimSuffix(long, "."))
		if !strings.HasPrefix(trimmedLong, trimmedShort) {
			if _, err := fmt.Fprintf(w, "%s\n\n", short); err != nil {
				return err
			}
		}

		if _, err := fmt.Fprintf(w, "### Synopsis\n\n"); err != nil {
			return err
		}

		// Escape any '<', '>' characters found in the long description
		long = strings.ReplaceAll(long, "<", "&lt;")
		long = strings.ReplaceAll(long, ">", "&gt;")
		if _, err := fmt.Fprintf(w, "\n%s\n\n", long); err != nil {
			return err
		}

		if cmd.Runnable() {
			if _, err := fmt.Fprintf(w, "```\n%s\n```\n\n", cmd.UseLine()); err != nil {
				return err
			}
		}
	}

	/* Examples */
	if len(cmd.Example) > 0 {
		if withTitle {
			if _, err := fmt.Fprintf(w, "%s\n\n", `## {{% heading "examples" %}}`); err != nil {
				return err
			}

			if _, err := fmt.Fprintf(w, "```\n%s\n```\n\n", cmd.Example); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(w, "### Examples\n\n"); err != nil {
				return err
			}

			if _, err := fmt.Fprintf(w, "```\n%s\n```\n\n", cmd.Example); err != nil {
				return err
			}
		}
	}

	/* Options */
	if err := printOptions(w, cmd, name, withTitle); err != nil {
		return err
	}

	/* SEE ALSO */
	if hasSeeAlso(cmd) {
		if withTitle {
			if _, err := fmt.Fprintf(w, "%s\n\n", `## {{% heading "seealso" %}}`); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(w, "%s\n\n", "SEE ALSO"); err != nil {
				return err
			}
		}

		if cmd.HasParent() {
			parent := cmd.Parent()
			pname := parent.CommandPath()
			link := "../" + pname + "/"
			if !indexFile {
				link = "../"
			}
			link = strings.ReplaceAll(link, " ", "_")
			if _, err := fmt.Fprintf(w, "* [%s](%s)\t - %s\n", pname, linkHandler(link), parent.Short); err != nil {
				return err
			}
			cmd.VisitParents(func(c *cobra.Command) {
				if c.DisableAutoGenTag {
					cmd.DisableAutoGenTag = c.DisableAutoGenTag
				}
			})
		}

		children := cmd.Commands()
		sort.Sort(byName(children))

		for _, child := range children {
			if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
				continue
			}
			cname := name + " " + child.Name()
			link := strings.ReplaceAll(cname, " ", "_") + "/"
			if !indexFile {
				link = "../" + link
			}

			// link = strings.ReplaceAll(link, " ", "_")
			if _, err := fmt.Fprintf(w, "* [%s](%s)\t - %s\n", cname, linkHandler(link), child.Short); err != nil {
				return err
			}
		}

		if _, err := fmt.Fprintf(w, "\n"); err != nil {
			return err
		}
	}

	return nil
}

func printOptions(w io.Writer, cmd *cobra.Command, name string, withTitle bool) error {
	flags := cmd.NonInheritedFlags()
	flags.SetOutput(w)

	if flags.HasFlags() {
		if withTitle {
			if _, err := fmt.Fprintf(w, "%s\n\n", `## {{% heading "options" %}}`); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprint(w, "### Options\n\n"); err != nil {
				return err
			}
		}
		usages := flagUsages(flags)
		fmt.Fprint(w, usages)
		if _, err := fmt.Fprint(w, "\n"); err != nil {
			return err
		}
	}

	parentFlags := cmd.InheritedFlags()
	parentFlags.SetOutput(w)
	if parentFlags.HasFlags() {

		if withTitle {
			if _, err := fmt.Fprintf(w, "%s\n\n", `## {{% heading "parentoptions" %}}`); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprint(w, "### Options inherited from parent commands\n\n"); err != nil {
				return err
			}
		}
		usages := flagUsages(parentFlags)
		fmt.Fprint(w, usages)

		if _, err := fmt.Fprint(w, "\n"); err != nil {
			return err
		}
	}
	return nil
}

func hasSeeAlso(cmd *cobra.Command) bool {
	if cmd.HasParent() {
		return true
	}
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		return true
	}
	return false
}

func flagUsages(f *pflag.FlagSet) string {
	x := new(bytes.Buffer)

	lines := make([]string, 0)

	lines = append(lines, "   <table style=\"width: 100%; table-layout: fixed;\">\n<colgroup>\n"+
		"<col span=\"1\" style=\"width: 10px;\" />\n"+
		"<col span=\"1\" />\n"+
		"</colgroup>\n"+
		"<tbody>\n")
	f.VisitAll(func(flag *pflag.Flag) {
		// if len(flag.Deprecated) > 0 || flag.Hidden {
		if flag.Hidden {
			return
		}

		line := "<tr>\n<td colspan=\"2\">"

		if len(flag.Shorthand) > 0 && len(flag.ShorthandDeprecated) == 0 {
			line += fmt.Sprintf("-%s, --%s", flag.Shorthand, flag.Name)
		} else {
			line += fmt.Sprintf("--%s", flag.Name)
		}

		varname, usage := unquoteUsage(flag)
		if len(varname) > 0 {
			line += " " + varname
		}
		if len(flag.NoOptDefVal) > 0 {
			switch flag.Value.Type() {
			case "string":
				line += fmt.Sprintf("[=\"%s\"]", flag.NoOptDefVal)
			case "bool":
				if flag.NoOptDefVal != "true" {
					line += fmt.Sprintf("[=%s]", flag.NoOptDefVal)
				}
			default:
				line += fmt.Sprintf("[=%s]", flag.NoOptDefVal)
			}
		}
		if !defaultIsZeroValue(flag) {
			defaultValue := flag.DefValue
			if flag.Value.Type() == "string" {
				// There are cases where the string is very very long, split
				// it to mutiple lines manually
				if len(defaultValue) > 40 {
					defaultValue = strings.ReplaceAll(defaultValue, ",", ",<br />")
				}
				// clean up kubectl cache-dir flag value
				if strings.Compare(flag.Name, "cache-dir") == 0 {
					myUser, err := user.Current()
					if err == nil {
						noprefix := strings.TrimPrefix(defaultValue, myUser.HomeDir)
						defaultValue = fmt.Sprintf("$HOME%s", noprefix)
					}
				}
				line += fmt.Sprintf("&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Default: \"%s\"", defaultValue)
			} else if flag.Value.Type() == "stringSlice" {
				defaultValue := flag.DefValue
				// For string slices, the default value should not contain '[' ]r ']'
				defaultValue = strings.TrimPrefix(defaultValue, "[")
				defaultValue = strings.TrimSuffix(defaultValue, "]")
				defaultValue = strings.ReplaceAll(defaultValue, " ", "")
				line += fmt.Sprintf("&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Default: \"%s\"", defaultValue)
			} else {
				line += fmt.Sprintf("&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Default: %s", defaultValue)
			}
		}
		line += "</td>\n</tr>\n<tr>\n<td></td><td style=\"line-height: 130%; word-wrap: break-word;\">"

		// process markdown in usage, force wrap for "\n"
		if len(flag.Deprecated) > 0 {
			usage = usage + " (DEPRECATED: " + flag.Deprecated + ")"
		}
		line += processUsage(usage) + "</td>\n</tr>\n"

		lines = append(lines, line)
	})
	lines = append(lines, "</tbody>\n</table>\n\n")

	for _, line := range lines {
		// fmt.Fprintln(x, line[:sidx], spacing, wrap(maxlen+2, cols, line[sidx+1:]))
		fmt.Fprintln(x, line)
	}

	return x.String()
}

func defaultIsZeroValue(f *pflag.Flag) bool {
	switch f.Value.Type() {
	case "bool":
		return f.DefValue == "false"
	case "duration":
		return f.DefValue == "0" || f.DefValue == "0s"
	case "int", "int8", "int32", "int64", "uint", "uint8", "uint16", "uint32", "count", "float32", "float64":
		return f.DefValue == "0"
	case "string":
		return f.DefValue == ""
	case "ip", "ipMask", "ipNet":
		return f.DefValue == "<nil>"
	case "intSlice", "stringSlice", "stringArray":
		return f.DefValue == "[]"
	case "namedCertKey":
		return f.DefValue == "[]"
	default:
		switch f.Value.String() {
		case "false":
			return true
		case "<nil>":
			return true
		case "":
			return true
		case "0":
			return true
		}
		return false
	}
}

// Adapted usage function from pflag
func unquoteUsage(flag *pflag.Flag) (name string, usage string) {
	// Look for a back-quoted name, but avoid the strings package.
	usage = flag.Usage
	for i := 0; i < len(usage); i++ {
		if usage[i] == '`' {
			for j := i + 1; j < len(usage); j++ {
				if usage[j] == '`' {
					name = usage[i+1 : j]
					usage = usage[:i] + name + usage[j+1:]
					return name, usage
				}
			}
			break // Only one back quote; use type name.
		}
	}

	name = flag.Value.Type()
	switch name {
	case "bool":
		name = ""
	case "float64", "float32":
		name = "float"
	case "int64", "severity", "Level":
		name = "int"
	case "uint64":
		name = "uint"
	case "stringSlice", "stringArray":
		name = "strings"
	case "intSlice":
		name = "ints"
	case "uintSlice":
		name = "uints"
	case "boolSlice":
		name = "bools"
	case "mapStringString":
		name = "&lt;comma-separated 'key=value' pairs&gt;"
	case "mapStringBool":
		name = "&lt;comma-separated 'key=True|False' pairs&gt;"
	case "namedCertKey", "ip":
		name = "string"
	case "moduleSpec":
		name = "&lt;comma-separated 'pattern=N' settings&gt;"
	case "portRange":
		name = "&lt;a string in the form 'N1-N2'&gt;"
	case "traceLocation":
		name = "&lt;a string in the form 'file:N'&gt;"
	}
	return
}

func processUsage(usage string) string {
	var buf bytes.Buffer
	var result string
	usage = strings.ReplaceAll(usage, "<", "&lt;")
	usage = strings.ReplaceAll(usage, ">", "&gt;")
	md := goldmark.New(goldmark.WithExtensions(highlighting.Highlighting))
	if err := md.Convert([]byte(usage), &buf); err != nil {
		result = usage
	} else {
		result = buf.String()
	}
	result = strings.TrimSuffix(result, "\n")
	result = strings.ReplaceAll(result, "\n", "<br/>")
	return result
}
