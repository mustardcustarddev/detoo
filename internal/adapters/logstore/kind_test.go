package logstore

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/mustardcustard/decision-records/internal/domain/record"
)

// kind is the kind of log the other tests store: adr's, whose marker file and
// default directory are these.
var kind = record.Kind{Name: "adr", Noun: "decision record"}

const (
	adrMarker = ".adr-dir"
	adrDir    = "docs/adr"
)

// The kind of log names its marker file and its default directory, so that one
// project can keep the logs of several tools without either finding the other.
func TestTheKindNamesTheLogsFiles(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, "docs", "adr"))
	s := New(root, record.Kind{Name: "rfc", Noun: "RFC"})

	dir, err := s.Dir(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if want := "docs/rfc"; dir != want {
		t.Errorf("Dir() = %q, want %q", dir, want)
	}

	if err := s.SetDir(context.Background(), "proposals"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, ".rfc-dir")); err != nil {
		t.Errorf("SetDir did not write .rfc-dir: %v", err)
	}
}
