// Package filesystem stores users in a JSON file.
//
// This is a driven adapter.  Note what it does not do: it names no port, and
// imports no slice's interface.  Go satisfies interfaces structurally, so
// *UserRepository implements create.Repository, get.Repository and
// delete.Repository at once, having said nothing about any of them.  Swapping
// this for package postgres is a one-line change in the composition root.
package filesystem

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/example/myapp/internal/features/users"
)

type UserRepository struct {
	path string
	mu   sync.Mutex
}

func NewUserRepository(path string) *UserRepository {
	return &UserRepository{path: path}
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	stored, err := r.load(ctx)
	if err != nil {
		return false, err
	}
	for _, user := range stored {
		if user.Email == email {
			return true, nil
		}
	}
	return false, nil
}

func (r *UserRepository) Create(ctx context.Context, user users.User) (users.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	stored, err := r.load(ctx)
	if err != nil {
		return users.User{}, err
	}
	user.ID = fmt.Sprintf("usr_%d", len(stored)+1)
	stored = append(stored, user)

	if err := r.save(stored); err != nil {
		return users.User{}, err
	}
	return user, nil
}

func (r *UserRepository) Get(ctx context.Context, id string) (users.User, error) {
	stored, err := r.load(ctx)
	if err != nil {
		return users.User{}, err
	}
	for _, user := range stored {
		if user.ID == id {
			return user, nil
		}
	}
	// The adapter reports that the record is absent.  It does not decide what
	// that means: get.Handler translates it into get.ErrUserNotFound, and the
	// CLI turns that into a message and an exit code.
	return users.User{}, os.ErrNotExist
}

func (r *UserRepository) load(ctx context.Context) ([]users.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(r.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var stored []users.User
	return stored, json.Unmarshal(data, &stored)
}

func (r *UserRepository) save(stored []users.User) error {
	data, err := json.MarshalIndent(stored, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(r.path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(r.path, data, 0o644)
}
