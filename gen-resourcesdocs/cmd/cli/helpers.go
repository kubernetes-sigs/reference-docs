package cli

import (
	"path"

	"github.com/kubernetes-sigs/reference-docs/gen-resourcesdocs/pkg/config"
	"github.com/kubernetes-sigs/reference-docs/gen-resourcesdocs/pkg/kubernetes"
	"github.com/spf13/cobra"
)

// prepareTOC loads Spec and Toc config, and completes TOC
// by adding associates resources and not specifed resources in TOC
func prepareTOC(cmd *cobra.Command) (*config.TOC, error) {
	file := cmd.Flag(fileOption).Value.String()
	spec, err := kubernetes.NewSpec(file)
	if err != nil {
		return nil, err
	}

	configDir := cmd.Flag(configDirOption).Value.String()
	toc, err := config.LoadTOC(path.Join(configDir, "toc.yaml"))
	if err != nil {
		return nil, err
	}
	err = toc.PopulateAssociates(spec)
	if err != nil {
		return nil, err
	}

	toc.AddOtherResources(spec)
	toc.Definitions = &spec.Swagger.Definitions
	toc.Actions = spec.Actions
	toc.Actions.Sort()

	// TODO browse directory
	categories, err := config.LoadCategories([]string{path.Join(configDir, "fields.yaml")})
	if err != nil {
		return nil, err
	}
	toc.Categories = categories

	return toc, nil
}
