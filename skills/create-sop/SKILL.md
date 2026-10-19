---
name: create-sop
description: Writes a standard operating procedure (SOP) with the sop command-line tool - the steps of a routine task, such as rolling back a release, rotating a key or restoring a backup, as an ordered list. Unlike the other detoo skills, it works from research rather than interview - it turns a short request, such as "write an SOP for rolling back a release", into a draft by reading the codebase (scripts, Makefiles, CI and deployment configuration, runbooks) and the official documentation of the tools involved, and asks the user only what research cannot settle, and to approve the draft. Every SOP it creates is a Draft. Use it to write, draft or document an SOP or a "how do we..." procedure, when a problem playbook needs a routine procedure that has no SOP, to supersede or update an SOP, or to start a log with sop init.
license: Apache-2.0
compatibility: Requires the sop command-line tool on PATH. Reading third-party documentation needs a web fetch or search tool.
---

# Create a standard operating procedure

A standard operating procedure (SOP) is the steps of a routine task, such as
rolling back a release, as an ordered list. It has no decision in it: it is
done the same way each time, by whoever is on hand, so it must be exact.

Problem playbooks, written with the `create-playbook` skill and the `plm`
tool, refer to SOPs for their routine steps: the immediate action, how to
collect the information, the steps of execution, and the checks of
evaluation. An SOP is written once, kept current once, and every playbook
that refers to it stays current with it.

Most SOPs can be found rather than asked for. The steps of a deployment are
in the Makefile, the CI workflow and the deployment configuration; the steps
through a provider's dashboard are in its documentation. So this skill works
differently from the other detoo skills: **research first, draft, and ask
last**. The user should be able to ask for an SOP in a sentence, and get a
sourced draft back with few questions or none.

## The rules that matter most

1. **Research, don't interview.** Read the code and the documentation
   before asking anything. Ask the user only what research cannot settle:
   sources that disagree, access or environments only the team knows, and
   whether to create the draft. Never ask the user to dictate the steps.
2. **Write nothing you cannot evidence.** Every step comes from one of these:
   * the repository: code, scripts, Makefiles, CI and deployment
     configuration, runbooks and documents, which you have read;
   * the official documentation of a tool or service involved, for the
     version the project uses, which you have read;
   * the output of a read-only command you ran, such as a tool's `--help`;
   * the user, in this conversation.

   What you remember of a tool, or think is common practice, is not a
   source. A command, a flag, a menu path or an environment name from memory
   sends the reader the wrong way when they can least afford it.
3. **Gaps are open questions, not guesses.** A step research cannot settle
   goes into the draft as an open question, where the reader will meet it.
   It does not stop the draft being written.
4. **Never run the procedure.** Read it; do not carry it out to test it. Run
   read-only commands only, such as `sop list`, `git log` or a tool's
   `--help`. Do not deploy, roll back, restart, delete, call a system's API
   or query a live system, unless the user asks you to.
5. **No decisions.** An SOP may check that a step worked, and say to stop
   and escalate if not. A choice between courses of action belongs in a
   problem playbook, which refers to the SOP.
6. **Every SOP you create is a Draft.** The tool writes `Draft`; never change
   it. Whether an SOP is approved is for the team to record.

When you must ask, ask through your question tool (`AskUserQuestion` in
Claude Code), and wait for the answer. [asking.md](references/asking.md) says
how to ask an open question with a tool that needs options, and when to ask
in the chat instead.

Before you start, read [writing-sops.md](references/writing-sops.md): what to
research, what counts as evidence, and how to write the steps.

## Workflow

### 1. Find the log

```bash
sop config        # sop_dir and sop_template
sop list          # existing SOPs
```

If there is no log yet, there is none when `sop list` fails with
`the docs/sop directory does not exist`. Ask the user whether to start one
with `sop init`, which creates `docs/sop` and writes SOP 1, which sets out how
to write an SOP.

Read the titles of the existing SOPs. Read in full any whose task is close to
this one: the new SOP may duplicate it, supersede it or update it. Read one or
two others for the house style.

### 2. Read the template, every time

