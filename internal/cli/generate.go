package cli

import (
	"github.com/spf13/cobra"

	"github.com/mustardcustard/decision-records/internal/app"
	"github.com/mustardcustard/decision-records/internal/cli/output"
	"github.com/mustardcustard/decision-records/internal/features/reports/graph"
	"github.com/mustardcustard/decision-records/internal/features/reports/toc"
)

func newGenerateCommand(application *app.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate [REPORT]",
		Short: "Write a report about the decision records",
		Long:  helpText("generate.md"),
		Args:  cobra.NoArgs,

		// Without a report, generate lists the reports it can write, which is
		// what the shell implementation did.
		RunE: func(cmd *cobra.Command, _ []string) error {
			for _, report := range cmd.Commands() {
				if report.IsAvailableCommand() {
					output.Path(cmd.OutOrStdout(), report.Name())
				}
			}
			return nil
		},
	}
	cmd.AddCommand(newGraphCommand(application), newTOCCommand(application))
	return cmd
}

func newGraphCommand(application *app.Application) *cobra.Command {
	var linkPrefix, linkExtension string

	cmd := &cobra.Command{
		Use:     "graph",
		Short:   "Write the links between records in Graphviz format",
		Long:    helpText("graph.md"),
		Example: helpText("graph.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := application.Reports.Graph.Handle(cmd.Context(), graph.Query{
				LinkPrefix:    linkPrefix,
				LinkExtension: linkExtension,
			})
			if err != nil {
				return err
			}
			return output.Document(cmd.OutOrStdout(), result.Graphviz)
		},
	}

	cmd.Flags().StringVarP(&linkExtension, "extension", "e", ".html",
		"the file `EXTENSION` of the documents the generated links refer to")
	cmd.Flags().StringVarP(&linkPrefix, "prefix", "p", "",
		"prefix each decision file link with `LINK_PREFIX`")
	return cmd
}

func newTOCCommand(application *app.Application) *cobra.Command {
	var intro, outro, linkPrefix string

	cmd := &cobra.Command{
		Use:   "toc",
		Short: "Write a table of contents in Markdown",
		Long:  helpText("toc.md"),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := application.Reports.TOC.Handle(cmd.Context(), toc.Query{
				Intro:      intro,
				Outro:      outro,
				LinkPrefix: linkPrefix,
			})
			if err != nil {
				return err
			}
			return output.Document(cmd.OutOrStdout(), result.Markdown)
		},
	}

	cmd.Flags().StringVarP(&intro, "intro", "i", "",
		"precede the table of contents with the Markdown in `INTRO`")
	cmd.Flags().StringVarP(&outro, "outro", "o", "",
		"follow the table of contents with the Markdown in `OUTRO`")
	cmd.Flags().StringVarP(&linkPrefix, "prefix", "p", "",
		"prefix each decision file link with `LINK_PREFIX`")
	return cmd
}
