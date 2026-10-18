package doctor

import (
	"context"
	"maps"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/mustardcustard/detoo/internal/adr"
	"github.com/mustardcustard/detoo/internal/domain/record"
	"github.com/mustardcustard/detoo/internal/plm"
	"github.com/mustardcustard/detoo/internal/rfc"
	"github.com/mustardcustard/detoo/internal/sop"
	"github.com/mustardcustard/detoo/internal/sty"
)

// fakeHistory is a git repository's history: the day of 1992 each file was
// first committed, by its path.
type fakeHistory map[string]int

func (h fakeHistory) Repository(context.Context) (bool, error) { return true, nil }

func (h fakeHistory) Added(_ context.Context, name string) (time.Time, bool, error) {
	day, ok := h[name]
	if !ok {
		return time.Time{}, false, nil
	}
	return time.Date(1992, 1, day, 0, 0, 0, 0, time.UTC), true, nil
}

var allKinds = []record.Kind{adr.Kind, rfc.Kind, sty.Kind, plm.Kind, sop.Kind}

func checkAll(t *testing.T, kind record.Kind, project fakeProject, history History) Log {
	t.Helper()
	result, err := Handler{Kinds: allKinds, Project: project, History: history}.Handle(context.Background(), Query{})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	for _, l := range result.Logs {
		if l.Kind.Name == kind.Name {
			return l
		}
	}
	t.Fatalf("no %s log", kind.Name)
	return Log{}
}

// A slip of a log's prefix is renamed to the prefix.  Another kind's prefix
// is not a slip, nor is a prefix far from any.
func TestASlipOfThePrefixIsMisnamed(t *testing.T) {
	log := checkAll(t, sty.Kind, fakeProject{dirs: map[string][]string{
		"docs/userstory": {"STORY-0001-first.md", "STOY-0002-second.md", "Stories-0003-third.md", "RFC-0001-elsewhere.md", "TICKET-0004-fourth.md"},
	}}, nil)
	want := []Rename{{Dir: "docs/userstory", From: "STOY-0002-second.md", To: "STORY-0002-second.md"}}
	if !reflect.DeepEqual(log.Misnamed, want) {
		t.Errorf("Misnamed = %+v, want %+v", log.Misnamed, want)
	}
	if log.Documents() != 2 {
		t.Errorf("Documents() = %d, want 2", log.Documents())
	}
	for _, prefix := range []string{"AD", "ADRS", "ADX"} {
		if got := nearest(prefix, allKinds); got != "adr" {
			t.Errorf("nearest(%q) = %q, want adr", prefix, got)
		}
	}
	for _, prefix := range []string{"ADR", "RFC", "X", "STORIES", "TICKET"} {
		if got := nearest(prefix, allKinds); got != "" {
			t.Errorf("nearest(%q) = %q, want none", prefix, got)
		}
	}
}

// Outside a git repository, which of two documents sharing a number came
// first is not known, so they need renumbering by hand.  Each but the first
// by name is shown with a number after the last, and a slip of the prefix
// counts by the name it is renamed to.
func TestASharedNumberWithoutHistoryNeedsFixingByHand(t *testing.T) {
	log := checkAll(t, sty.Kind, fakeProject{dirs: map[string][]string{
		"docs/userstory": {"STORY-0001-first.md", "STORY-0002-second.md", "STOY-0002-another.md", "STORY-0003-third.md"},
	}}, nil)
	wantDuplicates := []Duplicate{{
		Dir: "docs/userstory", Number: 2,
		Documents: []string{"STORY-0002-another.md", "STORY-0002-second.md"},
		Clash:     NoHistory,
	}}
	if !reflect.DeepEqual(log.Duplicates, wantDuplicates) {
		t.Errorf("Duplicates = %+v, want %+v", log.Duplicates, wantDuplicates)
	}
	wantRenumbers := []Renumber{{
		Dir: "docs/userstory", File: "STORY-0002-second.md",
		From: "STORY-0002-second.md", To: "STORY-0004-second.md",
		Duplicate: true, Clash: NoHistory,
	}}
	if !reflect.DeepEqual(log.Renumbers, wantRenumbers) {
		t.Errorf("Renumbers = %+v, want %+v", log.Renumbers, wantRenumbers)
	}
	if log.Problems() != 2 || log.ByHand() != 1 || log.Renumberable() != 0 {
		t.Errorf("Problems() = %d, ByHand() = %d, Renumberable() = %d, want 2, 1, 0", log.Problems(), log.ByHand(), log.Renumberable())
	}
}

