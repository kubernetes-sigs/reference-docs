# Config API generator

This is a Config API reference generator that uses the
[k8s.io/gengo](https://godoc.org/k8s.io/gengo) project to parse types and
generate API documentation from it.

This tool relies only on the Go source code containing the type definitions.
For references to external type definitions, especially those from the
`k8s.io/kubernetes` repo, you can tell the tool to generate links to the
existing docs, which can be the GoDocs if no better option exists.

## Try it out

1. Clone this repository.

2. Make sure you have go1.15+ installed. Run `make` to build the `genref`
   binary:

   ```
   cd genref
   make
   ```

3. Generate reference doc for some API types.

   ```shell
   ./genref -include kubelet-config -output output/html 
   ```

4. Visit `output/html/kubelet-config.v1beta1.html` to view the results.

## Customization

### Use a different config file

You can use the `-c` flag to specify a different configuration file for
generating the config API references. For example,

```shell
./genref -c myconfig.yaml
```

### Specify the package

You can modify the `config.yaml` file to customize the generated API
reference output.
You can also use `-include` and `-exclude` flags for the `genref` binary
to customize which package to include or exclude respectively.
The value for the `-include` and `-exclude` flags must be one of the
`name` of the `apis` listed in the `config.yaml` file, for example,
`kubelet-config`, `kube-scheduler-config`.

### Specify the output format

The tool can generate HTML pages directly or Markdown files if needed.
You can specify the output format using the `-f` flag. For example,

```shell
./genref -f markdown -include kubelet-config
```

### Customize the output template

The tool uses GoLang templates to generate HTML or Markdown.  The HTML
templates are under the `html/` subdirectory.  For Markdown, the templates are
under the `markdown/` subdirectory.

## Credit

This project is inspired and largely based on the
[gen-crd-api-reference-docs](https://github.com/ahmetb/gen-crd-api-reference-docs)
tool developed by Ahmet Alp Balkan (@ahmetb).

The tool was reworked for:

- refactored code structure;
- support to YAML config file so we can put comments in it;
- added (and tested) source for parsing Kubernetes unpublished API types;
- allow for parsing across packages which have emerged in many projects;
- use go.mod to manage package version for parsing;
- output styling changes to better align with Kubernetes API reference docs.
 
## TODOs

- [ ] Allow user to specify the top level structs.
- [ ] Add description to each API
- [ ] Mark packages that are used in the website so that they can be
      treated in a different way than other versions.

