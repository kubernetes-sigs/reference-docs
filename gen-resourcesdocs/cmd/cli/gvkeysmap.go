package cli

import (
	"fmt"
	"sort"

	"github.com/feloy/kubernetes-api-reference/pkg/kubernetes"
	"github.com/spf13/cobra"
)

// GVKeysMap defines the `gvkeysmap` subcommand
func GVKeysMap() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "gvkeysmap",
		Short:         "show the map between group/version and definition keys",
		Long:          "show the map between group/version and definition keys",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			file := cmd.Flag(fileOption).Value.String()
			spec, err := kubernetes.NewSpec(file)
			if err != nil {
				return err
			}
			gvs := make([]string, len(spec.GVToKey))
			i := 0
			for gv := range spec.GVToKey {
				gvs[i] = gv
				i++
			}
			sort.Strings(gvs)
			for _, gv := range gvs {
				keys := spec.GVToKey[gv]
				fmt.Printf("%s\n", gv)
				for _, key := range keys {
					fmt.Printf("\t%s\n", key)
				}
			}
			return nil
		},
	}
	return cmd
}
