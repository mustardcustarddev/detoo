package skills

import (
	"context"
	"fmt"
	"path"
	"slices"
	"strings"
)

// Install asks for skills from a source to be installed.
type Install struct {
	Source Source

	// Skills are the names of the skills to install, or none for every skill
	// the source offers.
	Skills []string

	// Agents are the names of the agents to install for, or none for Claude
	// Code and each agent that seems to be installed.
	Agents []string

	// Global installs for the user, in every project, rather than in the
	// project.
	Global bool
}

// Installed is what an Install installed, and where.
type Installed struct {
	Targets Targets

	// Installed are the names of the skills installed, in order.
	Installed []string
}

// Installer installs skills, fetched by Fetcher, in one of Installations.
type Installer struct {
	Fetcher       Fetcher
	Installations Installations
}

// Available returns the skills src offers, installing nothing.
func (a Installer) Available(ctx context.Context, src Source, progress Progress) ([]Skill, error) {
	return fetch(ctx, a.Fetcher, src, orSilent(progress))
}

// Handle installs the skills cmd asks for, telling progress of each step.
// It installs nothing unless every skill and agent named is one there is.
func (a Installer) Handle(ctx context.Context, cmd Install, progress Progress) (Installed, error) {
	progress = orSilent(progress)
	inst := a.Installations.pick(cmd.Global)
	// The agents are checked before the source is fetched, which can be slow.
	t, err := targets(ctx, inst, cmd.Agents)
	if err != nil {
		return Installed{}, err
	}
	offered, err := fetch(ctx, a.Fetcher, cmd.Source, progress)
	if err != nil {
		return Installed{}, err
	}
	chosen, err := choose(offered, cmd.Skills, cmd.Source)
	if err != nil {
		return Installed{}, err
	}

	progress.Installing(t)
	installed := Installed{Targets: t}
	for _, skill := range chosen {
		if err := install(ctx, inst, t, skill, cmd.Source); err != nil {
			return installed, fmt.Errorf("installing %s: %w", skill.Name, err)
		}
		installed.Installed = append(installed.Installed, skill.Name)
		progress.Installed(skill.Name)
	}
	return installed, nil
}

func fetch(ctx context.Context, f Fetcher, src Source, progress Progress) ([]Skill, error) {
	progress.Fetching(src)
	offered, err := f.Fetch(ctx, src)
	if err != nil {
		return nil, err
	}
	if len(offered) == 0 {
		return nil, fmt.Errorf("%s offers no skills: none of its directories under skills/ holds a SKILL.md", src)
	}
	progress.Fetched(src, offered)
	return offered, nil
}

// choose returns the skills named, or every skill offered when none is.
func choose(offered []Skill, names []string, src Source) ([]Skill, error) {
	if len(names) == 0 {
		return offered, nil
	}
	var chosen, unknown []string
	for _, name := range names {
		switch {
		case slices.Contains(chosen, name):
		case slices.ContainsFunc(offered, func(s Skill) bool { return s.Name == name }):
			chosen = append(chosen, name)
		default:
			unknown = append(unknown, name)
		}
	}
	if len(unknown) > 0 {
		var all []string
		for _, s := range offered {
			all = append(all, s.Name)
		}
		return nil, fmt.Errorf("%w skill %s: %s offers %s", ErrUnknown, strings.Join(unknown, ", "), src, strings.Join(all, ", "))
	}
	var skills []Skill
	for _, s := range offered {
		if slices.Contains(chosen, s.Name) {
			skills = append(skills, s)
		}
	}
	return skills, nil
}

// install installs skill in .agents/skills, links it into the directory of
// each agent in t that looks elsewhere, and records it in the lock file.
func install(ctx context.Context, inst Installation, t Targets, skill Skill, src Source) error {
	dir := path.Join(t.Dir, skill.Name)
	if err := inst.Install(ctx, dir, skill.Files); err != nil {
		return err
	}
	for _, l := range t.Links {
		if err := inst.Link(ctx, path.Join(l.Dir, skill.Name), dir); err != nil {
			return err
		}
	}
	return inst.Lock(ctx, skill.Name, Locked{Source: src, Path: skill.Path, Hash: inst.Hash(skill.Files)})
}
