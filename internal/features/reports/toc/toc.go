// Package toc writes a table of contents for the records of a log, in
// Markdown.
package toc

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// Query asks for a table of contents.
type Query struct {
	// Intro and Outro name files whose Markdown surrounds the list.
	Intro string
	Outro string

	// LinkPrefix precedes each record's file name in the links, so that the
	// contents can be published somewhere other than beside the records.
	LinkPrefix string
}

// Result is the table of contents.
//
// The document is the product of this slice, not a rendering of it.  Its
// Markdown is part of what the tool publishes — a caller redirects it into a
// file that other tools read — so composing it here is what lets an HTTP or MCP
// client ask for a table of contents and get the same one the CLI writes.  What
// the CLI still decides is where the text goes.
type Result struct {
	Markdown string
}

// Store is the port this slice needs.
type Store interface {
	Records(ctx context.Context) ([]record.Record, error)
	Lines(ctx context.Context, file string) ([]string, error)
	ReadFile(ctx context.Context, file string) (string, error)
}

// Handler writes a table of contents.
type Handler struct {
	Kind record.Kind
	Log  Store
}

func (h Handler) Handle(ctx context.Context, query Query) (Result, error) {
	records, err := h.Log.Records(ctx)
	if err != nil {
		return Result{}, err
	}

	// A log with no records has no contents, and a heading with nothing under
	// it would be published as though it were a document.  The intro and outro
	// are not read, because nothing is written for them to surround.
	if len(records) == 0 {
		return Result{}, nil
	}

	var out strings.Builder
	fmt.Fprintf(&out, "# %s\n", h.Kind.Heading)
	fmt.Fprintln(&out)

	if query.Intro != "" {
		text, err := h.Log.ReadFile(ctx, query.Intro)
		if err != nil {
			return Result{}, err
		}
		fmt.Fprintf(&out, "%s\n", strings.TrimSuffix(text, "\n"))
		fmt.Fprintln(&out)
	}

	for _, r := range records {
		lines, err := h.Log.Lines(ctx, r.File)
		if err != nil {
			return Result{}, err
		}
		fmt.Fprintf(&out, "* [%s](%s)\n", record.Title(lines), query.LinkPrefix+path.Base(r.File))
	}

	if query.Outro != "" {
		text, err := h.Log.ReadFile(ctx, query.Outro)
		if err != nil {
			return Result{}, err
		}
		fmt.Fprintln(&out)
		fmt.Fprintf(&out, "%s\n", strings.TrimSuffix(text, "\n"))
	}
	return Result{Markdown: out.String()}, nil
}
