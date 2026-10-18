package skills

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"path"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// memInstallation is a project held in memory: each skill's directory, by
// path, with the content of its SKILL.md, and each link, by path, with the
// directory it leads to.
type memInstallation struct {
	dirs     map[string]string
	links    map[string]string
	locked   map[string]Locked
	detected []string
}

func newMem() *memInstallation {
	return &memInstallation{dirs: map[string]string{}, links: map[string]string{}, locked: map[string]Locked{}}
}

func (m *memInstallation) Root(context.Context) (string, error) { return "/home/ann/app", nil }
func (m *memInstallation) Canonical() string                    { return Canonical }
func (m *memInstallation) SkillsDir(a Agent) string             { return a.Dir }
func (m *memInstallation) LockFile() string                     { return "skills-lock.json" }

func (m *memInstallation) Detected(_ context.Context, a Agent) (bool, error) {
	return slices.Contains(m.detected, a.Name), nil
}

// Hash is the content of SKILL.md, so that a test can tell versions apart.
func (m *memInstallation) Hash(files []File) string {
	for _, f := range files {
		if f.Path == "SKILL.md" {
			return string(f.Content)
		}
	}
	return ""
}

func (m *memInstallation) Install(_ context.Context, dir string, files []File) error {
	delete(m.links, dir)
	m.dirs[dir] = m.Hash(files)
	return nil
}

func (m *memInstallation) Link(_ context.Context, link, target string) error {
	delete(m.dirs, link)
	m.links[link] = target
	return nil
}

func (m *memInstallation) Entry(_ context.Context, name, target string) (EntryKind, error) {
	if to, ok := m.links[name]; ok {
		if to == target {
			return Linked, nil
		}
		return Foreign, nil
	}
	if _, ok := m.dirs[name]; ok {
		if name == target {
			return Linked, nil
		}
		return Copied, nil
	}
	return Missing, nil
}

func (m *memInstallation) Remove(_ context.Context, name string) error {
	delete(m.dirs, name)
	delete(m.links, name)
	return nil
}

// RemoveIfEmpty removes dir when no skill or link is left in it.
func (m *memInstallation) RemoveIfEmpty(_ context.Context, dir string) (bool, error) {
	for name := range m.dirs {
		if path.Dir(name) == dir {
			return false, nil
		}
	}
	for name := range m.links {
		if path.Dir(name) == dir {
			return false, nil
		}
	}
	return true, nil
}

func (m *memInstallation) Locked(context.Context) (map[string]Locked, error) {
	return maps.Clone(m.locked), nil
}

func (m *memInstallation) Lock(_ context.Context, name string, l Locked) error {
	m.locked[name] = l
	return nil
}

func (m *memInstallation) Unlock(_ context.Context, name string) error {
	delete(m.locked, name)
	return nil
}

// memFetcher offers skills by source, and fails for a source it has none of.
type memFetcher struct {
	offers  map[Source][]Skill
	fetched []Source
}

func (f *memFetcher) Fetch(_ context.Context, src Source) ([]Skill, error) {
	f.fetched = append(f.fetched, src)
	offered, ok := f.offers[src]
	if !ok {
		return nil, fmt.Errorf("%s cannot be reached", src)
	}
	return offered, nil
}

func skill(name, version string) Skill {
	return Skill{Name: name, Path: "skills/" + name + "/SKILL.md", Files: []File{{Path: "SKILL.md", Content: []byte(version)}}}
}

// heard records what progress was told, a line for each event.
type heard []string

func (h *heard) add(format string, a ...any)   { *h = append(*h, fmt.Sprintf(format, a...)) }
func (h *heard) Fetching(src Source)           { h.add("fetching %s", src) }
func (h *heard) Fetched(src Source, s []Skill) { h.add("fetched %d", len(s)) }
func (h *heard) Installing(t Targets) {
	var links []string
	for _, l := range t.Links {
		links = append(links, l.Agent.Name+"="+l.Dir)
	}
	h.add("installing into %s, linking %s", t.Dir, strings.Join(links, " "))
}
func (h *heard) Installed(name string)         { h.add("installed %s", name) }
func (h *heard) UpToDate(name string)          { h.add("up to date %s", name) }
func (h *heard) Gone(name string, src Source)  { h.add("gone %s", name) }
func (h *heard) Failed(name string, err error) { h.add("failed %s", name) }
func (h *heard) Removed(name string, paths []string) {
	h.add("removed %s from %s", name, strings.Join(paths, " "))
}

var detoo = Source{Kind: GitHub, Repo: DefaultSource}

