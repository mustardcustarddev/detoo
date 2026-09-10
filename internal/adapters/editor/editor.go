// Package editor opens a decision record in the user's editor.
//
// A driven adapter over another program.  The slice that creates a record says
// only that the record should be opened; which editor, and whether there is one
// at all, is settled here.
package editor

import (
	"context"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
)

// Editor runs the editor named by VISUAL, or by EDITOR if VISUAL is not set.
type Editor struct {
	workdir string
	getenv  func(string) string
	stdin   io.Reader
	stdout  io.Writer
	stderr  io.Writer
}

// New returns an Editor that runs in workdir, reads its configuration through
// getenv, and connects the editor to the given streams.
func New(workdir string, getenv func(string) string, stdin io.Reader, stdout, stderr io.Writer) *Editor {
	return &Editor{workdir: workdir, getenv: getenv, stdin: stdin, stdout: stdout, stderr: stderr}
}

// Edit opens file in the user's editor and waits for it to close.
//
// With neither VISUAL nor EDITOR set there is no editor to open, and the record
// is left for the user to find.  That is not an error: creating the record
// succeeded, and the path is reported either way.
func (e *Editor) Edit(ctx context.Context, file string) error {
	editor := e.getenv("VISUAL")
	if editor == "" {
		editor = e.getenv("EDITOR")
	}
	words := strings.Fields(editor)
	if len(words) == 0 {
		return nil
	}

	cmd := exec.CommandContext(ctx, words[0], append(words[1:], filepath.FromSlash(file))...)
	cmd.Dir = e.workdir
	cmd.Stdin = e.stdin
	cmd.Stdout = e.stdout
	cmd.Stderr = e.stderr
	return cmd.Run()
}
