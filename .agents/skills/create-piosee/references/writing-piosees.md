# Writing a PIOSEE record

A PIOSEE record is read twice. First by the people making the decision, who
need the same facts and options in front of them; then by whoever comes later,
to learn what was known, what was weighed, what was done, and whether it
worked. Both readers take each statement as true. A fact nobody checked, an
option nobody considered, or a reason nobody gave misleads the first reader
into a worse decision, and the second into a wrong lesson.

So the problem, the options, the choice and its reasons are the user's. You
gather the facts you can, and ask for the rest. See [piosee.md](piosee.md)
for what each step is for.

## Facts first

Before writing any prose, build a fact sheet: every claim the record will
make, and where it came from. Keep it in your working notes. It is shown to the
user with the draft, and is not written into the record.

| #   | Claim                                                         | Source                                    |
| --- | ------------------------------------------------------------- | ----------------------------------------- |
| 1   | Checkout errors rose from 0.1% to 4% at 14:05                 | `checkout_errors` dashboard, read 14:20   |
| 2   | Release 2.3 was deployed at 14:02                             | `git log --tags`, and the deploy log      |
| 3   | Release 2.3 includes a migration that adds a column           | `migrations/0042_add_tax_code.sql`, read  |
| 4   | Problem: customers cannot pay; the sale starts at 16:00       | User, in this conversation                |
| 5   | Options: roll back, fix forward, turn off the tax-code flag   | User, in this conversation                |
| 6   | Rolling back would need the migration reversing first         | User, in this conversation                |
| 7   | Selected: turn off the flag, because it is fastest and safest | User, in this conversation                |
| 8   | Who turns off the flag                                        | **Unknown — ask**                         |

### Acceptable sources

* **The user**, in this conversation: what they stated, or explicitly agreed
  to. A suggestion of yours that they did not take up is not a fact.
* **The repository**: code, configuration, runbooks, documents and earlier
  records that you have read. Note the path, and the line where it matters.
* **Commands you ran**, and their output: logs, metrics, status pages, the
  state of a deployment. Note the time for anything that changes.
* **A document the user pointed you to**, such as an incident ticket, a chat
  thread, a supplier's email or a timeline.

### Not acceptable

Do not write any of these unless a source above supports it:

* a fact about the problem: a cause, a time, a number, an impact, a number of
  people affected;
* an option the user did not consider, or a risk or benefit they did not see;
* the option selected, who selected it, or the reason;
* an action, who took it, or when;
* an outcome, or a lesson learned;
* a prediction: that an option "will fix it", "is low risk".

A cause you inferred from the evidence is a hypothesis, not a fact. Put it to
the user as one: "The errors began three minutes after the deploy. Do you
think the release caused them?" Their answer is the evidence, and the record
says it was their judgement, not a proven cause.

### Your own suggestions

The user may ask for your help: "what are our options?", "what are the risks
of rolling back?", "what would you do?". That is a fair request. Answer it in
the conversation, labelled as your suggestion, grounded in what you read and
ran.

