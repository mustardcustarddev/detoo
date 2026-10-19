package doctor

import (
	"context"
	"errors"
	"fmt"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// Command asks for every problem with the logs of the given kinds that can be
// fixed to be fixed.
type Command struct {
	// Git moves and renames with git, so that the repository records each as
	// a rename rather than a deletion and a new file.
	Git bool

	// Renumber gives each document that shares a number with one git
	// recorded before it a number of its own, and moves up the documents
	// after it to make room.
	Renumber bool
}

// ErrNotARepository is returned when a fix is asked to move with git in a
// project that is not in a git repository.
var ErrNotARepository = errors.New("the project is not in a git repository, so cannot move with git")

// Workspace is the port fixing needs: a view of the project's directory that
// can also be changed.  Paths are as Project takes them.
type Workspace interface {
	Project

	// Move moves the directory from to to, which must not exist, making the
	// directories above to that it needs.
	Move(ctx context.Context, from, to string) error

	// Rename renames the file from, in dir, to to, which must not be another
	// file.
	Rename(ctx context.Context, dir, from, to string) error

	// RemoveIfEmpty removes dir if nothing is left in it, and reports
	// whether it did.
	RemoveIfEmpty(ctx context.Context, dir string) (bool, error)

	// Markdown returns the Markdown files in the project's root, and in dirs
	// and every directory beneath them, leaving out hidden ones.
	Markdown(ctx context.Context, dirs []string) ([]string, error)

	// Write replaces the content of a file.
	Write(ctx context.Context, name, content string) error
}

// Git is a Workspace that moves and renames with git where it can: the files
// git tracks are moved with git mv, and the rest as Workspace moves them.  Its
// History orders the documents that share a number.
type Git interface {
	Workspace
	History
}

// Fixer fixes the logs of Kinds, in Project, or in Git when asked to move
// with git.
type Fixer struct {
	Kinds   []record.Kind
	Project Workspace
	Git     Git
}

// Repository reports whether the project is in a git repository, so that a
// fix can be offered the choice of moving with git.
func (f Fixer) Repository(ctx context.Context) (bool, error) {
	if f.Git == nil {
		return false, nil
	}
	return f.Git.Repository(ctx)
}

// Fixed is what a fix changed, and the health of the logs before and after.
type Fixed struct {
	Before, After Result

	// Git reports that the fix moved and renamed with git.
	Git bool

	// Logs are the fixes made to each log, in the order of the kinds.
	Logs []Fixes

	// Removed are the directories that moving logs out of left empty, and
	// which were removed: doc, when the logs under it have moved to docs.
	Removed []string

	// Relinked are the documents whose links were changed to follow the
	// documents that moved, in the order of their paths.
	Relinked []Relinked
}

// Fixes are the problems fixed in one log.  Renames name the directory the
// document was in before the fix, and the name it had there.
type Fixes struct {
	Moved   []Move
	Renamed []Rename

	// Duplicates are the numbers documents shared, and Renumbered the
	// documents renumbered so that each has its own.
	Duplicates []Duplicate
	Renumbered []Renumber
}

// Relinked is a document whose links were changed, by the path it has after
// the fix, and how many of its links were changed.  Retitled reports that the
// number in its title was changed, as the document was renumbered.
type Relinked struct {
	Path     string
	Links    int
	Retitled bool
}

// Count counts the problems fixed in the log.
func (f Fixes) Count() int { return len(f.Moved) + len(f.Renamed) + len(f.Duplicates) }

// Files counts the directories and documents moved or renamed in the log.
func (f Fixes) Files() int { return len(f.Moved) + len(f.Renamed) + len(f.Renumbered) }

// Count counts the problems fixed in every log.
func (f Fixed) Count() int {
	n := 0
	for _, l := range f.Logs {
		n += l.Count()
	}
	return n
}

// Renumbered counts the documents renumbered in every log.
func (f Fixed) Renumbered() int {
	n := 0
	for _, l := range f.Logs {
		n += len(l.Renumbered)
	}
	return n
}

// Changes counts what the fix changed: each directory moved or removed, each
// document renamed or renumbered, and each document whose links or title were
// changed.
func (f Fixed) Changes() int {
	n := len(f.Removed) + len(f.Relinked)
	for _, l := range f.Logs {
		n += l.Files()
	}
	return n
}

// Handle checks the logs, moves each misplaced directory and renames each
// misnamed document that would not clash, renumbers the documents that share
// a number if asked to, updates the links to what moved, and checks the logs
// again.  A clash is left as it is, to be fixed by hand.
func (f Fixer) Handle(ctx context.Context, cmd Command) (Fixed, error) {
	if cmd.Git {
		in, err := f.Repository(ctx)
		if err != nil {
			return Fixed{}, err
		}
		if !in {
			return Fixed{}, ErrNotARepository
		}
		f.Project = f.Git
	}
	checker := Handler{Kinds: f.Kinds, Project: f.Project}
	if f.Git != nil {
		checker.History = f.Git
	}
	before, err := checker.Handle(ctx, Query{})
	if err != nil {
		return Fixed{}, err
	}
	fixed := Fixed{Before: before, Git: cmd.Git, Logs: make([]Fixes, len(before.Logs))}
	moves := relocation{dirs: map[string]string{}, files: map[string]string{}, numbers: map[string][2]int{}}

	for i, l := range before.Logs {
		for _, m := range l.Misplaced {
			if m.Clash != NoClash {
				continue
			}
			if err := f.Project.Move(ctx, m.From, m.To); err != nil {
				return Fixed{}, err
			}
			moves.dirs[m.From] = m.To
			fixed.Logs[i].Moved = append(fixed.Logs[i].Moved, m)
		}
		// A document is renamed where it is now, which is where its
		// directory moved to if it moved.
		for _, rn := range l.Misnamed {
			if rn.Clash {
				continue
			}
			dir := moves.dir(rn.Dir)
			if err := f.Project.Rename(ctx, dir, rn.From, rn.To); err != nil {
				return Fixed{}, err
			}
			moves.files[path.Join(rn.Dir, rn.From)] = path.Join(dir, rn.To)
			fixed.Logs[i].Renamed = append(fixed.Logs[i].Renamed, rn)
		}
		if cmd.Renumber {
			if err := f.renumber(ctx, l, &fixed.Logs[i], moves); err != nil {
				return Fixed{}, err
			}
		}
	}

	if fixed.Removed, err = f.removeEmptied(ctx, fixed); err != nil {
		return Fixed{}, err
	}
	if fixed.Count() > 0 {
		if fixed.Relinked, err = f.relink(ctx, before, moves); err != nil {
			return Fixed{}, err
		}
	}

	if fixed.After, err = checker.Handle(ctx, Query{}); err != nil {
		return Fixed{}, err
	}
	return fixed, nil
}

// renumber renames each document of l that a fix renumbers, where it is now,
// the highest numbers first, so that no document is renamed onto another that
// has yet to move up.
func (f Fixer) renumber(ctx context.Context, l Log, fixes *Fixes, moves relocation) error {
	for _, d := range l.Duplicates {
		if d.Clash == NoClash {
			fixes.Duplicates = append(fixes.Duplicates, d)
		}
	}
	renumbers := slices.Clone(l.Renumbers)
	slices.SortStableFunc(renumbers, compareNumbers(l.Kind))
	for _, r := range renumbers {
		if r.Clash != NoClash {
			continue
		}
		dir := moves.dir(r.Dir)
		if err := f.Project.Rename(ctx, dir, r.From, r.To); err != nil {
			return err
		}
		moves.files[path.Join(r.Dir, r.File)] = path.Join(dir, r.To)
		moves.numbers[path.Join(dir, r.To)] = [2]int{l.Kind.Number(r.From), l.Kind.Number(r.To)}
		fixes.Renumbered = append(fixes.Renumbered, r)
	}
	// Reported in the order of the log.
	slices.Reverse(fixes.Renumbered)
	return nil
}

// removeEmptied removes each directory that held a directory moved away, such
// as doc, if the move left it empty.
func (f Fixer) removeEmptied(ctx context.Context, fixed Fixed) ([]string, error) {
	var removed []string
	for _, l := range fixed.Logs {
		for _, m := range l.Moved {
			parent := path.Dir(m.From)
			if parent == "." || parent == "/" || slices.Contains(removed, parent) {
				continue
			}
			ok, err := f.Project.RemoveIfEmpty(ctx, parent)
			if err != nil {
				return nil, err
			}
			if ok {
				removed = append(removed, parent)
			}
		}
	}
	return removed, nil
}

// relink changes the links in the project's Markdown documents that would
// break because of the moves: links to documents that moved, and the links of
// documents that moved themselves.
func (f Fixer) relink(ctx context.Context, before Result, moves relocation) ([]Relinked, error) {
	dirs := []string{"docs", "doc"}
	for _, l := range before.Logs {
		dirs = append(dirs, l.Dir)
	}
	for _, l := range before.Logs {
		for _, m := range l.Misplaced {
			dirs = append(dirs, m.To)
		}
	}
	files, err := f.Project.Markdown(ctx, dirs)
	if err != nil {
		return nil, err
	}
	var relinked []Relinked
	for _, name := range files {
		content, err := f.Project.Read(ctx, name)
		if err != nil {
			return nil, err
		}
		changed, n := relink(content, moves.origin(name), name, moves.to, moves.renumbered)
		changed, retitled := retitle(changed, moves.numbers[name])
		if n == 0 && !retitled {
			continue
		}
		if err := f.Project.Write(ctx, name, changed); err != nil {
			return nil, err
		}
		relinked = append(relinked, Relinked{Path: name, Links: n, Retitled: retitled})
	}
	return relinked, nil
}

// title matches the first line of a record, which begins with its number:
// "# 3. Use PostgreSQL".
var title = regexp.MustCompile(`^# ([0-9]+)\. `)

// retitle returns content with the number its title begins with changed from
// the first of numbers to the second, and whether it changed it.  Content
// whose title begins with another number, or with none, is left as it is.
func retitle(content string, numbers [2]int) (string, bool) {
	m := title.FindStringSubmatchIndex(content)
	if numbers == [2]int{} || m == nil {
		return content, false
	}
	if n, err := strconv.Atoi(content[m[2]:m[3]]); err != nil || n != numbers[0] {
		return content, false
	}
	return content[:m[2]] + strconv.Itoa(numbers[1]) + content[m[3]:], true
}

// relocation is where the fix moved things: each directory moved, and each
// document renamed, from its path before the fix to its path after, and the
// number each document renumbered had and has, by its path after.
type relocation struct {
	dirs, files map[string]string
	numbers     map[string][2]int
}

// renumbered returns the text a link to the document at name, its path after
// the fix, should have: a link whose text begins with the number the
// document had, as "3. Use PostgreSQL", begins with the one it has now.
func (r relocation) renumbered(name, text string) (string, bool) {
	numbers, ok := r.numbers[name]
	if !ok {
		return text, false
	}
	rest, ok := strings.CutPrefix(text, fmt.Sprintf("%d. ", numbers[0]))
	if !ok {
		return text, false
	}
	return fmt.Sprintf("%d. %s", numbers[1], rest), true
}

// dir returns where a directory is after the moves.
func (r relocation) dir(dir string) string {
	if to, ok := r.dirs[dir]; ok {
		return to
	}
	return dir
}

// to returns where a path before the fix is after it.
func (r relocation) to(name string) string {
	if to, ok := r.files[name]; ok {
		return to
	}
	return swapPrefix(name, r.dirs)
}

// origin returns where a path after the fix was before it.
func (r relocation) origin(name string) string {
	for from, to := range r.files {
		if to == name {
			return from
		}
	}
	reverse := make(map[string]string, len(r.dirs))
	for from, to := range r.dirs {
		reverse[to] = from
	}
	return swapPrefix(name, reverse)
}

// swapPrefix replaces the directory name is in, or is, by the one dirs maps
// it to.
func swapPrefix(name string, dirs map[string]string) string {
	for from, to := range dirs {
		if name == from {
			return to
		}
		if len(name) > len(from) && name[:len(from)+1] == from+"/" {
			return to + name[len(from):]
		}
	}
	return name
}
