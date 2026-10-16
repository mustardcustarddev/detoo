// Command adr records, links and reports on architecture decision records.
//
// Run "adr help" for the available commands.
//
// This package holds no behaviour of its own: it gives the adr tool's
// definition to the program every record tool shares, and exits with the
// status that returns.
package main

import (
	"os"

	"github.com/mustardcustard/detoo/internal/adr"
	"github.com/mustardcustard/detoo/internal/tool"
)

func main() {
	os.Exit(tool.Run(adr.Kind, adr.Files))
}
