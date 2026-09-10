// Package create implements the "create a user" capability.
//
// It imports no CLI package and no framework.  The invariant is checkable:
//
//	go list -deps ./internal/features/... | grep spf13/cobra   # must find nothing
package create

import (
	"context"
	"strings"
	"time"

	"github.com/example/myapp/internal/features/users"
)

// Repository is the port this slice needs, declared here by its consumer and
// kept to the two methods this use case actually calls.  users/get declares its
// own one-method Repository; the sibling declarations are not duplication, they
// are what stops this slice depending on reading and deleting users as well.
//
// No adapter names this interface.  Go satisfies it structurally, so
// *filesystem.UserRepository implements every slice's port at once.
type Repository interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user users.User) (users.User, error)
}

// Clock is a port too.  Calling time.Now directly would make the outcome depend
// on when the test runs, which is the same defect as calling os.Stdout.
type Clock interface {
	Now() time.Time
}

// Handler carries the slice's dependencies.  They are set once, by the
// composition root in internal/app, and never looked up from here.
type Handler struct {
	Users Repository
	Clock Clock
}

// Handle runs the use case.  Every slice in the application has this shape,
// which is what lets a second client reuse them without adaptation:
//
//	Handle(ctx context.Context, input) (Result, error)
func (h Handler) Handle(ctx context.Context, cmd Command) (Result, error) {
	// Validate the request, not the invocation.  Whether the name arrived as
	// one quoted argument or three loose ones is the CLI's problem; whether it
	// is empty is the application's, because every client must be told.
	name := strings.TrimSpace(cmd.Name)
	if name == "" {
		return Result{}, ErrNameRequired
	}
	email := strings.TrimSpace(cmd.Email)
	if !strings.Contains(email, "@") {
		return Result{}, ErrInvalidEmail
	}

	// Apply the business rules.
	role := cmd.Role
	if role == "" {
		role = users.RoleMember
	}
	if role == users.RoleAdmin && !strings.HasSuffix(email, "@example.com") {
		return Result{}, ErrExternalAdmin
	}

	taken, err := h.Users.ExistsByEmail(ctx, email)
	if err != nil {
		return Result{}, err
	}
	if taken {
		return Result{}, ErrEmailTaken
	}

	// The context comes from the caller, so a cancelled command stops the write.
	stored, err := h.Users.Create(ctx, users.User{
		Name:      name,
		Email:     email,
		Role:      role,
		CreatedAt: h.Clock.Now(),
	})
	if err != nil {
		return Result{}, err
	}

	return Result{
		ID:        stored.ID,
		Name:      stored.Name,
		Email:     stored.Email,
		Role:      stored.Role,
		CreatedAt: stored.CreatedAt,
	}, nil
}
