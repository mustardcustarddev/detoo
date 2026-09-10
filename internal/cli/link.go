package cli

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/mustardcustard/decision-records/internal/app"
	"github.com/mustardcustard/decision-records/internal/domain/adr"
	"github.com/mustardcustard/decision-records/internal/features/records/link"
)

func newLinkCommand(application *app.Application) *cobra.Command {
	var supercede, update bool

	cmd := &cobra.Command{
		Use:     "link (--supercede | --update) SOURCE TARGET",
		Short:   "Record that one decision record supercedes or updates another",
		Long:    helpText("link.md"),
		Example: helpText("link.example"),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			relation, err := relationFlag(supercede, update)
			if err != nil {
				return err
			}
			_, err = application.Records.Link.Handle(cmd.Context(), link.Command{
				Relation: relation,
				Source:   args[0],
				Target:   args[1],
			})
			return err
		},
	}

	cmd.Flags().BoolVarP(&supercede, "supercede", "s", false, "SOURCE supercedes TARGET, which is no longer accepted")
	cmd.Flags().BoolVarP(&update, "update", "u", false, "SOURCE makes substantive changes to TARGET")
	return cmd
}

// relationFlag returns the relation named by the --supercede and --update
// switches, exactly one of which must be set.
//
// This is CLI validation: which switches were typed.  That a link must have a
// relation is a rule about the request, and lives in the slice as
// link.ErrNoRelation.
func relationFlag(supercede, update bool) (adr.Relation, error) {
	switch {
	case supercede == update:
		return 0, errors.New("give either --supercede or --update")
	case supercede:
		return adr.Supercedes, nil
	}
	return adr.Updates, nil
}
