package cli

import (
	"github.com/spf13/cobra"

	"github.com/mustardcustard/decision-records/internal/app"
	"github.com/mustardcustard/decision-records/internal/features/records/link"
)

func newLinkCommand(application *app.Application) *cobra.Command {
	return &cobra.Command{
		Use:     "link SOURCE LINK TARGET REVERSE-LINK",
		Short:   "Link two decision records to each other",
		Long:    helpText("link.md"),
		Example: helpText("link.example"),

		// Four operands exactly, which is also what forces a link description
		// containing spaces to be quoted.
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := application.Records.Link.Handle(cmd.Context(), link.Command{
				Source:  args[0],
				Forward: args[1],
				Target:  args[2],
				Reverse: args[3],
			})
			return err
		},
	}
}
