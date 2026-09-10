---
name: feature-driven-vertical-slices
description: >-
  Use when building or refactoring Go command-line applications with Cobra, to organise code by
  capability rather than technical layer and keep the application headless, so the same behaviour
  can be driven by a CLI, an HTTP API, an MCP server, a worker, or a test.
license: Apache 2.0
compatibility: Go >= 1.21, github.com/spf13/cobra >= 1.8
metadata:
  author: tom regan
  version: 2.0.0
---

# Feature-Driven Vertical Slices for Cobra Applications

## Purpose

Build Go CLI applications as **feature-driven vertical slices** with
[Cobra](https://github.com/spf13/cobra).

The application must be **headless**: its behaviour must not depend on Cobra, on
terminal I/O, or on any other CLI concern.  Cobra is an adapter that translates
user input into calls on the application's feature slices.

```text
CLI / Cobra
     │
     ▼
Application features
     │
     ▼
Ports (interfaces)
     ▲
     │
Infrastructure adapters
```

The CLI is a client of the application, not the application itself.

## Prerequisites

- A Go module using `github.com/spf13/cobra` for its command tree.
- A single `package main` under `cmd/`, holding no application behaviour.

## Companion files

- `_examples/bad_command.go` — every violation in this document, in one `RunE`.
- `_examples/good_command.go` — the same command as a translation layer.
- `_templates/` — a working slice and the packages around it, laid out as the
  structure below: `cmd/myapp`, `cli`, `cli/output`, `app`, `features/users`,
  `features/users/create` (with its test), and `adapters/filesystem`.

The directory names begin with `_` because the Go tool ignores such directories
when matching `./...`.  That keeps illustrative code out of `go build`, `go vet`
and `go test` while it stays a real `.go` file for editors and `gofmt`.  Copy
the files without the underscore.

---

## Core Principles

### 1. Organise by capability, not by technical layer

Prefer:

```text
internal/
└── features/
    ├── users/
    │   ├── create/
    │   ├── delete/
    │   └── get/
    └── projects/
        ├── create/
        ├── delete/
        └── list/
```

Avoid organising the application primarily as:

```text
internal/
├── handlers/
├── services/
├── repositories/
├── models/
└── validators/
```

A slice holds the code needed to implement one coherent user capability, from
input to outcome.

### 2. Keep Cobra at the boundary

Cobra commands are adapters.  They are responsible for:

- defining commands, flags and arguments
- parsing and validating CLI input
- invoking application features
- rendering results
- translating application errors into CLI-facing errors

They must **not** contain:

- business rules
- database queries or filesystem manipulation
- HTTP requests
- domain workflows
- application state

Bad:

```go
RunE: func(cmd *cobra.Command, args []string) error {
    if name == "" {
        return errors.New("name required")
    }
    db.Exec("INSERT INTO users ...")
    return nil
},
```

Good:

```go
RunE: func(cmd *cobra.Command, args []string) error {
    result, err := app.Users.Create.Handle(cmd.Context(), create.Command{Name: name})
    if err != nil {
        return err
    }
    return out.Write(cmd.OutOrStdout(), result)
},
```

### 3. The application must be headless

Feature code must be usable without a terminal.  Inside `internal/features/`, do
not use:

```go
fmt.Println(...)      // the caller decides how to render
os.Exit(...)          // the caller decides the exit code
cobra.Command         // the caller decides how input arrives
```

Return values and errors instead:

```go
type Result struct {
    ID   string
    Name string
}

func (h Handler) Handle(ctx context.Context, cmd Command) (Result, error) {
    // application behaviour
}
```

The same slice can then be consumed by a Cobra CLI, an HTTP API, an MCP server,
a background worker, a test, or another Go program, with no logic duplicated.

This rule is mechanically checkable.  It should hold at all times:

```sh
go list -deps ./internal/features/... | grep -q spf13/cobra && echo "LEAK"
```

---

## Project Structure

```text
myapp/
├── cmd/
│   └── myapp/
│       └── main.go                 # build root, execute, map errors → exit code
│
└── internal/
    ├── cli/                        # Cobra adapter: the ONLY package importing cobra
    │   ├── root.go                 # assembles the command tree
    │   ├── users.go                # the "users" parent command
    │   ├── users_create.go         # flags/args → create.Command → render Result
    │   ├── users_get.go
    │   ├── projects_list.go
    │   └── output/                 # human / json / table renderers
    │
    ├── app/
    │   └── application.go          # Application: one field per slice, wiring only
    │
    ├── features/                   # The headless application
    │   ├── users/
    │   │   ├── user.go             # concepts every users slice shares
    │   │   ├── create/
    │   │   │   ├── command.go      # Command — the use-case input
    │   │   │   ├── handler.go      # Handle(ctx, Command) (Result, error)
    │   │   │   ├── result.go
    │   │   │   └── errors.go
    │   │   ├── delete/
    │   │   └── get/
    │   │       ├── query.go
    │   │       └── handler.go
    │   └── projects/
    │       ├── create/
    │       └── list/
    │
    └── adapters/                   # implement the ports slices declare
        ├── postgres/
        └── filesystem/
```

Everything but `cmd/` lives under `internal/`.  That is not a convention: it is
Go's only compiler-enforced boundary, and it stops another module importing the
command tree or a feature handler.  Put a package outside `internal/` only when
you intend to support third parties importing it forever.

**Scale the structure to the application.**  Below roughly two capabilities the
directories are ceremony, and a flat package with `new.go` beside `cmd_new.go`
is correct.  Introduce `features/` when the second capability arrives, not
before.  Do not create a directory merely to satisfy the diagram.

---

## Vertical Slices

A slice represents one capability from input to outcome:

```text
users/create

Command
    │
    ▼
Handler
    ├── validate
    ├── apply business rules
    ├── call a port
    └── construct Result
```

The input and output types:

```go
package create

type Command struct {
    Name string
}

type Result struct {
    ID   string
    Name string
}
```

The handler, with the port it needs declared beside it:

```go
package create

// Repository is the narrow port this slice requires.  It is declared here, by
// the consumer, so that the slice depends on nothing it does not use.
type Repository interface {
    Create(ctx context.Context, user users.User) error
}

type Handler struct {
    Users Repository
}

func (h Handler) Handle(ctx context.Context, cmd Command) (Result, error) {
    name := strings.TrimSpace(cmd.Name)
    if name == "" {
        return Result{}, ErrNameRequired
    }

    user := users.User{ID: newID(), Name: name}
    if err := h.Users.Create(ctx, user); err != nil {
        return Result{}, err
    }

    return Result{ID: user.ID, Name: user.Name}, nil
}
```

Every slice exposes the same small API, which is what makes slices
interchangeable between clients:

```go
Handle(ctx context.Context, input) (Result, error)
```

### Commands and queries

Use a **command** for an operation that changes state — `create`, `delete`,
`update`, `archive`, `restore` — and a **query** for one that retrieves
information — `get`, `list`, `search`, `status`.  A query's input type is named
`Query` rather than `Command`; the handler is otherwise identical.

Do not introduce a generic `Service` abstraction to conform to a pattern.
Prefer explicit use-case handlers.

---

## Ports and Adapters

**Declare a port in the slice that consumes it, and keep it narrow.**  Go
interfaces are satisfied structurally, so one concrete type implements every
slice's port without naming any of them:

```go
// internal/features/users/get/handler.go
type Repository interface {
    Get(ctx context.Context, id string) (users.User, error)
}

// internal/features/users/create/handler.go
type Repository interface {
    Create(ctx context.Context, user users.User) error
}

// internal/adapters/postgres/user_repository.go
// *postgres.UserRepository satisfies both, having said nothing about either.
```

The near-duplicate one-method interfaces are idiomatic Go, not redundancy.  They
are what lets `users/get` depend on reading a user without also depending on
creating and deleting one.

**Do not create a central `internal/ports/` package.**  It couples every slice
to one file, forces each slice to depend on the whole CRUD surface, and becomes
the shared bucket this skill otherwise warns against.

Create a port when it marks a meaningful dependency boundary — the filesystem,
a database, the network, the clock.  Do not create an interface for every type
merely for testability.

---

## Cobra Integration

Cobra commands map naturally onto slices:

```text
myapp users create   →   internal/features/users/create
myapp users delete   →   internal/features/users/delete
```

A command performs translation and nothing else:

```go
func newUsersCreateCommand(app *app.Application) *cobra.Command {
    var name string

    cmd := &cobra.Command{
        Use:   "create NAME",
        Short: "Create a user",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            result, err := app.Users.Create.Handle(cmd.Context(), create.Command{
                Name: args[0],
            })
            if err != nil {
                return err
            }
            fmt.Fprintln(cmd.OutOrStdout(), result.ID)
            return nil
        },
    }

    cmd.Flags().StringVar(&name, "email", "", "the user's email address")
    return cmd
}
```

The boundary is:

```text
Cobra flag → create.Command → create.Handler → Result → CLI output
```

The feature must not know that its input came from a flag.

### Where validation belongs

Both sides validate, and the cut between them is precise:

- **The CLI validates the shape of the invocation** — how many operands, whether
  a title was quoted, whether a flag parses as an integer, whether two flags
  conflict.  This belongs in a Cobra `Args` validator or in flag parsing, and it
  is meaningless to any other client.
- **The application validates the request** — whether a name is empty, whether a
  record exists, whether a transition is permitted.  This belongs in the
  handler, because an HTTP or MCP client needs exactly the same check.

The two look alike and are not:

```go
// CLI: quoting and arity. An HTTP client has no concept of either.
Args: func(cmd *cobra.Command, args []string) error {
    switch {
    case len(args) == 0:
        return errors.New("no name given")
    case len(args) > 1:
        return errors.New("quote the name so it is a single argument")
    }
    return nil
},

// Application: a rule about the request itself, true for every client.
if strings.TrimSpace(cmd.Name) == "" {
    return Result{}, ErrNameRequired
}
```

When in doubt, ask whether an MCP server would need the same check.  If it
would, the check belongs in the handler.

---

## Application Composition

Group the slice handlers behind one application object:

```go
package app

type Application struct {
    Users    Users
    Projects Projects
}

type Users struct {
    Create create.Handler
    Delete delete.Handler
    Get    get.Handler
}
```

The composition root constructs the dependencies:

```go
func New(users *postgres.UserRepository, projects *postgres.ProjectRepository) *Application {
    return &Application{
        Users: Users{
            Create: create.Handler{Users: users},
            Delete: delete.Handler{Users: users},
            Get:    get.Handler{Users: users},
        },
        Projects: Projects{
            Create: projectscreate.Handler{Projects: projects},
        },
    }
}
```

`Application` necessarily knows every slice, and will grow a field per use case.
That is correct for a composition root — but it is also how such a type becomes
a god object.  **The rule is that `Application` holds wiring only.**  It must
never acquire methods of its own; the moment behaviour is hung off it, callers
start depending on the whole application to reach one capability.

The exact composition mechanism can vary.  What matters is that dependency
construction happens outside the feature.

---

## Domain Models

Use domain types where they carry meaningful invariants or behaviour.  Do not
create a large model hierarchy automatically.  For a simple application a
capability package may reasonably hold:

```go
package users

type User struct {
    ID   string
    Name string
}
```

A complex domain may earn dedicated packages under `internal/domain/`.

Keep domain concepts separate from CLI representations.  In particular, never
use a struct of Cobra flag values as a domain model.

---

## CLI Output

Rendering is a presentation concern.  Features return structured results:

```go
type Result struct {
    ID        string
    Name      string
    CreatedAt time.Time
}
```

The CLI decides how to show them — human-readable, JSON, a table, CSV, or quiet
script-friendly output:

```sh
myapp users get 123                  # ID: 123 / Name: Alice / Created: 2026-09-09
myapp users get 123 --output json    # {"id":"123", ...}
```

A feature must not be responsible for either representation.  Write to
`cmd.OutOrStdout()` rather than `os.Stdout`, so that tests can capture output
without a process of their own.

---

## Context

Propagate the command's context:

```go
app.Users.Create.Handle(cmd.Context(), command)
```

Do not create a `context.Background()` inside a command without a specific
architectural reason.  Propagation is what allows cancellation, deadlines and
tracing to work.

---

## Error Handling

Features return meaningful application errors:

```go
var (
    ErrNameRequired = errors.New("name is required")
    ErrUserNotFound = errors.New("user not found")
)
```

The CLI translates them where a better message is possible:

```go
if errors.Is(err, get.ErrUserNotFound) {
    return fmt.Errorf("user %q not found", id)
}
```

Never call `os.Exit` from feature code.  Return errors to the process boundary
and let `main` choose the exit code.

---

## Testing

Test feature behaviour independently of Cobra:

```text
internal/features/users/create/
    handler.go
    handler_test.go
```

Feature tests need no temporary directory, no subprocess and no flag parsing —
substitute a fake or in-memory implementation for the slice's port.

CLI tests separately verify command registration, flags and arguments,
CLI-specific validation, output formatting, and exit behaviour.  Do not
reimplement the feature test suite through Cobra.

```text
Feature tests  →  application behaviour
CLI tests      →  translation and presentation
```

---

## Dependency Direction

```text
cmd/
  │
  ▼
internal/cli/
  │
  ▼
internal/app/
  │
  ▼
internal/features/          declares the interfaces it needs
  ▲
  │                         satisfied structurally, injected at the composition root
internal/adapters/
```

- Features must not import Cobra.
- Features must not import CLI packages.
- The CLI may import feature and application packages.
- Adapters implement the ports features declare, and are named by neither.

---

## Adding a New Feature

1.  **Identify the capability.**  `myapp projects archive` is the slice
    `projects/archive`.
2.  **Create the slice**, adding only the files it needs:

        internal/features/projects/archive/
            command.go
            handler.go
            errors.go

3.  **Define the use-case input.**

        type Command struct { ProjectID string }

4.  **Implement the handler**, declaring the narrow port it requires.  All
    application behaviour lives here.
5.  **Register it on the application.**

        type Projects struct { Archive archive.Handler }

6.  **Add the Cobra adapter** in `internal/cli/projects_archive.go`, mapping
    arguments and flags to the `Command`, and the `Result` to output.
7.  **Test the feature first**, then add focused CLI tests.

---

## What Not To Do

**Do not create a generic service layer.**  Avoid `handlers/ → services/ →
repositories/` where the service adds no abstraction.  Prefer
`users/create → Handler → port`.

**Do not put business logic in Cobra.**  A `RunE` that consults `user.IsAdmin`
or executes a query has taken over the application's job.

**Do not let features depend on each other casually.**  If `users/create` calls
`projects/create`, the slices are no longer independently understandable.
Extract the shared business concept into the capability package or a domain
package instead.

**Do not create an enormous shared package.**  `internal/common/`,
`internal/shared/`, `internal/utils/`, `internal/helpers/` accumulate anything
that resembles anything else.  Shared code should emerge from a genuine shared
concept, not from similarity.

---

## Additional Clients

The point of the boundary is that a second client is cheap:

```text
   Cobra CLI    HTTP API    MCP server    Worker
        │           │            │           │
        └───────────┴─────┬──────┴───────────┘
                          ▼
                 Headless application
                 (feature slices, domain)
                          │
                          ▼
                  Ports  →  Adapters
```

An MCP server is treated exactly as the CLI is:

```text
MCP tool invocation → command/query → handler → result → MCP response
```

Do not put application behaviour into MCP tool handlers.

---

## Architectural Test

When reviewing code, ask:

> Could I invoke this application's behaviour without running Cobra?

If not, application logic has leaked into the CLI.

> If I replaced Cobra with an HTTP API or an MCP server, would I have to rewrite
> the business logic?

If so, the boundary is in the wrong place.

And verify it, rather than trusting the reading:

```sh
go list -deps ./internal/features/... | grep spf13/cobra   # must find nothing
```

---

## Troubleshooting

**Import cycle between `internal/cli` and a feature.**  In Go this is a build
failure, not a smell.  It nearly always means shared context — an environment
struct, a clock, an output writer — was placed in `cli/`, which features then
needed.  Move it down to a package that imports nothing of yours.

**Two slices need the same repository interface.**  Leave the duplicate
declarations.  They are narrower than a shared one and cost nothing.  Hoist only
the *domain type* they exchange into the capability package.

**A slice has grown a dozen files.**  It is probably two capabilities.  Split it
by use case, not by technical layer.

**A use case is named after a Go builtin.**  `delete`, `copy`, `new`, `len`,
`max` and `min` are predeclared identifiers, not keywords, so `package delete`
is legal — but importing it shadows the builtin for the whole file, and
`delete(m, k)` then fails to compile in the composition root.  Either accept
the shadowing where the file never needs the builtin, or name the slice after
the domain verb: `remove`, `archive`, `revoke`.

**Flag names collide between commands.**  Flags are defined per command, so this
is only a problem for persistent flags on the root.  Keep the root's flags to
genuinely global concerns such as `--output` and `--verbose`.

---

## Guiding Rule

**Organise around capabilities, keep the application headless, and treat Cobra
as an adapter rather than as the application itself.**
