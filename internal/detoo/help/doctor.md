Checks the project's adr, rfc, sty and pse logs together, and reports any log
the tools would not find, or whose file names are out of date:

* a log kept under `doc/` instead of `docs/`;
* a log kept in a directory its tool no longer looks in, such as user stories
  in `docs/sty` instead of `docs/userstory`;
* documents whose file names do not begin with their log's prefix in upper
  case: `ADR-`, `RFC-`, `STORY-` or `PIOSEE-`.

When moving or renaming would clash, such as when both `docs/sty` and
`docs/userstory` exist, or `.sty-dir` names the old directory, the problem is
reported as needing fixing by hand.

The project is the nearest of the current directory and its parents that
holds a `docs` or `doc` directory, or a tool's marker file such as `.adr-dir`.

Without `--fix`, nothing is changed, and doctor says how many of the problems
`--fix` would fix.  With it, doctor moves each log to where its tool looks,
renames its documents to begin with the prefix, and changes the links to them
in the project's Markdown documents (those in the project's directory, and
under `docs`, `doc` and each log's directory), so that none break.  A directory that moving a log out of leaves empty, such as `doc`, is
removed.  Clashes are left as they are.  doctor then checks the logs again.

In a git repository, when there is something to move or rename, `--fix` asks
whether to move and rename with `git mv`, so that git records each as a rename
rather than a deletion and a new file.  Files git does not track are moved as
they would be without git.  `--git` moves with git without asking, and
`--git=false` without git; when doctor is not run at a terminal, it does not
ask, and moves without git unless given `--git`.

When nothing is wrong, or nothing is left wrong after a fix, doctor says so
and exits 0; otherwise it exits 1.
