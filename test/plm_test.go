package clitest

import (
	"strings"
	"testing"
)

// The plm tool shares its engine with the other record tools, and their tests
// cover the behaviour every tool shares.  These tests cover what the plm
// tool's definition changes: a directory named for problems rather than the
// tool, the prefix of its file names, its first playbook, the sections of a
// new playbook, the words and options for its links, the heading of its
// reports, its configuration, its help, and the playbooks named before it
// replaced pse.

func TestInitPlmRepository(t *testing.T) {
	w := newWorkspace(t, "plm")
	w.plm("init")
	w.ls("docs/problem")
	w.cat("docs/problem/PROBLEM-0001-plan-for-incidents-with-piosee.md")
	w.check()
}

// A new playbook is a draft, with a section for its scenario and for each step
// of PIOSEE, with the immediate action and the counter-indicated actions after
// the problem.
func TestCreatePlaybook(t *testing.T) {
	w := newWorkspace(t, "plm")
	w.plm("new", "Unexpected charge after deletion")
	w.cat("docs/problem/PROBLEM-0001-unexpected-charge-after-deletion.md")
	w.check()
}

// A playbook that is replaced loses its Draft status; one that is updated
// keeps it.
func TestReplaceAndUpdatePlaybooks(t *testing.T) {
	w := newWorkspace(t, "plm")
	w.plm("new", "Checkout errors after a release")
	w.plm("new", "Crossing the VAT threshold")
	w.plm("new", "-r", "1", "Checkout fails")
	w.plm("new", "--update", "2", "Registering for VAT abroad")
	w.plm("link", "-u", "3", "2")
	w.head(10, "docs/problem/PROBLEM-0001-checkout-errors-after-a-release.md")
	w.head(12, "docs/problem/PROBLEM-0002-crossing-the-vat-threshold.md")
	w.head(12, "docs/problem/PROBLEM-0003-checkout-fails.md")
	w.head(10, "docs/problem/PROBLEM-0004-registering-for-vat-abroad.md")
	w.plm("list")
	w.check()
}

// Playbooks named PIOSEE-, as pse named them, are still playbooks: plm lists
// them, numbers after them, and links to them, until detoo doctor renames
// them.
func TestPlmReadsPlaybooksNamedByPse(t *testing.T) {
	w := newWorkspace(t, "plm")
	w.write("docs/problem/PIOSEE-0001-checkout-fails.md",
		"# 1. Checkout fails\n\nDate: 2026-10-16\n\n## Status\n\nDraft\n\n## Scenario\n\nCheckout fails.\n")
	w.plm("list")
	w.plm("new", "-u", "1", "Checkout fails for one region")
	w.cat("docs/problem/PIOSEE-0001-checkout-fails.md")
	w.check()
}

// plm names its options after its own relations.  Nothing is written when an
// option is wrong.
func TestPlmLinkOptions(t *testing.T) {
	w := newWorkspace(t, "plm")
	w.plm("new", "First Playbook")
	w.plm("new", "Second Playbook")
	w.plm("link", "2", "1")
	w.plm("link", "-s", "2", "1")
	w.plm("new", "--reverse", "1", "Third Playbook")
	w.ls("docs/problem")
	w.head(8, "docs/problem/PROBLEM-0001-first-playbook.md")
	w.check()
}

func TestPlmReports(t *testing.T) {
	w := newWorkspace(t, "plm")
	w.plm("new", "First Playbook")
	w.plm("new", "-r", "1", "Second Playbook")
	w.plm("list", "--toc")
	w.plm("graph")
	w.plm("graph", "--graphviz")
	w.check()
}

// plm is configured by its own variables and marker file, and not by those of
// another tool.
func TestPlmConfiguration(t *testing.T) {
	w := newWorkspace(t, "plm")
	w.plm("config")
	w.plm("init", "decisions")
	w.plm("config")
	w.cat(".plm-dir")
	w.write("custom.md", "# NUMBER. TITLE, from custom.md\n\nDate: DATE\n\n## Status\n\nSTATUS\n\n## Problem\n")
	w.setenv("STY_TEMPLATE", "custom.md")
	w.setenv("STY_DATE", "2000-01-01")
	w.plm("new", "Ignores the sty settings")
	w.setenv("PLM_TEMPLATE", "custom.md")
	w.setenv("PLM_DATE", "2026-10-16")
	w.plm("new", "Uses the plm settings")
	w.cat("decisions/PROBLEM-0002-ignores-the-sty-settings.md")
	w.cat("decisions/PROBLEM-0003-uses-the-plm-settings.md")
	w.check()
}

// Two tools whose records are prefixed can share a directory, and neither
// sees the other's records.
func TestPrefixedLogsCanShareADirectory(t *testing.T) {
	w := newWorkspace(t, "plm")
	w.plm("init", "docs/work")
	w.sty("init", "docs/work")
	w.plm("new", "A Playbook")
	w.sty("new", "A Story")
	w.ls("docs/work")
	w.plm("list")
	w.sty("list")
	w.plm("link", "-u", "2", "a-story")
	w.check()
}

func TestPlmHelp(t *testing.T) {
	w := newWorkspace(t, "plm")
	w.plm("help")
	w.plm("help", "new")
	w.plm("help", "link")
	w.plm("nonsense")
	w.check()
}

// TestPlmCompletionScript checks that plm generates its own completion script,
// which "make deploy" installs.
func TestPlmCompletionScript(t *testing.T) {
	w := newWorkspace(t, "plm")
	if script := w.output("plm", "completion", "bash"); !strings.Contains(script, "__start_plm") {
		t.Errorf("completion script does not define __start_plm:\n%s", script)
	}
}
