// Package skills installs agent skills, such as those in detoo's own
// repository on GitHub, into a project or for the user, keeps them up to
// date, and removes them.
//
// It works as the skills CLI (npx skills) does, and shares its files, so that
// either can update or remove what the other installed.  A skill is installed
// once, in .agents/skills, which many agents read, and linked into the
// directory of each agent that looks elsewhere, such as .claude/skills.  What
// was installed, and from where, is recorded in a lock file: skills-lock.json
// in a project.
//
// Installer installs skills, Updater brings them up to date with their source,
// and Remover removes them.  Each tells Progress what it is doing as it does
// it, since fetching a source can take a while.
package skills

import (
	"context"
	"errors"
	"fmt"
	"path"
	"regexp"
	"slices"
	"strings"
)

// File is a file of a skill.
type File struct {
	// Path is where the file is in the skill's directory, slash separated.
	Path string

	// Executable reports whether the file may be run.
	Executable bool

	Content []byte
}

// Skill is a skill a source offers: a directory holding a SKILL.md and the
// files it refers to.
type Skill struct {
	// Name is the name its SKILL.md gives it, which it is installed under.
	Name string

	Description string

	// Path is where its SKILL.md is in the source, such as
	// skills/create-adr/SKILL.md.
	Path string

	Files []File
}

// SourceKind says where a source is.
type SourceKind int

const (
	// GitHub is a repository on GitHub.
	GitHub SourceKind = iota

	// Local is a directory on this machine.
	Local

	// Other is a kind of source that only the skills CLI installs from, such
	// as a git server other than GitHub, found in a lock file.
	Other
)

// Source is where skills are installed from.
type Source struct {
	Kind SourceKind

	// Repo is the repository on GitHub, such as mustardcustarddev/detoo.
	Repo string

	// Ref is the branch, tag or commit of Repo, or empty for its default
	// branch.
	Ref string

	// Dir is the directory of a Local source, as given.
	Dir string

	// Type and Name are what a lock file records of an Other source.
	Type, Name string
}

// DefaultSource is detoo's own repository, whose skills drive its tools.
const DefaultSource = "mustardcustarddev/detoo"

