# The adr tool

`adr` keeps a log of numbered, linked Markdown records in the repository. This
is what an agent needs to drive it without a terminal of its own.

## Running it non-interactively

Two things can block an agent:

* `adr new` opens the new record in `$VISUAL` or `$EDITOR`, and waits for the
  editor to close. With neither set, it skips the editor and only prints the
  path. **Always clear both:** `VISUAL= EDITOR= adr new "Title"`.
* On a terminal, `adr help` shows its text in a pager (`$ADR_PAGER`, then
  `$PAGER`, then `more`). **Always set one that exits:**
  `ADR_PAGER=cat adr help new`.

Options come before operands. An operand that begins with a hyphen must follow
`--`. A title is a single argument, so quote it.

## Commands

| Command | What it does |
| --- | --- |
| `adr config` | Prints `adr_dir="…"` and `adr_template="…"`: the log directory, relative to the working directory, and the project's template, or empty if the built-in template is used. |
| `adr list` | Prints the path of every record, one a line, in order. Prints nothing if there are none. |
| `adr init [DIRECTORY]` | Starts a log in `DIRECTORY`, default `docs/adr`, and writes record 1, "Record architecture decisions". Does nothing if the log already has records. |
| `adr new [-s REF]... [-u REF]... TITLE` | Creates the next numbered record from the template and prints its path. |
| `adr link (-s \| -u) SOURCE TARGET` | Records, in both records, that `SOURCE` supersedes or updates `TARGET`. Exactly one of `-s` and `-u`. |
| `adr list --toc [-i INTRO] [-o OUTRO] [-p PREFIX]` | Writes a Markdown table of contents to stdout. |
| `adr graph [--mermaid \| --graphviz] [-p PREFIX] [-e EXT]` | Draws the records and the links between them. |

A reference (`REF`, `SOURCE`, `TARGET`) is a record's number, or part of its
file name.

## What the tool writes

`adr new` copies the template and replaces the first `NUMBER`, `TITLE`, `DATE`
and `STATUS` on each line. The file name is the four-digit number and a slug of
the title, such as `docs/adr/0004-use-postgresql.md`.

* `DATE` is today in ISO 8601 (`2026-10-16`), or `$ADR_DATE` exactly as given.
* `STATUS` is `Accepted`.
* Everything else in the template, including its guidance text, is copied as
  it is. The guidance text must be replaced by the record's content.

Links are lines in the `## Status` section, written in both records:

```markdown
## Status

Accepted

Supercedes [2. Use MySQL](0002-use-mysql.md)
```

| Option | Written in the new record | Written in the old record | Old record's `Accepted` line |
| --- | --- | --- | --- |
| `-s`, `--supercede` | `Supercedes [N. Title](file)` | `Superceded by [N. Title](file)` | Removed |
| `-u`, `--update` | `Updates [N. Title](file)` | `Updated by [N. Title](file)` | Kept |

Note the tool's spelling, "Supercedes" and "Superceded by", when searching
records. Only a line reading exactly `Accepted` is removed. A status that has
been edited to something else, such as `Proposed` or `Deprecated`, stays beside
the link.

Never write link lines by hand. Use `-s` and `-u` on `adr new`, or `adr link`
afterwards, so that both records agree.

## Templates

The template in use is, in order of preference:

1. the file named by `$ADR_TEMPLATE`;
2. `templates/template.md` inside the log directory, which `adr config` reports
   as `adr_template`;
3. the built-in template, below.

The built-in template follows Michael Nygard's format:

```markdown
# NUMBER. TITLE

Date: DATE

## Status

STATUS

## Context

The issue motivating this decision, and any context that influences or constrains the decision.

## Decision

The change that we're proposing or have agreed to implement.

## Consequences

What becomes easier or more difficult to do and any risks introduced by the change that will need to be mitigated.
```

A project's template may have other sections. Fill in whatever sections the
template in use has, by reading the generated file after `adr new`.

## Numbering

A new record's number is one more than the highest number in the log
directory. Numbers are never reused while the records are kept, but deleting
the highest-numbered record means the next `adr new` reuses its number.

## Environment

| Variable | Effect |
| --- | --- |
| `VISUAL`, `EDITOR` | Editor opened on a new record. Clear both. |
| `ADR_TEMPLATE` | Template for new records. |
| `ADR_DATE` | Date written into a new record, as given. |
| `ADR_PAGER`, `PAGER` | Pager for help. Set to `cat`. |

The log directory is `docs/adr`, unless a `.adr-dir` file in the project root
names another. The working directory and each of its ancestors are searched in
turn for a `.adr-dir` file, then for a `docs/adr` directory, so the tool works
from anywhere in the project. Run `adr config` rather than assuming.

## Errors

The tool exits `0` on success and `1` on any failure, with a message on stderr
that begins `adr:`. Report the message to the user rather than working around
it.