func TestParseSource(t *testing.T) {
	for in, want := range map[string]Source{
		"mustardcustarddev/detoo":                        detoo,
		"mustardcustarddev/detoo#v1.2.0":                 {Kind: GitHub, Repo: DefaultSource, Ref: "v1.2.0"},
		"https://github.com/mustardcustarddev/detoo.git": detoo,
		"./skills":    {Kind: Local, Dir: "./skills"},
		"../detoo":    {Kind: Local, Dir: "../detoo"},
		"~/src/detoo": {Kind: Local, Dir: "~/src/detoo"},
		"/opt/skills": {Kind: Local, Dir: "/opt/skills"},
	} {
		got, err := ParseSource(in)
		if err != nil || got != want {
			t.Errorf("ParseSource(%q) = %+v, %v; want %+v", in, got, err, want)
		}
	}
	for _, in := range []string{"detoo", "a/b/c", "skills"} {
		if _, err := ParseSource(in); err == nil {
			t.Errorf("ParseSource(%q) succeeded, want an error", in)
		}
	}
}

// Every skill is installed in .agents/skills, linked for Claude Code and for
// each agent that seems installed and looks elsewhere, but for Kiro CLI and
// Junie, and recorded.
func TestInstallEverySkill(t *testing.T) {
	inst := newMem()
	inst.detected = []string{"cursor", "windsurf", "kiro-cli", "junie"}
	f := &memFetcher{offers: map[Source][]Skill{detoo: {skill("check-adrs", "1"), skill("create-adr", "1")}}}
	var h heard
	installed, err := Installer{Fetcher: f, Installations: Installations{Project: inst}}.Handle(context.Background(), Install{Source: detoo}, &h)
	if err != nil {
		t.Fatal(err)
	}
	want := heard{
		"fetching mustardcustarddev/detoo",
		"fetched 2",
		"installing into .agents/skills, linking claude-code=.claude/skills windsurf=.windsurf/skills",
		"installed check-adrs",
		"installed create-adr",
	}
	if !reflect.DeepEqual(h, want) {
		t.Errorf("heard\n%q\nwant\n%q", h, want)
	}
	if got := installed.Targets.Readers; len(got) != 1 || got[0].Name != "cursor" {
		t.Errorf("readers = %+v, want Cursor", got)
	}
	if got := inst.links[".windsurf/skills/create-adr"]; got != ".agents/skills/create-adr" {
		t.Errorf("windsurf's link leads to %q", got)
	}
	if got := inst.locked["check-adrs"]; got.Source != detoo || got.Hash != "1" || got.Path != "skills/check-adrs/SKILL.md" {
		t.Errorf("check-adrs is recorded as %+v", got)
	}
}

