package initialise

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/mustardcustard/detoo/internal/adapters/memlog"
	"github.com/mustardcustard/detoo/internal/adr"
	"github.com/mustardcustard/detoo/internal/rfc"
	"github.com/mustardcustard/detoo/internal/ust"
)

const firstTemplate = "# 1. Record architecture decisions\n\nDate: DATE\n\n## Status\n\nAccepted\n"

type stubTemplates struct{}

func (stubTemplates) First(context.Context) (string, error) { return firstTemplate, nil }

type stubClock struct{}

func (stubClock) Today(context.Context) string { return "1992-01-12" }

func newHandler(log *memlog.Log) Handler {
	return Handler{Kind: adr.Kind, Log: log, Templates: stubTemplates{}, Clock: stubClock{}}
}

// The first record is filed under the title the kind of log gives it.
func TestHandleFilesTheFirstRecordForTheKind(t *testing.T) {
	handler := newHandler(memlog.New())
	handler.Kind = rfc.Kind

	result, err := handler.Handle(context.Background(), Command{Dir: "docs/rfc"})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	if want := "docs/rfc/RFC-0001-use-requests-for-comments.md"; result.File != want {
		t.Errorf("File = %q, want %q", result.File, want)
	}
}

// The kind of log names the default directory of the first record, and the
// prefix of its file name.
func TestHandleFilesTheFirstRecordInTheKindsDirectory(t *testing.T) {
	handler := newHandler(memlog.Of(ust.Kind))
	handler.Kind = ust.Kind

	result, err := handler.Handle(context.Background(), Command{})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	if want := "docs/userstory/STORY-0001-describe-work-in-user-stories.md"; result.File != want {
		t.Errorf("File = %q, want %q", result.File, want)
	}
}

func TestHandleWritesTheFirstRecordInTheDefaultDirectory(t *testing.T) {
	log := memlog.New()

	result, err := newHandler(log).Handle(context.Background(), Command{})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if want := "docs/adr/ADR-0001-record-architecture-decisions.md"; result.File != want {
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

	if want := "architecture-log/ADR-0001-record-architecture-decisions.md"; result.File != want {
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
	if want := "architecture-log/ADR-0001-record-architecture-decisions.md"; result.File != want {
		t.Errorf("File = %q, want %q", result.File, want)
	}
}

// existingLog returns a log that has already been started.
func existingLog() *memlog.Log {
	return memlog.WithRecords(map[string]string{
		"docs/adr/0001-record-architecture-decisions.md": "# 1. Record architecture decisions\n",
	})
}

// Running init again leaves a log that already has records as it is, whether
// or not it is told the log's own directory.
func TestHandleDoesNothingWhenTheLogHasRecords(t *testing.T) {
	for _, test := range []struct {
		name string
		dir  string
	}{
		{"no directory", ""},
		{"the log's directory", "docs/adr"},
		{"the log's directory, untidily", "docs/adr/"},
	} {
		t.Run(test.name, func(t *testing.T) {
			log := existingLog()

			result, err := newHandler(log).Handle(context.Background(), Command{Dir: test.dir})
			if err != nil {
				t.Fatalf("Handle: %v", err)
			}
			if result.File != "" {
				t.Errorf("File = %q, want nothing written", result.File)
			}
			want := []string{"docs/adr/0001-record-architecture-decisions.md"}
			if got := log.Files(); !reflect.DeepEqual(got, want) {
				t.Errorf("the log holds %q, want %q", got, want)
			}
		})
	}
}

// Pointing init at another directory would start a second log and leave the
// first on disk, no longer used, so it is refused and nothing changes.
func TestHandleRefusesToStartASecondLog(t *testing.T) {
	log := existingLog()

	_, err := newHandler(log).Handle(context.Background(), Command{Dir: "other"})
	if !errors.Is(err, ErrLogExists) {
		t.Errorf("err = %v, want %v", err, ErrLogExists)
	}
	dir, err := log.Dir(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if want := "docs/adr"; dir != want {
		t.Errorf("the log directory is %q, want %q", dir, want)
	}
	want := []string{"docs/adr/0001-record-architecture-decisions.md"}
	if got := log.Files(); !reflect.DeepEqual(got, want) {
		t.Errorf("the log holds %q, want %q", got, want)
	}
}

// A log directory with no records in it has not been started, so init starts
// it: what matters is the records, not whether the directory exists.
func TestHandleStartsALogInAnEmptyDirectory(t *testing.T) {
	result, err := newHandler(memlog.WithRecords(map[string]string{})).Handle(context.Background(), Command{})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	if want := "docs/adr/ADR-0001-record-architecture-decisions.md"; result.File != want {
		t.Errorf("File = %q, want %q", result.File, want)
	}
}
