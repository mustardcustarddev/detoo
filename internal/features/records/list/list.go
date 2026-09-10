// Package list reports the decision records in the log.
package list

import (
	"context"

	"github.com/mustardcustard/decision-records/internal/domain/adr"
)

// Query asks for the log's records.  It has no fields yet, and exists so that
// the slice reads as every other one does, and so that a filter can be added
// later without changing the shape of the call.
type Query struct{}

// Result is the log, in order.
//
// It carries records rather than the lines a terminal would show, because the
// caller decides what to make of them: the CLI prints one path a line, and
// another client might render titles or serialise the numbers.
type Result struct {
	Records []adr.Record
}

// Store is the port this slice needs.
type Store interface {
	Records(ctx context.Context) ([]adr.Record, error)
}

// Handler lists the records in a log.
type Handler struct {
	Log Store
}

func (h Handler) Handle(ctx context.Context, _ Query) (Result, error) {
	records, err := h.Log.Records(ctx)
	if err != nil {
		return Result{}, err
	}
	return Result{Records: records}, nil
}