// repo matches a repository on GitHub given as OWNER/REPO.
var repo = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]*/[A-Za-z0-9._-]+$`)

// ParseSource reads a source as the user gives it: a repository on GitHub, as
// OWNER/REPO or its https://github.com/ address, followed by #REF for a
// branch, tag or commit, or a directory, as a path that begins with ., / or ~.
func ParseSource(s string) (Source, error) {
	if s == "." || s == "~" || strings.HasPrefix(s, "./") || strings.HasPrefix(s, "../") ||
		strings.HasPrefix(s, "/") || strings.HasPrefix(s, "~/") {
		return Source{Kind: Local, Dir: s}, nil
	}
	name, ref, _ := strings.Cut(s, "#")
	name = strings.TrimPrefix(name, "https://github.com/")
	name = strings.TrimSuffix(strings.TrimSuffix(name, "/"), ".git")
	if !repo.MatchString(name) {
		return Source{}, fmt.Errorf("%q is not a source: give a repository on GitHub, such as %s, or a directory, such as ./skills", s, DefaultSource)
	}
	return Source{Kind: GitHub, Repo: name, Ref: ref}, nil
}

// String is the source as the user would give it.
func (s Source) String() string {
	switch s.Kind {
	case GitHub:
		if s.Ref != "" {
			return s.Repo + "#" + s.Ref
		}
		return s.Repo
	case Local:
		return s.Dir
	}
	return s.Name
}

// Fetcher is the port that reads the skills a source offers.
type Fetcher interface {
	// Fetch returns the skills in src, each with its files, in order of name.
	Fetch(ctx context.Context, src Source) ([]Skill, error)
}

// Locked is what a lock file records of an installed skill.
type Locked struct {
	Source Source

	// Path is where its SKILL.md was in the source.
	Path string

	// Hash identifies the content of its files, as Installation.Hash computes
	// it.
	Hash string
}

// EntryKind is what is at a path where a skill could be.
type EntryKind int

const (
	// Missing is nothing.
	Missing EntryKind = iota

	// Linked is the skill installed in .agents/skills, through a link, or
	// because it is the same directory.
	Linked

	// Copied is a directory of its own: a copy of the skill made where a link
	// could not be.
	Copied

	// Foreign is something else, such as a file, or a link to elsewhere.
	Foreign
)

// Installation is the port to where skills are installed: a project, or the
// user's home directory.  The paths it takes and returns are slash
// separated: relative to the project's directory for a project, and absolute
// for the user.
type Installation interface {
	// Root returns the project's directory, or the user's home directory.
	Root(ctx context.Context) (string, error)

	// Canonical returns the directory each skill is installed in:
	// .agents/skills.
	Canonical() string

	// SkillsDir returns the directory agent looks for skills in.
	SkillsDir(agent Agent) string

	// Detected reports whether agent seems to be installed: whether its
	// configuration directory exists.
	Detected(ctx context.Context, agent Agent) (bool, error)

	// LockFile returns the lock file that records what is installed.
	LockFile() string

	// Hash identifies the content of a skill's files, as the lock file
	// records it.
	Hash(files []File) string

	// Install replaces whatever is at dir with a directory of files.
	Install(ctx context.Context, dir string, files []File) error

	// Link makes link a link to the directory target, replacing whatever is
	// there, or a copy of it where a link cannot be made.
	Link(ctx context.Context, link, target string) error

	// Entry reports what is at name, compared with the directory target.
	Entry(ctx context.Context, name, target string) (EntryKind, error)

	// Remove removes name, and everything in it.
	Remove(ctx context.Context, name string) error

	// RemoveIfEmpty removes dir if nothing is left in it, and reports whether
	// it did.
	RemoveIfEmpty(ctx context.Context, dir string) (bool, error)

	// Locked returns the skills the lock file records, by name.
	Locked(ctx context.Context) (map[string]Locked, error)

	// Lock records a skill in the lock file, and Unlock forgets one.
	Lock(ctx context.Context, name string, l Locked) error
	Unlock(ctx context.Context, name string) error
}

// Installations are the two places skills can be installed.
type Installations struct {
	Project, User Installation
}

func (i Installations) pick(global bool) Installation {
	if global {
		return i.User
	}
	return i.Project
}

// Progress hears of each step as it is taken, so that the user can follow
// what is being done.
type Progress interface {
	// Fetching is told before a source is read, and Fetched after.
	Fetching(src Source)
	Fetched(src Source, skills []Skill)

	// Installing is told where skills are about to be installed.
	Installing(t Targets)

	// Installed is told of a skill installed or updated.
	Installed(name string)

	// UpToDate is told of a skill that is the same as in its source.
	UpToDate(name string)

	// Gone is told of a skill its source no longer offers.
	Gone(name string, src Source)

	// Failed is told of a skill that could not be updated, and why.
	Failed(name string, err error)

	// Removed is told of a skill removed, from each path it was at.
	Removed(name string, paths []string)
}

// Targets are where skills are installed.
type Targets struct {
	// Root is the project's directory, or the user's home directory.
	Root string

	// Dir is the directory each skill is installed in.
	Dir string

	// Readers are the agents chosen that read Dir.
	Readers []Agent

	// Links are the agents chosen that look elsewhere, with the directory
	// each is given a link in.
	Links []Link

	// LockFile records what is installed.
	LockFile string
}

// Link is the directory an agent looks for skills in, other than
// .agents/skills.
type Link struct {
	Agent Agent
	Dir   string
}

// ErrUnknown is wrapped by the error returned when a skill or agent named
// is not one there is.
var ErrUnknown = errors.New("unknown")

// targets works out where skills go in inst for the agents named, or, when
// none is, for Claude Code and each agent that seems to be installed and is
// installed for automatically.
func targets(ctx context.Context, inst Installation, names []string) (Targets, error) {
	chosen, err := chooseAgents(ctx, inst, names)
	if err != nil {
		return Targets{}, err
	}
	root, err := inst.Root(ctx)
	if err != nil {
		return Targets{}, err
	}
	t := Targets{Root: root, Dir: inst.Canonical(), LockFile: inst.LockFile()}
	for _, agent := range chosen {
		dir := inst.SkillsDir(agent)
		switch {
		case dir == t.Dir:
			t.Readers = append(t.Readers, agent)
		case !slices.ContainsFunc(t.Links, func(l Link) bool { return l.Dir == dir }):
			t.Links = append(t.Links, Link{Agent: agent, Dir: dir})
		}
	}
	return t, nil
}

func chooseAgents(ctx context.Context, inst Installation, names []string) ([]Agent, error) {
	if len(names) > 0 {
		var chosen []Agent
		for _, name := range names {
			agent, ok := AgentNamed(name)
			if !ok {
				return nil, fmt.Errorf("%w agent %q: choose from %s", ErrUnknown, name, strings.Join(AgentNames(), ", "))
			}
			if !slices.Contains(chosen, agent) {
				chosen = append(chosen, agent)
			}
		}
		return chosen, nil
	}
	var chosen []Agent
	for _, agent := range Agents {
		if !agent.Automatic() {
			continue
		}
		found, err := inst.Detected(ctx, agent)
		if err != nil {
			return nil, err
		}
		if found || agent.Name == DefaultAgent {
			chosen = append(chosen, agent)
		}
	}
	return chosen, nil
}

// copies returns the directories, other than .agents/skills, that hold the
// skill called name, as a link to it or a copy of it: one in each agent's
// directory.
func copies(ctx context.Context, inst Installation, name string) (links, copied []string, err error) {
	target := path.Join(inst.Canonical(), name)
	seen := map[string]bool{inst.Canonical(): true}
	for _, agent := range Agents {
		dir := inst.SkillsDir(agent)
		if seen[dir] {
			continue
		}
		seen[dir] = true
		at := path.Join(dir, name)
		kind, err := inst.Entry(ctx, at, target)
		if err != nil {
			return nil, nil, err
		}
		switch kind {
		case Linked:
			links = append(links, at)
		case Copied:
			copied = append(copied, at)
		}
	}
	return links, copied, nil
}

// silent hears of nothing.
type silent struct{}

func (silent) Fetching(Source)          {}
func (silent) Fetched(Source, []Skill)  {}
func (silent) Installing(Targets)       {}
func (silent) Installed(string)         {}
func (silent) UpToDate(string)          {}
func (silent) Gone(string, Source)      {}
func (silent) Failed(string, error)     {}
func (silent) Removed(string, []string) {}

func orSilent(p Progress) Progress {
	if p == nil {
		return silent{}
	}
	return p
}
