package adr

import (
	"os"
	"path"
)

const configHelp = `usage: adr config

Prints the configuration that the tools are using, as shell variable
assignments, so that it can be inspected or evaluated by other scripts:

  adr_dir       the directory holding the decision records, relative to the
                current working directory
  adr_template  the project's own template for new records, or empty if the
                built-in template is used
`

func runConfig(e *Env, args []string) error {
	dir, err := e.ADRDir()
	if err != nil {
		return err
	}
	template := ""
	if project := path.Join(dir, ProjectTemplate); fileExists(e.path(project)) {
		template = project
	}
	e.printf("adr_dir=%q\n", dir)
	e.printf("adr_template=%q\n", template)
	return nil
}

// fileExists reports whether a file is there to be read.
func fileExists(file string) bool {
	info, err := os.Stat(file)
	return err == nil && !info.IsDir()
}
