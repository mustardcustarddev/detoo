package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mustardcustard/detoo/internal/app"
	"github.com/mustardcustard/detoo/internal/cli/output"
	"github.com/mustardcustard/detoo/internal/domain/record"
	"github.com/mustardcustard/detoo/internal/features/health/doctor"
)

// Confirm asks a question that is answered yes or no, and returns the answer.
type Confirm func(question string) (bool, error)

// NewDetooCommand returns the command tree of detoo, whose help is in files,
// wired to run against application.  terminal finds out about the terminal a
// report is written to; it is called only by a command that draws one, since
// it may query the terminal.  confirm asks the user a question, and is nil
// when there is no one to ask, as when detoo runs in a script.
func NewDetooCommand(files fs.FS, application *app.Detoo, terminal func() output.Terminal, confirm Confirm, stdout, stderr io.Writer, getenv func(string) string) *cobra.Command {
	// detoo keeps no log, but is named, documented and paged as a tool is.
	t := tool{Kind: record.Kind{Name: "detoo"}, files: files}
	root := &cobra.Command{
		Use:           t.Name,
		Short:         "Look after a project's logs as a whole",
		Long:          t.help("root.md"),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.Help()
			return exitWith(1)
		},
	}
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.AddCommand(newDoctorCommand(t, application, terminal, confirm))
	optionsFirst(root)
	usePagedHelp(root, t, stdout, stderr, getenv)
	return root
}

func newDoctorCommand(t tool, application *app.Detoo, terminal func() output.Terminal, confirm Confirm) *cobra.Command {
	var fix, git bool
	cmd := &cobra.Command{
		Use:     "doctor",
		Short:   "Check that each log is where its tool looks, and named as it names documents",
		Long:    t.help("doctor.md"),
		Example: t.help("doctor.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if git && !fix {
				return errors.New("give --fix to use --git")
			}
			healthy := false
			if fix {
				if !cmd.Flags().Changed("git") {
					var err error
					if git, err = offerGit(cmd, application, confirm); err != nil {
						return err
					}
				}
				fixed, err := application.Health.Fix.Handle(cmd.Context(), doctor.Command{Git: git})
				if err != nil {
					return err
				}
				if err := output.DoctorFix(cmd.OutOrStdout(), terminal(), fixed); err != nil {
					return err
				}
				healthy = fixed.After.Healthy()
			} else {
				result, err := application.Health.Doctor.Handle(cmd.Context(), doctor.Query{})
				if err != nil {
					return err
				}
				if err := output.Doctor(cmd.OutOrStdout(), terminal(), result); err != nil {
					return err
				}
				healthy = result.Healthy()
			}
			// The report has explained the problems; the exit status says
			// there were some.
			if !healthy {
				return exitWith(1)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&fix, "fix", false, "move and rename the logs to where the tools look, and update the links to them")
	cmd.Flags().BoolVar(&git, "git", false, "with --fix, move and rename with git mv; --git=false moves without git, and neither asks")
	return cmd
}

// offerGit asks whether to move with git, when there is someone to ask, the
// fix has something to move or rename, and the project is in a git
// repository.  Otherwise the fix moves without git, as it does when asked not
// to.
func offerGit(cmd *cobra.Command, application *app.Detoo, confirm Confirm) (bool, error) {
	if confirm == nil {
		return false, nil
	}
	result, err := application.Health.Doctor.Handle(cmd.Context(), doctor.Query{})
	if err != nil || result.Fixable() == 0 {
		return false, err
	}
	in, err := application.Health.Fix.Repository(cmd.Context())
	if err != nil || !in {
		return false, err
	}
	return confirm("This project is in a git repository. Move and rename with git mv, so that git records the renames?")
}

// Prompt returns a Confirm that writes each question to w, and reads the
// answer from r.  The answer is yes unless it begins with n; at the end of
// the input, it is no.
func Prompt(r io.Reader, w io.Writer) Confirm {
	answers := bufio.NewReader(r)
	return func(question string) (bool, error) {
		fmt.Fprintf(w, "%s [Y/n] ", question)
		line, err := answers.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return false, err
		}
		if errors.Is(err, io.EOF) && line == "" {
			fmt.Fprintln(w)
			return false, nil
		}
		answer := strings.ToLower(strings.TrimSpace(line))
		return !strings.HasPrefix(answer, "n"), nil
	}
}
