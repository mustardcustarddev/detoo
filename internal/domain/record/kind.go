package record

import "strings"

// Kind is the vocabulary of one kind of log: what its records are called, the
// status a new one is given, and the words for how one bears on another.
//
// A tool is a Kind and the files it ships with.  Everything that differs
// between adr and rfc is a value here or a file there, so that the slices, the
// adapters and the commands are written once and serve both.
type Kind struct {
	// Name is the tool's name, in lower case.  It names everything that
	// belongs to the tool, each by one convention applied where the thing is
	// used: the command, the NAME_DATE, NAME_TEMPLATE and NAME_PAGER
	// variables, the .NAME-dir marker file, the docs/NAME directory, and the
	// nodes of a graph.
	Name string

	// Noun and Plural are what one record, and several, are called in prose.
	Noun, Plural string

	// Heading titles a table of contents of the log.
	Heading string

	// Status is the status a new record is created with.  A record that is
	// replaced loses this line of its status.
	Status string

	// FirstRecord is the title the record that opens a log is filed under.
	FirstRecord string

	// Replaces and Updates are the words for each Relation.
	Replaces, Updates Words
}

// Words are what one Relation is called.
type Words struct {
	// Verb names the relation as an option does, and its first letter is the
	// option's shorthand: "supercede".
	Verb string

	// Forward is how the later record describes the relation: "Supercedes".
	Forward string

	// Reverse is how the earlier record describes it: "Superceded by".  It
	// ends in " by", which is how a graph tells a reverse link from the
	// forward one it repeats.
	Reverse string
}

// WordsFor returns the words for a relation.
func (k Kind) WordsFor(r Relation) Words {
	if r == Replaces {
		return k.Replaces
	}
	return k.Updates
}

// Variable returns the name of the environment variable that configures a
// setting of the tool: ADR_DATE, for adr's DATE.
func (k Kind) Variable(setting string) string {
	return strings.ToUpper(k.Name) + "_" + setting
}
