package cli

import (
	"errors"
	"io"
	"io/fs"

	"github.com/spf13/cobra"

	"github.com/mustardcustard/detoo/internal/app"
	"github.com/mustardcustard/detoo/internal/cli/output"
	"github.com/mustardcustard/detoo/internal/domain/record"
	"github.com/mustardcustard/detoo/internal/features/health/doctor"
)

// Confirm asks a question that is answered yes or no, and returns the answer.
type Confirm func(question output.Question) (bool, error)

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
		Short:         "Look after a project's logs as a whole, and the agent skills that drive the tools",
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
	root.AddCommand(newDoctorCommand(t, application, terminal, confirm), newSkillsCommand(t, application, terminal))
	optionsFirst(root)
	usePagedHelp(root, t, stdout, stderr, getenv)
	return root
}

func newDoctorCommand(t tool, application *app.Detoo, terminal func() output.Terminal, confirm Confirm) *cobra.Command {
	var fix, git, renumber bool
	cmd := &cobra.Command{
		Use:     "doctor",
		Short:   "Check that each log is where its tool looks, with documents named and numbered as it names them",
		Long:    t.help("doctor.md"),
		Example: t.help("doctor.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if git && !fix {
				return errors.New("give --fix to use --git")
			}
			if renumber && !fix {
				return errors.New("give --fix to use --renumber")
			}
			healthy := false
			if fix {
				if err := offer(cmd, application, terminal, confirm, &git, &renumber); err != nil {
					return err
				}
				fixed, err := application.Health.Fix.Handle(cmd.Context(), doctor.Command{Git: git, Renumber: renumber})
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
	cmd.Flags().BoolVar(&renumber, "renumber", false, "with --fix, renumber documents that share a number, by when git added them; --renumber=false does not, and neither asks")
	return cmd
}

// offer asks what a fix is to do that was not given as an option, when there
// is someone to ask: whether to renumber documents that share a number, and
// whether to move with git.  Otherwise the fix does neither, as it does when
// asked not to.  Before the first question, it says what the fix will do
// anyway, so that the questions are not taken for all of it.
func offer(cmd *cobra.Command, application *app.Detoo, terminal func() output.Terminal, confirm Confirm, git, renumber *bool) error {
	askGit, askRenumber := !cmd.Flags().Changed("git"), !cmd.Flags().Changed("renumber")
	if confirm == nil || !askGit && !askRenumber {
		return nil
	}
	result, err := application.Health.Doctor.Handle(cmd.Context(), doctor.Query{})
	if err != nil {
		return err
	}
	askRenumber = askRenumber && result.Renumberable() > 0
	if askGit {
		in, err := application.Health.Fix.Repository(cmd.Context())
		if err != nil {
			return err
		}
		askGit = in
	}
	plan := output.PlanItems(result)
	if len(plan) > 0 && (askRenumber || askGit) {
		if err := output.Plan(cmd.ErrOrStderr(), terminal(), result); err != nil {
			return err
		}
	}

	if askRenumber {
		// Renumbering can break what refers to the documents, so it is
		// done only when the user says yes.
		if *renumber, err = confirm(output.RenumberQuestion(result)); err != nil {
			return err
		}
	}
	// git mv is no use when there is nothing to move or rename.
	if askGit && (len(plan) > 0 || *renumber) {
		*git, err = confirm(output.GitQuestion)
	}
	return err
}
