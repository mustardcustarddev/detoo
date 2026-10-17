---
name: create-piosee
description: Writes a PIOSEE playbook with the pse command-line tool. A playbook is a plan, made in advance, for a situation the team expects to face, such as a customer reporting an unexpected charge, a tax threshold being crossed, or a failure that has happened before. It sets out the Scenario, the Problem to decide, the Information to gather, the Options, the rules for Selection, the steps of Execution and how to Evaluate the result. The skill turns a short request, such as "write a playbook for when a customer can't cancel", into a draft. It reads the project's template afresh, gathers the facts that the code, the documents and the systems hold, and interviews the user for the decision, the options and the rules for choosing, rather than guessing them. Every playbook it creates has the status Draft. Use it when the user asks to write, plan or draft a playbook or a PIOSEE, to plan for an incident or situation before it happens, to plan "what we do when…", to revise, replace or update a playbook (for example after one was used), or to start a PIOSEE log with pse init.
license: Apache-2.0
compatibility: Requires the pse command-line tool on PATH.
---

# Create a PIOSEE playbook

A PIOSEE playbook is a plan for a situation that has not happened yet. It
names the situation (the **Scenario**) and sets out how to handle it:

* **Problem**: the decision that will have to be made;
* **Information**: the investigation that decision needs;
* **Options**: the known options, whether one or several;
* **Selection**: the criteria for choosing between them;
* **Execution**: the steps to carry out the option chosen;
* **Evaluation**: the steps to judge the result.

When the situation arises, whoever is on hand opens the playbook and follows
it. They should not have to work out under pressure what was worked out
calmly in advance.

A playbook does not record an incident. It records no decision made at the
time, and nothing that happened. An incident can prompt a playbook, and what
was learned from it can shape one, but the playbook is written for the next
time.

What the situation is, what the team will decide, and how, belong to the user
and their team. Your job is to draw the plan out of them, to find the facts
that make it concrete, and to write it down faithfully. You are the
interviewer and the scribe, not the owner of the plan.

## The rules that matter most

1. **A playbook is a plan, not a record.** Write it for every occurrence of
   the situation, not for one of them. Write it in the present tense and the
   imperative: what to collect, which option to select when, what to do. If
   the user is describing an incident, write the plan that would serve the
   next one, and leave the incident's own details out.
2. **Write nothing you cannot evidence.** Every claim must come from one of
   three places:
   * the user, in this conversation;
   * the repository or documents, which you have read;
   * the output of a command you ran.

   This matters most for the concrete detail that makes a playbook usable:
   the names of services, logs, paths and fields, the steps through a
   dashboard, links, and reason codes. Policies, owners, authority, deadlines
   and targets matter just as much. Someone will follow these under pressure,
   and an invented one sends them the wrong way.
3. **The user sets the rules.** The options that are open, when each applies,
   the criteria for selecting one and the default are the team's. You may
   suggest options, criteria or steps when the user asks. Label them as
   suggestions, and write them in only if the user adopts them.
4. **Every playbook you create is a Draft.** The tool writes `Draft`; never
   change it. Whether a playbook is agreed is for the team to record.

When you want to make a claim and have no evidence for it, **stop and ask the
user**. Ask through your question tool (`AskUserQuestion` in Claude Code), and
wait for the answer. [asking.md](references/asking.md) says how to ask an open
question with a tool that needs options, and when to ask in the chat instead.
Do not carry on with an assumed answer.

Before you start, read:

* [piosee.md](references/piosee.md): what a PIOSEE playbook is, what each
  section is for, and what makes a good one;
* [interviewing.md](references/interviewing.md): how to plan and run the
  interview, with a worked example;
* [writing-piosees.md](references/writing-piosees.md): what counts as
  evidence, what to do when a fact is missing, and how to write each section.

## Workflow

### 1. Find the log

```bash
pse config        # pse_dir and pse_template
pse list          # existing playbooks
```

If there is no log yet, ask the user whether to start one, and where. There
is none when `pse list` fails with `the docs/piosee directory does not exist`,
or names another missing directory. `pse init [DIR]` creates the directory
and writes playbook 1, which plans for deciding when to write a playbook.

