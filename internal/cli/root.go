// Package cli is the Cobra adapter, and the only package in this module that
// imports Cobra.
//
// Every file here translates: it defines a command's arguments and flags,
// checks that the invocation is well shaped, builds a use-case input, calls one
// handler, and renders what comes back.  No business rule lives in this
// package, because a rule written here would be reachable from the command line
// and from nowhere else.
//
// The same command tree serves every record tool.  What a command is called,
// and the words its flags and help are written in, come from the tool.
package cli

import (
	"io"
	"io/fs"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mustardcustard/detoo/internal/app"
	"github.com/mustardcustard/detoo/internal/domain/record"
)

// NewRootCommand returns the command tree of the tool that keeps logs of kind,
// whose help is in files, wired to run against application.  The caller
// executes it: see package tool.
func NewRootCommand(kind record.Kind, files fs.FS, application *app.Application, stdout, stderr io.Writer, getenv func(string) string) *cobra.Command {
	t := tool{Kind: kind, files: files}
	root := &cobra.Command{
		Use:   t.Name,
		Short: "Record, link and report on " + strings.ToLower(t.Heading),
		Long:  t.help("root.md"),

		// Errors are reported by the caller, prefixed with the tool's name,
		// rather than by cobra, which would also print the usage after them.
		SilenceUsage:  true,
		SilenceErrors: true,

		// Run without a command, the tool explains itself and fails, so that a
		// script that calls it wrongly does not appear to succeed.
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.Help()
			return exitWith(1)
		},
	}
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.AddCommand(
		newConfigCommand(t, application),
		newGraphCommand(t, application),
		newInitCommand(t, application),
		newLinkCommand(t, application),
		newListCommand(t, application),
		newNewCommand(t, application),
	)
	optionsFirst(root)
	usePagedHelp(root, t, stdout, stderr, getenv)
	return root
}

// optionsFirst applies to every command the POSIX rule that options come before
// operands, so that everything after the first operand is an operand.  Without
// it pflag reads the options from among the operands.
func optionsFirst(cmd *cobra.Command) {
	cmd.Flags().SetInterspersed(false)
	for _, sub := range cmd.Commands() {
		optionsFirst(sub)
	}
}
