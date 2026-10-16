package create

// Command is the input to this use case.
//
// It describes the request in the application's own terms.  It is deliberately
// not a struct of flag values: the same Command is built by the Cobra adapter
// from arguments and flags, by an HTTP handler from a JSON body, and by an MCP
// server from a tool call.  Nothing here may hint at which.
type Command struct {
	Name  string
	Email string

	// Role is optional.  An empty Role means the default, which is the
	// application's rule to apply, not the CLI's.
	Role string
}
