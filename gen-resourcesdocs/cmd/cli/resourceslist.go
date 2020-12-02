package cli

import (
	"fmt"
	"sort"

	"github.com/feloy/kubernetes-api-reference/pkg/kubernetes"
	"github.com/spf13/cobra"
)

// ResourceslistCmd defines the `resourceslist` subcommand
func ResourceslistCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "resourceslist",
		Short:         "list k8s resources",
		Long:          "list Kubernetes resources in the specification",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			file := cmd.Flag(fileOption).Value.String()
			spec, err := kubernetes.NewSpec(file)
			if err != nil {
				return err
			}

			resources := spec.Resources
			i := 0
			keys := make([]string, len(*resources))
			for k := range *resources {
				keys[i] = k.String()
				i++
			}
			sort.Strings(keys)
			for _, k := range keys {
				rs := (*resources)[kubernetes.APIKind(k)]
				fmt.Println(k)
				for _, r := range rs {
					fmt.Println("\t" + r.GetGV())
				}
			}
			return nil
		},
	}
	return cmd
}
