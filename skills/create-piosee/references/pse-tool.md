# The pse tool

`pse` keeps a log of numbered, linked Markdown PIOSEE records in the
repository. This is what an agent needs to drive it without a terminal of its
own.

## Running it non-interactively

Two things can block an agent:

* `pse new` opens the new record in `$VISUAL` or `$EDITOR`, and waits for the
  editor to close. With neither set, it skips the editor and only prints the
  path. **Always clear both:** `VISUAL= EDITOR= pse new "Title"`.
* On a terminal, `pse help` shows its text in a pager (`$PSE_PAGER`, then
  `$PAGER`, then `more`). **Always set one that exits:**
  `PSE_PAGER=cat pse help new`.

Options come before operands. An operand that begins with a hyphen must follow
`--`. A title is a single argument, so quote it.

## Commands

* `pse config`: prints `pse_dir="…"` and `pse_template="…"`: the log
  directory, relative to the working directory, and the project's template in
  the log directory, or empty if there is none. Does not report
  `$PSE_TEMPLATE`. Works before the log exists.
* `pse list`: prints the path of every record, one a line, in order. Prints
  nothing if there are none. Fails with
  `pse: the docs/piosee directory does not exist` if there is no log.
* `pse init [DIRECTORY]`: starts a log in `DIRECTORY`, default `docs/piosee`,
  and writes record 1, "Make decisions with PIOSEE". Does nothing if the log
  already has records.
* `pse new [-r REF]... [-u REF]... TITLE`: creates the next numbered record
  from the template, in draft, and prints its path.
* `pse link (-r | -u) SOURCE TARGET`: records, in both records, that `SOURCE`
  reverses or updates `TARGET`. Exactly one of `-r` and `-u`.
* `pse list --toc [-i INTRO] [-o OUTRO] [-p PREFIX]`: writes a Markdown table
  of contents to stdout, headed `PIOSEE Decisions`.
* `pse graph [--mermaid | --graphviz] [-p PREFIX] [-e EXT]`: draws the
  records and the links between them.

A reference (`REF`, `SOURCE`, `TARGET`) is a record's number, its path, or
part of its file name. Digits alone are a number: `2` finds record 2 or
nothing, never record 12.

## What the tool writes

`pse new` copies the template and replaces the first `NUMBER`, `TITLE`,
`DATE` and `STATUS` on each line. The file name is `piosee-`, the four-digit
number and a slug of the title, such as
`docs/piosee/piosee-0004-roll-back-release-2-3.md`.

* `DATE` is today in ISO 8601 (`2026-10-16`), or `$PSE_DATE` exactly as given.
* `STATUS` is `Draft`.
* Everything else in the template, including its guidance text, is copied as
  it is. The guidance text must be replaced by the record's content.

Only files named `piosee-NNNN-….md` are records. Anything else in the log
directory, such as a `README.md`, is ignored.

## Links

Links are lines in the `## Status` section, written in both records:

```markdown
## Status

Draft

Reverses [4. Roll back release 2.3](piosee-0004-roll-back-release-2-3.md)
```

| Option            | Written in the new record   | Written in the old record      | Old record's `Draft` line |
| ----------------- | --------------------------- | ------------------------------ | ------------------------- |
| `-r`, `--reverse` | `Reverses [N. Title](file)` | `Reversed by [N. Title](file)` | Removed                   |
| `-u`, `--update`  | `Updates [N. Title](file)`  | `Updated by [N. Title](file)`  | Kept                      |

Only a line reading exactly `Draft` is removed. Any other status, such as
`Executed` or `Evaluated`, stays beside the link.

Use `-r` when a decision undoes or abandons an earlier one, typically because
its evaluation showed it had not worked. Use `-u` when a decision builds on an
earlier one without abandoning it, such as the next cycle of PIOSEE after an
evaluation found the problem had changed.

> **The link takes effect at once.** `pse new -r 4 "Title"` marks record 4
> "Reversed by" the new record, and removes record 4's `Draft` line, the
> moment the new record is created. `-u` adds "Updated by" to the old record
> at once, too.

Never write link lines by hand. Use `-r` and `-u` on `pse new`, or `pse link`
afterwards, so that both records agree.

## Templates

The template in use is, in order of preference:

1. the file named by `$PSE_TEMPLATE`, an absolute path or one relative to
   the working directory;
2. `templates/template.md` inside the log directory, which `pse config`
   reports as `pse_template`;
3. the built-in template, below.

The built-in template:

```markdown
# NUMBER. TITLE

Date: DATE

## Status

STATUS

## Problem

What is wrong, or what must be decided: what is happening, what it affects, how badly, and how long there is to act.

## Information

The facts that bear on the problem, and where each came from: what is known, what is not, and who was asked.

## Options

Each course of action considered, including doing nothing, with its risks and benefits.

## Select

The option chosen, who chose it, and why it was chosen over the others.

## Execute

How the option is carried out: the actions, who takes each, in what order, and what shows it is working.

## Evaluate

Whether the action solved the problem, what was learned, and whether the decision must be revisited. Until it is evaluated: how and when it will be.
```

A project's template may have other sections. Read the template in use for
every record, rather than assuming it is this one.

## Numbering

A new record's number is one more than the highest number in the log
directory. Numbers are never reused while the records are kept, but deleting
the highest-numbered record means the next `pse new` reuses its number.

## Environment

| Variable             | Effect                                     |
| -------------------- | ------------------------------------------ |
| `VISUAL`, `EDITOR`   | Editor opened on a new record. Clear both. |
| `PSE_TEMPLATE`       | Template for new records.                  |
| `PSE_DATE`           | Date written into a new record, as given.  |
| `PSE_PAGER`, `PAGER` | Pager for help. Set to `cat`.              |

The log directory is `docs/piosee`, unless a `.pse-dir` file in the project
root names another. The working directory and each of its ancestors are
searched in turn for a `.pse-dir` file, then for a `docs/piosee` directory,
so the tool works from anywhere in the project. Run `pse config` rather than
assuming.

## Errors

The tool exits `0` on success and `1` on any failure, with a message on stderr
that begins `pse:`. Report the message to the user rather than working around
it.
