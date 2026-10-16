// Package sty defines the sty tool: the vocabulary of a log of user stories,
// and the templates and help the tool ships with.
//
// It holds data and nothing else.  The behaviour is shared by every record
// tool, and cmd/sty puts the two together.
package sty

import (
	"embed"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// Kind is the vocabulary of a user story log.  A story that is split is
// replaced by each of the smaller stories it is split into.
var Kind = record.Kind{
	Name:        "sty",
	Noun:        "user story",
	Plural:      "user stories",
	Heading:     "User Stories",
	Status:      "Draft",
	FirstRecord: "describe-work-in-user-stories",
	Replaces:    record.Words{Verb: "replace", Forward: "Replaces", Reverse: "Replaced by"},
	Updates:     record.Words{Verb: "update", Forward: "Updates", Reverse: "Updated by"},
}

// Files holds the tool's built-in templates, under templates/, and what its
// commands print about themselves, under help/.
//
// The template for a new story follows the form described by Mike Cohn in
// https://www.mountaingoatsoftware.com/agile/user-stories
//
//go:embed templates help
var Files embed.FS
