package cli

import (
	"io"
	"testing"

	"github.com/spf13/cobra"

	"github.com/mustardcustard/decision-records/internal/app"
)

// These tests cover what belongs to the CLI: that the tree is assembled, that
// each command explains itself, and that its arguments parse. The behaviour
// behind the commands is tested in the slices, against fakes, and is not
// repeated here.
//
// A zero Application is enough, because no command is run: building the tree
// and parsing flags never reaches a handler.
func testRoot(t *testing.T) *cobra.Command {
	t.Helper()
	return NewRootCommand(&app.Application{}, io.Discard, io.Discard, func(string) string { return "" })
}

// commandsUnder returns cmd and every command beneath it.
func commandsUnder(cmd *cobra.Command) []*cobra.Command {
	found := []*cobra.Command{cmd}
	for _, sub := range cmd.Commands() {
		found = append(found, commandsUnder(sub)...)
	}
	return found
}

// TestEveryCommandIsDocumented checks that each command explains itself.  The
// help is embedded, and helpText panics on a file that is not there, so
// building the tree at all proves the files exist; this proves that each
// command names one, and that none of them is empty.
func TestEveryCommandIsDocumented(t *testing.T) {
	for _, cmd := range commandsUnder(testRoot(t)) {
		if cmd.Short == "" {
			t.Errorf("%q has no short description", cmd.CommandPath())
		}
		if cmd.Long == "" {
			t.Errorf("%q has no long description", cmd.CommandPath())
		}
	}
}

// TestEveryCommandTakesItsOptionsFirst guards the POSIX rule that options come
// before operands.  Without it an operand beginning with "-" — a record title,
// or a link description — is read as an option and rejected.
func TestEveryCommandTakesItsOptionsFirst(t *testing.T) {
	for _, cmd := range commandsUnder(testRoot(t)) {
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

func TestParseLinksSplitsTheFlagEncoding(t *testing.T) {
	links, err := parseLinks([]string{"1:Amends:Amended by"})
	if err != nil {
		t.Fatalf("parseLinks: %v", err)
	}
	if len(links) != 1 {
		t.Fatalf("parsed %d links, want 1", len(links))
	}
	if links[0].Target != "1" || links[0].Forward != "Amends" || links[0].Reverse != "Amended by" {
		t.Errorf("parseLinks gave %+v", links[0])
	}

	if _, err := parseLinks([]string{"1:Amends"}); err == nil {
		t.Error("parseLinks accepted a two-part link, want an error")
	}
}
