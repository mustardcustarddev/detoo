package link

import (
	"context"
	"errors"
	"testing"

	"github.com/mustardcustard/detoo/internal/adapters/memlog"
	"github.com/mustardcustard/detoo/internal/adr"
	"github.com/mustardcustard/detoo/internal/domain/record"
	"github.com/mustardcustard/detoo/internal/rfc"
)

const (
	first  = "# 1. First\n\n## Status\n\nAccepted\n\n## Context\n\nContext.\n"
	second = "# 2. Second\n\n## Status\n\nAccepted\n\n## Context\n\nContext.\n"
)

func newLog() *memlog.Log {
	return memlog.WithRecords(map[string]string{
		"docs/adr/0001-first.md":  first,
		"docs/adr/0002-second.md": second,
	})
}

// An update changes an earlier decision without replacing it, so the earlier
// record is still accepted.
func TestHandleUpdateLinksBothRecords(t *testing.T) {
	log := newLog()

	result, err := Handler{Kind: adr.Kind, Log: log}.Handle(context.Background(), Command{
		Relation: record.Updates,
		Source:   "2",
		Target:   "1",
	})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if result.Source != "docs/adr/0002-second.md" || result.Target != "docs/adr/0001-first.md" {
		t.Errorf("Result = %+v, want the two records", result)
	}
	wantSource := "# 2. Second\n\n## Status\n\nAccepted\n\nUpdates [1. First](0001-first.md)\n\n## Context\n\nContext.\n"
	if got := log.Content("docs/adr/0002-second.md"); got != wantSource {
		t.Errorf("source record is\n%s\nwant\n%s", got, wantSource)
	}
	wantTarget := "# 1. First\n\n## Status\n\nAccepted\n\nUpdated by [2. Second](0002-second.md)\n\n## Context\n\nContext.\n"
	if got := log.Content("docs/adr/0001-first.md"); got != wantTarget {
		t.Errorf("target record is\n%s\nwant\n%s", got, wantTarget)
	}
}

// A superceded record is no longer the decision, so it loses its accepted
// status, just as it does when a new record supercedes it.
func TestHandleSupercedeRetiresTheTarget(t *testing.T) {
	log := newLog()

	if _, err := (Handler{Kind: adr.Kind, Log: log}).Handle(context.Background(), Command{
		Relation: record.Replaces,
		Source:   "2",
		Target:   "1",
	}); err != nil {
		t.Fatalf("Handle: %v", err)
	}

	wantSource := "# 2. Second\n\n## Status\n\nAccepted\n\nSupercedes [1. First](0001-first.md)\n\n## Context\n\nContext.\n"
	if got := log.Content("docs/adr/0002-second.md"); got != wantSource {
		t.Errorf("source record is\n%s\nwant\n%s", got, wantSource)
	}
	wantTarget := "# 1. First\n\n## Status\n\nSuperceded by [2. Second](0002-second.md)\n\n## Context\n\nContext.\n"
	if got := log.Content("docs/adr/0001-first.md"); got != wantTarget {
		t.Errorf("target record is\n%s\nwant\n%s", got, wantTarget)
	}
}

// The words a link is written in, and the status a replaced record loses,
// belong to the kind of log: an RFC obsoletes a draft.
func TestHandleWritesInTheWordsOfTheKind(t *testing.T) {
	log := memlog.WithRecords(map[string]string{
		"docs/rfc/0001-first.md":  "# 1. First\n\n## Status\n\nDraft\n\n## Summary\n\nSummary.\n",
		"docs/rfc/0002-second.md": "# 2. Second\n\n## Status\n\nDraft\n\n## Summary\n\nSummary.\n",
	})
	if err := log.SetDir(context.Background(), "docs/rfc"); err != nil {
		t.Fatal(err)
	}

	if _, err := (Handler{Kind: rfc.Kind, Log: log}).Handle(context.Background(), Command{
		Relation: record.Replaces,
		Source:   "2",
		Target:   "1",
	}); err != nil {
		t.Fatalf("Handle: %v", err)
	}

	wantSource := "# 2. Second\n\n## Status\n\nDraft\n\nObsoletes [1. First](0001-first.md)\n\n## Summary\n\nSummary.\n"
	if got := log.Content("docs/rfc/0002-second.md"); got != wantSource {
		t.Errorf("source record is\n%s\nwant\n%s", got, wantSource)
	}
	wantTarget := "# 1. First\n\n## Status\n\nObsoleted by [2. Second](0002-second.md)\n\n## Summary\n\nSummary.\n"
	if got := log.Content("docs/rfc/0001-first.md"); got != wantTarget {
		t.Errorf("target record is\n%s\nwant\n%s", got, wantTarget)
	}
}

// A record is named by its number or by part of its file name, and either is
// resolved before anything is written.
func TestHandleResolvesRecordsByName(t *testing.T) {
	result, err := Handler{Kind: adr.Kind, Log: newLog()}.Handle(context.Background(), Command{
		Relation: record.Updates,
		Source:   "second",
		Target:   "first",
	})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	if result.Source != "docs/adr/0002-second.md" {
		t.Errorf("Source = %q, want the second record", result.Source)
	}
}

// A link that says neither how it bears on its target has nothing to write, and
// writes nothing.
func TestHandleRejectsALinkWithNoRelation(t *testing.T) {
	log := newLog()

	_, err := Handler{Kind: adr.Kind, Log: log}.Handle(context.Background(), Command{Source: "2", Target: "1"})
	if !errors.Is(err, ErrNoRelation) {
		t.Errorf("err = %v, want %v", err, ErrNoRelation)
	}
	if got := log.Content("docs/adr/0001-first.md"); got != first {
		t.Errorf("target record was changed to\n%s", got)
	}
	if got := log.Content("docs/adr/0002-second.md"); got != second {
		t.Errorf("source record was changed to\n%s", got)
	}
}

func TestHandleReportsAnUnknownRecord(t *testing.T) {
	if _, err := (Handler{Kind: adr.Kind, Log: newLog()}).Handle(context.Background(), Command{
		Relation: record.Updates,
		Source:   "2",
		Target:   "nothing-like-this",
	}); err == nil {
		t.Error("Handle succeeded, want an error naming the missing record")
	}
}
