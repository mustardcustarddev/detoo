# 2. Use a custom test harness

Date: 2026-09-10

## Status

Accepted

## Context

The tests in `test/` run the `adr` command the way a user does. Each test is a
series of commands in a new temporary directory. The test records everything
the commands print, plus listings and contents of files on disk, in one
transcript, and compares it with a golden file in `test/testdata`. The
`-update` flag rewrites the golden files. The tests were ported from the shell
test scripts of adr-tools. They use a custom test harness: about 260 lines in
`test/harness_test.go` and `test/main_test.go`, built only on the standard
`testing` package.

Because the harness is custom, we considered replacing it with a framework. The
unit tests under `internal/` were out of scope; they are fine as they are.

We considered three alternatives:

* **Ginkgo, with Gomega and gexec.** gexec would build the binaries, and each
  step would assert its own exit code and output. Ginkgo has no support for
  transcripts or golden files, though. We would either write every expected
  record in Go strings, or write a golden matcher, which is the code we already
  have. It also adds a large set of dependencies to a module that depends only
  on Cobra, and gives the repository two test styles.

* **testscript, from `github.com/rogpeppe/go-internal`.** Each test would be a
  `.txtar` script. testscript has commands to run programs, check exit codes,
  compare files and update expected output, and it runs inside `go test`. It
  is the closest fit, but it means rewriting every test and golden file, and it
  has no `head` command without writing one.

* **Commander, from `github.com/commander-cli/commander`.** Commander suits
  smoke tests of a single command, and has no concept of a scenario.
  * It runs one shell command per test, sorted by title.
  * It has no temporary directory for each test.
  * It can check only stdout, stderr and the exit code.
  * It trims output before comparing it, so it cannot tell a command that
    prints nothing from one that prints a blank line.
  * It runs outside `go test`, and has not been released since March 2023.

## Decision

We will keep the custom test harness in `test/`, which compares golden
transcripts using only the standard `testing` package.

## Consequences

There are no new dependencies. `go test ./...` still runs every test, and the
golden files still read like a terminal session.

We keep maintaining the harness, and its known weaknesses remain until someone
fixes them:

* A transcript does not record exit codes, so a command that prints the right
  error but exits 0 still passes.
* A failure prints the whole transcript and the whole golden file, without a
  diff.
* The tests run one at a time.

Each of these is a small change to the harness. Fixing them does not need this
decision to be revisited.
