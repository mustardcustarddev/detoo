# The rfc tool

`rfc` keeps a log of numbered, linked Markdown requests for comments in the
repository. This is what an agent needs to drive it without a terminal of its
own.

## Running it non-interactively

Two things can block an agent:

* `rfc new` opens the new RFC in `$VISUAL` or `$EDITOR`, and waits for the
  editor to close. With neither set, it skips the editor and only prints the
  path. **Always clear both:** `VISUAL= EDITOR= rfc new "Title"`.
* `rfc help` prints to stdout, unless `$RFC_PAGER` names a pager, which it
  runs on a terminal. `$PAGER` is ignored. **If it is set, clear it:**
  `RFC_PAGER= rfc help new`.

Options come before operands. An operand that begins with a hyphen must follow
`--`. A title is a single argument, so quote it.

## Commands

* `rfc config`: prints `rfc_dir="…"` and `rfc_template="…"`: the log
  directory, relative to the working directory, and the project's template in
  the log directory, or empty if there is none. Does not report
  `$RFC_TEMPLATE`. Works before the log exists.
* `rfc list`: prints the path of every RFC, one a line, in order. Prints
  nothing if there are none. Fails with
  `rfc: the docs/rfc directory does not exist` if there is no log.
* `rfc init`: starts a log in `docs/rfc`, and writes RFC 1, "Use requests
  for comments". Does nothing if the log already has RFCs.
* `rfc new [-o REF]... [-u REF]... TITLE`: creates the next numbered RFC from
  the template, in draft, and prints its path.
* `rfc link (-o | -u) SOURCE TARGET`: records, in both RFCs, that `SOURCE`
  obsoletes or updates `TARGET`. Exactly one of `-o` and `-u`.
* `rfc list --toc [-i INTRO] [-o OUTRO] [-p PREFIX]`: writes a Markdown table
  of contents to stdout. Here `-o` is the outro, not obsolete.
* `rfc graph [--mermaid | --graphviz] [-p PREFIX] [-e EXT]`: draws the RFCs
  and the links between them.

A reference (`REF`, `SOURCE`, `TARGET`) is an RFC's number, its path, or
part of its file name. Digits alone are a number: `2` finds RFC 2 or nothing,
never RFC 12.

## What the tool writes

`rfc new` copies the template and replaces the first `NUMBER`, `TITLE`, `DATE`
and `STATUS` on each line. The file name is `RFC-`, the four-digit number and a
slug of the title, such as
`docs/rfc/RFC-0002-migrate-to-hexagonal-architecture.md`. RFCs named before
the prefix (`NNNN-….md`) are still read; `detoo doctor` reports them.

* `DATE` is today in ISO 8601 (`2026-10-16`), or `$RFC_DATE` exactly as given.
* `STATUS` is `Draft`.
* Everything else in the template, including its guidance text, is copied as
  it is. The guidance text must be replaced by the RFC's content.

## Links

Links are lines in the `## Status` section, written in both RFCs:

```markdown
## Status

Draft

Obsoletes [2. Migrate to hexagonal architecture](RFC-0002-migrate-to-hexagonal-architecture.md)
```

| Option             | Written in the new RFC       | Written in the old RFC          | Old RFC's `Draft` line |
| ------------------ | ---------------------------- | ------------------------------- | ---------------------- |
| `-o`, `--obsolete` | `Obsoletes [N. Title](file)` | `Obsoleted by [N. Title](file)` | Removed                |
| `-u`, `--update`   | `Updates [N. Title](file)`   | `Updated by [N. Title](file)`   | Kept                   |

Only a line reading exactly `Draft` is removed. Any other status, such as
`Accepted`, stays beside the link.

> **The link takes effect at once.** `rfc new -o 4 "Title"` marks RFC 4
> "Obsoleted by" the new RFC, and removes RFC 4's `Draft` line, the moment the
> new RFC is created, although the new RFC is itself only a draft that may
> never be accepted. `-u` adds "Updated by" to the old RFC at once, too.

Never write link lines by hand. Use `-o` and `-u` on `rfc new`, or `rfc link`
afterwards, so that both RFCs agree.

## Templates

The template in use is, in order of preference:

1. the file named by `$RFC_TEMPLATE`, an absolute path or one relative to
   the working directory;
2. `templates/template.md` inside the log directory, which `rfc config`
   reports as `rfc_template`;
3. the built-in template, below.

The built-in template:

```markdown
# NUMBER. TITLE

Date: DATE

## Status

STATUS

## Summary

The proposal, in a paragraph.

## Motivation

Why the change should be made: the problem it solves, and who it matters to.

## Design

The change, in enough detail to comment on and to build.

## Drawbacks

Why the change might not be made.

## Alternatives

What else was considered, and why this design was chosen over it.

## Unresolved questions

What should be settled before the RFC is accepted, and what can wait.
```

A project's template may have other sections. Read the template in use for
every RFC, rather than assuming it is this one.

## Numbering

A new RFC's number is one more than the highest number in the log directory.
Numbers are never reused while the RFCs are kept, but deleting the
highest-numbered RFC means the next `rfc new` reuses its number.

## Environment

| Variable             | Effect                                  |
| -------------------- | --------------------------------------- |
| `VISUAL`, `EDITOR`   | Editor opened on a new RFC. Clear both. |
| `RFC_TEMPLATE`       | Template for new RFCs.                  |
| `RFC_DATE`           | Date written into a new RFC, as given.  |
| `RFC_PAGER`          | Pager for help, if set. Clear it.       |

The log directory is `docs/rfc`. The working directory and each of its
ancestors are searched in turn for a `docs/rfc` directory, so the tool works
from anywhere in the project. Run `rfc config` rather than assuming.

## Errors

The tool exits `0` on success and `1` on any failure, with a message on stderr
that begins `rfc:`. Report the message to the user rather than working around
it.
