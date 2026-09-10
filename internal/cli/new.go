package cli

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/mustardcustard/decision-records/internal/app"
	"github.com/mustardcustard/decision-records/internal/cli/output"
	"github.com/mustardcustard/decision-records/internal/features/records/create"
)

func newNewCommand(application *app.Application) *cobra.Command {
	var supercedes, updates []string

	cmd := &cobra.Command{
		Use:     "new TITLE",
		Short:   "Create the next numbered decision record",
		Long:    helpText("new.md"),
		Example: helpText("new.example"),
		Args:    titleArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := application.Records.Create.Handle(cmd.Context(), create.Command{
				Title:      args[0],
				Supercedes: supercedes,
				Updates:    updates,
			})
			if err != nil {
				return err
			}
			output.Path(cmd.OutOrStdout(), result.File)
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&supercedes, "supercede", "s", nil,
		"a `REFERENCE` (number or partial filename) to a decision this one supercedes")
	cmd.Flags().StringArrayVarP(&updates, "update", "u", nil,
		"a `REFERENCE` (number or partial filename) to a decision this one updates")
	return cmd
}

// titleArg accepts the single title argument that "adr new" takes.  Loose words
// are rejected rather than joined, so that a title is always quoted.
//
// This is CLI validation: how many operands there are, and whether the shell
// was told to keep them together.  That a title must not be empty is a rule
// about the request, and lives in the slice as create.ErrNoTitle.
func titleArg(_ *cobra.Command, args []string) error {
	switch {
	case len(args) == 0:
		return errors.New("no title given")
	case len(args) > 1:
		return errors.New("quote the title so it is a single argument")
	}
	return nil
}
