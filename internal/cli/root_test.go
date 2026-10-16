package cli

import (
	"io"
	"testing"

	"github.com/spf13/cobra"

	"github.com/mustardcustard/detoo/internal/adr"
	"github.com/mustardcustard/detoo/internal/app"
	"github.com/mustardcustard/detoo/internal/detoo"
	"github.com/mustardcustard/detoo/internal/domain/record"
	"github.com/mustardcustard/detoo/internal/features/reports/graph"
	"github.com/mustardcustard/detoo/internal/rfc"
	"github.com/mustardcustard/detoo/internal/sty"
)

// These tests cover what belongs to the CLI: that the tree is assembled, that
// each command explains itself, and that its arguments parse. The behaviour
// behind the commands is tested in the slices, against fakes, and is not
// repeated here.
//
// A zero Application is enough, because no command is run: building the tree
// and parsing flags never reaches a handler.
func testRoots(t *testing.T) []*cobra.Command {
	t.Helper()
	getenv := func(string) string { return "" }
	return []*cobra.Command{
		NewRootCommand(adr.Kind, adr.Files, &app.Application{}, io.Discard, io.Discard, getenv),
		NewRootCommand(rfc.Kind, rfc.Files, &app.Application{}, io.Discard, io.Discard, getenv),
		NewRootCommand(sty.Kind, sty.Files, &app.Application{}, io.Discard, io.Discard, getenv),
		NewDetooCommand(detoo.Files, &app.Detoo{}, nil, io.Discard, io.Discard, getenv),
	}
}

// commandsUnder returns cmd and every command beneath it.
func commandsUnder(cmd *cobra.Command) []*cobra.Command {
	found := []*cobra.Command{cmd}
	for _, sub := range cmd.Commands() {
		found = append(found, commandsUnder(sub)...)
	}
	return found
}

// TestEveryCommandIsDocumented checks that each command of each tool explains
// itself.  The help is embedded, and help panics on a file that is not there,
// so building the trees at all proves the files exist; this proves that each
// command names one, and that none of them is empty.
func TestEveryCommandIsDocumented(t *testing.T) {
	for _, root := range testRoots(t) {
		for _, cmd := range commandsUnder(root) {
			if cmd.Short == "" {
				t.Errorf("%q has no short description", cmd.CommandPath())
			}
			if cmd.Long == "" {
				t.Errorf("%q has no long description", cmd.CommandPath())
			}
		}
	}
}

// TestEveryCommandTakesItsOptionsFirst guards the POSIX rule that options come
// before operands.  Without it an operand beginning with "-" is read as an
// option and rejected.
func TestEveryCommandTakesItsOptionsFirst(t *testing.T) {
	for _, root := range testRoots(t) {
		for _, cmd := range commandsUnder(root) {
			// pflag does not report the setting, so ask it to parse an operand
			// followed by something option-shaped: it belongs to the operands.
			flags := cmd.Flags()
			if err := flags.Parse([]string{"operand", "-s"}); err != nil {
				t.Errorf("%q: parsing an operand before an option: %v", cmd.CommandPath(), err)
				continue
			}
			if got := flags.Args(); len(got) != 2 {
				t.Errorf("%q: operands are %q, want both of them", cmd.CommandPath(), got)
			}
		}
	}
}

// TestPagerIsTheTools covers where help is paged: through the tool's own PAGER
// variable, else PAGER, else more.  A pager runs only on a terminal, so the
// end-to-end tests cannot see which one was chosen.
func TestPagerIsTheTools(t *testing.T) {
	env := map[string]string{}
	getenv := func(name string) string { return env[name] }
	rfcTool := tool{Kind: rfc.Kind}

	if got := pager(rfcTool, getenv); got != "more" {
		t.Errorf("with nothing set, pager = %q, want more", got)
	}
	env["PAGER"] = "less"
	env["ADR_PAGER"] = "adr-pager"
	if got := pager(rfcTool, getenv); got != "less" {
		t.Errorf("with PAGER and ADR_PAGER set, rfc's pager = %q, want less", got)
	}
	env["RFC_PAGER"] = "rfc-pager"
	if got := pager(rfcTool, getenv); got != "rfc-pager" {
		t.Errorf("with RFC_PAGER set, pager = %q, want rfc-pager", got)
	}
}

