# The pse tool

`pse` keeps a log of numbered, linked Markdown PIOSEE playbooks in the
repository. This is what an agent needs to drive it without a terminal of its
own.

## Running it non-interactively

Two things can block an agent:

* `pse new` opens the new playbook in `$VISUAL` or `$EDITOR`, and waits for the
  editor to close. With neither set, it skips the editor and only prints the
  path. **Always clear both:** `VISUAL= EDITOR= pse new "Title"`.
* `pse help` prints to stdout, unless `$PSE_PAGER` names a pager, which it
  runs on a terminal. `$PAGER` is ignored. **If it is set, clear it:**
  `PSE_PAGER= pse help new`.

Options come before operands. An operand that begins with a hyphen must follow
`--`. A title is a single argument, so quote it.

## Commands

* `pse config`: prints `pse_dir="…"` and `pse_template="…"`: the log
  directory, relative to the working directory, and the project's template in
  the log directory, or empty if there is none. Does not report
  `$PSE_TEMPLATE`. Works before the log exists.
* `pse list`: prints the path of every playbook, one a line, in order.
  Prints nothing if there are none. Fails with
  `pse: the docs/piosee directory does not exist` if there is no log.
* `pse init [DIRECTORY]`: starts a log in `DIRECTORY`, default `docs/piosee`,
  and writes playbook 1, "Plan for incidents with PIOSEE". Does nothing if the
  log already has playbooks.
* `pse new [-r REF]... [-u REF]... TITLE`: creates the next numbered playbook
  from the template, in draft, and prints its path.
* `pse link (-r | -u) SOURCE TARGET`: records, in both playbooks, that
  `SOURCE` replaces or updates `TARGET`. Exactly one of `-r` and `-u`.
* `pse list --toc [-i INTRO] [-o OUTRO] [-p PREFIX]`: writes a Markdown table
  of contents to stdout, headed `PIOSEE Playbooks`.
* `pse graph [--mermaid | --graphviz] [-p PREFIX] [-e EXT]`: draws the
  playbooks and the links between them.

A reference (`REF`, `SOURCE`, `TARGET`) is a playbook's number, its path,
or part of its file name. Digits alone are a number: `2` finds playbook 2 or
nothing, never playbook 12.

## What the tool writes

`pse new` copies the template and replaces the first `NUMBER`, `TITLE`,
`DATE` and `STATUS` on each line. The file name is `PIOSEE-`, the four-digit
number and a slug of the title, such as
`docs/piosee/PIOSEE-0004-unexpected-charge-after-deletion.md`.

* `DATE` is today in ISO 8601 (`2026-10-16`), or `$PSE_DATE` exactly as given.
* `STATUS` is `Draft`.
* Everything else in the template, including its guidance text, is copied as
  it is. The guidance text must be replaced by the playbook's content.

Only files named `PIOSEE-NNNN-….md` are playbooks, along with playbooks named
before the prefix was upper case (`piosee-NNNN-….md`). Anything else in the
log directory, such as a `README.md`, is ignored. `detoo doctor` reports
playbooks with the older names, and `detoo doctor --fix` renames them.

## Links

Links are lines in the `## Status` section, written in both playbooks:

```markdown
## Status

Draft

Replaces [4. Unexpected charge after deletion](PIOSEE-0004-unexpected-charge-after-deletion.md)
```

| Option            | Written in the new playbook | Written in the old playbook    | Old playbook's `Draft` line |
| ----------------- | --------------------------- | ------------------------------ | --------------------------- |
| `-r`, `--replace` | `Replaces [N. Title](file)` | `Replaced by [N. Title](file)` | Removed                     |
| `-u`, `--update`  | `Updates [N. Title](file)`  | `Updated by [N. Title](file)`  | Kept                        |

Only a line reading exactly `Draft` is removed. Any other status, such as
`Agreed`, stays beside the link.

Use `-r` when a new playbook is a new plan for the situation an earlier one
planned for, such as a rewrite after the earlier one was used and found
wanting. Use `-u` when a new playbook changes part of an earlier one without
replacing it, such as a case the earlier one left out.

> **The link takes effect at once.** `pse new -r 4 "Title"` marks playbook 4
> "Replaced by" the new playbook, and removes playbook 4's `Draft` line, the
> moment the new playbook is created. `-u` adds "Updated by" to the old
> playbook at once, too.

Never write link lines by hand. Use `-r` and `-u` on `pse new`, or `pse link`
afterwards, so that both playbooks agree.

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

## Scenario

The situation this playbook is for, in one to three sentences: what happens, or is reported, that calls for it.

## Problem

The decision to make when the scenario arises, as one clear question, and the main risks if it is decided wrongly. No solutions here.

## Information

The facts to gather before deciding: what to collect, where each is found, and who owns it. Separate what is known from what is assumed, and say how to settle the critical unknowns quickly.

## Options

The realistic courses of action, distinct and actionable, including at least one conservative fallback. For each, when it applies and its tradeoffs.

## Selection

How to choose an option: the criteria, which option each finding of the investigation leads to, the default, and the confidence and assumptions behind the rule.

## Execution

The steps to carry out the option chosen, in order: who takes each, what each depends on, what to record and tell whom, and how to roll back.

## Evaluation

How to tell, after the playbook is used, whether it worked: the measurable outcomes, when and by whom they are reviewed, and what to do if they are missed.
```

A project's template may have other sections. Read the template in use for
every playbook, rather than assuming it is this one.

## Numbering

A new playbook's number is one more than the highest number in the log
directory. Numbers are never reused while the playbooks are kept, but
deleting the highest-numbered playbook means the next `pse new` reuses its
number.

## Environment

| Variable             | Effect                                       |
| -------------------- | -------------------------------------------- |
| `VISUAL`, `EDITOR`   | Editor opened on a new playbook. Clear both. |
| `PSE_TEMPLATE`       | Template for new playbooks.                  |
| `PSE_DATE`           | Date written into a new playbook, as given.  |
| `PSE_PAGER`          | Pager for help, if set. Clear it.            |

The log directory is `docs/piosee`, unless a `.pse-dir` file in the project
root names another. The working directory and each of its ancestors are
searched in turn for a `.pse-dir` file, then for a `docs/piosee` directory,
so the tool works from anywhere in the project. Run `pse config` rather than
assuming.

## Errors

The tool exits `0` on success and `1` on any failure, with a message on stderr
that begins `pse:`. Report the message to the user rather than working around
it.
