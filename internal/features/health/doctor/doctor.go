// Package doctor checks that a project's logs are where the tools look for
// them, and named as the tools name them today.
//
// It reports and changes nothing.  A log the tools have lost after a rename,
// such as user stories left in docs/sty, is found here, and so is a document
// named before its log's prefix was upper case.  When moving or renaming would
// clash with something already there, the problem is reported as needing
// fixing by hand.
package doctor

import (
	"context"
	"path"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// Query asks for the health of every log of the given kinds.
type Query struct{}

// Project is the port this slice needs: a view of the project's directory.
// Every path it takes and returns is slash separated and relative to the
// project's root, unless the project keeps a log at an absolute path.
type Project interface {
	// Root returns the project's directory.
	Root(ctx context.Context) (string, error)

	// Marker returns the directory that kind's marker file, such as
	// .adr-dir, names, and whether the project has one.
	Marker(ctx context.Context, kind record.Kind) (string, bool, error)

	// Files returns the names of the files in dir, and whether dir exists.
	Files(ctx context.Context, dir string) ([]string, bool, error)
}

// Handler checks the logs of Kinds.
type Handler struct {
	Kinds   []record.Kind
	Project Project
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

	// Dir is where the tool looks for the log: the directory its marker file
	// names, or docs/<Kind.Dir>.
	Dir string

	// Marker is the name of the marker file that names Dir, such as
	// .sty-dir, or empty when there is none.
	Marker string

	// Dirs are the directories that hold the log's documents, Dir first if it
	// exists, then each in Misplaced.
	Dirs []Dir

	// Misplaced are the directories the log is kept in that the tool no
	// longer looks in, or should not.
	Misplaced []Move

	// Misnamed are the documents not named with the log's prefix.
	Misnamed []Rename
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

	// MarkerNamesIt is a directory the log's marker file names, so that the
	// tool still looks there.
	MarkerNamesIt

	// SharedTarget is a directory whose target another misplaced directory
	// moves to as well.
	SharedTarget
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
// should have.  Clash reports that a file of that name is there already.
type Rename struct {
	Dir, From, To string
	Clash         bool
}

// Found reports whether the project keeps the log anywhere.
func (l Log) Found() bool { return len(l.Dirs) > 0 }

// Problems counts what is wrong with the log: each misplaced directory, and
// each misnamed document.
func (l Log) Problems() int { return len(l.Misplaced) + len(l.Misnamed) }

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
	for _, kind := range h.Kinds {
		log, err := h.check(ctx, kind)
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

func (h Handler) check(ctx context.Context, kind record.Kind) (Log, error) {
	standard := path.Join("docs", kind.Dir)
	log := Log{Kind: kind, Dir: standard}
	named, ok, err := h.Project.Marker(ctx, kind)
	if err != nil {
		return Log{}, err
	}
	if ok {
		log.Dir = path.Clean(named)
		log.Marker = "." + kind.Name + "-dir"
	}

	// Where the documents belong.  A marker file that names one of the wrong
	// places is itself the problem, so the log belongs in the standard place.
	home := log.Dir
	for _, c := range candidates(kind) {
		if c.dir == log.Dir {
			home = standard
		}
	}

	homeFiles, homeExists, err := h.Project.Files(ctx, home)
	if err != nil {
		return Log{}, err
	}
	if homeExists {
		log.Dirs = append(log.Dirs, Dir{Path: home, Documents: documents(kind, homeFiles)})
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
		docs := documents(kind, files)
		log.Dirs = append(log.Dirs, Dir{Path: c.dir, Documents: docs})

		move := Move{From: c.dir, To: home, Reason: c.reason, Documents: len(docs)}
		switch {
		case c.dir == log.Dir && log.Marker != "":
			move.Clash = MarkerNamesIt
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
		log.Misnamed = append(log.Misnamed, misnamed(kind, dir)...)
	}
	return log, nil
}

// documents returns the names among files that are documents of kind.
func documents(kind record.Kind, files []string) []string {
	var docs []string
	for _, name := range files {
		if kind.IsRecordName(name) {
			docs = append(docs, name)
		}
	}
	return docs
}

// misnamed returns a rename for each document in dir not named as kind names
// documents today.  A rename clashes when a file of the new name is there
// already: compared by exact name, so that a filesystem that ignores case does
// not see a document as clashing with itself.
func misnamed(kind record.Kind, dir Dir) []Rename {
	present := map[string]bool{}
	for _, name := range dir.Documents {
		present[name] = true
	}
	var renames []Rename
	for _, name := range dir.Documents {
		canonical, _ := kind.CanonicalName(name)
		if canonical != name {
			renames = append(renames, Rename{Dir: dir.Path, From: name, To: canonical, Clash: present[canonical]})
		}
	}
	return renames
}
