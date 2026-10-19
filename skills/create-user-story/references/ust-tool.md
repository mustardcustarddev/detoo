# The ust tool

`ust` keeps a log of numbered, linked Markdown user stories in the repository.
This is what an agent needs to drive it without a terminal of its own.

## Running it non-interactively

Two things can block an agent:

* `ust new` opens the new story in `$VISUAL` or `$EDITOR`, and waits for the
  editor to close. With neither set, it skips the editor and only prints the
  path. **Always clear both:** `VISUAL= EDITOR= ust new "Title"`.
* `ust help` prints to stdout, unless `$UST_PAGER` names a pager, which it
  runs on a terminal. `$PAGER` is ignored. **If it is set, clear it:**
  `UST_PAGER= ust help new`.

Options come before operands. An operand that begins with a hyphen must follow
`--`. A title is a single argument, so quote it.

## Commands

* `ust config`: prints `ust_dir="…"` and `ust_template="…"`: the log
  directory, relative to the working directory, and the project's template in
  the log directory, or empty if there is none. Does not report
  `$UST_TEMPLATE`. Works before the log exists.
* `ust list`: prints the path of every story, one a line, in order. Prints
  nothing if there are none. Fails with
  `ust: the docs/userstory directory does not exist` if there is no log.
* `ust init`: starts a log in `docs/userstory`, and writes story 1, "Describe
  work in user stories". Does nothing if the log already has stories.
* `ust new [-r REF]... [-u REF]... TITLE`: creates the next numbered story
  from the template, in draft, and prints its path.
* `ust link (-r | -u) SOURCE TARGET`: records, in both stories, that `SOURCE`
  replaces or updates `TARGET`. Exactly one of `-r` and `-u`.
* `ust list --toc [-i INTRO] [-o OUTRO] [-p PREFIX]`: writes a Markdown table
  of contents to stdout.
* `ust graph [--mermaid | --graphviz] [-p PREFIX] [-e EXT]`: draws the
  stories and the links between them.

A reference (`REF`, `SOURCE`, `TARGET`) is a story's number, its path, or
part of its file name. Digits alone are a number: `2` finds story 2 or nothing,
never story 12.

## What the tool writes

`ust new` copies the template and replaces the first `NUMBER`, `TITLE`, `DATE`
and `STATUS` on each line. The file name is `STORY-`, the four-digit number and
a slug of the title, such as
`docs/userstory/STORY-0004-reorder-a-previous-purchase.md`.

* `DATE` is today in ISO 8601 (`2026-10-16`), or `$UST_DATE` exactly as given.
* `STATUS` is `Draft`.
* Everything else in the template, including its guidance text, is copied as
  it is. The guidance text must be replaced by the story's content.

Only files named `STORY-NNNN-….md` are stories, along with stories named
before the prefix was upper case (`story-NNNN-….md`) or before there was one
(`NNNN-….md`). Anything else in the log directory, such as a `README.md`, is
ignored. `detoo doctor` reports stories with the older names, and stories left
in `docs/sty`, where ust no longer looks.

## Links

Links are lines in the `## Status` section, written in both stories:

```markdown
## Status

Draft

Replaces [4. Manage saved sessions](STORY-0004-manage-saved-sessions.md)
```

| Option            | Written in the new story    | Written in the old story       | Old story's `Draft` line |
| ----------------- | --------------------------- | ------------------------------ | ------------------------ |
| `-r`, `--replace` | `Replaces [N. Title](file)` | `Replaced by [N. Title](file)` | Removed                  |
| `-u`, `--update`  | `Updates [N. Title](file)`  | `Updated by [N. Title](file)`  | Kept                     |

Only a line reading exactly `Draft` is removed. Any other status, such as
`Ready` or `Done`, stays beside the link.

Use `-r` when a story is split: each of the smaller stories replaces the large
one, which is then marked "Replaced by" every one of them. Use `-u` when a
story changes what an earlier story asked for, such as a new rule for a
feature that is already built.

> **The link takes effect at once.** `ust new -r 4 "Title"` marks story 4
> "Replaced by" the new story, and removes story 4's `Draft` line, the moment
> the new story is created. `-u` adds "Updated by" to the old story at once,
> too.

Never write link lines by hand. Use `-r` and `-u` on `ust new`, or `ust link`
afterwards, so that both stories agree.

## Templates

The template in use is, in order of preference:

1. the file named by `$UST_TEMPLATE`, an absolute path or one relative to
   the working directory;
2. `templates/template.md` inside the log directory, which `ust config`
   reports as `ust_template`;
3. the built-in template, below.

The built-in template:

```markdown
# NUMBER. TITLE

Date: DATE

## Status

STATUS

## Story

As a [type of user], I [need/want/am required] to [do something], so that [reason or benefit].

## Acceptance criteria

What must be true for the story to be complete: the examples, rules and tests that will confirm it.

## Notes

What the conversation about the story has settled so far, and what it has still to settle.
```

A project's template may have other sections. Read the template in use for
every story, rather than assuming it is this one.

## Numbering

A new story's number is one more than the highest number in the log
directory. Numbers are never reused while the stories are kept, but deleting
the highest-numbered story means the next `ust new` reuses its number.

## Environment

| Variable             | Effect                                    |
| -------------------- | ----------------------------------------- |
| `VISUAL`, `EDITOR`   | Editor opened on a new story. Clear both. |
| `UST_TEMPLATE`       | Template for new stories.                 |
| `UST_DATE`           | Date written into a new story, as given.  |
| `UST_PAGER`          | Pager for help, if set. Clear it.         |

The log directory is `docs/userstory`. The working directory and each of its
ancestors are searched in turn for a `docs/userstory` directory, so the tool
works from anywhere in the project. Run `ust config` rather than assuming. If
the stories are in `docs/sty`, from when the tool was called `sty`, suggest
`detoo doctor`, which finds them and, with `--fix`, moves them to
`docs/userstory`.

## Errors

The tool exits `0` on success and `1` on any failure, with a message on stderr
that begins `ust:`. Report the message to the user rather than working around
it.
