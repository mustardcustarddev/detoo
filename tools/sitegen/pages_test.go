package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestEveryCommandHasAPage guards against a command being added to a tool and
// left off the site: the pages are found by walking each tool's commands, so
// nothing needs editing when one is added, and this test proves it.
func TestEveryCommandHasAPage(t *testing.T) {
	pages, err := Pages()
	if err != nil {
		t.Fatal(err)
	}
	have := map[string]bool{}
	for _, p := range pages {
		have[p.Tool+" "+p.Command] = true
	}

	for _, tl := range tools() {
		if !have[tl.name+" "] {
			t.Errorf("no overview page for %s", tl.name)
		}
		for _, cmd := range tl.root(&strings.Builder{}).Commands() {
			if documented(cmd) && !have[tl.name+" "+cmd.Name()] {
				t.Errorf("no page for %s %s", tl.name, cmd.Name())
			}
		}
	}

	// The walk above proves nothing if there are no tools to walk.  Adding a
	// tool means editing tools(), so naming them here costs nothing extra;
	// adding a command needs no edit at all.
	for _, name := range []string{"adr", "rfc", "sty", "plm", "sop", "detoo"} {
		if !have[name+" "] {
			t.Errorf("no overview page for %s", name)
		}
	}
}

func TestPageShowsBuiltInHelp(t *testing.T) {
	tests := []struct {
		tool, command, prose string
		want                 []string
	}{
		{"adr", "new", "../../internal/adr/help/new.md", []string{"Usage:\n  adr new TITLE [flags]", "--supersede REFERENCE"}},
		{"detoo", "doctor", "../../internal/detoo/help/doctor.md", []string{"Usage:\n  detoo doctor"}},
	}
	for _, tt := range tests {
		p := page(t, tt.tool, tt.command)
		prose, err := os.ReadFile(tt.prose)
		if err != nil {
			t.Fatal(err)
		}
		first, _, _ := strings.Cut(strings.TrimSpace(string(prose)), "\n")
		for _, want := range append(tt.want, first) {
			if !strings.Contains(p.Help, want) {
				t.Errorf("%s %s: help does not contain %q:\n%s", tt.tool, tt.command, want, p.Help)
			}
		}
	}
}

// TestHelpMatchesBinary is the story's acceptance criterion: the page shows
// the same text as the tool's built-in help.
func TestHelpMatchesBinary(t *testing.T) {
	cmd := exec.Command("go", "run", "./cmd/adr", "help", "new")
	cmd.Dir = "../.."
	cmd.Env = append(os.Environ(), "ADR_PAGER=")
	want, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	if got := page(t, "adr", "new").Help; got != string(want) {
		t.Errorf("page differs from `adr help new`\n--- page\n%s\n--- adr help new\n%s", got, want)
	}
}

func TestFenceOutlastsBackticks(t *testing.T) {
	md := Page{Tool: "adr", Title: "adr", Help: "a ```` b"}.Markdown()
	if !strings.Contains(md, "\n`````text\na ```` b\n`````\n") {
		t.Errorf("help is not fenced by a longer run of backticks:\n%s", md)
	}
}

func TestFrontmatterQuotesValues(t *testing.T) {
	md := Page{Tool: "adr", Title: "x", Description: `a: "b"`}.Markdown()
	if !strings.Contains(md, "\ndescription: \"a: \\\"b\\\"\"\n") {
		t.Errorf("description is not quoted:\n%s", md)
	}
}

func TestPagePaths(t *testing.T) {
	if got := (Page{Tool: "adr"}).Path(); got != "adr/index.md" {
		t.Errorf("overview path = %q", got)
	}
	if got := (Page{Tool: "adr", Command: "new"}).Path(); got != "adr/new.md" {
		t.Errorf("command path = %q", got)
	}
}

func page(t *testing.T, tool, command string) Page {
	t.Helper()
	pages, err := Pages()
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range pages {
		if p.Tool == tool && p.Command == command {
			return p
		}
	}
	t.Fatalf("no page for %s %s", tool, command)
	return Page{}
}
