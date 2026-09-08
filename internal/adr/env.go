// Package adr implements the architecture decision record tooling: a library of
// operations over a directory of decision records, and the commands that expose
// them on the command line.
package adr

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Env is the context a command runs in: the working directory it resolves
// relative paths against, where it writes its output, and how it reads the
// environment.  Everything a command touches goes through here, so commands can
// be driven from a test without a process of their own.
type Env struct {
	Dir    string
	Stdout io.Writer
	Stderr io.Writer
	Getenv func(string) string
}

// NewEnv returns an Env wired to the current process.
func NewEnv() (*Env, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &Env{Dir: dir, Stdout: os.Stdout, Stderr: os.Stderr, Getenv: os.Getenv}, nil
}

// path resolves a path reported by this package (always relative to the working
// directory, always slash separated) to a path on disk.
func (e *Env) path(rel string) string {
	return filepath.Join(e.Dir, filepath.FromSlash(rel))
}

func (e *Env) getenv(name string) string {
	if e.Getenv == nil {
		return ""
	}
	return e.Getenv(name)
}

func (e *Env) print(a ...any) {
	fmt.Fprintln(e.Stdout, a...)
}

func (e *Env) printf(format string, a ...any) {
	fmt.Fprintf(e.Stdout, format, a...)
}

// ExitError reports that a command has already explained itself and the process
// should terminate with Code.
type ExitError struct {
	Code int
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("exit status %d", e.Code)
}

func exitWith(code int) error {
	return &ExitError{Code: code}
}

// UsageError reports a command invoked with the wrong arguments.
type UsageError struct {
	Usage string
}

func (e *UsageError) Error() string {
	return "usage: " + e.Usage
}
