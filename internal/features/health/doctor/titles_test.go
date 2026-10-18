package doctor

import (
	"reflect"
	"testing"

	"github.com/mustardcustard/detoo/internal/adr"
)

// titled is a project with the decision records in docs/adr, each with the
// content given, by file name.
func titled(contents map[string]string) fakeProject {
	project := fakeProject{dirs: map[string][]string{"docs/adr": nil}, contents: map[string]string{}}
	for name, content := range contents {
		project.dirs["docs/adr"] = append(project.dirs["docs/adr"], name)
		project.contents["docs/adr/"+name] = content
	}
	return project
}

// A title matches its file name when it has the file's number, and its words
// are the file's once case and punctuation are set aside, as the tools name a
// document from its title.  A document with no title, such as an empty one,
// has none to match.
func TestTitlesThatMatchTheirFileNamesAreHealthy(t *testing.T) {
	log := check(t, adr.Kind, titled(map[string]string{
		"ADR-0001-do-stuff.md":       "# 1. Do stuff\n\nDate: 2026-10-18\n",
		"ADR-0002-use-postgresql.md": "# 2. Use PostgreSQL!\n",
		"ADR-0003-empty.md":          "",
		"ADR-0004-notes.md":          "Notes, with no title.\n",
	}))
	if log.Problems() != 0 || log.Mistitled != nil {
		t.Errorf("Problems() = %d, Mistitled = %v, want 0, none", log.Problems(), log.Mistitled)
	}
}

// A title numbered other than its file is mistitled, and needs fixing by
// hand, since which of the two is wrong is not known.
func TestATitleNumberedOtherThanItsFileIsMistitled(t *testing.T) {
	log := check(t, adr.Kind, titled(map[string]string{
		"ADR-0001-do-stuff.md": "# 3. Do stuff\n",
	}))
	want := []Mistitle{{Dir: "docs/adr", File: "ADR-0001-do-stuff.md", Title: "3. Do stuff", Number: true}}
	if !reflect.DeepEqual(log.Mistitled, want) {
		t.Errorf("Mistitled = %v, want %v", log.Mistitled, want)
	}
	if log.Problems() != 1 || log.ByHand() != 1 {
		t.Errorf("Problems() = %d, ByHand() = %d, want 1, 1", log.Problems(), log.ByHand())
	}
}

// A title whose words are not its file's is mistitled.
func TestATitleWithOtherWordsThanItsFileIsMistitled(t *testing.T) {
	log := check(t, adr.Kind, titled(map[string]string{
		"ADR-0001-do-stuff.md": "# 1. Do other stuff\n",
	}))
	want := []Mistitle{{Dir: "docs/adr", File: "ADR-0001-do-stuff.md", Title: "1. Do other stuff", Words: true}}
	if !reflect.DeepEqual(log.Mistitled, want) {
		t.Errorf("Mistitled = %v, want %v", log.Mistitled, want)
	}
}

// A title without a number is mistitled, as the tools number every title.
func TestATitleWithoutANumberIsMistitled(t *testing.T) {
	log := check(t, adr.Kind, titled(map[string]string{
		"ADR-0001-do-stuff.md": "# Do stuff\n",
	}))
	want := []Mistitle{{Dir: "docs/adr", File: "ADR-0001-do-stuff.md", Title: "Do stuff", Number: true}}
	if !reflect.DeepEqual(log.Mistitled, want) {
		t.Errorf("Mistitled = %v, want %v", log.Mistitled, want)
	}
}

// A document named before its log's prefix is checked by the number and the
// words in its name, so that it is mistitled as well as misnamed.
func TestADocumentNamedWithoutThePrefixIsCheckedToo(t *testing.T) {
	log := check(t, adr.Kind, titled(map[string]string{
		"0001-do-stuff.md": "# 2. Do things\n",
	}))
	want := []Mistitle{{Dir: "docs/adr", File: "0001-do-stuff.md", Title: "2. Do things", Number: true, Words: true}}
	if !reflect.DeepEqual(log.Mistitled, want) {
		t.Errorf("Mistitled = %v, want %v", log.Mistitled, want)
	}
	if log.Problems() != 2 {
		t.Errorf("Problems() = %d, want 2: misnamed and mistitled", log.Problems())
	}
}
