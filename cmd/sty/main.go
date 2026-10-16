// Command sty records, links and reports on user stories.
//
// Run "sty help" for the available commands.
//
// This package holds no behaviour of its own: it gives the sty tool's
// definition to the program every record tool shares, and exits with the
// status that returns.
package main

import (
	"os"

	"github.com/mustardcustard/detoo/internal/sty"
	"github.com/mustardcustard/detoo/internal/tool"
)

func main() {
	os.Exit(tool.Run(sty.Kind, sty.Files))
}
