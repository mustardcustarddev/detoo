// Package pse defines the pse tool: the vocabulary of a log of decisions made
// with PIOSEE, and the templates and help the tool ships with.
//
// It holds data and nothing else.  The behaviour is shared by every record
// tool, and cmd/pse puts the two together.
package pse

import (
	"embed"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// Kind is the vocabulary of a PIOSEE log, which is kept in docs/piosee, with
// each record's file name beginning "PIOSEE-".  A decision whose evaluation
// shows it did not work is reversed by the decision that undoes it.
var Kind = record.Kind{
	Name:        "pse",
	Dir:         "piosee",
	Prefix:      "PIOSEE",
	Noun:        "PIOSEE record",
	Plural:      "PIOSEE records",
	Heading:     "PIOSEE Decisions",
	Status:      "Draft",
	Retires:     []string{"Draft"},
	FirstRecord: "make-decisions-with-piosee",
	Replaces:    record.Words{Verb: "reverse", Forward: "Reverses", Reverse: "Reversed by"},
	Updates:     record.Words{Verb: "update", Forward: "Updates", Reverse: "Updated by"},
}

// Files holds the tool's built-in templates, under templates/, and what its
// commands print about themselves, under help/.
//
// The template for a new record has a section for each step of PIOSEE, as
// listed in https://en.wikipedia.org/wiki/Pilot_decision_making
//
//go:embed templates help
var Files embed.FS
