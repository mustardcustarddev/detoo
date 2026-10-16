package clitest

import "testing"

// detoo looks across the logs the record tools keep.  These tests cover its
// doctor: what it reports of logs in their place and out of it, that it
// changes nothing, and the status it exits with.

// Logs started and added to by the tools are where the tools look, and named
// with their prefixes, so doctor finds nothing wrong.
func TestDoctorFindsNothingWrong(t *testing.T) {
	w := newWorkspace(t, "detoo")
	w.adr("init")
	w.adr("new", "Use PostgreSQL for persistence")
	w.sty("init")
	w.detoo("doctor")
	w.check()
}

// The logs of the design: decision records named before their prefix, RFCs
// under doc/, stories split between the directory sty used to look in and
// the one it looks in now, and PIOSEE records in their place.
func TestDoctorReportsLogsInTheWrongPlace(t *testing.T) {
	w := newWorkspace(t, "detoo")
	w.touch(
		"docs/adr/0001-record-architecture-decisions.md",
		"docs/adr/0002-use-markdown-for-documents.md",
		"docs/adr/README.md",
		"doc/rfc/0001-adopt-a-single-detoo-command.md",
		"docs/userstory/story-0001-describe-work-in-user-stories.md",
		"docs/sty/0002-find-documents-in-the-wrong-place.md",
		"docs/piosee/PIOSEE-0001-make-decisions-with-piosee.md",
	)
	w.detoo("doctor")
	w.ls("docs/sty")
	w.ls("doc/rfc")
	w.check()
}

// A marker file naming a former directory, and two directories that would
// move to the same place, need fixing by hand.  A log moved alone does not.
func TestDoctorReportsClashesAsNeedingFixingByHand(t *testing.T) {
	w := newWorkspace(t, "detoo")
	w.write(".sty-dir", "docs/sty\n")
	w.touch(
		"docs/sty/STORY-0001-first.md",
		"doc/userstory/STORY-0002-second.md",
		"doc/sty/STORY-0003-third.md",
		"doc/adr/ADR-0001-record-architecture-decisions.md",
	)
	w.detoo("doctor")
	w.check()
}

// doctor finds the project from anywhere inside it, as the tools find their
// logs.
func TestDoctorFindsTheProjectFromASubdirectory(t *testing.T) {
	w := newWorkspace(t, "detoo")
	w.touch("docs/rfc/0001-use-requests-for-comments.md")
	w.mkdir("src/app")
	w.cd("src/app")
	w.detoo("doctor")
	w.check()
}

func TestDetooHelp(t *testing.T) {
	w := newWorkspace(t, "detoo")
	w.detoo("help")
	w.detoo("help", "doctor")
	w.detoo("doctor", "extra")
	w.detoo("nonsense")
	w.check()
}
