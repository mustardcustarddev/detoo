package skills

import (
	"context"
	"fmt"
	"maps"
	"path"
	"slices"
)

// Remove asks for installed skills to be removed.
type Remove struct {
	// Skills are the names of the skills to remove.
	Skills []string

	// All removes every skill the lock file records, instead of Skills.
	All bool

	// Global removes skills installed for the user.
	Global bool
}

// Removed is what a Remove removed.
type Removed struct {
	// LockFile recorded the skills.
	LockFile string

	// Removed are the names of the skills removed, in order.
	Removed []string

	// Emptied are the skills directories removed because removing the
	// skills left them empty, such as .claude/skills.
	Emptied []string
}

// Remover removes skills from one of Installations.
type Remover struct {
	Installations Installations
}

// Handle removes the skills cmd asks for from .agents/skills and from the
// directory of every agent that has a link to it or a copy of it, and from the
// lock file, telling progress of each.  Then it removes each of those
// directories that is left empty.  It removes nothing unless every skill
// named is installed.
func (r Remover) Handle(ctx context.Context, cmd Remove, progress Progress) (Removed, error) {
	progress = orSilent(progress)
	inst := r.Installations.pick(cmd.Global)
	removed := Removed{LockFile: inst.LockFile()}
	locked, err := inst.Locked(ctx)
	if err != nil {
		return removed, err
	}
	var names []string
	for _, name := range cmd.Skills {
		if !slices.Contains(names, name) {
			names = append(names, name)
		}
	}
	if cmd.All {
		names = slices.Sorted(maps.Keys(locked))
	}

	// Every skill is found before any is removed.
	found := map[string][]string{}
	for _, name := range names {
		paths, err := installed(ctx, inst, name)
		if err != nil {
			return removed, err
		}
		if _, ok := locked[name]; !ok && len(paths) == 0 {
			return removed, fmt.Errorf("%w skill %s: it is not installed in %s", ErrUnknown, name, inst.Canonical())
		}
		found[name] = paths
	}

	var dirs []string
	for _, name := range names {
		for _, p := range found[name] {
			if err := inst.Remove(ctx, p); err != nil {
				return removed, fmt.Errorf("removing %s: %w", name, err)
			}
			if dir := path.Dir(p); !slices.Contains(dirs, dir) {
				dirs = append(dirs, dir)
			}
		}
		if _, ok := locked[name]; ok {
			if err := inst.Unlock(ctx, name); err != nil {
				return removed, err
			}
		}
		removed.Removed = append(removed.Removed, name)
		progress.Removed(name, found[name])
	}

	for _, dir := range dirs {
		emptied, err := inst.RemoveIfEmpty(ctx, dir)
		if err != nil {
			return removed, err
		}
		if emptied {
			removed.Emptied = append(removed.Emptied, dir)
		}
	}
	return removed, nil
}

// installed returns where the skill called name is installed: in
// .agents/skills, and linked or copied into the directories of the agents.
// The links come first, so that none is left pointing at nothing.
func installed(ctx context.Context, inst Installation, name string) ([]string, error) {
	links, copied, err := copies(ctx, inst, name)
	if err != nil {
		return nil, err
	}
	paths := append(links, copied...)
	dir := path.Join(inst.Canonical(), name)
	kind, err := inst.Entry(ctx, dir, dir)
	if err != nil {
		return nil, err
	}
	if kind != Missing {
		paths = append(paths, dir)
	}
	return paths, nil
}
