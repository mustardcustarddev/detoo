// Package initialise starts a log of records, whose first record is the
// decision to keep one.
package initialise

import (
	"context"
	"errors"
	"fmt"
	"path"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// ErrLogExists reports that init was asked to start a log in one directory
// while the project's log, which already has records, is in another.
var ErrLogExists = errors.New("a log already exists")

// logExistsError is ErrLogExists, told in the words of the kind of log and
// naming the directory the log is in.
type logExistsError struct {
	plural, dir string
}

func (e *logExistsError) Error() string {
	return fmt.Sprintf("a log of %s already exists in %s", e.plural, e.dir)
}

func (e *logExistsError) Is(target error) bool { return target == ErrLogExists }

// Command asks for a log.
type Command struct {
	// Dir is where the log should live, relative to the working directory.
	// Empty leaves the location as it is, so that a project which already
	// chose one keeps it.
	Dir string
}

// Result names the record that was created.  File is empty when the log
// already had records, and nothing was written.
type Result struct {
	File string
}

// Store is the port this slice needs.
type Store interface {
	SetDir(ctx context.Context, dir string) error
	Dir(ctx context.Context) (string, error)
	Records(ctx context.Context) ([]record.Record, error)
	NextNumber(ctx context.Context, dir string) (int, error)
	Write(ctx context.Context, file, content string) error
}

// Templates is the port that supplies the first record's Markdown.
type Templates interface {
	First(ctx context.Context) (string, error)
}

// Clock is the port that supplies the record's date.
type Clock interface {
	Today(ctx context.Context) string
}

// Handler starts a log.
//
// It writes its record itself rather than calling the create slice.  The two
// look alike for a few lines and then diverge: this record comes from a fixed
// template, has nothing to link to, and is not opened for editing, because the
// user did not sit down to write it.  Sharing the code would mean options that
// exist only to switch that behaviour off.
type Handler struct {
	Kind      record.Kind
	Log       Store
	Templates Templates
	Clock     Clock
}

func (h Handler) Handle(ctx context.Context, cmd Command) (Result, error) {
	// A log that already has records is left as it is, so that init can be run
	// again safely.  Starting one somewhere else would strand those records in
	// a directory the project no longer uses, so that is refused.
	current, err := h.Log.Dir(ctx)
	if err != nil {
		return Result{}, err
	}
	records, err := h.Log.Records(ctx)
	var missing *record.MissingDirError
	switch {
	case errors.As(err, &missing):
		// No directory, no log: start one.
	case err != nil:
		return Result{}, err
	case len(records) > 0:
		if cmd.Dir != "" && path.Clean(cmd.Dir) != path.Clean(current) {
			return Result{}, &logExistsError{plural: h.Kind.Plural, dir: current}
		}
		return Result{}, nil
	}

	if cmd.Dir != "" {
		if err := h.Log.SetDir(ctx, cmd.Dir); err != nil {
			return Result{}, err
		}
	}

	dir, err := h.Log.Dir(ctx)
	if err != nil {
		return Result{}, err
	}
	template, err := h.Templates.First(ctx)
	if err != nil {
		return Result{}, err
	}
	number, err := h.Log.NextNumber(ctx, dir)
	if err != nil {
		return Result{}, err
	}

	// The title only names the file: the record's heading is in the template.
	file := h.Kind.FileName(dir, number, h.Kind.FirstRecord)
	content := record.Render(template, number, h.Kind.FirstRecord, h.Clock.Today(ctx), h.Kind.Status)
	if err := h.Log.Write(ctx, file, content); err != nil {
		return Result{}, err
	}
	return Result{File: file}, nil
}
