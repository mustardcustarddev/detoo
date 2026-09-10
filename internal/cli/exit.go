package cli

import "fmt"

// ExitError reports that a command has already explained itself and the process
// should terminate with Code.
//
// It belongs to this package rather than to a slice: what a failure should do
// to a process is a property of running as a command, and an HTTP or MCP client
// consuming the same handlers has no use for it.
type ExitError struct {
	Code int
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("exit status %d", e.Code)
}

func exitWith(code int) error {
	return &ExitError{Code: code}
}
