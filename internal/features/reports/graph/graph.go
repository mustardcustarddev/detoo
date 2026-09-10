// Package graph draws the decision records and the links between them, in
// Mermaid or Graphviz format.
package graph

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/mustardcustard/decision-records/internal/domain/adr"
)

// Format is the language a graph is written in.
type Format int

const (
	// Mermaid is the zero Format, so a graph is drawn in Mermaid unless another
	// format is asked for.
	Mermaid Format = iota

	// Graphviz is the format of the dot tools.
	Graphviz
)

// ErrUnknownFormat reports a Format that is none of the above.
var ErrUnknownFormat = errors.New("unknown graph format")

// Query asks for a graph of the log.
type Query struct {
	Format Format

	// LinkPrefix and LinkExtension surround each record's name to give the
	// address of the document the node links to.
	LinkPrefix    string
	LinkExtension string
}

// Result is the graph, as a Mermaid flowchart or a Graphviz digraph.
//
// As with a table of contents, the document is the product: it is piped into
// mmdc or dot, so it has to be valid in its format whichever client asked for
// it.  That is a contract, not a choice of presentation, and so it is settled
// here.
type Result struct {
	Document string
}

// Store is the port this slice needs.
type Store interface {
	Records(ctx context.Context) ([]adr.Record, error)
	Lines(ctx context.Context, file string) ([]string, error)
}

// Handler draws the graph.
type Handler struct {
	Log Store
}

// node is a record, as the graph draws it.
type node struct {
	number int
	title  string
	url    string
}

// edge is a link from one record to another.
type edge struct {
	from, to int
	label    string
}

func (h Handler) Handle(ctx context.Context, query Query) (Result, error) {
	var draw func(nodes []node, edges []edge) string
	switch query.Format {
	case Mermaid:
		draw = mermaid
	case Graphviz:
		draw = graphviz
	default:
		return Result{}, ErrUnknownFormat
	}

	records, err := h.Log.Records(ctx)
	if err != nil {
		return Result{}, err
	}

	// A log with no records has nothing to draw, and an empty flowchart or
	// digraph is not a graph worth rendering.
	if len(records) == 0 {
		return Result{}, nil
	}

	var nodes []node
	var edges []edge
	for _, record := range records {
		lines, err := h.Log.Lines(ctx, record.File)
		if err != nil {
			return Result{}, err
		}
		name := strings.TrimSuffix(path.Base(record.File), ".md")
		nodes = append(nodes, node{
			number: record.Number,
			title:  adr.Title(lines),
			url:    query.LinkPrefix + name + query.LinkExtension,
		})
		for _, link := range adr.Links(lines) {
			// Reverse links duplicate a link recorded by the other record.
			if strings.HasSuffix(link.Type, " by") {
				continue
			}
			edges = append(edges, edge{from: record.Number, to: link.Target, label: link.Type})
		}
	}
	return Result{Document: draw(nodes, edges)}, nil
}

// graphviz draws a digraph.  Each record is chained to the one before it, so
// that the graph reads as a log.
func graphviz(nodes []node, edges []edge) string {
	var out strings.Builder
	fmt.Fprintln(&out, "digraph {")
	fmt.Fprintln(&out, "  node [shape=plaintext];")
	fmt.Fprintln(&out, "  subgraph {")
	for _, n := range nodes {
		fmt.Fprintf(&out, "    _%d [label=%q; URL=%q];\n", n.number, n.title, n.url)
		if n.number > 1 {
			fmt.Fprintf(&out, "    _%d -> _%d [style=\"dotted\", weight=1];\n", n.number-1, n.number)
		}
	}
	fmt.Fprintln(&out, "  }")
	for _, e := range edges {
		fmt.Fprintf(&out, "  _%d -> _%d [label=%q, weight=0]\n", e.from, e.to, e.label)
	}
	fmt.Fprintln(&out, "}")
	return out.String()
}

// mermaid draws a flowchart, laid out as the digraph is: each record chained to
// the one before it by a dotted arrow, and linked to its document.
func mermaid(nodes []node, edges []edge) string {
	var out strings.Builder
	fmt.Fprintln(&out, "flowchart TD")
	for _, n := range nodes {
		fmt.Fprintf(&out, "  adr%d[\"%s\"]\n", n.number, mermaidText(n.title))
		if n.number > 1 {
			fmt.Fprintf(&out, "  adr%d -.-> adr%d\n", n.number-1, n.number)
		}
	}
	for _, e := range edges {
		fmt.Fprintf(&out, "  adr%d -->|\"%s\"| adr%d\n", e.from, mermaidText(e.label), e.to)
	}
	for _, n := range nodes {
		fmt.Fprintf(&out, "  click adr%d \"%s\"\n", n.number, n.url)
	}
	return out.String()
}

// mermaidText makes text safe inside a quoted Mermaid label.  Mermaid has no
// backslash escapes, so a quote is written as an entity.
func mermaidText(text string) string {
	return strings.ReplaceAll(text, `"`, "#quot;")
}
