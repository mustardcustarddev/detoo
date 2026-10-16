package clitest

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// workspace is a throwaway project directory and a transcript of what the
// commands run in it printed.  Its methods stand in for the shell commands the
// original test scripts used, and each one records what it did in the
// transcript, so that a golden file reads like a terminal session.
type workspace struct {
	t     *testing.T
	suite string            // the tool whose golden files hold this transcript
	root  string            // the project directory
	dir   string            // the current directory, relative to root
	env   map[string]string // environment variables set by the test
	log   bytes.Buffer      // the transcript
}

// newWorkspace returns an empty project, with the environment the tests share:
// a fixed date for each tool, so that records are reproducible, and no editor,
// so that creating a record does not try to open one.  The transcript is
// compared against a golden file among those of the suite's tool.
func newWorkspace(t *testing.T, suite string) *workspace {
	env := map[string]string{"VISUAL": "", "EDITOR": ""}
	for _, name := range names {
		env[variable(name, "DATE")] = "1992-01-12"
	}
	return &workspace{
		t:     t,
		suite: suite,
		root:  t.TempDir(),
		dir:   ".",
		env:   env,
	}
}

// variable returns the name of a tool's environment variable: ADR_DATE, for
// adr's DATE.
func variable(tool, setting string) string {
	return strings.ToUpper(tool) + "_" + setting
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

// adr runs the adr command.
func (w *workspace) adr(args ...string) {
	w.t.Helper()
	w.run("adr", args...)
}

// rfc runs the rfc command.
func (w *workspace) rfc(args ...string) {
	w.t.Helper()
	w.run("rfc", args...)
}

// sty runs the sty command.
func (w *workspace) sty(args ...string) {
	w.t.Helper()
	w.run("sty", args...)
}

// pse runs the pse command.
func (w *workspace) pse(args ...string) {
	w.t.Helper()
	w.run("pse", args...)
}

// detoo runs the detoo command, and records its exit status, as echo $? would
// show it.  The user's home directory is the project's, so that the project
// is shown as ~ whichever temporary directory it is in.
func (w *workspace) detoo(args ...string) {
	w.t.Helper()
	w.env["HOME"] = w.root
	code := w.run("detoo", args...)
	w.echo("echo $?")
	w.printf("%d\n", code)
}

// run runs a tool, recording its output — on both streams, in the order it was
// written — in the transcript, and returns its exit status.
func (w *workspace) run(tool string, args ...string) int {
	w.t.Helper()
	w.echo(tool, args...)

	cmd := exec.Command(tools[tool], args...)
	cmd.Dir = w.cwd()
	cmd.Env = w.environ()
	cmd.Stdout, cmd.Stderr = &w.log, &w.log
	if err := cmd.Run(); err != nil {
		var exit *exec.ExitError
		if !errors.As(err, &exit) {
			w.t.Fatalf("running %s %s: %v", tool, strings.Join(args, " "), err)
		}
		return exit.ExitCode()
	}
	return 0
}

// output runs a tool and returns everything it wrote, leaving the transcript
// alone.  It is for output that is too long, or too much cobra's own, to keep
// in a golden file.
func (w *workspace) output(tool string, args ...string) string {
	w.t.Helper()
	cmd := exec.Command(tools[tool], args...)
	cmd.Dir = w.cwd()
	cmd.Env = w.environ()
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out
	if err := cmd.Run(); err != nil {
		w.t.Fatalf("running %s %s: %v\n%s", tool, strings.Join(args, " "), err, &out)
	}
	return out.String()
}

// environ is the environment commands run in: the test's own, without any
// template a tool might be pointed at, with help that is not paged, with
// reports drawn 80 columns wide and without colour, and with the variables
// the test has set.
func (w *workspace) environ() []string {
	env := os.Environ()
	for _, tool := range names {
		env = append(env, variable(tool, "TEMPLATE")+"=", variable(tool, "PAGER")+"=cat")
	}
	env = append(env, variable("detoo", "PAGER")+"=cat", "COLUMNS=", "CLICOLOR_FORCE=")
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

// touch creates empty files, and any directory they need.
func (w *workspace) touch(names ...string) {
	w.t.Helper()
	w.echo("touch", names...)
	for _, name := range names {
		if err := os.MkdirAll(filepath.Dir(w.path(name)), 0o755); err != nil {
			w.t.Fatal(err)
		}
		if err := os.WriteFile(w.path(name), nil, 0o644); err != nil {
			w.t.Fatal(err)
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

func (w *workspace) read(name string) string {
	w.t.Helper()
	content, err := os.ReadFile(w.path(name))
	if err != nil {
		w.t.Fatal(err)
	}
	return string(content)
}

// check compares the transcript against the test's golden file, in the
// directory of the suite's tool, which -update rewrites.
func (w *workspace) check() {
	w.t.Helper()
	golden := filepath.Join("testdata", w.suite, w.t.Name()+".golden")
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