Read the template in use, afresh, for every SOP. It is the file named by
`$SOP_TEMPLATE` if that is set (`sop config` does not report it). Otherwise it
is the file in `sop_template`, if that is set. Otherwise it is the built-in
one in [sop-tool.md](references/sop-tool.md#templates): a title, a date, a
status, and the steps. Do not work from memory of it; projects change their
templates.

### 3. Check that it is an SOP

* **A routine task**, done the same way each time, is an SOP. Its title is
  the task, as an instruction: "Roll back a release".
* **A task that needs a decision**, such as whether to roll back or fix
  forward, is a problem playbook: suggest the `create-playbook` skill, if it
  is available. Offer to write SOPs for its routine steps.
* **Several tasks** need several SOPs. List them, and confirm which to write,
  in one question.
* **A task done differently in each environment** may be one SOP with a step
  per environment, or one SOP for each. Follow the house style, or ask.

### 4. Research

Work through [writing-sops.md](references/writing-sops.md#where-to-look), and
build the fact sheet it describes: each step, its source, and anything left
unknown.

1. **Find how the project does the task.** Read the Makefile, scripts, CI
   workflows, deployment configuration, READMEs, runbooks and documents
   that touch it. Find the exact commands, targets, arguments, environments
   and paths.
2. **Find what the tools and services require.** Read the official
   documentation of each one involved, for the version the project uses, such
   as a version pinned in a lock file or configuration. Run a tool's `--help`
   where it is installed.
3. **Find the preconditions and the checks.** What access, tools and
   approvals a step needs; how to tell that each step worked; and what to do
   if it did not, as far as a source says.
4. **Find who uses it.** Read any playbook or SOP that refers to this task,
   with `plm list` and `sop list`, so that the SOP fits how they use it.

### 5. Draft

Write the draft in the conversation, in the template's sections, from the
fact sheet alone. Write each step as [writing-sops.md](references/writing-sops.md#the-steps)
describes. Where research left a gap, write the step as an open question:
"Open question: the name of the production cluster is not recorded in the
repository."

After the draft, list the sources: each step, and where it came from.

### 6. Ask once

Ask the user, in one call of your question tool:

* **whether to create the draft**: **Create** it as it is, with any open
  questions left in for the team, or **Correct** it first, with their
  correction typed in as Other;
* and, in the same call, up to three questions research could not settle,
  each specific and citing what you did find: sources that disagree, an
  environment or access only the team knows. Leave anything else as an open
  question in the draft.

If the request came from the `create-playbook` skill, or the user asked for
the SOP without review, still show the draft, and ask the one question
whether to create it.

Revise for any correction, and play it back if it changes more than a
detail.

### 7. Create the SOP

```bash
VISUAL= EDITOR= sop new "Title of the SOP"
```

Clearing `VISUAL` and `EDITOR` stops the tool waiting on an editor. It prints
the path of the new file.

A new SOP may supersede or update an existing one. If so, tell the user that
the tool records the link in **both** SOPs at once, and that with `-s` the old
SOP loses its `Draft` or `Approved` line. Then:

* **supersede**: `VISUAL= EDITOR= sop new -s 4 "Title"`, for a new procedure
  for the task SOP 4 set out;
* **update**: `VISUAL= EDITOR= sop new -u 4 "Title"`, for an SOP that changes
  part of SOP 4, such as a step for a case it left out.

Then edit the new file in place: keep the lines the tool wrote (the title,
the date, `Draft`, and any link lines in the Status section), and replace the
guidance text with the approved steps. Do not rewrite the whole file, and
never retype the title: its number is the tool's, and must stay the number in
the file name. Never write link lines by hand. See
[sop-tool.md](references/sop-tool.md) for every command and its effects.

### 8. Finish

* Read the file back, and check it against
  [writing-sops.md](references/writing-sops.md#before-you-finish).
* Check that the number in the title, `# N.`, is the number in the file
  name, `SOP-000N-...`. If they differ, put the tool's title line back.
* Run `sop list` to confirm the SOP is in the log.
* The project may keep a table of contents made by `sop list --toc`, such as
  a `README.md` in the log directory. If so, offer to regenerate it.
* Tell the user:
  * the path, and that the SOP is a draft;
  * any SOPs the tool changed;
  * the open questions left for the team;
  * that someone who does the task should follow it once, and correct it.

  Leave committing to them. If a playbook asked for the SOP, give it the
  path and title to link to.

## Revising an SOP

* **A Draft** is revised in place, to the same standard of evidence.
* **When a fact changes**, such as a renamed target or a new flag, research
  it again, and propose the correction in place. Git keeps the history.
* **When the procedure changes**, such as a new deployment tool, offer a new
  SOP that supersedes the old one (`-s`), or one that updates it (`-u`) for
  a case it left out.
* **After the SOP has been used**, ask what was missing, wrong or out of
  order, research it, and revise as above.

## Edge cases

* **Research finds nothing**: no script, no configuration, no documentation.
  Say so, with where you looked. Ask the user for a document or a person to
  ask, or whether to interview them for the steps. Do not invent them.
* **The sources disagree**, such as a README and the CI workflow. Show both,
  with their sources, and ask which holds. The configuration that runs is
  usually the better evidence; say so, but let the user choose.
* **You cannot read the documentation**, because you have no web tool, or
  the pages are not public. Say so, and ask for the documentation or a link,
  or leave the step as an open question.
* **The user says the SOP is approved.** It is still created as a Draft. Tell
  them they can mark it approved through their process.
* **The `sop` command is missing.** Tell the user. Do not create SOPs by hand,
  since the numbering and links would not be maintained.
* **The tool reports an error.** Show the user the message, which begins
  `sop:`, rather than working around it.
