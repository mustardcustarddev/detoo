# The plm tool

`plm` keeps a log of numbered, linked Markdown problem playbooks in the
repository. It replaced `pse`, which kept them in `docs/piosee`, named
`PIOSEE-`. This is what an agent needs to drive it without a terminal of its
own.

## Running it non-interactively

Two things can block an agent:

* `plm new` opens the new playbook in `$VISUAL` or `$EDITOR`, and waits for the
  editor to close. With neither set, it skips the editor and only prints the
  path. **Always clear both:** `VISUAL= EDITOR= plm new "Title"`.
* `plm help` prints to stdout, unless `$PLM_PAGER` names a pager, which it
  runs on a terminal. `$PAGER` is ignored. **If it is set, clear it:**
  `PLM_PAGER= plm help new`.

Options come before operands. An operand that begins with a hyphen must follow
`--`. A title is a single argument, so quote it.

## Commands

* `plm config`: prints `plm_dir="…"` and `plm_template="…"`: the log
  directory, relative to the working directory, and the project's template in
  the log directory, or empty if there is none. Does not report
  `$PLM_TEMPLATE`. Works before the log exists.
* `plm list`: prints the path of every playbook, one a line, in order.
  Prints nothing if there are none. Fails with
  `plm: the docs/problem directory does not exist` if there is no log.
* `plm init [DIRECTORY]`: starts a log in `DIRECTORY`, default `docs/problem`,
  and writes playbook 1, "Plan for incidents with PIOSEE". Does nothing if the
  log already has playbooks.
* `plm new [-r REF]... [-u REF]... TITLE`: creates the next numbered playbook
  from the template, in draft, and prints its path.
* `plm link (-r | -u) SOURCE TARGET`: records, in both playbooks, that
  `SOURCE` replaces or updates `TARGET`. Exactly one of `-r` and `-u`.
* `plm list --toc [-i INTRO] [-o OUTRO] [-p PREFIX]`: writes a Markdown table
  of contents to stdout, headed `Problem Playbooks`.
* `plm graph [--mermaid | --graphviz] [-p PREFIX] [-e EXT]`: draws the
  playbooks and the links between them.

A reference (`REF`, `SOURCE`, `TARGET`) is a playbook's number, its path,
or part of its file name. Digits alone are a number: `2` finds playbook 2 or
nothing, never playbook 12.

## What the tool writes

`plm new` copies the template and replaces the first `NUMBER`, `TITLE`,
`DATE` and `STATUS` on each line. The file name is `PROBLEM-`, the four-digit
number and a slug of the title, such as
`docs/problem/PROBLEM-0004-unexpected-charge-after-deletion.md`.

* `DATE` is today in ISO 8601 (`2026-10-16`), or `$PLM_DATE` exactly as given.
* `STATUS` is `Draft`.
* Everything else in the template, including its guidance text, is copied as
  it is. The guidance text must be replaced by the playbook's content.

Only files named `PROBLEM-NNNN-….md` are playbooks, along with playbooks named
as `pse` named them (`PIOSEE-NNNN-….md`), or before the prefix was upper
case. Anything else in the log directory, such as a `README.md`, is ignored.
`detoo doctor` reports playbooks with the older names, and a log left in
`docs/piosee`, and `detoo doctor --fix` moves and renames them.

## Links

Links are lines in the `## Status` section, written in both playbooks:

```markdown
## Status

Draft

Replaces [4. Unexpected charge after deletion](PROBLEM-0004-unexpected-charge-after-deletion.md)
```

| Option            | Written in the new playbook | Written in the old playbook    | Old playbook's `Draft` line |
| ----------------- | --------------------------- | ------------------------------ | --------------------------- |
| `-r`, `--replace` | `Replaces [N. Title](file)` | `Replaced by [N. Title](file)` | Removed                     |
| `-u`, `--update`  | `Updates [N. Title](file)`  | `Updated by [N. Title](file)`  | Kept                        |

Only a line reading exactly `Draft` is removed. Any other status, such as
`Agreed`, stays beside the link.

Use `-r` when a new playbook is a new plan for the problem an earlier one
planned for, such as a rewrite after the earlier one was used and found
wanting. Use `-u` when a new playbook changes part of an earlier one without
replacing it, such as a case the earlier one left out.

> **The link takes effect at once.** `plm new -r 4 "Title"` marks playbook 4
> "Replaced by" the new playbook, and removes playbook 4's `Draft` line, the
> moment the new playbook is created. `-u` adds "Updated by" to the old
> playbook at once, too.

Never write link lines by hand. Use `-r` and `-u` on `plm new`, or `plm link`
afterwards, so that both playbooks agree.

## Templates

The template in use is, in order of preference:

1. the file named by `$PLM_TEMPLATE`, an absolute path or one relative to
   the working directory;
2. `templates/template.md` inside the log directory, which `plm config`
   reports as `plm_template`;
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

## Immediate action

The action to take at once, before gathering information, to make the situation safe or stop it getting worse, and who takes it, or the standard operating procedure (SOP) that sets it out. If there is none, write "No immediate action is indicated."

## Counter-indicated actions

Optional. The actions known to make the situation worse, or to waste time, that someone under pressure may be tempted to take, and why each is wrong. Delete this section if none are known.

## Information

The facts to gather before deciding: what to collect, where each is found, and who owns it, or the SOP that says how to collect it. Separate what is known from what is assumed, and say how to settle the critical unknowns quickly.

## Options

The realistic courses of action, distinct and actionable, including at least one conservative fallback. For each, when it applies and its tradeoffs.

## Selection

How to choose an option: the criteria, which option each finding of the investigation leads to, the default, and the confidence and assumptions behind the rule.

## Execution

The steps to carry out the option chosen, in order: who takes each, what each depends on, what to record and tell whom, and how to roll back. Refer to an SOP for a routine procedure, rather than repeating its steps.

## Evaluation

How to tell, after the playbook is used, whether it worked: the measurable outcomes, the SOP for any routine check, when and by whom they are reviewed, and what to do if they are missed.
```

The **Counter-indicated actions** section is optional: delete it when none
are known. A playbook links to an SOP by a relative link to its file, such
as `[SOP 3. Roll back a release](../sop/SOP-0003-roll-back-a-release.md)`;
`plm` does not check these links, so check them yourself.

A project's template may have other sections, or follow another method
altogether. Read the template in use for every playbook, rather than
assuming it is this one.

## Numbering

A new playbook's number is one more than the highest number in the log
directory. Numbers are never reused while the playbooks are kept, but
deleting the highest-numbered playbook means the next `plm new` reuses its
number.

## Environment

| Variable             | Effect                                       |
| -------------------- | -------------------------------------------- |
| `VISUAL`, `EDITOR`   | Editor opened on a new playbook. Clear both. |
| `PLM_TEMPLATE`       | Template for new playbooks.                  |
| `PLM_DATE`           | Date written into a new playbook, as given.  |
| `PLM_PAGER`          | Pager for help, if set. Clear it.            |

The log directory is `docs/problem`, unless a `.plm-dir` file in the project
root names another. The working directory and each of its ancestors are
searched in turn for a `.plm-dir` file, then for a `docs/problem` directory,
so the tool works from anywhere in the project. Run `plm config` rather than
assuming.

## Errors

The tool exits `0` on success and `1` on any failure, with a message on stderr
that begins `plm:`. Report the message to the user rather than working around
it.
