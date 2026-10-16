// Package create adds the next numbered record to the log, together with its
// links to the records it replaces or updates.
package create

import (
	"context"
	"errors"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// ErrNoTitle reports that no title was given for the record.
//
// It is an application error, not a CLI one: every client must supply a title,
// however its input arrives.  How many arguments the user typed, and whether
// they quoted them, is the CLI's business and is checked there.
var ErrNoTitle = errors.New("no title given")

// Command asks for a new record.
type Command struct {
	Title string

	// Replaces holds references to records this one replaces.  Each is marked
	// as replaced, in the words of the kind of log, and loses the status a new
	// record is given.
	Replaces []string

	// Updates holds references to records this one makes substantive changes
	// to.  Each is marked updated, and keeps its status.
	Updates []string
}

// Result names the record that was created.
type Result struct {
	File string
}

// Store is the port this slice needs.
type Store interface {
	Dir(ctx context.Context) (string, error)
	NextNumber(ctx context.Context, dir string) (int, error)
	Write(ctx context.Context, file, content string) error
	Resolve(ctx context.Context, ref string) (record.Record, error)
	Lines(ctx context.Context, file string) ([]string, error)
	WriteLines(ctx context.Context, file string, lines []string) error
}

// Templates is the port that supplies the Markdown a record starts from.
type Templates interface {
	Record(ctx context.Context, dir string) (string, error)
}

// Clock is the port that supplies the record's date.
type Clock interface {
	Today(ctx context.Context) string
}

// Editor is the port that opens the new record for the user to write.
type Editor interface {
	Edit(ctx context.Context, file string) error
}

// Handler creates a record.
type Handler struct {
	Kind      record.Kind
	Log       Store
	Templates Templates
	Clock     Clock
	Editor    Editor
}

func (h Handler) Handle(ctx context.Context, cmd Command) (Result, error) {
	if cmd.Title == "" {
		return Result{}, ErrNoTitle
	}

	dir, err := h.Log.Dir(ctx)
	if err != nil {
		return Result{}, err
	}
	template, err := h.Templates.Record(ctx, dir)
	if err != nil {
		return Result{}, err
	}
	number, err := h.Log.NextNumber(ctx, dir)
	if err != nil {
		return Result{}, err
	}

	created := record.Record{Number: number, File: record.FileName(dir, number, cmd.Title)}
	content := record.Render(template, number, cmd.Title, h.Clock.Today(ctx), h.Kind.Status)
	if err := h.Log.Write(ctx, created.File, content); err != nil {
		return Result{}, err
	}

	// The record is written before its links, so that a link to it resolves.
	for _, ref := range cmd.Replaces {
		if err := h.relate(ctx, created, record.Replaces, ref); err != nil {
			return Result{}, err
		}
	}
	for _, ref := range cmd.Updates {
		if err := h.relate(ctx, created, record.Updates, ref); err != nil {
			return Result{}, err
		}
	}

	// Opening the editor is the last thing done, and its failure is reported,
	// but the record exists either way and the caller is told where it is.
	if err := h.Editor.Edit(ctx, created.File); err != nil {
		return Result{}, err
	}
	return Result{File: created.File}, nil
}

// relate records that the new record replaces or updates an existing one: the
// old record gains a link to the new one, and loses its status if the relation
// retires it, and the new record records the old.
func (h Handler) relate(ctx context.Context, created record.Record, relation record.Relation, ref string) error {
	target, err := h.Log.Resolve(ctx, ref)
	if err != nil {
		return err
	}
	words := h.Kind.WordsFor(relation)
	if err := h.addLink(ctx, target, words.Reverse, created); err != nil {
		return err
	}
	if relation.Retires() {
		if err := h.removeStatus(ctx, target, h.Kind.Status); err != nil {
			return err
		}
	}
	return h.addLink(ctx, created, words.Forward, target)
}

// addLink appends a link from source to target, labelled with the target's
// title.  The link slice keeps its own copy: see the note there.
func (h Handler) addLink(ctx context.Context, source record.Record, linkType string, target record.Record) error {
	targetLines, err := h.Log.Lines(ctx, target.File)
	if err != nil {
		return err
	}
	sourceLines, err := h.Log.Lines(ctx, source.File)
	if err != nil {
		return err
	}
	return h.Log.WriteLines(ctx, source.File, record.AddLink(sourceLines, linkType, record.Title(targetLines), target.File))
}

func (h Handler) removeStatus(ctx context.Context, target record.Record, status string) error {
	lines, err := h.Log.Lines(ctx, target.File)
	if err != nil {
		return err
	}
	return h.Log.WriteLines(ctx, target.File, record.RemoveStatus(lines, status))
}