// Naming skills installs only those, for only the agents named; naming one
// there is not installs nothing, and so does naming an agent there is not,
// before anything is fetched.
func TestInstallNamedSkillsForNamedAgents(t *testing.T) {
	inst := newMem()
	f := &memFetcher{offers: map[Source][]Skill{detoo: {skill("check-adrs", "1"), skill("create-adr", "1")}}}
	a := Installer{Fetcher: f, Installations: Installations{Project: inst}}
	ctx := context.Background()

	_, err := a.Handle(ctx, Install{Source: detoo, Skills: []string{"create-adr", "make-tea"}}, nil)
	if !errors.Is(err, ErrUnknown) || !strings.Contains(err.Error(), "make-tea") || !strings.Contains(err.Error(), "check-adrs, create-adr") {
		t.Errorf("installing make-tea = %v, want an error naming it and what there is", err)
	}
	_, err = a.Handle(ctx, Install{Source: detoo, Agents: []string{"vim"}}, nil)
	if !errors.Is(err, ErrUnknown) || len(f.fetched) != 1 {
		t.Errorf("installing for vim = %v, after %d fetches; want an error before fetching", err, len(f.fetched))
	}
	if len(inst.dirs) > 0 {
		t.Fatalf("installed %v", inst.dirs)
	}

	installed, err := a.Handle(ctx, Install{Source: detoo, Skills: []string{"create-adr"}, Agents: []string{"roo", "codex", "junie"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(installed.Installed, []string{"create-adr"}) {
		t.Errorf("installed %v", installed.Installed)
	}
	if got := slices.Sorted(maps.Keys(inst.links)); !reflect.DeepEqual(got, []string{".junie/skills/create-adr", ".roo/skills/create-adr"}) {
		t.Errorf("links %v, want only Junie's and Roo's", got)
	}
}

// Update installs again only what changed, fetching each source once,
// reports what its source no longer offers, and goes on past a source that
// cannot be fetched.  A copy made where a link could not be is replaced.
func TestUpdate(t *testing.T) {
	inst := newMem()
	mine := Source{Kind: Local, Dir: "/home/ann/skills"}
	gone := Source{Kind: GitHub, Repo: "acme/gone"}
	for name, src := range map[string]Source{"check-adrs": detoo, "create-adr": detoo, "review-adrs": detoo, "tidy": mine, "lost": gone} {
		inst.locked[name] = Locked{Source: src, Path: "skills/" + name + "/SKILL.md", Hash: "1"}
		inst.dirs[path.Join(Canonical, name)] = "1"
	}
	inst.dirs[".claude/skills/create-adr"] = "1"
	inst.locked["npx-only"] = Locked{Source: Source{Kind: Other, Type: "well-known", Name: "example.com"}}
	f := &memFetcher{offers: map[Source][]Skill{
		detoo: {skill("check-adrs", "1"), skill("create-adr", "2")},
		mine:  {skill("tidy", "1")},
	}}
	var h heard
	updated, err := Updater{Fetcher: f, Installations: Installations{Project: inst}}.Handle(context.Background(), Update{}, &h)
	if err != nil {
		t.Fatal(err)
	}
	want := heard{
		"fetching mustardcustarddev/detoo",
		"fetched 2",
		"up to date check-adrs",
		"installed create-adr",
		"gone review-adrs",
		"fetching acme/gone",
		"failed lost",
		"failed npx-only",
		"fetching /home/ann/skills",
		"fetched 1",
		"up to date tidy",
	}
	if !reflect.DeepEqual(h, want) {
		t.Errorf("heard\n%q\nwant\n%q", h, want)
	}
	if !reflect.DeepEqual(updated.Updated, []string{"create-adr"}) || !reflect.DeepEqual(updated.Failed, []string{"lost", "npx-only"}) {
		t.Errorf("updated %+v", updated)
	}
	if inst.dirs[".agents/skills/create-adr"] != "2" || inst.dirs[".claude/skills/create-adr"] != "2" || inst.locked["create-adr"].Hash != "2" {
		t.Errorf("create-adr was not updated everywhere: %v, %+v", inst.dirs, inst.locked["create-adr"])
	}
	if _, ok := inst.dirs[".agents/skills/review-adrs"]; !ok {
		t.Errorf("review-adrs, gone from its source, was removed")
	}
}

// A skill recorded but missing from .agents/skills is installed again, and
// one not recorded cannot be updated.
func TestUpdateRestoresAndNeedsARecord(t *testing.T) {
	inst := newMem()
	inst.locked["create-adr"] = Locked{Source: detoo, Path: "skills/create-adr/SKILL.md", Hash: "1"}
	f := &memFetcher{offers: map[Source][]Skill{detoo: {skill("create-adr", "1")}}}
	u := Updater{Fetcher: f, Installations: Installations{Project: inst}}
	if _, err := u.Handle(context.Background(), Update{Skills: []string{"create-adr", "make-tea"}}, nil); !errors.Is(err, ErrUnknown) {
		t.Errorf("updating make-tea = %v, want ErrUnknown", err)
	}
	updated, err := u.Handle(context.Background(), Update{Skills: []string{"create-adr"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(updated.Updated, []string{"create-adr"}) || inst.dirs[".agents/skills/create-adr"] != "1" {
		t.Errorf("updated %+v, dirs %v", updated, inst.dirs)
	}
}

// Remove removes the links first, then the skill, then its record, then each
// directory it leaves empty; a skill not installed is refused before anything
// is removed.
func TestRemove(t *testing.T) {
	inst := newMem()
	for _, name := range []string{"check-adrs", "create-adr"} {
		inst.dirs[".agents/skills/"+name] = "1"
		inst.links[".claude/skills/"+name] = ".agents/skills/" + name
		inst.locked[name] = Locked{Source: detoo}
	}
	inst.dirs[".roo/skills/create-adr"] = "1"
	inst.links[".kiro/skills/create-adr"] = "/elsewhere/create-adr"
	r := Remover{Installations: Installations{Project: inst}}
	ctx := context.Background()

	if _, err := r.Handle(ctx, Remove{Skills: []string{"create-adr", "make-tea"}}, nil); !errors.Is(err, ErrUnknown) {
		t.Errorf("removing make-tea = %v, want ErrUnknown", err)
	}
	if len(inst.locked) != 2 {
		t.Fatalf("removed something: %v", inst.locked)
	}

	var h heard
	removed, err := r.Handle(ctx, Remove{Skills: []string{"create-adr", "create-adr"}}, &h)
	if err != nil {
		t.Fatal(err)
	}
	want := heard{"removed create-adr from .claude/skills/create-adr .roo/skills/create-adr .agents/skills/create-adr"}
	if !reflect.DeepEqual(h, want) || !reflect.DeepEqual(removed.Removed, []string{"create-adr"}) {
		t.Errorf("heard %q, removed %v", h, removed.Removed)
	}
	if !reflect.DeepEqual(removed.Emptied, []string{".roo/skills"}) {
		t.Errorf("emptied %v, want only .roo/skills", removed.Emptied)
	}
	if _, ok := inst.links[".kiro/skills/create-adr"]; !ok {
		t.Errorf("a link to elsewhere was removed")
	}

	removed, err = r.Handle(ctx, Remove{All: true}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(removed.Removed, []string{"check-adrs"}) || len(inst.locked) != 0 || len(inst.dirs) != 0 {
		t.Errorf("removed %v, leaving %v and %v", removed.Removed, inst.locked, inst.dirs)
	}
	if want := []string{".claude/skills", ".agents/skills"}; !reflect.DeepEqual(removed.Emptied, want) {
		t.Errorf("emptied %v, want %v", removed.Emptied, want)
	}
}
