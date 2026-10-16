// Package tool puts a record tool together and runs it as a command.
//
// Every record tool is the same program given a different definition, so the
// program is written once, here, and each package main is a single call: see
// cmd/adr, cmd/rfc and cmd/sty.  Nothing below this package calls os.Exit, and
// nothing below it chooses an adapter.
package tool

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/mustardcustard/detoo/internal/adapters/clock"
	"github.com/mustardcustard/detoo/internal/adapters/editor"
	"github.com/mustardcustard/detoo/internal/adapters/logstore"
	"github.com/mustardcustard/detoo/internal/adapters/templates"
	"github.com/mustardcustard/detoo/internal/app"
	"github.com/mustardcustard/detoo/internal/cli"
	"github.com/mustardcustard/detoo/internal/domain/record"
)

// Run runs the tool that keeps logs of kind, whose templates and help are in
// files, with the process's arguments, and returns its exit status.
func Run(kind record.Kind, files fs.FS) int {
	workdir, err := os.Getwd()
	if err != nil {
		return fail(kind, err)
	}

	// The composition root.  Every dependence on the outside world — the disk,
	// the environment, the clock, the user's editor — is chosen here, once.
	application := app.New(
		kind,
		logstore.New(workdir, kind),
		templates.New(workdir, os.Getenv, kind, files),
		clock.New(os.Getenv, kind),
		editor.New(workdir, os.Getenv, os.Stdin, os.Stdout, os.Stderr),
	)

	if err := cli.NewRootCommand(kind, files, application, os.Stdout, os.Stderr, os.Getenv).Execute(); err != nil {
		// A command that has already explained itself says only how to exit.
		if exit, ok := errors.AsType[*cli.ExitError](err); ok {
			return exit.Code
		}
		return fail(kind, err)
	}
	return 0
}

// fail reports an error the way the Go tools do, on stderr and named after the
// tool.
func fail(kind record.Kind, err error) int {
	fmt.Fprintf(os.Stderr, "%s: %v\n", kind.Name, err)
	return 1
}
