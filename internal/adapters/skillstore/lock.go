package skillstore

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mustardcustard/detoo/internal/features/agents/skills"
)

// The lock files, by name, and the version of each that is written.
const (
	projectLock = "skills-lock.json"
	userLock    = ".skill-lock.json"

	projectVersion = 1
	userVersion    = 3
)

// Locked returns the skills the lock file records, by name.  A lock file
// that does not exist records none.
func (s *Store) Locked(ctx context.Context) (map[string]skills.Locked, error) {
	lock, err := s.read(ctx)
	if err != nil {
		return nil, err
	}
	locked := map[string]skills.Locked{}
	for _, f := range lock.skills {
		var e entry
		if err := json.Unmarshal(f.value, &e); err != nil {
			return nil, fmt.Errorf("%s: the skill %s: %w", s.LockFile(), f.key, err)
		}
		l := skills.Locked{Path: e.SkillPath, Hash: e.ComputedHash}
		if s.global {
			l.Hash = e.SkillFolderHash
		}
		switch e.SourceType {
		case "github":
			l.Source = skills.Source{Kind: skills.GitHub, Repo: e.Source, Ref: e.Ref}
		case "local":
			l.Source = skills.Source{Kind: skills.Local, Dir: s.fromLock(ctx, e.Source)}
		default:
			l.Source = skills.Source{Kind: skills.Other, Type: e.SourceType, Name: e.Source}
		}
		locked[f.key] = l
	}
	return locked, nil
}

// entry is what both lock files record of a skill, that detoo reads.
type entry struct {
	Source          string `json:"source"`
	SourceType      string `json:"sourceType"`
	Ref             string `json:"ref"`
	SkillPath       string `json:"skillPath"`
	ComputedHash    string `json:"computedHash"`
	SkillFolderHash string `json:"skillFolderHash"`
}

// Lock records the skill called name in the lock file, keeping what the file
// records of it that detoo does not know.
func (s *Store) Lock(ctx context.Context, name string, l skills.Locked) error {
	lock, err := s.read(ctx)
	if err != nil {
		return err
	}
	var e object
	if i := lock.skills.index(name); i >= 0 {
		if e, err = parseObject(lock.skills[i].value); err != nil {
			return fmt.Errorf("%s: the skill %s: %w", s.LockFile(), name, err)
		}
	}

	source, sourceType := l.Source.Name, l.Source.Type
	switch l.Source.Kind {
	case skills.GitHub:
		source, sourceType = l.Source.Repo, "github"
	case skills.Local:
		source, sourceType = s.toLock(ctx, l.Source.Dir), "local"
	}
	// A field already recorded keeps its place, and a new one goes in the
	// place the skills CLI gives it: the user's lock file records the source
	// type before the ref, and a project's after it.
	e.set("source", source)
	if s.global {
		e.set("sourceType", sourceType)
		if l.Source.Kind == skills.GitHub {
			e.set("sourceUrl", "https://github.com/"+l.Source.Repo+".git")
		}
	}
	if l.Source.Ref != "" {
		e.set("ref", l.Source.Ref)
	} else {
		e.remove("ref")
	}
	e.set("sourceType", sourceType)
	e.set("skillPath", l.Path)
	if s.global {
		e.set("skillFolderHash", l.Hash)
		now := s.now().UTC().Format("2006-01-02T15:04:05.000Z")
		if e.index("installedAt") < 0 {
			e.set("installedAt", now)
		}
		e.set("updatedAt", now)
	} else {
		e.set("computedHash", l.Hash)
	}

	value, err := e.MarshalJSON()
	if err != nil {
		return err
	}
	lock.skills.put(name, value)
	return s.write(ctx, lock)
}

// Unlock forgets the skill called name.
func (s *Store) Unlock(ctx context.Context, name string) error {
	lock, err := s.read(ctx)
	if err != nil {
		return err
	}
	i := lock.skills.index(name)
	if i < 0 {
		return nil
	}
	lock.skills = slices.Delete(lock.skills, i, i+1)
	return s.write(ctx, lock)
}

// lockFile is a lock file: its fields, in order, among them its skills.
type lockFile struct {
	fields object
	skills object
}

// lockPath returns where the lock file is on disk.
func (s *Store) lockPath(ctx context.Context) (string, error) {
	return s.path(ctx, s.LockFile())
}

func (s *Store) version() int {
	if s.global {
		return userVersion
	}
	return projectVersion
}

