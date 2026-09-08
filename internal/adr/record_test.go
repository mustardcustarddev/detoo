package adr

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// testEnv returns an Env on an empty directory, with output discarded.
func testEnv(t *testing.T) *Env {
	t.Helper()
	return &Env{
		Dir:    t.TempDir(),
		Stdout: io.Discard,
		Stderr: io.Discard,
		Getenv: func(string) string { return "" },
	}
}

// writeRecord creates a record, and the directories leading to it.
func writeRecord(t *testing.T, e *Env, name, content string) {
	t.Helper()
	file := e.path(name)
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestSlug(t *testing.T) {
	for _, test := range []struct {
		title string
		want  string
	}{
		{"The First Decision", "the-first-decision"},
		{"Something About Node.JS", "something-about-node-js"},
		{"Slash/Slash/Slash/", "slash-slash-slash"},
		{"-Bar-", "bar"},
		{"  spaced   out  ", "spaced-out"},
		{"Use PostgreSQL 15", "use-postgresql-15"},
		{"!?!", ""},
	} {
		if got := Slug(test.title); got != test.want {
			t.Errorf("Slug(%q) = %q, want %q", test.title, got, test.want)
		}
	}
}

func TestNumber(t *testing.T) {
	for _, test := range []struct {
		file string
		want int
	}{
		{"doc/adr/0001-first.md", 1},
		{"doc/adr/0009-ninth.md", 9},
		{"doc/adr/0010-tenth.md", 10},
		{"doc/adr/templates", 0},
	} {
		if got := Number(test.file); got != test.want {
			t.Errorf("Number(%q) = %d, want %d", test.file, got, test.want)
		}
	}
}

func TestListReturnsOnlyRecords(t *testing.T) {
	e := testEnv(t)
	writeRecord(t, e, "doc/adr/0002-second.md", "# 2. Second\n")
	writeRecord(t, e, "doc/adr/0001-first.md", "# 1. First\n")
	writeRecord(t, e, "doc/adr/README.md", "not a record\n")
	writeRecord(t, e, "doc/adr/templates/template.md", "not a record\n")

	records, err := e.List()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"doc/adr/0001-first.md", "doc/adr/0002-second.md"}
	if !reflect.DeepEqual(records, want) {
		t.Errorf("List() = %q, want %q", records, want)
	}
}

func TestListReportsAMissingDirectory(t *testing.T) {
	e := testEnv(t)
	_, err := e.List()
	var missing *MissingDirError
	if !errors.As(err, &missing) {
		t.Fatalf("List() error = %v, want a *MissingDirError", err)
	}
	if missing.Dir != DefaultDir {
		t.Errorf("missing directory = %q, want %q", missing.Dir, DefaultDir)
	}
}

func TestResolveMatchesNumberOrName(t *testing.T) {
	e := testEnv(t)
	writeRecord(t, e, "doc/adr/0001-first.md", "# 1. First\n")
	writeRecord(t, e, "doc/adr/0002-second.md", "# 2. Second\n")

	for _, test := range []struct{ ref, want string }{
		{"2", "doc/adr/0002-second.md"},
		{"second", "doc/adr/0002-second.md"},
		{"doc/adr/0001-first.md", "doc/adr/0001-first.md"},
	} {
		got, err := e.Resolve(test.ref)
		if err != nil {
			t.Fatalf("Resolve(%q): %v", test.ref, err)
		}
		if got != test.want {
			t.Errorf("Resolve(%q) = %q, want %q", test.ref, got, test.want)
		}
	}
	if _, err := e.Resolve("nothing"); err == nil {
		t.Error("Resolve(\"nothing\") succeeded, want an error")
	}
}

func TestTitle(t *testing.T) {
	e := testEnv(t)
	writeRecord(t, e, "doc/adr/0001-first.md", "# 1. First Record\n\nDate: 1992-01-12\n")
	title, err := e.Title("1")
	if err != nil {
		t.Fatal(err)
	}
	if want := "1. First Record"; title != want {
		t.Errorf("Title(\"1\") = %q, want %q", title, want)
	}
}
