// Package pse defines the pse tool: the vocabulary of a log of PIOSEE
// playbooks, each a plan for deciding what to do when a foreseen situation
// arises, and the templates and help the tool ships with.
//
// It holds data and nothing else.  The behaviour is shared by every record
// tool, and cmd/pse puts the two together.
package pse

import (
	"embed"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// Kind is the vocabulary of a PIOSEE log, which is kept in docs/piosee, with
// each playbook's file name beginning "PIOSEE-".  A playbook rewritten for the
// same situation, such as after it was used and found wanting, replaces the
// one it rewrites.
var Kind = record.Kind{
	Name:        "pse",
	Dir:         "piosee",
	Prefix:      "PIOSEE",
	Noun:        "PIOSEE playbook",
	Plural:      "PIOSEE playbooks",
	Heading:     "PIOSEE Playbooks",
	Status:      "Draft",
	Retires:     []string{"Draft"},
	FirstRecord: "plan-for-incidents-with-piosee",
	Replaces:    record.Words{Verb: "replace", Forward: "Replaces", Reverse: "Replaced by"},
	Updates:     record.Words{Verb: "update", Forward: "Updates", Reverse: "Updated by"},
}

// Files holds the tool's built-in templates, under templates/, and what its
// commands print about themselves, under help/.
//
// The template for a new playbook has a section for the scenario it plans
// for, and one for each step of PIOSEE: Problem, Information, Options,
// Selection, Execution and Evaluation.
//
//go:embed templates help
var Files embed.FS
