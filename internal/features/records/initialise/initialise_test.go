package initialise

import (
	"context"
	"reflect"
	"testing"

	"github.com/mustardcustard/detoo/internal/adapters/memlog"
	"github.com/mustardcustard/detoo/internal/adr"
	"github.com/mustardcustard/detoo/internal/fyr"
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
	handler := newHandler(memlog.Of(rfc.Kind))
	handler.Kind = rfc.Kind

	result, err := handler.Handle(context.Background(), Command{})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	if want := "docs/rfc/RFC-0001-use-requests-for-comments.md"; result.File != want {
		t.Errorf("File = %q, want %q", result.File, want)
	}
}

// The kind of log names the directory of the first record, and the prefix of
// its file name.
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

// existingLog returns a log that has already been started.
func existingLog() *memlog.Log {
	return memlog.WithRecords(map[string]string{
		"docs/adr/0001-record-architecture-decisions.md": "# 1. Record architecture decisions\n",
	})
}

// Running init again leaves a log that already has records as it is.
func TestHandleDoesNothingWhenTheLogHasRecords(t *testing.T) {
	log := existingLog()

	result, err := newHandler(log).Handle(context.Background(), Command{})
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

// A kind with no first record opens its log with a README, written as the
// template is, so that the first record the user writes is number 1.
func TestHandleOpensALogWithAReadmeForAKindWithNoFirstRecord(t *testing.T) {
	log := memlog.Of(fyr.Kind)
	handler := newHandler(log)
	handler.Kind = fyr.Kind

	result, err := handler.Handle(context.Background(), Command{})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	if want := "docs/firedrill/README.md"; result.File != want {
		t.Errorf("File = %q, want %q", result.File, want)
	}
	if got := log.Content(result.File); got != firstTemplate {
		t.Errorf("README is\n%s\nwant\n%s", got, firstTemplate)
	}
}

// Running init again leaves a README as it is, whether init wrote it or the
// project did, such as a table of contents.
func TestHandleLeavesAReadmeAsItIs(t *testing.T) {
	log := memlog.Of(fyr.Kind)
	log.Write(context.Background(), "docs/firedrill/README.md", "# Our drills\n")
	handler := newHandler(log)
	handler.Kind = fyr.Kind

	result, err := handler.Handle(context.Background(), Command{})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	if result.File != "" {
		t.Errorf("File = %q, want nothing written", result.File)
	}
	if got, want := log.Content("docs/firedrill/README.md"), "# Our drills\n"; got != want {
		t.Errorf("README is %q, want %q", got, want)
	}
}
