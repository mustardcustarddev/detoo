# Writing a fire drill

A fire drill is run by one person, who did not write it, alone. No one is on
hand to play the reporter, hand out the evidence or say what the playbook
expected, so the drill does all of that itself: it is read like a Choose Your
Own Adventure book. Its logs should read like the system's own, so that its
reader learns to read the system's own. Its symptoms should be the ones the
system would really show, so that the reader practises on what they will
meet. Its choices should be the ones a person under pressure would really be
tempted by. And nothing in it should reach anything real.

So the scenario, the decisions and the response come from the playbook. The
local setup, the formats and the messages come from the code. What the other
systems report comes from their documentation. Only the particulars of the
incident are made up.

## Facts first

Before writing the document or any page, build a fact sheet: each fact the
drill relies on, and where it came from. Keep it in your working notes, and
show it with the draft. The drill itself has no sources section.

| #   | Fact                                                                    | Source                                                    |
| --- | ----------------------------------------------------------------------- | --------------------------------------------------------- |
| 1   | Scenario: a customer reports two charges for one order                  | `docs/problem/PROBLEM-0004-...md`, Scenario               |
| 2   | Run locally with `make dev`, which starts `api` and `worker`            | `Makefile`, target `dev`; `docker-compose.yml`            |
| 3   | Logs are JSON, with `time` (RFC 3339, UTC), `level`, `msg`, `order_id`  | `internal/log/log.go`, lines 12–30                        |
| 4   | A charge timeout logs `charge timed out, retrying` at `warn`            | `checkout/service.go`, line 61                            |
| 5   | The provider sends `charge.succeeded` for each charge, with its payload | Provider's webhook reference, for the pinned API version  |
| 6   | Charge IDs begin `ch_`                                                  | Provider's API reference, the charge object               |
| 7   | Nothing logs when a webhook is received twice                           | `webhooks/handler.go`, read: no log call on that path     |
| 8   | Where support sees the customer's charges                               | **Unknown — open question**                               |

Every line of every page traces to an entry in it, except the particulars
of the incident.

## The layout

`fyr new` creates the document, and a directory beside it, named as it is
without `.md`, with a `problem` and a `solution` directory, each with a
`README.md`. Leave the READMEs as the tool wrote them. Write the rest:

```text
FIREDRILL-0003-customer-charged-twice.md
FIREDRILL-0003-customer-charged-twice/
  problem/
    incident.md             the incident that starts the drill
    logs/api.log            the evidence, one file for each source
    refund-at-once.md       a page for each choice
    ...
  solution/
    expected-response.md
    debrief.md
```

From the document, link a page by its path under the directory, such as
`FIREDRILL-0003-customer-charged-twice/problem/incident.md`.
From a page, link the document as `../../FIREDRILL-0003-...md`, another page
in `problem` by its name, and the solution as `../solution/...`.

## The plot

Plot the drill before writing it: the steps, the choices, and the page each
choice leads to.

* **The steps** follow the incident as it would unfold: it starts, the
  reader looks for the evidence, more arrives, the reader decides, acts, and
  checks. One step for each, numbered.
* **A choice** comes wherever the playbook makes a decision: whether to take
  the immediate action, whether the evidence is enough to select an option,
  which option it selects, and anywhere the case tempts the reader off the
  playbook, such as a reporter who demands an action.
* **The options** at a choice are two to four, each something a competent
  person under pressure might do:
  * the one the playbook leads to;
  * a primary risk the playbook warns of, or a counter-indicated action;
  * another of the playbook's options, which the evidence does not yet
    support;
  * a shortcut, such as acting on the reporter's word, or skipping a step
    of the Information.

  Write each in the reader's terms, as an action, such as "Revoke the
  release", not "Option 1". At the choice of an option, name the playbook's
  options, since choosing between them is the skill practised. Do not always
  put the playbook's option first, and give no hint of which it is.
* **Each option leads to one page.** A page may continue the story to
  another choice, send the reader on to a step of the document, send them
  back to the choice to choose again, or end the drill. Every path ends at
  the solution. A path the playbook would not take may end the drill, when
  what follows is the harm it does, or send the reader back, when the
  playbook gives them a way to recover.

## The fire drill document

The document `fyr new` creates is where the reader starts. Write it to them,
in the second person and the present tense, in the sections of the template
in use. These are the built-in template's; a project's template may differ.

