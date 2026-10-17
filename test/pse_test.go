package clitest

import (
	"strings"
	"testing"
)

// The pse tool shares its engine with the other record tools, and their tests
// cover the behaviour every tool shares.  These tests cover what the pse
// tool's definition changes: a directory named for the method rather than the
// tool, the prefix of its file names, its first playbook, the sections of a
// new playbook, the words and options for its links, the heading of its
// reports, its configuration, and its help.

func TestInitPseRepository(t *testing.T) {
	w := newWorkspace(t, "pse")
	w.pse("init")
	w.ls("docs/piosee")
	w.cat("docs/piosee/PIOSEE-0001-plan-for-incidents-with-piosee.md")
	w.check()
}

// A new playbook is a draft, with a section for its scenario and for each step
// of PIOSEE.
func TestCreatePiosee(t *testing.T) {
	w := newWorkspace(t, "pse")
	w.pse("new", "Unexpected charge after deletion")
	w.cat("docs/piosee/PIOSEE-0001-unexpected-charge-after-deletion.md")
	w.check()
}

// A playbook that is replaced loses its Draft status; one that is updated
// keeps it.
func TestReplaceAndUpdatePiosees(t *testing.T) {
	w := newWorkspace(t, "pse")
	w.pse("new", "Checkout errors after a release")
	w.pse("new", "Crossing the VAT threshold")
	w.pse("new", "-r", "1", "Checkout fails")
	w.pse("new", "--update", "2", "Registering for VAT abroad")
	w.pse("link", "-u", "3", "2")
	w.head(10, "docs/piosee/PIOSEE-0001-checkout-errors-after-a-release.md")
	w.head(12, "docs/piosee/PIOSEE-0002-crossing-the-vat-threshold.md")
	w.head(12, "docs/piosee/PIOSEE-0003-checkout-fails.md")
	w.head(10, "docs/piosee/PIOSEE-0004-registering-for-vat-abroad.md")
	w.pse("list")
	w.check()
}

// pse names its options after its own relations.  Nothing is written when an
// option is wrong.
func TestPseLinkOptions(t *testing.T) {
	w := newWorkspace(t, "pse")
	w.pse("new", "First Playbook")
	w.pse("new", "Second Playbook")
	w.pse("link", "2", "1")
	w.pse("link", "-s", "2", "1")
	w.pse("new", "--reverse", "1", "Third Playbook")
	w.ls("docs/piosee")
	w.head(8, "docs/piosee/PIOSEE-0001-first-playbook.md")
	w.check()
}

func TestPseReports(t *testing.T) {
	w := newWorkspace(t, "pse")
	w.pse("new", "First Playbook")
	w.pse("new", "-r", "1", "Second Playbook")
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
	w.cat("decisions/PIOSEE-0002-ignores-the-sty-settings.md")
	w.cat("decisions/PIOSEE-0003-uses-the-pse-settings.md")
	w.check()
}

// Two tools whose records are prefixed can share a directory, and neither
// sees the other's records.
func TestPrefixedLogsCanShareADirectory(t *testing.T) {
	w := newWorkspace(t, "pse")
	w.pse("init", "docs/work")
	w.sty("init", "docs/work")
	w.pse("new", "A Playbook")
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
