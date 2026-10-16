---
name: create-piosee
description: Works a decision through PIOSEE (Problem, Information, Options, Select, Execute, Evaluate) and records it with the pse command-line tool. Turns a short request, such as "help us decide whether to roll back the release", into a draft record by reading the project's template afresh, gathering facts from the code and the systems, and interviewing the user for the problem, the options and their choice, rather than guessing or choosing for them. Every record it creates has the status Draft. Use when the user asks to make, work through, record or write up a decision with PIOSEE or pse, to decide what to do about an incident, an outage, a failed release or another problem that needs a decision, to evaluate a decision already recorded, to reverse or update one, or to start a PIOSEE log with pse init.
license: Apache-2.0
compatibility: Requires the pse command-line tool on PATH.
---

# Create a PIOSEE record

PIOSEE is a six-step method for making a decision under pressure, taught to
airline pilots: define the **Problem**, gather the **Information**, generate
the **Options**, **Select** one, **Execute** it, and **Evaluate** whether it
worked, returning to the Problem if it did not. A PIOSEE record is one
decision written down in those six steps, so that the people making it see
the same facts, and those who come later can see what was known, what was
weighed, what was done, and how it turned out.

The problem, the options that are open and the choice belong to the user and
their team. Your job is to help them work through the steps, and to write down
faithfully what they decide: you are the first officer and the scribe, not the
captain.

## The rules that matter most

1. **Time comes first.** Ask how long there is before anything else. If harm
   is happening now and the action that stops it is obvious, such as taking a
   failing server out of service, tell the user to act, and offer to record
   the decision afterwards. The record must never delay the action.
2. **Write nothing you cannot evidence.** Every claim in the record must come
   from the user, in this conversation; from the repository, which you have
   read; or from the output of a command you ran. Facts, options, risks,
   reasons, actions and outcomes are not to be invented, inferred as if
   known, or embellished.
3. **The user selects.** Never choose the option, and never write a choice
   the user did not make. You may compare the options, and say which you
   would choose and why, as a labelled suggestion, if the user asks.
4. **Every record you create is a Draft.** The tool writes `Draft`; never
   change it. What becomes of a decision is for the team to record.

When you want to make a claim and have no evidence for it, **stop and ask the
user**. Ask through your question tool (`AskUserQuestion` in Claude Code), and
wait for the answer. [asking.md](references/asking.md) says how to ask an open
question with a tool that needs options, and when to ask in the chat instead.
Do not carry on with an assumed answer.

Before you start, read:

* [piosee.md](references/piosee.md): what PIOSEE is, each step, and its
  strengths and weaknesses, from the sources;
* [interviewing.md](references/interviewing.md): how to plan and run the
  interview, step by step, with a worked example;
* [writing-piosees.md](references/writing-piosees.md): what counts as
  evidence, what to do when a fact is missing, and how to write each section.

## Workflow

### 1. Ask how long there is

Before reading anything, ask the user how long they have to decide, and what
happens if they do not. Then:

* **Harm now, obvious action**: rule 1. Tell them to act; record later.
* **Minutes**: keep to the steps, one short question at a time, and draft as
  you go. Gather only the facts that would change the choice.
* **Hours or more**, or a decision already made: follow the whole workflow.

### 2. Find the log

```bash
pse config        # pse_dir and pse_template
pse list          # existing records
```

If there is no log yet (`pse list` fails with `the docs/piosee directory does
not exist`, or names another missing directory), ask the user whether to start
one, and where. `pse init [DIR]` creates the directory and writes record 1,
which describes the practice of making decisions with PIOSEE.

Read the titles of the existing records, and in full any whose problem is
close to this one: this decision may reverse or update one of them.

### 3. Read the template, every time

