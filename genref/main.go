package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	texttemplate "text/template"

	"github.com/pkg/errors"
	"k8s.io/gengo/parser"
	"k8s.io/gengo/types"

	"sigs.k8s.io/yaml"
)

var (
	flConfig  = flag.String("c", "config.yaml", "path to config file")
	flFormat  = flag.String("f", "markdown", "format for output, one of 'html' and 'markdown'.")
	flInclude = flag.String("include", "", "API definitions to include, comma-separated list")
	flExclude = flag.String("exclude", "", "API definitions to exclude, comma-separated list")
	flPath    = flag.String("o", ".", "path for the output files")
	flVerbose = flag.Bool("v", false, "turn on verbose output")
)

const (
	docCommentForceIncludes = "// +gencrdrefdocs:force"
)

type generatorConfig struct {
	// HiddenMemberFields hides fields with specified names on all types.
	HiddenMemberFields []string `json:"hiddenMemberFields"`

	// HideTypePatterns hides types matching the specified patterns from the
	// output.
	HideTypePatterns []string `json:"hideTypePatterns"`

	// ExternalPackages lists recognized external package references and how to
	// link to them.
	ExternalPackages []externalPackage `json:"externalPackages"`

	// StripPrefix is a list of type name prefixes that will be stripped
	StripPrefix []string `json:"stripPrefix"`

	// MarkdownDisabled controls markdown rendering for comment lines.
	MarkdownDisabled bool `json:"markdownDisabled"`

	// APIs to process
	Definitions []apiDefinition `json:"apis"`
}

type externalPackage struct {
	// Match is a reqular expression for matching type names which are defined
	// and documented externally.
	Match string `json:"match"`

	// Target provides a text template string for building the link to the
	// external documentation for a type.
	Target string `json:"target"`
}

// apiDefinition is a local struct for specifying the API type definitions for
// which reference documentations are to be generated. These definitions are
// provided and customized in the configuration YAML as well.
type apiDefinition struct {
	// Name is the key string that represents a specific package
	Name string `json:"name"`

	// Title is the string that will appear as the title of the generated page
	Title string `json:"title"`

	// Package is the import path for the API package where a type is defined.
	Package string `json:"package"`

	// Path is the path for an API type/resource definition. Each package has
	// a different convention of defining its API types.
	Path string `json:"path"`

	// Skip is a boolean flag indicating whether the package currently has
	// some problems in generating reference docs. We tag a package as
	// skipped if the current generator doesn't work on it.
	Skip bool `json:"skip,omitempty"`

	// Includes is list of packages that are designed for shared type
	// definitions.
	Includes []string `json:"includes"`
}

// Global vars
// Map from type definition to the API package
var typePkgMap map[string]*apiPackage
var config generatorConfig
var references map[string][]*apiType

func init() {
	flag.Parse()

	if *flConfig == "" {
		panic("-config not specified")
	}
	var path string
	var err error
	if *flFormat == "html" || *flFormat == "markdown" {
		path, err = filepath.Abs(*flFormat)
	} else {
		panic(errors.Errorf("unsupported format '%s' specified", *flFormat))
	}

	if err != nil {
		panic(errors.Wrapf(err, "template directory '%s' is not found", path))
	}
	fi, err := os.Stat(path)
	if err != nil {
		panic(errors.Wrapf(err, "cannot read the %s directory", path))
	}
	if !fi.IsDir() {
		panic(errors.Errorf("%s path is not a directory", path))
	}

	typePkgMap = make(map[string]*apiPackage)
	references = make(map[string][]*apiType)
}

// processAPIPath processes a path for package enumeration and processing.
func processAPIPath(path string, includes []string, title string) ([]*apiPackage, error) {
	pinfo("Parsing go packages in %s", path)
	gopkgs, err := parseAPIPackages(path)
	if err != nil {
		return nil, err
	}
	if len(gopkgs) == 0 {
		return nil, errors.Errorf("no API packages found in %s", path)
	}

	for _, p := range includes {
		extra, err := parseAPIPackages(p)
		if err != nil {
			return nil, err
		}
		for _, e := range extra {
			gopkgs = append(gopkgs, e)
		}
	}

	pkgs, err := combineAPIPackages(gopkgs, title)
	if err != nil {
		return nil, err
	}

	// Update typePkgMap and references map
	for _, p := range pkgs {
		for _, t := range p.Types {
			typePkgMap[t.String()] = p
			for _, m := range t.Members {
				mType := &apiType{*m.Type}
				rt := mType.deref().String()
				references[rt] = append(references[rt], t)
			}
		}
	}

	return pkgs, nil
}

// parseAPIPackages scans a given directory for packages.
func parseAPIPackages(dir string) ([]*types.Package, error) {
	b := parser.New()
	// the following will silently fail (turn on -v=4 to see logs)
	if err := b.AddDirRecursive(dir); err != nil {
		return nil, err
	}
	scan, err := b.FindTypes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse pkgs and types")
	}
	var pkgNames []string
	for p := range scan {
		gopkg := scan[p]
		gname := groupName(gopkg)
		pverbose("trying package=%s groupName=%s", p, gname)
		pverbose("num types=%d", len(gopkg.Types))
		// Do not pick up packages that are in vendor/ as API packages.
		if isVendorPackage(gopkg) {
			pwarning("Ignoring vendor package '%v'", p)
			continue
		}

		if len(gopkg.Types) > 0 || containsString(gopkg.DocComments, docCommentForceIncludes) {
			pverbose("Package=%v has group name and has types", p)
			pkgNames = append(pkgNames, p)
		}
	}
	sort.Strings(pkgNames)
	var pkgs []*types.Package
	for _, p := range pkgNames {
		pverbose("Using package=%s", p)
		if p == dir {
			pkgs = append(pkgs, scan[p])
		}
	}
	return pkgs, nil
}

