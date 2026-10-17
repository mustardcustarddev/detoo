---
name: create-piosee
description: Writes a PIOSEE playbook with the pse command-line tool. A playbook is a plan, made in advance, for a situation the team expects to face, such as a customer reporting an unexpected charge, a tax threshold being crossed, or a failure that has happened before. It sets out the Scenario, the Problem to decide, the Information to gather, the Options, the rules for Selection, the steps of Execution and how to Evaluate the result. The skill turns a short request, such as "write a playbook for when a Stripe payment is handled incorrectly", into a draft. It interviews the user about the nature of the problem, then researches the codebase, the project's documents and the documentation of the systems involved, such as the payment provider's, to propose the investigation, the options and the steps. It asks the user to confirm each proposal, and interviews them step by step only where its research finds nothing. Every playbook it creates has the status Draft. Use it when the user asks to write, plan or draft a playbook or a PIOSEE, to plan for an incident or situation before it happens, to plan "what we do when…", to revise, replace or update a playbook (for example after one was used), or to start a PIOSEE log with pse init.
license: Apache-2.0
compatibility: Requires the pse command-line tool on PATH. Reading third-party documentation needs a web fetch or search tool.
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

The user knows what the problem is. Most of the rest can be found out: the
code shows where the evidence is and what the system can do, and the
documentation of the systems involved shows the options and the steps. Your
job is to learn the problem from the user, research the rest, propose it with
its sources, and have the user confirm it. Ask the user to work through a
section step by step only when your research finds nothing for it.

## The rules that matter most

1. **A playbook is a plan, not a record.** Write it for every occurrence of
   the situation, not for one of them. Write it in the present tense and the
   imperative: what to collect, which option to select when, what to do. If
   the user is describing an incident, write the plan that would serve the
   next one, and leave the incident's own details out.
2. **Interview for the problem; research the rest.** Ask the user what
   happens, how it shows itself, and what has to be decided. Then research
   the context, the investigation, the options and the steps, and propose
   them. Ask the user to confirm each proposal: yes, no, or a correction in
   their own words.
3. **Write nothing you cannot evidence.** Every claim must come from one of
   four places:
   * the user, in this conversation;
   * the repository or the project's documents, which you have read;
   * the documentation of a system involved, such as a payment provider's,
     which you have read;
   * the output of a command you ran.

   A proposal cites its source, and enters the playbook only when the user
   confirms it. What you know from memory, or think is common practice, is
   not a source: read the documentation, or ask. This matters most for the
   concrete detail that makes a playbook usable: the names of services,
   logs, paths and fields, the steps through a dashboard or an API, links,
   and reason codes. Someone will follow these under pressure, and an
   invented one sends them the wrong way.
4. **Policy is the team's.** Authority limits, owners, deadlines, targets and
   the default option are the team's to set. Propose one only when a policy,
   a contract or the system's documentation states it, and say which. Where
   none does, ask.
5. **Every playbook you create is a Draft.** The tool writes `Draft`; never
   change it. Whether a playbook is agreed is for the team to record.

When you want to make a claim and have no evidence for it, **stop and ask the
user**. Ask through your question tool (`AskUserQuestion` in Claude Code), and
wait for the answer. [asking.md](references/asking.md) says how to ask an open
question with a tool that needs options, and when to ask in the chat instead.
Do not carry on with an assumed answer.

Before you start, read:

* [piosee.md](references/piosee.md): what a PIOSEE playbook is, what each
  section is for, and what makes a good one;
* [interviewing.md](references/interviewing.md): how to interview for the
  problem, research the rest, and confirm it, with a worked example;
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

List its sections, and the guidance under each. These are what the research
and the confirmations must fill.

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

### 4. Interview for the problem

A request such as "a Stripe payment handled incorrectly" names an area, not a
situation. Before researching, learn from the user, in one or two rounds:

* **what happens**, or is reported: a customer charged twice, charged with no
  order, or paid for an order that stays unpaid;
* **how the team comes to know of it**: a customer contact, an alert, a
  reconciliation report;
* **the systems involved**, as the user knows them: the provider, the
  services, the tools support uses;
* **what has to be decided** when it happens, and what goes wrong if it is
  decided wrongly;
* **anything to read first**: an incident, a ticket, a log, a document.

Play the Scenario and the Problem back in a sentence or two, and have the
user confirm them. These two cannot come from research: without them there is
nothing to research.

### 5. Research

