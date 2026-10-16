package clitest

import (
	"strings"
	"testing"
)

// The sty tool shares its engine with adr and rfc.  The adr tests cover the
// behaviour every tool shares, and the rfc tests cover how a tool's definition
// names its directory, its marker file and its variables.  These tests cover
// what the sty tool's definition changes: its first record, the status and
// sections of a new story, the words and options for its links, the heading of
// its reports, and its help.

func TestInitStyRepository(t *testing.T) {
	w := newWorkspace(t, "sty")
	w.sty("init")
	w.ls("docs/userstory")
	w.cat("docs/userstory/STORY-0001-describe-work-in-user-stories.md")
	w.check()
}

// A new story is a draft, in the sections of the story template.
func TestCreateStory(t *testing.T) {
	w := newWorkspace(t, "sty")
	w.sty("new", "Export invoices for a month")
	w.cat("docs/userstory/STORY-0001-export-invoices-for-a-month.md")
	w.check()
}

// A story that is split is replaced by each of the smaller stories, and loses
// its Draft status; a story that is updated keeps it.
func TestSplitAndUpdateStories(t *testing.T) {
	w := newWorkspace(t, "sty")
	w.sty("new", "Reorder")
	w.sty("new", "Save sessions")
	w.sty("new", "-r", "1", "Reorder a previous purchase")
	w.sty("new", "--replace", "1", "Adjust quantities when reordering")
	w.sty("new", "Warn of clashing sessions")
	w.sty("link", "-u", "5", "2")
	w.head(12, "docs/userstory/STORY-0001-reorder.md")
	w.head(10, "docs/userstory/STORY-0002-save-sessions.md")
	w.head(10, "docs/userstory/STORY-0004-adjust-quantities-when-reordering.md")
	w.head(10, "docs/userstory/STORY-0005-warn-of-clashing-sessions.md")
	w.check()
}

// sty names its options after its own relations, and has neither --supersede
// nor --obsolete.  Nothing is written when an option is wrong.
func TestStyLinkOptions(t *testing.T) {
	w := newWorkspace(t, "sty")
	w.sty("new", "First Story")
	w.sty("new", "Second Story")
	w.sty("link", "2", "1")
	w.sty("link", "-o", "2", "1")
	w.sty("new", "--supersede", "1", "Third Story")
	w.ls("docs/userstory")
	w.head(8, "docs/userstory/STORY-0001-first-story.md")
	w.check()
}

func TestStyReports(t *testing.T) {
	w := newWorkspace(t, "sty")
	w.sty("new", "First Story")
	w.sty("new", "-r", "1", "Second Story")
	w.sty("list", "--toc")
	w.sty("graph")
	w.check()
}

// Stories named before the prefix was upper case are still found, in order of
// their numbers, and a new story takes the next number and today's prefix.
func TestStyReadsStoriesNamedInOlderWays(t *testing.T) {
	w := newWorkspace(t, "sty")
	w.mkdir("docs/userstory")
	w.write("docs/userstory/story-0001-reorder.md", "# 1. Reorder\n\n## Status\n\nDraft\n\n## Story\n\nStory.\n")
	w.write("docs/userstory/0002-save-sessions.md", "# 2. Save sessions\n\n## Status\n\nDraft\n\n## Story\n\nStory.\n")
	w.sty("new", "--replace", "1", "Reorder a previous purchase")
	w.sty("list")
	w.head(6, "docs/userstory/story-0001-reorder.md")
	w.check()
}

func TestStyHelp(t *testing.T) {
	w := newWorkspace(t, "sty")
	w.sty("help")
	w.sty("help", "new")
	w.sty("help", "link")
	w.sty("nonsense")
	w.check()
}

// TestStyCompletionScript checks that sty generates its own completion script,
// which "make deploy" installs.
func TestStyCompletionScript(t *testing.T) {
	w := newWorkspace(t, "sty")
	if script := w.output("sty", "completion", "bash"); !strings.Contains(script, "__start_sty") {
		t.Errorf("completion script does not define __start_sty:\n%s", script)
	}
}
