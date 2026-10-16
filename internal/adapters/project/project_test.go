package project

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mustardcustard/detoo/internal/adr"
	"github.com/mustardcustard/detoo/internal/domain/record"
	"github.com/mustardcustard/detoo/internal/sty"
)

var kinds = []record.Kind{adr.Kind, sty.Kind}

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
// docs directory, a doc directory or a marker file.
func TestRootIsTheNearestDirectoryHoldingALog(t *testing.T) {
	for _, test := range []struct {
		name, file string
	}{
		{"docs", "docs/adr/ADR-0001-first.md"},
		{"doc", "doc/adr/ADR-0001-first.md"},
		{"marker", ".sty-dir"},
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
	write(t, filepath.Join(root, ".sty-dir"), "docs/sty/\n")
	p := New(root, kinds)

	dir, ok, err := p.Marker(context.Background(), sty.Kind)
	if err != nil || !ok || dir != "docs/sty" {
		t.Errorf("Marker(sty) = %q, %v, %v, want docs/sty, true, nil", dir, ok, err)
	}
	if _, ok, err := p.Marker(context.Background(), adr.Kind); err != nil || ok {
		t.Errorf("Marker(adr) = %v, %v, want no marker", ok, err)
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
