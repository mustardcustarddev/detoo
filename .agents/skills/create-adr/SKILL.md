---
name: create-adr
description: Creates Architecture Decision Records (ADRs) with the adr command-line tool, filling in the log's template strictly from evidence - what the user has said, the code, and documents in the repository - and interviewing the user for anything that cannot be evidenced, rather than inventing context, alternatives or justifications. Every record it creates has the status Draft. Use when the user asks to record, write, document or add an architecture decision or ADR, to supersede or update an existing ADR, or to start an ADR log with adr init.
license: Apache-2.0
compatibility: Requires the adr command-line tool on PATH.
---

# Create an ADR

Record one architecture decision in the project's ADR log, using the `adr`
tool, so that a future developer can understand it and judge whether it still
holds.

## The rule that matters most

**Write nothing you cannot evidence.** Every claim in the record must come
from the user, in this conversation; from the repository, which you have read;
or from the output of a command you ran. Context, alternatives, reasons,
numbers, dates and consequences are not to be invented, inferred as if known,
or embellished.

When you want to make a claim and have no evidence for it, **stop and ask the
user**. Ask through your question tool (`AskUserQuestion` in Claude Code), and
wait for the answer. [asking.md](references/asking.md) says how to ask an open
question with a tool that needs options, and when to ask in the chat instead.
Do not carry on with an assumed answer. If the user does not know, state the
gap plainly, or leave the point out; never fill it with something plausible.

Read [writing-adrs.md](references/writing-adrs.md) before you draft. It
explains what counts as evidence, and how to write each section.

## Every record you create is a Draft

The tool writes `Draft`; never change it. Whether a decision is accepted is for
people to record, by editing the line themselves once they have reviewed the
record, even when the user tells you the decision is already agreed.

## Workflow

### 1. Find the log

```bash
adr config        # adr_dir and adr_template
adr list          # existing records
```

If there is no log yet (`adr list` fails or prints nothing, and `adr_dir` does
not exist), ask the user whether to start one, and where. `adr init [DIR]`
creates the directory and writes record 1, the decision to keep ADRs.

Read the template in use: the file in `adr_template` if it is set, otherwise
the built-in one in [adr-tool.md](references/adr-tool.md#templates). Read the
existing records' titles, and any whose subject is close to this decision.

### 2. Check that it belongs in the log

Judge whether the decision is architectural, using
[significance.md](references/significance.md). If it looks like a local design
choice or a trivial one, tell the user why, in a sentence or two, and ask
whether to record it anyway. It is their log; your job is to make the choice
an informed one.

Look for existing records this one would supersede (replace) or update
(change substantively, while the old one stands). Confirm any such relation
with the user.

### 3. Gather the facts

Start the fact sheet described in
[writing-adrs.md](references/writing-adrs.md#facts-first). Collect what is
already evidenced:

* what the user has said in this conversation;
* the code and configuration the decision concerns: read it, and measure what
  you would otherwise estimate;
* related documents and existing ADRs.

Mark everything the template needs that you do not have as **Unknown**.

### 4. Interview the user

Ask about every Unknown. Each round is one call to your question tool, of a
few related questions, in the order the record needs them:

1. The problem: what made a decision necessary, and why then.
2. The constraints, and the forces in tension.
3. The alternatives that were actually considered, if any, and why each was
   or was not chosen.
4. Exactly what was decided.
5. The consequences: what becomes easier, what harder, and what risks.
6. The date, and relations to existing records.

Ask open questions about reasons. Offer as options only answers you have
evidence for, such as the date or which record is superseded; an open question
still goes through the question tool, with the answer typed in, as
[asking.md](references/asking.md) describes. Do not offer guessed reasons as
choices: a user picking the nearest option is not evidence.

If an answer is vague, ask a follow-up rather than sharpening it yourself. If
the user asks you to fill the gaps yourself, explain that the record must be
true to be useful, and offer to write only what is known, with the gaps stated.

### 5. Draft, and get approval

Write the draft in the conversation, in the template's sections. After it,
list the sources: each claim, and where it came from. Ask the user to approve
it or correct it. Revise until they approve.

Do not create the file before approval: `adr new` takes the next number and,
with `-s` or `-u`, changes other records too.

### 6. Create the record

```bash
VISUAL= EDITOR= adr new "Title of the decision"
```

Clearing `VISUAL` and `EDITOR` stops the tool waiting on an editor. It prints
the path of the new file. Set `ADR_DATE=YYYY-MM-DD` if the user chose a date
other than today.

If the record supersedes or updates an existing one, tell the user that the
tool records the link in **both** records at once. With `-s`, the old record
is marked "Superseded by" the new one, and loses its `Accepted` or `Draft`
line, although the new record is only a draft. Ask whether to:

* **link now**: `VISUAL= EDITOR= adr new -s 4 "Title"` (or `-u 4`); or
* **link on acceptance**: name the old record in the text of the new one, and
  have whoever accepts it run `adr link -s NEW OLD` (or `-u`) then.

Then edit that file:

* Keep the lines the tool wrote: the title, the date, `Draft`, and any link
  lines in the Status section.
* Replace the guidance text under each heading with the approved content.

Never write link lines by hand; use the tool:

```bash
adr link -s SOURCE TARGET   # SOURCE supersedes TARGET
adr link -u SOURCE TARGET   # SOURCE updates TARGET
```

See [adr-tool.md](references/adr-tool.md) for every command and its effects.

### 7. Finish

* Read the file back, and check it against the list at the end of
  [writing-adrs.md](references/writing-adrs.md#before-you-finish).
* Run `adr list` to confirm the record is in the log.
* If the project keeps a generated table of contents (look for one made by
  `adr list --toc`, such as a `README.md` in the log directory), offer to
  regenerate it.
* Tell the user the path, that the record is a draft for them to accept by
  changing `Draft` to `Accepted`, and any records the tool changed. Leave
  committing to them.

## Changing an existing decision

An accepted ADR is history: do not rewrite what it decided. Fixing a typo or
a broken link is fine. To change or reverse a decision, write a new record
that supersedes or updates the old one, following this workflow.

## Edge cases

* **Recording a decision made long ago.** The code may show what was done, and
  `git log` when; only the user can say why. Ask, and ask which date to use.
* **Several decisions at once.** Suggest one record each, and confirm the
  split with the user.
* **The `adr` command is missing.** Tell the user. Do not create records by
  hand, since the numbering and links would not be maintained.
* **The tool reports an error.** Show the user the message, which begins
  `adr:`, rather than working around it.
