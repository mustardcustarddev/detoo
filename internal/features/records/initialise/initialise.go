// Package initialise starts a log of decision records, whose first record is
// the decision to keep one.
package initialise

import (
	"context"

	"github.com/mustardcustard/decision-records/internal/domain/adr"
)

// firstRecordTitle names the record that opens a log.  It is the title the file
// is named after; the record's own heading comes from the template.
const firstRecordTitle = "record-architecture-decisions"

// Command asks for a log.
type Command struct {
	// Dir is where the log should live, relative to the working directory.
	// Empty leaves the location as it is, so that a project which already
	// chose one keeps it.
	Dir string
}

// Result names the record that was created.
type Result struct {
	File string
}

// Store is the port this slice needs.
type Store interface {
	SetDir(ctx context.Context, dir string) error
	Dir(ctx context.Context) (string, error)
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
	Log       Store
	Templates Templates
	Clock     Clock
}

func (h Handler) Handle(ctx context.Context, cmd Command) (Result, error) {
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

	file := adr.FileName(dir, number, firstRecordTitle)
	content := adr.Render(template, number, firstRecordTitle, h.Clock.Today(ctx), adr.StatusAccepted)
	if err := h.Log.Write(ctx, file, content); err != nil {
		return Result{}, err
	}
	return Result{File: file}, nil
}
