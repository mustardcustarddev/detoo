package adr

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const helpHelp = `usage: adr help [COMMAND [SUBCOMMAND...]]

If COMMAND and any SUBCOMMANDs are given, displays help for that command.
Otherwise lists the available commands.

Uses the environment variables ADR_PAGER or PAGER (in that order) to
specify the command to be used to display the help text.  The default
pager is: more.
`

func runHelp(e *Env, args []string) error {
	cmd := find(commands, args)
	if cmd == nil || cmd.Hidden {
		printCommands(e)
		return nil
	}
	help := cmd.Help
	if len(cmd.Subs) > 0 {
		help += "\nSubcommands:\n"
		for _, sub := range visible(cmd.Subs) {
			help += fmt.Sprintf("  adr %s %s\n", strings.Join(args, " "), sub.Name)
		}
	}
	return e.page(help)
}

// find walks a path of command names, returning the command it names.
func find(in []*Command, names []string) *Command {
	var cmd *Command
	for _, name := range names {
		if cmd = lookup(in, name); cmd == nil {
			return nil
		}
		in = cmd.Subs
	}
	return cmd
}

func printCommands(e *Env) {
	e.print("usage: adr help COMMAND [ARG] ...")
	e.print("COMMAND is one of:")
	for _, cmd := range visible(commands) {
		e.printf("  %s\n", cmd.Name)
	}
	e.print("Run 'adr help COMMAND' for help on a specific command.")
}

// page writes text through the user's pager when the output is a terminal, and
// straight out when it is not.
func (e *Env) page(text string) error {
	pager := e.getenv("ADR_PAGER")
	if pager == "" {
		pager = e.getenv("PAGER")
	}
	if pager == "" {
		pager = "more"
	}
	words := strings.Fields(pager)
	if len(words) > 0 && isTerminal(e.Stdout) {
		cmd := exec.Command(words[0], words[1:]...)
		cmd.Stdin = strings.NewReader(text)
		cmd.Stdout = e.Stdout
		cmd.Stderr = e.Stderr
		if err := cmd.Run(); err == nil {
			return nil
		}
	}
	_, err := io.WriteString(e.Stdout, text)
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

// runAutocomplete suggests completions for a partially typed command line,
// which is passed as the words typed so far, starting with "adr".
func runAutocomplete(e *Env, words []string) error {
	word := func(i int) string {
		if i < len(words) {
			return words[i]
		}
		return ""
	}
	if len(words) <= 2 {
		suggestions := matching(names(visible(commands)), word(1))
		if !only(suggestions, word(1)) {
			e.print(strings.Join(suggestions, " "))
			return nil
		}
	}
	if cmd := lookup(commands, word(1)); cmd != nil && len(cmd.Subs) > 0 {
		suggestions := matching(names(visible(cmd.Subs)), word(2))
		if !only(suggestions, word(2)) {
			e.print(strings.Join(suggestions, " "))
		}
	}
	return nil
}

// only reports whether suggestions is just the word already typed, in which case
// there is nothing left to suggest.
func only(suggestions []string, word string) bool {
	return len(suggestions) == 1 && suggestions[0] == word
}
