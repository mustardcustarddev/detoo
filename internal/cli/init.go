package cli

import (
	"github.com/spf13/cobra"

	"github.com/mustardcustard/decision-records/internal/app"
	"github.com/mustardcustard/decision-records/internal/cli/output"
	"github.com/mustardcustard/decision-records/internal/features/records/initialise"
)

func newInitCommand(application *app.Application) *cobra.Command {
	return &cobra.Command{
		Use:   "init [DIRECTORY]",
		Short: "Start a log of decision records",
		Long:  helpText("init.md"),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// No directory means the project keeps whichever it already uses,
			// which the slice reads as an empty Dir.
			var dir string
			if len(args) > 0 {
				dir = args[0]
			}
			result, err := application.Records.Initialise.Handle(cmd.Context(), initialise.Command{Dir: dir})
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
