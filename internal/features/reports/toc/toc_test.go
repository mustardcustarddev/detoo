package toc

import (
	"context"
	"strings"
	"testing"

	"github.com/mustardcustard/detoo/internal/adapters/memlog"
	"github.com/mustardcustard/detoo/internal/adr"
	"github.com/mustardcustard/detoo/internal/rfc"
)

func newLog() *memlog.Log {
	return memlog.WithRecords(map[string]string{
		"docs/adr/0001-first.md":  "# 1. First Decision\n",
		"docs/adr/0002-second.md": "# 2. Second Decision\n",
	})
}

func TestHandleWritesAListOfRecords(t *testing.T) {
	result, err := Handler{Kind: adr.Kind, Log: newLog()}.Handle(context.Background(), Query{})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	want := "# Architecture Decision Records\n\n" +
		"* [1. First Decision](0001-first.md)\n" +
		"* [2. Second Decision](0002-second.md)\n"
	if result.Markdown != want {
		t.Errorf("Markdown =\n%s\nwant\n%s", result.Markdown, want)
	}
}

// The contents are headed with what the kind of log holds.
func TestHandleHeadsTheListForTheKind(t *testing.T) {
	result, err := Handler{Kind: rfc.Kind, Log: newLog()}.Handle(context.Background(), Query{})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	if want := "# Requests for Comments\n\n"; !strings.HasPrefix(result.Markdown, want) {
		t.Errorf("Markdown =\n%s\nwant it to begin %q", result.Markdown, want)
	}
}

// The prefix lets the contents be published somewhere other than beside the
// records it links to.
func TestHandlePrefixesTheLinks(t *testing.T) {
	result, err := Handler{Kind: adr.Kind, Log: newLog()}.Handle(context.Background(), Query{LinkPrefix: "foo/docs/adr/"})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	want := "# Architecture Decision Records\n\n" +
		"* [1. First Decision](foo/docs/adr/0001-first.md)\n" +
		"* [2. Second Decision](foo/docs/adr/0002-second.md)\n"
	if result.Markdown != want {
		t.Errorf("Markdown =\n%s\nwant\n%s", result.Markdown, want)
	}
}

// A log with no records has no contents, so nothing is written — not even the
// heading, or an intro and outro, which are not read.  The intro and outro
// named here do not exist, and asking for them is not an error.
func TestHandleWritesNothingForAnEmptyLog(t *testing.T) {
	for _, test := range []struct {
		name  string
		query Query
	}{
		{"alone", Query{}},
		{"with an intro and outro", Query{Intro: "intro.md", Outro: "outro.md"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			empty := memlog.WithRecords(map[string]string{})
			result, err := Handler{Kind: adr.Kind, Log: empty}.Handle(context.Background(), test.query)
			if err != nil {
				t.Fatalf("Handle: %v", err)
			}
			if result.Markdown != "" {
				t.Errorf("Markdown =\n%s\nwant nothing", result.Markdown)
			}
		})
	}
}

func TestHandleSurroundsTheListWithAnIntroAndOutro(t *testing.T) {
	log := newLog()
	if err := log.Write(context.Background(), "intro.md", "An intro.\n\nMultiple paragraphs.\n"); err != nil {
		t.Fatal(err)
	}
	if err := log.Write(context.Background(), "outro.md", "An outro.\n"); err != nil {
		t.Fatal(err)
	}

	result, err := Handler{Kind: adr.Kind, Log: log}.Handle(context.Background(), Query{Intro: "intro.md", Outro: "outro.md"})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	want := "# Architecture Decision Records\n\n" +
		"An intro.\n\nMultiple paragraphs.\n\n" +
		"* [1. First Decision](0001-first.md)\n" +
		"* [2. Second Decision](0002-second.md)\n\n" +
		"An outro.\n"
	if result.Markdown != want {
		t.Errorf("Markdown =\n%s\nwant\n%s", result.Markdown, want)
	}
}
