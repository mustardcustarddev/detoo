// Command plm records, links and reports on problem playbooks, each a plan,
// made with PIOSEE, for responding to a problem when it arises.
//
// Run "plm help" for the available commands.
//
// This package holds no behaviour of its own: it gives the plm tool's
// definition to the program every record tool shares, and exits with the
// status that returns.
package main

import (
	"os"

	"github.com/mustardcustard/detoo/internal/plm"
	"github.com/mustardcustard/detoo/internal/tool"
)

func main() {
	os.Exit(tool.Run(plm.Kind, plm.Files))
}
