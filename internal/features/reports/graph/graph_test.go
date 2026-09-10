package graph

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/mustardcustard/decision-records/internal/adapters/memlog"
)

func newLog() *memlog.Log {
	return memlog.WithRecords(map[string]string{
		"docs/adr/0001-first.md": "# 1. First\n\n## Status\n\nSuperceded by [2. Second](0002-second.md)\n\n## Context\n\nC.\n",
		"docs/adr/0002-second.md": "# 2. Second\n\n## Status\n\nAccepted\n\n" +
			"Supercedes [1. First](0001-first.md)\n\n## Context\n\nC.\n",
	})
}

// A graph is drawn in Mermaid unless another format is asked for.
func TestHandleDrawsTheLogInMermaidByDefault(t *testing.T) {
	result, err := Handler{Log: newLog()}.Handle(context.Background(), Query{LinkExtension: ".html"})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	want := "flowchart TD\n" +
		"  adr1[\"1. First\"]\n" +
		"  adr2[\"2. Second\"]\n" +
		"  adr1 -.-> adr2\n" +
		"  adr2 -->|\"Supercedes\"| adr1\n" +
		"  click adr1 \"0001-first.html\"\n" +
		"  click adr2 \"0002-second.html\"\n"
	if result.Document != want {
		t.Errorf("Document =\n%s\nwant\n%s", result.Document, want)
	}
}

func TestHandleWritesADigraphOfTheLog(t *testing.T) {
	result, err := Handler{Log: newLog()}.Handle(context.Background(), Query{Format: Graphviz, LinkExtension: ".html"})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	want := "digraph {\n" +
		"  node [shape=plaintext];\n" +
		"  subgraph {\n" +
		"    _1 [label=\"1. First\"; URL=\"0001-first.html\"];\n" +
		"    _2 [label=\"2. Second\"; URL=\"0002-second.html\"];\n" +
		"    _1 -> _2 [style=\"dotted\", weight=1];\n" +
		"  }\n" +
		"  _2 -> _1 [label=\"Supercedes\", weight=0]\n" +
		"}\n"
	if result.Document != want {
		t.Errorf("Document =\n%s\nwant\n%s", result.Document, want)
	}
}

// A reverse link is the other record's link seen from the far end, so drawing
// it too would give every relationship two arrows.
func TestHandleDrawsEachLinkOnce(t *testing.T) {
	for _, test := range []struct {
		name   string
		format Format
		arrow  string // what each link, and nothing else, is drawn with
	}{
		{"Mermaid", Mermaid, "-->|"},
		{"Graphviz", Graphviz, "weight=0"},
	} {
		t.Run(test.name, func(t *testing.T) {
			result, err := Handler{Log: newLog()}.Handle(context.Background(), Query{Format: test.format})
			if err != nil {
				t.Fatalf("Handle: %v", err)
			}
			if strings.Contains(result.Document, "Superceded by") {
				t.Errorf("the graph draws a reverse link:\n%s", result.Document)
			}
			if got := strings.Count(result.Document, test.arrow); got != 1 {
				t.Errorf("the graph has %d links, want 1:\n%s", got, result.Document)
			}
		})
	}
}

func TestHandlePrefixesAndExtendsTheNodeLinks(t *testing.T) {
	for _, test := range []struct {
		name   string
		format Format
		want   string
	}{
		{"Mermaid", Mermaid, `click adr1 "http://example.com/0001-first.xxx"`},
		{"Graphviz", Graphviz, `URL="http://example.com/0001-first.xxx"`},
	} {
		t.Run(test.name, func(t *testing.T) {
			result, err := Handler{Log: newLog()}.Handle(context.Background(), Query{
				Format:        test.format,
				LinkPrefix:    "http://example.com/",
				LinkExtension: ".xxx",
			})
			if err != nil {
				t.Fatalf("Handle: %v", err)
			}
			if !strings.Contains(result.Document, test.want) {
				t.Errorf("the graph does not contain %s:\n%s", test.want, result.Document)
			}
		})
	}
}

// Mermaid has no backslash escapes, so a quote in a title would end the node's
// label early and break the diagram.
func TestHandleEscapesQuotesForMermaid(t *testing.T) {
	log := memlog.WithRecords(map[string]string{
		"docs/adr/0001-use-quotes.md": "# 1. Use \"quotes\"\n\n## Status\n\nAccepted\n\n## Context\n\nC.\n",
	})

	result, err := Handler{Log: log}.Handle(context.Background(), Query{})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	if want := `adr1["1. Use #quot;quotes#quot;"]`; !strings.Contains(result.Document, want) {
		t.Errorf("the graph does not contain %s:\n%s", want, result.Document)
	}
}

// A log with no records has nothing to draw, so nothing is written in either
// format — not even an empty flowchart or digraph.
func TestHandleDrawsNothingForAnEmptyLog(t *testing.T) {
	for _, test := range []struct {
		name   string
		format Format
	}{
		{"Mermaid", Mermaid},
		{"Graphviz", Graphviz},
	} {
		t.Run(test.name, func(t *testing.T) {
			empty := memlog.WithRecords(map[string]string{})
			result, err := Handler{Log: empty}.Handle(context.Background(), Query{Format: test.format})
			if err != nil {
				t.Fatalf("Handle: %v", err)
			}
			if result.Document != "" {
				t.Errorf("Document =\n%s\nwant nothing", result.Document)
			}
		})
	}
}

// The format is checked before the log is read, so an unknown one is reported
// whether or not there is anything to draw.
func TestHandleRejectsAnUnknownFormat(t *testing.T) {
	for _, test := range []struct {
		name string
		log  *memlog.Log
	}{
		{"records", newLog()},
		{"empty", memlog.WithRecords(map[string]string{})},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := Handler{Log: test.log}.Handle(context.Background(), Query{Format: Format(99)})
			if !errors.Is(err, ErrUnknownFormat) {
				t.Errorf("err = %v, want %v", err, ErrUnknownFormat)
			}
		})
	}
}
