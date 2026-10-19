---
name: create-firedrill
description: Writes a fire drill from a problem playbook - an exercise in which the team practises responding to the problem the playbook plans for, on a local copy of the system, without causing any real problem. It reads the playbook and its SOPs, researches how the system runs locally, what it logs and what the other systems involved report, such as a payment provider's events and errors, and writes the fire drill document - the setup, the symptoms the team will meet, the response the playbook leads to, and a debrief - with realistic assets, such as the incident description and logs, based on the real system. Use it to write, plan or prepare a fire drill, game day, incident simulation or tabletop exercise, to practise or test a playbook, or to rehearse "what we do when...".
license: Apache-2.0
compatibility: Requires the plm command-line tool on PATH to find the playbooks, and the sop tool to find the SOPs they refer to. Reading third-party documentation needs a web fetch or search tool.
---

# Create a fire drill

A fire drill is an exercise: the team meets a simulated incident, and
responds to it as it would to a real one, so that it is ready when a real
one comes. Here it practises one problem playbook, written with the
`create-playbook` skill and the `plm` tool. The drill shows whether the team
can follow the playbook, and whether the playbook is right: a log it names
that says nothing useful, an option it lacks, a step that does not work. See
[Fire Drills](https://openpracticelibrary.com/practice/fire-drills/) in the
Open Practice Library for the practice.

The drill runs on a local copy of the system, run the way the project runs
it locally. Nothing is broken to make the incident. The fire drill document
describes the symptoms, and the drill's assets stand in for the evidence: an
incident description, such as the alert or the customer's message that
starts the drill, and logs like those the real system would write.

Almost everything the drill needs can be found rather than asked for. The
playbook gives the scenario and the response. The code shows how the system
runs locally and what it logs, in what format. The documentation of the
other systems involved shows what they report: a payment provider's
documentation, not the user, says what a failed payment's event and error
look like. So, as with `create-sop`: **research first, draft, and ask last**.

## The rules that matter most

1. **Never cause a real problem.** Do not break, stop, overload or
   misconfigure anything, locally or anywhere else, to make the incident.
   Describe the symptoms in the document, and put the evidence in the
   assets. Run read-only commands only: do not run the system, call a
   system's API or query a live system, unless the user asks you to.
2. **The drill does no harm when it is run.** Wherever the playbook's
   response would act on something outside the local system, such as
   production, a provider's live account, a release, a customer or a public
   issue, the drill says that the team describes the step rather than takes
   it. Every message the drill has the team send, and the incident
   description, is marked as a drill.
3. **Base every asset on the real system.** Every log line has the format,
   the fields and the message the code really writes. Every event, error
   code and payload from another system is the one its documentation gives,
   for the version the project uses. What you remember of a log format or
   an API, or think is typical, is not a source.
4. **Invent only the incident's particulars.** The times, identifiers,
   amounts and people of this one incident are made up, in the real formats,
   and consistent across every asset. Never copy real customer data, real
   secrets or real production logs into an asset.
5. **Follow the playbook; do not rewrite it.** The expected response is the
   playbook's, as it stands. Where research shows that the playbook is
   wrong or missing something, note it for the debrief, and leave revising
   the playbook to the `create-playbook` skill.

Ask the user only what research cannot settle, through your question tool
(`AskUserQuestion` in Claude Code), and wait for the answer.
[asking.md](references/asking.md) says how to ask an open question with a
tool that needs options, and when to ask in the chat instead.

Before you start, read
[writing-firedrills.md](references/writing-firedrills.md): the sections of
the fire drill document, how to make the assets realistic, and the checks
before you finish.

## Workflow

### 1. Read and understand the playbook

```bash
plm config        # plm_dir
plm list          # the playbooks
sop config        # sop_dir
```

Find the playbook the user names. If they name none, or a problem with no
playbook, list the playbooks and ask which to drill. A problem with no
playbook needs one first: suggest the `create-playbook` skill, if it is
available.

Read the playbook in full, and every SOP it links to. Note, in your working
notes:

* the **Scenario**: how the problem shows itself, and how the team comes to
  know of it;
* the **Immediate action**, or that none is indicated, and the
  **Counter-indicated actions**;
* each source of evidence the **Information** section names: the logs, the
  records, the dashboards, the people, and what is looked for in each;
* the **Options**, and the **Selection** rules that map the evidence onto
  them;
* each step of **Execution** and **Evaluation**, and whether it acts on the
  local system or on something outside it.

A playbook whose status is not `Draft` or `Accepted`, such as one that has
been replaced, is the wrong one to drill: say so, and point to the playbook
that replaced it.

### 2. Choose the case

A playbook plans for every occurrence of a problem; a drill is one of them.
Choose the case the drill presents: the particular facts that the evidence
will show, and so the option the Selection rules lead to. Prefer the case
that exercises most of the playbook: the immediate action, most of the
Information, and an option with real steps. Where the Selection has a
conservative fallback for unclear evidence, a case that leads to it is a good
second drill.

Propose the case in step 6, with the option it leads to. Do not ask for it
first.

### 3. Research the local setup

Find how the project runs the system locally: the README, the contributing
guide, the Makefile, scripts, `docker-compose` or devcontainer files, a
`Procfile`, and the configuration for the local or test environment. Note
the exact commands, the services they start, the ports, and any seed or
fixture data. For another system involved, find whether the project uses its
test mode, sandbox or emulator locally, and how.

The setup the drill describes is only the system run as it normally is
locally. It adds no fault. If the system cannot be run locally, say so in
the setup, and run the drill on the assets alone, as a tabletop exercise.

### 4. Research the evidence

For each source of evidence the playbook's Information names, and each
symptom the Scenario describes, find what the real system would show.

1. **What the code writes.** Find the logger and its configuration: the
   format (JSON or text), the field names, the timestamp format and time
   zone, the levels, and where the output goes. Then follow the code path of
   the problem, and note each log line, error and output it writes, verbatim,
   with its file and line. For a command-line tool, the output and its
   error messages are its logs.
2. **What the other systems report.** Read the official documentation of
   each system involved, for the version the project uses, such as an API
   version pinned in its configuration: the events it sends, their payloads,
   its error codes and messages, the formats of its identifiers, and what
   its dashboard shows for this case.
3. **What the records hold.** Find the tables, files or fields the playbook
   says to check, and what they would hold in this case.
4. **What there is not.** Note any evidence the playbook looks for that the
   system does not produce: a log line that is never written, an alert that
   does not exist. Do not fake it in the assets. It is a finding for the
   debrief.

Note each finding, with its source, in a fact sheet, as
[writing-firedrills.md](references/writing-firedrills.md#facts-first)
describes.

### 5. Draft the drill

Draft, from the fact sheet alone, as
[writing-firedrills.md](references/writing-firedrills.md) describes:

* the **fire drill document**, in its sections;
* the **assets**: the incident description, the logs, and any other
  evidence the Information section needs, such as an event payload, a record
  or a dashboard's view, written as a file.

Write a gap research could not fill as an open question in the document,
not as a guess in an asset.

### 6. Ask once

Write in the chat: the case, and the option it leads to; the setup; the
symptoms; the list of assets, with the first lines of each log; the expected
response; and the open questions. Then ask, in one call of your question
tool:

* **whether to write the drill**: **Write it** as it is, or **Change it**,
  with the change typed in as Other;
* and, in the same call, up to three questions research could not settle,
  such as which of two sources holds, or a fact about the team's tools that
  no document records.

Revise for any change, and play it back if it changes more than a detail.

### 7. Write the drill

Write the drill in a `firedrill` directory beside the problem log: if
`plm_dir` is `docs/problem`, in `docs/firedrill`. Give each drill its own
directory, named for the playbook's number and the case, such as
`docs/firedrill/FIREDRILL-0004-duplicate-charge-refunded/`:

```text
FIREDRILL-0004-duplicate-charge-refunded/
  FIREDRILL.md          the fire drill document
  assets/
    incident.md         the incident description
    logs/...            the logs, one file for each source
    ...                 any other evidence
```

If the directory exists, show what is in it, and ask whether to replace it
or choose another name. Never overwrite a drill unasked.

### 8. Finish

* Read every file back, and check the drill against the list at the end of
  [writing-firedrills.md](references/writing-firedrills.md#before-you-finish).
* Tell the user:
  * the path of the document, and that it is a draft;
  * that the document holds the expected response, so the team should see
    only `assets/` until the debrief;
  * the open questions, and anything research found wrong or missing in
    the playbook;
  * that after the drill, `create-playbook` can revise the playbook from
    the debrief.

  Leave committing to them.

## Edge cases

* **The user wants a drill for a problem with no playbook.** Suggest writing
  the playbook first. A drill tests a plan; without one there is nothing to
  follow, and nothing to evaluate.
* **The user asks to make the problem happen for real**, such as stopping a
  service or sending bad data. Decline: describe the symptoms in the document
  and the assets. If they want the system itself to show the symptom, leave
  that to them, as a note in the setup, and do not do it.
* **The playbook names evidence the system does not produce.** Say so in the
  document's debrief notes, and in your summary. Do not invent the evidence.
* **Research finds no local setup.** Say where you looked. Run the drill as a
  tabletop exercise on the assets, and say so in the setup.
* **You cannot read the documentation of a system involved**, because you
  have no web tool, or the pages are not public. Say so, ask the user for
  the documentation or a link, or leave its evidence as an open question.
  Do not write its events or errors from memory.
* **The sources disagree**, such as the code and the playbook on what is
  logged. The code is the better evidence of what the system writes; build
  the assets from it, and note the playbook's error for the debrief.
* **The `plm` command is missing.** Ask the user for the path of the
  playbook, and read it from there.