// In a git repository, the document committed first keeps the number, and
// the other goes after the last committed before it, from there on.  The
// documents after it move up only as far as they must, so a gap in the
// numbers takes up the shift.
func TestASharedNumberIsPlacedByHistory(t *testing.T) {
	log := checkAll(t, adr.Kind, fakeProject{dirs: map[string][]string{
		"docs/adr": {"ADR-0001-a.md", "ADR-0002-b.md", "ADR-0002-c.md", "ADR-0003-d.md", "ADR-0005-e.md"},
	}}, fakeHistory{
		"docs/adr/ADR-0001-a.md": 1,
		"docs/adr/ADR-0002-c.md": 2,
		"docs/adr/ADR-0002-b.md": 3,
		"docs/adr/ADR-0003-d.md": 4,
		"docs/adr/ADR-0005-e.md": 5,
	})
	wantDuplicates := []Duplicate{{Dir: "docs/adr", Number: 2, Documents: []string{"ADR-0002-c.md", "ADR-0002-b.md"}}}
	if !reflect.DeepEqual(log.Duplicates, wantDuplicates) {
		t.Errorf("Duplicates = %+v, want %+v", log.Duplicates, wantDuplicates)
	}
	wantRenumbers := []Renumber{
		{Dir: "docs/adr", File: "ADR-0002-b.md", From: "ADR-0002-b.md", To: "ADR-0003-b.md", Duplicate: true, After: "ADR-0002-c.md", Before: "ADR-0004-d.md"},
		{Dir: "docs/adr", File: "ADR-0003-d.md", From: "ADR-0003-d.md", To: "ADR-0004-d.md"},
	}
	if !reflect.DeepEqual(log.Renumbers, wantRenumbers) {
		t.Errorf("Renumbers = %+v, want %+v", log.Renumbers, wantRenumbers)
	}
	if log.ByHand() != 0 || log.Renumberable() != 2 {
		t.Errorf("ByHand() = %d, Renumberable() = %d, want 0, 2", log.ByHand(), log.Renumberable())
	}
}

// A document git has not committed came last, so it takes the number after
// the last, and nothing else moves.
func TestAnUncommittedDocumentSharingANumberGoesLast(t *testing.T) {
	log := checkAll(t, sty.Kind, fakeProject{dirs: map[string][]string{
		"docs/userstory": {"STORY-0001-first.md", "STORY-0002-second.md", "STOY-0001-new.md"},
	}}, fakeHistory{
		"docs/userstory/STORY-0001-first.md":  1,
		"docs/userstory/STORY-0002-second.md": 2,
	})
	want := []Renumber{{
		Dir: "docs/userstory", File: "STOY-0001-new.md", From: "STORY-0001-new.md", To: "STORY-0003-new.md",
		Duplicate: true, After: "STORY-0002-second.md",
	}}
	if !reflect.DeepEqual(log.Renumbers, want) {
		t.Errorf("Renumbers = %+v, want %+v", log.Renumbers, want)
	}
}

// historyWorkspace is a project in memory, in a git repository with a
// history.
type historyWorkspace struct {
	*memWorkspace
	fakeHistory
}

