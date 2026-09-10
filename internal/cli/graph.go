package cli

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/mustardcustard/decision-records/internal/app"
	"github.com/mustardcustard/decision-records/internal/cli/output"
	"github.com/mustardcustard/decision-records/internal/features/reports/graph"
)

func newGraphCommand(application *app.Application) *cobra.Command {
	var graphviz, mermaid bool
	var linkPrefix, linkExtension string

	cmd := &cobra.Command{
		Use:     "graph",
		Short:   "Draw the links between records in Mermaid or Graphviz format",
		Long:    helpText("graph.md"),
		Example: helpText("graph.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			format, err := formatFlag(graphviz, mermaid)
			if err != nil {
				return err
			}
			result, err := application.Reports.Graph.Handle(cmd.Context(), graph.Query{
				Format:        format,
				LinkPrefix:    linkPrefix,
				LinkExtension: linkExtension,
			})
			if err != nil {
				return err
			}
			return output.Document(cmd.OutOrStdout(), result.Document)
		},
	}

	cmd.Flags().BoolVar(&mermaid, "mermaid", false, "draw the graph in Mermaid format, the default")
	cmd.Flags().BoolVar(&graphviz, "graphviz", false, "draw the graph in Graphviz format")
	cmd.Flags().StringVarP(&linkExtension, "extension", "e", ".html",
		"the file `EXTENSION` of the documents the generated links refer to")
	cmd.Flags().StringVarP(&linkPrefix, "prefix", "p", "",
		"prefix each decision file link with `LINK_PREFIX`")
	return cmd
}

// formatFlag returns the format named by the --graphviz and --mermaid switches.
// A graph is drawn in Mermaid unless --graphviz is given.
//
// This is CLI validation: which switches were typed.  That a format is one the
// graph can be drawn in is a rule about the request, and lives in the slice as
// graph.ErrUnknownFormat.
func formatFlag(graphviz, mermaid bool) (graph.Format, error) {
	switch {
	case graphviz && mermaid:
		return graph.Mermaid, errors.New("give either --graphviz or --mermaid, not both")
	case graphviz:
		return graph.Graphviz, nil
	}
	return graph.Mermaid, nil
}
