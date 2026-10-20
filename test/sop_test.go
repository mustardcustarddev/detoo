package clitest

import (
	"strings"
	"testing"
)

// The sop tool shares its engine with the other record tools, and their tests
// cover the behaviour every tool shares.  These tests cover what the sop
// tool's definition changes: its directory, the prefix of its file names, the
// README its log opens with, the sections of a new one, the words and options
// for its links, the heading of its reports, its configuration, and its help.

// A log of SOPs opens with a README, not an SOP, so the first SOP is
// number 1.  Running init again leaves the README alone.
func TestInitSopRepository(t *testing.T) {
	w := newWorkspace(t, "sop")
	w.sop("init")
	w.ls("docs/sop")
	w.cat("docs/sop/README.md")
	w.sop("init")
	w.sop("new", "Roll back a release")
	w.sop("list")
	w.check()
}

// A new SOP is a draft: a title, a date, a status and its steps.
func TestCreateSop(t *testing.T) {
	w := newWorkspace(t, "sop")
	w.sop("new", "Roll back a release")
	w.cat("docs/sop/SOP-0001-roll-back-a-release.md")
	w.check()
}

// An SOP that is superseded loses its Draft or Approved status; one that is
// updated keeps it.
func TestSupersedeAndUpdateSops(t *testing.T) {
	w := newWorkspace(t, "sop")
	w.sop("new", "Roll back a release")
	w.sop("new", "Rotate the signing key")
	w.write("docs/sop/SOP-0002-rotate-the-signing-key.md",
		"# 2. Rotate the signing key\n\nDate: 2026-10-16\n\n## Status\n\nApproved\n\n## Steps\n\n1. Rotate it.\n")
	w.sop("new", "-s", "1", "Roll back a release")
	w.sop("new", "--supersede", "2", "Rotate the signing keys")
	w.sop("link", "-u", "4", "3")
	w.head(10, "docs/sop/SOP-0001-roll-back-a-release.md")
	w.head(10, "docs/sop/SOP-0002-rotate-the-signing-key.md")
	w.head(12, "docs/sop/SOP-0003-roll-back-a-release.md")
	w.sop("list")
	w.check()
}

func TestSopReports(t *testing.T) {
	w := newWorkspace(t, "sop")
	w.sop("new", "Roll back a release")
	w.sop("new", "-s", "1", "Roll back a release with the deploy tool")
	w.sop("list", "--toc")
	w.sop("graph")
	w.check()
}

// sop is configured by its own variables.
func TestSopConfiguration(t *testing.T) {
	w := newWorkspace(t, "sop")
	w.sop("config")
	w.sop("init")
	w.sop("config")
	w.write("custom.md", "# NUMBER. TITLE, from custom.md\n\nDate: DATE\n\n## Status\n\nSTATUS\n\n## Steps\n")
	w.setenv("SOP_TEMPLATE", "custom.md")
	w.setenv("SOP_DATE", "2026-10-16")
	w.sop("new", "Uses the sop settings")
	w.cat("docs/sop/SOP-0001-uses-the-sop-settings.md")
	w.check()
}

func TestSopHelp(t *testing.T) {
	w := newWorkspace(t, "sop")
	w.sop("help")
	w.sop("help", "new")
	w.sop("help", "link")
	w.check()
}

// TestSopCompletionScript checks that sop generates its own completion script,
// which "make deploy" installs.
func TestSopCompletionScript(t *testing.T) {
	w := newWorkspace(t, "sop")
	if script := w.output("sop", "completion", "bash"); !strings.Contains(script, "__start_sop") {
		t.Errorf("completion script does not define __start_sop:\n%s", script)
	}
}
