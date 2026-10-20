package clitest

import (
	"strings"
	"testing"
)

// The fyr tool shares its engine with the other record tools, and their tests
// cover the behaviour every tool shares.  These tests cover what the fyr
// tool's definition changes: its directory, the prefix of its file names, the
// README its log opens with, the sections of a new one and the directory beside it,
// the words and options for its links, the heading of its reports, its
// configuration, and its help.

// A log of fire drills opens with a README, not a drill, so the first drill
// is number 1.  Running init again leaves the README alone.
func TestInitFyrRepository(t *testing.T) {
	w := newWorkspace(t, "fyr")
	w.fyr("init")
	w.ls("docs/firedrill")
	w.cat("docs/firedrill/README.md")
	w.fyr("init")
	w.fyr("new", "Customer charged twice")
	w.fyr("list")
	w.check()
}

// A new fire drill is a draft: a title, a date, a status, and a section for
// each part of the drill.  Beside it is a directory named as it is, holding a
// problem directory and a solution directory, each with a README that says
// what goes in it.
func TestCreateFireDrill(t *testing.T) {
	w := newWorkspace(t, "fyr")
	w.fyr("new", "Customer charged twice")
	w.ls("docs/firedrill")
	w.cat("docs/firedrill/FIREDRILL-0001-customer-charged-twice.md")
	w.ls("docs/firedrill/FIREDRILL-0001-customer-charged-twice")
	w.cat("docs/firedrill/FIREDRILL-0001-customer-charged-twice/problem/README.md")
	w.cat("docs/firedrill/FIREDRILL-0001-customer-charged-twice/solution/README.md")
	w.fyr("list")
	w.check()
}

// A fire drill that is replaced loses its Draft status; one that is updated
// keeps it.
func TestReplaceAndUpdateFireDrills(t *testing.T) {
	w := newWorkspace(t, "fyr")
	w.fyr("new", "Customer charged twice")
	w.fyr("new", "Release breaks checkout")
	w.fyr("new", "-r", "1", "Customer charged twice")
	w.fyr("link", "-u", "3", "2")
	w.head(10, "docs/firedrill/FIREDRILL-0001-customer-charged-twice.md")
	w.head(10, "docs/firedrill/FIREDRILL-0002-release-breaks-checkout.md")
	w.head(12, "docs/firedrill/FIREDRILL-0003-customer-charged-twice.md")
	w.fyr("list")
	w.check()
}

func TestFyrReports(t *testing.T) {
	w := newWorkspace(t, "fyr")
	w.fyr("new", "Customer charged twice")
	w.fyr("new", "-r", "1", "Customer charged twice")
	w.fyr("list", "--toc")
	w.fyr("graph")
	w.check()
}

// fyr is configured by its own variables.
func TestFyrConfiguration(t *testing.T) {
	w := newWorkspace(t, "fyr")
	w.fyr("config")
	w.fyr("init")
	w.fyr("config")
	w.write("custom.md", "# NUMBER. TITLE, from custom.md\n\nDate: DATE\n\n## Status\n\nSTATUS\n\n## Playbook\n")
	w.setenv("FYR_TEMPLATE", "custom.md")
	w.setenv("FYR_DATE", "2026-10-19")
	w.fyr("new", "Uses the fyr settings")
	w.cat("docs/firedrill/FIREDRILL-0001-uses-the-fyr-settings.md")
	w.ls("docs/firedrill/FIREDRILL-0001-uses-the-fyr-settings")
	w.check()
}

func TestFyrHelp(t *testing.T) {
	w := newWorkspace(t, "fyr")
	w.fyr("help")
	w.fyr("help", "new")
	w.fyr("help", "link")
	w.check()
}

// TestFyrCompletionScript checks that fyr generates its own completion script,
// which "make deploy" installs.
func TestFyrCompletionScript(t *testing.T) {
	w := newWorkspace(t, "fyr")
	if script := w.output("fyr", "completion", "bash"); !strings.Contains(script, "__start_fyr") {
		t.Errorf("completion script does not define __start_fyr:\n%s", script)
	}
}
