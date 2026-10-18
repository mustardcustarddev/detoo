// Command ust records, links and reports on user stories.
//
// Run "ust help" for the available commands.
//
// This package holds no behaviour of its own: it gives the ust tool's
// definition to the program every record tool shares, and exits with the
// status that returns.
package main

import (
	"os"

	"github.com/mustardcustard/detoo/internal/tool"
	"github.com/mustardcustard/detoo/internal/ust"
)

func main() {
	os.Exit(tool.Run(ust.Kind, ust.Files))
}
