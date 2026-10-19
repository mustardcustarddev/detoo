package cli

import (
	"github.com/spf13/cobra"

	"github.com/mustardcustard/detoo/internal/app"
	"github.com/mustardcustard/detoo/internal/cli/output"
	"github.com/mustardcustard/detoo/internal/features/records/initialise"
)

func newInitCommand(t tool, application *app.Application) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Start a log of " + t.Plural,
		Long:  t.help("init.md"),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := application.Records.Initialise.Handle(cmd.Context(), initialise.Command{})
			if err != nil {
				return err
			}
			// A log that already existed is left alone, and there is no new
			// record to name.
			if result.File != "" {
				output.Path(cmd.OutOrStdout(), result.File)
			}
			return nil
		},
	}
}
