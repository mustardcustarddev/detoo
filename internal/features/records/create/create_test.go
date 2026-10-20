package create

import (
	"context"
	"errors"
	"testing"

	"github.com/mustardcustard/detoo/internal/adapters/memlog"
	"github.com/mustardcustard/detoo/internal/adr"
	"github.com/mustardcustard/detoo/internal/rfc"
	"github.com/mustardcustard/detoo/internal/ust"
)

// Nothing here builds a binary, parses a flag or writes a file.  The slice's
// ports are narrow enough that the whole of creating a record — numbering,
// templating, superceding, linking, and opening an editor — can be exercised
// from memory.

const template = "# NUMBER. TITLE\n\nDate: DATE\n\n## Status\n\nSTATUS\n\n## Context\n\nContext.\n"

type stubTemplates struct {
	record    string
	directory map[string]string
}

func (s stubTemplates) Record(context.Context, string) (string, error) { return s.record, nil }

func (s stubTemplates) Directory(context.Context) (map[string]string, error) {
	return s.directory, nil
}

type stubClock struct{}

func (stubClock) Today(context.Context) string { return "1992-01-12" }

// spyEditor records the file it was asked to open, so that a test can check
// the record is offered to the user without an editor ever running.
type spyEditor struct {
	edited string
	err    error
}

func (e *spyEditor) Edit(_ context.Context, file string) error {
	e.edited = file
	return e.err
}

func newHandler(log *memlog.Log, editor *spyEditor) Handler {
	return Handler{
		Kind:      adr.Kind,
		Log:       log,
		Templates: stubTemplates{record: template},
		Clock:     stubClock{},
		Editor:    editor,
	}
}

func TestHandleCreatesTheNextNumberedRecord(t *testing.T) {
	log := memlog.WithRecords(map[string]string{
		"docs/adr/0001-first.md": "# 1. First\n\n## Status\n\nAccepted\n\n## Context\n\nContext.\n",
	})
	editor := &spyEditor{}

	result, err := newHandler(log, editor).Handle(context.Background(), Command{Title: "Use MySQL Database"})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if want := "docs/adr/ADR-0002-use-mysql-database.md"; result.File != want {
		t.Errorf("File = %q, want %q", result.File, want)
	}
	want := "# 2. Use MySQL Database\n\nDate: 1992-01-12\n\n## Status\n\nDraft\n\n## Context\n\nContext.\n"
	if got := log.Content(result.File); got != want {
		t.Errorf("record is\n%s\nwant\n%s", got, want)
	}
	if editor.edited != result.File {
		t.Errorf("the editor was asked for %q, want %q", editor.edited, result.File)
	}
}

// A kind whose records have files of their own, such as a fire drill's pages,
// creates them in the directory beside the record, filled in as it is.
func TestHandleWritesTheFilesBesideTheRecord(t *testing.T) {
	log := memlog.New()
	handler := newHandler(log, &spyEditor{})
	handler.Templates = stubTemplates{record: template, directory: map[string]string{
		"problem/README.md":  "# Problem\n\nThe pages of NUMBER. TITLE.\n",
		"solution/README.md": "# Solution\n",
	}}

	result, err := handler.Handle(context.Background(), Command{Title: "Use MySQL Database"})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	want := map[string]string{
		"docs/adr/ADR-0001-use-mysql-database/problem/README.md":  "# Problem\n\nThe pages of 1. Use MySQL Database.\n",
		"docs/adr/ADR-0001-use-mysql-database/solution/README.md": "# Solution\n",
	}
	for file, content := range want {
		if got := log.Content(file); got != content {
			t.Errorf("%s is %q, want %q", file, got, content)
		}
	}
	if got, want := len(log.Files()), 3; got != want {
		t.Errorf("wrote %d files, want %d: %v", got, want, log.Files())
	}
	if result.File != "docs/adr/ADR-0001-use-mysql-database.md" {
		t.Errorf("File = %q, want the record", result.File)
	}
}

func TestHandleRejectsARequestWithNoTitle(t *testing.T) {
	_, err := newHandler(memlog.New(), &spyEditor{}).Handle(context.Background(), Command{})
	if !errors.Is(err, ErrNoTitle) {
		t.Errorf("err = %v, want %v", err, ErrNoTitle)
	}
}

func TestHandleSupersedesAnExistingRecord(t *testing.T) {
	log := memlog.WithRecords(map[string]string{
		"docs/adr/0001-first.md": "# 1. First\n\n## Status\n\nAccepted\n\n## Context\n\nContext.\n",
	})

	result, err := newHandler(log, &spyEditor{}).Handle(context.Background(), Command{
		Title:    "Second",
		Replaces: []string{"1"},
	})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	// The superseded record gains the link and loses its accepted status.
	wantOld := "# 1. First\n\n## Status\n\nSuperseded by [2. Second](ADR-0002-second.md)\n\n## Context\n\nContext.\n"
	if got := log.Content("docs/adr/0001-first.md"); got != wantOld {
		t.Errorf("superseded record is\n%s\nwant\n%s", got, wantOld)
	}
	// The new record records what it replaced.
	wantNew := "# 2. Second\n\nDate: 1992-01-12\n\n## Status\n\nDraft\n\nSupersedes [1. First](0001-first.md)\n\n## Context\n\nContext.\n"
	if got := log.Content(result.File); got != wantNew {
		t.Errorf("new record is\n%s\nwant\n%s", got, wantNew)
	}
}

