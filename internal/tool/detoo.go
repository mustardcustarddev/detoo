package tool

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strconv"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/term"

	"github.com/mustardcustard/detoo/internal/adapters/project"
	"github.com/mustardcustard/detoo/internal/app"
	"github.com/mustardcustard/detoo/internal/cli"
	"github.com/mustardcustard/detoo/internal/cli/output"
	"github.com/mustardcustard/detoo/internal/domain/record"
)

// RunDetoo runs detoo, which looks across the logs of kinds and whose help is
// in files, with the process's arguments, and returns its exit status.
func RunDetoo(kinds []record.Kind, files fs.FS) int {
	workdir, err := os.Getwd()
	if err != nil {
		return failDetoo(err)
	}

	// The composition root, as in Run: the disk and the terminal are chosen
	// here, once.
	p := project.New(workdir, kinds)
	application := app.NewDetoo(kinds, p, project.NewGit(p))
	// A question is asked only of someone at a terminal, never of a script.
	var confirm cli.Confirm
	if term.IsTerminal(os.Stdin.Fd()) && term.IsTerminal(os.Stdout.Fd()) {
		confirm = func(q output.Question) (bool, error) {
			return output.Ask(os.Stdin, os.Stderr, terminal(), q)
		}
	}
	root := cli.NewDetooCommand(files, application, terminal, confirm, os.Stdout, os.Stderr, os.Getenv)
	if err := root.Execute(); err != nil {
		if exit, ok := errors.AsType[*cli.ExitError](err); ok {
			return exit.Code
		}
		return failDetoo(err)
	}
	return 0
}

func failDetoo(err error) int {
	fmt.Fprintf(os.Stderr, "detoo: %v\n", err)
	return 1
}

// terminal finds out about the terminal on stdout: how wide it is, whether its
// background is dark, and what colours it shows.  Output that is not to a
// terminal is drawn without colour, $COLUMNS wide, or else 80.
func terminal() output.Terminal {
	t := output.Terminal{
		Dark:    true,
		Profile: colorprofile.Detect(os.Stdout, os.Environ()),
		Version: version(),
	}
	// By its real path, as the project's directory is reported.
	if home, err := os.UserHomeDir(); err == nil {
		t.Home, _ = filepath.EvalSymlinks(home)
	}
	t.Width, _ = strconv.Atoi(os.Getenv("COLUMNS"))
	if term.IsTerminal(os.Stdout.Fd()) {
		if width, _, err := term.GetSize(os.Stdout.Fd()); err == nil {
			t.Width = width
		}
		if t.Profile > colorprofile.ASCII {
			t.Dark = lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
		}
	}
	return t
}

// version is the version of the release detoo was built from, such as
// "1.0.0", or else the commit it was built from, or "dev" when that is not
// known.
func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}
	return versionOf(info)
}

// release matches the version Go records for a build of a commit that is
// tagged as a release, such as v1.0.0, with nothing changed.  A build of any
// other commit has a pseudo-version, such as
// v0.0.0-20261018151304-eb82bb5ed344, and one with changes ends in +dirty.
var release = regexp.MustCompile(`^v([0-9]+\.[0-9]+\.[0-9]+)$`)

func versionOf(info *debug.BuildInfo) string {
	if m := release.FindStringSubmatch(info.Main.Version); m != nil {
		return m[1]
	}
	var revision string
	var modified bool
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.modified":
			modified = setting.Value == "true"
		}
	}
	switch {
	case revision != "" && modified:
		return revision[:min(7, len(revision))] + "-dirty"
	case revision != "":
		return revision[:min(7, len(revision))]
	case info.Main.Version != "" && info.Main.Version != "(devel)":
		return info.Main.Version
	}
	return "dev"
}
