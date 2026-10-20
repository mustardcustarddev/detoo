// Command fyr records, links and reports on fire drills.
//
// Run "fyr help" for the available commands.
//
// This package holds no behaviour of its own: it gives the fyr tool's
// definition to the program every record tool shares, and exits with the
// status that returns.
package main

import (
	"os"

	"github.com/mustardcustard/detoo/internal/fyr"
	"github.com/mustardcustard/detoo/internal/tool"
)

func main() {
	os.Exit(tool.Run(fyr.Kind, fyr.Files))
}
