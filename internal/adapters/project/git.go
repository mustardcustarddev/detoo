package project

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Git is a project in a git repository, which moves and renames with git mv,
// so that the repository records each as a rename.  A directory or a file git
// does not track is moved as Project moves it, since git mv refuses to.
type Git struct {
	*Project
}

// NewGit returns project, moving with git.
func NewGit(project *Project) *Git {
	return &Git{Project: project}
}

// Repository reports whether the project is in a git working tree.  Without
// git installed, it is not.
func (g *Git) Repository(ctx context.Context) (bool, error) {
	out, err := g.git(ctx, "rev-parse", "--is-inside-work-tree")
	if errors.Is(err, exec.ErrNotFound) {
		return false, nil
	}
	if _, failed := errors.AsType[*exec.ExitError](err); failed {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) == "true", nil
}

// Move moves the directory from to to with git mv, which takes the files git
// does not track along with those it does, or as Project does when git tracks
// nothing in it.
func (g *Git) Move(ctx context.Context, from, to string) error {
	src, err := g.path(ctx, from)
	if err != nil {
		return err
	}
	tracked, err := g.tracked(ctx, src)
	if err != nil {
		return err
	}
	if !tracked {
		return g.Project.Move(ctx, from, to)
	}
	dst, err := g.path(ctx, to)
	if err != nil {
		return err
	}
	if _, err := os.Lstat(dst); err == nil {
		return fmt.Errorf("cannot move %s to %s: %s exists", from, to, to)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	// git mv does not make the directory it moves into.
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	_, err = g.git(ctx, "mv", "--", src, dst)
	return err
}

// Rename renames the file from, in dir, to to with git mv, or as Project does
// when git does not track it.
func (g *Git) Rename(ctx context.Context, dir, from, to string) error {
	full, err := g.path(ctx, dir)
	if err != nil {
		return err
	}
	src, dst := filepath.Join(full, from), filepath.Join(full, to)
	tracked, err := g.tracked(ctx, src)
	if err != nil {
		return err
	}
	if !tracked {
		return g.Project.Rename(ctx, dir, from, to)
	}
	// As Project does: never over another file, but a change of case on a
	// filesystem that ignores case, which git mv makes too.
	srcInfo, err := os.Lstat(src)
	if err != nil {
		return err
	}
	if dstInfo, err := os.Lstat(dst); err == nil {
		if !os.SameFile(srcInfo, dstInfo) {
			return fmt.Errorf("cannot rename %s to %s in %s: %s exists", from, to, dir, to)
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	_, err = g.git(ctx, "mv", "--", src, dst)
	return err
}

// tracked reports whether git tracks the file name, or any file under the
// directory name.
func (g *Git) tracked(ctx context.Context, name string) (bool, error) {
	out, err := g.git(ctx, "ls-files", "-z", "--", name)
	if err != nil {
		return false, err
	}
	return out != "", nil
}

// git runs git in the project's directory, and returns what it printed.  A
// failure is reported with what git said about it.
func (g *Git) git(ctx context.Context, args ...string) (string, error) {
	root, err := g.Root(ctx)
	if err != nil {
		return "", err
	}
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = root
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return "", fmt.Errorf("git %s: %s: %w", args[0], msg, err)
		}
		return "", fmt.Errorf("git %s: %w", args[0], err)
	}
	return stdout.String(), nil
}
