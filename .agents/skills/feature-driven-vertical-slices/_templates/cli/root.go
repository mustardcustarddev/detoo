// Package cli is the Cobra adapter, and the only package in the module that
// imports Cobra.  It assembles the command tree and owns the concerns common to
// every command — global flags, output format, error and exit behaviour.
//
// Keep it a chassis, not a garage.  Files here are named after the command they
// translate (users_create.go), and they contain translation only.  The moment a
// business rule lands in this package it is unreachable from every other client.
package cli

import (
	"github.com/spf13/cobra"

	"github.com/example/myapp/internal/app"
	"github.com/example/myapp/internal/cli/output"
)

// NewRootCommand returns the command tree, wired to run against application.
// The caller executes it: see cmd/myapp.
func NewRootCommand(application *app.Application) *cobra.Command {
	format := &output.Format{Kind: output.Human}

	root := &cobra.Command{
		Use:   "myapp",
		Short: "Manage users and projects",

		// Errors are reported once, by main, prefixed with the program name.
		// Left to Cobra they would be printed here and followed by the usage.
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Persistent flags are reserved for genuinely global concerns.  A flag that
	// means something to one command belongs on that command.
	root.PersistentFlags().VarP(format, "output", "o", "output `FORMAT`: human or json")

	root.AddCommand(newUsersCommand(application, format))
	return root
}

// newUsersCommand groups the "users" slices under one parent command.  It has no
// RunE of its own: it exists to give the subcommands a namespace.
func newUsersCommand(application *app.Application, format *output.Format) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Create, inspect and remove users",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(
		newUsersCreateCommand(application, format),
		newUsersGetCommand(application, format),
		newUsersDeleteCommand(application),
	)
	return cmd
}
