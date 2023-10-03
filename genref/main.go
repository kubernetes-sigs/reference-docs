package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	texttemplate "text/template"

	"k8s.io/gengo/parser"
	"k8s.io/gengo/types"
	"k8s.io/klog/v2"

	"sigs.k8s.io/yaml"
)

var (
	flConfig  = flag.String("c", "config.yaml", "path to config file")
	flFormat  = flag.String("f", "markdown", "format for output, one of 'html' and 'markdown'.")
	flInclude = flag.String("include", "", "API definitions to include, comma-separated list")
	flExclude = flag.String("exclude", "", "API definitions to exclude, comma-separated list")
	flPath    = flag.String("o", ".", "path for the output files")
	// flVerbose = flag.Bool("verbose", false, "turn on verbose output")
)

const (
	docCommentForceIncludes = "// +gencrdrefdocs:force"

	CRED    = "\033[31m"
	CGREEN  = "\033[32m"
	CYELLOW = "\033[33m"
	CEND    = "\033[0m"
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

	// MainPackage is an override for API definitions that involves more
	// than one package.
	MainPackage string `json:"mainPackage"`
}

// Global vars
// Map from type definition to the API package
var typePkgMap map[string]*apiPackage
var config generatorConfig
var references map[string][]*apiType

func init() {
	klog.InitFlags(nil)

	flag.Set("logtostderr", "true")
	flag.Parse()

	var path string
	if *flFormat == "html" || *flFormat == "markdown" {
		var err error
		path, err = filepath.Abs(*flFormat)
		if err != nil {
			klog.Fatalf("template directory '%s' is not found: %w", path, err)
		}
	} else {
		klog.Fatalf("unsupported format '%s' specified", *flFormat)
	}

	fi, err := os.Stat(path)
	if err != nil {
		klog.Fatalf("cannot read the %s directory: %w", path, err)
	}
	if !fi.IsDir() {
		klog.Fatalf("%s path is not a directory", path)
	}

	typePkgMap = make(map[string]*apiPackage)
	references = make(map[string][]*apiType)
}

// processAPIPath processes a path for package enumeration and processing.
func processAPIPath(path string, includes []string, title string, mainPkg string) ([]*apiPackage, error) {
	klog.V(0).Infof("Parsing go packages in %s", path)
	gopkgs, err := parseAPIPackages(path)
	if err != nil {
		return nil, err
	}
	if len(gopkgs) == 0 {
		return nil, fmt.Errorf("no API packages found in %s", path)
	}

	for _, p := range includes {
		extra, err := parseAPIPackages(p)
		if err != nil {
			return nil, err
		}
		gopkgs = append(gopkgs, extra...)
	}

	pkgs, err := combineAPIPackages(gopkgs, title, mainPkg)
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
		return nil, fmt.Errorf("failed to parse pkgs and types: %w", err)
	}
	var pkgNames []string
	for p := range scan {
		gopkg := scan[p]
		gname := groupName(gopkg)
		klog.V(5).Infof("trying package=%s groupName=%s", p, gname)
		klog.V(6).Infof("num types=%d", len(gopkg.Types))
		// Do not pick up packages that are in vendor/ as API packages.
		if isVendorPackage(gopkg) {
			klog.Warningf("Ignoring vendor package '%v'", p)
			continue
		}

		if len(gopkg.Types) > 0 || containsString(gopkg.DocComments, docCommentForceIncludes) {
			klog.V(5).Infof("Package=%v has group name and has types", p)
			pkgNames = append(pkgNames, p)
		}
	}
	sort.Strings(pkgNames)
	var pkgs []*types.Package
	for _, p := range pkgNames {
		klog.V(5).Infof("Using package=%s", p)
		if p == dir {
			pkgs = append(pkgs, scan[p])
		}
	}
	return pkgs, nil
}

// combineAPIPackages groups the Go packages by the <apiGroup+apiVersion> they
// offer, and combines the types in them.
func combineAPIPackages(pkgs []*types.Package, title string, mainPkg string) ([]*apiPackage, error) {
	pkgMap := make(map[string]*apiPackage)
	re := `^v\d+((alpha|beta)\d+)?$`

	for _, gopkg := range pkgs {
		group := groupName(gopkg)
		// assumes basename (i.e. "v1" in "core/v1") is apiVersion
		version := gopkg.Name

		if !regexp.MustCompile(re).MatchString(version) {
			return nil, fmt.Errorf("cannot infer apiVersion for package %s (basename '%q' is not recognizable)", gopkg.Path, version)
		}

		typeList := make([]*apiType, 0, len(gopkg.Types))
		for _, t := range gopkg.Types {
			typeList = append(typeList, &apiType{*t})
		}

		id := fmt.Sprintf("%s/%s", group, version)
		v, ok := pkgMap[id]
		if !ok {
			isMain := true
			if len(mainPkg) > 0 && group != mainPkg {
				isMain = false
			}
			pkgMap[id] = &apiPackage{
				apiGroup:   group,
				apiVersion: version,
				Types:      typeList,
				GoPackages: []*types.Package{gopkg},
				Title:      title,
				IsMain:     isMain,
			}
		} else {
			v.Types = append(v.Types, typeList...)
			v.GoPackages = append(v.GoPackages, gopkg)
		}
	}
	// Sort this map
	packageIds := make([]string, 0, len(pkgMap))
	out := make([]*apiPackage, 0, len(pkgMap))
	for k := range pkgMap {
		packageIds = append(packageIds, k)
	}
	sort.Strings(packageIds)
	for _, key := range packageIds {
		out = append(out, pkgMap[key])
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
		var tmpl *template.Template
		tmpl, err = template.New("").ParseGlob(glob)
		if err != nil {
			return fmt.Errorf("parse error: %w", err)
		}

		err = tmpl.ExecuteTemplate(w, "packages", params)
	} else {
		var tmpl *texttemplate.Template
		tmpl, err = texttemplate.New("").ParseGlob(glob)
		if err != nil {
			return fmt.Errorf("parse error: %w", err)
		}

		err = tmpl.ExecuteTemplate(w, "packages", params)
	}

	if err != nil {
		return fmt.Errorf("template execution error: %w", err)
	}

	return nil
}

// writeFile creates the output file at the specified output path.
func writeFile(pkgs []*apiPackage, outputPath string) error {
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create dir %s: %w", dir, err)
	}
	var b bytes.Buffer
	if err := render(&b, pkgs); err != nil {
		return fmt.Errorf("failed to render the result: %w", err)
	}
	// s := regexp.MustCompile(`(?m)^\s+`).ReplaceAllString(b.String(), "")
	s := b.String()

	if err := os.WriteFile(outputPath, []byte(s), 0644); err != nil {
		return fmt.Errorf(CRED+"Failed to write output file: %w"+CEND, err)
	}

	klog.Infof(CGREEN+"Output written to %s"+CEND, outputPath)
	return nil
}

func main() {
	f, err := os.ReadFile(*flConfig)
	if err != nil {
		klog.Fatalf("Failed to open config file: %v", err)
	}

	if err = yaml.UnmarshalStrict(f, &config); err != nil {
		klog.Fatalf("Failed to parse config file: %v", err)
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
		pkgs, err := processAPIPath(apiDir, item.Includes, item.Title, item.MainPackage)
		if err != nil {
			klog.ErrorS(err, "cannot process API path")
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
			klog.ErrorS(err, "cannot write file")
			continue
		}
	}
}
