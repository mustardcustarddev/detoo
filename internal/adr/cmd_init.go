package adr

import (
	"os"
	"path/filepath"
)

const initHelp = `usage: adr init [DIRECTORY]

Initialises the directory of architecture decision records:

 * creates a subdirectory of the current working directory
 * creates the first ADR in that subdirectory, recording the decision to
   record architectural decisions with ADRs.

If the DIRECTORY is not given, the ADRs are stored in the directory ` + "`doc/adr`" + `.
`

func runInit(e *Env, args []string) error {
	if len(args) > 0 && args[0] != "" {
		dir := args[0]
		if err := os.MkdirAll(e.path(dir), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(e.Dir, DirMarkerFile), []byte(dir+"\n"), 0o644); err != nil {
			return err
		}
	}
	file, err := e.New("record-architecture-decisions", NewOptions{Template: InitTemplate, NoEdit: true})
	if err != nil {
		return err
	}
	e.print(file)
	return nil
}
