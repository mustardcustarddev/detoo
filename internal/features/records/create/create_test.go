package create

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/mustardcustard/decision-records/internal/adapters/memlog"
)

// Nothing here builds a binary, parses a flag or writes a file.  The slice's
// ports are narrow enough that the whole of creating a record — numbering,
// templating, superceding, linking, and opening an editor — can be exercised
// from memory.

const template = "# NUMBER. TITLE\n\nDate: DATE\n\n## Status\n\nSTATUS\n\n## Context\n\nContext.\n"

type stubTemplates struct{ record string }

func (s stubTemplates) Record(context.Context, string) (string, error) { return s.record, nil }

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

	if want := "docs/adr/0002-use-mysql-database.md"; result.File != want {
		t.Errorf("File = %q, want %q", result.File, want)
	}
	want := "# 2. Use MySQL Database\n\nDate: 1992-01-12\n\n## Status\n\nAccepted\n\n## Context\n\nContext.\n"
	if got := log.Content(result.File); got != want {
		t.Errorf("record is\n%s\nwant\n%s", got, want)
	}
	if editor.edited != result.File {
		t.Errorf("the editor was asked for %q, want %q", editor.edited, result.File)
	}
}

func TestHandleRejectsARequestWithNoTitle(t *testing.T) {
	_, err := newHandler(memlog.New(), &spyEditor{}).Handle(context.Background(), Command{})
	if !errors.Is(err, ErrNoTitle) {
		t.Errorf("err = %v, want %v", err, ErrNoTitle)
	}
}

func TestHandleSupercedesAnExistingRecord(t *testing.T) {
	log := memlog.WithRecords(map[string]string{
		"docs/adr/0001-first.md": "# 1. First\n\n## Status\n\nAccepted\n\n## Context\n\nContext.\n",
	})

	result, err := newHandler(log, &spyEditor{}).Handle(context.Background(), Command{
		Title:      "Second",
		Supercedes: []string{"1"},
	})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	// The superceded record gains the link and loses its accepted status.
	wantOld := "# 1. First\n\n## Status\n\nSuperceded by [2. Second](0002-second.md)\n\n## Context\n\nContext.\n"
	if got := log.Content("docs/adr/0001-first.md"); got != wantOld {
		t.Errorf("superceded record is\n%s\nwant\n%s", got, wantOld)
	}
	// The new record records what it replaced.
	wantNew := "# 2. Second\n\nDate: 1992-01-12\n\n## Status\n\nAccepted\n\nSupercedes [1. First](0001-first.md)\n\n## Context\n\nContext.\n"
	if got := log.Content(result.File); got != wantNew {
		t.Errorf("new record is\n%s\nwant\n%s", got, wantNew)
	}
}

func TestHandleLinksBothWays(t *testing.T) {
	log := memlog.WithRecords(map[string]string{
		"docs/adr/0001-first.md": "# 1. First\n\n## Status\n\nAccepted\n\n## Context\n\nContext.\n",
	})

	result, err := newHandler(log, &spyEditor{}).Handle(context.Background(), Command{
		Title: "Second",
		Links: []Link{{Target: "1", Forward: "Amends", Reverse: "Amended by"}},
	})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if got, want := log.Content(result.File), "Amends [1. First](0001-first.md)"; !strings.Contains(got, want) {
		t.Errorf("new record does not contain %q:\n%s", want, got)
	}
	if got, want := log.Content("docs/adr/0001-first.md"), "Amended by [2. Second](0002-second.md)"; !strings.Contains(got, want) {
		t.Errorf("target record does not contain %q:\n%s", want, got)
	}
}

// A record that cannot be found is reported rather than silently skipped.
func TestHandleReportsAnUnknownRecordToSupercede(t *testing.T) {
	log := memlog.WithRecords(map[string]string{
		"docs/adr/0001-first.md": "# 1. First\n\n## Status\n\nAccepted\n\n## Context\n\nContext.\n",
	})

	_, err := newHandler(log, &spyEditor{}).Handle(context.Background(), Command{
		Title:      "Second",
		Supercedes: []string{"nothing-like-this"},
	})
	if err == nil {
		t.Error("Handle succeeded, want an error naming the missing record")
	}
}
