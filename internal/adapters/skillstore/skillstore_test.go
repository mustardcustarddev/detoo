package skillstore

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/mustardcustard/detoo/internal/features/agents/skills"
)

// demo is a skill with a reference, a script, and names whose order differs
// between byte order and localeCompare.
var demo = []skills.File{
	{Path: "SKILL.md", Content: []byte("---\nname: demo\n---\n# Demo\n")},
	{Path: "references/asking.md", Content: []byte("Ask one question at a time.\n")},
	{Path: "scripts/run.sh", Executable: true, Content: []byte("#!/bin/sh\necho hi\n")},
	{Path: "README.md", Content: []byte("Readme\n")},
}

// The hashes are those the skills CLI computed of demo, and git write-tree
// gave it, so that each tool sees the other's skills as up to date.
func TestHashesAreTheSkillsCLIs(t *testing.T) {
	if got, want := folderHash(demo), "c283db0fd59563bfc8fb6468ed1f61d2ac3f2df310e78189316f4d82881bab29"; got != want {
		t.Errorf("folderHash = %s, want %s", got, want)
	}
	if got, want := treeHash(demo), "fcb5f9b74bf69897c343b03fa9d95a8bcb80c827"; got != want {
		t.Errorf("treeHash = %s, want %s", got, want)
	}
}

// The order is the one Node's localeCompare sorted these names in.
func TestCollateSortsAsLocaleCompare(t *testing.T) {
	want := []string{"_x", "-x", "a b", "a_b", "a-b", "a-c", "a,b", "a;b", "a:b", "a!b", "a?b", "a.b", "a.c",
		"a.md", "A.md", "a'b", "a\"b", "a(b", "a[b", "a{b", "a@b", "a*b", "a/b", "a/bb", "a/c", "a\\b", "a&b", "a#b",
		"a%b", "a`b", "a^b", "a<b", "a=b", "a>b", "a|b", "a0", "a1.md", "a10.md", "a2.md", "a9", "aa/b", "ab", "aB",
		"Ab", "b", "B", "readme.md", "README.md", "references/adr-tool.md", "references/asking.md", "SKILL.md", "x+",
		"x~", "x$", "z", "Z"}
	got := slices.Clone(want)
	slices.Reverse(got)
	slices.SortFunc(got, collate)
	if !slices.Equal(got, want) {
		t.Errorf("sorted\n%q\nwant\n%q", got, want)
	}
}

func newProject(t *testing.T) (*Store, string) {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	s := Project(root, filepath.Join(root, "home"), func(string) string { return "" })
	return s, root
}

