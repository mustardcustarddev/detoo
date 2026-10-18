// Command sop records, links and reports on standard operating procedures.
//
// Run "sop help" for the available commands.
//
// This package holds no behaviour of its own: it gives the sop tool's
// definition to the program every record tool shares, and exits with the
// status that returns.
package main

import (
	"os"

	"github.com/mustardcustard/detoo/internal/sop"
	"github.com/mustardcustard/detoo/internal/tool"
)

func main() {
	os.Exit(tool.Run(sop.Kind, sop.Files))
}
