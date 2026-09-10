// ❌ BAD — internal/cli/users_create.go
//
// Every common vertical-slice violation in one Cobra command.  There is no
// application here: the command *is* the application, so none of this behaviour
// can be reached by an HTTP API, an MCP server, or a test that does not parse
// flags.
package cli

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func newUsersCreateCommand(db *sql.DB) *cobra.Command {
	var name, email, role string

	return &cobra.Command{
		Use:   "create",
		Short: "Create a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Application validation in the CLI.  An HTTP client needs exactly
			// these checks and cannot reach them.
			if strings.TrimSpace(name) == "" {
				fmt.Println("error: name is required")
				os.Exit(1) // the CLI decides the exit code mid-behaviour
			}
			if !strings.Contains(email, "@") {
				fmt.Println("error: invalid email")
				os.Exit(1)
			}

			// A business rule, stranded in a flag parser.
			if role == "" {
				role = "member"
			}
			if role == "admin" && !strings.HasSuffix(email, "@example.com") {
				fmt.Println("error: admins must have an internal address")
				os.Exit(1)
			}

			// Infrastructure reached directly from the command, with a context
			// invented on the spot, so ^C cannot cancel the write.
			ctx := context.Background()
			row := db.QueryRowContext(ctx,
				"INSERT INTO users (name, email, role, created_at) VALUES ($1,$2,$3,$4) RETURNING id",
				name, email, role, time.Now())

			var id string
			if err := row.Scan(&id); err != nil {
				fmt.Println("error:", err)
				os.Exit(1)
			}

			// Audit trail written by the presentation layer.
			f, _ := os.OpenFile("audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
			defer f.Close()
			fmt.Fprintf(f, "%s created %s\n", time.Now().Format(time.RFC3339), id)

			// Rendering hard-wired to stdout, so tests must capture the process
			// rather than the command.
			out, _ := json.Marshal(map[string]string{"id": id, "name": name})
			fmt.Println(string(out))
			return nil
		},
	}
}
