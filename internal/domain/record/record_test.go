package record

import (
	"reflect"
	"testing"
)

// These tests need no temporary directory, no binary and no fakes: the package
// they cover is pure, so the whole of the record format can be exercised from
// literals.

// Only replacing a record retires it, and the zero Relation is no relation at
// all, so that a request which names none is recognised.
func TestRelations(t *testing.T) {
	if Relation(0).Valid() {
		t.Error("the zero Relation is valid")
	}
	if !Replaces.Valid() || !Replaces.Retires() {
		t.Error("Replaces is not a valid relation that retires its target")
	}
	if !Updates.Valid() || Updates.Retires() {
		t.Error("Updates is not a valid relation that leaves its target alone")
	}
}

func TestWordsForARelation(t *testing.T) {
	kind := Kind{
		Replaces: Words{Verb: "obsolete", Forward: "Obsoletes", Reverse: "Obsoleted by"},
		Updates:  Words{Verb: "update", Forward: "Updates", Reverse: "Updated by"},
	}
	if got := kind.WordsFor(Replaces); got != kind.Replaces {
		t.Errorf("WordsFor(Replaces) = %v, want %v", got, kind.Replaces)
	}
	if got := kind.WordsFor(Updates); got != kind.Updates {
		t.Errorf("WordsFor(Updates) = %v, want %v", got, kind.Updates)
	}
}

func TestVariableIsNamedAfterTheTool(t *testing.T) {
	if got, want := (Kind{Name: "rfc"}).Variable("DATE"), "RFC_DATE"; got != want {
		t.Errorf("Variable(DATE) = %q, want %q", got, want)
	}
}

func TestSlug(t *testing.T) {
	for _, test := range []struct {
		title string
		want  string
	}{
		{"The First Decision", "the-first-decision"},
		{"Something About Node.JS", "something-about-node-js"},
		{"Slash/Slash/Slash/", "slash-slash-slash"},
		{"-Bar-", "bar"},
		{"  spaced   out  ", "spaced-out"},
		{"Use PostgreSQL 15", "use-postgresql-15"},
		{"!?!", ""},
	} {
		if got := Slug(test.title); got != test.want {
			t.Errorf("Slug(%q) = %q, want %q", test.title, got, test.want)
		}
	}
}

func TestNumber(t *testing.T) {
	for _, test := range []struct {
		file string
		want int
	}{
		{"docs/adr/0001-first.md", 1},
		{"docs/adr/0009-ninth.md", 9},
		{"docs/adr/0010-tenth.md", 10},
		{"docs/adr/templates", 0},
	} {
		if got := Number(test.file); got != test.want {
			t.Errorf("Number(%q) = %d, want %d", test.file, got, test.want)
		}
	}
}

func TestFileNamePadsTheNumberToFourDigits(t *testing.T) {
	for _, test := range []struct {
		number int
		title  string
		want   string
	}{
		{1, "First Decision", "docs/adr/0001-first-decision.md"},
		{10, "Tenth Decision", "docs/adr/0010-tenth-decision.md"},
		{9, "Something About Node.JS", "docs/adr/0009-something-about-node-js.md"},
	} {
		if got := FileName("docs/adr", test.number, test.title); got != test.want {
			t.Errorf("FileName(%d, %q) = %q, want %q", test.number, test.title, got, test.want)
		}
	}
}

func TestIsRecordName(t *testing.T) {
	for _, test := range []struct {
		name string
		want bool
	}{
		{"0001-first.md", true},
		{"0010-tenth-decision.md", true},
		{"README.md", false},
		{"templates", false},
		{"template.md", false},
	} {
		if got := IsRecordName(test.name); got != test.want {
			t.Errorf("IsRecordName(%q) = %v, want %v", test.name, got, test.want)
		}
	}
}

func TestTitleIsTakenFromTheFirstLine(t *testing.T) {
	lines := []string{"# 1. First Record", "", "Date: 1992-01-12"}
	if got, want := Title(lines), "1. First Record"; got != want {
		t.Errorf("Title() = %q, want %q", got, want)
	}
	if got := Title(nil); got != "" {
		t.Errorf("Title(nil) = %q, want an empty string", got)
	}
}

var recordWithLinks = []string{
	"# 3. Third Record",
	"",
	"Date: 1992-01-12",
	"",
	"## Status",
	"",
	"Accepted",
	"",
	"Supercedes [1. First Record](0001-first-record.md)",
	"",
	"Updated by [4. Fourth Record](0004-fourth-record.md)",
	"",
	"## Context",
	"",
	"Some context.",
}