Read the template in use, afresh, for every record: the file named by
`$PSE_TEMPLATE` if it is set (`pse config` does not report it), otherwise the
file in `pse_template` if it is set, otherwise the built-in one in
[pse-tool.md](references/pse-tool.md#templates). Do not work from memory of
the built-in template or of an earlier record; projects change their
templates.

### 4. Check that it is a PIOSEE decision

* PIOSEE is for one problem that needs a decision and an action: an
  incident, a failed release, a supplier who has let the team down, a
  deadline that cannot be met. The record lives through the decision, from
  the problem to the evaluation.
* A long-lived choice about how the system is built is an architecture
  decision: suggest an ADR, with the `create-adr` skill if it is available. A
  change others would want a say in first is a proposal: suggest an RFC, with
  `create-rfc`. Work for someone who benefits is a user story, with
  `create-user-story`. The user chooses.
* Establish which of three situations this is:
  * **a live decision**, not yet made: work it through with the user;
  * **a decision already made**, to be written up, such as after an
    incident: record what was known and weighed at the time;
  * **an existing record to evaluate**: go to step 10.

### 5. Gather the facts

Information is the step where you can help most, and where invention does the
most harm. Before asking, learn what the repository and the systems can tell
you about the problem: the code involved, recent changes, logs, metrics,
alerts, tickets, runbooks, and earlier records. Run read-only commands only;
anything that changes a system is an action, and actions are the user's to
select. Note where each fact came from, with the time for anything that
changes.

Start the fact sheet described in
[writing-piosees.md](references/writing-piosees.md#facts-first). Mark
everything the template needs that you do not have as **Unknown**.

### 6. Plan the interview

Go through the template's sections, and under each write the questions this
problem raises, using the question bank in
[interviewing.md](references/interviewing.md). Cross off what the facts
already answer.

### 7. Work through the steps with the user

Ask in rounds, one step at a time, in order. Each round is one call to your
question tool; under time pressure, one question a call. Do not move on until
the user agrees the step is done; returning to an earlier step when something
new comes to light is part of the method.

1. **Problem**: what is wrong, what it affects, how badly, and how long there
   is. Play it back in a sentence or two, and get the user's agreement.
2. **Information**: show the facts you found, with their sources, and ask
   what else they know, what is still unknown, and who else should be asked.
3. **Options**: ask what courses of action they see, including doing nothing
   and buying time. For each, the risks and benefits as they see them. If the
   user asks for more options, offer yours, labelled as suggestions.
4. **Select**: ask which option they choose, and why. Their reason is the
   evidence. If they ask what you would choose, give your view, labelled as
   yours, and still let them choose.
5. **Execute**: who does what, in what order, and what will show it is
   working.
6. **Evaluate**: how and when they will judge whether it worked, and what they
   will do if it did not.

For **a decision already made**, ask the same questions in the past tense,
about what was known and considered at the time. Record options that were
seen only afterwards as such, in Evaluate, or leave them out. A record written
after the event must show how the decision was made, not justify it.

When the user cannot answer a point, offer the choices in
[writing-piosees.md](references/writing-piosees.md#when-a-fact-is-missing).
The problem, the option selected, and the reason for it cannot be deferred:
without them there is no decision to record.

### 8. Draft, and get approval

Write the draft in the conversation, in the template's sections. After it,
list the sources: each claim, and where it came from. Ask the user to approve
it or correct it, and revise until they do.

Do not create a file before approval: `pse new` takes the next number and,
with `-r` or `-u`, changes other records too. Under time pressure, the
approval can be quick, but it is still the user's.

### 9. Create the record

```bash
VISUAL= EDITOR= pse new "Title of the decision"
```

Clearing `VISUAL` and `EDITOR` stops the tool waiting on an editor. It prints
the path of the new file.

If the decision reverses or updates an earlier one, tell the user that the
tool records the link in **both** records at once, and that with `-r` the old
record loses its `Draft` line. Then:

* **reverse**: `VISUAL= EDITOR= pse new -r 4 "Title"`, for a decision that
  undoes or abandons the one in record 4;
* **update**: `VISUAL= EDITOR= pse new -u 4 "Title"`, for a decision that
  builds on record 4, such as the next cycle after an evaluation.

Then edit the new file: keep the lines the tool wrote (the title, the date,
`Draft`, and any link lines in the Status section), and replace the guidance
text under each heading with the approved content. Never write link lines by
hand. See [pse-tool.md](references/pse-tool.md) for every command and its
effects.

### 10. Evaluate

A decision cannot be evaluated before it is carried out. Until then, the
Evaluate section says how and when it will be, and that it has not been yet.

When the user returns to evaluate a record, gather the facts again (what
happened, from the same kinds of source), ask whether it solved the problem,
and write the answer into the Evaluate section of the same record. Leave the
earlier sections as they were: they are what was known and decided at the
time, and hindsight belongs in Evaluate.

If the evaluation shows that the problem remains, or has changed, PIOSEE
starts again from the Problem. That is a new decision, and a new record: `-u`
on the old one if it continues the course, `-r` if it reverses it.

### 11. Finish

* Read each file back, and check it against the list at the end of
  [writing-piosees.md](references/writing-piosees.md#before-you-finish).
* Run `pse list` to confirm the record is in the log.
* If the project keeps a table of contents made by `pse list --toc`, offer to
  regenerate it.
* Tell the user the path, that the record is a draft, any records the tool
  changed, and when the decision is to be evaluated. Leave committing to them.

## Edge cases

* **The user asks you to decide.** Explain that the record will say who chose
  and why, and that the choice must be theirs; offer your comparison of the
  options, and your view, labelled as yours.
* **The user wants to skip steps.** Under real time pressure, keep each step
  to a sentence rather than skipping it; the ones most often skipped are
  Information and Evaluate. If the user still wants to skip one, say what is
  lost, and write "Not recorded" under it rather than inventing it.
* **The decision was made on instinct.** Record that, in the user's words.
  PIOSEE does not capture intuition well; do not invent the analysis that did
  not happen.
* **The `pse` command is missing.** Tell the user. Do not create records by
  hand, since the numbering and links would not be maintained.
* **The tool reports an error.** Show the user the message, which begins
  `pse:`, rather than working around it.
