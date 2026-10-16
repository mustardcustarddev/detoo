package clitest

import (
	"strings"
	"testing"
)

// The pse tool shares its engine with the other record tools, and their tests
// cover the behaviour every tool shares.  These tests cover what the pse
// tool's definition changes: a directory named for the method rather than the
// tool, the prefix of its file names, its first record, the sections of a new
// record, the words and options for its links, the heading of its reports, its
// configuration, and its help.

func TestInitPseRepository(t *testing.T) {
	w := newWorkspace(t, "pse")
	w.pse("init")
	w.ls("docs/piosee")
	w.cat("docs/piosee/piosee-0001-make-decisions-with-piosee.md")
	w.check()
}

// A new record is a draft, in a section for each step of PIOSEE.
func TestCreatePiosee(t *testing.T) {
	w := newWorkspace(t, "pse")
	w.pse("new", "Roll back release 2.3")
	w.cat("docs/piosee/piosee-0001-roll-back-release-2-3.md")
	w.check()
}

// A decision that is reversed loses its Draft status; one that is updated
// keeps it.
func TestReverseAndUpdatePiosees(t *testing.T) {
	w := newWorkspace(t, "pse")
	w.pse("new", "Turn off tax codes")
	w.pse("new", "Move the launch to March")
	w.pse("new", "-r", "1", "Roll back release 2.3")
	w.pse("new", "--update", "2", "Move the launch to April")
	w.pse("link", "-u", "3", "2")
	w.head(10, "docs/piosee/piosee-0001-turn-off-tax-codes.md")
	w.head(12, "docs/piosee/piosee-0002-move-the-launch-to-march.md")
	w.head(12, "docs/piosee/piosee-0003-roll-back-release-2-3.md")
	w.head(10, "docs/piosee/piosee-0004-move-the-launch-to-april.md")
	w.pse("list")
	w.check()
}

// pse names its options after its own relations.  Nothing is written when an
// option is wrong.
func TestPseLinkOptions(t *testing.T) {
	w := newWorkspace(t, "pse")
	w.pse("new", "First Decision")
	w.pse("new", "Second Decision")
	w.pse("link", "2", "1")
	w.pse("link", "-s", "2", "1")
	w.pse("new", "--replace", "1", "Third Decision")
	w.ls("docs/piosee")
	w.head(8, "docs/piosee/piosee-0001-first-decision.md")
	w.check()
}

func TestPseReports(t *testing.T) {
	w := newWorkspace(t, "pse")
	w.pse("new", "First Decision")
	w.pse("new", "-r", "1", "Second Decision")
	w.pse("list", "--toc")
	w.pse("graph")
	w.pse("graph", "--graphviz")
	w.check()
}

// pse is configured by its own variables and marker file, and not by those of
// another tool.
func TestPseConfiguration(t *testing.T) {
	w := newWorkspace(t, "pse")
	w.pse("config")
	w.pse("init", "decisions")
	w.pse("config")
	w.cat(".pse-dir")
	w.write("custom.md", "# NUMBER. TITLE, from custom.md\n\nDate: DATE\n\n## Status\n\nSTATUS\n\n## Problem\n")
	w.setenv("STY_TEMPLATE", "custom.md")
	w.setenv("STY_DATE", "2000-01-01")
	w.pse("new", "Ignores the sty settings")
	w.setenv("PSE_TEMPLATE", "custom.md")
	w.setenv("PSE_DATE", "2026-10-16")
	w.pse("new", "Uses the pse settings")
	w.cat("decisions/piosee-0002-ignores-the-sty-settings.md")
	w.cat("decisions/piosee-0003-uses-the-pse-settings.md")
	w.check()
}

// Two tools whose records are prefixed can share a directory, and neither
// sees the other's records.
func TestPrefixedLogsCanShareADirectory(t *testing.T) {
	w := newWorkspace(t, "pse")
	w.pse("init", "docs/work")
	w.sty("init", "docs/work")
	w.pse("new", "A Decision")
	w.sty("new", "A Story")
	w.ls("docs/work")
	w.pse("list")
	w.sty("list")
	w.pse("link", "-u", "2", "a-story")
	w.check()
}

func TestPseHelp(t *testing.T) {
	w := newWorkspace(t, "pse")
	w.pse("help")
	w.pse("help", "new")
	w.pse("help", "link")
	w.pse("nonsense")
	w.check()
}

// TestPseCompletionScript checks that pse generates its own completion script,
// which "make deploy" installs.
func TestPseCompletionScript(t *testing.T) {
	w := newWorkspace(t, "pse")
	if script := w.output("pse", "completion", "bash"); !strings.Contains(script, "__start_pse") {
		t.Errorf("completion script does not define __start_pse:\n%s", script)
	}
}
