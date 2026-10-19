# Writing a fire drill

A fire drill is run by someone who did not write it, for a team that should
not have to wonder whether the evidence in front of it is real. Its logs
should read like the system's own, so that the team learns to read the
system's own. Its symptoms should be the ones the system would really show,
so that the team practises on what it will meet. And nothing in it should
reach anything real.

So the scenario and the response come from the playbook. The local setup,
the formats and the messages come from the code. What the other systems
report comes from their documentation. Only the particulars of the incident
are made up.

## Facts first

Before writing the document or any asset, build a fact sheet: each fact the
drill relies on, and where it came from. Keep it in your working notes, show
it with the draft, and put its sources in the document's **Sources**
section.

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

Every line of every asset traces to an entry in it, except the particulars
of the incident.

## The fire drill document

`FIREDRILL.md` is for whoever runs the drill. Write it in the present tense
and the imperative, in these sections.

### Title, date and status

```markdown
# Fire drill: <title of the playbook>, <the case>

Date: <today, YYYY-MM-DD>

Status: Draft
```

Below the status, link the playbook, relative to the document, with its
number and title as its own title line gives them, such as
`[4. Customer charged twice](../../problem/PROBLEM-0004-customer-charged-twice.md)`.

### Purpose

What the drill practises, in two or three sentences: the playbook, the case,
and the option the evidence should lead to. Name any part of the playbook
the drill sets out to test, such as a new immediate action.

### Setup

How to prepare, before the team arrives:

1. **The local system.** The exact commands that run it locally, from the
   project's own instructions, with their source, and any seed data the
   case needs, loaded the way the project loads it. No fault is added. If
   the system cannot be run locally, say that the drill is a tabletop
   exercise on the assets.
2. **The materials.** Which assets to hand out, and when: the incident
   description at the start; the others when the team looks for them, or all
   at once, as the setup suits.
3. **The team.** Who takes part: those who would follow the playbook.
4. **The time.** How long to allow, and when to stop.

### Ground rules

* This is a drill. Every message sent during it begins `FIRE DRILL:`, and
  goes only to those taking part.
* Act only on the local system. A step that would act on anything else is
  described, not taken; the document says which steps those are.
* Anyone can stop the drill at any time, for a real incident or any other
  reason.

### Scenario

How the drill starts: the incident description, as the team receives it,
and how it arrives, such as an alert, a customer's email or a GitHub issue.
Link `assets/incident.md`.

### Symptoms

What the team will find wherever it looks, one entry for each source the
playbook's Information names and each the team might reasonably try:

* **Where**: the log, record, dashboard, page or person.
* **What it shows**: the symptom, in the system's own terms.
* **The asset**, if one stands in for it, linked; or, where there is none,
  what to tell the team if it asks, such as what a provider's dashboard
  would show.

Name, too, any source that shows nothing, and say so: "The webhook log has no
entry for the second delivery: the handler does not log it."

### Expected response

The response the playbook leads to in this case, section by section:

* **Immediate action**: what the team should do at once, or that none is
  indicated.
* **Counter-indicated actions**: what to watch for the team doing, and why
  it would make things worse.
* **Information**: what the team should find, and where, in the symptoms
  and assets.
* **Selection**: which option the evidence leads to, by which rule.
* **Execution**: each step of that option, and for each, whether the team
  **takes it** on the local system or **describes it**, because it would act
  on something outside it.
* **Evaluation**: how the team should judge the result, and what the drill
  can show of it.

Quote or link the playbook; do not rewrite its rules.

### Debrief

Questions for the team after the drill, following the Detect, Identify,
Communicate and Resolve loop of the practice:

* **Detect**: how and when did the team learn of the problem? Would it have
  learned sooner in production?
* **Identify**: did the playbook's Information lead to the evidence? What
  was missing, misleading or slow to find?
* **Communicate**: who was told, how, and when? Who should have been?
* **Resolve**: did the Selection rules lead to one option, clearly? Did each
  step of Execution work, or would it have?

Then list what research already found wrong or missing in the playbook: a
source that does not hold what it says, evidence the system does not
produce, a step that would not work. These, and what the team found, are for
revising the playbook with `create-playbook`.

### Open questions

Each fact research could not settle, and what it would change in the drill.
Leave the section out if there are none.

### Sources

The fact sheet's sources: each fact the drill relies on, and where it came
from, with links to the documentation read.

## The assets

The assets are what the team works from, so they carry no answers: no note
of the expected option, no comment on what a log line means.

### The incident description

`assets/incident.md` is the alert, message or report that starts the drill,
as the team would receive it, in the form the Scenario names: an alert's
text, a customer's email, an issue. It begins `FIRE DRILL:`, so that it
cannot be mistaken for a real one if it is forwarded. It holds what the
reporter would really know, in their own words, and no more: a customer
knows their bank statement, not the charge ID.

### The logs

One file for each source, named for it, such as `assets/logs/api.log`.

* **The real format.** The format, field names, timestamp format, time zone
  and levels of the logger's configuration. For a command-line tool, its
  output and error messages, as a terminal shows them.
* **The real messages.** Only lines the code writes, verbatim, with the
  fields it attaches. A message the code does not write does not appear.
* **The other systems' own words.** Event types, payloads, error codes and
  messages as their documentation gives them, for the version the project
  uses.
* **One timeline.** One clock across every asset. The symptoms begin before
  the incident description arrives; the lines that matter sit among
  ordinary traffic from the same code, such as other requests and health
  checks, enough that the team has to search, not so much that the search
  is the exercise.
* **One incident.** The same order, customer, charge and request
  identifiers wherever they appear, so that the team can follow one from
  system to system, as the playbook's Information says to.

### Other evidence

Anything else the Information section says to check, written as a file in
`assets/`: a record as a query would return it, an event as the provider's
dashboard shows it, a file as it would be on disk. Use the same formats and
identifiers as the logs.

### The particulars

* Identifiers in the real formats, such as the provider's `ch_` prefix or a
  UUID, but made up.
* People and addresses that cannot be real: `example.com` addresses, and
  names that are plainly made up.
* Card numbers, account numbers and other payment details only from the
  provider's documented test values.
* No real customer data, secrets, tokens or production log lines, even
  redacted.

## Before you finish

* Nothing in the drill breaks, stops or changes any system, and nothing was
  run to make it.
* Every step of Execution that would act outside the local system is marked
  **describes it**.
* The incident description, and every message the drill has the team send,
  is marked `FIRE DRILL:`.
* The setup's commands are the project's own, with their source.
* Every log line's format and message traces to the code, and every other
  system's event or error to its documentation.
* Identifiers, times and amounts agree across every asset.
* The assets hold no answer, and no real data.
* The expected response follows the playbook as it stands, and leads to one
  option by its Selection rules.
* The link to the playbook works, and its text matches the playbook's title
  line.
* The debrief lists what research found wrong or missing in the playbook.
