# 3. Run a detoo release without installing it

Date: 2026-10-20

## Status

Accepted

## Steps

1. Make sure you have `git`, `make`, the Go version that the release's
   `go.mod` names, and a clone of `mustardcustarddev/detoo`. This procedure
   builds a published release, such as the one a reporter has, from its tag,
   and runs it beside any detoo you have installed, without changing either.
2. Write down the release to run, `<version>`, such as `v0.2.0`, and the same
   version without its `v`, `<x.y.z>`. `brew list --versions detoo` prints
   the installed release as `detoo <x.y.z>`.
3. In the clone, fetch the tags with `git fetch --tags`, and check that
   `git tag --list <version>` prints `<version>`. A revoked release keeps its
   tag, as [SOP 2](SOP-0002-revoke-a-detoo-release.md) keeps it, so it can be
   run too.
   * If it prints nothing, stop: there is no release `<version>`.
4. Check the release out into a new worktree outside the clone, `<dir>`, such
   as `~/detoo-<x.y.z>`, with `git worktree add --detach <dir> <version>`. The
   clone's own files do not change.
5. Build it with `make --directory=<dir> build`, which vets, tests and
   compiles every command into `<dir>/build/`. If it fails, stop, and
   escalate.
6. Make a scratch project outside the clone and outside `<dir>`, with its own
   `docs` or `doc` directory and the files the commands need, such as a log
   like the reporter's, and go into it with `cd`. `detoo doctor` acts on the
   nearest of the current directory and its parents that holds a `docs` or
   `doc` directory (`detoo help doctor`), and both the clone and `<dir>` hold
   one.
7. Check the version with `<dir>/build/detoo doctor`, which changes nothing
   without `--fix`. Its header shows `<x.y.z>` beside the detoo badge. If it
   shows anything else, such as a commit hash, stop, and escalate.
8. Run each command, in the scratch project, by its path in `<dir>/build/`,
   such as `<dir>/build/detoo doctor --fix`. Do not install the release, with
   `make deploy` or Homebrew.
9. Afterwards, from the clone, remove the worktree with
   `git worktree remove <dir>`, and remove the scratch project.
