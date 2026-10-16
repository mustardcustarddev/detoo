// Package app is the composition root: it constructs every slice handler and
// hands the adapters they depend on to them.
//
// This is the one package that necessarily knows the whole application, and it
// will grow a field per use case.  That is correct for a composition root, and
// it is also how such a type becomes a god object.  The rule that prevents it:
// Application holds wiring only.  It must never acquire methods.  The moment
// behaviour is hung off it, a caller that wants one capability starts depending
// on all of them — which is exactly the coupling the slices were for.
package app

import (
	"github.com/example/myapp/internal/features/users/create"
	"github.com/example/myapp/internal/features/users/delete"
	"github.com/example/myapp/internal/features/users/get"
)

// Application groups the slice handlers by capability, so that a caller reads
// as application.Users.Create.Handle(ctx, cmd).
type Application struct {
	Users Users
}

type Users struct {
	Create create.Handler
	Delete delete.Handler
	Get    get.Handler
}

// New wires the application.  Its parameters are the concrete adapters; each
// slice sees only the narrow interface it declared, so nothing below this line
// knows that users are stored in a file rather than a database.
func New(repo UserRepository, clock Clock) *Application {
	return &Application{
		Users: Users{
			Create: create.Handler{Users: repo, Clock: clock},
			Delete: delete.Handler{Users: repo},
			Get:    get.Handler{Users: repo},
		},
	}
}

// UserRepository is the union of what the users slices ask for, named here only
// so that New has a signature.  The adapters satisfy it without importing this
// package, and no slice ever refers to it.
type UserRepository interface {
	create.Repository
	get.Repository
	delete.Repository
}

// Clock is likewise a convenience for the wiring, not a port a slice depends on.
type Clock = create.Clock
