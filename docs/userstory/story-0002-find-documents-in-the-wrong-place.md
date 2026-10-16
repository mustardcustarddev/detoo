# 2. Find documents in the wrong place

Date: 2026-10-16

## Status

Draft

## Story

As a member of the team, I want `detoo doctor` to tell me which of my logs are
in the wrong directory or have out-of-date file names, so that I can find logs
the tools have lost after a rename and keep my documentation up to date.

## Acceptance criteria

* `detoo doctor` checks the project's adr, rfc, sty and pse logs together.
* It reports a log kept under `doc/` instead of `docs/`.
* It reports a log kept in a former directory, such as stories in `docs/sty`
  instead of `docs/userstory`.
* It reports stories whose file names lack the `story-` prefix, such as
  `0001-describe-work-in-user-stories.md`.
* When a fix would clash, such as when both `docs/sty` and `docs/userstory`
  exist, or `.sty-dir` names the old directory, it reports that the problem
  needs fixing by hand.
* It changes nothing.
* When nothing is wrong, it says so and exits 0. When it finds problems, it
  exits 1.

## Notes

There is no `detoo` command today: the project has four tools, `adr`, `rfc`,
`sty` and `pse`.

A tool loses a renamed log because it looks only for its `.<tool>-dir` file
and then for `docs/<dir>`: stories left in `docs/sty`, for example, are no
longer found. Stories were kept in `docs/sty`, with no file-name prefix, until
the default changed to `docs/userstory` and the `story-` prefix.

Out of scope: what is inside a document, such as missing sections, or broken
links that a rename did not cause.

Story 3, "Fix documents in the wrong place", fixes the problems this story
finds.
