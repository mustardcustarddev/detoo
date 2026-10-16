// Package adr defines the adr tool: the vocabulary of a log of architecture
// decision records, and the templates and help the tool ships with.
//
// It holds data and nothing else.  The behaviour is shared by every record
// tool, and cmd/adr puts the two together.
package adr

import (
	"embed"

	"github.com/mustardcustard/decision-records/internal/domain/record"
)

// Kind is the vocabulary of an architecture decision record log.
var Kind = record.Kind{
	Name:        "adr",
	Noun:        "decision record",
	Plural:      "decision records",
	Heading:     "Architecture Decision Records",
	Status:      "Accepted",
	FirstRecord: "record-architecture-decisions",
	Replaces:    record.Words{Verb: "supercede", Forward: "Supercedes", Reverse: "Superceded by"},
	Updates:     record.Words{Verb: "update", Forward: "Updates", Reverse: "Updated by"},
}

// Files holds the tool's built-in templates, under templates/, and what its
// commands print about themselves, under help/.
//
// The template for a new record follows the style described by Michael Nygard
// in http://thinkrelevance.com/blog/2011/11/15/documenting-architecture-decisions
//
//go:embed templates help
var Files embed.FS
