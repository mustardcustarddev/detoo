// Package templates supplies the Markdown a new decision record is created
// from.
//
// A driven adapter: it decides where a template comes from — an environment
// variable, the project's own file, or the one built into the binary — so that
// a slice can ask for "the template for a record in this log" and be given the
// text.
package templates

import (
	"context"
	_ "embed"
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/mustardcustard/decision-records/internal/domain/adr"
)

// defaultRecord is used for a new record when the project provides no template
// of its own.  It follows the style described by Michael Nygard in
// http://thinkrelevance.com/blog/2011/11/15/documenting-architecture-decisions
//
//go:embed templates/template.md
var defaultRecord string

// firstRecord is the template for the first record of a log, which records the
// decision to record decisions.
//
//go:embed templates/init.md
var firstRecord string

// Templates resolves record templates under a working directory.
type Templates struct {
	workdir string
	getenv  func(string) string
}

// New returns a Templates that reads the project's files under workdir and its
// configuration through getenv.
func New(workdir string, getenv func(string) string) *Templates {
	return &Templates{workdir: workdir, getenv: getenv}
}

func (t *Templates) path(rel string) string {
	return filepath.Join(t.workdir, filepath.FromSlash(rel))
}

// Record returns the template a new record in dir is created from: the one
// named by ADR_TEMPLATE, else the project's own, else the built-in default.
func (t *Templates) Record(_ context.Context, dir string) (string, error) {
	if file := t.getenv("ADR_TEMPLATE"); file != "" {
		return readFile(t.path(file))
	}
	project, err := readFile(t.path(path.Join(dir, adr.ProjectTemplate)))
	if errors.Is(err, fs.ErrNotExist) {
		return defaultRecord, nil
	}
	return project, err
}

// First returns the template for the record that opens a log.
//
// It is deliberately not affected by ADR_TEMPLATE: that variable overrides the
// template for the records a user writes, whereas this one states what the log
// is for, and is the same in every project.
func (t *Templates) First(context.Context) (string, error) {
	return firstRecord, nil
}

// ProjectFile returns the path of the project's own template for records in
// dir, relative to the working directory, or "" if the project has none.
func (t *Templates) ProjectFile(_ context.Context, dir string) (string, error) {
	file := path.Join(dir, adr.ProjectTemplate)
	if info, err := os.Stat(t.path(file)); err != nil || info.IsDir() {
		return "", nil
	}
	return file, nil
}

func readFile(file string) (string, error) {
	data, err := os.ReadFile(file)
	return string(data), err
}