// TestTitleArgRequiresExactlyOneOperand covers the CLI half of the validation
// boundary: arity and quoting are checked here, while an empty title is a rule
// the create slice enforces for every client.
func TestTitleArgRequiresExactlyOneOperand(t *testing.T) {
	for _, test := range []struct {
		name string
		args []string
		want string
	}{
		{"none", nil, "no title given"},
		{"one", []string{"A Decision"}, ""},
		{"several", []string{"A", "Decision"}, "quote the title so it is a single argument"},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := titleArg(nil, test.args)
			switch {
			case test.want == "" && err != nil:
				t.Errorf("titleArg(%q) = %v, want no error", test.args, err)
			case test.want != "" && (err == nil || err.Error() != test.want):
				t.Errorf("titleArg(%q) = %v, want %q", test.args, err, test.want)
			}
		})
	}
}

// TestRelationFlagNeedsExactlyOneSwitch covers the CLI half of validating a
// link: which switch was typed is checked here, and named in the tool's words,
// while a request with no relation is a rule the link slice enforces for every
// client.
func TestRelationFlagNeedsExactlyOneSwitch(t *testing.T) {
	for _, test := range []struct {
		name            string
		kind            record.Kind
		replace, update bool
		want            record.Relation
		wantErr         string
	}{
		{"supersede", adr.Kind, true, false, record.Replaces, ""},
		{"update", adr.Kind, false, true, record.Updates, ""},
		{"neither", adr.Kind, false, false, 0, "give either --supersede or --update"},
		{"both", adr.Kind, true, true, 0, "give either --supersede or --update"},
		{"obsolete", rfc.Kind, true, false, record.Replaces, ""},
		{"neither, for an RFC", rfc.Kind, false, false, 0, "give either --obsolete or --update"},
		{"replace", sty.Kind, true, false, record.Replaces, ""},
		{"neither, for a user story", sty.Kind, false, false, 0, "give either --replace or --update"},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := relationFlag(test.kind, test.replace, test.update)
			switch {
			case test.wantErr == "" && err != nil:
				t.Errorf("relationFlag() = %v, want no error", err)
			case test.wantErr != "" && (err == nil || err.Error() != test.wantErr):
				t.Errorf("relationFlag() = %v, want %q", err, test.wantErr)
			case got != test.want:
				t.Errorf("relationFlag() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestFormatFlagDefaultsToMermaid covers the switches "graph" takes: the graph
// is Mermaid unless --graphviz is given, and cannot be both.
func TestFormatFlagDefaultsToMermaid(t *testing.T) {
	for _, test := range []struct {
		name              string
		graphviz, mermaid bool
		want              graph.Format
		wantErr           string
	}{
		{"neither", false, false, graph.Mermaid, ""},
		{"mermaid", false, true, graph.Mermaid, ""},
		{"graphviz", true, false, graph.Graphviz, ""},
		{"both", true, true, graph.Mermaid, "give either --graphviz or --mermaid, not both"},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := formatFlag(test.graphviz, test.mermaid)
			switch {
			case test.wantErr == "" && err != nil:
				t.Errorf("formatFlag() = %v, want no error", err)
			case test.wantErr != "" && (err == nil || err.Error() != test.wantErr):
				t.Errorf("formatFlag() = %v, want %q", err, test.wantErr)
			case got != test.want:
				t.Errorf("formatFlag() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestTOCOptionsNeedTheTOCSwitch covers "list": the options that shape a table
// of contents are rejected without --toc, rather than silently ignored.
func TestTOCOptionsNeedTheTOCSwitch(t *testing.T) {
	for _, test := range []struct {
		name                 string
		toc                  bool
		intro, outro, prefix string
		wantErr              bool
	}{
		{"list", false, "", "", "", false},
		{"toc", true, "intro.md", "outro.md", "docs/", false},
		{"intro without toc", false, "intro.md", "", "", true},
		{"outro without toc", false, "", "outro.md", "", true},
		{"prefix without toc", false, "", "", "docs/", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := tocOptions(test.toc, test.intro, test.outro, test.prefix)
			if (err != nil) != test.wantErr {
				t.Errorf("tocOptions() = %v, want an error: %t", err, test.wantErr)
			}
		})
	}
}
