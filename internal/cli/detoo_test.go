package cli

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/mustardcustard/detoo/internal/adapters/project"
	"github.com/mustardcustard/detoo/internal/app"
	"github.com/mustardcustard/detoo/internal/detoo"
)

// The answer to a question is yes unless it begins with n, and no when there
// is no answer at all.
func TestPromptReadsYesOrNo(t *testing.T) {
	for _, test := range []struct {
		input string
		want  bool
	}{
		{"\n", true},
		{"y\n", true},
		{"Yes\n", true},
		{"n\n", false},
		{"  No \n", false},
		{"", false},
	} {
		var asked strings.Builder
		got, err := Prompt(strings.NewReader(test.input), &asked)("Move with git?")
		if err != nil || got != test.want {
			t.Errorf("answer %q = %v, %v, want %v", test.input, got, err, test.want)
		}
		if !strings.HasPrefix(asked.String(), "Move with git? [Y/n] ") {
			t.Errorf("the question was written as %q", asked.String())
		}
	}
}

// --git says how to fix, so it is no use without --fix.  No handler is
// reached, so a zero Detoo is enough.
func TestGitNeedsFix(t *testing.T) {
	root := NewDetooCommand(detoo.Files, &app.Detoo{}, nil, nil, io.Discard, io.Discard, func(string) string { return "" })
	root.SetArgs([]string{"doctor", "--git"})
	if err := root.Execute(); err == nil || !strings.Contains(err.Error(), "--fix") {
		t.Errorf("doctor --git = %v, want an error asking for --fix", err)
	}
}

// A fix offers to move with git only when it has something to move or
// rename, in a git repository.  These run the real handlers against a
// temporary repository, since the question depends on what they find.
func TestFixOffersGitOnlyWhenThereIsSomethingToMove(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	for _, test := range []struct {
		name, file string
		repository bool
		want       bool
	}{
		{"nothing to move", "docs/adr/ADR-0001-first.md", true, false},
		{"a log to move", "doc/adr/0001-first.md", true, true},
		{"not a repository", "doc/adr/0001-first.md", false, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			t.Setenv("GIT_CEILING_DIRECTORIES", filepath.Dir(root))
			if err := os.MkdirAll(filepath.Join(root, filepath.Dir(test.file)), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(root, test.file), nil, 0o644); err != nil {
				t.Fatal(err)
			}
			if test.repository {
				if out, err := exec.Command("git", "init", "--quiet", root).CombinedOutput(); err != nil {
					t.Fatalf("git init: %v\n%s", err, out)
				}
			}
			p := project.New(root, detoo.Kinds)
			application := app.NewDetoo(detoo.Kinds, p, project.NewGit(p))
			asked := false
			confirm := func(string) (bool, error) {
				asked = true
				return false, nil
			}
			cmd := &cobra.Command{}
			cmd.SetContext(context.Background())
			if _, err := offerGit(cmd, application, confirm); err != nil {
				t.Fatal(err)
			}
			if asked != test.want {
				t.Errorf("asked = %v, want %v", asked, test.want)
			}
		})
	}
}
