---
name: create-rfc
description: Drafts Requests for Comments (RFCs) with the rfc command-line tool. Turns a short request, such as "create an RFC to migrate the system to hexagonal architecture", into a draft others can comment on, by reading the project's RFC template afresh, gathering facts from the code and documents, and interviewing the user section by section for their real intent and reasoning, rather than guessing. Every RFC it creates has the status Draft. Use when the user asks to write, draft, propose, create or start an RFC or request for comments, to put a proposed change out for feedback, to obsolete or update an existing RFC, or to start an RFC log with rfc init.
license: Apache-2.0
compatibility: Requires the rfc command-line tool on PATH.
---

# Create an RFC

An RFC proposes a change to the people it affects, before it is made, so that
they can comment on it. A good one says exactly what is proposed, why, what it
costs, what else could be done, and what the author wants reviewers to settle.
Only the author knows most of that. Your job is to draw it out of them, and to
write it down faithfully: you are the author's interviewer and scribe, not the
author.

## The rules that matter most

1. **Write nothing you cannot evidence.** Every claim in the RFC must come
   from the user, in this conversation; from the repository, which you have
   read; or from the output of a command you ran. Problems, goals, scope,
   design choices, alternatives, drawbacks, estimates and predictions are not
   to be invented, inferred as if known, or embellished.
2. **A one-line request is a theme, not a brief.** "Create an RFC to migrate
   the system to hexagonal architecture" does not say what problem it solves,
   how far it goes, in what order, at what cost, or what was weighed against
   it. Do not draft from it. Gather the facts, then interview the user.
3. **Every RFC you create is a Draft.** The tool writes `Draft`; never change
   it. Accepting, rejecting or withdrawing an RFC is the outcome of the
   comment process, and is recorded by people.

When you want to make a claim and have no evidence for it, **stop and ask the
user**. Ask through your question tool (`AskUserQuestion` in Claude Code), and
wait for the answer. [asking.md](references/asking.md) says how to ask an open
question with a tool that needs options, and when to ask in the chat instead.
Do not carry on with an assumed answer.

Before you start, read:

* [interviewing.md](references/interviewing.md): how to plan and run the
  interview, with a worked example;
* [writing-rfcs.md](references/writing-rfcs.md): what counts as evidence, what
  to do when a fact is missing, and how to write each section.

## Workflow

### 1. Find the log

```bash
rfc config        # rfc_dir and rfc_template
rfc list          # existing RFCs
```

If there is no log yet (`rfc list` fails with `the docs/rfc directory does not
exist`), ask the user whether to start one. `rfc init` creates `docs/rfc` and
writes RFC 1, the decision to propose significant changes in RFCs.

Read RFC 1, or whichever RFC sets out the project's process: it may say what
needs an RFC, who comments, and where. Follow it. Read the titles of the
existing RFCs, and in full any whose subject is close to this proposal. If the
project keeps an ADR log (`adr list`), read the records the proposal touches:
an RFC often proposes to change a recorded decision.

### 2. Read the template, every time

