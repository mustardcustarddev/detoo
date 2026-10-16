package cli

import (
	"io"
	"io/fs"

	"github.com/spf13/cobra"

	"github.com/mustardcustard/detoo/internal/app"
	"github.com/mustardcustard/detoo/internal/cli/output"
	"github.com/mustardcustard/detoo/internal/domain/record"
	"github.com/mustardcustard/detoo/internal/features/health/doctor"
)

// NewDetooCommand returns the command tree of detoo, whose help is in files,
// wired to run against application.  terminal finds out about the terminal a
// report is written to; it is called only by a command that draws one, since
// it may query the terminal.
func NewDetooCommand(files fs.FS, application *app.Detoo, terminal func() output.Terminal, stdout, stderr io.Writer, getenv func(string) string) *cobra.Command {
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
	root.AddCommand(newDoctorCommand(t, application, terminal))
	optionsFirst(root)
	usePagedHelp(root, t, stdout, stderr, getenv)
	return root
}

func newDoctorCommand(t tool, application *app.Detoo, terminal func() output.Terminal) *cobra.Command {
	return &cobra.Command{
		Use:     "doctor",
		Short:   "Check that each log is where its tool looks, and named as it names documents",
		Long:    t.help("doctor.md"),
		Example: t.help("doctor.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := application.Health.Doctor.Handle(cmd.Context(), doctor.Query{})
			if err != nil {
				return err
			}
			if err := output.Doctor(cmd.OutOrStdout(), terminal(), result); err != nil {
				return err
			}
			// The report has explained the problems; the exit status says
			// there were some.
			if !result.Healthy() {
				return exitWith(1)
			}
			return nil
		},
	}
}
