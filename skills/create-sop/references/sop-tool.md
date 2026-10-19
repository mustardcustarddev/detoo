# The sop tool

`sop` keeps a log of numbered, linked Markdown standard operating procedures
(SOPs) in the repository. This is what an agent needs to drive it without a
terminal of its own.

## Running it non-interactively

Two things can block an agent:

* `sop new` opens the new SOP in `$VISUAL` or `$EDITOR`, and waits for the
  editor to close. With neither set, it skips the editor and only prints the
  path. **Always clear both:** `VISUAL= EDITOR= sop new "Title"`.
* `sop help` prints to stdout, unless `$SOP_PAGER` names a pager, which it
  runs on a terminal. `$PAGER` is ignored. **If it is set, clear it:**
  `SOP_PAGER= sop help new`.

Options come before operands. An operand that begins with a hyphen must follow
`--`. A title is a single argument, so quote it.

## Commands

* `sop config`: prints `sop_dir="…"` and `sop_template="…"`: the log
  directory, relative to the working directory, and the project's template in
  the log directory, or empty if there is none. Does not report
  `$SOP_TEMPLATE`. Works before the log exists.
* `sop list`: prints the path of every SOP, one a line, in order. Prints
  nothing if there are none. Fails with
  `sop: the docs/sop directory does not exist` if there is no log.
* `sop init`: starts a log in `docs/sop`, and writes SOP 1, "Write a
  standard operating procedure". Does nothing if the log already has SOPs.
* `sop new [-s REF]... [-u REF]... TITLE`: creates the next numbered SOP from
  the template, in draft, and prints its path.
* `sop link (-s | -u) SOURCE TARGET`: records, in both SOPs, that `SOURCE`
  supersedes or updates `TARGET`. Exactly one of `-s` and `-u`.
* `sop list --toc [-i INTRO] [-o OUTRO] [-p PREFIX]`: writes a Markdown table
  of contents to stdout, headed `Standard Operating Procedures`.
* `sop graph [--mermaid | --graphviz] [-p PREFIX] [-e EXT]`: draws the SOPs
  and the links between them.

A reference (`REF`, `SOURCE`, `TARGET`) is an SOP's number, its path, or part
of its file name. Digits alone are a number: `2` finds SOP 2 or nothing, never
SOP 12.

## What the tool writes

`sop new` copies the template and replaces the first `NUMBER`, `TITLE`,
`DATE` and `STATUS` on each line. The file name is `SOP-`, the four-digit
number and a slug of the title, such as
`docs/sop/SOP-0004-roll-back-a-release.md`.

* `DATE` is today in ISO 8601 (`2026-10-16`), or `$SOP_DATE` exactly as given.
* `STATUS` is `Draft`.
* Everything else in the template, including its guidance text, is copied as
  it is. The guidance text must be replaced by the SOP's steps.

Only files named `SOP-NNNN-….md` are SOPs. Anything else in the log
directory, such as a `README.md`, is ignored.

## Links

Links are lines in the `## Status` section, written in both SOPs:

```markdown
## Status

Draft

Supersedes [4. Roll back a release](SOP-0004-roll-back-a-release.md)
```

| Option              | Written in the new SOP        | Written in the old SOP           | Old SOP's status line |
| ------------------- | ----------------------------- | -------------------------------- | --------------------- |
| `-s`, `--supersede` | `Supersedes [N. Title](file)` | `Superseded by [N. Title](file)` | Removed               |
| `-u`, `--update`    | `Updates [N. Title](file)`    | `Updated by [N. Title](file)`    | Kept                  |

Only a line reading exactly `Draft` or `Approved` is removed. Any other
status stays beside the link.

Use `-s` when a new SOP is a new procedure for the task an earlier one set
out, such as after the deployment changed. Use `-u` when a new SOP changes
part of an earlier one without replacing it, such as a step for a case the
earlier one left out.

> **The link takes effect at once.** `sop new -s 4 "Title"` marks SOP 4
> "Superseded by" the new SOP, and removes SOP 4's `Draft` or `Approved`
> line, the moment the new SOP is created. `-u` adds "Updated by" to the old
> SOP at once, too.

Never write link lines by hand. Use `-s` and `-u` on `sop new`, or `sop link`
afterwards, so that both SOPs agree.

## Templates

The template in use is, in order of preference:

1. the file named by `$SOP_TEMPLATE`, an absolute path or one relative to
   the working directory;
2. `templates/template.md` inside the log directory, which `sop config`
   reports as `sop_template`;
3. the built-in template, below.

The built-in template:

```markdown
# NUMBER. TITLE

Date: DATE

## Status

STATUS

## Steps

1. The first step, as an instruction: what to do, and where, naming the exact command, page or setting.
2. The next step, with how to check that it worked before going on.
```

The link lines go in the Status section, which ends at the next `##`
heading, so keep a heading, such as `## Steps`, after it. A project's
template may have other sections. Read the template in use for every SOP,
rather than assuming it is this one.

## Numbering

A new SOP's number is one more than the highest number in the log directory.
Numbers are never reused while the SOPs are kept, but deleting the
highest-numbered SOP means the next `sop new` reuses its number.

## Environment

| Variable           | Effect                                     |
| ------------------ | ------------------------------------------ |
| `VISUAL`, `EDITOR` | Editor opened on a new SOP. Clear both.    |
| `SOP_TEMPLATE`     | Template for new SOPs.                     |
| `SOP_DATE`         | Date written into a new SOP, as given.     |
| `SOP_PAGER`        | Pager for help, if set. Clear it.          |

The log directory is `docs/sop`. The working directory and each of its
ancestors are searched in turn for a `docs/sop` directory, so the tool works
from anywhere in the project. Run `sop config` rather than assuming.

## Errors

The tool exits `0` on success and `1` on any failure, with a message on stderr
that begins `sop:`. Report the message to the user rather than working around
it.
