// ✅ GOOD — internal/cli/users_create.go
//
// The command is a translation layer: it validates the *shape* of the
// invocation, maps arguments and flags onto a use-case input, calls one
// handler, and renders the result.  Every rule in bad_command.go now lives in
// internal/features/users/create, where an HTTP API or an MCP server can reach
// it too.  See ../_templates for the slice this calls.
package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/example/myapp/internal/app"
	"github.com/example/myapp/internal/cli/output"
	"github.com/example/myapp/internal/features/users/create"
)

func newUsersCreateCommand(application *app.Application, format *output.Format) *cobra.Command {
	var email, role string

	cmd := &cobra.Command{
		Use:   "create NAME",
		Short: "Create a user",

		// CLI-shaped validation only: how many operands, and whether a name
		// with spaces was quoted.  Neither question means anything to another
		// client, so neither belongs in the handler.
		Args: func(_ *cobra.Command, args []string) error {
			switch {
			case len(args) == 0:
				return errors.New("no name given")
			case len(args) > 1:
				return errors.New("quote the name so it is a single argument")
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := application.Users.Create.Handle(cmd.Context(), create.Command{
				Name:  args[0],
				Email: email,
				Role:  role,
			})
			if err != nil {
				// Translate an application error into a better CLI message.
				if errors.Is(err, create.ErrEmailTaken) {
					return fmt.Errorf("a user with the address %q already exists", email)
				}
				return err
			}
			return format.Write(cmd.OutOrStdout(), result)
		},
	}

	cmd.Flags().StringVar(&email, "email", "", "the user's email address")
	cmd.Flags().StringVar(&role, "role", "", "the user's role; defaults to member")
	return cmd
}
