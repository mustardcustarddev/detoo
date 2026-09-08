package adrtest

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// workspace is a throwaway project directory and a transcript of what the
// commands run in it printed.  Its methods stand in for the shell commands the
// original test scripts used, and each one records what it did in the
// transcript, so that a golden file reads like a terminal session.
type workspace struct {
	t    *testing.T
	root string            // the project directory
	dir  string            // the current directory, relative to root
	env  map[string]string // environment variables set by the test
	log  bytes.Buffer      // the transcript
}

// newWorkspace returns an empty project, with the environment the tests share:
// a fixed date, so that records are reproducible, and no editor, so that
// creating a record does not try to open one.
func newWorkspace(t *testing.T) *workspace {
	return &workspace{
		t:    t,
		root: t.TempDir(),
		dir:  ".",
		env:  map[string]string{"ADR_DATE": "1992-01-12", "VISUAL": "", "EDITOR": ""},
	}
}

// cwd is the directory commands run in.
func (w *workspace) cwd() string {
	return filepath.Join(w.root, filepath.FromSlash(w.dir))
}

// path resolves a path given by a test against the current directory.
func (w *workspace) path(name string) string {
	return filepath.Join(w.cwd(), filepath.FromSlash(name))
}

// echo records a command in the transcript, as it would have been typed.
func (w *workspace) echo(command string, args ...string) {
	quoted := make([]string, 0, len(args))
	for _, arg := range args {
		if strings.ContainsAny(arg, " \t") {
			arg = fmt.Sprintf("%q", arg)
		}
		quoted = append(quoted, arg)
	}
	fmt.Fprintln(&w.log, strings.TrimSpace(command+" "+strings.Join(quoted, " ")))
}

// printf records output in the transcript.
func (w *workspace) printf(format string, a ...any) {
	fmt.Fprintf(&w.log, format, a...)
}

// adr runs the adr command, recording its output — on both streams, in the
// order it was written — in the transcript.
func (w *workspace) adr(args ...string) {
	w.t.Helper()
	w.echo("adr", args...)

	cmd := exec.Command(adrBin, args...)
	cmd.Dir = w.cwd()
	cmd.Env = w.environ()
	cmd.Stdout, cmd.Stderr = &w.log, &w.log
	if err := cmd.Run(); err != nil {
		var exit *exec.ExitError
		if !errors.As(err, &exit) {
			w.t.Fatalf("running adr %s: %v", strings.Join(args, " "), err)
		}
	}
}

func (w *workspace) environ() []string {
	env := append(os.Environ(), "ADR_TEMPLATE=", "ADR_PAGER=cat")
	for name, value := range w.env {
		env = append(env, name+"="+value)
	}
	return env
}

// setenv sets an environment variable for the commands that follow.
func (w *workspace) setenv(name, value string) {
	w.env[name] = value
	w.echo(name + "=" + value)
}

// useEditor points VISUAL or EDITOR at the fake editor, which reports how it was
// called.  The label distinguishes the two in the transcript.
func (w *workspace) useEditor(name, label string) {
	w.env[name] = fakeEditorBin + " " + label
	w.echo(name + "=fake-" + strings.ToLower(label))
}

// cd changes the current directory.
func (w *workspace) cd(dir string) {
	w.echo("cd", dir)
	w.dir = filepath.ToSlash(filepath.Join(w.dir, dir))
}

// mkdir creates a directory, and any parent it needs.
func (w *workspace) mkdir(dir string) {
	w.t.Helper()
	w.echo("mkdir -p", dir)
	if err := os.MkdirAll(w.path(dir), 0o755); err != nil {
		w.t.Fatal(err)
	}
}

// write creates a file with the given content.
func (w *workspace) write(name, content string) {
	w.t.Helper()
	w.echo("cat >", name, "<<EOF")
	w.printf("%sEOF\n", content)
	if err := os.WriteFile(w.path(name), []byte(content), 0o644); err != nil {
		w.t.Fatal(err)
	}
}

// ls lists the visible names in a directory.
func (w *workspace) ls(dir string) {
	w.t.Helper()
	w.echo("ls", dir)
	entries, err := os.ReadDir(w.path(dir))
	if err != nil {
		w.t.Fatal(err)
	}
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), ".") {
			w.printf("%s\n", entry.Name())
		}
	}
}

// cat shows a file.
func (w *workspace) cat(name string) {
	w.t.Helper()
	w.echo("cat", name)
	w.printf("%s", w.read(name))
}

// head shows the first n lines of a file.
func (w *workspace) head(n int, name string) {
	w.t.Helper()
	w.echo(fmt.Sprintf("head -%d", n), name)
	lines := strings.SplitAfter(w.read(name), "\n")
	for i, line := range lines {
		if i == n || line == "" {
			break
		}
		w.printf("%s", line)
	}
}

// grep shows the lines of the files matching a glob that contain a substring,
// each prefixed by the file it came from.
func (w *workspace) grep(substring, glob string) {
	w.t.Helper()
	w.echo("grep", substring, glob)
	matches, err := filepath.Glob(w.path(glob))
	if err != nil {
		w.t.Fatal(err)
	}
	sort.Strings(matches)
	for _, match := range matches {
		name, err := filepath.Rel(w.cwd(), match)
		if err != nil {
			w.t.Fatal(err)
		}
		for _, line := range strings.Split(strings.TrimSuffix(w.read(name), "\n"), "\n") {
			if strings.Contains(line, substring) {
				w.printf("%s:%s\n", filepath.ToSlash(name), line)
			}
		}
	}
}

func (w *workspace) read(name string) string {
	w.t.Helper()
	content, err := os.ReadFile(w.path(name))
	if err != nil {
		w.t.Fatal(err)
	}
	return string(content)
}

// check compares the transcript against the test's golden file, which -update
// rewrites.
func (w *workspace) check() {
	w.t.Helper()
	golden := filepath.Join("testdata", w.t.Name()+".golden")
	if *update {
		if err := os.WriteFile(golden, w.log.Bytes(), 0o644); err != nil {
			w.t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		w.t.Fatalf("%v (run: go test ./test -update)", err)
	}
	if got := w.log.String(); got != string(want) {
		w.t.Errorf("transcript does not match %s\n\n--- got ---\n%s\n--- want ---\n%s", golden, got, want)
	}
}