### Title, date and status

`fyr new` writes these. Keep them as the tool wrote them: the title is the
drill's number and the playbook's title, exactly as the playbook's title line
gives it after its number, such as `# 3. Customer charged twice` for a drill
of `# 4. Customer charged twice`. Nothing is added to it, not even the case.

### Playbook

Link the playbook, relative to the document, with its number and title as
its own title line gives them, such as
`[4. Customer charged twice](../problem/PROBLEM-0004-customer-charged-twice.md)`.
Say that the reader should have it open throughout.

### How to play

Keep the template's text. Add the case, in a sentence that gives nothing
away, such as "A customer writes that they were charged twice", since the
title cannot hold it; who the reader plays, such as the maintainer or the
engineer on call; and how long the drill takes.

### Ground rules

Keep the template's rules, and add any the case needs, such as where a
command that changes files may be run safely.

### Setup

What the reader prepares before the first step, in order:

1. **The local system.** The exact commands that run it locally, from the
   project's own instructions, and any seed data the case needs, loaded the
   way the project loads it. No fault is added. If the system cannot be run
   locally, say that the drill is a tabletop exercise on the pages.
2. **What to have to hand**: the playbook, the SOPs it refers to, and
   somewhere to write down each choice and each message.

Name no version: where the setup needs one, such as the release a reporter
has, give the command that finds it, such as the latest release tag.

### The drill

The script: numbered steps, each saying what happens and what the reader
does.

* **A step that happens** says what arrives or what the reader finds, and
  links the page that holds it: "A new issue is opened. Read
  [the issue](FIREDRILL-.../problem/incident.md)."
* **A step that asks for a choice** says what the reader faces, then gives
  the options as a list, each linking its page: "What do you do first?"
* **A step the reader takes on the local system**, such as reproducing a
  problem, says what to do, not what they will find: that is theirs to see.
* **A step that would act outside the local system** has the reader write it
  down, such as the comment they would post, and keep it.

A step reached only from a page says so, such as "From the reporter's
reply", so that the reader who reads on does not take it in order. Give no
answers in the document: what the evidence shows, and which choice is the
playbook's, are for the pages and the solution.

## The problem pages

The pages are what the reader works from, and what happens to them. Name each
for what it holds or what happens, such as `incident.md`, `logs/api.log` or
`revoke-the-release.md`, never for whether it is right, such as `wrong-1.md`.

Every page ends with where to go next. On a page of evidence, put it after a
rule (`---`), apart from the evidence: "Go on to step 2 of the drill." A log
file, which cannot hold a link, needs none: the step that sent the reader
there says where to go after it.

Files the reader copies to the local system, such as records to reproduce a
problem with, go in a directory of their own under `problem`, and the step
that uses them gives the commands to copy them somewhere safe first.

### The incident

`incident.md` is the alert, message or report that starts the drill, as the
reader would receive it, in the form the Scenario names: an alert's text, a
customer's email, an issue. It begins `FIRE DRILL:`, so that it cannot be
mistaken for a real one. It holds what the reporter would really know, in
their own words, and no more: a customer knows their bank statement, not the
charge ID.

### The evidence

The logs, records and other evidence the Information section names, each a
page the drill sends the reader to when they look for it, or that a choice
page reveals, such as the reporter's reply to a question. They carry no
answers: no note of the expected option, no comment on what a line means.
How much evidence there is, and how to make it, depends on the kind of drill:
see [kinds-of-drill.md](kinds-of-drill.md). For good and bad examples, see
[examples.md](examples.md).

* **The real format.** The format, field names, timestamp format, time zone
  and levels of the logger's configuration. For a command-line tool, its
  output and error messages, as a terminal shows them. Leave out any line
  of the output that would name an exact version, as a person pasting it
  might.
* **The real messages.** Only lines the code writes, verbatim, with the
  fields it attaches. A message the code does not write does not appear.
* **The other systems' own words.** Event types, payloads, error codes and
  messages as their documentation gives them, for the version the project
  uses.
* **One logs file for each source**, named for it, such as `logs/api.log`:
  a plain file, as the logger writes it, not a Markdown page.
* **One timeline.** One clock across every page. The symptoms begin before
  the incident arrives; the lines that matter sit among ordinary traffic
  from the same code, such as other requests and health checks, enough that
  the reader has to search, not so much that the search is the exercise.
