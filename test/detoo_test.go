package clitest

import "testing"

// detoo looks across the logs the record tools keep.  These tests cover its
// doctor: what it reports of logs in their place and out of it, that it
// changes nothing unless asked to fix them, what a fix moves, renames and
// relinks, and the status it exits with.

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

// A fix moves RFCs from doc/ to docs/, and stories from the directory sty
// used to look in, renames every document to begin with its prefix, and
// changes the links to them, so that the tools find them again and doctor
// finds nothing wrong.
func TestDoctorFixMovesAndRenamesLogs(t *testing.T) {
	w := newWorkspace(t, "detoo")
	w.write("docs/adr/0001-record-architecture-decisions.md",
		"# 1. Record architecture decisions\n\n## Status\n\nAccepted\n\n## Decision\n\nWe propose it in [RFC 1](../../doc/rfc/0001-adopt-a-single-detoo-command.md).\n")
	w.write("doc/rfc/0001-adopt-a-single-detoo-command.md",
		"# 1. Adopt a single detoo command\n\n## Status\n\nDraft\n\n## Summary\n\nAs [ADR 1](../../docs/adr/0001-record-architecture-decisions.md) says.\n")
	w.write("docs/sty/story-0001-describe-work-in-user-stories.md",
		"# 1. Describe work in user stories\n\nDate: 1992-01-12\n\n## Status\n\nDraft\n\nUpdated by [2. Find documents in the wrong place](0002-find-documents-in-the-wrong-place.md)\n")
	w.write("docs/sty/0002-find-documents-in-the-wrong-place.md",
		"# 2. Find documents in the wrong place\n\nDate: 1992-01-12\n\n## Status\n\nDraft\n\nUpdates [1. Describe work in user stories](story-0001-describe-work-in-user-stories.md)\n")
	w.write("README.md", "See [the stories](docs/sty/) and [the RFCs](doc/rfc/).\n")
	w.detoo("doctor", "--fix")
	w.ls(".")
	w.ls("docs/rfc")
	w.ls("docs/userstory")
	w.cat("README.md")
	w.cat("docs/adr/ADR-0001-record-architecture-decisions.md")
	w.cat("docs/rfc/RFC-0001-adopt-a-single-detoo-command.md")
	w.cat("docs/userstory/STORY-0002-find-documents-in-the-wrong-place.md")
	w.sty("list")
	w.detoo("doctor")
	w.check()
}

// A fix leaves the clashes as they are, fixes what it can around them, and
// exits 1, since problems remain.
func TestDoctorFixLeavesClashesAlone(t *testing.T) {
	w := newWorkspace(t, "detoo")
	w.write(".sty-dir", "docs/sty\n")
	w.touch(
		"docs/sty/STORY-0001-first.md",
		"doc/userstory/STORY-0002-second.md",
		"doc/sty/story-0003-third.md",
		"doc/adr/0001-first.md",
		"doc/adr/adr-0001-first.md",
	)
	w.detoo("doctor", "--fix")
	w.ls("docs/sty")
	w.ls("doc/sty")
	w.ls("docs/adr")
	w.check()
}

// A log kept by adr-tools is under doc/adr, numbered without a prefix, and
// linked by file name.  A fix makes it a log adr keeps.
func TestDoctorFixMigratesFromAdrTools(t *testing.T) {
	w := newWorkspace(t, "detoo")
	w.write("doc/adr/0001-record-architecture-decisions.md",
		"# 1. Record architecture decisions\n\nDate: 2016-02-12\n\n## Status\n\nAccepted\n\n## Context\n\nWe need to record the architectural decisions made on this project.\n")
	w.write("doc/adr/0002-use-mysql.md",
		"# 2. Use MySQL\n\nDate: 2016-02-12\n\n## Status\n\nSuperseded by [3. Use PostgreSQL](0003-use-postgresql.md)\n\n## Context\n\nWe need a database.\n")
	w.write("doc/adr/0003-use-postgresql.md",
		"# 3. Use PostgreSQL\n\nDate: 2016-02-13\n\n## Status\n\nAccepted\n\nSupersedes [2. Use MySQL](0002-use-mysql.md)\n\n## Context\n\nWe need a better database.\n")
	w.detoo("doctor", "--fix")
	w.ls("docs/adr")
	w.cat("docs/adr/ADR-0002-use-mysql.md")
	w.adr("list")
	w.adr("new", "--supersede", "3", "Use SQLite")
	w.cat("docs/adr/ADR-0003-use-postgresql.md")
	w.check()
}

// In a git repository, a fix asked to move with git does so with git mv, so
// that git's index follows the files it tracks, and moves the files it does
// not track as it would without git.  Nothing is committed: git mv works on
// the index alone.
func TestDoctorFixMovesWithGit(t *testing.T) {
	w := newWorkspace(t, "detoo")
	w.touch(
		"doc/adr/0001-record-architecture-decisions.md",
		"doc/adr/0002-use-markdown-for-documents.md",
		"docs/userstory/story-0001-describe-work-in-user-stories.md",
	)
	w.git("init", "--quiet", "--initial-branch=main")
	w.git("add", "doc/adr/0001-record-architecture-decisions.md", "docs/userstory")
	w.detoo("doctor", "--fix", "--git")
	w.git("ls-files")
	w.ls("docs/adr")
	w.check()
}

// Moving with git needs a git repository, and a fix to make.
func TestDoctorFixWithGitNeedsARepository(t *testing.T) {
	w := newWorkspace(t, "detoo")
	w.touch("doc/adr/0001-record-architecture-decisions.md")
	w.detoo("doctor", "--fix", "--git")
	w.detoo("doctor", "--git")
	w.ls("doc/adr")
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
