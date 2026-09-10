// Package cli is the Cobra adapter, and the only package in this module that
// imports Cobra.
//
// Every file here translates: it defines a command's arguments and flags,
// checks that the invocation is well shaped, builds a use-case input, calls one
// handler, and renders what comes back.  No business rule lives in this
// package, because a rule written here would be reachable from the command line
// and from nowhere else.
package cli

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/mustardcustard/decision-records/internal/app"
)

// NewRootCommand returns the adr command tree, wired to run against
// application.  The caller executes it: see cmd/adr.
func NewRootCommand(application *app.Application, stdout, stderr io.Writer, getenv func(string) string) *cobra.Command {
	root := &cobra.Command{
		Use:   "adr",
		Short: "Record, link and report on architecture decision records",
		Long:  helpText("root.md"),

		// Errors are reported by the caller, prefixed with "adr:", rather than
		// by cobra, which would also print the usage after them.
		SilenceUsage:  true,
		SilenceErrors: true,

		// Run without a command, adr explains itself and fails, so that a
		// script that calls it wrongly does not appear to succeed.
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.Help()
			return exitWith(1)
		},
	}
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.AddCommand(
		newConfigCommand(application),
		newGenerateCommand(application),
		newInitCommand(application),
		newLinkCommand(application),
		newListCommand(application),
		newNewCommand(application),
	)
	optionsFirst(root)
	usePagedHelp(root, stdout, stderr, getenv)
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
