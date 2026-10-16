// Package link records that one record replaces or updates another, in both
// records.
package link

import (
	"context"
	"errors"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// ErrNoRelation reports that a link was asked for without saying whether it
// replaces or updates its target.
//
// It is an application error: a link with no relation has nothing to write,
// however the request arrived.  Which flags the user typed is the CLI's
// business, and is checked there.
var ErrNoRelation = errors.New("a link must either replace or update a record")

// Command asks for a link from one record to another.
//
// Source and Target are references — a record's number, or any part of its file
// name — because that is how a user names a record.  Resolving one is the
// store's work, not the caller's.
type Command struct {
	Relation record.Relation // how Source bears on Target
	Source   string          // the record that replaces or updates
	Target   string          // the record replaced or updated
}

// Result names the records that were linked.
type Result struct {
	Source string
	Target string
}

// Store is the port this slice needs.
type Store interface {
	Resolve(ctx context.Context, ref string) (record.Record, error)
	Lines(ctx context.Context, file string) ([]string, error)
	WriteLines(ctx context.Context, file string, lines []string) error
}

// Handler links two records to each other.
type Handler struct {
	Kind record.Kind
	Log  Store
}

func (h Handler) Handle(ctx context.Context, cmd Command) (Result, error) {
	if !cmd.Relation.Valid() {
		return Result{}, ErrNoRelation
	}
	source, err := h.Log.Resolve(ctx, cmd.Source)
	if err != nil {
		return Result{}, err
	}
	target, err := h.Log.Resolve(ctx, cmd.Target)
	if err != nil {
		return Result{}, err
	}

	// Both records are written: a link the reader can only follow one way is
	// half a link, and the reverse is what makes the log navigable.
	words := h.Kind.WordsFor(cmd.Relation)
	if err := h.addLink(ctx, source, words.Forward, target); err != nil {
		return Result{}, err
	}
	if err := h.addLink(ctx, target, words.Reverse, source); err != nil {
		return Result{}, err
	}
	if cmd.Relation.Retires() {
		if err := h.removeStatus(ctx, target, h.Kind.Status); err != nil {
			return Result{}, err
		}
	}
	return Result{Source: source.File, Target: target.File}, nil
}

// addLink appends a link from source to target, labelled with the target's
// title, to the source record's status section.
//
// Creating a record does the same thing when it replaces or updates an existing
// one, and has its own copy of these few lines, and of removeStatus.  The rules
// they share — where in the Markdown a link goes, and which relations retire a
// record — are in package record, and are not duplicated.  What is repeated is
// the store choreography around them, which is small, and the price of keeping
// either slice readable on its own.  Exporting this so that create could call
// it would make one slice depend on another.
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