Learn everything you can before asking anything more. Work through the
playbook's concerns in order, and note each finding with its source in the
fact sheet described in
[writing-piosees.md](references/writing-piosees.md#facts-first).

1. **The context of the problem.** Find the code that handles the situation,
   and how it can go wrong: where payments are created, which events are
   handled and how, retries, idempotency, and the states a record passes
   through. Read anything the user pointed you to, and the existing
   playbooks, runbooks, ADRs and policies that touch it.
2. **The course of investigation (Information).** Find where the evidence
   lives, and how to reach it: the logs and the events to look for in them,
   the records by table or path and field, the identifiers that link one
   system to another (such as an order ID kept in a payment's metadata),
   the admin tools and dashboards the project uses, and who can reach each.
3. **The options and the execution (Options, Execution).** Read the
   documentation of each system involved: the provider's guides and API
   reference for this situation, such as refunds, disputes and resending
   events. Use the official documentation, and match it to the version the
   project uses, such as an API version pinned in its configuration. Find
   what the code can do too: admin endpoints, scripts, jobs that can be run
   again. From these, work out the realistic options, when each applies, its
   constraints (time limits, whether it can be undone), and its steps.
4. **Selection and Evaluation.** Find what decides between the options: the
   conditions the documentation sets, and any policy, limit, service level
   or metric in the repository. Draft rules that map what the investigation
   finds onto an option, with a conservative fallback where the evidence is
   unclear, and mark the source of each.

Run read-only commands only. Do not call a system's API, or query a live
system, unless the user asks you to. Cite a link only if you have read the
page, or the user gave it to you.

Mark everything the template needs that your research did not settle as
**Unknown**.

### 6. Confirm what you found

Go through the sections in the template's order. For each, write the
proposal in the chat, in the playbook's own words, with the source of each
point. Then ask the user to confirm it through your question tool, with the
options **Yes** and **No**, and their correction typed in as **Other**.
Confirm up to four sections in one call. See
[interviewing.md](references/interviewing.md#confirming-a-proposal).

* **Yes**: the proposal is the user's now. Mark it confirmed in the fact
  sheet.
* **Other**: use the correction in place of the proposal, and play it back
  if it changes more than a detail.
* **No**: ask what is wrong, as an open question, and treat the section as
  one research did not settle.

Where a section offers real choices, such as which of the options the team
would use, ask about them in one multiple-choice question instead.

### 7. Interview for what research did not find

For each section, or each point, still marked Unknown, interview the user
step by step, using the question bank in
[interviewing.md](references/interviewing.md#question-bank). Make each
question specific, and cite what you did find. This always includes any
policy no source states: authority limits, owners, the default option,
deadlines and targets.

If an answer is vague, ask a follow-up rather than sharpening it yourself.
When the user cannot answer a point, offer the choices in
[writing-piosees.md](references/writing-piosees.md#when-a-fact-is-missing).

Stop when every section of the template can be written from confirmed
entries in the fact sheet.

### 8. Draft, and get approval

Write the draft in the conversation, in the template's sections, adding any
section the user asked for or the house style uses, such as **Reference
Links** to the documentation you read. After it, list the sources: each
claim, where it came from, and whether the user confirmed it or gave it. Ask
the user to approve the draft or correct it. Revise until they approve.

Do not create a file before approval. `pse new` takes the next number, and
with `-r` or `-u` it changes other playbooks too.

### 9. Create the playbook

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

### 10. Finish

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
  owner, research it again, propose the correction, and ask the user whether
  to correct the playbook in place. Git keeps the history.
* **When the plan changes**, ask the user which they prefer:
  * a new playbook that replaces the old one (`-r`), for a new set of
    options or selection rules for the same situation;
  * a new playbook that updates it (`-u`), for a case it left out;
  * an edit in place.
* **After the playbook has been used**, the user may want to review it
  against its own Evaluation section. Ask what happened, what the playbook
  got right, and where it was missing, wrong or slow. Research what it
  lacked, and revise the plan as above. Leave the incident itself out of the
  playbook; it belongs where the team records its incidents.

## Edge cases

* **The user says the playbook is agreed.** It is still created as a Draft.
  Tell them they can mark it agreed through their process.
* **The user wants a quick playbook now, and the details later.** Write what
  is confirmed, and mark the rest as open questions. The scenario and the
  decision must still come from the user, and at least one option must be
  confirmed.
* **There is only one known option.** That is a valid playbook. Propose the
  fallback for when it does not apply, or the evidence is unclear, if a
  source gives one; otherwise ask for it. Do not invent one.
* **You cannot read the documentation**, because you have no web tool, or
  the pages are not public. Say so, ask the user for the documentation or a
  link, and interview for what it would have answered.
* **The sources disagree**, such as the code and the documentation, or two
  versions of an API. Show both, with their sources, and ask the user which
  holds.
* **The user asks you to decide the selection rule.** Explain that the
  playbook will be followed as the team's policy, so the rule must be theirs
  to confirm. Propose one from your sources, labelled with them.
* **The `pse` command is missing.** Tell the user. Do not create playbooks by
  hand, since the numbering and links would not be maintained.
* **The tool reports an error.** Show the user the message, which begins
  `pse:`, rather than working around it.
