// Package skillsource reads the skills a source offers: a repository on
// GitHub, downloaded as a tarball, or a directory on this machine.
//
// This is a driven adapter, serving the skills slice.  A source's skills are
// the directories under its skills/ directory that hold a SKILL.md, or, when
// it has no skills/ directory, its own directories that do.
package skillsource

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/mustardcustard/detoo/internal/features/agents/skills"
)

// Source fetches skills over HTTP from GitHub's API, and from directories.
type Source struct {
	// API is the address of GitHub's API: https://api.github.com.
	API string

	Client *http.Client

	// workdir is where a relative directory is found from, and home where
	// ~ is.
	workdir, home string
	getenv        func(string) string
}

// New returns a Source that finds a relative directory from workdir, and ~ in
// home, and reads the token it gives GitHub, if any, from $GITHUB_TOKEN or
// $GH_TOKEN with getenv.
func New(workdir, home string, getenv func(string) string) *Source {
	return &Source{API: "https://api.github.com", Client: http.DefaultClient, workdir: workdir, home: home, getenv: getenv}
}

// Fetch returns the skills in src, in order of name.
func (s *Source) Fetch(ctx context.Context, src skills.Source) ([]skills.Skill, error) {
	var files map[string]file
	var err error
	switch src.Kind {
	case skills.GitHub:
		files, err = s.download(ctx, src)
	case skills.Local:
		files, err = s.read(src.Dir)
	default:
		return nil, fmt.Errorf("cannot fetch from %s, a %s source", src, src.Type)
	}
	if err != nil {
		return nil, err
	}
	return discover(files)
}

// file is a file of a source, by its path from the source's root.
type file struct {
	executable bool
	content    []byte
}

// download fetches the repository's tarball, keeping the files of its skills.
func (s *Source) download(ctx context.Context, src skills.Source) (map[string]file, error) {
	url := s.API + "/repos/" + src.Repo + "/tarball"
	if src.Ref != "" {
		url += "/" + src.Ref
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "detoo")
	token := s.getenv("GITHUB_TOKEN")
	if token == "" {
		token = s.getenv("GH_TOKEN")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("downloading %s from GitHub: %w", src, err)
	}
	defer resp.Body.Close()
	switch {
	case resp.StatusCode == http.StatusNotFound && src.Ref != "":
		return nil, fmt.Errorf("%s was not found on GitHub: check the repository and %s exist, or, if it is private, set GITHUB_TOKEN", src.Repo, src.Ref)
	case resp.StatusCode == http.StatusNotFound:
		return nil, fmt.Errorf("%s was not found on GitHub: check its name, or, if it is private, set GITHUB_TOKEN", src.Repo)
	case resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests:
		return nil, fmt.Errorf("GitHub refused to send %s (%s): it may be limiting requests, which setting GITHUB_TOKEN raises", src, resp.Status)
	case resp.StatusCode != http.StatusOK:
		return nil, fmt.Errorf("downloading %s from GitHub: %s", src, resp.Status)
	}

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading %s from GitHub: %w", src, err)
	}
	files, err := untar(gz)
	if err != nil {
		return nil, fmt.Errorf("reading %s from GitHub: %w", src, err)
	}
	return files, nil
}

// untar reads the regular files of a tarball whose every path is under one
// top directory, as GitHub's are, by their paths beneath it.  Links are left
// out, so that none can lead out of a skill.
func untar(r io.Reader) (map[string]file, error) {
	files := map[string]file{}
	tr := tar.NewReader(r)
	for {
		h, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return files, nil
		}
		if err != nil {
			return nil, err
		}
		if h.Typeflag != tar.TypeReg {
			continue
		}
		_, name, ok := strings.Cut(h.Name, "/")
		if !ok || !keep(name) {
			continue
		}
		content, err := io.ReadAll(tr)
		if err != nil {
			return nil, err
		}
		files[name] = file{executable: h.Mode&0o111 != 0, content: content}
	}
}

// keep reports whether the file name, from a source's root, could be part of
// a skill: it is in a directory, safely, and not in one of git's or npm's.
func keep(name string) bool {
	if !fs.ValidPath(name) || !strings.Contains(name, "/") {
		return false
	}
	for part := range strings.SplitSeq(name, "/") {
		if part == ".git" || part == "node_modules" {
			return false
		}
	}
	return true
}

