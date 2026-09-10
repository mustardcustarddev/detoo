package cli

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// helpFiles holds what the commands print about themselves: a ".md" file of
// prose for each, and a ".example" file for those that show worked examples.
// The text is kept out of the Go source, as the record templates are.
//
//go:embed help
var helpFiles embed.FS

// helpText returns an embedded help file, without the blank lines a text file
// carries at its start and end.  Only newlines are trimmed: an ".example" file
// indents every line, and that indentation is part of the text.  The files ship
// inside the binary, so a missing one is a mistake in this package rather than
// a runtime condition; TestEveryCommandIsDocumented catches it before it
// reaches a user.
func helpText(name string) string {
	text, err := helpFiles.ReadFile("help/" + name)
	if err != nil {
		panic(err)
	}
	return strings.Trim(string(text), "\n")
}

// usePagedHelp sends the help of root, and of every command under it, through
// the user's pager.  Cobra renders the help as it always does; only where it
// ends up changes.
func usePagedHelp(root *cobra.Command, stdout, stderr io.Writer, getenv func(string) string) {
	render := root.HelpFunc() // cobra's own, captured before it is replaced
	root.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		var help bytes.Buffer
		cmd.SetOut(&help)
		render(cmd, args)
		cmd.SetOut(stdout)

		// Cobra's help is written, not returned, so a failure to write it is
		// reported here rather than travelling back up to main.
		if err := page(help.String(), stdout, stderr, getenv); err != nil {
			fmt.Fprintf(stderr, "adr: %v\n", err)
		}
	})
}

// page writes text through the user's pager when the output is a terminal, and
// straight out when it is not.
func page(text string, stdout, stderr io.Writer, getenv func(string) string) error {
	pager := getenv("ADR_PAGER")
	if pager == "" {
		pager = getenv("PAGER")
	}
	if pager == "" {
		pager = "more"
	}
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
