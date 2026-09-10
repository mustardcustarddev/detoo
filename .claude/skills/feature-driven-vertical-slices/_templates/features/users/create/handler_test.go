package create

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/example/myapp/internal/features/users"
)

// These tests are the point of the architecture.  There is no binary, no flag
// parsing, no temporary directory and no subprocess: the fakes below are a few
// lines each because the ports the slice declares are narrow.  Cobra tests
// cover translation and rendering separately, and do not repeat any of this.

type fakeRepo struct {
	created []users.User
	taken   bool
	err     error
}

func (f *fakeRepo) ExistsByEmail(context.Context, string) (bool, error) {
	return f.taken, f.err
}

func (f *fakeRepo) Create(_ context.Context, user users.User) (users.User, error) {
	if f.err != nil {
		return users.User{}, f.err
	}
	user.ID = "usr_1"
	f.created = append(f.created, user)
	return user, nil
}

type fixedClock struct{ t time.Time }

func (c fixedClock) Now() time.Time { return c.t }

func newHandler(repo *fakeRepo) Handler {
	return Handler{Users: repo, Clock: fixedClock{t: time.Date(2026, 9, 9, 0, 0, 0, 0, time.UTC)}}
}

func TestHandleCreatesAUserWithTheDefaultRole(t *testing.T) {
	repo := &fakeRepo{}

	result, err := newHandler(repo).Handle(context.Background(), Command{
		Name:  "Ada Lovelace",
		Email: "ada@example.com",
	})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if result.ID != "usr_1" {
		t.Errorf("ID = %q, want usr_1", result.ID)
	}
	if result.Role != users.RoleMember {
		t.Errorf("Role = %q, want %q", result.Role, users.RoleMember)
	}
	if len(repo.created) != 1 {
		t.Fatalf("stored %d users, want 1", len(repo.created))
	}
}

func TestHandleRejectsARequestWithNoName(t *testing.T) {
	_, err := newHandler(&fakeRepo{}).Handle(context.Background(), Command{
		Name:  "   ",
		Email: "ada@example.com",
	})
	if !errors.Is(err, ErrNameRequired) {
		t.Errorf("err = %v, want %v", err, ErrNameRequired)
	}
}

func TestHandleRejectsAnAdministratorWithAnExternalAddress(t *testing.T) {
	_, err := newHandler(&fakeRepo{}).Handle(context.Background(), Command{
		Name:  "Ada Lovelace",
		Email: "ada@elsewhere.net",
		Role:  users.RoleAdmin,
	})
	if !errors.Is(err, ErrExternalAdmin) {
		t.Errorf("err = %v, want %v", err, ErrExternalAdmin)
	}
}

func TestHandleRejectsAnAddressAlreadyRegistered(t *testing.T) {
	_, err := newHandler(&fakeRepo{taken: true}).Handle(context.Background(), Command{
		Name:  "Ada Lovelace",
		Email: "ada@example.com",
	})
	if !errors.Is(err, ErrEmailTaken) {
		t.Errorf("err = %v, want %v", err, ErrEmailTaken)
	}
}
