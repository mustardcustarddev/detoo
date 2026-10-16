// Package project looks at a project's directory as a whole: where it is, the
// marker files that say where its logs are kept, and the files in a directory.
//
// This is a driven adapter.  Where logstore serves the one log a record tool
// keeps, this serves a command that looks across every log in a project.
package project

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// Project is the project reachable from a working directory.  Every path it
// accepts and reports is relative to the project's root, unless it is
// absolute, and is slash separated, whatever the host filesystem uses.
type Project struct {
	workdir string
	kinds   []record.Kind
	root    string // found by Root, once
}

// New returns the project holding workdir, whose logs are of the given kinds.
func New(workdir string, kinds []record.Kind) *Project {
	return &Project{workdir: workdir, kinds: kinds}
}

// Root returns the project's directory, by its real path: the nearest of the
// working directory and its ancestors that holds a marker file of one of the
// kinds, a docs directory or a doc directory.  A working directory with none
// of them above it is the project's directory itself.
func (p *Project) Root(context.Context) (string, error) {
	if p.root != "" {
		return p.root, nil
	}
	start, err := filepath.Abs(p.workdir)
	if err != nil {
		return "", err
	}
	// The directory is reported by its real path, so that it compares
	// equal to others whatever links led to it.
	if real, err := filepath.EvalSymlinks(start); err == nil {
		start = real
	}
	p.root = start
	for dir := start; ; dir = filepath.Dir(dir) {
		if p.isRoot(dir) {
			p.root = dir
			break
		}
		if filepath.Dir(dir) == dir {
			break
		}
	}
	return p.root, nil
}

func (p *Project) isRoot(dir string) bool {
	for _, kind := range p.kinds {
		if _, err := os.Stat(filepath.Join(dir, markerFile(kind))); err == nil {
			return true
		}
	}
	for _, name := range []string{"docs", "doc"} {
		if info, err := os.Stat(filepath.Join(dir, name)); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

// markerFile names the file that records where a project keeps a log of
// kind: .adr-dir for adr's.
func markerFile(kind record.Kind) string {
	return "." + kind.Name + "-dir"
}

// path resolves a path relative to the project's root to a path on disk.  An
// absolute path is kept as it is.
func (p *Project) path(ctx context.Context, name string) (string, error) {
	name = filepath.FromSlash(name)
	if filepath.IsAbs(name) {
		return name, nil
	}
	root, err := p.Root(ctx)
	if err != nil {
		return "", err
	}
	return filepath.Join(root, name), nil
}

// Marker returns the directory the marker file of kind, in the project's
// root, names.
func (p *Project) Marker(ctx context.Context, kind record.Kind) (string, bool, error) {
	file, err := p.path(ctx, markerFile(kind))
	if err != nil {
		return "", false, err
	}
	data, err := os.ReadFile(file)
	if errors.Is(err, fs.ErrNotExist) {
		return "", false, nil
	} else if err != nil {
		return "", false, err
	}
	return path.Clean(filepath.ToSlash(strings.TrimSpace(string(data)))), true, nil
}

// Files returns the names of the files in dir, leaving out directories and
// hidden files, and whether dir exists.
func (p *Project) Files(ctx context.Context, dir string) ([]string, bool, error) {
	full, err := p.path(ctx, dir)
	if err != nil {
		return nil, false, err
	}
	info, err := os.Stat(full)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	if !info.IsDir() {
		return nil, false, nil
	}
	entries, err := os.ReadDir(full)
	if err != nil {
		return nil, false, err
	}
	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			names = append(names, entry.Name())
		}
	}
	return names, true, nil
}
