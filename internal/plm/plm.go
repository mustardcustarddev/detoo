// Package plm defines the plm tool: the vocabulary of a log of problem
// playbooks, each a plan for responding to a problem, such as an incident,
// when it arises, and the templates and help the tool ships with.
//
// It holds data and nothing else.  The behaviour is shared by every record
// tool, and cmd/plm puts the two together.
package plm

import (
	"embed"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// Kind is the vocabulary of a problem log, which is kept in docs/problem, with
// each playbook's file name beginning "PROBLEM-".  Logs were kept in
// docs/piosee before, by the pse tool, with each file name beginning
// "PIOSEE-".  A playbook rewritten for the same problem, such as after it was
// used and found wanting, replaces the one it rewrites.
var Kind = record.Kind{
	Name:           "plm",
	Dir:            "problem",
	Prefix:         "PROBLEM",
	FormerPrefixes: []string{"PIOSEE"},
	Former:         []string{"piosee"},
	Noun:           "playbook",
	Plural:         "playbooks",
	Heading:        "Problem Playbooks",
	Status:         "Draft",
	Retires:        []string{"Draft"},
	FirstRecord:    "plan-for-incidents-with-piosee",
	Replaces:       record.Words{Verb: "replace", Forward: "Replaces", Reverse: "Replaced by"},
	Updates:        record.Words{Verb: "update", Forward: "Updates", Reverse: "Updated by"},
}

// Files holds the tool's built-in templates, under templates/, and what its
// commands print about themselves, under help/.
//
// The template for a new playbook follows PIOSEE: a section for the scenario
// it plans for, the Problem, the Immediate action to take before anything
// else, the Counter-indicated actions that would make things worse, and then
// Information, Options, Selection, Execution and Evaluation.
//
//go:embed templates help
var Files embed.FS
