// Package ust defines the ust tool: the vocabulary of a log of user stories,
// and the templates and help the tool ships with.
//
// It holds data and nothing else.  The behaviour is shared by every record
// tool, and cmd/ust puts the two together.
package ust

import (
	"embed"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// Kind is the vocabulary of a user story log, which is kept in docs/userstory,
// with each story's file name beginning "STORY-".  Logs were kept in docs/sty
// before, and the tool was called sty.  A story that is split is
// replaced by each of the smaller stories it is split into.
var Kind = record.Kind{
	Name:        "ust",
	Dir:         "userstory",
	Prefix:      "STORY",
	Former:      []string{"sty"},
	Noun:        "user story",
	Plural:      "user stories",
	Heading:     "User Stories",
	Status:      "Draft",
	Retires:     []string{"Draft"},
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
