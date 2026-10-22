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
		for _, cmd := range commandsUnder(tl.root(&strings.Builder{})) {
			if command := strings.TrimPrefix(cmd.CommandPath(), tl.name+" "); !have[tl.name+" "+command] {
				t.Errorf("no page for %s %s", tl.name, command)
			}
		}
	}

	// The walk above proves nothing if there are no tools to walk.  Adding a
	// tool means editing tools(), so naming them here costs nothing extra;
	// adding a command needs no edit at all.
	for _, name := range []string{"adr", "rfc", "ust", "plm", "sop", "fyr", "detoo"} {
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
		{"adr", "new", "../../internal/adr/help/new.md", []string{
			"## Usage\n\n```bash\nadr new TITLE [flags]\n```",
			"## Examples\n\n```bash\n# Create a new ADR with the title \"Use MySQL Database\":\nadr new \"Use MySQL Database\"\n",
			"| `--supersede REFERENCE`, `-s` | a REFERENCE (number or partial filename) to a record this one supersedes |",
		}},
		{"adr", "", "../../internal/adr/help/root.md", []string{
			"```bash\nadr [command]\n```",
			"| [`new`](new/) | Create the next numbered decision record |",
		}},
		{"detoo", "doctor", "../../internal/detoo/help/doctor.md", []string{"```bash\ndetoo doctor [flags]\n```"}},
		{"detoo", "skills install", "../../internal/detoo/help/skills-install.md", []string{
			"```bash\ndetoo skills install [SKILL...] [flags]\n```",
			"| `--source SOURCE` |",
			"(default `mustardcustarddev/detoo`) |",
		}},
	}
	for _, tt := range tests {
		p := page(t, tt.tool, tt.command)
		prose, err := os.ReadFile(tt.prose)
		if err != nil {
			t.Fatal(err)
		}
		// The prose is Markdown already, so it is on the page as it is, but
		// for an overview's first paragraph, which is the page's description.
		text := strings.TrimSpace(string(prose))
		if tt.command == "" {
			_, text, _ = strings.Cut(text, "\n\n")
		}
		for _, want := range append(tt.want, text) {
			if !strings.Contains(p.Body, want) {
				t.Errorf("%s %s: page does not contain %q:\n%s", tt.tool, tt.command, want, p.Body)
			}
		}
		for _, unwanted := range []string{"--help", "Use \""} {
			if strings.Contains(p.Body, unwanted) {
				t.Errorf("%s %s: page contains %q:\n%s", tt.tool, tt.command, unwanted, p.Body)
			}
		}
	}
}

// TestPageSaysWhatHelpSays is the story's acceptance criterion: the page says
// what the tool's built-in help says.  The page is laid out as a page rather
// than a terminal, so each part of the help is looked for on it: the prose,
// each line of the examples, and what each option does.
func TestPageSaysWhatHelpSays(t *testing.T) {
	cmd := exec.Command("go", "run", "./cmd/adr", "help", "new")
	cmd.Dir = "../.."
	cmd.Env = append(os.Environ(), "ADR_PAGER=")
	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	help := string(out)
	body := page(t, "adr", "new").Body

	prose, rest, _ := strings.Cut(help, "\n\nUsage:\n")
	want := []string{prose}
	section := ""
	for _, line := range strings.Split(rest, "\n") {
		switch {
		case line == "Examples:" || line == "Flags:":
			section = line
		case strings.TrimSpace(line) == "" || strings.Contains(line, "--help"):
		case section == "Examples:":
			want = append(want, strings.TrimSpace(line))
		case section == "Flags:":
			// "  -s, --supersede REFERENCE   what it does": what it does follows
			// the first run of two spaces after the option.
			_, usage, _ := strings.Cut(strings.TrimSpace(line), "  ")
			want = append(want, strings.TrimSpace(usage))
		}
	}
	if len(want) < 6 {
		t.Fatalf("found too little in `adr help new` to compare:\n%s", help)
	}
	for _, w := range want {
		if !strings.Contains(body, w) {
			t.Errorf("page does not contain %q from `adr help new`:\n%s", w, body)
		}
	}
}

func TestFenceOutlastsBackticks(t *testing.T) {
	if got := codeBlock("a ```` b"); got != "`````bash\na ```` b\n`````" {
		t.Errorf("code is not fenced by a longer run of backticks:\n%s", got)
	}
}

func TestDedentKeepsRelativeIndentation(t *testing.T) {
	if got := dedent("\n  # a\n  b\n\n    c\n"); got != "# a\nb\n\n  c" {
		t.Errorf("dedent = %q", got)
	}
}

func TestOverviewDoesNotRepeatItsDescription(t *testing.T) {
	p := page(t, "adr", "")
	if strings.Contains(p.Body, p.Description) {
		t.Errorf("overview repeats its description %q:\n%s", p.Description, p.Body)
	}
	if got := withoutLead("A b.\n\nC.\n", "A b."); got != "C.\n" {
		t.Errorf("withoutLead = %q", got)
	}
	if got := withoutLead("A b. C.\n\nD.\n", "A b."); got != "A b. C.\n\nD.\n" {
		t.Errorf("withoutLead removed a paragraph that says more than the lead: %q", got)
	}
}

func TestCellEscapesPipes(t *testing.T) {
	if got := cell("a | b\n  c"); got != `a \| b c` {
		t.Errorf("cell = %q", got)
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
