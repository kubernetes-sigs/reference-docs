package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// ShowTOCCmd defines the `showtoc` subcommand
func ShowTOCCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "showtoc",
		Short:         "show the table of contents",
		Long:          "list the parts and chapter of the documentation",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			toc, err := prepareTOC(cmd)
			if err != nil {
				return fmt.Errorf("Unable to load specs and/or toc config: %v", err)
			}
			toc.ToMarkdown(os.Stdout)
			return nil
		},
	}
	cmd.Flags().StringP(configDirOption, "c", "", "Directory containing documentation configuration")
	cmd.MarkFlagRequired(configDirOption)

	return cmd
}