* **One incident.** The same order, customer, charge and request
  identifiers wherever they appear, so that the reader can follow one from
  system to system, as the playbook's Information says to.
* **A record or a dashboard** is written as a query would return it, or as
  the dashboard shows it, in the same formats and identifiers as the logs.

### The choice pages

A choice page says what happens after the choice, in the system's own terms
and the people's own words: what the reader does, what it changes, what
arrives next. It is the story, not the verdict:

* A choice the playbook would make shows its result, and moves the story on.
* A choice it would not shows what that would really cost, such as a
  release withdrawn from everyone, or a customer refunded twice, plainly and
  without blame. It may name the part of the playbook to read again, such as
  its primary risks; the full reasoning is for the solution.
* A step that would act outside the local system has the reader write it
  down, and the page says what would happen if it were taken.

A choice page ends with what to do next: make another choice, go to a step
of the document, go back to the choice, or, at the end of a path, "The drill is
over. Read [the solution](../solution/expected-response.md), then
[the debrief](../solution/debrief.md)."

## The solution

### expected-response.md

The response the playbook leads to in this case, section by section, for the
reader to check their own against:

* **Immediate action**: what to do at once, or that none is indicated.
* **Counter-indicated actions**: what the drill tempted the reader to do,
  and why it would make things worse.
* **Information**: what to find, and where, among the pages and on the local
  system.
* **Selection**: which option the evidence leads to, by which rule.
* **Execution**: each step of that option, and for each, whether the reader
  **takes it** on the local system or **writes it down**, because it would
  act on something outside it. A message, such as a reply to a customer, is
  sketched, with what a good one says.
* **Evaluation**: how to judge the result, and what the drill can show of
  it.

Then **the path**: the playbook's choice at each step of the drill, and for
each other option, why the playbook does not take it, linking its page.
Quote or link the playbook; do not rewrite its rules.

### debrief.md

Questions for the reader to look back on the drill with, following the
Detect, Identify, Communicate and Resolve loop of the practice, and asking
about this drill's own choices:

* **Detect**: how and when did you learn of the problem? Would you have
  learned sooner in production?
* **Identify**: did the playbook's Information lead you to the evidence?
  What was missing, misleading or slow to find?
* **Communicate**: who would you have told, how, and when? Did what you
  wrote say what the solution's sketch says?
* **Resolve**: did the Selection rules lead you to one option, clearly? Did
  each step of Execution work, or would it have? Which choices did you get
  wrong, and why?

Then **What the drill found in the playbook**: what research found wrong or
missing in it: a source that does not hold what it says, evidence the
system does not produce, a step that would not work. These, and what the
reader found, are for revising the playbook with `create-playbook`.

Then **Open questions**: each fact research could not settle, and what it
would change in the drill. Leave the section out if there are none.

## The particulars

* Identifiers in the real formats, such as the provider's `ch_` prefix or a
  UUID, but made up.
* People and addresses that cannot be real: `example.com` addresses, and
  names that are plainly made up.
* Card numbers, account numbers and other payment details only from the
  provider's documented test values.
* No real customer data, secrets, tokens or production log lines, even
  redacted.
* No exact version: "the current release", "the release before it".

## Before you finish

* Nothing in the drill breaks, stops or changes any system, and nothing was
  run to make it.
* Every step that would act outside the local system has the reader write it
  down, not take it.
* The incident, and every page that looks like a real message, begins
  `FIRE DRILL:`.
* The setup's commands are the project's own.
* Every log line's format and message traces to the code, and every other
  system's event or error to its documentation.
* Identifiers, times and amounts agree across every page, and every metric
  agrees with the logs it was made from.
* Every log file is a plain file in the logger's format, with ordinary
  traffic around the incident, and no line in it says what it means.
* An alert is the project's own rule, and each stage's files end at its
  time.
* You solved the drill from its own pages, and the searches the expected
  response describes find what it says.
* No exact version appears anywhere in the drill.
* There is no sources section.
* The document and the evidence hold no answer, and no page holds real
  data.
* Every choice leads to a page that exists, every page ends with what to do
  next, and every path ends at the solution.
* The expected response follows the playbook as it stands, and leads to one
  option by its Selection rules.
* The drill's title, after its number, is the playbook's, exactly.
* The link to the playbook works, and its text matches the playbook's title
  line.
* The debrief lists what research found wrong or missing in the playbook.
