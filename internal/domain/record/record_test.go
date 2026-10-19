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

// A kind may name its records by number alone, or begin each with a prefix,
// and may have begun them with another prefix before.
var (
	unprefixed = Kind{Name: "adr"}
	prefixed   = Kind{Name: "ust", Prefix: "STORY"}
	renamed    = Kind{Name: "plm", Prefix: "PROBLEM", FormerPrefixes: []string{"PIOSEE"}}
)

func TestNumber(t *testing.T) {
	for _, test := range []struct {
		kind Kind
		file string
		want int
	}{
		{unprefixed, "docs/adr/0001-first.md", 1},
		{unprefixed, "docs/adr/0009-ninth.md", 9},
		{unprefixed, "docs/adr/0010-tenth.md", 10},
		{unprefixed, "docs/adr/templates", 0},
		{prefixed, "docs/userstory/STORY-0001-first.md", 1},
		{prefixed, "docs/userstory/STORY-0010-tenth.md", 10},
		{prefixed, "docs/userstory/story-0002-lower-case.md", 2},
		{prefixed, "docs/userstory/0003-unprefixed.md", 3},
		{prefixed, "docs/userstory/epic-0004-other-prefix.md", 0},
		{prefixed, "STORY-0005/epic-0006-in-a-prefixed-directory.md", 0},
		{renamed, "docs/problem/PROBLEM-0007-current-prefix.md", 7},
		{renamed, "docs/piosee/PIOSEE-0008-former-prefix.md", 8},
	} {
		if got := test.kind.Number(test.file); got != test.want {
			t.Errorf("%s: Number(%q) = %d, want %d", test.kind.Name, test.file, got, test.want)
		}
	}
}

func TestFileNamePadsTheNumberToFourDigits(t *testing.T) {
	for _, test := range []struct {
		kind   Kind
		dir    string
		number int
		title  string
		want   string
	}{
		{unprefixed, "docs/adr", 1, "First Decision", "docs/adr/0001-first-decision.md"},
		{unprefixed, "docs/adr", 10, "Tenth Decision", "docs/adr/0010-tenth-decision.md"},
		{unprefixed, "docs/adr", 9, "Something About Node.JS", "docs/adr/0009-something-about-node-js.md"},
		{prefixed, "docs/userstory", 1, "Export Invoices", "docs/userstory/STORY-0001-export-invoices.md"},
		{prefixed, "docs/userstory", 12, "Reorder", "docs/userstory/STORY-0012-reorder.md"},
	} {
		if got := test.kind.FileName(test.dir, test.number, test.title); got != test.want {
			t.Errorf("%s: FileName(%d, %q) = %q, want %q", test.kind.Name, test.number, test.title, got, test.want)
		}
	}
}

// A record's name is recognised with its kind's prefix in any case, or with
// none, as records were named before the prefixes were upper case, or with a
// prefix the kind had before, but not with another prefix.
func TestIsRecordName(t *testing.T) {
	for _, test := range []struct {
		kind Kind
		name string
		want bool
	}{
		{unprefixed, "0001-first.md", true},
		{unprefixed, "0010-tenth-decision.md", true},
		{unprefixed, "story-0001-first.md", false},
		{unprefixed, "README.md", false},
		{unprefixed, "templates", false},
		{unprefixed, "template.md", false},
		{prefixed, "STORY-0001-first.md", true},
		{prefixed, "STORY-0010-tenth.md", true},
		{prefixed, "story-0001-first.md", true},
		{prefixed, "0001-first.md", true},
		{prefixed, "epic-0001-first.md", false},
		{prefixed, "story-first.md", false},
		{prefixed, "storyline-0001-first.md", false},
		{prefixed, "README.md", false},
		{renamed, "PROBLEM-0001-first.md", true},
		{renamed, "PIOSEE-0001-first.md", true},
		{renamed, "piosee-0001-first.md", true},
		{renamed, "STORY-0001-first.md", false},
	} {
		if got := test.kind.IsRecordName(test.name); got != test.want {
			t.Errorf("%s: IsRecordName(%q) = %v, want %v", test.kind.Name, test.name, got, test.want)
		}
	}
}

// A record named in an older way has the name the kind gives it today, and a
// file that is not a record has none.
func TestCanonicalName(t *testing.T) {
	for _, test := range []struct {
		kind Kind
		name string
		want string
		ok   bool
	}{
		{prefixed, "STORY-0001-first.md", "STORY-0001-first.md", true},
		{prefixed, "story-0001-first.md", "STORY-0001-first.md", true},
		{prefixed, "Story-0001-first.md", "STORY-0001-first.md", true},
		{prefixed, "0001-first.md", "STORY-0001-first.md", true},
		{prefixed, "epic-0001-first.md", "", false},
		{prefixed, "README.md", "", false},
		{unprefixed, "0001-first.md", "0001-first.md", true},
		{renamed, "PIOSEE-0001-first.md", "PROBLEM-0001-first.md", true},
		{renamed, "piosee-0001-first.md", "PROBLEM-0001-first.md", true},
		{renamed, "PROBLEM-0001-first.md", "PROBLEM-0001-first.md", true},
	} {
		got, ok := test.kind.CanonicalName(test.name)
		if got != test.want || ok != test.ok {
			t.Errorf("%s: CanonicalName(%q) = %q, %v, want %q, %v", test.kind.Name, test.name, got, ok, test.want, test.ok)
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
	"Supersedes [1. First Record](0001-first-record.md)",
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
		"Supersedes [1. First Record](0001-first-record.md)",
		"Updated by [4. Fourth Record](0004-fourth-record.md)",
	}
	if got := Status(recordWithLinks); !reflect.DeepEqual(got, want) {
		t.Errorf("Status() = %q, want %q", got, want)
	}
}

func TestLinksReadsLinksFromTheStatusSection(t *testing.T) {
	want := []Link{{Type: "Supersedes", Target: 1}, {Type: "Updated by", Target: 4}}
	if got := unprefixed.Links(recordWithLinks); !reflect.DeepEqual(got, want) {
		t.Errorf("Links() = %v, want %v", got, want)
	}
}

// A link's target is numbered by its file name, read as its kind names files,
// so that a link to a record named in any of the kind's ways is found, and one
// to a file that is not a record of the kind is not.
func TestLinksReadTheTargetsNumberAfterThePrefix(t *testing.T) {
	lines := []string{
		"# 3. Third Story",
		"",
		"## Status",
		"",
		"Draft",
		"",
		"Replaces [1. First Story](STORY-0001-first-story.md)",
		"",
		"Updates [2. Second Story](story-0002-second-story.md)",
		"",
		"Updates [4. An Epic](epic-0004-an-epic.md)",
		"",
		"## Story",
	}
	want := []Link{{Type: "Replaces", Target: 1}, {Type: "Updates", Target: 2}}
	if got := prefixed.Links(lines); !reflect.DeepEqual(got, want) {
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
		"Superseded by [2. Second](0002-second.md)", "",
		"## Context", "", "Accepted",
	}
	// Only the status section is touched: the "Accepted" in the context stays.
	want := []string{
		"# 1. First", "", "## Status", "",
		"Superseded by [2. Second](0002-second.md)", "",
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
