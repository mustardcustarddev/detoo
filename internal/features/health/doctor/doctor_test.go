package doctor

import (
	"context"
	"reflect"
	"testing"

	"github.com/mustardcustard/detoo/internal/adr"
	"github.com/mustardcustard/detoo/internal/domain/record"
	"github.com/mustardcustard/detoo/internal/ust"
)

// fakeProject is a project held in memory: the files in each directory, and
// the content of the files that have any, by path.
type fakeProject struct {
	dirs     map[string][]string
	contents map[string]string
}

func (p fakeProject) Root(context.Context) (string, error) { return "/home/ann/handbook", nil }

func (p fakeProject) Files(_ context.Context, dir string) ([]string, bool, error) {
	files, ok := p.dirs[dir]
	return files, ok, nil
}

func (p fakeProject) Read(_ context.Context, name string) (string, error) {
	return p.contents[name], nil
}

func check(t *testing.T, kind record.Kind, project fakeProject) Log {
	t.Helper()
	result, err := Handler{Kinds: []record.Kind{kind}, Project: project}.Handle(context.Background(), Query{})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	if result.Root != "/home/ann/handbook" {
		t.Errorf("Root = %q, want the project's", result.Root)
	}
	return result.Logs[0]
}

// A log in its place, with every document named with the prefix, is healthy,
// and files that are not documents are not counted.
func TestAHealthyLogHasNoProblems(t *testing.T) {
	log := check(t, adr.Kind, fakeProject{dirs: map[string][]string{
		"docs/adr": {"ADR-0001-first.md", "ADR-0002-second.md", "README.md"},
	}})
	if log.Problems() != 0 || !log.Found() || log.Documents() != 2 {
		t.Errorf("Problems() = %d, Found() = %v, Documents() = %d, want 0, true, 2", log.Problems(), log.Found(), log.Documents())
	}
}

// A project that does not keep the log has no problems with it.
func TestAMissingLogIsNotAProblem(t *testing.T) {
	log := check(t, adr.Kind, fakeProject{})
	if log.Found() || log.Problems() != 0 || log.Dir != "docs/adr" {
		t.Errorf("Found() = %v, Problems() = %d, Dir = %q, want false, 0, docs/adr", log.Found(), log.Problems(), log.Dir)
	}
}

// Documents named before the prefix was upper case, or before there was one,
// are reported with the names they should have.
func TestDocumentsWithoutTheUpperCasePrefixAreMisnamed(t *testing.T) {
	log := check(t, ust.Kind, fakeProject{dirs: map[string][]string{
		"docs/userstory": {"story-0001-first.md", "0002-second.md", "STORY-0003-third.md"},
	}})
	want := []Rename{
		{Dir: "docs/userstory", From: "story-0001-first.md", To: "STORY-0001-first.md"},
		{Dir: "docs/userstory", From: "0002-second.md", To: "STORY-0002-second.md"},
	}
	if !reflect.DeepEqual(log.Misnamed, want) {
		t.Errorf("Misnamed = %v, want %v", log.Misnamed, want)
	}
	if log.ByHand() != 0 {
		t.Errorf("ByHand() = %d, want 0", log.ByHand())
	}
}

// A rename to a name that is there already clashes.
func TestARenameOntoAnExistingFileClashes(t *testing.T) {
	log := check(t, adr.Kind, fakeProject{dirs: map[string][]string{
		"docs/adr": {"0001-first.md", "ADR-0001-first.md"},
	}})
	want := []Rename{{Dir: "docs/adr", From: "0001-first.md", To: "ADR-0001-first.md", Clash: true, ClashWith: "ADR-0001-first.md"}}
	if !reflect.DeepEqual(log.Misnamed, want) {
		t.Errorf("Misnamed = %v, want %v", log.Misnamed, want)
	}
	if log.ByHand() != 1 {
		t.Errorf("ByHand() = %d, want 1", log.ByHand())
	}
}

// Two documents that would be renamed to the same name clash with each other,
// so that neither is renamed over the other.
func TestTwoRenamesToOneNameClash(t *testing.T) {
	log := check(t, adr.Kind, fakeProject{dirs: map[string][]string{
		"docs/adr": {"0001-first.md", "adr-0001-first.md", "0002-second.md"},
	}})
	want := []Rename{
		{Dir: "docs/adr", From: "0001-first.md", To: "ADR-0001-first.md", Clash: true, ClashWith: "adr-0001-first.md"},
		{Dir: "docs/adr", From: "adr-0001-first.md", To: "ADR-0001-first.md", Clash: true, ClashWith: "0001-first.md"},
		{Dir: "docs/adr", From: "0002-second.md", To: "ADR-0002-second.md"},
	}
	if !reflect.DeepEqual(log.Misnamed, want) {
		t.Errorf("Misnamed = %v, want %v", log.Misnamed, want)
	}
	if log.ByHand() != 2 {
		t.Errorf("ByHand() = %d, want 2", log.ByHand())
	}
}

