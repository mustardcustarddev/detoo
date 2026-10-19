package project

import (
	"context"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// newRepository returns a project in a new git repository, whose index holds
// the files tracked, and which also holds the files untracked.  Nothing is
// committed: git mv works on the index alone.
func newRepository(t *testing.T, tracked, untracked []string) (*Git, string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	root := resolved(t, t.TempDir())
	for _, name := range append(append([]string{}, tracked...), untracked...) {
		write(t, filepath.Join(root, name), "")
	}
	gitIn(t, root, "init", "--quiet")
	if len(tracked) > 0 {
		gitIn(t, root, append([]string{"add", "--"}, tracked...)...)
	}
	return NewGit(New(root)), root
}

func gitIn(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return string(out)
}

// index lists the files git tracks.
func index(t *testing.T, root string) []string {
	t.Helper()
	return strings.Fields(gitIn(t, root, "ls-files"))
}

func TestRepositoryIsAGitWorkingTree(t *testing.T) {
	g, _ := newRepository(t, nil, []string{"docs/adr/ADR-0001-first.md"})
	if in, err := g.Repository(context.Background()); err != nil || !in {
		t.Errorf("Repository() = %v, %v, want true", in, err)
	}
	dir := resolved(t, t.TempDir())
	// Whatever repository the temporary directory is in, this one is not.
	t.Setenv("GIT_CEILING_DIRECTORIES", filepath.Dir(dir))
	outside := NewGit(New(dir))
	if in, err := outside.Repository(context.Background()); err != nil || in {
		t.Errorf("Repository() outside a repository = %v, %v, want false", in, err)
	}
}

// A directory git tracks files in is moved with git mv, which takes the
// files it does not track along.
func TestGitMovesATrackedDirectory(t *testing.T) {
	g, root := newRepository(t, []string{"doc/adr/0001-first.md"}, []string{"doc/adr/0002-second.md"})
	if err := g.Move(context.Background(), "doc/adr", "docs/adr"); err != nil {
		t.Fatal(err)
	}
	if got := index(t, root); !reflect.DeepEqual(got, []string{"docs/adr/0001-first.md"}) {
		t.Errorf("index = %v, want the tracked file moved", got)
	}
	if !exists(filepath.Join(root, "docs/adr/0002-second.md")) || exists(filepath.Join(root, "doc/adr")) {
		t.Errorf("the untracked file was not moved with the directory")
	}
}

// A directory git tracks nothing in is moved all the same.
func TestGitMovesAnUntrackedDirectory(t *testing.T) {
	g, root := newRepository(t, []string{"README.md"}, []string{"doc/rfc/0001-first.md"})
	if err := g.Move(context.Background(), "doc/rfc", "docs/rfc"); err != nil {
		t.Fatal(err)
	}
	if !exists(filepath.Join(root, "docs/rfc/0001-first.md")) {
		t.Errorf("doc/rfc was not moved to docs/rfc")
	}
}

// A tracked file is renamed with git mv, even when only its case changes, and
// an untracked one is renamed all the same.
func TestGitRenamesFiles(t *testing.T) {
	g, root := newRepository(t,
		[]string{"docs/userstory/story-0001-first.md"},
		[]string{"docs/userstory/0002-second.md"})
	ctx := context.Background()
	if err := g.Rename(ctx, "docs/userstory", "story-0001-first.md", "STORY-0001-first.md"); err != nil {
		t.Fatal(err)
	}
	if err := g.Rename(ctx, "docs/userstory", "0002-second.md", "STORY-0002-second.md"); err != nil {
		t.Fatal(err)
	}
	if got := index(t, root); !reflect.DeepEqual(got, []string{"docs/userstory/STORY-0001-first.md"}) {
		t.Errorf("index = %v, want the tracked file renamed", got)
	}
	names, _, _ := g.Files(ctx, "docs/userstory")
	if want := []string{"STORY-0001-first.md", "STORY-0002-second.md"}; !reflect.DeepEqual(names, want) {
		t.Errorf("Files() = %v, want %v", names, want)
	}
}

// A file git does not track, or has staged but never committed, has no
// history, even before the repository's first commit.
func TestAddedIsNotKnownBeforeACommit(t *testing.T) {
	g, _ := newRepository(t, []string{"docs/adr/ADR-0001-first.md"}, []string{"docs/adr/ADR-0002-second.md"})
	for _, name := range []string{"docs/adr/ADR-0001-first.md", "docs/adr/ADR-0002-second.md", "docs/adr/adr-0001-FIRST.md"} {
		if _, ok, err := g.Added(context.Background(), name); err != nil || ok {
			t.Errorf("Added(%s) = %v, %v, want not added", name, ok, err)
		}
	}
}
