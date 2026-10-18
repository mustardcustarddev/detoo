// Package skillstore installs skills on disk, where the skills CLI (npx
// skills) installs them, and records them in its lock files: in a project,
// in .agents/skills and skills-lock.json; for the user, in ~/.agents/skills
// and ~/.agents/.skill-lock.json.
//
// This is a driven adapter, serving the skills slice.  It writes the lock
// files as the skills CLI does, keeping what it does not know of them, so
// that either tool can update or remove what the other installed.
package skillstore

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/mustardcustard/detoo/internal/features/agents/skills"
)

// Store is where skills are installed: a project, or the user's home
// directory.
type Store struct {
	global  bool
	workdir string
	home    string
	getenv  func(string) string
	now     func() time.Time
	root    string // found by Root, once
}

// Project returns the store of the project holding workdir, for the user
// whose home directory is home, with their environment read by getenv.
func Project(workdir, home string, getenv func(string) string) *Store {
	return &Store{workdir: workdir, home: realPath(home), getenv: getenv, now: time.Now}
}

// User returns the store of the user whose home directory is home, with
// their environment read by getenv.
func User(home string, getenv func(string) string) *Store {
	return &Store{global: true, home: realPath(home), getenv: getenv, now: time.Now}
}

// realPath returns dir by its real path, as the project's directory is
// found, so that the two are compared alike whatever links lead to them.
func realPath(dir string) string {
	if real, err := filepath.EvalSymlinks(dir); err == nil {
		return real
	}
	return dir
}

// Root returns the project's directory, by its real path, or the user's home
// directory.  The project's is the nearest of the working directory and its
// ancestors, below the home directory, that holds skills-lock.json, .agents
// or .git.  A working directory with none of them above it is the project's
// directory itself.
func (s *Store) Root(context.Context) (string, error) {
	if s.global {
		return s.home, nil
	}
	if s.root != "" {
		return s.root, nil
	}
	start, err := filepath.Abs(s.workdir)
	if err != nil {
		return "", err
	}
	if real, err := filepath.EvalSymlinks(start); err == nil {
		start = real
	}
	s.root = start
	for dir := start; dir != s.home; dir = filepath.Dir(dir) {
		if isRoot(dir) {
			s.root = dir
			break
		}
		if filepath.Dir(dir) == dir {
			break
		}
	}
	return s.root, nil
}

func isRoot(dir string) bool {
	for _, name := range []string{projectLock, ".agents", ".git"} {
		if _, err := os.Lstat(filepath.Join(dir, name)); err == nil {
			return true
		}
	}
	return false
}

// path resolves a path the slice gives to one on disk: relative to the
// project's directory, unless it is absolute.
func (s *Store) path(ctx context.Context, name string) (string, error) {
	name = filepath.FromSlash(name)
	if filepath.IsAbs(name) {
		return name, nil
	}
	root, err := s.Root(ctx)
	if err != nil {
		return "", err
	}
	return filepath.Join(root, name), nil
}

// Canonical returns .agents/skills, in the project or the home directory.
func (s *Store) Canonical() string {
	if s.global {
		return s.place(skills.Global)
	}
	return skills.Canonical
}

// SkillsDir returns the directory agent reads, in the project or for the
// user.
func (s *Store) SkillsDir(agent skills.Agent) string {
	if s.global {
		return s.place(agent.Global)
	}
	return agent.Dir
}

// Detected reports whether agent's configuration directory exists.
func (s *Store) Detected(_ context.Context, agent skills.Agent) (bool, error) {
	_, err := os.Stat(filepath.FromSlash(s.place(agent.Config)))
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return err == nil, err
}

// place returns where p is, slash separated, as the skills CLI finds it.
func (s *Store) place(p skills.Place) string {
	base := s.home
	switch p.Base {
	case skills.ConfigHome:
		base = s.env("XDG_CONFIG_HOME", ".config")
	case skills.ClaudeHome:
		base = s.env("CLAUDE_CONFIG_DIR", ".claude")
	case skills.CodexHome:
		base = s.env("CODEX_HOME", ".codex")
	}
	return filepath.ToSlash(filepath.Join(base, p.Path))
}

// env returns the directory the variable called name gives, or else dir in
// the home directory.
func (s *Store) env(name, dir string) string {
	if v := s.getenv(name); v != "" {
		return v
	}
	return filepath.Join(s.home, dir)
}