func TestALogUnderDocBelongsUnderDocs(t *testing.T) {
	log := check(t, adr.Kind, fakeProject{dirs: map[string][]string{
		"doc/adr": {"ADR-0001-first.md"},
	}})
	want := []Move{{From: "doc/adr", To: "docs/adr", Reason: UnderDoc, Documents: 1}}
	if !reflect.DeepEqual(log.Misplaced, want) {
		t.Errorf("Misplaced = %v, want %v", log.Misplaced, want)
	}
}

// Stories left in docs/sty, where ust used to keep them, belong in
// docs/userstory, and are misnamed there too.
func TestALogInAFormerDirectoryBelongsInTheCurrentOne(t *testing.T) {
	log := check(t, ust.Kind, fakeProject{dirs: map[string][]string{
		"docs/sty": {"0001-first.md"},
	}})
	want := []Move{{From: "docs/sty", To: "docs/userstory", Reason: Former, Documents: 1}}
	if !reflect.DeepEqual(log.Misplaced, want) {
		t.Errorf("Misplaced = %v, want %v", log.Misplaced, want)
	}
	if len(log.Misnamed) != 1 || log.Misnamed[0].Dir != "docs/sty" {
		t.Errorf("Misnamed = %v, want the story in docs/sty", log.Misnamed)
	}
	if log.Problems() != 2 {
		t.Errorf("Problems() = %d, want 2", log.Problems())
	}
}

// When both docs/sty and docs/userstory exist, moving one onto the other
// would clash.
func TestAMoveOntoAnExistingDirectoryClashes(t *testing.T) {
	log := check(t, ust.Kind, fakeProject{dirs: map[string][]string{
		"docs/sty":       {"STORY-0001-first.md"},
		"docs/userstory": {"STORY-0002-second.md"},
	}})
	want := []Move{{From: "docs/sty", To: "docs/userstory", Reason: Former, Documents: 1, Clash: TargetExists}}
	if !reflect.DeepEqual(log.Misplaced, want) {
		t.Errorf("Misplaced = %v, want %v", log.Misplaced, want)
	}
	if got := []string{log.Dirs[0].Path, log.Dirs[1].Path}; !reflect.DeepEqual(got, []string{"docs/userstory", "docs/sty"}) {
		t.Errorf("Dirs = %v, want the current directory first", got)
	}
}

func TestTwoMovesToOnePlaceClash(t *testing.T) {
	log := check(t, ust.Kind, fakeProject{dirs: map[string][]string{
		"doc/userstory": {"STORY-0001-first.md"},
		"docs/sty":      {"STORY-0002-second.md"},
	}})
	want := []Move{
		{From: "doc/userstory", To: "docs/userstory", Reason: UnderDoc, Documents: 1},
		{From: "docs/sty", To: "docs/userstory", Reason: Former, Documents: 1, Clash: SharedTarget, ClashWith: "doc/userstory"},
	}
	if !reflect.DeepEqual(log.Misplaced, want) {
		t.Errorf("Misplaced = %v, want %v", log.Misplaced, want)
	}
}

// The result counts across every log, in the order of the kinds.
func TestTheResultCountsAcrossLogs(t *testing.T) {
	result, err := Handler{Kinds: []record.Kind{adr.Kind, ust.Kind}, Project: fakeProject{dirs: map[string][]string{
		"docs/adr":       {"0001-first.md"},
		"docs/sty":       {"STORY-0001-first.md"},
		"docs/userstory": {"STORY-0002-second.md"},
	}}}.Handle(context.Background(), Query{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Logs[0].Kind.Name != "adr" || result.Logs[1].Kind.Name != "ust" {
		t.Errorf("logs are not in the order of the kinds")
	}
	if result.Problems() != 2 || result.ByHand() != 1 || result.Documents() != 3 || result.Directories() != 3 || result.Healthy() {
		t.Errorf("Problems() = %d, ByHand() = %d, Documents() = %d, Directories() = %d, Healthy() = %v, want 2, 1, 3, 3, false",
			result.Problems(), result.ByHand(), result.Documents(), result.Directories(), result.Healthy())
	}
}
