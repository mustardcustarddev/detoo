package cli

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/mustardcustard/decision-records/internal/app"
	"github.com/mustardcustard/decision-records/internal/cli/output"
	"github.com/mustardcustard/decision-records/internal/features/records/list"
	"github.com/mustardcustard/decision-records/internal/features/reports/toc"
)

func newListCommand(t tool, application *app.Application) *cobra.Command {
	var contents bool
	var intro, outro, linkPrefix string

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List the " + t.Plural + ", or write a table of contents",
		Long:    t.help("list.md"),
		Example: t.help("list.example"),
		Args:    cobra.NoArgs,

		// One command, two slices: the records are either listed for a script,
		// one path a line, or published as a Markdown document.
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := tocOptions(contents, intro, outro, linkPrefix); err != nil {
				return err
			}
			if contents {
				result, err := application.Reports.TOC.Handle(cmd.Context(), toc.Query{
					Intro:      intro,
					Outro:      outro,
					LinkPrefix: linkPrefix,
				})
				if err != nil {
					return err
				}
				return output.Document(cmd.OutOrStdout(), result.Markdown)
			}

			result, err := application.Records.List.Handle(cmd.Context(), list.Query{})
			if err != nil {
				return err
			}
			output.Records(cmd.OutOrStdout(), result.Records)
			return nil
		},
	}

	cmd.Flags().BoolVar(&contents, "toc", false, "write a table of contents in Markdown")
	cmd.Flags().StringVarP(&intro, "intro", "i", "",
		"precede the table of contents with the Markdown in `INTRO`")
	cmd.Flags().StringVarP(&outro, "outro", "o", "",
		"follow the table of contents with the Markdown in `OUTRO`")
	cmd.Flags().StringVarP(&linkPrefix, "prefix", "p", "",
		"prefix each link to a record with `LINK_PREFIX`")
	return cmd
}

// tocOptions rejects the options that shape a table of contents when --toc was
// not given, rather than silently ignoring them.
func tocOptions(contents bool, intro, outro, linkPrefix string) error {
	if !contents && (intro != "" || outro != "" || linkPrefix != "") {
		return errors.New("give --toc to use --intro, --outro or --prefix")
	}
	return nil
}
