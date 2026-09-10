// Package graph writes the links between the decision records in Graphviz
// format.
package graph

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/mustardcustard/decision-records/internal/domain/adr"
)

// Query asks for a graph of the log.
type Query struct {
	// LinkPrefix and LinkExtension surround each record's name to give the
	// address of the document the node links to.
	LinkPrefix    string
	LinkExtension string
}

// Result is the graph, as a Graphviz digraph.
//
// As with a table of contents, the document is the product: it is piped into
// dot, so it has to be valid Graphviz whichever client asked for it.  That is a
// contract, not a choice of presentation, and so it is settled here.
type Result struct {
	Graphviz string
}

// Store is the port this slice needs.
type Store interface {
	Records(ctx context.Context) ([]adr.Record, error)
	Lines(ctx context.Context, file string) ([]string, error)
}

// Handler writes the graph.
type Handler struct {
	Log Store
}

func (h Handler) Handle(ctx context.Context, query Query) (Result, error) {
	records, err := h.Log.Records(ctx)
	if err != nil {
		return Result{}, err
	}

	var out strings.Builder
	fmt.Fprintln(&out, "digraph {")
	fmt.Fprintln(&out, "  node [shape=plaintext];")
	fmt.Fprintln(&out, "  subgraph {")

	// One node a record, chained in order, so that the graph reads as a log.
	for _, record := range records {
		lines, err := h.Log.Lines(ctx, record.File)
		if err != nil {
			return Result{}, err
		}
		name := strings.TrimSuffix(path.Base(record.File), ".md")
		fmt.Fprintf(&out, "    _%d [label=%q; URL=%q];\n",
			record.Number, adr.Title(lines), query.LinkPrefix+name+query.LinkExtension)
		if record.Number > 1 {
			fmt.Fprintf(&out, "    _%d -> _%d [style=\"dotted\", weight=1];\n", record.Number-1, record.Number)
		}
	}
	fmt.Fprintln(&out, "  }")

	for _, record := range records {
		lines, err := h.Log.Lines(ctx, record.File)
		if err != nil {
			return Result{}, err
		}
		for _, link := range adr.Links(lines) {
			// Reverse links duplicate a link recorded by the other record.
			if strings.HasSuffix(link.Type, " by") {
				continue
			}
			fmt.Fprintf(&out, "  _%d -> _%d [label=%q, weight=0]\n", record.Number, link.Target, link.Type)
		}
	}
	fmt.Fprintln(&out, "}")
	return Result{Graphviz: out.String()}, nil
}
