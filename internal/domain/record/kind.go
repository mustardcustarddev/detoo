package record

import "strings"

// Kind is the vocabulary of one kind of log: what its records are called, the
// status a new one is given, and the words for how one bears on another.
//
// A tool is a Kind and the files it ships with.  Everything that differs
// between the tools is a value here or a file there, so that the slices, the
// adapters and the commands are written once and serve them all.
type Kind struct {
	// Name is the tool's name, in lower case.  It names everything that
	// belongs to the tool, each by one convention applied where the thing is
	// used: the command, the NAME_DATE, NAME_TEMPLATE and NAME_PAGER
	// variables, and the nodes of a graph.
	Name string

	// Dir names the directory, under docs, that a project keeps its log in:
	// "adr", for docs/adr.  It need not be the tool's name.
	Dir string

	// Prefix begins the file name of every record, in upper case, followed by
	// a hyphen and the record's number: "STORY", for STORY-0001-reorder.md.
	//
	// Records named before the prefixes were upper case, such as
	// story-0001-reorder.md, or before there was a prefix, such as
	// 0001-reorder.md, are still records of the kind: see CanonicalName.
	Prefix string

	// FormerPrefixes are the prefixes the log's records began with before
	// Prefix: "PIOSEE", for PIOSEE-0001-outage.md.  A record named with one
	// is still a record of the kind, and its canonical name has Prefix.
	FormerPrefixes []string

	// Former names the directories, under docs, that the log was kept in by
	// default before Dir: "sty", for docs/sty.  A tool no longer looks there.
	Former []string

	// Noun and Plural are what one record, and several, are called in prose.
	Noun, Plural string

	// Heading titles a table of contents of the log.
	Heading string

	// Status is the status a new record is created with.
	Status string

	// Retires are the statuses a record loses when another replaces it,
	// those that say it is in force or about to be: a superseded decision is
	// neither a draft nor accepted.  Any other status, such as Deprecated,
	// is the author's to change.
	Retires []string

	// FirstRecord is the title the record that opens a log is filed under.
	FirstRecord string

	// Replaces and Updates are the words for each Relation.
	Replaces, Updates Words
}

// Words are what one Relation is called.
type Words struct {
	// Verb names the relation as an option does, and its first letter is the
	// option's shorthand: "supersede".
	Verb string

	// Forward is how the later record describes the relation: "Supersedes".
	Forward string

	// Reverse is how the earlier record describes it: "Superseded by".  It
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
