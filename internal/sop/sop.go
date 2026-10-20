// Package sop defines the sop tool: the vocabulary of a log of standard
// operating procedures, each the steps of a routine task, such as rolling
// back a release, and the templates and help the tool ships with.
//
// It holds data and nothing else.  The behaviour is shared by every record
// tool, and cmd/sop puts the two together.
package sop

import (
	"embed"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// Kind is the vocabulary of a log of standard operating procedures, which is
// kept in docs/sop, with each procedure's file name beginning "SOP-".  A
// procedure rewritten for the same task supersedes the one it rewrites, which
// is then neither a draft nor approved.  The log opens with a README rather
// than a first procedure, which would be numbered as one of the team's
// procedures without being one.
var Kind = record.Kind{
	Name:     "sop",
	Dir:      "sop",
	Prefix:   "SOP",
	Noun:     "SOP",
	Plural:   "SOPs",
	Heading:  "Standard Operating Procedures",
	Status:   "Draft",
	Retires:  []string{"Draft", "Approved"},
	Replaces: record.Words{Verb: "supersede", Forward: "Supersedes", Reverse: "Superseded by"},
	Updates:  record.Words{Verb: "update", Forward: "Updates", Reverse: "Updated by"},
}

// Files holds the tool's built-in templates, under templates/, and what its
// commands print about themselves, under help/.
//
// The template for a new procedure is a title, a date, a status and its
// steps, as an ordered list.
//
//go:embed templates help
var Files embed.FS
