package clitest

import (
	"strings"
	"testing"
)

// The rfc tool shares its engine with adr, whose tests cover numbering, slugs,
// finding the log, project templates and the editor.  These tests cover what
// the rfc tool's definition changes: its directory and first record, the status
// and sections of a new RFC, the words and options for its links, its reports,
// its configuration, and its help.

func TestInitRfcRepository(t *testing.T) {
	w := newWorkspace(t, "rfc")
	w.rfc("init")
	w.ls("docs/rfc")
	w.cat("docs/rfc/RFC-0001-use-requests-for-comments.md")
	w.check()
}

// A new RFC is a draft, in the sections of the RFC template.
func TestCreateRfc(t *testing.T) {
	w := newWorkspace(t, "rfc")
	w.rfc("new", "Adopt gRPC for internal services")
	w.cat("docs/rfc/RFC-0001-adopt-grpc-for-internal-services.md")
	w.check()
}

// Obsoleting an RFC removes its Draft status; updating one leaves it alone.
func TestRfcLinking(t *testing.T) {
	w := newWorkspace(t, "rfc")
	w.rfc("new", "First Proposal")
	w.rfc("new", "Second Proposal")
	w.rfc("new", "Third Proposal")
	w.rfc("link", "--update", "3", "1")
	w.rfc("link", "-o", "3", "2")
	w.head(12, "docs/rfc/RFC-0001-first-proposal.md")
	w.head(10, "docs/rfc/RFC-0002-second-proposal.md")
	w.head(14, "docs/rfc/RFC-0003-third-proposal.md")
	w.check()
}

func TestObsoleteAndUpdateNewRfcs(t *testing.T) {
	w := newWorkspace(t, "rfc")
	w.rfc("new", "First Proposal")
	w.rfc("new", "Second Proposal")
	w.rfc("new", "-o", "1", "--update", "2", "Third Proposal")
	w.head(10, "docs/rfc/RFC-0001-first-proposal.md")
	w.head(12, "docs/rfc/RFC-0002-second-proposal.md")
	w.head(14, "docs/rfc/RFC-0003-third-proposal.md")
	w.check()
}

// rfc names its options after its own relations, and has no --supersede.
// Nothing is written when an option is wrong.
func TestRfcLinkOptions(t *testing.T) {
	w := newWorkspace(t, "rfc")
	w.rfc("new", "First Proposal")
	w.rfc("new", "Second Proposal")
	w.rfc("link", "2", "1")
	w.rfc("link", "-s", "2", "1")
	w.rfc("new", "--supersede", "1", "Third Proposal")
	w.ls("docs/rfc")
	w.head(8, "docs/rfc/RFC-0001-first-proposal.md")
	w.check()
}

func TestRfcReports(t *testing.T) {
	w := newWorkspace(t, "rfc")
	w.rfc("new", "First Proposal")
	w.rfc("new", "-o", "1", "Second Proposal")
	w.rfc("list")
	w.rfc("list", "--toc")
	w.rfc("graph")
	w.rfc("graph", "--graphviz")
	w.check()
}

// rfc is configured by its own directory and variables, and ignores adr's.
func TestRfcConfiguration(t *testing.T) {
	w := newWorkspace(t, "rfc")
	w.rfc("config")
	w.rfc("init")
	w.rfc("config")
	w.mkdir("subdir")
	w.cd("subdir")
	w.rfc("config")
	w.cd("..")
	w.mkdir("docs/rfc/templates")
	w.write("docs/rfc/templates/template.md", "# NUMBER. TITLE\n\nDate: DATE\n\n## Status\n\nSTATUS\n\n## Proposal\n")
	w.rfc("config")
	w.write("custom.md", "# NUMBER. TITLE, from custom.md\n\nDate: DATE\n\n## Status\n\nSTATUS\n\n## Proposal\n")
	w.setenv("ADR_TEMPLATE", "custom.md")
	w.setenv("ADR_DATE", "2000-01-01")
	w.rfc("new", "Ignores the ADR settings")
	w.setenv("RFC_TEMPLATE", "custom.md")
	w.setenv("RFC_DATE", "2026-09-10")
	w.rfc("new", "Uses the RFC settings")
	w.cat("docs/rfc/RFC-0002-ignores-the-adr-settings.md")
	w.cat("docs/rfc/RFC-0003-uses-the-rfc-settings.md")
	w.check()
}

// adr and rfc keep separate logs in one project, and neither sees the other's
// records.
func TestAdrAndRfcLogsAreIndependent(t *testing.T) {
	w := newWorkspace(t, "rfc")
	w.adr("init")
	w.rfc("init")
	w.adr("new", "A Decision")
	w.rfc("new", "A Proposal")
	w.ls("docs")
	w.adr("list")
	w.rfc("list")
	w.rfc("link", "-o", "2", "a-decision")
	w.check()
}

func TestRfcHelp(t *testing.T) {
	w := newWorkspace(t, "rfc")
	w.rfc("help")
	w.rfc("help", "new")
	w.rfc("help", "link")
	w.rfc("nonsense")
	w.check()
}

// TestRfcCompletionScript checks that rfc generates its own completion script,
// which "make deploy" installs.  Its content is cobra's, and is not compared
// against a golden file.
func TestRfcCompletionScript(t *testing.T) {
	w := newWorkspace(t, "rfc")
	if script := w.output("rfc", "completion", "bash"); !strings.Contains(script, "__start_rfc") {
		t.Errorf("completion script does not define __start_rfc:\n%s", script)
	}
}