// combineAPIPackages groups the Go packages by the <apiGroup+apiVersion> they
// offer, and combines the types in them.
func combineAPIPackages(pkgs []*types.Package, title string) ([]*apiPackage, error) {
	pkgMap := make(map[string]*apiPackage)
	re := `^v\d+((alpha|beta)\d+)?$`

	for _, gopkg := range pkgs {
		group := groupName(gopkg)
		// assumes basename (i.e. "v1" in "core/v1") is apiVersion
		version := gopkg.Name

		if !regexp.MustCompile(re).MatchString(version) {
			return nil, errors.Errorf("cannot infer apiVersion for package %s (basename '%q' is not recognizable)", gopkg.Path, version)
		}

		typeList := make([]*apiType, 0, len(gopkg.Types))
		for _, t := range gopkg.Types {
			typeList = append(typeList, &apiType{*t})
		}

		id := fmt.Sprintf("%s/%s", group, version)
		v, ok := pkgMap[id]
		if !ok {
			pkgMap[id] = &apiPackage{
				apiGroup:   group,
				apiVersion: version,
				Types:      typeList,
				GoPackages: []*types.Package{gopkg},
				Title:      title,
			}
		} else {
			v.Types = append(v.Types, typeList...)
			v.GoPackages = append(v.GoPackages, gopkg)
		}
	}
	out := make([]*apiPackage, 0, len(pkgMap))
	for _, v := range pkgMap {
		out = append(out, v)
	}
	return out, nil
}

// render is the render procedure for templating.
func render(w io.Writer, pkgs []*apiPackage) error {
	var err error

	gitCommit, _ := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	params := map[string]interface{}{
		"packages":  pkgs,
		"config":    config,
		"gitCommit": strings.TrimSpace(string(gitCommit)),
	}

	glob := filepath.Join(*flFormat, "*.tpl")
	if *flFormat == "html" {
		tmpl, err := template.New("").ParseGlob(glob)
		if err != nil {
			return errors.Wrap(err, "parse error")
		}

		err = tmpl.ExecuteTemplate(w, "packages", params)
	} else {
		tmpl, err := texttemplate.New("").ParseGlob(glob)
		if err != nil {
			return errors.Wrap(err, "parse error")
		}

		err = tmpl.ExecuteTemplate(w, "packages", params)
	}

	return errors.Wrap(err, "template execution error")
}

// writeFile creates the output file at the specified output path.
func writeFile(pkgs []*apiPackage, outputPath string) error {
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.Errorf("failed to create dir %s: %v", dir, err)
	}
	var b bytes.Buffer
	if err := render(&b, pkgs); err != nil {
		return errors.Wrap(err, "failed to render the result")
	}
	// s := regexp.MustCompile(`(?m)^\s+`).ReplaceAllString(b.String(), "")
	s := b.String()

	if err := ioutil.WriteFile(outputPath, []byte(s), 0644); err != nil {
		return errors.Errorf("failed to write to out file: %v", err)
	}

	pinfo("Output written to %s", outputPath)
	return nil
}

func main() {
	f, err := ioutil.ReadFile(*flConfig)
	if err != nil {
		perror("Failed to open config file: %+v", err)
		os.Exit(-1)
	}

	if err = yaml.UnmarshalStrict(f, &config); err != nil {
		perror("Failed to parse config file: %+v", err)
		os.Exit(-1)
	}

	pkgInclude := []string{}
	pkgExclude := []string{}
	if *flInclude != "" {
		pkgInclude = strings.Split(*flInclude, ",")
	}
	if *flExclude != "" {
		pkgExclude = strings.Split(*flExclude, ",")
	}

	for _, item := range config.Definitions {
		if item.Skip {
			continue
		}

		parts := []string{item.Package, item.Path}
		apiDir := strings.Join(parts, "/")
		// determine package to explicitly exclude, or include
		if len(pkgExclude) > 0 && containsString(pkgExclude, item.Name) {
			continue
		}
		if len(pkgInclude) > 0 && !containsString(pkgInclude, item.Name) {
			continue
		}
		pkgs, err := processAPIPath(apiDir, item.Includes, item.Title)
		if err != nil {
			perror("%+v", err)
			continue
		}

		segments := strings.Split(item.Path, "/")
		version := segments[len(segments)-1]
		fn := fmt.Sprintf("%s/%s.%s", *flPath, item.Name, version)
		if *flFormat == "html" {
			fn = fn + ".html"
		} else if *flFormat == "markdown" {
			fn = fn + ".md"
		}
		if err = writeFile(pkgs, fn); err != nil {
			perror("%+v", err)
			continue
		}
	}
}
