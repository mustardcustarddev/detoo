// Package create adds the next numbered decision record to the log, together
// with any links it has to the records it supercedes or amends.
package create

import (
	"context"
	"errors"

	"github.com/mustardcustard/decision-records/internal/domain/adr"
)

// ErrNoTitle reports that no title was given for the record.
//
// It is an application error, not a CLI one: every client must supply a title,
// however its input arrives.  How many arguments the user typed, and whether
// they quoted them, is the CLI's business and is checked there.
var ErrNoTitle = errors.New("no title given")

// Link is a link from the new record to an existing one.
//
// It is structured rather than the "TARGET:LINK:REVERSE-LINK" string the -l
// option takes, because that encoding belongs to the command line.  The CLI
// splits it; this slice never learns that a colon was involved.
type Link struct {
	Target  string // a reference to the record linked to
	Forward string // how the new record describes the link
	Reverse string // how the target describes it
}

// Command asks for a new record.
type Command struct {
	Title string

	// Supercedes holds references to records this one replaces.  Each is
	// marked superceded, and loses its "Accepted" status.
	Supercedes []string

	// Links holds further links from the new record to existing ones.
	Links []Link
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
	Resolve(ctx context.Context, ref string) (adr.Record, error)
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

// Handler creates a decision record.
type Handler struct {
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

	record := adr.Record{Number: number, File: adr.FileName(dir, number, cmd.Title)}
	content := adr.Render(template, number, cmd.Title, h.Clock.Today(ctx), adr.StatusAccepted)
	if err := h.Log.Write(ctx, record.File, content); err != nil {
		return Result{}, err
	}

	// The record is written before its links, so that a link to it resolves.
	for _, ref := range cmd.Supercedes {
		if err := h.supercede(ctx, record, ref); err != nil {
			return Result{}, err
		}
	}
	for _, link := range cmd.Links {
		if err := h.link(ctx, record, link); err != nil {
			return Result{}, err
		}
	}

	// Opening the editor is the last thing done, and its failure is reported,
	// but the record exists either way and the caller is told where it is.
	if err := h.Editor.Edit(ctx, record.File); err != nil {
		return Result{}, err
	}
	return Result{File: record.File}, nil
}

// supercede marks an existing record replaced by the new one: the old record
// gains a link forward and loses its "Accepted" status, and the new record
// records what it replaced.
func (h Handler) supercede(ctx context.Context, record adr.Record, ref string) error {
	target, err := h.Log.Resolve(ctx, ref)
	if err != nil {
		return err
	}
	if err := h.addLink(ctx, target, "Superceded by", record); err != nil {
		return err
	}
	if err := h.removeStatus(ctx, target, adr.StatusAccepted); err != nil {
		return err
	}
	return h.addLink(ctx, record, "Supercedes", target)
}

// link records a link both ways between the new record and an existing one.
func (h Handler) link(ctx context.Context, record adr.Record, link Link) error {
	target, err := h.Log.Resolve(ctx, link.Target)
	if err != nil {
		return err
	}
	if err := h.addLink(ctx, record, link.Forward, target); err != nil {
		return err
	}
	return h.addLink(ctx, target, link.Reverse, record)
}

// addLink appends a link from source to target, labelled with the target's
// title.  The link slice keeps its own copy: see the note there.
func (h Handler) addLink(ctx context.Context, source adr.Record, linkType string, target adr.Record) error {
	targetLines, err := h.Log.Lines(ctx, target.File)
	if err != nil {
		return err
	}
	sourceLines, err := h.Log.Lines(ctx, source.File)
	if err != nil {
		return err
	}
	return h.Log.WriteLines(ctx, source.File, adr.AddLink(sourceLines, linkType, adr.Title(targetLines), target.File))
}

func (h Handler) removeStatus(ctx context.Context, record adr.Record, status string) error {
	lines, err := h.Log.Lines(ctx, record.File)
	if err != nil {
		return err
	}
	return h.Log.WriteLines(ctx, record.File, adr.RemoveStatus(lines, status))
}
