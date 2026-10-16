package cli

import (
	"github.com/spf13/cobra"

	"github.com/mustardcustard/detoo/internal/app"
	"github.com/mustardcustard/detoo/internal/cli/output"
	"github.com/mustardcustard/detoo/internal/features/records/config"
)

func newConfigCommand(t tool, application *app.Application) *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Print the configuration in use",
		Long:  t.help("config.md"),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := application.Records.Config.Handle(cmd.Context(), config.Query{})
			if err != nil {
				return err
			}
			output.Config(cmd.OutOrStdout(), t.Name, result.Dir, result.Template)
			return nil
		},
	}
}
