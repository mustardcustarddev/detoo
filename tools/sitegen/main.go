// Command sitegen writes the documentation site's command reference: a
// Markdown page for each tool, and for each of its commands, holding exactly
// what the tool prints as its help.
//
// The pages are generated rather than written, so that the site cannot say
// something different from the tools.  They are built in, not committed: the
// site's build runs this first.
//
// Usage:
//
//	go run ./tools/sitegen [DIR]
//
// DIR defaults to docs/site/src/content/reference, and is emptied first, so
// that a command that has gone leaves no page behind.
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	dir := filepath.FromSlash("docs/site/src/content/reference")
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	if err := write(dir); err != nil {
		fmt.Fprintf(os.Stderr, "sitegen: %v\n", err)
		os.Exit(1)
	}
}

// write replaces the contents of dir with the reference pages.
func write(dir string) error {
	pages, err := Pages()
	if err != nil {
		return err
	}
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	for _, p := range pages {
		path := filepath.Join(dir, filepath.FromSlash(p.Path()))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(p.Markdown()), 0o644); err != nil {
			return err
		}
	}
	return nil
}
