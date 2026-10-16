// Command myapp manages users and projects.
//
// This is the only package main, and it holds no application behaviour.  Its
// whole job is the process boundary: build the adapters, wire the application,
// execute the command tree, and turn an error into an exit code.  That is why
// no other package needs to call os.Exit.
package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/example/myapp/internal/adapters/filesystem"
	"github.com/example/myapp/internal/app"
	"github.com/example/myapp/internal/cli"
	"github.com/example/myapp/internal/features/users/create"
)

func main() {
	os.Exit(run())
}

func run() int {
	// The composition root: the concrete adapters are chosen here, once, and
	// nothing downstream knows which they are.
	application := app.New(filesystem.NewUserRepository("users.json"), systemClock{})

	if err := cli.NewRootCommand(application).Execute(); err != nil {
		// A few application errors deserve a particular exit code, so that
		// scripts can distinguish them.  The application said what went wrong;
		// this is where the process decides what to do about it.
		if errors.Is(err, create.ErrEmailTaken) {
			return fail(err, 2)
		}
		return fail(err, 1)
	}
	return 0
}

// fail reports an error the way the Go tools do: on stderr, named.
func fail(err error, code int) int {
	fmt.Fprintf(os.Stderr, "myapp: %v\n", err)
	return code
}

// systemClock is the real implementation of the clock port the slices declare.
type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now() }
