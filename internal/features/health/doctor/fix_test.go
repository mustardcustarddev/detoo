package doctor

import (
	"context"
	"fmt"
	"maps"
	"path"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/mustardcustard/detoo/internal/adr"
	"github.com/mustardcustard/detoo/internal/domain/record"
	"github.com/mustardcustard/detoo/internal/sty"
)

// memWorkspace is a project held in memory as its files, by path, and the
// directories each marker file names.
type memWorkspace struct {
	files   map[string]string
	markers map[string]string
}

func (w *memWorkspace) Root(context.Context) (string, error) { return "/home/ann/handbook", nil }

func (w *memWorkspace) Marker(_ context.Context, kind record.Kind) (string, bool, error) {
	dir, ok := w.markers[kind.Name]
	return dir, ok, nil
}

func (w *memWorkspace) Files(_ context.Context, dir string) ([]string, bool, error) {
	var names []string
	exists := false
	for _, name := range slices.Sorted(maps.Keys(w.files)) {
		if rest, ok := strings.CutPrefix(name, dir+"/"); ok {
			exists = true
			if !strings.Contains(rest, "/") {
				names = append(names, rest)
			}
		}
	}
	return names, exists, nil
}

func (w *memWorkspace) Move(_ context.Context, from, to string) error {
	if _, exists, _ := w.Files(context.Background(), to); exists {
		return fmt.Errorf("%s exists", to)
	}
	for _, name := range slices.Collect(maps.Keys(w.files)) {
		if rest, ok := strings.CutPrefix(name, from+"/"); ok {
			w.files[to+"/"+rest] = w.files[name]
			delete(w.files, name)
		}
	}
	return nil
}

func (w *memWorkspace) Rename(_ context.Context, dir, from, to string) error {
	if _, ok := w.files[path.Join(dir, to)]; ok {
		return fmt.Errorf("%s exists", to)
	}
	w.files[path.Join(dir, to)] = w.files[path.Join(dir, from)]
	delete(w.files, path.Join(dir, from))
	return nil
}

// RemoveIfEmpty removes nothing, as a directory in memory exists only while
// it holds a file.
func (w *memWorkspace) RemoveIfEmpty(_ context.Context, dir string) (bool, error) {
	_, exists, _ := w.Files(context.Background(), dir)
	return !exists, nil
}

func (w *memWorkspace) Markdown(_ context.Context, dirs []string) ([]string, error) {
	var found []string
	for _, name := range slices.Sorted(maps.Keys(w.files)) {
		if !strings.HasSuffix(name, ".md") {
			continue
		}
		in := !strings.Contains(name, "/")
		for _, dir := range dirs {
			in = in || strings.HasPrefix(name, dir+"/")
		}
		if in {
			found = append(found, name)
		}
	}
	return found, nil
}

func (w *memWorkspace) Read(_ context.Context, name string) (string, error) {
	return w.files[name], nil
}

func (w *memWorkspace) Write(_ context.Context, name, content string) error {
	w.files[name] = content
	return nil
}

func fix(t *testing.T, w *memWorkspace, kinds ...record.Kind) Fixed {
	t.Helper()
	fixed, err := Fixer{Kinds: kinds, Project: w}.Handle(context.Background(), Command{})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	return fixed
}

// A log of adr-tools: under doc/, named without a prefix, and linked by
// number alone.  The fix moves and renames it, and the links follow.
func TestFixMovesAndRenamesALogAndItsLinks(t *testing.T) {
	w := &memWorkspace{files: map[string]string{
		"doc/adr/0001-first.md":  "# 1. First\n\nSuperseded by [2. Second](0002-second.md)\n",
		"doc/adr/0002-second.md": "# 2. Second\n\nSupersedes [1. First](0001-first.md)\n",
		"doc/adr/README.md":      "* [1. First](0001-first.md)\n* [2. Second](0002-second.md)\n",
		"README.md":              "Read [the decisions](doc/adr/) and [the first](doc/adr/0001-first.md).\n",
	}}
	fixed := fix(t, w, adr.Kind)

	want := map[string]string{
		"docs/adr/ADR-0001-first.md":  "# 1. First\n\nSuperseded by [2. Second](ADR-0002-second.md)\n",
		"docs/adr/ADR-0002-second.md": "# 2. Second\n\nSupersedes [1. First](ADR-0001-first.md)\n",
		"docs/adr/README.md":          "* [1. First](ADR-0001-first.md)\n* [2. Second](ADR-0002-second.md)\n",
		"README.md":                   "Read [the decisions](docs/adr/) and [the first](docs/adr/ADR-0001-first.md).\n",
	}
	if !reflect.DeepEqual(w.files, want) {
		t.Errorf("files are\n%v\nwant\n%v", w.files, want)
	}
	if fixed.Count() != 3 || !fixed.After.Healthy() || fixed.Before.Healthy() {
		t.Errorf("Count() = %d, After.Healthy() = %v, Before.Healthy() = %v, want 3, true, false",
			fixed.Count(), fixed.After.Healthy(), fixed.Before.Healthy())
	}
	wantRelinked := []Relinked{
		{"README.md", 2, false},
		{"docs/adr/ADR-0001-first.md", 1, false},
		{"docs/adr/ADR-0002-second.md", 1, false},
		{"docs/adr/README.md", 2, false},
	}
	if !reflect.DeepEqual(fixed.Relinked, wantRelinked) {
		t.Errorf("Relinked = %v, want %v", fixed.Relinked, wantRelinked)
	}
	if fixed.Removed == nil || fixed.Removed[0] != "doc" {
		t.Errorf("Removed = %v, want doc", fixed.Removed)
	}
	if fixed.Changes() != 3+1+4 {
		t.Errorf("Changes() = %d, want 8", fixed.Changes())
	}
}

// Clashes are left as they are, and reported by the check after the fix,
// while the documents that can be renamed are.
func TestFixLeavesClashesAlone(t *testing.T) {
	w := &memWorkspace{files: map[string]string{
		"docs/sty/0001-first.md":              "",
		"docs/userstory/STORY-0002-second.md": "",
		"docs/userstory/0002-second.md":       "",
	}}
	fixed := fix(t, w, sty.Kind)

	want := []string{"docs/sty/STORY-0001-first.md", "docs/userstory/0002-second.md", "docs/userstory/STORY-0002-second.md"}
	if got := slices.Sorted(maps.Keys(w.files)); !reflect.DeepEqual(got, want) {
		t.Errorf("files are %v, want %v", got, want)
	}
	if len(fixed.Logs[0].Moved) != 0 || len(fixed.Logs[0].Renamed) != 1 {
		t.Errorf("Logs = %+v, want one rename and no move", fixed.Logs)
	}
	if fixed.After.Healthy() || fixed.After.Problems() != fixed.After.ByHand() || fixed.After.ByHand() != 2 {
		t.Errorf("After has %d problems, %d by hand, want 2 by hand", fixed.After.Problems(), fixed.After.ByHand())
	}
}

// A healthy project is left as it is.
func TestFixChangesNothingInAHealthyProject(t *testing.T) {
	w := &memWorkspace{files: map[string]string{
		"docs/adr/ADR-0001-first.md": "[x](../../README.md)\n",
		"README.md":                  "[y](docs/adr/ADR-0001-first.md)\n",
	}}
	before := maps.Clone(w.files)
	fixed := fix(t, w, adr.Kind)
	if fixed.Changes() != 0 || !fixed.After.Healthy() || !reflect.DeepEqual(w.files, before) {
		t.Errorf("Changes() = %d, files = %v, want nothing changed", fixed.Changes(), w.files)
	}
}
