package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Hugo defines the `kwebsite` subcommand
func KWebsite() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "kwebsite",
		Short:         "output specification for k/website",
		Long:          "output the specification in a format usable for the Kubernetes website",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			toc, err := prepareTOC(cmd)
			if err != nil {
				return fmt.Errorf("Unable to load specs and/or toc config: %v", err)
			}

			outputDir := cmd.Flag(outputDirOption).Value.String()
			templatesDir := cmd.Flag(templatesDirOption).Value.String()
			err = toc.ToKWebsite(outputDir, templatesDir)
			if err != nil {
				return err
			}

			show, err := cmd.Flags().GetBool(showDefinitionsOption)
			if err != nil {
				return err
			}
			if show {
				toc.OutputDocumentedDefinitions()
			}
			return nil
		},
	}
	cmd.Flags().StringP(configDirOption, "c", "", "Directory containing documentation configuration")
	cmd.MarkFlagRequired(configDirOption)
	cmd.Flags().StringP(outputDirOption, "o", "", "Directory to write markdown files")
	cmd.MarkFlagRequired(outputDirOption)
	cmd.Flags().StringP(templatesDirOption, "t", "", "Directory containing go templates for output")
	cmd.MarkFlagRequired(templatesDirOption)
	cmd.Flags().Bool(showDefinitionsOption, false, "Show where definitions are defined on output")
	return cmd
}
