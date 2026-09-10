package cli

import (
	"github.com/spf13/cobra"

	"github.com/mustardcustard/decision-records/internal/app"
	"github.com/mustardcustard/decision-records/internal/cli/output"
	"github.com/mustardcustard/decision-records/internal/features/records/list"
)

func newListCommand(application *app.Application) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List the decision records",
		Long:  helpText("list.md"),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := application.Records.List.Handle(cmd.Context(), list.Query{})
			if err != nil {
				return err
			}
			output.Records(cmd.OutOrStdout(), result.Records)
			return nil
		},
	}
}
