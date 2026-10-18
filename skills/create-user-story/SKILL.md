---
name: create-user-story
description: Writes user stories with the ust command-line tool. Turns a short subject, such as "write a user story for exporting invoices", into a draft story with acceptance criteria, by reading the project's story template afresh, gathering facts from the code, and interviewing the user for who the story is for, what they need and why, rather than guessing. Every story it creates has the status Draft. Use when the user asks to write, draft, create or add a user story, a backlog item or acceptance criteria, to split an epic or a large story into smaller ones, to replace or update an existing story, or to start a story log with ust init.
license: Apache-2.0
compatibility: Requires the ust command-line tool on PATH.
---

# Create a user story

A user story is a short description of something a person wants to do with
the software, from their point of view, and why: "As a [type of user], I
[need/want/am required] to [do something], so that [reason or benefit]." With
its acceptance criteria, it tells the team what to build and how they will
know it is done. Only the user, or the people they speak for, know who wants
it and why. Your job is to draw that out, and write it down faithfully: you
are the interviewer and scribe, not the product owner.

## The rules that matter most

1. **Write nothing you cannot evidence.** Every claim in the story must come
   from the user, in this conversation; from the repository, which you have
   read; or from the output of a command you ran. Users, reasons, acceptance
   criteria, scope, sizes and priorities are not to be invented, inferred as if
   known, or embellished.
2. **A subject is not a story.** "A user story for exporting invoices" does
   not say who exports them, what they do with them, or what must be true when
   it is done. Do not draft from it. Gather the facts, then interview the
   user.
3. **Every story you create is a Draft.** The tool writes `Draft`; never
   change it. Whether a story is ready, in progress, done or dropped is for
   the team to record.

When you want to make a claim and have no evidence for it, **stop and ask the
user**. Ask through your question tool (`AskUserQuestion` in Claude Code), and
wait for the answer. [asking.md](references/asking.md) says how to ask an open
question with a tool that needs options, and when to ask in the chat instead.
Do not carry on with an assumed answer.

Before you start, read:

* [user-stories.md](references/user-stories.md): what a user story is, and
  what makes a good one, summarised from Mike Cohn's guide;
* [interviewing.md](references/interviewing.md): how to plan and run the
  interview, with a worked example;
* [writing-user-stories.md](references/writing-user-stories.md): what counts
  as evidence, what to do when a fact is missing, and how to write each
  section.

## Workflow

### 1. Find the log

```bash
ust config        # ust_dir and ust_template
ust list          # existing stories
```

If there is no log yet (`ust list` fails with `the docs/userstory directory
does not exist`, or names another missing directory), ask the user whether to
start one, and where. `ust init [DIR]` creates the directory and writes story
1, which describes the practice of keeping stories.

Read story 1, and any story or document that sets out the project's
conventions. Read the titles of the existing stories, and in full any whose
subject is close to this one: the new story may replace or update one of
them, or duplicate it.

### 2. Read the template, every time