A suggestion enters the record only when the user adopts it; it is then
theirs, and the fact sheet says so ("User adopted the agent's suggestion, in
this conversation"). An option the user adopted from your list is an option
they considered; it does not make your view of it their reason for selecting
it. Ask for their reason.

Do not volunteer suggestions in place of asking. Ask the open question first;
offer ideas only if the user asks for them, or says they do not know and would
like some.

### When a fact is missing

Ask the user. If they do not know, or prefer not to say, do one of these, and
let them choose:

* **Record it as unknown**, in Information: "The cause of the errors was not
  known when the decision was made." What was not known is part of the
  decision, and often the most useful thing a later reader learns.
* **Find it out**, if there is time: who could be asked, or what could be run.
* **Leave it out**, where the template does not need it.

Never fill the gap with something plausible.

Three gaps cannot be filled this way. If the user cannot say **what the
problem is**, **which option was selected**, or **why**, there is no decision
to record yet. Tell them so. If the problem is unclear, working out what it
is, is the first decision.

## The sections

The built-in template has a title, a date, a status, and one section for each
step. A project's own template may differ: fill in each of its sections by
the guidance written in it, with the same standard of evidence.

### Title

A short phrase naming the decision: "Roll back release 2.3", "Turn off tax
codes during the sale", "Move the launch to March". It names what was chosen,
if that is known when the record is made; otherwise the problem: "Respond to
checkout errors after release 2.3". Run `pse list` and follow the style of
the existing titles. The title becomes the file name, so keep it short.

### Date

The tool writes today's date. For a decision already made, ask whether the
user wants the date it was made instead, and if so set `PSE_DATE`.

### Status

Always `Draft`. The tool writes it; leave it. What becomes of the decision is
for the team to record.

### Problem

What is wrong, or what must be decided, in a few sentences: what is happening,
what it affects, how badly, and how long there is to act. In the user's words,
agreed by them.

* State the problem, not the solution. "Customers cannot pay", not "we need to
  roll back".
* Include the time available, and what happens when it runs out: "The sale
  starts at 16:00, in under two hours."
* If the cause is not known, say so. Do not put a guessed cause in the problem.

### Information

The facts that bear on the problem, each with its source, as a short list.

* What is known: measurements, times, recent changes, what the runbook says.
* What is not known, and whether anyone is finding out.
* Who was asked, or consulted, by role: "the payments team", "the on-call
  database engineer".
* Facts only. Interpretation belongs in Options and Select.

### Options

Each course of action considered, as a list, each with the risks and benefits
the user saw.

* Include doing nothing, and buying time, if they were considered. If they
  were not, do not add them.
* Write each option as an action: "Roll back to release 2.2."
* Risks and benefits are the user's assessment. Keep their weight: "probably"
  stays "probably".
* For a decision already made, only the options considered at the time.

### Select

The option chosen, who chose it, and why it was chosen over the others.

* Name the option exactly as it appears under Options.
* Who chose it, by role, if the user wants that recorded: "the incident lead".
* The reason, in the user's words. If it was a judgement call, or instinct,
  say so.

### Execute

How the option is carried out.

* The actions, in order, with who takes each, by role.
* What shows that it is working: the signal the team is watching.
* For a decision already made, what was actually done, with times where they
  are known. If what was done differed from what was chosen, say so.

### Evaluate

Whether the action solved the problem.

* **Before it can be evaluated**, write how and when it will be: "Evaluate at
  16:30: checkout errors below 0.5% for 30 minutes." Then: "Not yet
  evaluated."
* **When it is evaluated**, write what happened, with its sources; whether the
  problem was solved; and what was learned. Options seen only in hindsight go
  here, labelled as such.
* If the problem remains, or has changed, say that the decision is to be
  revisited, and in which record, once it exists. The link is made by the
  tool; see [pse-tool.md](pse-tool.md#links).

If there is nothing to say under a heading the template does not require,
leave the heading out. If a step was skipped, write "Not recorded."

## Style

* Plain words, in the team's terms. "The payments service", not
  "`svc-pay-prod-eu1`", unless the name is what readers will search for.
* Short. Under time pressure, a sentence per step is enough.
* Past tense for what has happened, present for what is true now, future only
  for what the team has agreed to do.
* No blame. Record what was known and done, not who should have known better.
* No placeholders and no template guidance left behind.

## Examples

Each pair shows what to write instead, if the evidence supports it.

* Instead of the problem "The release broke checkout", write "Since 14:05,
  4% of checkout attempts fail, against 0.1% before. The cause is not yet
  known." and let the user's judgement of the cause go in Select.
* Instead of adding "Do nothing" to the options because every record has it,
  ask whether it was considered. If it was not, leave it out.
* Instead of the reason "Turning off the flag was the lowest-risk option",
  ask why they chose it, and write theirs: "The incident lead chose it because
  it can be undone in a minute, and does not touch the database."
* Instead of evaluating "The fix worked" before it has run, write when and how
  it will be evaluated, and "Not yet evaluated."
* Instead of listing, after the event, an option nobody thought of at the
  time, put it in Evaluate: "In hindsight, the tax-code service could have
  been put into read-only mode."

## Before you finish

* Every sentence traces to a line of the fact sheet.
* The problem, the option selected, and the reason all come from the user.
* Every "Unknown — ask" line was asked about, and the answer is used, recorded
  as unknown, or left out, by the user's choice.
* No guessed cause is stated as a fact.
* The options are those considered at the time; hindsight is in Evaluate.
* Evaluate does not claim an outcome that has not happened.
* Every section of the template in use is filled in, and no guidance text is
  left.
* The status is `Draft`.
* Links between records were made by the tool, not by hand.
* The user has seen the draft with its sources, and approved it.