func TestStatusReadsTheStatusSection(t *testing.T) {
	want := []string{
		"Accepted",
		"Supercedes [1. First Record](0001-first-record.md)",
		"Updated by [4. Fourth Record](0004-fourth-record.md)",
	}
	if got := Status(recordWithLinks); !reflect.DeepEqual(got, want) {
		t.Errorf("Status() = %q, want %q", got, want)
	}
}

func TestLinksReadsLinksFromTheStatusSection(t *testing.T) {
	want := []Link{{Type: "Supercedes", Target: 1}, {Type: "Updated by", Target: 4}}
	if got := Links(recordWithLinks); !reflect.DeepEqual(got, want) {
		t.Errorf("Links() = %v, want %v", got, want)
	}
}

func TestAddLinkAppendsToTheStatusSection(t *testing.T) {
	lines := []string{"# 1. First", "", "## Status", "", "Accepted", "", "## Context", "", "Context."}
	want := []string{
		"# 1. First", "", "## Status", "", "Accepted", "",
		"Updated by [2. Second](0002-second.md)", "",
		"## Context", "", "Context.",
	}
	got := AddLink(lines, "Updated by", "2. Second", "docs/adr/0002-second.md")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("AddLink() =\n%q\nwant\n%q", got, want)
	}
}

// A record whose status section is the last in the file is left alone.  See the
// note on AddLink: the behaviour is preserved deliberately.
func TestAddLinkLeavesATrailingStatusSectionAlone(t *testing.T) {
	lines := []string{"# 2. Second", "", "## Status", "", "Accepted"}
	if got := AddLink(lines, "Updates", "1. First", "0001-first.md"); !reflect.DeepEqual(got, lines) {
		t.Errorf("AddLink() = %q, want it unchanged", got)
	}
}

func TestRemoveStatusClosesTheGapItLeaves(t *testing.T) {
	lines := []string{
		"# 1. First", "", "## Status", "", "Accepted", "",
		"Superceded by [2. Second](0002-second.md)", "",
		"## Context", "", "Accepted",
	}
	// Only the status section is touched: the "Accepted" in the context stays.
	want := []string{
		"# 1. First", "", "## Status", "",
		"Superceded by [2. Second](0002-second.md)", "",
		"## Context", "", "Accepted",
	}
	if got := RemoveStatus(lines, "Accepted"); !reflect.DeepEqual(got, want) {
		t.Errorf("RemoveStatus() =\n%q\nwant\n%q", got, want)
	}
}

func TestRenderReplacesEachPlaceholderOnceALine(t *testing.T) {
	template := "# NUMBER. TITLE\n\nDate: DATE\n\n## Status\n\nSTATUS\n"
	want := "# 7. Use MySQL\n\nDate: 1992-01-12\n\n## Status\n\nAccepted\n"
	if got := Render(template, 7, "Use MySQL", "1992-01-12", "Accepted"); got != want {
		t.Errorf("Render() =\n%s\nwant\n%s", got, want)
	}
}

// A title that happens to contain a placeholder is inserted, not rescanned.
func TestRenderDoesNotRescanTheTitle(t *testing.T) {
	want := "# 1. The NUMBER of things\n"
	if got := Render("# NUMBER. TITLE\n", 1, "The NUMBER of things", "", ""); got != want {
		t.Errorf("Render() = %q, want %q", got, want)
	}
}

// TestFindMatchesTheNumberOrTheFileName covers what a reference to a record
// may be: its number, its path, or part of its file name, but never part of
// the directory the log is in.
func TestFindMatchesTheNumberOrTheFileName(t *testing.T) {
	records := []Record{
		{Number: 1, File: "/home/x2/docs/adr/0001-first.md"},
		{Number: 3, File: "/home/x2/docs/adr/0003-third.md"},
		{Number: 12, File: "/home/x2/docs/adr/0012-twelfth.md"},
	}
	for _, test := range []struct {
		ref  string
		want int // the number of the record found, or 0 for none
	}{
		{"3", 3},
		{"0003", 3},
		{"12", 12},
		{"2", 0}, // not record 12, whose name contains a 2, nor the directory x2
		{"third", 3},
		{"0012-tw", 12},
		{"/home/x2/docs/adr/0001-first.md", 1},
		{"adr", 0},
		{"docs", 0},
		{"nothing", 0},
	} {
		got, ok := Find(records, test.ref)
		switch {
		case test.want == 0 && ok:
			t.Errorf("Find(%q) = %v, want none", test.ref, got.File)
		case test.want != 0 && (!ok || got.Number != test.want):
			t.Errorf("Find(%q) = %v, %t, want record %d", test.ref, got.File, ok, test.want)
		}
	}
}
