// Command adr records, links and reports on architecture decision records.
//
// Run "adr help" for the available commands.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/mustardcustard/decision-records/internal/adr"
)

func main() {
	os.Exit(run())
}

func run() int {
	env, err := adr.NewEnv()
	if err == nil {
		err = adr.Run(env, os.Args[1:])
	}
	if err == nil {
		return 0
	}
	var exit *adr.ExitError
	if errors.As(err, &exit) {
		return exit.Code
	}
	fmt.Fprintf(os.Stderr, "adr: %v\n", err)
	return 1
}
