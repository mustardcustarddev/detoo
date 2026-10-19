package logstore

import (
	"context"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// kind is the kind of log the other tests store: adr's, whose directory is
// adrDir.
var kind = record.Kind{Name: "adr", Dir: "adr", Noun: "decision record"}

const adrDir = "docs/adr"

// The kind of log names its directory, so that one project can keep the logs
// of several tools without either finding the other.  The directory need not
// be named after the tool.
func TestTheKindNamesTheLogsFiles(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, "docs", "adr"))
	s := New(root, record.Kind{Name: "rfc", Dir: "requests", Noun: "RFC"})

	dir, err := s.Dir(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if want := "docs/requests"; dir != want {
		t.Errorf("Dir() = %q, want %q", dir, want)
	}
}

// A kind whose records are prefixed sees the files named in any of its ways,
// in order of their numbers, and none that carry another prefix.
func TestAPrefixedKindSeesOnlyItsOwnRecords(t *testing.T) {
	s := New(t.TempDir(), record.Kind{Name: "ust", Dir: "userstory", Prefix: "STORY", Noun: "user story"})
	writeRecord(t, s, "docs/userstory/story-0001-first.md", "# 1. First\n")
	writeRecord(t, s, "docs/userstory/STORY-0002-second.md", "# 2. Second\n")
	writeRecord(t, s, "docs/userstory/0003-unprefixed.md", "# 3. Unprefixed\n")
	writeRecord(t, s, "docs/userstory/epic-0007-an-epic.md", "# 7. An epic\n")
	writeRecord(t, s, "docs/userstory/README.md", "# User Stories\n")

	records, err := s.Records(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := []record.Record{
		{Number: 1, File: "docs/userstory/story-0001-first.md"},
		{Number: 2, File: "docs/userstory/STORY-0002-second.md"},
		{Number: 3, File: "docs/userstory/0003-unprefixed.md"},
	}
	if !reflect.DeepEqual(records, want) {
		t.Errorf("Records() = %v, want %v", records, want)
	}
	if got, err := s.NextNumber(context.Background(), "docs/userstory"); err != nil || got != 4 {
		t.Errorf("NextNumber() = %d, %v, want 4, nil", got, err)
	}
}
