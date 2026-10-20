# The fyr tool

`fyr` keeps a log of numbered, linked Markdown fire drills in the repository.
This is what an agent needs to drive it without a terminal of its own.

## Running it non-interactively

Two things can block an agent:

* `fyr new` opens the new drill in `$VISUAL` or `$EDITOR`, and waits for the
  editor to close. With neither set, it skips the editor and only prints the
  path. **Always clear both:** `VISUAL= EDITOR= fyr new "Title"`.
* `fyr help` prints to stdout, unless `$FYR_PAGER` names a pager, which it
  runs on a terminal. `$PAGER` is ignored. **If it is set, clear it:**
  `FYR_PAGER= fyr help new`.

Options come before operands. An operand that begins with a hyphen must follow
`--`. A title is a single argument, so quote it.

## Commands

* `fyr config`: prints `fyr_dir="…"` and `fyr_template="…"`: the log
  directory, relative to the working directory, and the project's template in
  the log directory, or empty if there is none. Does not report
  `$FYR_TEMPLATE`. Works before the log exists.
* `fyr list`: prints the path of every drill, one a line, in order. Prints
  nothing if there are none. Fails with
  `fyr: the docs/firedrill directory does not exist` if there is no log.
* `fyr init`: starts a log in `docs/firedrill`, and writes its
  `README.md`, which sets out how to practise a playbook in a fire drill. It
  writes no drill, so the first drill `fyr new` creates is number 1. Does
  nothing if the log already has drills, or a `README.md`.
* `fyr new [-r REF]... [-u REF]... TITLE`: creates the next numbered drill
  from the template, in draft, with its directory beside it, and prints the
  drill's path.
* `fyr link (-r | -u) SOURCE TARGET`: records, in both drills, that `SOURCE`
  replaces or updates `TARGET`. Exactly one of `-r` and `-u`.
* `fyr list --toc [-i INTRO] [-o OUTRO] [-p PREFIX]`: writes a Markdown table
  of contents to stdout, headed `Fire Drills`.
* `fyr graph [--mermaid | --graphviz] [-p PREFIX] [-e EXT]`: draws the drills
  and the links between them.

A reference (`REF`, `SOURCE`, `TARGET`) is a drill's number, its path, or part
of its file name. Digits alone are a number: `2` finds drill 2 or nothing,
never drill 12.

## What the tool writes

`fyr new` copies the template and replaces the first `NUMBER`, `TITLE`,
`DATE` and `STATUS` on each line. The file name is `FIREDRILL-`, the
four-digit number and a slug of the title, such as
`docs/firedrill/FIREDRILL-0004-customer-charged-twice.md`.

* `DATE` is today in ISO 8601 (`2026-10-19`), or `$FYR_DATE` exactly as
  given.
* `STATUS` is `Draft`.
* Everything else in the template, including its guidance text, is copied as
  it is. The guidance text must be replaced by the drill.

Beside the drill, it creates a directory named as the drill is without
`.md`, holding two directories, each with a `README.md` that says what goes
in it:

```text
docs/firedrill/
  FIREDRILL-0004-customer-charged-twice.md
  FIREDRILL-0004-customer-charged-twice/
    problem/README.md     the pages the drill sends its reader to
    solution/README.md    the expected response and the debrief
```

The READMEs are the tool's own, whatever the template: `$FYR_TEMPLATE` and
`fyr_template` change only the drill. The tool writes nothing else in the
directory; the pages, `expected-response.md` and `debrief.md` are yours to
write.

Only files named `FIREDRILL-NNNN-….md` are drills. Anything else in the log
directory, such as a `README.md`, or a drill's directory, is ignored.
`detoo doctor --fix`, when it renames or renumbers a drill, renames its
directory with it, and changes the links into it.

## Links

Links are lines in the `## Status` section, written in both drills:

```markdown
## Status

Draft

Replaces [4. Customer charged twice](FIREDRILL-0004-customer-charged-twice.md)
```

| Option            | Written in the new drill    | Written in the old drill       | Old drill's status line |
| ----------------- | --------------------------- | ------------------------------ | ----------------------- |
| `-r`, `--replace` | `Replaces [N. Title](file)` | `Replaced by [N. Title](file)` | Removed                 |
| `-u`, `--update`  | `Updates [N. Title](file)`  | `Updated by [N. Title](file)`  | Kept                    |

Only a line reading exactly `Draft` is removed. Any other status stays beside
the link.

Use `-r` when a new drill is a new drill for the case an earlier one
practised, such as after the playbook changed. Use `-u` when a new drill
changes part of an earlier one without replacing it, such as a choice it
left out.

> **The link takes effect at once.** `fyr new -r 4 "Title"` marks drill 4
> "Replaced by" the new drill, and removes drill 4's `Draft` line, the moment
> the new drill is created. `-u` adds "Updated by" to the old drill at once,
> too.

Never write link lines by hand. Use `-r` and `-u` on `fyr new`, or `fyr link`
afterwards, so that both drills agree.

## Templates

The template in use is, in order of preference:

1. the file named by `$FYR_TEMPLATE`, an absolute path or one relative to
   the working directory;
2. `templates/template.md` inside the log directory, which `fyr config`
   reports as `fyr_template`;
3. the built-in template, below.

The built-in template:

```markdown
# NUMBER. TITLE

Date: DATE

## Status

STATUS

## Playbook

The playbook this drill practises, linked.

## How to play

This drill is for one person. Set it up, then read the drill a step at a time, and do what each step says. Where a step asks you to choose, write down your choice and why, then read the page it names, in the `problem` directory beside this document. When a page says the drill is over, read the `solution` directory.

## Ground rules

* Act only on the local system. A step that would act on anything else is written down, not taken.
* Send nothing. Write any message the drill has you send, and keep it for the debrief.
* Stop at any time, for a real incident or any other reason.

## Setup

What to have to hand, and how to run the system locally, before the drill starts.

## The drill

1. What happens first, and what to read.
2. What happens next, and the choice it asks for: each option, and the page in `problem` it leads to.
```

The link lines go in the Status section, which ends at the next `##`
heading, so keep a heading, such as `## Playbook`, after it. A project's
template may have other sections. Read the template in use for every drill,
rather than assuming it is this one.

## Numbering

A new drill's number is one more than the highest number in the log
directory. Numbers are never reused while the drills are kept, but deleting
the highest-numbered drill means the next `fyr new` reuses its number.

## Environment

| Variable           | Effect                                     |
| ------------------ | ------------------------------------------ |
| `VISUAL`, `EDITOR` | Editor opened on a new drill. Clear both.  |
| `FYR_TEMPLATE`     | Template for new drills.                   |
| `FYR_DATE`         | Date written into a new drill, as given.   |
| `FYR_PAGER`        | Pager for help, if set. Clear it.          |

The log directory is `docs/firedrill`. The working directory and each of its
ancestors are searched in turn for a `docs/firedrill` directory, so the tool
works from anywhere in the project. Run `fyr config` rather than assuming.

## Errors

The tool exits `0` on success and `1` on any failure, with a message on stderr
that begins `fyr:`. Report the message to the user rather than working around
it.
