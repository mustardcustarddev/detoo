package skills

import (
	"context"
	"fmt"
	"maps"
	"path"
	"slices"
)

// Update asks for installed skills to be brought up to date with their
// sources.
type Update struct {
	// Skills are the names of the skills to update, or none for every skill
	// the lock file records.
	Skills []string

	// Global updates the skills installed for the user.
	Global bool
}

// Updated is what an Update found of each skill.
type Updated struct {
	// LockFile records the skills.
	LockFile string

	// Updated were installed again from their sources, and UpToDate were the
	// same as their sources.
	Updated, UpToDate []string

	// Gone are no longer offered by their sources, and are left as they are.
	Gone []string

	// Failed could not be checked or updated.
	Failed []string
}

// Updater updates skills in one of Installations from their sources, fetched
// by Fetcher.
type Updater struct {
	Fetcher       Fetcher
	Installations Installations
}

// Handle updates the skills cmd asks for, fetching each source once, and
// telling progress of each skill.  A skill whose source cannot be fetched is
// reported as failed, and the others are updated still.
func (u Updater) Handle(ctx context.Context, cmd Update, progress Progress) (Updated, error) {
	progress = orSilent(progress)
	inst := u.Installations.pick(cmd.Global)
	updated := Updated{LockFile: inst.LockFile()}
	locked, err := inst.Locked(ctx)
	if err != nil {
		return updated, err
	}
	names := cmd.Skills
	if len(names) == 0 {
		names = slices.Sorted(maps.Keys(locked))
	}
	for _, name := range names {
		if _, ok := locked[name]; !ok {
			return updated, fmt.Errorf("%w skill %s: %s does not record it", ErrUnknown, name, inst.LockFile())
		}
	}

	// Skills from the same source are updated together, from one fetch, in
	// the order the first of each was named.
	var sources []Source
	bySource := map[Source][]string{}
	for _, name := range names {
		src := locked[name].Source
		if _, ok := bySource[src]; !ok {
			sources = append(sources, src)
		}
		if !slices.Contains(bySource[src], name) {
			bySource[src] = append(bySource[src], name)
		}
	}
	for _, src := range sources {
		if err := u.updateFrom(ctx, inst, src, bySource[src], locked, progress, &updated); err != nil {
			return updated, err
		}
	}
	return updated, nil
}

// updateFrom updates the skills called names from src.  It returns an error
// only when the installation cannot be changed.
func (u Updater) updateFrom(ctx context.Context, inst Installation, src Source, names []string, locked map[string]Locked, progress Progress, updated *Updated) error {
	fail := func(err error) {
		for _, name := range names {
			updated.Failed = append(updated.Failed, name)
			progress.Failed(name, err)
		}
	}
	if src.Kind == Other {
		fail(fmt.Errorf("installed from a %s source, which only the skills CLI updates", src.Type))
		return nil
	}
	offered, err := fetch(ctx, u.Fetcher, src, progress)
	if err != nil {
		fail(err)
		return nil
	}

	for _, name := range names {
		i := slices.IndexFunc(offered, func(s Skill) bool { return s.Name == name })
		if i < 0 {
			updated.Gone = append(updated.Gone, name)
			progress.Gone(name, src)
			continue
		}
		skill := offered[i]
		dir := path.Join(inst.Canonical(), name)
		kind, err := inst.Entry(ctx, dir, dir)
		if err != nil {
			return err
		}
		hash := inst.Hash(skill.Files)
		if l := locked[name]; l.Hash == hash && l.Path == skill.Path && kind != Missing {
			updated.UpToDate = append(updated.UpToDate, name)
			progress.UpToDate(name)
			continue
		}

		// Links lead to the skill in .agents/skills, so are up to date once
		// it is; a copy made where a link could not be is replaced too.
		if err := inst.Install(ctx, dir, skill.Files); err != nil {
			return fmt.Errorf("updating %s: %w", name, err)
		}
		_, copied, err := copies(ctx, inst, name)
		if err != nil {
			return err
		}
		for _, c := range copied {
			if err := inst.Install(ctx, c, skill.Files); err != nil {
				return fmt.Errorf("updating %s: %w", name, err)
			}
		}
		if err := inst.Lock(ctx, name, Locked{Source: src, Path: skill.Path, Hash: hash}); err != nil {
			return err
		}
		updated.Updated = append(updated.Updated, name)
		progress.Installed(name)
	}
	return nil
}
