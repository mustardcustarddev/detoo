package initialise

import (
	"context"
	"testing"

	"github.com/mustardcustard/decision-records/internal/adapters/memlog"
)

const firstTemplate = "# 1. Record architecture decisions\n\nDate: DATE\n\n## Status\n\nAccepted\n"

type stubTemplates struct{}

func (stubTemplates) First(context.Context) (string, error) { return firstTemplate, nil }

type stubClock struct{}

func (stubClock) Today(context.Context) string { return "1992-01-12" }

func newHandler(log *memlog.Log) Handler {
	return Handler{Log: log, Templates: stubTemplates{}, Clock: stubClock{}}
}

func TestHandleWritesTheFirstRecordInTheDefaultDirectory(t *testing.T) {
	log := memlog.New()

	result, err := newHandler(log).Handle(context.Background(), Command{})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if want := "docs/adr/0001-record-architecture-decisions.md"; result.File != want {
		t.Errorf("File = %q, want %q", result.File, want)
	}
	want := "# 1. Record architecture decisions\n\nDate: 1992-01-12\n\n## Status\n\nAccepted\n"
	if got := log.Content(result.File); got != want {
		t.Errorf("record is\n%s\nwant\n%s", got, want)
	}
}

func TestHandleMovesTheLogWhenGivenADirectory(t *testing.T) {
	log := memlog.New()

	result, err := newHandler(log).Handle(context.Background(), Command{Dir: "architecture-log"})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if want := "architecture-log/0001-record-architecture-decisions.md"; result.File != want {
		t.Errorf("File = %q, want %q", result.File, want)
	}
	dir, err := log.Dir(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if want := "architecture-log"; dir != want {
		t.Errorf("the log directory is %q, want %q", dir, want)
	}
}

// An empty directory means "wherever this project already keeps its records",
// so the log is not moved back to the default.
func TestHandleLeavesTheDirectoryAloneWhenNoneIsGiven(t *testing.T) {
	log := memlog.New()
	if err := log.SetDir(context.Background(), "architecture-log"); err != nil {
		t.Fatal(err)
	}

	result, err := newHandler(log).Handle(context.Background(), Command{})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	if want := "architecture-log/0001-record-architecture-decisions.md"; result.File != want {
		t.Errorf("File = %q, want %q", result.File, want)
	}
}
