// Package initialise starts a log of records, whose first record is the
// decision to keep one, or, for a kind with no first record, whose README says
// what the log is for.
package initialise

import (
	"context"
	"errors"
	"path"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// Command asks for a log.
type Command struct{}

// Result names the record, or the README, that was created.  File is empty
// when the log already had records, or its README, and nothing was written.
type Result struct {
	File string
}

// Store is the port this slice needs.
type Store interface {
	Dir(ctx context.Context) (string, error)
	Records(ctx context.Context) ([]record.Record, error)
	NextNumber(ctx context.Context, dir string) (int, error)
	Write(ctx context.Context, file, content string) error
	Exists(ctx context.Context, file string) (bool, error)
}

// Templates is the port that supplies the first record's Markdown, or the
// README's.
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

func (h Handler) Handle(ctx context.Context, _ Command) (Result, error) {
	// A log that already has records is left as it is, so that init can be run
	// again safely.
	records, err := h.Log.Records(ctx)
	var missing *record.MissingDirError
	switch {
	case errors.As(err, &missing):
		// No directory, no log: start one.
	case err != nil:
		return Result{}, err
	case len(records) > 0:
		return Result{}, nil
	}

	dir, err := h.Log.Dir(ctx)
	if err != nil {
		return Result{}, err
	}
	template, err := h.Templates.First(ctx)
	if err != nil {
		return Result{}, err
	}
	if h.Kind.FirstRecord == "" {
		return h.readme(ctx, dir, template)
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

// readme opens a log with a README, unless it has one, which may be the
// project's own, such as a table of contents.  The README is not a record, so
// it has no number, and its date would say only when init ran.
func (h Handler) readme(ctx context.Context, dir, content string) (Result, error) {
	file := path.Join(dir, record.Readme)
	exists, err := h.Log.Exists(ctx, file)
	if err != nil || exists {
		return Result{}, err
	}
	if err := h.Log.Write(ctx, file, content); err != nil {
		return Result{}, err
	}
	return Result{File: file}, nil
}
