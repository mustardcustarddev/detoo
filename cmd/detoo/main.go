// Command detoo looks after a project's logs as a whole.
//
// Run "detoo help" for the available commands.
//
// This package holds no behaviour of its own: it gives detoo's definition to
// the program that runs it, and exits with the status that returns.
package main

import (
	"os"

	"github.com/mustardcustard/detoo/internal/detoo"
	"github.com/mustardcustard/detoo/internal/tool"
)

func main() {
	os.Exit(tool.RunDetoo(detoo.Kinds, detoo.Files))
}
