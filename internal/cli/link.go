package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mustardcustard/detoo/internal/app"
	"github.com/mustardcustard/detoo/internal/domain/record"
	"github.com/mustardcustard/detoo/internal/features/records/link"
)

func newLinkCommand(t tool, application *app.Application) *cobra.Command {
	var replace, update bool

	cmd := &cobra.Command{
		Use: fmt.Sprintf("link (--%s | --%s) SOURCE TARGET", t.Replaces.Verb, t.Updates.Verb),
		Short: fmt.Sprintf("Record that one %s %s or %s another",
			t.Noun, strings.ToLower(t.Replaces.Forward), strings.ToLower(t.Updates.Forward)),
		Long:    t.help("link.md"),
		Example: t.help("link.example"),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			relation, err := relationFlag(t.Kind, replace, update)
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

	// Each switch is named with the tool's verb for its relation, such as
	// --supersede, and abbreviated to the verb's first letter.
	cmd.Flags().BoolVarP(&replace, t.Replaces.Verb, t.Replaces.Verb[:1], false,
		fmt.Sprintf("SOURCE %s TARGET, which loses its %s status", strings.ToLower(t.Replaces.Forward), strings.Join(t.Retires, " or ")))
	cmd.Flags().BoolVarP(&update, t.Updates.Verb, t.Updates.Verb[:1], false,
		"SOURCE makes substantive changes to TARGET")
	return cmd
}

// relationFlag returns the relation named by the switches for replacing and
// updating a record, such as --supersede and --update, exactly one of which
// must be set.
//
// This is CLI validation: which switches were typed.  That a link must have a
// relation is a rule about the request, and lives in the slice as
// link.ErrNoRelation.
func relationFlag(kind record.Kind, replace, update bool) (record.Relation, error) {
	switch {
	case replace == update:
		return 0, fmt.Errorf("give either --%s or --%s", kind.Replaces.Verb, kind.Updates.Verb)
	case replace:
		return record.Replaces, nil
	}
	return record.Updates, nil
}