Read the titles of the existing playbooks. Read in full any whose scenario is
close to this one: the new playbook may duplicate it, replace it or update it.
Read one or two others too, for the house style: how concrete they are, and
which extra sections they add, such as communication templates or reference
links.

### 2. Read the template, every time

Read the template in use, afresh, for every playbook. It is the file named by
`$PSE_TEMPLATE` if that is set (`pse config` does not report it). Otherwise it
is the file in `pse_template`, if that is set. Otherwise it is the built-in
one in [pse-tool.md](references/pse-tool.md#templates). Do not work from
memory of the built-in template or of an earlier playbook; projects change
their templates.

List its sections, and the guidance under each. These are what the interview
must fill.

### 3. Check that it is a playbook

* **A situation the team expects to face** and must decide how to handle is
  a playbook. It may be one that recurs, such as a kind of customer
  complaint. It may be one that is scheduled or foreseen, such as a
  regulatory deadline, or one that has happened once and could happen again.
  There may be one known option or several.
* **An incident happening now** is not the time to write a playbook. Tell the
  user to deal with it first. If an existing playbook covers it, point them
  to it. Offer to write or revise the playbook afterwards, from what they
  learned.
* **A record of what happened in an incident**, such as a postmortem, is not a
  playbook. Say so. Offer to write a playbook for the next time, and leave the
  record to wherever the team keeps them.
* **A long-lived choice about how the system is built** is an architecture
  decision: suggest an ADR, with the `create-adr` skill if it is available. A
  change others would want a say in first is a proposal: suggest an RFC, with
  `create-rfc`. Work for someone who benefits is a user story, with
  `create-user-story`. The user chooses.
* **A request that covers several situations** may need several playbooks.
  List the situations, and confirm which to write.

### 4. Gather the facts

A playbook is only as useful as it is concrete. The person following it needs
to know exactly where to look and what to click. Much of that is in the
repository, so learn what it can tell you before asking anything:

* the code that handles the situation, and the names it uses: services and
  components, the logs they write, records they store (by path and field),
  events they emit or receive, and configuration;
* the third-party systems involved, such as a payment provider or a cloud
  console, and the documentation the project links to for them;
* runbooks, policies, support macros, earlier playbooks, ADRs, and anything
  the user points you to;
* how the situation showed itself before, if the user points you to an
  incident, a ticket or a log, as evidence for the Scenario and the
  Information to gather.

Run read-only commands only. Note where each fact came from. Cite a link only
if you have read the page, or the user gave it to you.

The facts are yours to establish. The plan is the user's. Start the fact
sheet described in
[writing-piosees.md](references/writing-piosees.md#facts-first), and mark
everything the template needs that you do not have as **Unknown**.

### 5. Plan the interview

Go through the template's sections, and under each write the questions this
situation raises, using the question bank in
[interviewing.md](references/interviewing.md). Make them specific to the
situation and to what you found. Cross off what the facts already answer.
Keep the plan in your working notes, and revise it as answers come in.

### 6. Interview the user

Ask in rounds, a section or two at a time, in the template's order. Each
round is one call to your question tool. For the built-in template:

1. **Scenario**: what happens, or is reported, and how the team comes to know
   of it. How often it happens, or when it is expected, and the standing
   conditions that frame it.
2. **Problem**: the decision that has to be made, as one question, and what
   goes wrong if it is decided wrongly, in each direction.
3. **Information**: show the sources you found, and ask what else the person
   handling it must find out, and from whom. Ask what is known and what is
   assumed, and how to settle the critical unknowns quickly.
4. **Options**: the courses of action open to them, when each applies, and
   the tradeoffs of each. Ask what to do when the evidence is unclear: the
   conservative fallback.
5. **Selection**: the criteria, which findings lead to which option, the
   default, how confident they are in the rule, and what it assumes.
6. **Execution**: the steps for each option, in order, by role. Cover what
   each step depends on, what to record, whom to tell and how, how to roll
   back, and what to do to stop it happening again.
7. **Evaluation**: how they will tell the playbook worked, with measures and
   targets. Ask who reviews it, and when, and what they do if it falls short.

After each round, play the section back in a sentence or two, so the user can
correct it before the next section builds on it. If an answer is vague, ask a
follow-up rather than sharpening it yourself.

When the user cannot answer a point, offer the choices in
[writing-piosees.md](references/writing-piosees.md#when-a-fact-is-missing).
The scenario, the decision and the options cannot be deferred: without them
there is no plan to write. Say so.

If the user asks you to fill the gaps yourself, explain that someone will
follow this playbook under pressure, so it must say what the team means.
Offer to write only what is known, and mark the rest as open questions. If
the user asks for your ideas, give them as labelled suggestions; see
[writing-piosees.md](references/writing-piosees.md#your-own-suggestions).

Stop when every section of the template can be written from the fact sheet.

### 7. Draft, and get approval

Write the draft in the conversation, in the template's sections, adding any
section the user asked for or the house style uses. After it, list the
sources: each claim, and where it came from. Ask the user to approve it or
correct it. Revise until they approve.

Do not create a file before approval. `pse new` takes the next number, and
with `-r` or `-u` it changes other playbooks too.

### 8. Create the playbook

```bash
VISUAL= EDITOR= pse new "Title of the playbook"
```

Clearing `VISUAL` and `EDITOR` stops the tool waiting on an editor. It prints
the path of the new file.

A new playbook may replace or update an existing one. If so, tell the user
that the tool records the link in **both** playbooks at once, and that with
`-r` the old playbook loses its `Draft` line. Then:

* **replace**: `VISUAL= EDITOR= pse new -r 4 "Title"`, for a new plan for the
  situation playbook 4 planned for;
* **update**: `VISUAL= EDITOR= pse new -u 4 "Title"`, for a playbook that
  changes part of playbook 4, such as a case it left out.

Then edit the new file:

* Keep the lines the tool wrote: the title, the date, `Draft`, and any link
  lines in the Status section.
* Replace the guidance text under each heading with the approved content.

Never write link lines by hand. See [pse-tool.md](references/pse-tool.md) for
every command and its effects.

### 9. Finish

* Read the file back, and check it against the list at the end of
  [writing-piosees.md](references/writing-piosees.md#before-you-finish).
* Run `pse list` to confirm the playbook is in the log.
* The project may keep a table of contents made by `pse list --toc`, such as
  a `README.md` in the log directory. If so, offer to regenerate it.
* Tell the user:
  * the path, and that the playbook is a draft;
  * any playbooks the tool changed;
  * the open questions left for the team;
  * that the people who will follow the playbook should review it.

  Leave committing to them.

## Revising a playbook

A playbook is meant to be kept current, so that it works when it is needed.

* **A Draft** is revised in place, to the same standard of evidence.
* **When a fact changes**, such as a renamed service, a moved menu or a new
  owner, ask the user whether to correct the playbook in place. Git keeps
  the history.
* **When the plan changes**, ask the user which they prefer:
  * a new playbook that replaces the old one (`-r`), for a new set of
    options or selection rules for the same situation;
  * a new playbook that updates it (`-u`), for a case it left out;
  * an edit in place.
* **After the playbook has been used**, the user may want to review it
  against its own Evaluation section. Ask what happened, what the playbook
  got right, and where it was missing, wrong or slow. Revise the plan from
  their answers, as above. Leave the incident itself out of the playbook; it
  belongs where the team records its incidents.

## Edge cases

* **The user says the playbook is agreed.** It is still created as a Draft.
  Tell them they can mark it agreed through their process.
* **The user wants a quick playbook now, and the details later.** Write what
  is known, and mark the rest as open questions. The scenario, the decision
  and at least one option must still come from the user.
* **There is only one known option.** That is a valid playbook. Ask what to
  do when that option does not apply, or the evidence is unclear, which is
  often an escalation. Record the answer as the fallback. Do not invent one.
* **The user asks you to decide the selection rule.** Explain that the
  playbook will be followed as the team's policy, so the rule must be theirs.
  Offer your comparison of the options, and your view, labelled as yours.
* **The `pse` command is missing.** Tell the user. Do not create playbooks by
  hand, since the numbering and links would not be maintained.
* **The tool reports an error.** Show the user the message, which begins
  `pse:`, rather than working around it.