// read reads the lock file.  One that does not exist, or that an older
// version of the skills CLI wrote, which it would start afresh, is empty.  One
// a newer version wrote is refused, so that it is not overwritten with what
// detoo does not understand.
func (s *Store) read(ctx context.Context) (lockFile, error) {
	empty := lockFile{fields: object{
		{"version", json.RawMessage(fmt.Sprint(s.version()))},
		{"skills", json.RawMessage("{}")},
	}}
	if s.global {
		empty.fields = append(empty.fields, field{"dismissed", json.RawMessage("{}")})
	}
	file, err := s.lockPath(ctx)
	if err != nil {
		return lockFile{}, err
	}
	data, err := os.ReadFile(file)
	if errors.Is(err, fs.ErrNotExist) {
		return empty, nil
	} else if err != nil {
		return lockFile{}, err
	}
	fields, err := parseObject(data)
	if err != nil {
		return lockFile{}, fmt.Errorf("%s: %w", s.LockFile(), err)
	}
	var version int
	if i := fields.index("version"); i >= 0 {
		if err := json.Unmarshal(fields[i].value, &version); err != nil {
			return lockFile{}, fmt.Errorf("%s: its version: %w", s.LockFile(), err)
		}
	}
	i := fields.index("skills")
	switch {
	case version > s.version():
		return lockFile{}, fmt.Errorf("%s is version %d, newer than the version %d detoo writes: update it with the skills CLI", s.LockFile(), version, s.version())
	case version < s.version() || i < 0:
		return empty, nil
	}
	installed, err := parseObject(fields[i].value)
	if err != nil {
		return lockFile{}, fmt.Errorf("%s: its skills: %w", s.LockFile(), err)
	}
	return lockFile{fields: fields, skills: installed}, nil
}

// write writes the lock file as the skills CLI does: indented by two spaces,
// with characters such as < as they are.  In a project, the skills are in
// order of name, and the file ends in a newline.
func (s *Store) write(ctx context.Context, lock lockFile) error {
	if !s.global {
		slices.SortFunc(lock.skills, func(a, b field) int { return strings.Compare(a.key, b.key) })
	}
	installed, err := lock.skills.MarshalJSON()
	if err != nil {
		return err
	}
	lock.fields.put("skills", installed)
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(lock.fields); err != nil {
		return err
	}
	data := buf.Bytes()
	if s.global {
		data = bytes.TrimSuffix(data, []byte("\n"))
	}
	file, err := s.lockPath(ctx)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		return err
	}
	return os.WriteFile(file, data, 0o644)
}

// toLock returns how the lock file records a local directory: from the
// project's directory, as ./dir or ../dir, so that the project can move, or
// as an absolute path for the user.
func (s *Store) toLock(ctx context.Context, dir string) string {
	abs := realPath(s.local(dir))
	if s.global {
		return filepath.ToSlash(abs)
	}
	root, err := s.Root(ctx)
	if err != nil {
		return filepath.ToSlash(abs)
	}
	rel, err := filepath.Rel(root, abs)
	switch {
	case err != nil:
		return filepath.ToSlash(abs)
	case rel == ".":
		return "."
	case rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)):
		return filepath.ToSlash(rel)
	}
	return "./" + filepath.ToSlash(rel)
}

// fromLock returns the absolute path of a local directory the lock file
// records.
func (s *Store) fromLock(ctx context.Context, dir string) string {
	if filepath.IsAbs(filepath.FromSlash(dir)) {
		return dir
	}
	root, err := s.Root(ctx)
	if err != nil {
		return dir
	}
	return filepath.ToSlash(filepath.Join(root, filepath.FromSlash(dir)))
}

// local returns the absolute path of a directory given by the user: from the
// working directory, or under ~.
func (s *Store) local(dir string) string {
	if rest, ok := strings.CutPrefix(dir, "~"); ok && (rest == "" || strings.HasPrefix(rest, "/")) {
		return filepath.Join(s.home, filepath.FromSlash(rest))
	}
	if filepath.IsAbs(filepath.FromSlash(dir)) {
		return filepath.FromSlash(dir)
	}
	workdir := s.workdir
	if workdir == "" {
		workdir, _ = os.Getwd()
	}
	return filepath.Join(workdir, filepath.FromSlash(dir))
}

// object is a JSON object whose fields keep their order, so that rewriting a
// lock file changes only what it must.
type object []field

type field struct {
	key   string
	value json.RawMessage
}

func (o object) index(key string) int {
	return slices.IndexFunc(o, func(f field) bool { return f.key == key })
}

// put sets the field called key, where it is, or else last.
func (o *object) put(key string, value json.RawMessage) {
	if i := o.index(key); i >= 0 {
		(*o)[i].value = value
		return
	}
	*o = append(*o, field{key, value})
}

// set sets the field called key to a string.
func (o *object) set(key, value string) {
	data, _ := json.Marshal(value)
	o.put(key, data)
}

func (o *object) remove(key string) {
	if i := o.index(key); i >= 0 {
		*o = slices.Delete(*o, i, i+1)
	}
}

func (o object) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, f := range o {
		if i > 0 {
			buf.WriteByte(',')
		}
		key, err := json.Marshal(f.key)
		if err != nil {
			return nil, err
		}
		buf.Write(key)
		buf.WriteByte(':')
		buf.Write(f.value)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// parseObject reads a JSON object, keeping the order of its fields.
func parseObject(data []byte) (object, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	if t, err := dec.Token(); err != nil {
		return nil, err
	} else if t != json.Delim('{') {
		return nil, errors.New("not a JSON object")
	}
	var o object
	for dec.More() {
		t, err := dec.Token()
		if err != nil {
			return nil, err
		}
		key, _ := t.(string)
		var value json.RawMessage
		if err := dec.Decode(&value); err != nil {
			return nil, err
		}
		o = append(o, field{key, value})
	}
	return o, nil
}
