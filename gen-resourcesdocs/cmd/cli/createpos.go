package cli

import (
	"github.com/kubernetes-sigs/reference-docs/gen-resourcesdocs/pkg/gettext"
	"github.com/kubernetes-sigs/reference-docs/gen-resourcesdocs/pkg/kubernetes"
	"github.com/spf13/cobra"
)

// ShowTOCCmd defines the `showtoc` subcommand
func CreatePoFilesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "po",
		Short:         "create PO file with swagger descriptions",
		Long:          "create PO file with swagger descriptions",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createPotFiles(cmd)
		},
	}

	cmd.Flags().StringP("po-directory", "p", "", "Directory containing PO files")
	cmd.MarkFlagRequired(poDirectory)

	return cmd
}

func createPotFiles(cmd *cobra.Command) error {
	file := cmd.Flag(fileOption).Value.String()
	poPath := cmd.Flag(poDirectory).Value.String()

	spec, err := kubernetes.NewSpec(file)
	if err != nil {
		return err
	}

	definitions := &spec.Swagger.Definitions

	potFiles := gettext.NewPotFiles(poPath)

	for k, def := range *definitions {
		potFiles.Add(kubernetes.Key(k), def)
	}

	err = potFiles.CreateFiles()
	if err != nil {
		return err
	}
	return nil
}