// A project's lock file is written as the skills CLI writes it: its skills in
// order of name, each with its fields in the CLI's order, and what detoo does
// not know of kept.
func TestProjectLockKeepsWhatItDoesNotKnow(t *testing.T) {
	s, root := newProject(t)
	ctx := context.Background()
	write(t, filepath.Join(root, "skills-lock.json"), `{
  "version": 1,
  "skills": {
    "zeta": {
      "source": "acme/tools",
      "sourceType": "github",
      "skillPath": "skills/zeta/SKILL.md",
      "computedHash": "abc",
      "pluginName": "zeta"
    },
    "gamma": {
      "source": "https://git.example.com/gamma.git",
      "sourceUrl": "https://git.example.com/gamma.git",
      "sourceType": "git",
      "computedHash": "def"
    }
  }
}
`)
	src := skills.Source{Kind: skills.GitHub, Repo: "mustardcustarddev/detoo", Ref: "v1.0.0"}
	if err := s.Lock(ctx, "create-adr", skills.Locked{Source: src, Path: "skills/create-adr/SKILL.md", Hash: "123"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Lock(ctx, "zeta", skills.Locked{Source: skills.Source{Kind: skills.GitHub, Repo: "acme/tools"}, Path: "skills/zeta/SKILL.md", Hash: "456"}); err != nil {
		t.Fatal(err)
	}
	want := `{
  "version": 1,
  "skills": {
    "create-adr": {
      "source": "mustardcustarddev/detoo",
      "ref": "v1.0.0",
      "sourceType": "github",
      "skillPath": "skills/create-adr/SKILL.md",
      "computedHash": "123"
    },
    "gamma": {
      "source": "https://git.example.com/gamma.git",
      "sourceUrl": "https://git.example.com/gamma.git",
      "sourceType": "git",
      "computedHash": "def"
    },
    "zeta": {
      "source": "acme/tools",
      "sourceType": "github",
      "skillPath": "skills/zeta/SKILL.md",
      "computedHash": "456",
      "pluginName": "zeta"
    }
  }
}
`
	if got := read(t, filepath.Join(root, "skills-lock.json")); got != want {
		t.Errorf("skills-lock.json =\n%s\nwant\n%s", got, want)
	}

	locked, err := s.Locked(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got := locked["create-adr"]; got.Source != src || got.Hash != "123" || got.Path != "skills/create-adr/SKILL.md" {
		t.Errorf("create-adr = %+v", got)
	}
	if got := locked["gamma"].Source; got.Kind != skills.Other || got.Type != "git" {
		t.Errorf("gamma's source = %+v, want another kind", got)
	}

	if err := s.Unlock(ctx, "create-adr"); err != nil {
		t.Fatal(err)
	}
	if got := read(t, filepath.Join(root, "skills-lock.json")); strings.Contains(got, "create-adr") {
		t.Errorf("create-adr is still recorded:\n%s", got)
	}
}

// A local directory is recorded from the project's directory, wherever it was
// given from, and read back as an absolute path.
func TestProjectLockRecordsLocalDirectoriesRelatively(t *testing.T) {
	_, root := newProject(t)
	ctx := context.Background()
	mkdir(t, filepath.Join(root, "project", ".git"))
	mkdir(t, filepath.Join(root, "project", "src"))
	s := Project(filepath.Join(root, "project", "src"), filepath.Join(root, "home"), func(string) string { return "" })
	src := skills.Source{Kind: skills.Local, Dir: "../../detoo"}
	if err := s.Lock(ctx, "demo", skills.Locked{Source: src, Path: "skills/demo/SKILL.md", Hash: "1"}); err != nil {
		t.Fatal(err)
	}
	if got := read(t, filepath.Join(root, "project", "skills-lock.json")); !strings.Contains(got, `"source": "../detoo"`) {
		t.Errorf("skills-lock.json =\n%s\nwant the source ../detoo", got)
	}
	locked, err := s.Locked(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := locked["demo"].Source.Dir, filepath.ToSlash(filepath.Join(root, "detoo")); got != want {
		t.Errorf("source = %s, want %s", got, want)
	}

	// A home directory reached through a link is compared by its real path.
	mkdir(t, filepath.Join(root, "detoo"))
	if err := os.Symlink(root, filepath.Join(root, "linked")); err != nil {
		t.Fatal(err)
	}
	s = Project(filepath.Join(root, "project"), filepath.Join(root, "linked"), func(string) string { return "" })
	if err := s.Lock(ctx, "demo", skills.Locked{Source: skills.Source{Kind: skills.Local, Dir: "~/detoo"}}); err != nil {
		t.Fatal(err)
	}
	if got := read(t, filepath.Join(root, "project", "skills-lock.json")); !strings.Contains(got, `"source": "../detoo"`) {
		t.Errorf("skills-lock.json =\n%s\nwant the source ../detoo", got)
	}
}

// A lock file a newer skills CLI wrote is not overwritten.
func TestNewerLockFilesAreRefused(t *testing.T) {
	s, root := newProject(t)
	write(t, filepath.Join(root, "skills-lock.json"), `{"version": 2, "skills": {}}`)
	err := s.Lock(context.Background(), "demo", skills.Locked{Source: skills.Source{Kind: skills.GitHub, Repo: "a/b"}})
	if err == nil || !strings.Contains(err.Error(), "version 2") {
		t.Errorf("Lock = %v, want an error naming version 2", err)
	}
}

// The user's lock file keeps the fields the skills CLI keeps there, and
// records when each skill was installed and last updated.
func TestUserLock(t *testing.T) {
	home := t.TempDir()
	s := User(home, func(string) string { return "" })
	ctx := context.Background()
	write(t, filepath.Join(home, ".agents", ".skill-lock.json"), `{
  "version": 3,
  "skills": {
    "find-skills": {
      "source": "vercel-labs/skills",
      "sourceType": "github",
      "sourceUrl": "https://github.com/vercel-labs/skills.git",
      "skillPath": "skills/find-skills/SKILL.md",
      "skillFolderHash": "3013",
      "installedAt": "2026-04-14T12:05:25.895Z",
      "updatedAt": "2026-04-14T12:05:25.895Z"
    }
  },
  "dismissed": {
    "findSkillsPrompt": true
  },
  "lastSelectedAgents": [
    "amp"
  ]
}`)
	s.now = func() time.Time { return time.Date(2026, 10, 18, 9, 30, 0, 0, time.UTC) }
	src := skills.Source{Kind: skills.GitHub, Repo: "mustardcustarddev/detoo"}
	if err := s.Lock(ctx, "create-adr", skills.Locked{Source: src, Path: "skills/create-adr/SKILL.md", Hash: "9c14"}); err != nil {
		t.Fatal(err)
	}
	s.now = func() time.Time { return time.Date(2026, 10, 19, 9, 30, 0, 0, time.UTC) }
	if err := s.Lock(ctx, "create-adr", skills.Locked{Source: src, Path: "skills/create-adr/SKILL.md", Hash: "9c15"}); err != nil {
		t.Fatal(err)
	}
	want := `{
  "version": 3,
  "skills": {
    "find-skills": {
      "source": "vercel-labs/skills",
      "sourceType": "github",
      "sourceUrl": "https://github.com/vercel-labs/skills.git",
      "skillPath": "skills/find-skills/SKILL.md",
      "skillFolderHash": "3013",
      "installedAt": "2026-04-14T12:05:25.895Z",
      "updatedAt": "2026-04-14T12:05:25.895Z"
    },
    "create-adr": {
      "source": "mustardcustarddev/detoo",
      "sourceType": "github",
      "sourceUrl": "https://github.com/mustardcustarddev/detoo.git",
      "skillPath": "skills/create-adr/SKILL.md",
      "skillFolderHash": "9c15",
      "installedAt": "2026-10-18T09:30:00.000Z",
      "updatedAt": "2026-10-19T09:30:00.000Z"
    }
  },
  "dismissed": {
    "findSkillsPrompt": true
  },
  "lastSelectedAgents": [
    "amp"
  ]
}`
	if got := read(t, filepath.Join(home, ".agents", ".skill-lock.json")); got != want {
		t.Errorf(".skill-lock.json =\n%s\nwant\n%s", got, want)
	}
	locked, err := s.Locked(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got := locked["find-skills"]; got.Hash != "3013" || got.Source.Repo != "vercel-labs/skills" {
		t.Errorf("find-skills = %+v", got)
	}
}

// The user's lock file is under $XDG_STATE_HOME when it is set, and the
// agents' directories under theirs.
func TestUserDirectoriesFollowTheEnvironment(t *testing.T) {
	env := map[string]string{"XDG_STATE_HOME": "/state", "XDG_CONFIG_HOME": "/config", "CLAUDE_CONFIG_DIR": "/claude"}
	s := User("/home/ann", func(name string) string { return env[name] })
	if got, want := s.LockFile(), "/state/skills/.skill-lock.json"; got != want {
		t.Errorf("LockFile = %s, want %s", got, want)
	}
	if got, want := s.Canonical(), "/home/ann/.agents/skills"; got != want {
		t.Errorf("Canonical = %s, want %s", got, want)
	}
	for name, want := range map[string]string{
		"claude-code": "/claude/skills",
		"opencode":    "/config/opencode/skills",
		"codex":       "/home/ann/.agents/skills",
		"windsurf":    "/home/ann/.codeium/windsurf/skills",
	} {
		agent, _ := skills.AgentNamed(name)
		if got := s.SkillsDir(agent); got != want {
			t.Errorf("SkillsDir(%s) = %s, want %s", name, got, want)
		}
	}
}

// A skill is installed in .agents/skills, and linked into another agent's
// directory relatively; a link or copy already there is replaced, and a link
// that already leads there is kept.
func TestInstallLinkAndRemove(t *testing.T) {
	s, root := newProject(t)
	ctx := context.Background()
	write(t, filepath.Join(root, ".agents/skills/demo/stale.md"), "old\n")
	if err := s.Install(ctx, ".agents/skills/demo", demo); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, ".agents/skills/demo/stale.md")); !os.IsNotExist(err) {
		t.Errorf("a file of the old skill was left: %v", err)
	}
	if info, err := os.Stat(filepath.Join(root, ".agents/skills/demo/scripts/run.sh")); err != nil || info.Mode().Perm() != 0o755 {
		t.Errorf("run.sh = %v, %v; want it executable", info, err)
	}

	write(t, filepath.Join(root, ".claude/skills/demo/SKILL.md"), "a copy\n")
	if got, _ := s.Entry(ctx, ".claude/skills/demo", ".agents/skills/demo"); got != skills.Copied {
		t.Errorf("Entry of a copy = %v, want Copied", got)
	}
	for range 2 {
		if err := s.Link(ctx, ".claude/skills/demo", ".agents/skills/demo"); err != nil {
			t.Fatal(err)
		}
	}
	if dest, err := os.Readlink(filepath.Join(root, ".claude/skills/demo")); err != nil || dest != "../../.agents/skills/demo" {
		t.Errorf("link = %q, %v; want ../../.agents/skills/demo", dest, err)
	}
	if got, _ := s.Entry(ctx, ".claude/skills/demo", ".agents/skills/demo"); got != skills.Linked {
		t.Errorf("Entry of a link = %v, want Linked", got)
	}
	if got, _ := s.Entry(ctx, ".roo/skills/demo", ".agents/skills/demo"); got != skills.Missing {
		t.Errorf("Entry of nothing = %v, want Missing", got)
	}

	if err := s.Remove(ctx, ".agents/skills/demo"); err != nil {
		t.Fatal(err)
	}
	// A link left leading nowhere is still the skill's.
	if got, _ := s.Entry(ctx, ".claude/skills/demo", ".agents/skills/demo"); got != skills.Linked {
		t.Errorf("Entry of a broken link = %v, want Linked", got)
	}
	if err := s.Remove(ctx, ".claude/skills/demo"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(filepath.Join(root, ".claude/skills/demo")); !os.IsNotExist(err) {
		t.Errorf("the link was left: %v", err)
	}
}

// An agent whose directory is itself a link to .agents/skills already has the
// skill, and linking it does not replace the skill with a link to itself.
func TestLinkingThroughALinkedDirectory(t *testing.T) {
	s, root := newProject(t)
	ctx := context.Background()
	if err := s.Install(ctx, ".agents/skills/demo", demo); err != nil {
		t.Fatal(err)
	}
	mkdir(t, filepath.Join(root, ".claude"))
	if err := os.Symlink("../.agents/skills", filepath.Join(root, ".claude/skills")); err != nil {
		t.Fatal(err)
	}
	if err := s.Link(ctx, ".claude/skills/demo", ".agents/skills/demo"); err != nil {
		t.Fatal(err)
	}
	if got := read(t, filepath.Join(root, ".agents/skills/demo/README.md")); got != "Readme\n" {
		t.Errorf("README.md = %q", got)
	}
}

// The project is the nearest directory with a lock file, .agents or .git,
// but never the home directory, where the user's skills are.
func TestRoot(t *testing.T) {
	_, root := newProject(t)
	home := filepath.Join(root, "home")
	mkdir(t, filepath.Join(home, ".agents"))
	mkdir(t, filepath.Join(home, "work", "app", ".git"))
	mkdir(t, filepath.Join(home, "work", "app", "src"))
	mkdir(t, filepath.Join(home, "notes"))
	for dir, want := range map[string]string{
		"work/app/src": "work/app",
		"notes":        "notes",
		".":            ".",
	} {
		s := Project(filepath.Join(home, dir), home, func(string) string { return "" })
		got, err := s.Root(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if want := filepath.Join(home, want); got != want {
			t.Errorf("Root from %s = %s, want %s", dir, got, want)
		}
	}
}

func write(t *testing.T, name, content string) {
	t.Helper()
	mkdir(t, filepath.Dir(name))
	if err := os.WriteFile(name, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func read(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func mkdir(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
}
