---
name: create-firedrill
description: Writes a fire drill with the fyr command-line tool, from a problem playbook - a self-led exercise in which one person practises responding to the problem the playbook plans for, on a local copy of the system, without causing any real problem. The drill reads like a Choose Your Own Adventure book - it sets out the setup and what happens, step by step, each choice leads to a page of what happens next, and the solution holds the expected response and a debrief. It reads the playbook and its SOPs, researches how the system runs locally, what it logs and what the other systems involved report, such as a payment provider's events and errors, and writes realistic pages, such as the incident, the logs, an alert or a metric, based on the real system. Use it to write, plan or prepare a fire drill, game day, incident simulation or tabletop exercise, to practise or test a playbook, or to rehearse "what we do when...".
license: Apache-2.0
compatibility: Requires the fyr command-line tool on PATH to write the drill, the plm tool to find the playbooks, and the sop tool to find the SOPs they refer to. Reading third-party documentation needs a web fetch or search tool.
---

# Create a fire drill

A fire drill is an exercise: you meet a simulated incident, and respond to
it as you would to a real one, so that you are ready when a real one comes.
Here it practises one problem playbook, written with the `create-playbook`
skill and the `plm` tool, and is kept in a log with the `fyr` tool. The drill
shows whether its reader can follow the playbook, and whether the playbook
is right: a log it names that says nothing useful, an option it lacks, a
step that does not work. See
[Fire Drills](https://openpracticelibrary.com/practice/fire-drills/) in the
Open Practice Library for the practice.

Each drill is **self-led**: one person runs it alone, with no one to play
the reporter or hand out the evidence. So it is written as a Choose Your Own
Adventure book for software engineers. The fire drill document sets out the
setup, then what happens, step by step, like a script: first the alert
arrives, then the customer writes again. Where the playbook asks for a
decision, the document asks the reader to choose, and each option names a
page in the drill's `problem` directory, which says what happens next. When
the drill is over, the `solution` directory holds the expected response and
the debrief:

```text
docs/firedrill/
  FIREDRILL-0003-customer-charged-twice.md   the setup and the script
  FIREDRILL-0003-customer-charged-twice/
    problem/
      README.md               written by fyr
      incident.md             the incident that starts the drill
      logs/api.log            the evidence, one file for each source
      refund-at-once.md       what happens after each choice
      ...
    solution/
      README.md               written by fyr
      expected-response.md    the response the playbook leads to
      debrief.md              questions to look back with, and the findings
```

`fyr new` creates the document and the two READMEs. This skill writes
everything else.

The drill runs on a local copy of the system, run the way the project runs
it locally. Nothing is broken to make the incident. The pages describe the
symptoms, and stand in for the evidence: the incident, such as the alert or
the customer's message that starts the drill, and logs like those the real
system would write.

Almost everything the drill needs can be found rather than asked for. The
playbook gives the scenario, the decisions and the response. The code shows
how the system runs locally and what it logs, in what format. The
documentation of the other systems involved shows what they report: a
payment provider's documentation, not the user, says what a failed payment's
event and error look like. So, as with `create-sop`: **research first,
draft, and ask last**.

## The rules that matter most

1. **Never cause a real problem.** Do not break, stop, overload or
   misconfigure anything, locally or anywhere else, to make the incident.
   Describe the symptoms in the pages, and put the evidence in them. Run
   read-only commands only: do not run the system, call a system's API or
   query a live system, unless the user asks you to.
2. **The drill does no harm when it is run.** Wherever a choice would act on
   something outside the local system, such as production, a provider's live
   account, a release, a customer or a public issue, the drill has the
   reader write the step down rather than take it, and the page says what
   would happen. The reader sends nothing. Every page that looks like a real
   message, such as the incident, begins `FIRE DRILL:`.
3. **Base every page on the real system.** Every log line has the format,
   the fields and the message the code really writes. Every event, error
   code and payload from another system is the one its documentation gives,
   for the version the project uses. What you remember of a log format or
   an API, or think is typical, is not a source.
4. **Invent only the incident's particulars.** The times, identifiers,
   amounts and people of this one incident are made up, in the real formats,
   and consistent across every page. Never copy real customer data, real
   secrets or real production logs into a page.
5. **Follow the playbook; do not rewrite it.** The expected response is the
   playbook's, as it stands. Where research shows that the playbook is
   wrong or missing something, note it for the debrief, and leave revising
   the playbook to the `create-playbook` skill.
6. **Name no exact version.** Write "the current release", or "the version
   the reporter has", never `2.4.1`; and give the command that finds a
   version, such as the latest release tag, rather than the version it
   finds. A drill that names a version is out of date with the next release.
   Research still reads the code and documentation of the version the
   project uses.
7. **The drill has the playbook's title, exactly.** Its title is the
   playbook's, as the playbook's title line gives it after the number,
   character for character: a drill of `# 4. Customer charged twice` is
   created with `fyr new "Customer charged twice"`, and its title line is
   `# 3. Customer charged twice`, with the drill's own number. Add nothing to
   it, such as the case, and change nothing, not even the case of a letter
   or a comma. The case goes in the drill, under How to play. Two drills of
   one playbook share its title, and differ in number.

Ask the user only what research cannot settle, through your question tool
(`AskUserQuestion` in Claude Code), and wait for the answer.
[asking.md](references/asking.md) says how to ask an open question with a
tool that needs options, and when to ask in the chat instead.

Before you start, read
[writing-firedrills.md](references/writing-firedrills.md): the fire drill
document, the choices, the pages, the solution, and the checks before you
finish. [kinds-of-drill.md](references/kinds-of-drill.md) says how a drill
that is reproduced, one driven by logs, and one driven by an alert each
differ, and how to make their evidence.
[examples.md](references/examples.md) has good and bad examples of each
part, and the mistakes drills most often make.
[fyr-tool.md](references/fyr-tool.md) has every `fyr` command and its
effects.

## Workflow

### 1. Read and understand the playbook

```bash
plm config        # plm_dir
plm list          # the playbooks
sop config        # sop_dir
fyr config        # fyr_dir and fyr_template
fyr list          # the drills already written
```

Find the playbook the user names. If they name none, or a problem with no
playbook, list the playbooks and ask which to drill. A problem with no
playbook needs one first: suggest the `create-playbook` skill, if it is
available.

Read the playbook in full, and every SOP it links to. Note, in your working
notes:

* the **Scenario**: how the problem shows itself, and how its reader comes
  to know of it;
* the **Primary risks**, the **Immediate action**, or that none is
  indicated, and the **Counter-indicated actions**;
* each source of evidence the **Information** section names: the logs, the
  records, the dashboards, the people, and what is looked for in each;
* the **Options**, and the **Selection** rules that map the evidence onto
  them;
* each step of **Execution** and **Evaluation**, and whether it acts on the
  local system or on something outside it.

A playbook whose status is not `Draft` or `Accepted`, such as one that has
been replaced, is the wrong one to drill: say so, and point to the playbook
that replaced it.

Read any drill already written for the playbook, with its directory: the new
one may repeat it, replace it or update it. If there is no log of drills
yet, `fyr list` fails with `the docs/firedrill directory does not exist`;
ask, in step 7, whether to start one with `fyr init`, which creates
`docs/firedrill` and writes its `README.md`, which sets out how to practise a
playbook in a fire drill. The README is not a drill: the new drill is
number 1.

Read the template in use, afresh, for every drill. It is the file named by
`$FYR_TEMPLATE` if that is set (`fyr config` does not report it). Otherwise
it is the file in `fyr_template`, if that is set. Otherwise it is the
built-in one in [fyr-tool.md](references/fyr-tool.md#templates). Do not work
from memory of it; projects change their templates.

### 2. Choose the kind, and the case

Choose the kind of drill, from where the playbook's evidence is, as
[kinds-of-drill.md](references/kinds-of-drill.md) describes:

* **reproduced**, when a local copy can show the problem without breaking
  anything, such as a command-line tool run on the reporter's files;
* **log-driven**, when the problem happened in a running service, and the
  reader works from logs the drill holds;
* **alert-driven**, when an alert or a metric starts it, and the evidence
  arrives as time passes.

A playbook plans for every occurrence of a problem; a drill is one of them.
Choose the case the drill presents: the particular facts that the evidence
will show, and so the option the Selection rules lead to. Prefer the case
that exercises most of the playbook: the immediate action, most of the
Information, and an option with real steps. Where the Selection has a
conservative fallback for unclear evidence, a case that leads to it is a good
second drill.

Propose the kind and the case in step 7, with the option the case leads to.
Do not ask for them first.

### 3. Research the local setup

Find how the project runs the system locally: the README, the contributing
guide, the Makefile, scripts, `docker-compose` or devcontainer files, a
`Procfile`, and the configuration for the local or test environment. Note
the exact commands, the services they start, the ports, and any seed or
fixture data. For another system involved, find whether the project uses its
test mode, sandbox or emulator locally, and how.

The setup the drill describes is only the system run as it normally is
locally. It adds no fault. If the system cannot be run locally, say so in
the setup, and run the drill on the pages alone, as a tabletop exercise.

### 4. Research the evidence

For each source of evidence the playbook's Information names, and each
symptom the Scenario describes, find what the real system would show.

1. **What the code writes.** Find the logger and its configuration: the
   format (JSON or text), the field names, the timestamp format and time
   zone, the levels, and where the output goes. Then follow the code path of
   the problem, and note each log line, error and output it writes, verbatim,
   with its file and line. For a command-line tool, the output and its
   error messages are its logs. A project's tests often hold the exact
   output, in golden files or fixtures.
2. **What the other systems report.** Read the official documentation of
   each system involved, for the version the project uses, such as an API
   version pinned in its configuration: the events it sends, their payloads,
   its error codes and messages, the formats of its identifiers, and what
   its dashboard shows for this case.
3. **What the records hold.** Find the tables, files or fields the playbook
   says to check, and what they would hold in this case.
4. **What there is not.** Note any evidence the playbook looks for that the
   system does not produce: a log line that is never written, an alert that
   does not exist. Do not fake it in a page. It is a finding for the
   debrief.

Note each finding, with its source, in a fact sheet, as
[writing-firedrills.md](references/writing-firedrills.md#facts-first)
describes. The fact sheet is for you and the user; the drill has no sources
section.

### 5. Plot the drill

Plot the path through the drill before writing any of it, as
[writing-firedrills.md](references/writing-firedrills.md#the-plot)
describes:

* **the steps**: what happens, in order, from the incident to the end;
* **the choices**: at each decision the playbook makes, such as whether to
  take the immediate action, whether the evidence is enough, and which
  option to select, the option the playbook leads to and two or three
  others a person under pressure would plausibly take, drawn from the
  primary risks, the counter-indicated actions, the other options, and the
  shortcuts;
* **the pages**: for each choice, the page it leads to, what happens there,
  and where it sends the reader next: on to a step, back to the choice, or
  to the end of the drill. Every path ends at the solution.

### 6. Draft the drill

Draft, from the fact sheet and the plot alone, as
[writing-firedrills.md](references/writing-firedrills.md) describes:

* the **fire drill document**, in the sections of the template in use: the
  setup, and the drill as a script, with its choices;
* the **problem pages**: the incident, the evidence the reader looks for,
  such as logs, an event payload, a record or a dashboard's view, and a
  page for each choice;
* the **solution**: `expected-response.md` and `debrief.md`.

Write a gap research could not fill as an open question in the debrief, not
as a guess in a page.

For a drill driven by logs or an alert, generate the ordinary traffic with a
small script in your scratchpad, and place the incident's lines by hand, as
[kinds-of-drill.md](references/kinds-of-drill.md#log-driven) describes.
Make the metrics from the logs. Keep the script out of the drill.

### 7. Ask once

Write in the chat: the kind of drill, and why; the case, and the option it
leads to; the setup; the
plot, with each choice, its options, and where each leads; the list of
pages, with the first lines of each log; the expected response; and the
open questions. Then ask, in one call of your question tool:

* **whether to write the drill**: **Write it** as it is, or **Change it**,
  with the change typed in as Other;
* if there is no log of drills, whether to start one with `fyr init`;
* if a drill of the same case exists, whether the new one replaces it,
  updates it, or stands beside it;
* and, in the same call, up to three questions research could not settle,
  such as which of two sources holds, or a fact about the team's tools that
  no document records.

Revise for any change, and play it back if it changes more than a detail.

### 8. Write the drill

Create the drill with the playbook's title, exactly as its title line gives
it after the number, as rule 7 says. For `# 4. Customer charged twice`:

```bash
VISUAL= EDITOR= fyr new "Customer charged twice"
```

Clearing `VISUAL` and `EDITOR` stops the tool waiting on an editor. It prints
the path of the new file, such as
`docs/firedrill/FIREDRILL-0003-customer-charged-twice.md`,
and creates the directory beside it, with `problem/README.md` and
`solution/README.md`.

A new drill may replace or update one already written. If so, tell the user
that the tool records the link in **both** drills at once, and that with `-r`
the old drill loses its `Draft` line. Then:

* **replace**: `VISUAL= EDITOR= fyr new -r 2 "Customer charged twice"`, for a new drill of the
  case drill 2 practised, such as after the playbook changed;
* **update**: `VISUAL= EDITOR= fyr new -u 2 "Customer charged twice"`, for a drill that changes
  part of drill 2, such as a choice it left out.

Then:

1. Edit the new document in place: keep the lines the tool wrote (the title,
   the date, `Draft`, and any link lines in the Status section), and replace
   the guidance text with the drill. Never retype the title: its number is
   the tool's, and must stay the number in the file name. Never write link
   lines by hand.
2. Write the problem pages in `problem`, and the expected response and the
   debrief in `solution`, beside the READMEs.
3. Leave the READMEs as the tool wrote them.

### 9. Finish

* Read every file back, and check the drill against the list at the end of
  [writing-firedrills.md](references/writing-firedrills.md#before-you-finish).
* Follow every path from the first step: each choice leads to a page that
  exists, and every path ends at the solution.
* Solve the drill yourself, from the incident, knowing only what each step
  gives. Run the searches the expected response describes, such as `grep`
  for the order number, on the drill's own files, and check that they find
  what it says. Then check that nothing finds the answer sooner: a file
  name, an option's wording, a line that says what it means.
* Check that the number in the title, `# N.`, is the number in the file
  name, `FIREDRILL-000N-...`, and that the rest of the title is the
  playbook's, exactly. If either differs, put the tool's title line back.
* Run `fyr list` to confirm the drill is in the log.
* The project may keep a table of contents made by `fyr list --toc`, such as
  a `CONTENTS.md` in the log directory. If so, offer to regenerate it. Never
  write one over the `README.md` that `fyr init` wrote.
* Tell the user:
  * the path of the document, and that it is a draft;
  * any drills the tool changed;
  * that the reader should start at the document, and read the `solution`
    only when the drill is over;
  * the open questions, and anything research found wrong or missing in
    the playbook;
  * that after the drill, `create-playbook` can revise the playbook from
    the debrief.

  Leave committing to them.

## Tips

These are the mistakes drills most often make. See
[examples.md](references/examples.md) for each, good and bad.

* **Hide the answer from everything but the evidence.** Not in the title,
  the file names, the document's steps, the options' order or wording, or a
  comment in a log. The reader finds it, or does not.
* **Write what the system does, not what is typical.** A log line no code
  writes, an alert no rule defines, or a provider's error no documentation
  gives makes a drill that practises a system that does not exist.
* **Real incidents do not announce themselves.** No line says "duplicate
  charge". The evidence is a retry, two events seconds apart, a request that
  never finished, or a line that is missing.
* **Make the evidence big enough to search.** A dozen lines, all relevant,
  is a reading exercise. Logs are plain files, with ordinary traffic and the
  system's everyday warnings around the incident.
* **The reporter knows what they saw.** A customer has their bank app, not
  the charge IDs. A user says "deleted" when it was moved. Their words are
  the pressure; what they paste is the evidence.
* **Every option is tempting.** Write each one as a competent person under
  pressure would do it, alike in length and tone. An absurd option is a free
  answer.
* **A page tells the story, not the verdict.** Say what happens next, in the
  system's terms, and what a choice costs. Leave "why" to the solution.
* **One clock, one incident.** Identifiers, times and amounts agree on every
  page, and metrics are made from the logs, so that a reader who checks one
  against another finds them consistent.

## Edge cases

* **The user wants a drill for a problem with no playbook.** Suggest writing
  the playbook first. A drill tests a plan; without one there is nothing to
  follow, and nothing to evaluate.
* **The user wants a drill for a team**, with someone to run it. The drill
  is still written for one reader; a team can read it together, choosing as
  one. Say so.
* **The user asks to make the problem happen for real**, such as stopping a
  service or sending bad data. Decline: describe the symptoms in the pages.
  If they want the system itself to show the symptom, leave that to them, as
  a note in the setup, and do not do it.
* **The playbook names evidence the system does not produce.** Say so in the
  debrief, and in your summary. Do not invent the evidence.
* **Research finds no local setup.** Say where you looked. Run the drill as a
  tabletop exercise on the pages, and say so in the setup.
* **You cannot read the documentation of a system involved**, because you
  have no web tool, or the pages are not public. Say so, ask the user for
  the documentation or a link, or leave its evidence as an open question.
  Do not write its events or errors from memory.
* **The sources disagree**, such as the code and the playbook on what is
  logged. The code is the better evidence of what the system writes; build
  the pages from it, and note the playbook's error for the debrief.
* **The `plm` command is missing.** Ask the user for the path of the
  playbook, and read it from there.
* **The `fyr` command is missing.** Tell the user. Do not create drills by
  hand, since the numbering and links would not be maintained.
* **The tool reports an error.** Show the user the message, which begins
  `fyr:`, rather than working around it.
