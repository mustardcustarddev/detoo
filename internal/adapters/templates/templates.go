// Package templates supplies the Markdown a new record is created from.
//
// A driven adapter: it decides where a template comes from — an environment
// variable, the project's own file, or the one the tool ships with — so that a
// slice can ask for "the template for a record in this log" and be given the
// text.
package templates

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// The templates a tool ships with, within its files.
const (
	// builtinRecord is used for a new record when the project provides no
	// template of its own.
	builtinRecord = "templates/template.md"

	// builtinFirst is the template for the record that opens a log, which
	// records the decision to keep one.
	builtinFirst = "templates/init.md"
)

// Templates resolves the templates of one kind of log under a working
// directory.
type Templates struct {
	workdir string
	getenv  func(string) string
	kind    record.Kind
	files   fs.FS
}

// New returns a Templates for logs of kind.  It reads the project's files under
// workdir, its configuration through getenv, and the templates the tool ships
// with from files.
func New(workdir string, getenv func(string) string, kind record.Kind, files fs.FS) *Templates {
	return &Templates{workdir: workdir, getenv: getenv, kind: kind, files: files}
}

// path resolves a file named relative to the working directory, or by an
// absolute path, which is kept as it is.
func (t *Templates) path(name string) string {
	name = filepath.FromSlash(name)
	if filepath.IsAbs(name) {
		return name
	}
	return filepath.Join(t.workdir, name)
}

// Record returns the template a new record in dir is created from: the one
// named by the tool's TEMPLATE variable, such as ADR_TEMPLATE, else the
// project's own, else the one the tool ships with.
func (t *Templates) Record(_ context.Context, dir string) (string, error) {
	if file := t.getenv(t.kind.Variable("TEMPLATE")); file != "" {
		return readFile(t.path(file))
	}
	project, err := readFile(t.path(path.Join(dir, record.ProjectTemplate)))
	if errors.Is(err, fs.ErrNotExist) {
		return t.builtin(builtinRecord)
	}
	return project, err
}

// First returns the template for the record that opens a log.
//
// It is deliberately not affected by the TEMPLATE variable: that overrides the
// template for the records a user writes, whereas this one states what the log
// is for, and is the same in every project.
func (t *Templates) First(context.Context) (string, error) {
	return t.builtin(builtinFirst)
}

// ProjectFile returns the path of the project's own template for records in
// dir, relative to the working directory, or "" if the project has none.
func (t *Templates) ProjectFile(_ context.Context, dir string) (string, error) {
	file := path.Join(dir, record.ProjectTemplate)
	if info, err := os.Stat(t.path(file)); err != nil || info.IsDir() {
		return "", nil
	}
	return file, nil
}

// builtin returns one of the templates the tool ships with.
func (t *Templates) builtin(name string) (string, error) {
	data, err := fs.ReadFile(t.files, name)
	return string(data), err
}

func readFile(file string) (string, error) {
	data, err := os.ReadFile(file)
	return string(data), err
}