// A fix asked to renumber renames the documents that share a number, and
// those moved up to make room, changes the number in their titles, and the
// links to them, text and all.  The slip of the prefix is renamed first.
func TestFixRenumbersByHistory(t *testing.T) {
	w := &memWorkspace{files: map[string]string{
		"docs/adr/ADR-0001-a.md":  "# 1. A\n",
		"docs/adr/ADR-0002-b.md":  "# 2. B\n\nSee [3. D](ADR-0003-d.md).\n",
		"docs/adr/ADRS-0002-c.md": "# 2. C\n",
		"docs/adr/ADR-0003-d.md":  "# 3. D\n\nAfter [2. B](ADR-0002-b.md) and [2. C](ADRS-0002-c.md).\n",
		"README.md":               "Read [D](docs/adr/ADR-0003-d.md).\n",
	}}
	g := historyWorkspace{w, fakeHistory{
		"docs/adr/ADR-0001-a.md":  1,
		"docs/adr/ADR-0002-b.md":  2,
		"docs/adr/ADRS-0002-c.md": 3,
		"docs/adr/ADR-0003-d.md":  4,
	}}
	fixed, err := Fixer{Kinds: []record.Kind{adr.Kind}, Project: w, Git: g}.Handle(context.Background(), Command{Renumber: true})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	want := map[string]string{
		"docs/adr/ADR-0001-a.md": "# 1. A\n",
		"docs/adr/ADR-0002-b.md": "# 2. B\n\nSee [4. D](ADR-0004-d.md).\n",
		"docs/adr/ADR-0003-c.md": "# 3. C\n",
		"docs/adr/ADR-0004-d.md": "# 4. D\n\nAfter [2. B](ADR-0002-b.md) and [3. C](ADR-0003-c.md).\n",
		"README.md":              "Read [D](docs/adr/ADR-0004-d.md).\n",
	}
	if !reflect.DeepEqual(w.files, want) {
		t.Errorf("files are\n%v\nwant\n%v", w.files, want)
	}
	if !fixed.After.Healthy() || fixed.Count() != 2 || fixed.Renumbered() != 2 {
		t.Errorf("After.Healthy() = %v, Count() = %d, Renumbered() = %d, want true, 2, 2", fixed.After.Healthy(), fixed.Count(), fixed.Renumbered())
	}
	wantRelinked := []Relinked{
		{"README.md", 1, false},
		{"docs/adr/ADR-0002-b.md", 1, false},
		{"docs/adr/ADR-0003-c.md", 0, true},
		{"docs/adr/ADR-0004-d.md", 1, true},
	}
	if !reflect.DeepEqual(fixed.Relinked, wantRelinked) {
		t.Errorf("Relinked = %v, want %v", fixed.Relinked, wantRelinked)
	}
	// One rename, two renumbers, and four documents changed.
	if fixed.Changes() != 1+2+4 {
		t.Errorf("Changes() = %d, want 7", fixed.Changes())
	}
}

// A fix not asked to renumber leaves documents that share a number as they
// are, and fixes the rest.
func TestFixRenumbersOnlyWhenAsked(t *testing.T) {
	w := &memWorkspace{files: map[string]string{
		"docs/adr/ADR-0001-a.md": "",
		"docs/adr/ADR-0001-b.md": "",
		"docs/adr/0002-c.md":     "",
	}}
	g := historyWorkspace{w, fakeHistory{}}
	fixed, err := Fixer{Kinds: []record.Kind{adr.Kind}, Project: w, Git: g}.Handle(context.Background(), Command{})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	want := []string{"docs/adr/ADR-0001-a.md", "docs/adr/ADR-0001-b.md", "docs/adr/ADR-0002-c.md"}
	if got := slices.Sorted(maps.Keys(w.files)); !reflect.DeepEqual(got, want) {
		t.Errorf("files are %v, want %v", got, want)
	}
	if fixed.After.Problems() != 1 || fixed.After.ByHand() != 0 || fixed.After.Renumberable() != 1 {
		t.Errorf("After has %d problems, %d by hand, %d to renumber, want 1, 0, 1",
			fixed.After.Problems(), fixed.After.ByHand(), fixed.After.Renumberable())
	}
}
