// Package logstore stores a log of records in a directory of Markdown files.
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

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// Store is a log of records under a working directory.  Every path it accepts
// and reports is relative to that directory, unless the project keeps its log
// at an absolute path, and is slash separated, whatever the host filesystem
// uses.
type Store struct {
	workdir string
	kind    record.Kind
}

// New returns a Store on the log of the given kind reachable from workdir.
func New(workdir string, kind record.Kind) *Store {
	return &Store{workdir: workdir, kind: kind}
}

// markerFile names the file that records where a project keeps its log:
// .adr-dir for adr's.  Each tool has its own, so that one project can keep the
// logs of several.
func (s *Store) markerFile() string {
	return "." + s.kind.Name + "-dir"
}

// defaultDir is the directory a project without a marker file keeps its log in:
// docs/adr for adr's.
func (s *Store) defaultDir() string {
	return path.Join("docs", s.kind.Dir)
}

// path resolves a path reported by this package to a path on disk.  An
// absolute path is kept as it is.
func (s *Store) path(name string) string {
	name = filepath.FromSlash(name)
	if filepath.IsAbs(name) {
		return name
	}
	return filepath.Join(s.workdir, name)
}

// Dir returns the directory holding the records.  The working
// directory and each of its ancestors is searched, in turn, for a marker file
// naming the directory, and then for a directory at the default location.
// Neither needs to exist: projects that have no records yet get the default.
func (s *Store) Dir(context.Context) (string, error) {
	dir, err := filepath.Abs(s.workdir)
	if err != nil {
		return "", err
	}
	for up := ""; ; dir, up = filepath.Dir(dir), path.Join(up, "..") {
		marker, err := os.ReadFile(filepath.Join(dir, s.markerFile()))
		switch {
		case err == nil:
			// A relative directory is relative to the marker file, and an
			// absolute one is the same from anywhere.
			named := strings.TrimSpace(string(marker))
			if filepath.IsAbs(named) {
				return filepath.ToSlash(named), nil
			}
			return path.Join(up, filepath.ToSlash(named)), nil
		case !errors.Is(err, fs.ErrNotExist):
			return "", err
		}
		if info, err := os.Stat(filepath.Join(dir, filepath.FromSlash(s.defaultDir()))); err == nil && info.IsDir() {
			return path.Join(up, s.defaultDir()), nil
		}
		if parent := filepath.Dir(dir); parent == dir {
			return s.defaultDir(), nil
		}
	}
}

// SetDir records dir as where this project keeps its log, and
// creates it.  Dir then finds it from anywhere in the project.
func (s *Store) SetDir(_ context.Context, dir string) error {
	if err := os.MkdirAll(s.path(dir), 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(s.workdir, s.markerFile()), []byte(dir+"\n"), 0o644)
}

// Records returns the records, in order.
//
// A log directory that is not there is reported as a *record.MissingDirError
// rather than as an fs error: that the project has no log is an application
// condition, and every client needs to recognise it.
func (s *Store) Records(ctx context.Context) ([]record.Record, error) {
	dir, err := s.Dir(ctx)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(s.path(dir))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, &record.MissingDirError{Dir: dir}
	} else if err != nil {
		return nil, err
	}

	var records []record.Record
	for _, entry := range entries {
		if !entry.IsDir() && s.kind.IsRecordName(entry.Name()) {
			file := path.Join(dir, entry.Name())
			records = append(records, record.Record{Number: s.kind.Number(file), File: file})
		}
	}

	// Ordered by number rather than by name, because a log may hold records
	// named before today's prefix, such as 0001-first.md beside
	// ADR-0002-second.md.
	sort.Slice(records, func(i, j int) bool {
		if records[i].Number != records[j].Number {
			return records[i].Number < records[j].Number
		}
		return records[i].File < records[j].File
	})
	if records == nil {
		records = []record.Record{}
	}
	return records, nil
}

// Resolve returns the record referred to by ref: see record.Find.
func (s *Store) Resolve(ctx context.Context, ref string) (record.Record, error) {
	records, err := s.Records(ctx)
	if err != nil {
		return record.Record{}, err
	}
	if r, ok := record.Find(records, ref); ok {
		return r, nil
	}
	return record.Record{}, fmt.Errorf("no %s matches %q", s.kind.Noun, ref)
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
		if number := s.kind.Number(entry.Name()); number > highest {
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
