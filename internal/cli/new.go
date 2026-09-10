package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mustardcustard/decision-records/internal/app"
	"github.com/mustardcustard/decision-records/internal/cli/output"
	"github.com/mustardcustard/decision-records/internal/features/records/create"
)

func newNewCommand(application *app.Application) *cobra.Command {
	var supercedes, links []string

	cmd := &cobra.Command{
		Use:     "new TITLE",
		Short:   "Create the next numbered decision record",
		Long:    helpText("new.md"),
		Example: helpText("new.example"),
		Args:    titleArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			parsed, err := parseLinks(links)
			if err != nil {
				return err
			}
			result, err := application.Records.Create.Handle(cmd.Context(), create.Command{
				Title:      args[0],
				Supercedes: supercedes,
				Links:      parsed,
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
	cmd.Flags().StringArrayVarP(&links, "link", "l", nil,
		"a link to a previous decision, as `TARGET:LINK:REVERSE-LINK`")
	return cmd
}

// titleArg accepts the single title argument that "adr new" takes.  Loose words
// are rejected rather than joined, so that a title is quoted the way a link
// description already must be.
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

// parseLinks splits the "TARGET:LINK:REVERSE-LINK" specifications the --link
// option takes.
//
// The colons are an encoding for a flag, so they are decoded here.  What
// reaches the slice is a list of links whose fields are named, and which says
// nothing about how the user typed them.
func parseLinks(specs []string) ([]create.Link, error) {
	links := make([]create.Link, 0, len(specs))
	for _, spec := range specs {
		parts := strings.Split(spec, ":")
		if len(parts) != 3 {
			return nil, fmt.Errorf("link %q is not of the form TARGET:LINK:REVERSE-LINK", spec)
		}
		links = append(links, create.Link{Target: parts[0], Forward: parts[1], Reverse: parts[2]})
	}
	return links, nil
}
