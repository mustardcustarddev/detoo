package doctor

import (
	"cmp"
	"context"
	"fmt"
	"maps"
	"path"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// prefixed matches the file name of a record after a prefix of letters, such
// as STOY-0003-title.md, capturing the prefix and the rest.
var prefixed = regexp.MustCompile(`^([A-Za-z]+)-([0-9]+-[^/]*\.md)$`)

// leadingDigits matches the number a record's file name starts with, after
// its prefix.
var leadingDigits = regexp.MustCompile(`^[0-9]+`)

// canonical returns the name a document of kind should have, and whether name
// is a document of kind at all: named as kind names its records, now or
// before, or with a slip of its prefix, as STOY-0003-title.md is a user
// story's.  A prefix is a slip of kind's when kind's is the nearest to it of
// the prefixes of kinds, and near enough; the prefix of another kind, such as
// RFC- in a directory that adr and rfc share, is not a slip.
func canonical(kind record.Kind, kinds []record.Kind, name string) (string, bool) {
	if to, ok := kind.CanonicalName(name); ok {
		return to, true
	}
	m := prefixed.FindStringSubmatch(name)
	if m == nil || kind.Prefix == "" || nearest(strings.ToUpper(m[1]), kinds) != kind.Name {
		return "", false
	}
	return kind.Prefix + "-" + m[2], true
}

// nearest returns the name of the kind whose prefix prefix is a slip of, or ""
// when it is none's: when it is a kind's prefix itself, when two kinds' are
// as near, or when none is near enough.  A slip is a few letters added,
// dropped or changed: at most one for a prefix of up to three letters, and
// half the letters of a longer one.
func nearest(prefix string, kinds []record.Kind) string {
	best, bestDistance, tie := "", -1, false
	for _, k := range kinds {
		if k.Prefix == "" {
			continue
		}
		d := distance(prefix, k.Prefix)
		if d == 0 {
			return ""
		}
		if d > max(1, len(k.Prefix)/2) {
			continue
		}
		switch {
		case bestDistance < 0 || d < bestDistance:
			best, bestDistance, tie = k.Name, d, false
		case d == bestDistance:
			tie = true
		}
	}
	if tie {
		return ""
	}
	return best
}

// distance counts the letters to add, drop or change to turn a into b.
func distance(a, b string) int {
	prev := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		cur := make([]int, len(b)+1)
		cur[0] = i
		for j := 1; j <= len(b); j++ {
			change := prev[j-1]
			if a[i-1] != b[j-1] {
				change++
			}
			cur[j] = min(prev[j]+1, cur[j-1]+1, change)
		}
		prev = cur
	}
	return prev[len(b)]
}

// numbered is a document of a directory, by the name it has on disk and the
// name a fix gives it, and its number.
type numbered struct {
	file, name string
	number     int

	// added is when git first recorded the document, if dated.
	added time.Time
	dated bool
}

// earlier reports whether git recorded a before b.  A document git has not
// recorded is the latest of all.
func (a numbered) earlier(b numbered) bool {
	return a.dated && (!b.dated || a.added.Before(b.added))
}

// order sorts documents that share a number: the earliest recorded first,
// then by name.
func order(a, b numbered) int {
	switch {
	case a.earlier(b):
		return -1
	case b.earlier(a):
		return 1
	}
	return strings.Compare(a.name, b.name)
}

