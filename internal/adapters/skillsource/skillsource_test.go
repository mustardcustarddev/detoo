package skillsource

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mustardcustard/detoo/internal/features/agents/skills"
)

// tarball makes a tarball as GitHub's API sends one: every path under a
// directory named after the repository and commit.
func tarball(t *testing.T, entries ...*tar.Header) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for _, h := range entries {
		h.Name = "acme-tools-1a2b3c/" + h.Name
		content := []byte(h.Linkname)
		if h.Typeflag == tar.TypeReg {
			h.Linkname = ""
			h.Size = int64(len(content))
		}
		if err := tw.WriteHeader(h); err != nil {
			t.Fatal(err)
		}
		if h.Typeflag == tar.TypeReg {
			tw.Write(content)
		}
	}
	tw.Close()
	gz.Close()
	return buf.Bytes()
}

// regular is a regular file in a tarball; its content is carried in Linkname
// until tarball writes it.
func regular(name, content string, mode int64) *tar.Header {
	return &tar.Header{Typeflag: tar.TypeReg, Name: name, Linkname: content, Mode: mode}
}

// The skills are read from the tarball of the ref asked for, with the token
// given, leaving out links and anything under .git.
func TestFetchFromGitHub(t *testing.T) {
	body := tarball(t,
		&tar.Header{Typeflag: tar.TypeDir, Name: "skills/", Mode: 0o755},
		regular("README.md", "# Tools\n", 0o644),
		regular("skills/write-docs/SKILL.md", "---\nname: write-docs\ndescription: >\n  Writes the docs,\n  carefully.\n---\n# Write docs\n", 0o644),
		regular("skills/write-docs/scripts/check.sh", "#!/bin/sh\n", 0o755),
		&tar.Header{Typeflag: tar.TypeSymlink, Name: "skills/write-docs/secrets", Linkname: "/etc", Mode: 0o777},
		regular("skills/write-docs/.git/config", "[core]\n", 0o644),
		regular("skills/review/SKILL.md", "---\nname: \"review-code\"\ndescription: 'Reviews code.'\n---\n", 0o644),
		regular("skills/notes.md", "Not a skill.\n", 0o644),
	)
	var gotPath, gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotAuth = r.URL.Path, r.Header.Get("Authorization")
		w.Write(body)
	}))
	defer server.Close()

	s := New(t.TempDir(), t.TempDir(), func(name string) string {
		if name == "GH_TOKEN" {
			return "secret"
		}
		return ""
	})
	s.API = server.URL
	found, err := s.Fetch(context.Background(), skills.Source{Kind: skills.GitHub, Repo: "acme/tools", Ref: "v2"})
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/repos/acme/tools/tarball/v2" || gotAuth != "Bearer secret" {
		t.Errorf("asked for %s with %q, want /repos/acme/tools/tarball/v2 with the token", gotPath, gotAuth)
	}
	if len(found) != 2 {
		t.Fatalf("found %d skills, want 2: %+v", len(found), found)
	}
	review, docs := found[0], found[1]
	if review.Name != "review-code" || review.Description != "Reviews code." || review.Path != "skills/review/SKILL.md" {
		t.Errorf("review = %+v", review)
	}
	if docs.Name != "write-docs" || docs.Description != "Writes the docs, carefully." {
		t.Errorf("write-docs = %+v", docs)
	}
	var paths []string
	for _, f := range docs.Files {
		paths = append(paths, f.Path)
		if f.Path == "scripts/check.sh" && !f.Executable {
			t.Errorf("check.sh is not executable")
		}
	}
	if got := strings.Join(paths, " "); got != "SKILL.md scripts/check.sh" {
		t.Errorf("files = %s, want SKILL.md scripts/check.sh", got)
	}
}

// What GitHub refuses is explained.
func TestGitHubErrors(t *testing.T) {
	for status, want := range map[int]string{
		http.StatusNotFound:   "acme/tools was not found on GitHub",
		http.StatusForbidden:  "GITHUB_TOKEN",
		http.StatusBadGateway: "502",
	} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
		}))
		s := New(t.TempDir(), t.TempDir(), func(string) string { return "" })
		s.API = server.URL
		_, err := s.Fetch(context.Background(), skills.Source{Kind: skills.GitHub, Repo: "acme/tools"})
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("status %d: %v, want an error mentioning %q", status, err, want)
		}
		server.Close()
	}
}

// A directory is read from the working directory, or the home directory, and
// a source with no skills/ directory offers its own directories that hold a
// SKILL.md.  A skill with no name in its SKILL.md is named after its
// directory.
func TestFetchFromADirectory(t *testing.T) {
	work, home := t.TempDir(), t.TempDir()
	write(t, filepath.Join(work, "mine", "tidy", "SKILL.md"), "# Tidy\n")
	write(t, filepath.Join(work, "mine", "tidy", "references", "rules.md"), "Rules\n")
	write(t, filepath.Join(work, "mine", "docs", "README.md"), "Not a skill.\n")
	write(t, filepath.Join(home, "repo", "skills", "plan", "SKILL.md"), "---\nname: plan\n---\n")
	write(t, filepath.Join(home, "repo", "other", "SKILL.md"), "---\nname: other\n---\n")
	s := New(work, home, func(string) string { return "" })

	found, err := s.Fetch(context.Background(), skills.Source{Kind: skills.Local, Dir: "./mine"})
	if err != nil {
		t.Fatal(err)
	}
	if len(found) != 1 || found[0].Name != "tidy" || len(found[0].Files) != 2 || found[0].Path != "tidy/SKILL.md" {
		t.Errorf("./mine offers %+v, want tidy with two files", found)
	}

	found, err = s.Fetch(context.Background(), skills.Source{Kind: skills.Local, Dir: "~/repo"})
	if err != nil {
		t.Fatal(err)
	}
	if len(found) != 1 || found[0].Name != "plan" {
		t.Errorf("~/repo offers %+v, want only plan, under skills/", found)
	}

	if _, err := s.Fetch(context.Background(), skills.Source{Kind: skills.Local, Dir: "./missing"}); err == nil {
		t.Errorf("a missing directory was read")
	}
}

// A skill must be named safely, and once.
func TestSkillsMustBeNamedSafelyAndOnce(t *testing.T) {
	for _, test := range []struct {
		files map[string]string
		want  string
	}{
		{map[string]string{"skills/a/SKILL.md": "---\nname: ../evil\n---\n"}, "cannot be a directory's name"},
		{map[string]string{"skills/a/SKILL.md": "---\nname: same\n---\n", "skills/b/SKILL.md": "---\nname: same\n---\n"}, "both name"},
	} {
		files := map[string]file{}
		for name, content := range test.files {
			files[name] = file{content: []byte(content)}
		}
		if _, err := discover(files); err == nil || !strings.Contains(err.Error(), test.want) {
			t.Errorf("discover(%v) = %v, want an error mentioning %q", test.files, err, test.want)
		}
	}
}

func write(t *testing.T, name, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(name, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
