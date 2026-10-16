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

Nothing is changed.  When nothing is wrong, doctor says so and exits 0; when
it finds problems, it exits 1.
