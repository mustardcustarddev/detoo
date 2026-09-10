// Command adr records, links and reports on architecture decision records.
//
// Run "adr help" for the available commands.
//
// This is the only package main, and it holds no behaviour of its own: it
// chooses the adapters, wires the application, runs the command tree, and turns
// an error into an exit status.  That is why nothing below it calls os.Exit.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/mustardcustard/decision-records/internal/adapters/clock"
	"github.com/mustardcustard/decision-records/internal/adapters/editor"
	"github.com/mustardcustard/decision-records/internal/adapters/logstore"
	"github.com/mustardcustard/decision-records/internal/adapters/templates"
	"github.com/mustardcustard/decision-records/internal/app"
	"github.com/mustardcustard/decision-records/internal/cli"
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
		logstore.New(workdir),
		templates.New(workdir, os.Getenv),
		clock.New(os.Getenv),
		editor.New(workdir, os.Getenv, os.Stdin, os.Stdout, os.Stderr),
	)

	if err := cli.NewRootCommand(application, os.Stdout, os.Stderr, os.Getenv).Execute(); err != nil {
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
