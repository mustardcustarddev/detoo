// Command adr records, links and reports on architecture decision records.
//
// Run "adr help" for the available commands.
//
// This package holds no behaviour of its own: it chooses the adapters, wires
// the application for the adr tool, runs the command tree, and turns an error
// into an exit status.  That is why nothing below it calls os.Exit.  Every
// other record tool has a package main just like it.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/mustardcustard/detoo/internal/adapters/clock"
	"github.com/mustardcustard/detoo/internal/adapters/editor"
	"github.com/mustardcustard/detoo/internal/adapters/logstore"
	"github.com/mustardcustard/detoo/internal/adapters/templates"
	"github.com/mustardcustard/detoo/internal/adr"
	"github.com/mustardcustard/detoo/internal/app"
	"github.com/mustardcustard/detoo/internal/cli"
)

func main() {
	os.Exit(run())
}

func run() int {
	workdir, err := os.Getwd()
	if err != nil {
		return fail(err)
	}

	// The composition root.  Every dependence on the outside world — the disk,
	// the environment, the clock, the user's editor — is chosen here, once.
	application := app.New(
		adr.Kind,
		logstore.New(workdir, adr.Kind),
		templates.New(workdir, os.Getenv, adr.Kind, adr.Files),
		clock.New(os.Getenv, adr.Kind),
		editor.New(workdir, os.Getenv, os.Stdin, os.Stdout, os.Stderr),
	)

	if err := cli.NewRootCommand(adr.Kind, adr.Files, application, os.Stdout, os.Stderr, os.Getenv).Execute(); err != nil {
		// A command that has already explained itself says only how to exit.
		if exit, ok := errors.AsType[*cli.ExitError](err); ok {
			return exit.Code
		}
		return fail(err)
	}
	return 0
}

// fail reports an error the way the Go tools do, on stderr and named.
func fail(err error) int {
	fmt.Fprintf(os.Stderr, "adr: %v\n", err)
	return 1
}
