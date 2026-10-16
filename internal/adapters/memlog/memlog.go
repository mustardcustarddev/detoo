// Package memlog keeps a log of records in memory.
//
// It is a second driven adapter, satisfying the same ports as logstore, and it
// exists to make the boundary demonstrable rather than merely claimed: the
// slices are tested against this, with no temporary directory and no editor,
// and neither they nor their tests change when the real store is swapped in.
//
// A client that wanted to serve a log from somewhere other than a disk — an
// MCP server holding one per session, say — would write an adapter like this
// one and change nothing else.
package memlog

import (
	"context"
	"fmt"
	"maps"
	"path"
	"sort"
	"strings"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// Log is an in-memory log of records.  Its zero value is not usable: call
// New.
type Log struct {
	dir     string
	created bool
	files   map[string]string
}

// New returns an empty log that has not been initialised, so that Records
// reports a missing directory exactly as a store on an empty project would.
func New() *Log {
	return &Log{dir: "docs/adr", files: map[string]string{}}
}

// WithRecords returns a log holding the given records, keyed by path.
func WithRecords(records map[string]string) *Log {
	log := New()
	log.created = true
	maps.Copy(log.files, records)
	return log
}

// Content returns a record's text, for a test to assert on.
func (l *Log) Content(file string) string { return l.files[file] }

// Files returns the paths in the log, in order.
func (l *Log) Files() []string {
	files := make([]string, 0, len(l.files))
	for file := range l.files {
		files = append(files, file)
	}
	sort.Strings(files)
	return files
}

func (l *Log) Dir(context.Context) (string, error) { return l.dir, nil }

func (l *Log) SetDir(_ context.Context, dir string) error {
	l.dir, l.created = dir, true
	return nil
}

func (l *Log) Records(context.Context) ([]record.Record, error) {
	if !l.created {
		return nil, &record.MissingDirError{Dir: l.dir}
	}
	var files []string
	for file := range l.files {
		if path.Dir(file) == l.dir && record.IsRecordName(path.Base(file)) {
			files = append(files, file)
		}
	}
	sort.Strings(files)

	records := make([]record.Record, 0, len(files))
	for _, file := range files {
		records = append(records, record.Record{Number: record.Number(file), File: file})
	}
	return records, nil
}

func (l *Log) Resolve(ctx context.Context, ref string) (record.Record, error) {
	records, err := l.Records(ctx)
	if err != nil {
		return record.Record{}, err
	}
	for _, r := range records {
		if strings.Contains(r.File, ref) {
			return r, nil
		}
	}
	return record.Record{}, fmt.Errorf("no record matches %q", ref)
}

func (l *Log) NextNumber(_ context.Context, dir string) (int, error) {
	highest := 0
	for file := range l.files {
		if path.Dir(file) == dir {
			if number := record.Number(file); number > highest {
				highest = number
			}
		}
	}
	return highest + 1, nil
}

func (l *Log) Lines(_ context.Context, file string) ([]string, error) {
	content, ok := l.files[file]
	if !ok {
		return nil, fmt.Errorf("no such file %q", file)
	}
	text := strings.TrimSuffix(content, "\n")
	if text == "" {
		return nil, nil
	}
	return strings.Split(text, "\n"), nil
}

func (l *Log) WriteLines(_ context.Context, file string, lines []string) error {
	var text strings.Builder
	for _, line := range lines {
		text.WriteString(line)
		text.WriteByte('\n')
	}
	l.files[file] = text.String()
	return nil
}

func (l *Log) Write(_ context.Context, file, content string) error {
	l.files[file] = content
	l.created = true
	return nil
}

func (l *Log) ReadFile(_ context.Context, file string) (string, error) {
	content, ok := l.files[file]
	if !ok {
		return "", fmt.Errorf("no such file %q", file)
	}
	return content, nil
}
