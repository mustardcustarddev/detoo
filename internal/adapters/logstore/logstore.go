// Package logstore stores a log of decision records in a directory of Markdown
// files.
//
// This is a driven adapter.  It names no port: every slice declares the handful
// of methods it needs, and *Store satisfies each of them structurally, so
// nothing here knows which use case is calling.  Replacing it with a store
// backed by a database or a remote API is a change to the composition root
// alone.
package logstore

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

	"github.com/mustardcustard/decision-records/internal/domain/adr"
)

const (
	// DirMarkerFile names the file that records where a project keeps its
	// decision records.
	DirMarkerFile = ".adr-dir"

	// DefaultDir is the directory used by projects without a marker file.
	DefaultDir = "docs/adr"
)

// Store is a log of decision records under a working directory.  Every path it
// accepts and reports is relative to that directory, and slash separated,
// whatever the host filesystem uses.
type Store struct {
	workdir string
}

// New returns a Store on the log reachable from workdir.
func New(workdir string) *Store {
	return &Store{workdir: workdir}
}

// path resolves a path reported by this package to a path on disk.
func (s *Store) path(rel string) string {
	return filepath.Join(s.workdir, filepath.FromSlash(rel))
}

// Dir returns the directory holding the decision records.  The working
// directory and each of its ancestors is searched, in turn, for a marker file
// naming the directory, and then for a directory at the default location.
// Neither needs to exist: projects that have no records yet get the default.
func (s *Store) Dir(context.Context) (string, error) {
	dir, err := filepath.Abs(s.workdir)
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

// SetDir records dir as where this project keeps its decision records, and
// creates it.  Dir then finds it from anywhere in the project.
func (s *Store) SetDir(_ context.Context, dir string) error {
	if err := os.MkdirAll(s.path(dir), 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(s.workdir, DirMarkerFile), []byte(dir+"\n"), 0o644)
}

// Records returns the decision records, in order.
//
// A log directory that is not there is reported as an *adr.MissingDirError
// rather than as an fs error: that the project has no decision log is an
// application condition, and every client needs to recognise it.
func (s *Store) Records(ctx context.Context) ([]adr.Record, error) {
	dir, err := s.Dir(ctx)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(s.path(dir))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, &adr.MissingDirError{Dir: dir}
	} else if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && adr.IsRecordName(entry.Name()) {
			files = append(files, path.Join(dir, entry.Name()))
		}
	}
	sort.Strings(files)

	records := make([]adr.Record, 0, len(files))
	for _, file := range files {
		records = append(records, adr.Record{Number: adr.Number(file), File: file})
	}
	return records, nil
}

// Resolve returns the record referred to by ref, which is either a record's
// number or any part of its file name.  The first matching record wins.
func (s *Store) Resolve(ctx context.Context, ref string) (adr.Record, error) {
	records, err := s.Records(ctx)
	if err != nil {
		return adr.Record{}, err
	}
	for _, record := range records {
		if strings.Contains(record.File, ref) {
			return record, nil
		}
	}
	return adr.Record{}, fmt.Errorf("no decision record matches %q", ref)
}

// NextNumber returns the number for the next record in dir, one higher than the
// highest already there.
func (s *Store) NextNumber(_ context.Context, dir string) (int, error) {
	entries, err := os.ReadDir(s.path(dir))
	if errors.Is(err, fs.ErrNotExist) {
		return 1, nil
	} else if err != nil {
		return 0, err
	}
	highest := 0
	for _, entry := range entries {
		if number := adr.Number(entry.Name()); number > highest {
			highest = number
		}
	}
	return highest + 1, nil
}

// Lines reads a record as lines, dropping the final line terminator so that
// WriteLines round-trips it.
func (s *Store) Lines(_ context.Context, file string) ([]string, error) {
	data, err := os.ReadFile(s.path(file))
	if err != nil {
		return nil, err
	}
	text := strings.TrimSuffix(string(data), "\n")
	if text == "" {
		return nil, nil
	}
	return strings.Split(text, "\n"), nil
}

// WriteLines writes lines to a record, each terminated by a newline.
func (s *Store) WriteLines(_ context.Context, file string, lines []string) error {
	var text strings.Builder
	for _, line := range lines {
		text.WriteString(line)
		text.WriteByte('\n')
	}
	return os.WriteFile(s.path(file), []byte(text.String()), 0o644)
}

// Write creates a record, and the log directory if it is not there yet.
func (s *Store) Write(_ context.Context, file, content string) error {
	if err := os.MkdirAll(filepath.Dir(s.path(file)), 0o755); err != nil {
		return err
	}
	return os.WriteFile(s.path(file), []byte(content), 0o644)
}

// ReadFile reads a file in the project that is not a record, such as the
// Markdown that surrounds a table of contents.
func (s *Store) ReadFile(_ context.Context, file string) (string, error) {
	data, err := os.ReadFile(s.path(file))
	return string(data), err
}
