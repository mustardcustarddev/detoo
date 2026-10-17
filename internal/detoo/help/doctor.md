Checks the project's adr, rfc, sty and pse logs together, and reports any log
the tools would not find, or whose file names are out of date:

* a log kept under `doc/` instead of `docs/`;
* a log kept in a directory its tool no longer looks in, such as user stories
  in `docs/sty` instead of `docs/userstory`;
* documents whose file names do not begin with their log's prefix in upper
  case: `ADR-`, `RFC-`, `STORY-` or `PIOSEE-`, including a slip of the prefix,
  such as `STOY-0003-title.md` among the user stories;
* documents that share a number, such as `STORY-0003-first.md` and
  `STORY-0003-second.md`.

When moving or renaming would clash, such as when both `docs/sty` and
`docs/userstory` exist, or `.sty-dir` names the old directory, the problem is
reported as needing fixing by hand.

In a git repository, doctor works out how to renumber documents that share a
number from when git first committed each.  The first keeps the number.  Each
other goes after the last document committed before it, and before the one
after that, and the documents after it move up as far as they must to make
room.  A document not yet committed goes last.  Outside a git repository,
which came first is not known, so doctor suggests new numbers, and leaves the
renumbering to be done by hand.

The project is the nearest of the current directory and its parents that
holds a `docs` or `doc` directory, or a tool's marker file such as `.adr-dir`.

Without `--fix`, nothing is changed, and doctor says how many of the problems
`--fix` would fix.  With it, doctor moves each log to where its tool looks,
renames its documents to begin with the prefix, and changes the links to them
in the project's Markdown documents (those in the project's directory, and
under `docs`, `doc` and each log's directory), so that none break.  A
directory that moving a log out of leaves empty, such as `doc`, is removed.
Clashes are left as they are.  doctor then checks the logs again.

`--fix` renumbers documents that share a number only when asked to.  It
renames each document it renumbers, changes the number in its title, and
changes the links to it, and their text when it begins with the number, such
as "3. Use PostgreSQL".  Anything else that refers to a renumbered document
could be affected: links from outside the project's Markdown, other branches,
and mentions of its number in text, tickets and commit messages.  So at a
terminal, doctor warns of this, and asks first, renumbering only if the answer
is yes.  `--renumber` renumbers without asking, and `--renumber=false` does
not renumber; when doctor is not run at a terminal, it does not ask, and does
not renumber unless given `--renumber`.

In a git repository, when there is something to move or rename, `--fix` also
asks whether to move and rename with `git mv`, so that git records each as a
rename rather than a deletion and a new file.  Files git does not track are moved as
they would be without git.  `--git` moves with git without asking, and
`--git=false` without git; when doctor is not run at a terminal, it does not
ask, and moves without git unless given `--git`.

Before asking anything, `--fix` lists what it will do whatever the answers,
such as the documents it renames to begin with their prefix, so that a
question is not taken for the whole of the fix.

When nothing is wrong, or nothing is left wrong after a fix, doctor says so
and exits 0; otherwise it exits 1.
