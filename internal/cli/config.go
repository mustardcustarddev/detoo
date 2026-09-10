package cli

import (
	"github.com/spf13/cobra"

	"github.com/mustardcustard/decision-records/internal/app"
	"github.com/mustardcustard/decision-records/internal/cli/output"
	"github.com/mustardcustard/decision-records/internal/features/records/config"
)

func newConfigCommand(application *app.Application) *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Print the configuration in use",
		Long:  helpText("config.md"),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := application.Records.Config.Handle(cmd.Context(), config.Query{})
			if err != nil {
				return err
			}
			output.Config(cmd.OutOrStdout(), result.Dir, result.Template)
			return nil
		},
	}
}
