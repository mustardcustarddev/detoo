// Package rfc defines the rfc tool: the vocabulary of a log of requests for
// comments, and the templates and help the tool ships with.
//
// It holds data and nothing else.  The behaviour is shared by every record
// tool, and cmd/rfc puts the two together.
package rfc

import (
	"embed"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// Kind is the vocabulary of a request for comments log.  The words for its
// relations are the ones the IETF uses between RFCs.
var Kind = record.Kind{
	Name:        "rfc",
	Noun:        "RFC",
	Plural:      "RFCs",
	Heading:     "Requests for Comments",
	Status:      "Draft",
	FirstRecord: "use-requests-for-comments",
	Replaces:    record.Words{Verb: "obsolete", Forward: "Obsoletes", Reverse: "Obsoleted by"},
	Updates:     record.Words{Verb: "update", Forward: "Updates", Reverse: "Updated by"},
}

// Files holds the tool's built-in templates, under templates/, and what its
// commands print about themselves, under help/.
//
//go:embed templates help
var Files embed.FS