Read the template in use, afresh, for every RFC: the file named by
`$RFC_TEMPLATE` if it is set (`rfc config` does not report it), otherwise the
file in `rfc_template` if it is set, otherwise the built-in one in
[rfc-tool.md](references/rfc-tool.md#templates). Do not work from memory of
the built-in template or of an earlier RFC; projects change their templates.

List its sections, and the guidance under each. These are what the interview
must fill.

### 3. Check that it needs an RFC

* An RFC is for a change that others would want a say in, before it is made.
  If the change is small and local, say so in a sentence or two, and ask
  whether to go ahead anyway. It is the user's call.
* If the change is already decided, or already made, it is a decision to
  record rather than a proposal to comment on: suggest an ADR instead, with
  the `create-adr` skill if it is available.
* If the request holds several independent proposals, suggest one RFC each,
  and confirm the split.
* Look for existing RFCs this one would obsolete (replace) or update (change
  substantively). Confirm any such relation with the user.

### 4. Gather the facts

Before asking anything, learn what the repository can tell you, so that you
ask only what you could not read, and can ask precisely:

* the current state of what the proposal would change: its structure, its
  dependencies, the code involved. Measure what you would otherwise estimate,
  with commands, and note them;
* related RFCs, ADRs and documents, and anything the user points you to;
* what the user has already said.

Read in proportion to the change: a proposal to restructure the whole system
needs a survey of the whole system; one to change a single interface needs
that interface and its callers.

The current state is yours to establish. The proposal, its reasons and its
priorities are the user's. Start the fact sheet described in
[writing-rfcs.md](references/writing-rfcs.md#facts-first), and mark everything
the template needs that you do not have as **Unknown**.

### 5. Plan the interview

Go through the template's sections, and under each write the questions this
proposal raises, using the question bank and theme lenses in
[interviewing.md](references/interviewing.md). Make them specific to the theme
and to what you found. Cross off what the facts already answer. Keep the plan
in your working notes, and revise it as answers come in.

### 6. Interview the user

Ask in rounds, a few related questions at a time, in an order where each
builds on the last. Each round is one call to your question tool. For the
built-in template:

1. Intent and scope: what exactly is proposed, what the theme's terms mean
   here, the end state, what is in and out, and who should comment.
2. Motivation: the problem, with an example; who feels it; why now; the cost
   of doing nothing.
3. Design: what changes and how; the path and the transition; what stays the
   same; what is undecided.
4. Drawbacks: costs, risks, expected objections.
5. Alternatives: what else was weighed, including doing nothing.
6. Unresolved questions: what reviewers should settle, before acceptance or
   later; and relations to existing RFCs.

Ask open questions about intent and reasons. Offer as options only answers you
have evidence for; an open question still goes through the question tool, with
the answer typed in, as [asking.md](references/asking.md) describes. Never
offer guessed reasons as choices: a user picking the nearest option is not
evidence. You may put an observation from the code to the user and ask whether
it matters; their answer is the evidence.

After each round, play back in two or three sentences what you understood, so
the user can correct it. If an answer is vague, ask a follow-up rather than
sharpening it yourself.

When the user cannot answer a point, offer the choices in
[writing-rfcs.md](references/writing-rfcs.md#when-a-fact-is-missing): make it
an unresolved question for reviewers, state the gap, or leave it out. What is
proposed, and why, cannot be deferred: without them, the RFC is not ready to
draft. Say so.

If the user asks you to fill the gaps yourself, explain that reviewers will
comment on what the RFC says, so it must say what the author thinks; offer to
write only what is known, with the open points as unresolved questions. If the
user asks for your ideas, give them as labelled suggestions; see
[writing-rfcs.md](references/writing-rfcs.md#your-own-suggestions).

Stop when every section of the template can be written from the fact sheet.

### 7. Draft, and get approval

Write the draft in the conversation, in the template's sections. After it,
list the sources: each claim, and where it came from. Ask the user to approve
it or correct it. Revise until they approve.

Do not create the file before approval: `rfc new` takes the next number and,
with `-o` or `-u`, changes other RFCs too.

### 8. Create the RFC

```bash
VISUAL= EDITOR= rfc new "Title of the proposal"
```

Clearing `VISUAL` and `EDITOR` stops the tool waiting on an editor. It prints
the path of the new file.

If the proposal obsoletes or updates an existing RFC, tell the user that the
tool records the link in **both** RFCs at once. With `-o`, the old RFC is
marked "Obsoleted by" the new one, and loses its `Draft` line, although the
new RFC is only a draft that may never be accepted. Ask whether to:

* **link now**: `VISUAL= EDITOR= rfc new -o 4 "Title"` (or `-u 4`); or
* **link on acceptance**: name the old RFC in the text of the new one, and
  have whoever accepts it run `rfc link -o NEW OLD` (or `-u`) then.

Then edit the new file in place:

* Keep the lines the tool wrote: the title, the date, `Draft`, and any link
  lines in the Status section. Do not rewrite the whole file, and never
  retype the title: its number is the tool's, and must stay the number in
  the file name.
* Replace the guidance text under each heading with the approved content,
  one section at a time.

Never write link lines by hand. See [rfc-tool.md](references/rfc-tool.md) for
every command and its effects.

### 9. Finish

* Read the file back, and check it against the list at the end of
  [writing-rfcs.md](references/writing-rfcs.md#before-you-finish).
* Check that the number in the title, `# N.`, is the number in the file
  name, `RFC-000N-...`. If they differ, put the tool's title line back.
* Run `rfc list` to confirm the RFC is in the log.
* If the project keeps a generated table of contents (look for one made by
  `rfc list --toc`, such as a `README.md` in the log directory), offer to
  regenerate it.
* Tell the user the path, that the RFC is a draft, and any RFCs the tool
  changed. Its next step is the project's comment process, such as a pull
  request if RFC 1 says so. Leave committing, and opening the pull request, to
  them.

## Revising an RFC

* **A Draft** is meant to change in response to comments. Revise it in place,
  to the same standard of evidence: changes come from the user, or from
  comments they show you and ask you to address. It stays a Draft.
* **An accepted, rejected or withdrawn RFC** is history: do not rewrite what
  it proposed. Fixing a typo or a broken link is fine. To change course, write
  a new RFC that obsoletes or updates it, following this workflow.

## Edge cases

* **The proposal is already agreed.** It is still created as a Draft. Tell the
  user they can accept it through their process; if the work is already done,
  suggest an ADR instead.
* **The user wants a quick draft now and the details later.** Draft what is
  known, with every open point as an unresolved question, provided what is
  proposed and why came from the user.
* **The `rfc` command is missing.** Tell the user. Do not create RFCs by hand,
  since the numbering and links would not be maintained.
* **The tool reports an error.** Show the user the message, which begins
  `rfc:`, rather than working around it.
