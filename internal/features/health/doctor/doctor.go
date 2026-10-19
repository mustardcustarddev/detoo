// Package doctor checks that a project's logs are where the tools look for
// them, and named as the tools name them today, and fixes them.
//
// Handler reports and changes nothing.  A log the tools have lost after a
// rename, such as user stories left in docs/sty, is found here, and so is a
// document named before its log's prefix was upper case, or with a slip of
// it, two documents that share a number, and a document whose title does not
// match its file name.  When moving or renaming would clash with something
// already there, the problem is reported as needing fixing by hand, and so
// are documents that share a number when there is no git history to say which
// came first, and every title that does not match.
//
// Fixer moves and renames what Handler finds, renumbers the documents that
// share a number when asked to, and updates the links to what it moved,
// leaving the clashes alone.
package doctor

import (
	"context"
	"path"
	"time"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// Query asks for the health of every log of the given kinds.
type Query struct{}

// Project is the port this slice needs: a view of the project's directory.
// Every path it takes and returns is slash separated and relative to the
// project's root.
type Project interface {
	// Root returns the project's directory.
	Root(ctx context.Context) (string, error)

	// Files returns the names of the files in dir, and whether dir exists.
	Files(ctx context.Context, dir string) ([]string, bool, error)

	// Read returns the content of a file.
	Read(ctx context.Context, name string) (string, error)
}

// History is the port that orders documents that share a number: when git
// first recorded each.
type History interface {
	// Repository reports whether the project is in a git repository.
	Repository(ctx context.Context) (bool, error)

	// Added returns when the file name was first committed, following it
	// through renames, and whether it has been.  A file git does not track,
	// or has not committed yet, has not.
	Added(ctx context.Context, name string) (time.Time, bool, error)
}

// Handler checks the logs of Kinds, ordering documents that share a number
// by History.  Without History, as outside a git repository, they need
// renumbering by hand.
type Handler struct {
	Kinds   []record.Kind
	Project Project
	History History
}

// Result is the health of each log, in the order of the kinds.
type Result struct {
	// Root is the project's directory, which every path is relative to.
	Root string

	Logs []Log
}

// Log is the health of one log.
type Log struct {
	Kind record.Kind

	// Dir is where the tool looks for the log: docs/<Kind.Dir>.
	Dir string

	// Dirs are the directories that hold the log's documents, Dir first if it
	// exists, then each in Misplaced.
	Dirs []Dir

	// Misplaced are the directories the log is kept in that the tool no
	// longer looks in, or should not.
	Misplaced []Move

	// Misnamed are the documents not named with the log's prefix.
	Misnamed []Rename

	// Duplicates are the numbers that more than one document in a directory
	// of the log has.
	Duplicates []Duplicate

	// Mistitled are the documents whose titles do not match their file
	// names.
	Mistitled []Mistitle

	// Renumbers give each document a number of its own: each document that
	// shares a number but the first, and each that moves up to make room for
	// one.  When the duplicates need fixing by hand, they are only a
	// suggestion.
	Renumbers []Renumber
}

// Dir is a directory that holds documents of a log.
type Dir struct {
	Path string

	// Documents are the file names of the documents in the directory, in the
	// order the directory lists them.
	Documents []string
}

// Reason is why a directory is in the wrong place.
type Reason int

const (
	// UnderDoc is a log kept under doc/ rather than docs/.
	UnderDoc Reason = iota

	// Former is a log kept in a directory the tool used to look in, such as
	// docs/sty.
	Former
)

// Clash is why a problem cannot be fixed by moving or renaming.
type Clash int

const (
	// NoClash is a problem that a move or a rename fixes.
	NoClash Clash = iota

	// TargetExists is a directory whose target is there already, such as
	// docs/sty when docs/userstory exists.
	TargetExists

	// SharedTarget is a directory whose target another misplaced directory
	// moves to as well.
	SharedTarget

	// NoHistory is documents that share a number in a project that is not in
	// a git repository, so that which came first is not known.
	NoHistory

	// NamesClash is documents that share a number in a directory where a
	// rename clashes, which is to be fixed first.
	NamesClash
)

// Move is a log kept in the wrong directory, and where it belongs.
type Move struct {
	From, To  string
	Reason    Reason
	Documents int
	Clash     Clash

	// ClashWith is the other directory moving to To, when Clash is
	// SharedTarget.
	ClashWith string
}

// Rename is a document not named with its log's prefix, and the name it
// should have.  Clash reports that a document of that name is there already,
// or that another would be renamed to it too: ClashWith names that document.
type Rename struct {
	Dir, From, To string
	Clash         bool
	ClashWith     string
}

// Duplicate is a number that more than one document in a directory of a log
// has.
type Duplicate struct {
	Dir    string
	Number int

	// Documents are the names of the documents numbered Number, as Misnamed
	// renames them: the one that keeps the number first, then the others in
	// the order git first recorded them.  Without that order, they are in
	// the order of their names.
	Documents []string

	// Clash is why the documents need renumbering by hand, or NoClash.
	Clash Clash
}

// Renumber is a document given a number of its own, or moved up to make room
// for one that is.
type Renumber struct {
	// Dir is the directory the document is in, and File its name there.
	Dir, File string

	// From is the document's name as Misnamed renames it, and To its name
	// with its new number.
	From, To string

	// Duplicate reports that the document shares its number with another
	// that came first.  Otherwise it moves up to make room for one that did.
	Duplicate bool

	// After and Before are the documents a duplicate goes between, by their
	// names once renumbered: the last added before it, and the one after
	// that in the log, if any.
	After, Before string

	// Clash is why the document needs renumbering by hand, or NoClash.
	Clash Clash
}

// Mistitle is a document whose title does not match its file name, as the
// tools name a document from its number and title.  Which of the two is
// wrong is not known, so it is fixed by hand.
type Mistitle struct {
	// Dir is the directory the document is in, and File its name there.
	Dir, File string

	// Title is the document's title, from its first line, without the "# ".
	Title string

	// Number reports that the title is not numbered as the file is, or not
	// numbered at all.
	Number bool

	// Words reports that the rest of the title is not the rest of the file
	// name, even with its case and punctuation set aside.
	Words bool
}

// Found reports whether the project keeps the log anywhere.
func (l Log) Found() bool { return len(l.Dirs) > 0 }

// Problems counts what is wrong with the log: each misplaced directory, each
// misnamed document, each number that documents share, and each mistitled
// document.
func (l Log) Problems() int {
	return len(l.Misplaced) + len(l.Misnamed) + len(l.Duplicates) + len(l.Mistitled)
}

// ByHand counts the problems that need fixing by hand.
func (l Log) ByHand() int {
	n := 0
	for _, m := range l.Misplaced {
		if m.Clash != NoClash {
			n++
		}
	}
	for _, r := range l.Misnamed {
		if r.Clash {
			n++
		}
	}
	for _, d := range l.Duplicates {
		if d.Clash != NoClash {
			n++
		}
	}
	return n + len(l.Mistitled)
}

// Renumberable counts the documents a fix would renumber, if asked to.
func (l Log) Renumberable() int {
	n := 0
	for _, r := range l.Renumbers {
		if r.Clash == NoClash {
			n++
		}
	}
	return n
}

// Documents counts the log's documents, wherever they are kept.
func (l Log) Documents() int {
	n := 0
	for _, d := range l.Dirs {
		n += len(d.Documents)
	}
	return n
}

// DocumentsIn counts the log's documents in dir.
func (l Log) DocumentsIn(dir string) int {
	for _, d := range l.Dirs {
		if d.Path == dir {
			return len(d.Documents)
		}
	}
	return 0
}

// Problems counts what is wrong with every log.
func (r Result) Problems() int {
	n := 0
	for _, l := range r.Logs {
		n += l.Problems()
	}
	return n
}

// ByHand counts the problems with every log that need fixing by hand.
func (r Result) ByHand() int {
	n := 0
	for _, l := range r.Logs {
		n += l.ByHand()
	}
	return n
}

// Fixable counts the problems with every log that a fix moves or renames,
// those that do not need fixing by hand.
func (r Result) Fixable() int { return r.Problems() - r.ByHand() }

// Duplicates counts the numbers that documents share, in every log, which a
// fix renumbers if asked to.
func (r Result) Duplicates() int {
	n := 0
	for _, l := range r.Logs {
		for _, d := range l.Duplicates {
			if d.Clash == NoClash {
				n++
			}
		}
	}
	return n
}

// Renumberable counts the documents a fix would renumber, in every log.
func (r Result) Renumberable() int {
	n := 0
	for _, l := range r.Logs {
		n += l.Renumberable()
	}
	return n
}

// Documents counts the documents of every log.
func (r Result) Documents() int {
	n := 0
	for _, l := range r.Logs {
		n += l.Documents()
	}
	return n
}

// Directories counts the directories that hold a log.
func (r Result) Directories() int {
	n := 0
	for _, l := range r.Logs {
		n += len(l.Dirs)
	}
	return n
}

// Healthy reports whether nothing is wrong.
func (r Result) Healthy() bool { return r.Problems() == 0 }

func (h Handler) Handle(ctx context.Context, _ Query) (Result, error) {
	root, err := h.Project.Root(ctx)
	if err != nil {
		return Result{}, err
	}
	result := Result{Root: root}
	history, err := h.history(ctx)
	if err != nil {
		return Result{}, err
	}
	for _, kind := range h.Kinds {
		log, err := h.check(ctx, kind, history)
		if err != nil {
			return Result{}, err
		}
		result.Logs = append(result.Logs, log)
	}
	return result, nil
}

// candidate is a directory a log may have been left in.
type candidate struct {
	dir    string
	reason Reason
}

// candidates are the directories a log of kind may have been left in: under
// doc/ rather than docs/, and in each directory the tool looked in before.
func candidates(kind record.Kind) []candidate {
	found := []candidate{{path.Join("doc", kind.Dir), UnderDoc}}
	for _, former := range kind.Former {
		found = append(found,
			candidate{path.Join("docs", former), Former},
			candidate{path.Join("doc", former), Former})
	}
	return found
}

// history returns the history that orders the documents sharing a number, or
// nil when the project is not in a git repository.
func (h Handler) history(ctx context.Context) (History, error) {
	if h.History == nil {
		return nil, nil
	}
	in, err := h.History.Repository(ctx)
	if err != nil || !in {
		return nil, err
	}
	return h.History, nil
}

func (h Handler) check(ctx context.Context, kind record.Kind, history History) (Log, error) {
	home := path.Join("docs", kind.Dir)
	log := Log{Kind: kind, Dir: home}

	homeFiles, homeExists, err := h.Project.Files(ctx, home)
	if err != nil {
		return Log{}, err
	}
	if homeExists {
		log.Dirs = append(log.Dirs, Dir{Path: home, Documents: documents(kind, h.Kinds, homeFiles)})
	}

	movingTo := map[string]string{} // target → the directory moving there
	for _, c := range candidates(kind) {
		files, exists, err := h.Project.Files(ctx, c.dir)
		if err != nil {
			return Log{}, err
		}
		if !exists {
			continue
		}
		docs := documents(kind, h.Kinds, files)
		log.Dirs = append(log.Dirs, Dir{Path: c.dir, Documents: docs})

		move := Move{From: c.dir, To: home, Reason: c.reason, Documents: len(docs)}
		switch {
		case homeExists:
			move.Clash = TargetExists
		case movingTo[home] != "":
			move.Clash, move.ClashWith = SharedTarget, movingTo[home]
		default:
			movingTo[home] = c.dir
		}
		log.Misplaced = append(log.Misplaced, move)
	}

	for _, dir := range log.Dirs {
		log.Misnamed = append(log.Misnamed, misnamed(kind, h.Kinds, dir)...)
	}
	for _, dir := range log.Dirs {
		titles, err := mistitled(ctx, h.Project, kind, h.Kinds, dir)
		if err != nil {
			return Log{}, err
		}
		log.Mistitled = append(log.Mistitled, titles...)
	}
	for _, dir := range log.Dirs {
		duplicates, renumbers, err := shared(ctx, kind, dir, log.Misnamed, history)
		if err != nil {
			return Log{}, err
		}
		log.Duplicates = append(log.Duplicates, duplicates...)
		log.Renumbers = append(log.Renumbers, renumbers...)
	}
	return log, nil
}

// documents returns the names among files that are documents of kind, of
// the kinds a project keeps.
func documents(kind record.Kind, kinds []record.Kind, files []string) []string {
	var docs []string
	for _, name := range files {
		if _, ok := canonical(kind, kinds, name); ok {
			docs = append(docs, name)
		}
	}
	return docs
}

// misnamed returns a rename for each document in dir not named as kind names
// documents today, of the kinds a project keeps.  A rename clashes when a
// document of the new name is there already, or when another document would
// be renamed to it too, as 0001-first.md and adr-0001-first.md would both be
// ADR-0001-first.md.  Names are compared exactly, so that a filesystem that
// ignores case does not see a document as clashing with itself.
func misnamed(kind record.Kind, kinds []record.Kind, dir Dir) []Rename {
	present := map[string]bool{}
	claimants := map[string][]string{} // a name → the documents that would have it
	for _, name := range dir.Documents {
		present[name] = true
		to, _ := canonical(kind, kinds, name)
		claimants[to] = append(claimants[to], name)
	}
	var renames []Rename
	for _, name := range dir.Documents {
		to, _ := canonical(kind, kinds, name)
		if to == name {
			continue
		}
		rename := Rename{Dir: dir.Path, From: name, To: to}
		if present[to] {
			rename.Clash, rename.ClashWith = true, to
		} else {
			for _, other := range claimants[to] {
				if other != name {
					rename.Clash, rename.ClashWith = true, other
					break
				}
			}
		}
		renames = append(renames, rename)
	}
	return renames
}
