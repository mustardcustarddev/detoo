package clitest

import (
	"strings"
	"testing"
)

// The ust tool shares its engine with adr and rfc.  The adr tests cover the
// behaviour every tool shares, and the rfc tests cover how a tool's definition
// names its directory and its variables.  These tests cover what the ust
// tool's definition changes: the README its log opens with, the status and
// sections of a new story, the words and options for its links, the heading of
// its reports, and its help.

// A log of user stories opens with a README, not a story, so the first story
// is number 1.  Running init again leaves the README alone.
func TestInitUstRepository(t *testing.T) {
	w := newWorkspace(t, "ust")
	w.ust("init")
	w.ls("docs/userstory")
	w.cat("docs/userstory/README.md")
	w.ust("init")
	w.ust("new", "Export invoices for a month")
	w.ust("list")
	w.check()
}

// A new story is a draft, in the sections of the story template.
func TestCreateStory(t *testing.T) {
	w := newWorkspace(t, "ust")
	w.ust("new", "Export invoices for a month")
	w.cat("docs/userstory/STORY-0001-export-invoices-for-a-month.md")
	w.check()
}

// A story that is split is replaced by each of the smaller stories, and loses
// its Draft status; a story that is updated keeps it.
func TestSplitAndUpdateStories(t *testing.T) {
	w := newWorkspace(t, "ust")
	w.ust("new", "Reorder")
	w.ust("new", "Save sessions")
	w.ust("new", "-r", "1", "Reorder a previous purchase")
	w.ust("new", "--replace", "1", "Adjust quantities when reordering")
	w.ust("new", "Warn of clashing sessions")
	w.ust("link", "-u", "5", "2")
	w.head(12, "docs/userstory/STORY-0001-reorder.md")
	w.head(10, "docs/userstory/STORY-0002-save-sessions.md")
	w.head(10, "docs/userstory/STORY-0004-adjust-quantities-when-reordering.md")
	w.head(10, "docs/userstory/STORY-0005-warn-of-clashing-sessions.md")
	w.check()
}

// ust names its options after its own relations, and has neither --supersede
// nor --obsolete.  Nothing is written when an option is wrong.
func TestUstLinkOptions(t *testing.T) {
	w := newWorkspace(t, "ust")
	w.ust("new", "First Story")
	w.ust("new", "Second Story")
	w.ust("link", "2", "1")
	w.ust("link", "-o", "2", "1")
	w.ust("new", "--supersede", "1", "Third Story")
	w.ls("docs/userstory")
	w.head(8, "docs/userstory/STORY-0001-first-story.md")
	w.check()
}

func TestUstReports(t *testing.T) {
	w := newWorkspace(t, "ust")
	w.ust("new", "First Story")
	w.ust("new", "-r", "1", "Second Story")
	w.ust("list", "--toc")
	w.ust("graph")
	w.check()
}

// Stories named before the prefix was upper case are still found, in order of
// their numbers, and a new story takes the next number and today's prefix.
func TestUstReadsStoriesNamedInOlderWays(t *testing.T) {
	w := newWorkspace(t, "ust")
	w.mkdir("docs/userstory")
	w.write("docs/userstory/story-0001-reorder.md", "# 1. Reorder\n\n## Status\n\nDraft\n\n## Story\n\nStory.\n")
	w.write("docs/userstory/0002-save-sessions.md", "# 2. Save sessions\n\n## Status\n\nDraft\n\n## Story\n\nStory.\n")
	w.ust("new", "--replace", "1", "Reorder a previous purchase")
	w.ust("list")
	w.head(6, "docs/userstory/story-0001-reorder.md")
	w.check()
}

func TestUstHelp(t *testing.T) {
	w := newWorkspace(t, "ust")
	w.ust("help")
	w.ust("help", "new")
	w.ust("help", "link")
	w.ust("nonsense")
	w.check()
}

// TestUstCompletionScript checks that ust generates its own completion script,
// which "make deploy" installs.
func TestUstCompletionScript(t *testing.T) {
	w := newWorkspace(t, "ust")
	if script := w.output("ust", "completion", "bash"); !strings.Contains(script, "__start_ust") {
		t.Errorf("completion script does not define __start_ust:\n%s", script)
	}
}
