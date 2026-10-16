package create

import "time"

// Result is the outcome of the use case, as structured data.
//
// The slice does not decide how this is shown.  The CLI renders it as text,
// JSON or a table; an HTTP API marshals it; a worker may only log it.  Adding a
// String method here, or a `json:` tag chosen to suit one client, moves a
// presentation decision into the application.
type Result struct {
	ID        string
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
}
