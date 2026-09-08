package adr

import (
	"reflect"
	"testing"
)

const recordWithLinks = `# 3. Third Record

Date: 1992-01-12

## Status

Accepted

Supercedes [1. First Record](0001-first-record.md)

Amended by [4. Fourth Record](0004-fourth-record.md)

## Context

Some context.

## Decision

Some decision.
`

func TestStatusReadsTheStatusSection(t *testing.T) {
	e := testEnv(t)
	writeRecord(t, e, "doc/adr/0003-third-record.md", recordWithLinks)

	status, err := e.Status("3")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"Accepted",
		"Supercedes [1. First Record](0001-first-record.md)",
		"Amended by [4. Fourth Record](0004-fourth-record.md)",
	}
	if !reflect.DeepEqual(status, want) {
		t.Errorf("Status(\"3\") = %q, want %q", status, want)
	}
}

func TestLinksReadsLinksFromTheStatusSection(t *testing.T) {
	e := testEnv(t)
	writeRecord(t, e, "doc/adr/0003-third-record.md", recordWithLinks)

	links, err := e.Links("3")
	if err != nil {
		t.Fatal(err)
	}
	want := []Link{{Type: "Supercedes", Target: 1}, {Type: "Amended by", Target: 4}}
	if !reflect.DeepEqual(links, want) {
		t.Errorf("Links(\"3\") = %v, want %v", links, want)
	}
}

func TestAddLinkAppendsToTheStatusSection(t *testing.T) {
	e := testEnv(t)
	writeRecord(t, e, "doc/adr/0001-first.md", "# 1. First\n\n## Status\n\nAccepted\n\n## Context\n\nContext.\n")
	writeRecord(t, e, "doc/adr/0002-second.md", "# 2. Second\n\n## Status\n\nAccepted\n")

	if err := e.AddLink("1", "Amended by", "2"); err != nil {
		t.Fatal(err)
	}
	want := "# 1. First\n\n## Status\n\nAccepted\n\nAmended by [2. Second](0002-second.md)\n\n## Context\n\nContext.\n"
	if got := readRecord(t, e, "doc/adr/0001-first.md"); got != want {
		t.Errorf("record is\n%s\nwant\n%s", got, want)
	}
}

func TestRemoveStatusClosesTheGapItLeaves(t *testing.T) {
	e := testEnv(t)
	writeRecord(t, e, "doc/adr/0001-first.md",
		"# 1. First\n\n## Status\n\nAccepted\n\nSuperceded by [2. Second](0002-second.md)\n\n## Context\n\nAccepted\n")

	if err := e.RemoveStatus("Accepted", "1"); err != nil {
		t.Fatal(err)
	}
	// Only the status section is touched: the "Accepted" in the context stays.
	want := "# 1. First\n\n## Status\n\nSuperceded by [2. Second](0002-second.md)\n\n## Context\n\nAccepted\n"
	if got := readRecord(t, e, "doc/adr/0001-first.md"); got != want {
		t.Errorf("record is\n%s\nwant\n%s", got, want)
	}
}

func readRecord(t *testing.T, e *Env, name string) string {
	t.Helper()
	content, err := readFile(e.path(name))
	if err != nil {
		t.Fatal(err)
	}
	return content
}
