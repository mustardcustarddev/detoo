package cli

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/mustardcustard/detoo/internal/app"
	"github.com/mustardcustard/detoo/internal/cli/output"
	"github.com/mustardcustard/detoo/internal/features/agents/skills"
)

// newSkillsCommand returns detoo skills, whose commands list the agent skills
// that drive the tools, install them, update them and remove them.  Each says
// what it is doing as it does it, on stderr; only the list of skills a source
// offers is written to stdout, for a script to read.
func newSkillsCommand(t tool, application *app.Detoo, terminal func() output.Terminal) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "skills",
		Short:   "List, install, update and remove the agent skills that drive the tools",
		Long:    t.help("skills.md"),
		Example: t.help("skills.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.Help()
			return exitWith(1)
		},
	}
	cmd.AddCommand(
		newSkillsListCommand(t, application, terminal),
		newSkillsInstallCommand(t, application, terminal),
		newSkillsUpdateCommand(t, application, terminal),
		newSkillsRemoveCommand(t, application, terminal),
	)
	return cmd
}

func newSkillsListCommand(t tool, application *app.Detoo, terminal func() output.Terminal) *cobra.Command {
	var source string
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List the agent skills a source offers, by default those detoo offers",
		Long:    t.help("skills-list.md"),
		Example: t.help("skills-list.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			src, err := skills.ParseSource(source)
			if err != nil {
				return err
			}
			term := terminal()
			offered, err := application.Skills.Install.Available(cmd.Context(), src, output.NewSkillsProgress(cmd.ErrOrStderr(), term))
			if err != nil {
				return err
			}
			return output.SkillsAvailable(cmd.OutOrStdout(), term, offered)
		},
	}
	cmd.Flags().StringVar(&source, "source", skills.DefaultSource,
		"list the skills of `SOURCE`: a repository on GitHub, as OWNER/REPO or OWNER/REPO#REF, or a directory, as a path beginning with ., / or ~")
	return cmd
}

func newSkillsInstallCommand(t tool, application *app.Detoo, terminal func() output.Terminal) *cobra.Command {
	var source string
	var agents []string
	var global bool
	cmd := &cobra.Command{
		Use:     "install [SKILL...]",
		Short:   "Install agent skills, by default every skill detoo offers",
		Long:    t.help("skills-install.md"),
		Example: t.help("skills-install.example"),
		RunE: func(cmd *cobra.Command, names []string) error {
			src, err := skills.ParseSource(source)
			if err != nil {
				return err
			}
			progress := output.NewSkillsProgress(cmd.ErrOrStderr(), terminal())
			installed, err := application.Skills.Install.Handle(cmd.Context(), skills.Install{
				Source: src,
				Skills: names,
				Agents: agents,
				Global: global,
			}, progress)
			if err != nil {
				return err
			}
			progress.SkillsInstalled(installed, global)
			return nil
		},
	}
	cmd.Flags().StringVar(&source, "source", skills.DefaultSource,
		"install from `SOURCE`: a repository on GitHub, as OWNER/REPO or OWNER/REPO#REF, or a directory, as a path beginning with ., / or ~")
	cmd.Flags().StringArrayVarP(&agents, "agent", "a", nil,
		"install for `AGENT`, such as claude-code or codex, rather than for Claude Code and each agent found; repeat for more")
	cmd.Flags().BoolVarP(&global, "global", "g", false, "install for your user, in every project, rather than in this one")
	return cmd
}

func newSkillsUpdateCommand(t tool, application *app.Detoo, terminal func() output.Terminal) *cobra.Command {
	var global bool
	cmd := &cobra.Command{
		Use:     "update [SKILL...]",
		Short:   "Update installed detoo skills from detoo's repository, by default every one",
		Long:    t.help("skills-update.md"),
		Example: t.help("skills-update.example"),
		RunE: func(cmd *cobra.Command, names []string) error {
			progress := output.NewSkillsProgress(cmd.ErrOrStderr(), terminal())
			progress.Updating = true
			updated, err := application.Skills.Update.Handle(cmd.Context(), skills.Update{Skills: names, Global: global}, progress)
			if err != nil {
				return err
			}
			if !progress.SkillsUpdated(updated) {
				return exitWith(1)
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&global, "global", "g", false, "update the skills installed for your user, rather than in this project")
	return cmd
}

func newSkillsRemoveCommand(t tool, application *app.Detoo, terminal func() output.Terminal) *cobra.Command {
	var global, all bool
	cmd := &cobra.Command{
		Use:     "remove SKILL...",
		Short:   "Remove installed agent skills, or with --all, every detoo skill",
		Long:    t.help("skills-remove.md"),
		Example: t.help("skills-remove.example"),
		RunE: func(cmd *cobra.Command, names []string) error {
			switch {
			case all && len(names) > 0:
				return errors.New("give --all, or name the skills to remove, not both")
			case !all && len(names) == 0:
				return errors.New("name the skills to remove, or give --all")
			}
			progress := output.NewSkillsProgress(cmd.ErrOrStderr(), terminal())
			removed, err := application.Skills.Remove.Handle(cmd.Context(), skills.Remove{Skills: names, All: all, Global: global}, progress)
			if err != nil {
				return err
			}
			progress.SkillsRemoved(removed)
			return nil
		},
	}
	cmd.Flags().BoolVar(&all, "all", false, "remove every detoo skill the lock file records")
	cmd.Flags().BoolVarP(&global, "global", "g", false, "remove skills installed for your user, rather than in this project")
	return cmd
}