// shared returns the numbers that more than one document in dir has, by the
// names renames give them, and the renumbering that gives each a number of
// its own.
//
// With history, the document git recorded first keeps the number.  Each other
// goes after the last document git recorded before it, from the one that
// keeps the number on, and before the one after that: a document git has not
// recorded yet goes last.  The documents after it move up as far as they must
// to make room, so that a gap in the numbers takes up the shift.  Without
// history, or when a rename in dir clashes, the duplicates need renumbering by
// hand, and each but the first by name is shown with a number after the last.
// A document whose rename clashes is a problem already, and is left out.
func shared(ctx context.Context, kind record.Kind, dir Dir, renames []Rename, history History) ([]Duplicate, []Renumber, error) {
	to := map[string]string{}
	clashing := map[string]bool{}
	clash := NoClash
	for _, rn := range renames {
		if rn.Dir != dir.Path {
			continue
		}
		if rn.Clash {
			clash = NamesClash
			clashing[rn.From] = true
			continue
		}
		to[rn.From] = rn.To
	}
	docs := make([]numbered, 0, len(dir.Documents))
	byNumber := map[int][]numbered{}
	for _, file := range dir.Documents {
		if clashing[file] {
			continue
		}
		name := file
		if renamed, ok := to[file]; ok {
			name = renamed
		}
		d := numbered{file: file, name: name, number: kind.Number(name)}
		docs = append(docs, d)
		byNumber[d.number] = append(byNumber[d.number], d)
	}
	if len(docs) == len(byNumber) {
		return nil, nil, nil
	}
	if clash == NoClash && history == nil {
		clash = NoHistory
	}

	if clash == NoClash {
		for i := range docs {
			var err error
			docs[i].added, docs[i].dated, err = history.Added(ctx, path.Join(dir.Path, docs[i].file))
			if err != nil {
				return nil, nil, err
			}
		}
		clear(byNumber)
		for _, d := range docs {
			byNumber[d.number] = append(byNumber[d.number], d)
		}
	}

	// The log in order of number, each number's documents in order.
	var log, moving []numbered
	var duplicates []Duplicate
	for _, number := range slices.Sorted(maps.Keys(byNumber)) {
		group := byNumber[number]
		slices.SortFunc(group, order)
		log = append(log, group[0])
		if len(group) == 1 {
			continue
		}
		moving = append(moving, group[1:]...)
		d := Duplicate{Dir: dir.Path, Number: number, Clash: clash}
		for _, g := range group {
			d.Documents = append(d.Documents, g.name)
		}
		duplicates = append(duplicates, d)
	}
	slices.SortFunc(moving, order)

	if clash != NoClash {
		last := log[len(log)-1].number
		renumbers := make([]Renumber, 0, len(moving))
		for i, d := range moving {
			renumbers = append(renumbers, Renumber{
				Dir: dir.Path, File: d.file, From: d.name, To: renumbered(kind, d.name, last+1+i),
				Duplicate: true, Clash: clash,
			})
		}
		return duplicates, renumbers, nil
	}

	isMoving := map[string]bool{}
	for _, d := range moving {
		isMoving[d.name] = true
		keeps := slices.IndexFunc(log, func(k numbered) bool { return k.number == d.number })
		at := keeps
		for i := keeps + 1; i < len(log); i++ {
			if !d.earlier(log[i]) {
				at = i
			}
		}
		log = slices.Insert(log, at+1, d)
	}

	names := make([]string, len(log))
	numbers := make([]int, len(log))
	prev := 0
	for i, d := range log {
		numbers[i] = max(d.number, prev+1)
		names[i] = renumbered(kind, d.name, numbers[i])
		prev = numbers[i]
	}
	var renumbers []Renumber
	for i, d := range log {
		if numbers[i] == d.number {
			continue
		}
		r := Renumber{Dir: dir.Path, File: d.file, From: d.name, To: names[i], Duplicate: isMoving[d.name]}
		if r.Duplicate {
			r.After = names[i-1]
			if i+1 < len(log) {
				r.Before = names[i+1]
			}
		}
		renumbers = append(renumbers, r)
	}
	return duplicates, renumbers, nil
}

// renumbered returns a document's name, as kind names documents, with another
// number.
func renumbered(kind record.Kind, name string, number int) string {
	prefix := ""
	if kind.Prefix != "" {
		prefix = kind.Prefix + "-"
	}
	rest := strings.TrimPrefix(name, prefix)
	digits := leadingDigits.FindString(rest)
	return fmt.Sprintf("%s%0*d%s", prefix, max(4, len(digits)), number, rest[len(digits):])
}

// compareNumbers orders renumbers so that each is made before the document it
// moves onto is moved away: the highest numbers first.
func compareNumbers(kind record.Kind) func(a, b Renumber) int {
	return func(a, b Renumber) int {
		return cmp.Or(
			cmp.Compare(kind.Number(b.From), kind.Number(a.From)),
			cmp.Compare(kind.Number(b.To), kind.Number(a.To)),
		)
	}
}
