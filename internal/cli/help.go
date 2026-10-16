package cli

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mustardcustard/decision-records/internal/domain/record"
)

// tool is what the commands know of the tool they belong to: the vocabulary of
// its log, and the files it ships with.
type tool struct {
	record.Kind
	files fs.FS
}

// help returns what a command prints about itself: the tool's ".md" file of
// prose for the command, or its ".example" file of worked examples.  The text
// is kept out of the Go source, as the record templates are.
//
// The blank lines a text file carries at its start and end are trimmed, but
// only newlines: an ".example" file indents every line, and that indentation is
// part of the text.  The files ship inside the binary, so a missing one is a
// mistake in the tool rather than a runtime condition;
// TestEveryCommandIsDocumented catches it before it reaches a user.
func (t tool) help(name string) string {
	text, err := fs.ReadFile(t.files, "help/"+name)
	if err != nil {
		panic(err)
	}
	return strings.Trim(string(text), "\n")
}

// usePagedHelp sends the help of root, and of every command under it, through
// the user's pager.  Cobra renders the help as it always does; only where it
// ends up changes.
func usePagedHelp(root *cobra.Command, t tool, stdout, stderr io.Writer, getenv func(string) string) {
	render := root.HelpFunc() // cobra's own, captured before it is replaced
	root.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		var help bytes.Buffer
		cmd.SetOut(&help)
		render(cmd, args)
		cmd.SetOut(stdout)

		// Cobra's help is written, not returned, so a failure to write it is
		// reported here rather than travelling back up to main.
		if err := page(help.String(), pager(t, getenv), stdout, stderr); err != nil {
			fmt.Fprintf(stderr, "%s: %v\n", t.Name, err)
		}
	})
}

// pager returns the command that shows help: the one named by the tool's PAGER
// variable, such as ADR_PAGER, else by PAGER, else more.
func pager(t tool, getenv func(string) string) string {
	if pager := getenv(t.Variable("PAGER")); pager != "" {
		return pager
	}
	if pager := getenv("PAGER"); pager != "" {
		return pager
	}
	return "more"
}

// page writes text through pager when the output is a terminal, and straight
// out when it is not.
func page(text, pager string, stdout, stderr io.Writer) error {
	words := strings.Fields(pager)
	if len(words) > 0 && isTerminal(stdout) {
		cmd := exec.Command(words[0], words[1:]...)
		cmd.Stdin = strings.NewReader(text)
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		if err := cmd.Run(); err == nil {
			return nil
		}
	}
	_, err := io.WriteString(stdout, text)
	return err
}

func isTerminal(w io.Writer) bool {
	file, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}
