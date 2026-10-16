// Command pse records, links and reports on decisions made with PIOSEE.
//
// Run "pse help" for the available commands.
//
// This package holds no behaviour of its own: it gives the pse tool's
// definition to the program every record tool shares, and exits with the
// status that returns.
package main

import (
	"os"

	"github.com/mustardcustard/detoo/internal/pse"
	"github.com/mustardcustard/detoo/internal/tool"
)

func main() {
	os.Exit(tool.Run(pse.Kind, pse.Files))
}
