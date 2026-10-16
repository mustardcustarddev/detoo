package create

import "errors"

// The errors this use case can fail with.
//
// They are sentinel values so that a client can test for them with errors.Is
// and choose its own response: the CLI rewrites them into a better message and
// an exit code, an HTTP API maps them onto status codes.  The application says
// what went wrong; it never says what should happen next.
var (
	ErrNameRequired = errors.New("name is required")
	ErrInvalidEmail = errors.New("email is not a valid address")
	ErrEmailTaken   = errors.New("email is already registered")

	// ErrExternalAdmin reports the rule that an administrator must hold an
	// internal address.  A rule like this is the reason the slice exists: put
	// it in a Cobra RunE and no other client can enforce it.
	ErrExternalAdmin = errors.New("an administrator must have an internal address")
)