Read the template in use, afresh, for every story: the file named by
`$UST_TEMPLATE` if it is set (`ust config` does not report it), otherwise the
file in `ust_template` if it is set, otherwise the built-in one in
[ust-tool.md](references/ust-tool.md#templates). Do not work from memory of
the built-in template or of an earlier story; projects change their templates.

List its sections, and the guidance under each. These are what the interview
must fill.

### 3. Check that it is a story

* A story describes an outcome for someone who benefits. A task ("add a
  column to the invoices table") is work towards one: ask what the outcome is,
  and who it is for.
* Technical work with no one who benefits may be better as a technical
  backlog item than a story. Say so, and let the user choose.
* A change that others would want a say in before it is made, such as a new
  architecture, is a proposal: suggest an RFC, with the `create-rfc` skill if
  it is available. A decision already made is an ADR, with `create-adr`.
* A large or vague subject is an epic. Plan to split it; see
  [writing-user-stories.md](references/writing-user-stories.md#size).
* If the request holds several stories, list them, and confirm which to write.

### 4. Gather the facts

Before asking anything, learn what the repository can tell you, so that you
ask only what you could not read, and can ask precisely:

* what the software does now in the area the story touches: its screens,
  commands or endpoints, its data, and its rules;
* the roles or kinds of user the software knows about;
* related stories, RFCs, ADRs and documents, and anything the user points you
  to;
* what the user has already said.

Read in proportion to the story; one small story needs the code it touches,
not a survey of the system.

The current state is yours to establish. Who wants the change, what they
need and why are the user's. Start the fact sheet described in
[writing-user-stories.md](references/writing-user-stories.md#facts-first), and
mark everything the template needs that you do not have as **Unknown**.

### 5. Plan the interview

Go through the template's sections, and under each write the questions this
story raises, using the question bank and subject lenses in
[interviewing.md](references/interviewing.md). Make them specific to the
subject and to what you found. Cross off what the facts already answer. Keep
the plan in your working notes, and revise it as answers come in.

### 6. Interview the user

Ask in rounds, a few related questions at a time, in an order where each
builds on the last. Each round is one call to your question tool. For the
built-in template:

1. Who and what: the type of user, and what they need to do.
2. Why: what they do with it, the decision it helps them make, why now.
3. Confirmation: examples, rules, unhappy paths and limits, as acceptance
   criteria.
4. Size: whether it can be finished in a sprint, and if not, how to split it.
5. Notes and links: scope, dependencies, open questions, and existing stories
   it replaces or updates.

Ask open questions about users and reasons. Offer as options only answers you
have evidence for, such as the roles you found in the code; an open question
still goes through the question tool, with the answer typed in, as
[asking.md](references/asking.md) describes. Never offer guessed reasons as
choices: a user picking the nearest option is not evidence. You may put an
observation from the code to the user and ask whether it matters; their answer
is the evidence.

After each round, play the story back in its one sentence, with the criteria
so far, so the user can correct it. If an answer is vague, ask a follow-up
rather than sharpening it yourself.

When the user cannot answer a point, offer the choices in
[writing-user-stories.md](references/writing-user-stories.md#when-a-fact-is-missing):
make it an open question for refinement, or leave it out. Who the story is
for, what they need, and why, cannot be deferred: without them, the story is
not ready to write. Say so.

If the user asks you to fill the gaps yourself, explain that the team will
build what the story says, so it must say what the user means; offer to write
only what is known, with the rest as open questions. If the user asks for
your ideas, such as acceptance criteria or a split, give them as labelled
suggestions; see
[writing-user-stories.md](references/writing-user-stories.md#your-own-suggestions).

Stop when every section of the template can be written from the fact sheet.

### 7. Draft, and get approval

Write the draft in the conversation, in the template's sections. After it,
list the sources: each claim, and where it came from. Ask the user to approve
it or correct it. Revise until they approve. When there are several stories,
draft them all, and get approval for the set.

Do not create a file before approval: `ust new` takes the next number and,
with `-r` or `-u`, changes other stories too.

### 8. Create the story

```bash
VISUAL= EDITOR= ust new "Title of the story"
```

Clearing `VISUAL` and `EDITOR` stops the tool waiting on an editor. It prints
the path of the new file.

If the story replaces or updates an existing one, tell the user that the tool
records the link in **both** stories at once, and that with `-r` the old story
loses its `Draft` line. Then:

* **replace**: `VISUAL= EDITOR= ust new -r 4 "Title"`, for each story that
  replaces story 4, such as each part of a split;
* **update**: `VISUAL= EDITOR= ust new -u 4 "Title"`, for a story that
  changes what story 4 asked for.

Then edit each new file in place:

* Keep the lines the tool wrote: the title, the date, `Draft`, and any link
  lines in the Status section. Do not rewrite the whole file, and never
  retype the title: its number is the tool's, and must stay the number in
  the file name.
* Replace the guidance text under each heading with the approved content,
  one section at a time.

Never write link lines by hand. See [ust-tool.md](references/ust-tool.md) for
every command and its effects.

### 9. Finish

* Read each file back, and check it against the list at the end of
  [writing-user-stories.md](references/writing-user-stories.md#before-you-finish).
* Check that the number in each title, `# N.`, is the number in its file
  name, `STORY-000N-...`. If they differ, put the tool's title line back.
* Run `ust list` to confirm the stories are in the log.
* If the project keeps a generated table of contents (look for one made by
  `ust list --toc`, such as a `README.md` in the log directory), offer to
  regenerate it.
* Tell the user the paths, that the stories are drafts, any stories the tool
  changed, and the open questions left for refinement. Leave committing to
  them.

## Revising a story

* **A Draft**, or any story not yet built, is meant to change as the
  conversation goes on. Revise it in place, to the same standard of evidence:
  changes come from the user, or from notes of conversations they show you.
* **A story that is done** is history: do not rewrite what it asked for.
  Fixing a typo or a broken link is fine. To change the behaviour it
  delivered, write a new story that updates it, following this workflow.

## Edge cases

* **The user says the story is ready for the sprint.** It is still created as
  a Draft. Tell them they can mark it ready through their process.
* **The user wants a quick story now and the details later.** Write what is
  known, with the rest as open questions, provided who, what and why came from
  the user.
* **The user wants many stories at once**, such as from an epic or a
  document. Interview once for what they share, then for each story what is
  its own. Keep a fact sheet per story.
* **The `ust` command is missing.** Tell the user. Do not create stories by
  hand, since the numbering and links would not be maintained.
* **The tool reports an error.** Show the user the message, which begins
  `ust:`, rather than working around it.