// LockFile returns skills-lock.json, in the project, or the user's lock file:
// in $XDG_STATE_HOME/skills, or else ~/.agents.
func (s *Store) LockFile() string {
	if !s.global {
		return projectLock
	}
	if state := s.getenv("XDG_STATE_HOME"); state != "" {
		return filepath.ToSlash(filepath.Join(state, "skills", userLock))
	}
	return filepath.ToSlash(filepath.Join(s.home, ".agents", userLock))
}

// Hash identifies files as the lock file records them.
func (s *Store) Hash(files []skills.File) string {
	if s.global {
		return treeHash(files)
	}
	return folderHash(files)
}

// Install replaces whatever is at dir with a directory of files.
func (s *Store) Install(ctx context.Context, dir string, files []skills.File) error {
	full, err := s.path(ctx, dir)
	if err != nil {
		return err
	}
	if err := os.RemoveAll(full); err != nil {
		return err
	}
	if err := os.MkdirAll(full, 0o755); err != nil {
		return err
	}
	for _, f := range files {
		name := filepath.Join(full, filepath.FromSlash(f.Path))
		if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
			return err
		}
		mode := os.FileMode(0o644)
		if f.Executable {
			mode = 0o755
		}
		if err := os.WriteFile(name, f.Content, mode); err != nil {
			return err
		}
	}
	return nil
}

// Link makes link a relative link to the directory target, replacing whatever
// is there, unless link already leads to target.  Where a link cannot be
// made, it copies target instead.
func (s *Store) Link(ctx context.Context, link, target string) error {
	from, err := s.path(ctx, link)
	if err != nil {
		return err
	}
	to, err := s.path(ctx, target)
	if err != nil {
		return err
	}
	if same(from, to) {
		return nil
	}
	if err := os.RemoveAll(from); err != nil {
		return err
	}
	dir := filepath.Dir(from)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	// Relative to where the link really is, so that it holds when the
	// project moves.
	realDir, realTo := dir, to
	if r, err := filepath.EvalSymlinks(dir); err == nil {
		realDir = r
	}
	if r, err := filepath.EvalSymlinks(to); err == nil {
		realTo = r
	}
	rel, err := filepath.Rel(realDir, realTo)
	if err != nil {
		rel = realTo
	}
	if err := os.Symlink(rel, from); err == nil {
		return nil
	}
	return copyDir(to, from)
}

// same reports whether a and b are the same directory, by their real paths.
func same(a, b string) bool {
	ra, err := filepath.EvalSymlinks(a)
	if err != nil {
		return false
	}
	rb, err := filepath.EvalSymlinks(b)
	return err == nil && ra == rb
}

func copyDir(from, to string) error {
	return filepath.WalkDir(from, func(p string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(from, p)
		if err != nil {
			return err
		}
		dst := filepath.Join(to, rel)
		if entry.IsDir() {
			return os.MkdirAll(dst, 0o755)
		}
		if !entry.Type().IsRegular() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		content, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		return os.WriteFile(dst, content, info.Mode().Perm())
	})
}

// Entry reports what is at name: Linked when it is target, or a link to it,
// even one target is missing from; Copied when it is a directory of its own;
// and Foreign when it is anything else.
func (s *Store) Entry(ctx context.Context, name, target string) (skills.EntryKind, error) {
	full, err := s.path(ctx, name)
	if err != nil {
		return skills.Missing, err
	}
	to, err := s.path(ctx, target)
	if err != nil {
		return skills.Missing, err
	}
	info, err := os.Lstat(full)
	if errors.Is(err, fs.ErrNotExist) {
		return skills.Missing, nil
	} else if err != nil {
		return skills.Missing, err
	}
	if same(full, to) {
		return skills.Linked, nil
	}
	switch {
	case info.Mode()&fs.ModeSymlink != 0:
		dest, err := os.Readlink(full)
		if err != nil {
			return skills.Missing, err
		}
		if !filepath.IsAbs(dest) {
			dest = filepath.Join(filepath.Dir(full), dest)
		}
		if filepath.Clean(dest) == filepath.Clean(to) {
			return skills.Linked, nil
		}
	case info.IsDir():
		return skills.Copied, nil
	}
	return skills.Foreign, nil
}

// RemoveIfEmpty removes dir if nothing is left in it, and reports whether it
// did.
func (s *Store) RemoveIfEmpty(ctx context.Context, dir string) (bool, error) {
	full, err := s.path(ctx, dir)
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

// Remove removes name, and everything in it.  A link is removed, and not
// what it leads to.
func (s *Store) Remove(ctx context.Context, name string) error {
	full, err := s.path(ctx, name)
	if err != nil {
		return err
	}
	return os.RemoveAll(full)
}
