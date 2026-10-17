// Package project looks at a project's directory as a whole: where it is, the
// marker files that say where its logs are kept, and the files in a directory,
// and changes it: it moves directories, renames files, and rewrites
// documents.
//
// This is a driven adapter.  Where logstore serves the one log a record tool
// keeps, this serves a command that looks after every log in a project.
package project

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
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

// Move moves the directory from to to, making the directories above to that
// it needs.  It refuses to move onto anything already at to.
func (p *Project) Move(ctx context.Context, from, to string) error {
	src, err := p.path(ctx, from)
	if err != nil {
		return err
	}
	dst, err := p.path(ctx, to)
	if err != nil {
		return err
	}
	if _, err := os.Lstat(dst); err == nil {
		return fmt.Errorf("cannot move %s to %s: %s exists", from, to, to)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return os.Rename(src, dst)
}

// Rename renames the file from, in dir, to to.  It refuses to rename over
// another file, but renames a file whose name differs only in case on a
// filesystem that ignores case, where the two names are the same file.
func (p *Project) Rename(ctx context.Context, dir, from, to string) error {
	full, err := p.path(ctx, dir)
	if err != nil {
		return err
	}
	src, dst := filepath.Join(full, from), filepath.Join(full, to)
	srcInfo, err := os.Lstat(src)
	if err != nil {
		return err
	}
	if dstInfo, err := os.Lstat(dst); err == nil {
		if !os.SameFile(srcInfo, dstInfo) {
			return fmt.Errorf("cannot rename %s to %s in %s: %s exists", from, to, dir, to)
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return os.Rename(src, dst)
}

// RemoveIfEmpty removes dir if nothing is left in it, and reports whether it
// did.
func (p *Project) RemoveIfEmpty(ctx context.Context, dir string) (bool, error) {
	full, err := p.path(ctx, dir)
	if err != nil {
		return false, err
	}
	entries, err := os.ReadDir(full)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	if len(entries) > 0 {
		return false, nil
	}
	return true, os.Remove(full)
}

// Markdown returns the Markdown files in the project's root, and in dirs and
// every directory beneath them, in order and each once.  Hidden files and
// directories are left out, and so is a directory that does not exist.
func (p *Project) Markdown(ctx context.Context, dirs []string) ([]string, error) {
	seen := map[string]bool{}
	var found []string
	add := func(name string) {
		if !seen[name] {
			seen[name] = true
			found = append(found, name)
		}
	}

	root, err := p.Root(ctx)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.Type().IsRegular() && isMarkdown(entry.Name()) {
			add(entry.Name())
		}
	}

	for _, dir := range dirs {
		full, err := p.path(ctx, dir)
		if err != nil {
			return nil, err
		}
		err = filepath.WalkDir(full, func(file string, entry fs.DirEntry, err error) error {
			if err != nil {
				if file == full && errors.Is(err, fs.ErrNotExist) {
					return fs.SkipAll
				}
				return err
			}
			if file != full && strings.HasPrefix(entry.Name(), ".") {
				if entry.IsDir() {
					return fs.SkipDir
				}
				return nil
			}
			if entry.Type().IsRegular() && isMarkdown(entry.Name()) {
				rel, err := filepath.Rel(full, file)
				if err != nil {
					return err
				}
				add(path.Join(path.Clean(dir), filepath.ToSlash(rel)))
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	sort.Strings(found)
	return found, nil
}

func isMarkdown(name string) bool {
	return !strings.HasPrefix(name, ".") && strings.EqualFold(filepath.Ext(name), ".md")
}

// Read returns the content of a file.
func (p *Project) Read(ctx context.Context, name string) (string, error) {
	full, err := p.path(ctx, name)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(full)
	return string(data), err
}

// Write replaces the content of a file, keeping its permissions.
func (p *Project) Write(ctx context.Context, name, content string) error {
	full, err := p.path(ctx, name)
	if err != nil {
		return err
	}
	info, err := os.Stat(full)
	if err != nil {
		return err
	}
	return os.WriteFile(full, []byte(content), info.Mode().Perm())
}