// read reads the files of a directory on this machine, given as the user gave
// it.
func (s *Source) read(dir string) (map[string]file, error) {
	root := dir
	if rest, ok := strings.CutPrefix(dir, "~"); ok && (rest == "" || strings.HasPrefix(rest, "/")) {
		root = filepath.Join(s.home, rest)
	} else if !filepath.IsAbs(dir) {
		root = filepath.Join(s.workdir, dir)
	}
	if info, err := os.Stat(root); errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("cannot read skills from %s: there is no such directory", dir)
	} else if err != nil {
		return nil, fmt.Errorf("cannot read skills from %s: %w", dir, err)
	} else if !info.IsDir() {
		return nil, fmt.Errorf("cannot read skills from %s: it is not a directory", dir)
	}

	files := map[string]file{}
	err := filepath.WalkDir(root, func(p string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() && (entry.Name() == ".git" || entry.Name() == "node_modules") {
			return fs.SkipDir
		}
		if !entry.Type().IsRegular() {
			return nil
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		name := filepath.ToSlash(rel)
		if !keep(name) {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		content, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		files[name] = file{executable: info.Mode()&0o111 != 0, content: content}
		return nil
	})
	return files, err
}

// validName matches the name of a skill that is safe to make a directory of.
var validName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

// discover finds the skills among a source's files: the directories under
// skills/ that hold a SKILL.md, or, when there are none, the source's own
// directories that do.
func discover(files map[string]file) ([]skills.Skill, error) {
	found, err := skillsIn(files, "skills")
	if err != nil || len(found) > 0 {
		return found, err
	}
	return skillsIn(files, "")
}

func skillsIn(files map[string]file, parent string) ([]skills.Skill, error) {
	var found []skills.Skill
	byName := map[string]string{}
	for name := range files {
		// A skill's directory is directly under parent, or the source's root.
		dir, ok := strings.CutSuffix(name, "/SKILL.md")
		if !ok || path.Dir(dir) != path.Clean("./"+parent) {
			continue
		}
		meta := frontmatter(string(files[name].content))
		skill := skills.Skill{Name: meta["name"], Description: meta["description"], Path: name}
		if skill.Name == "" {
			skill.Name = path.Base(dir)
		}
		if !validName.MatchString(skill.Name) {
			return nil, fmt.Errorf("%s names its skill %q, which cannot be a directory's name", name, skill.Name)
		}
		if other, ok := byName[skill.Name]; ok {
			return nil, fmt.Errorf("%s and %s both name their skill %s", other, name, skill.Name)
		}
		byName[skill.Name] = name
		for p, f := range files {
			if rel, ok := strings.CutPrefix(p, dir+"/"); ok {
				skill.Files = append(skill.Files, skills.File{Path: rel, Executable: f.executable, Content: f.content})
			}
		}
		sort.Slice(skill.Files, func(i, j int) bool { return skill.Files[i].Path < skill.Files[j].Path })
		found = append(found, skill)
	}
	sort.Slice(found, func(i, j int) bool { return found[i].Name < found[j].Name })
	return found, nil
}

// frontmatter reads the simple YAML fields between the --- lines at the top of
// a SKILL.md: name: value, with the value quoted or not, or folded over the
// indented lines that follow a > or |.
func frontmatter(text string) map[string]string {
	fields := map[string]string{}
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return fields
	}
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "---" {
			break
		}
		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		value = strings.TrimSpace(value)
		if strings.HasPrefix(value, ">") || strings.HasPrefix(value, "|") {
			var folded []string
			for i+1 < len(lines) && (strings.HasPrefix(lines[i+1], " ") || strings.HasPrefix(lines[i+1], "\t")) {
				i++
				folded = append(folded, strings.TrimSpace(lines[i]))
			}
			value = strings.Join(folded, " ")
		} else if len(value) >= 2 && (value[0] == '"' || value[0] == '\'') && value[len(value)-1] == value[0] {
			value = value[1 : len(value)-1]
		}
		fields[strings.TrimSpace(key)] = value
	}
	return fields
}
