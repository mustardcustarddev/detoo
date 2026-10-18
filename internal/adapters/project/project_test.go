package project

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"

	"github.com/mustardcustard/detoo/internal/adr"
	"github.com/mustardcustard/detoo/internal/domain/record"
	"github.com/mustardcustard/detoo/internal/ust"
)

var kinds = []record.Kind{adr.Kind, ust.Kind}

func write(t *testing.T, name, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(name, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// resolved returns dir with its symbolic links resolved, as Root reports it
// on a system whose temporary directory is behind one.
func resolved(t *testing.T, dir string) string {
	t.Helper()
	dir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

// The project is the nearest directory above the working one that holds a
// docs directory, a doc directory or a marker file, under its name or a
// former one.
func TestRootIsTheNearestDirectoryHoldingALog(t *testing.T) {
	for _, test := range []struct {
		name, file string
	}{
		{"docs", "docs/adr/ADR-0001-first.md"},
		{"doc", "doc/adr/ADR-0001-first.md"},
		{"marker", ".ust-dir"},
		{"former marker", ".sty-dir"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := resolved(t, t.TempDir())
			write(t, filepath.Join(root, test.file), "docs/sty\n")
			work := filepath.Join(root, "src", "app")
			if err := os.MkdirAll(work, 0o755); err != nil {
				t.Fatal(err)
			}
			got, err := New(work, kinds).Root(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			if got != root {
				t.Errorf("Root() = %q, want %q", got, root)
			}
		})
	}
}

func TestRootIsTheWorkingDirectoryWhenNothingIsFound(t *testing.T) {
	work := resolved(t, t.TempDir())
	got, err := New(work, kinds).Root(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	// An ancestor of the temporary directory may hold a docs directory, but
	// none should; if one does, the test says so.
	if got != work {
		t.Errorf("Root() = %q, want %q", got, work)
	}
}

func TestMarkerReadsTheNamedDirectory(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, ".ust-dir"), "docs/sty/\n")
	p := New(root, kinds)

	dir, ok, err := p.Marker(context.Background(), ust.Kind.Marker())
	if err != nil || !ok || dir != "docs/sty" {
		t.Errorf("Marker(.ust-dir) = %q, %v, %v, want docs/sty, true, nil", dir, ok, err)
	}
	if _, ok, err := p.Marker(context.Background(), adr.Kind.Marker()); err != nil || ok {
		t.Errorf("Marker(.adr-dir) = %v, %v, want no marker", ok, err)
	}
}

// Files lists the files in a directory, leaving out subdirectories and hidden
// files, and reports a directory that is not there, or is a file, as missing.
func TestFilesListsTheFilesInADirectory(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "docs/adr/ADR-0001-first.md"), "")
	write(t, filepath.Join(root, "docs/adr/.DS_Store"), "")
	write(t, filepath.Join(root, "docs/adr/templates/template.md"), "")
	write(t, filepath.Join(root, "docs/rfc"), "")
	p := New(root, kinds)

	files, ok, err := p.Files(context.Background(), "docs/adr")
	if err != nil || !ok || !reflect.DeepEqual(files, []string{"ADR-0001-first.md"}) {
		t.Errorf("Files(docs/adr) = %q, %v, %v, want [ADR-0001-first.md], true, nil", files, ok, err)
	}
	for _, dir := range []string{"docs/sty", "docs/rfc"} {
		if _, ok, err := p.Files(context.Background(), dir); err != nil || ok {
			t.Errorf("Files(%s) = %v, %v, want missing", dir, ok, err)
		}
	}
}

func exists(name string) bool {
	_, err := os.Lstat(name)
	return err == nil
}

// Move makes the directories it needs, and refuses to move onto anything
// there already.
func TestMoveMovesADirectory(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "doc/adr/0001-first.md"), "first")
	write(t, filepath.Join(root, "doc/rfc/0001-first.md"), "")
	write(t, filepath.Join(root, "docs/rfc/README.md"), "")
	p := New(root, kinds)

	if err := p.Move(context.Background(), "doc/adr", "docs/adr"); err != nil {
		t.Fatal(err)
	}
	if exists(filepath.Join(root, "doc/adr")) || !exists(filepath.Join(root, "docs/adr/0001-first.md")) {
		t.Errorf("doc/adr was not moved to docs/adr")
	}
	if err := p.Move(context.Background(), "doc/rfc", "docs/rfc"); err == nil {
		t.Errorf("Move onto an existing directory succeeded")
	}
}

// Rename refuses to rename over another file, but changes the case of a name
// on a filesystem that ignores case.
func TestRenameRenamesAFile(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "docs/userstory/story-0001-first.md"), "")
	write(t, filepath.Join(root, "docs/userstory/0002-second.md"), "")
	write(t, filepath.Join(root, "docs/userstory/STORY-0003-third.md"), "")
	p := New(root, kinds)

	if err := p.Rename(context.Background(), "docs/userstory", "story-0001-first.md", "STORY-0001-first.md"); err != nil {
		t.Fatal(err)
	}
	names, _, _ := p.Files(context.Background(), "docs/userstory")
	if !slices.Contains(names, "STORY-0001-first.md") || slices.Contains(names, "story-0001-first.md") {
		t.Errorf("Files() = %v, want STORY-0001-first.md in place of story-0001-first.md", names)
	}
	if err := p.Rename(context.Background(), "docs/userstory", "0002-second.md", "STORY-0003-third.md"); err == nil {
		t.Errorf("Rename over another file succeeded")
	}
}

func TestRemoveIfEmptyRemovesOnlyAnEmptyDirectory(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "docs/adr/ADR-0001-first.md"), "")
	if err := os.Mkdir(filepath.Join(root, "doc"), 0o755); err != nil {
		t.Fatal(err)
	}
	p := New(root, kinds)

	for _, test := range []struct {
		dir     string
		removed bool
	}{{"docs", false}, {"doc", true}, {"missing", false}} {
		removed, err := p.RemoveIfEmpty(context.Background(), test.dir)
		if err != nil || removed != test.removed {
			t.Errorf("RemoveIfEmpty(%s) = %v, %v, want %v", test.dir, removed, err, test.removed)
		}
	}
	if exists(filepath.Join(root, "doc")) || !exists(filepath.Join(root, "docs")) {
		t.Errorf("RemoveIfEmpty removed the wrong directory")
	}
}

// Markdown finds the documents in the root and under the directories given,
// once each, leaving out hidden ones and directories that are not there.
func TestMarkdownFindsTheDocuments(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{
		"README.md", "main.go", ".hidden.md", "src/notes.md",
		"docs/adr/ADR-0001-first.md", "docs/adr/templates/template.md",
		"docs/.cache/x.md", "docs/rfc/diagram.png",
	} {
		write(t, filepath.Join(root, name), "")
	}
	p := New(root, kinds)

	got, err := p.Markdown(context.Background(), []string{"docs", "docs/adr", "doc"})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"README.md", "docs/adr/ADR-0001-first.md", "docs/adr/templates/template.md"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Markdown() = %v, want %v", got, want)
	}
}

func TestWriteReplacesAFile(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "README.md"), "before")
	p := New(root, kinds)
	if err := p.Write(context.Background(), "README.md", "after"); err != nil {
		t.Fatal(err)
	}
	if got, err := p.Read(context.Background(), "README.md"); err != nil || got != "after" {
		t.Errorf("Read() = %q, %v, want after", got, err)
	}
}
