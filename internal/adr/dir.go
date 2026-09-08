package adr

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	// DirMarkerFile names the file that records where a project keeps its
	// decision records.
	DirMarkerFile = ".adr-dir"

	// DefaultDir is the directory used by projects without a marker file.
	DefaultDir = "doc/adr"
)

// ADRDir returns the directory holding the decision records, as a path relative
// to the working directory.  The working directory and each of its ancestors is
// searched, in turn, for a marker file naming the directory, and then for a
// directory at the default location.  Neither needs to exist: projects that have
// no records yet get the default.
func (e *Env) ADRDir() (string, error) {
	dir, err := filepath.Abs(e.Dir)
	if err != nil {
		return "", err
	}
	for up := ""; ; dir, up = filepath.Dir(dir), path.Join(up, "..") {
		marker, err := os.ReadFile(filepath.Join(dir, DirMarkerFile))
		switch {
		case err == nil:
			return path.Join(up, filepath.ToSlash(strings.TrimSpace(string(marker)))), nil
		case !errors.Is(err, fs.ErrNotExist):
			return "", err
		}
		if info, err := os.Stat(filepath.Join(dir, filepath.FromSlash(DefaultDir))); err == nil && info.IsDir() {
			return path.Join(up, DefaultDir), nil
		}
		if parent := filepath.Dir(dir); parent == dir {
			return DefaultDir, nil
		}
	}
}
