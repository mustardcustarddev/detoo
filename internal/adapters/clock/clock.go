// Package clock supplies the date written into a new record.
//
// It is a port rather than a call to time.Now inside a slice, because a handler
// that reads the clock directly produces a different record every day and
// cannot be tested without freezing time.
package clock

import (
	"context"
	"time"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// Clock reports today's date as it will appear in a record.
type Clock struct {
	getenv   func(string) string
	variable string
}

// New returns a Clock for logs of kind, which reads its configuration through
// getenv.
func New(getenv func(string) string, kind record.Kind) *Clock {
	return &Clock{getenv: getenv, variable: kind.Variable("DATE")}
}

// Today returns the date for a new record: the value of the tool's DATE
// variable, such as ADR_DATE, if it is set, and otherwise the current date in
// ISO 8601.
//
// The date is a string rather than a time.Time because the variable is passed
// through exactly as the user wrote it.  Parsing and reformatting it here would
// reject dates this tool has always accepted.
func (c *Clock) Today(context.Context) string {
	if date := c.getenv(c.variable); date != "" {
		return date
	}
	return time.Now().Format(time.DateOnly)
}
