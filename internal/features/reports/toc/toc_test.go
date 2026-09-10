package toc

import (
	"context"
	"testing"

	"github.com/mustardcustard/decision-records/internal/adapters/memlog"
)

func newLog() *memlog.Log {
	return memlog.WithRecords(map[string]string{
		"docs/adr/0001-first.md":  "# 1. First Decision\n",
		"docs/adr/0002-second.md": "# 2. Second Decision\n",
	})
}

func TestHandleWritesAListOfRecords(t *testing.T) {
	result, err := Handler{Log: newLog()}.Handle(context.Background(), Query{})
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

// The prefix lets the contents be published somewhere other than beside the
// records it links to.
func TestHandlePrefixesTheLinks(t *testing.T) {
	result, err := Handler{Log: newLog()}.Handle(context.Background(), Query{LinkPrefix: "foo/docs/adr/"})
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

func TestHandleSurroundsTheListWithAnIntroAndOutro(t *testing.T) {
	log := newLog()
	if err := log.Write(context.Background(), "intro.md", "An intro.\n\nMultiple paragraphs.\n"); err != nil {
		t.Fatal(err)
	}
	if err := log.Write(context.Background(), "outro.md", "An outro.\n"); err != nil {
		t.Fatal(err)
	}

	result, err := Handler{Log: log}.Handle(context.Background(), Query{Intro: "intro.md", Outro: "outro.md"})
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
