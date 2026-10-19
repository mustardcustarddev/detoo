package cli

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/mustardcustard/detoo/internal/adapters/project"
	"github.com/mustardcustard/detoo/internal/app"
	"github.com/mustardcustard/detoo/internal/cli/output"
	"github.com/mustardcustard/detoo/internal/detoo"
	"github.com/mustardcustard/detoo/internal/features/agents/skills"
)

// --git and --renumber say how to fix, so they are no use without --fix.  No
// handler is reached, so a zero Detoo is enough.
func TestGitAndRenumberNeedFix(t *testing.T) {
	for _, flag := range []string{"--git", "--renumber"} {
		root := NewDetooCommand(detoo.Files, &app.Detoo{}, nil, nil, io.Discard, io.Discard, func(string) string { return "" })
		root.SetArgs([]string{"doctor", flag})
		if err := root.Execute(); err == nil || !strings.Contains(err.Error(), "--fix") {
			t.Errorf("doctor %s = %v, want an error asking for --fix", flag, err)
		}
	}
}

// A fix offers to renumber only when documents share a number in a git
// repository, whose history says which came first, and to move with git only
// when it has something to move or rename there.  Renumbering, once declined,
// is nothing to move.  Before asking, it says what it will do anyway.  These
// run the real handlers against a temporary repository, since the questions
// depend on what they find.
func TestFixOffersOnlyWhatItCanDo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	renumber, git := "Renumber documents that share a number?", output.GitQuestion.Title
	shared := []string{"docs/adr/ADR-0001-first.md", "docs/adr/ADR-0001-second.md"}
	for _, test := range []struct {
		name       string
		files      []string
		repository bool
		yes        bool
		want       []string
		plan       string
	}{
		{"nothing to move", []string{"docs/adr/ADR-0001-first.md"}, true, true, nil, ""},
		{"a log to move", []string{"doc/adr/0001-first.md"}, true, true, []string{git}, "Move doc/adr/ → docs/adr/"},
		{"not a repository", []string{"doc/adr/0001-first.md"}, false, true, nil, ""},
		{"a shared number", shared, true, true, []string{renumber, git}, ""},
		{"a shared number not renumbered", shared, true, false, []string{renumber}, ""},
		{"a shared number outside a repository", shared, false, true, nil, ""},
		{"a shared number and a rename", append([]string{"docs/adr/ADX-0002-third.md"}, shared...), true, false,
			[]string{renumber, git}, "Rename docs/adr/ADX-0002-third.md → ADR-0002-third.md"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			t.Setenv("GIT_CEILING_DIRECTORIES", filepath.Dir(root))
			for _, file := range test.files {
				if err := os.MkdirAll(filepath.Join(root, filepath.Dir(file)), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(root, file), nil, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			if test.repository {
				if out, err := exec.Command("git", "init", "--quiet", root).CombinedOutput(); err != nil {
					t.Fatalf("git init: %v\n%s", err, out)
				}
			}
			p := project.New(root)
			application := app.NewDetoo(detoo.Kinds, p, project.NewGit(p), nil, skills.Installations{})
			var asked []string
			confirm := func(q output.Question) (bool, error) {
				asked = append(asked, q.Title)
				return test.yes, nil
			}
			cmd := &cobra.Command{}
			cmd.SetContext(context.Background())
			var said strings.Builder
			cmd.SetErr(&said)
			terminal := func() output.Terminal { return output.Terminal{} }
			var git, renumber bool
			if err := offer(cmd, application, terminal, confirm, &git, &renumber); err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(asked, test.want) {
				t.Errorf("asked %q, want %q", asked, test.want)
			}
			if test.plan == "" && said.Len() > 0 || !strings.Contains(said.String(), test.plan) {
				t.Errorf("said %q, want the plan %q", said.String(), test.plan)
			}
		})
	}
}
