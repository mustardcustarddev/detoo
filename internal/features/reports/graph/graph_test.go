package graph

import (
	"context"
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

func TestHandleWritesADigraphOfTheLog(t *testing.T) {
	result, err := Handler{Log: newLog()}.Handle(context.Background(), Query{LinkExtension: ".html"})
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
	if result.Graphviz != want {
		t.Errorf("Graphviz =\n%s\nwant\n%s", result.Graphviz, want)
	}
}

// A reverse link is the other record's link seen from the far end, so drawing
// it too would give every relationship two arrows.
func TestHandleDrawsEachLinkOnce(t *testing.T) {
	result, err := Handler{Log: newLog()}.Handle(context.Background(), Query{})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	if strings.Contains(result.Graphviz, "Superceded by") {
		t.Errorf("the graph draws a reverse link:\n%s", result.Graphviz)
	}
	if got := strings.Count(result.Graphviz, "weight=0"); got != 1 {
		t.Errorf("the graph has %d links, want 1:\n%s", got, result.Graphviz)
	}
}

func TestHandlePrefixesAndExtendsTheNodeLinks(t *testing.T) {
	result, err := Handler{Log: newLog()}.Handle(context.Background(), Query{
		LinkPrefix:    "http://example.com/",
		LinkExtension: ".xxx",
	})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	if want := `URL="http://example.com/0001-first.xxx"`; !strings.Contains(result.Graphviz, want) {
		t.Errorf("the graph does not contain %s:\n%s", want, result.Graphviz)
	}
}
