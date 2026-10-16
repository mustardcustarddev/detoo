package cli

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mustardcustard/detoo/internal/app"
	"github.com/mustardcustard/detoo/internal/cli/output"
	"github.com/mustardcustard/detoo/internal/features/records/create"
)

func newNewCommand(t tool, application *app.Application) *cobra.Command {
	var replaces, updates []string

	cmd := &cobra.Command{
		Use:     "new TITLE",
		Short:   "Create the next numbered " + t.Noun,
		Long:    t.help("new.md"),
		Example: t.help("new.example"),
		Args:    titleArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := application.Records.Create.Handle(cmd.Context(), create.Command{
				Title:    args[0],
				Replaces: replaces,
				Updates:  updates,
			})
			if err != nil {
				return err
			}
			output.Path(cmd.OutOrStdout(), result.File)
			return nil
		},
	}

	// Each option is named with the tool's verb for its relation, such as
	// --supercede, and abbreviated to the verb's first letter.
	cmd.Flags().StringArrayVarP(&replaces, t.Replaces.Verb, t.Replaces.Verb[:1], nil,
		"a `REFERENCE` (number or partial filename) to a record this one "+strings.ToLower(t.Replaces.Forward))
	cmd.Flags().StringArrayVarP(&updates, t.Updates.Verb, t.Updates.Verb[:1], nil,
		"a `REFERENCE` (number or partial filename) to a record this one "+strings.ToLower(t.Updates.Forward))
	return cmd
}

// titleArg accepts the single title argument that "new" takes.  Loose words are
// rejected rather than joined, so that a title is always quoted.
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