// A decision superseded while it is still a draft is no longer one.
func TestHandleSupersedesADraft(t *testing.T) {
	log := memlog.WithRecords(map[string]string{
		"docs/adr/0001-first.md": "# 1. First\n\n## Status\n\nDraft\n\n## Context\n\nContext.\n",
	})

	if _, err := newHandler(log, &spyEditor{}).Handle(context.Background(), Command{
		Title:    "Second",
		Replaces: []string{"1"},
	}); err != nil {
		t.Fatalf("Handle: %v", err)
	}

	wantOld := "# 1. First\n\n## Status\n\nSuperseded by [2. Second](ADR-0002-second.md)\n\n## Context\n\nContext.\n"
	if got := log.Content("docs/adr/0001-first.md"); got != wantOld {
		t.Errorf("superseded record is\n%s\nwant\n%s", got, wantOld)
	}
}

func TestHandleUpdatesAnExistingRecord(t *testing.T) {
	log := memlog.WithRecords(map[string]string{
		"docs/adr/0001-first.md": "# 1. First\n\n## Status\n\nAccepted\n\n## Context\n\nContext.\n",
	})

	result, err := newHandler(log, &spyEditor{}).Handle(context.Background(), Command{
		Title:   "Second",
		Updates: []string{"1"},
	})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	// The updated record gains the link, and is still accepted.
	wantOld := "# 1. First\n\n## Status\n\nAccepted\n\nUpdated by [2. Second](ADR-0002-second.md)\n\n## Context\n\nContext.\n"
	if got := log.Content("docs/adr/0001-first.md"); got != wantOld {
		t.Errorf("updated record is\n%s\nwant\n%s", got, wantOld)
	}
	wantNew := "# 2. Second\n\nDate: 1992-01-12\n\n## Status\n\nDraft\n\nUpdates [1. First](0001-first.md)\n\n## Context\n\nContext.\n"
	if got := log.Content(result.File); got != wantNew {
		t.Errorf("new record is\n%s\nwant\n%s", got, wantNew)
	}
}

// The status a record starts with, the words its links are written in, and the
// status a replaced record loses all belong to the kind of log.
func TestHandleWritesInTheWordsOfTheKind(t *testing.T) {
	log := memlog.Of(rfc.Kind)
	if err := log.Write(context.Background(), "docs/rfc/0001-first.md", "# 1. First\n\n## Status\n\nDraft\n\n## Context\n\nContext.\n"); err != nil {
		t.Fatal(err)
	}
	handler := newHandler(log, &spyEditor{})
	handler.Kind = rfc.Kind

	result, err := handler.Handle(context.Background(), Command{Title: "Second", Replaces: []string{"1"}})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	wantOld := "# 1. First\n\n## Status\n\nObsoleted by [2. Second](RFC-0002-second.md)\n\n## Context\n\nContext.\n"
	if got := log.Content("docs/rfc/0001-first.md"); got != wantOld {
		t.Errorf("obsoleted record is\n%s\nwant\n%s", got, wantOld)
	}
	wantNew := "# 2. Second\n\nDate: 1992-01-12\n\n## Status\n\nDraft\n\nObsoletes [1. First](0001-first.md)\n\n## Context\n\nContext.\n"
	if got := log.Content(result.File); got != wantNew {
		t.Errorf("new record is\n%s\nwant\n%s", got, wantNew)
	}
}

// A kind with a prefix begins each new record's file name with it, in upper
// case, numbers its records from those named in any of its ways, such as the
// lower-case story-0001-first.md, and links them by those names.
func TestHandleNamesRecordsWithTheKindsPrefix(t *testing.T) {
	log := memlog.Of(ust.Kind)
	if err := log.Write(context.Background(), "docs/userstory/story-0001-first.md", "# 1. First\n\n## Status\n\nDraft\n\n## Story\n\nStory.\n"); err != nil {
		t.Fatal(err)
	}
	if err := log.Write(context.Background(), "docs/userstory/epic-0009-not-a-story.md", "# Not a story\n"); err != nil {
		t.Fatal(err)
	}
	handler := newHandler(log, &spyEditor{})
	handler.Kind = ust.Kind

	result, err := handler.Handle(context.Background(), Command{Title: "Second", Updates: []string{"1"}})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if want := "docs/userstory/STORY-0002-second.md"; result.File != want {
		t.Errorf("File = %q, want %q", result.File, want)
	}
	wantOld := "# 1. First\n\n## Status\n\nDraft\n\nUpdated by [2. Second](STORY-0002-second.md)\n\n## Story\n\nStory.\n"
	if got := log.Content("docs/userstory/story-0001-first.md"); got != wantOld {
		t.Errorf("updated record is\n%s\nwant\n%s", got, wantOld)
	}
}

// A record that cannot be found is reported rather than silently skipped.
func TestHandleReportsAnUnknownRecordToSupersede(t *testing.T) {
	log := memlog.WithRecords(map[string]string{
		"docs/adr/0001-first.md": "# 1. First\n\n## Status\n\nAccepted\n\n## Context\n\nContext.\n",
	})

	_, err := newHandler(log, &spyEditor{}).Handle(context.Background(), Command{
		Title:    "Second",
		Replaces: []string{"nothing-like-this"},
	})
	if err == nil {
		t.Error("Handle succeeded, want an error naming the missing record")
	}
}
