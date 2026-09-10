package link

import (
	"context"
	"testing"

	"github.com/mustardcustard/decision-records/internal/adapters/memlog"
)

func newLog() *memlog.Log {
	return memlog.WithRecords(map[string]string{
		"docs/adr/0001-first.md":  "# 1. First\n\n## Status\n\nAccepted\n\n## Context\n\nContext.\n",
		"docs/adr/0002-second.md": "# 2. Second\n\n## Status\n\nAccepted\n\n## Context\n\nContext.\n",
	})
}

func TestHandleRecordsTheLinkOnBothRecords(t *testing.T) {
	log := newLog()

	result, err := Handler{Log: log}.Handle(context.Background(), Command{
		Source:  "1",
		Forward: "Amends",
		Target:  "2",
		Reverse: "Amended by",
	})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if result.Source != "docs/adr/0001-first.md" || result.Target != "docs/adr/0002-second.md" {
		t.Errorf("Result = %+v, want the two records", result)
	}
	wantSource := "# 1. First\n\n## Status\n\nAccepted\n\nAmends [2. Second](0002-second.md)\n\n## Context\n\nContext.\n"
	if got := log.Content("docs/adr/0001-first.md"); got != wantSource {
		t.Errorf("source record is\n%s\nwant\n%s", got, wantSource)
	}
	wantTarget := "# 2. Second\n\n## Status\n\nAccepted\n\nAmended by [1. First](0001-first.md)\n\n## Context\n\nContext.\n"
	if got := log.Content("docs/adr/0002-second.md"); got != wantTarget {
		t.Errorf("target record is\n%s\nwant\n%s", got, wantTarget)
	}
}

// A record is named by its number or by part of its file name, and either is
// resolved before anything is written.
func TestHandleResolvesRecordsByName(t *testing.T) {
	log := newLog()

	result, err := Handler{Log: log}.Handle(context.Background(), Command{
		Source:  "first",
		Forward: "Clarifies",
		Target:  "second",
		Reverse: "Clarified by",
	})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	if result.Source != "docs/adr/0001-first.md" {
		t.Errorf("Source = %q, want the first record", result.Source)
	}
}

func TestHandleReportsAnUnknownRecord(t *testing.T) {
	if _, err := (Handler{Log: newLog()}).Handle(context.Background(), Command{
		Source:  "1",
		Forward: "Amends",
		Target:  "nothing-like-this",
		Reverse: "Amended by",
	}); err == nil {
		t.Error("Handle succeeded, want an error naming the missing record")
	}
}
