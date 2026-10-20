// Package fyr defines the fyr tool: the vocabulary of a log of fire drills,
// each an exercise in which a team practises a problem playbook on a local
// copy of the system, and the templates and help the tool ships with.
//
// It holds data and nothing else.  The behaviour is shared by every record
// tool, and cmd/fyr puts the two together.
package fyr

import (
	"embed"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// Kind is the vocabulary of a log of fire drills, which is kept in
// docs/firedrill, with each drill's file name beginning "FIREDRILL-".  A drill
// rewritten for the same case replaces the one it rewrites, which is then no
// longer a draft.  The log opens with a README rather than a first drill,
// since a drill that sets out how to drill is not one.
var Kind = record.Kind{
	Name:     "fyr",
	Dir:      "firedrill",
	Prefix:   "FIREDRILL",
	Noun:     "fire drill",
	Plural:   "fire drills",
	Heading:  "Fire Drills",
	Status:   "Draft",
	Retires:  []string{"Draft"},
	Replaces: record.Words{Verb: "replace", Forward: "Replaces", Reverse: "Replaced by"},
	Updates:  record.Words{Verb: "update", Forward: "Updates", Reverse: "Updated by"},
}

// Files holds the tool's built-in templates, under templates/, and what its
// commands print about themselves, under help/.
//
// The template for a new drill is a title, a date, a status, and a heading
// for each part of the drill, from the playbook it practises to the debrief.
//
//go:embed templates help
var Files embed.FS
