# 3. Fix documents in the wrong place

Date: 2026-10-16

## Status

Draft

## Story

As a member of the team, I want `detoo doctor --fix` to move and rename my logs
to where the tools expect them, so that the tools find them again and I can
keep my documentation up to date.

## Acceptance criteria

* `detoo doctor --fix` moves a log from `doc/` to `docs/`.
* It moves a log from a former directory to the current one, such as
  `docs/sty` to `docs/userstory`.
* It renames stories to the current file names, such as
  `0001-describe-work-in-user-stories.md` to
  `story-0001-describe-work-in-user-stories.md`.
* When it renames a file, it updates the links to it in other documents, so
  that none break.
* It leaves clashes unchanged, and reports that they need fixing by hand.
* After a fix, the tools find the log again: `sty list`, for example, lists
  the stories that were in `docs/sty`.
* If problems remain afterwards, it exits 1. Otherwise it says all is well.

## Notes

Depends on story 2, "Find documents in the wrong place", which finds the
problems this story fixes.

Out of scope: what is inside a document.
