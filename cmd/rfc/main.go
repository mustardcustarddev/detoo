// Command rfc records, links and reports on requests for comments.
//
// Run "rfc help" for the available commands.
//
// This package holds no behaviour of its own: it gives the rfc tool's
// definition to the program every record tool shares, and exits with the
// status that returns.
package main

import (
	"os"

	"github.com/mustardcustard/detoo/internal/rfc"
	"github.com/mustardcustard/detoo/internal/tool"
)

func main() {
	os.Exit(tool.Run(rfc.Kind, rfc.Files))
}
