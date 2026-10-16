package logstore

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// The store is the one place that touches the disk, so these are the only tests
// in the application that need a temporary directory.  The slices above it are
// tested against fakes.

func newStore(t *testing.T) *Store {
	t.Helper()
	return New(t.TempDir(), kind)
}

func writeRecord(t *testing.T, s *Store, name, content string) {
	t.Helper()
	file := s.path(name)
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestDirDefaultsWhenNothingIsFound(t *testing.T) {
	dir, err := newStore(t).Dir(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if dir != adrDir {
		t.Errorf("Dir() = %q, want %q", dir, adrDir)
	}
}

func TestDirFindsAMarkerFile(t *testing.T) {
	s := newStore(t)
	write(t, filepath.Join(s.workdir, adrMarker), "architecture-log\n")

	dir, err := s.Dir(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if want := "architecture-log"; dir != want {
		t.Errorf("Dir() = %q, want %q", dir, want)
	}
}

func TestDirSearchesAncestors(t *testing.T) {
	for _, test := range []struct {
		name  string
		setup func(root string)
		want  string
	}{
		{
			name:  "default directory in an ancestor",
			setup: func(root string) { mkdirAll(t, filepath.Join(root, "docs", "adr")) },
			want:  "../../docs/adr",
		},
		{
			name:  "marker file in an ancestor",
			setup: func(root string) { write(t, filepath.Join(root, adrMarker), "architecture-log\n") },
			want:  "../../architecture-log",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			test.setup(root)
			deep := filepath.Join(root, "src", "main")
			mkdirAll(t, deep)

			dir, err := New(deep, kind).Dir(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			if dir != test.want {
				t.Errorf("Dir() = %q, want %q", dir, test.want)
			}
		})
	}
}

// TestDirPrefersTheNearestMarker checks that a marker file closer to the working
// directory wins over one further up.
func TestDirPrefersTheNearestMarker(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, adrMarker), "outer\n")
	inner := filepath.Join(root, "component")
	mkdirAll(t, inner)
	write(t, filepath.Join(inner, adrMarker), "inner\n")

	dir, err := New(inner, kind).Dir(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if want := "inner"; dir != want {
		t.Errorf("Dir() = %q, want %q", dir, want)
	}
}

func TestSetDirCreatesTheDirectoryAndTheMarker(t *testing.T) {
	s := newStore(t)
	if err := s.SetDir(context.Background(), "architecture-log"); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(s.path("architecture-log")); err != nil || !info.IsDir() {
		t.Errorf("the log directory was not created: %v", err)
	}
	dir, err := s.Dir(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if want := "architecture-log"; dir != want {
		t.Errorf("Dir() = %q, want %q", dir, want)
	}
}

func TestRecordsReturnsOnlyRecordsInOrder(t *testing.T) {
	s := newStore(t)
	writeRecord(t, s, "docs/adr/0002-second.md", "# 2. Second\n")
	writeRecord(t, s, "docs/adr/0001-first.md", "# 1. First\n")
	writeRecord(t, s, "docs/adr/README.md", "not a record\n")
	writeRecord(t, s, "docs/adr/templates/template.md", "not a record\n")

	records, err := s.Records(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := []record.Record{
		{Number: 1, File: "docs/adr/0001-first.md"},
		{Number: 2, File: "docs/adr/0002-second.md"},
	}
	if !reflect.DeepEqual(records, want) {
		t.Errorf("Records() = %v, want %v", records, want)
	}
}

func TestRecordsReportsAMissingDirectory(t *testing.T) {
	_, err := newStore(t).Records(context.Background())

	var missing *record.MissingDirError
	if !errors.As(err, &missing) {
		t.Fatalf("Records() error = %v, want an *record.MissingDirError", err)
	}
	if missing.Dir != adrDir {
		t.Errorf("missing directory = %q, want %q", missing.Dir, adrDir)
	}
}

func TestResolveMatchesNumberOrName(t *testing.T) {
	s := newStore(t)
	writeRecord(t, s, "docs/adr/0001-first.md", "# 1. First\n")
	writeRecord(t, s, "docs/adr/0002-second.md", "# 2. Second\n")

	for _, test := range []struct{ ref, want string }{
		{"2", "docs/adr/0002-second.md"},
		{"second", "docs/adr/0002-second.md"},
		{"docs/adr/0001-first.md", "docs/adr/0001-first.md"},
	} {
		got, err := s.Resolve(context.Background(), test.ref)
		if err != nil {
			t.Fatalf("Resolve(%q): %v", test.ref, err)
		}
		if got.File != test.want {
			t.Errorf("Resolve(%q) = %q, want %q", test.ref, got.File, test.want)
		}
	}
	if _, err := s.Resolve(context.Background(), "nothing"); err == nil {
		t.Error(`Resolve("nothing") succeeded, want an error`)
	}
}

func TestNextNumberFollowsTheHighestRecord(t *testing.T) {
	s := newStore(t)
	if got, err := s.NextNumber(context.Background(), "docs/adr"); err != nil || got != 1 {
		t.Errorf("NextNumber() on an absent log = %d, %v, want 1, nil", got, err)
	}

	writeRecord(t, s, "docs/adr/0001-first.md", "# 1. First\n")
	writeRecord(t, s, "docs/adr/0009-ninth.md", "# 9. Ninth\n")
	if got, err := s.NextNumber(context.Background(), "docs/adr"); err != nil || got != 10 {
		t.Errorf("NextNumber() = %d, %v, want 10, nil", got, err)
	}
}

func TestLinesRoundTripsThroughWriteLines(t *testing.T) {
	s := newStore(t)
	writeRecord(t, s, "docs/adr/0001-first.md", "# 1. First\n\n## Status\n\nAccepted\n")

	lines, err := s.Lines(context.Background(), "docs/adr/0001-first.md")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.WriteLines(context.Background(), "docs/adr/0001-first.md", lines); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(s.path("docs/adr/0001-first.md"))
	if err != nil {
		t.Fatal(err)
	}
	if want := "# 1. First\n\n## Status\n\nAccepted\n"; string(content) != want {
		t.Errorf("round trip produced %q, want %q", content, want)
	}
}

func TestWriteCreatesTheLogDirectory(t *testing.T) {
	s := newStore(t)
	if err := s.Write(context.Background(), "docs/adr/0001-first.md", "# 1. First\n"); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(s.path("docs/adr/0001-first.md"))
	if err != nil {
		t.Fatal(err)
	}
	if want := "# 1. First\n"; string(content) != want {
		t.Errorf("Write() wrote %q, want %q", content, want)
	}
}

func mkdirAll(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
}

func write(t *testing.T, file, content string) {
	t.Helper()
	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
