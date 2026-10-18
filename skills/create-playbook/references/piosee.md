# PIOSEE problem playbooks

## What a playbook is

A problem playbook is a plan, written in advance, for responding to a
problem when it arises, such as an incident. It follows PIOSEE, a
decision-making method from aviation, named for its steps: **P**roblem,
**I**nformation, **O**ptions, **S**election, **E**xecution and
**E**valuation. They are prefaced by the **Scenario** the playbook is for.
After the Problem come the **Immediate action**, and, when any are known,
the **Counter-indicated actions**.

PIOSEE is a way of acting, not of documenting. A crew whose engine fails
after take-off first keeps the aircraft flying, by drills learned for that
moment, and only then works out where to land. A playbook does the same: its
immediate action makes the situation safe, or stops it getting worse, before
anyone gathers information. It says which actions would make things worse.
Then it leads the reader through deciding and acting.

The problem may be one that recurs, such as an outage, or a customer who
cannot cancel a subscription. It may be one that is scheduled or foreseen,
such as crossing a tax registration threshold on a known date. It may be one
that has happened once and could happen again. What these have in common is
that someone will one day have to act, and decide what to do, often under
pressure. The playbook means the thinking has already been done, calmly, by
the people who know the system and the policy.

When the problem arises, the person on hand does the following:

1. Recognises the scenario.
2. Takes the immediate action, if there is one, and avoids the
   counter-indicated actions throughout.
3. Gathers the information it lists.
4. Applies its selection rules to choose an option.
5. Follows the execution steps.
6. Later, judges the result by its evaluation measures.

They decide in the moment, but by rules agreed beforehand. So similar cases
are decided alike, and each step leaves an audit trail.

## Playbooks and SOPs

A **standard operating procedure** (SOP), kept with the `sop` tool, is the
steps of a routine task, such as rolling back a release, pausing a queue or
reading a service's logs. It has no decision in it: it is done the same way
each time.

A playbook decides; an SOP does. Where a playbook needs a routine task done,
it refers to the SOP for it, by a link, rather than repeating the steps:

* the **Immediate action** may be an SOP, such as one that turns on a
  maintenance mode;
* an **Information** item may say to collect the evidence by an SOP, such as
  one for searching a service's logs;
* an **Execution** step may be an SOP, such as one for rolling back a
  release;
* an **Evaluation** check may be an SOP, such as one for confirming that
  payments are flowing again.

So an SOP is written, reviewed and corrected once, and every playbook that
uses it stays current with it.

## What a playbook is not

* **Not a record of an incident.** It does not say what happened on a given
  day, who decided what, or how it turned out. A postmortem or an incident
  ticket does that. A playbook is written for every occurrence of a problem.
* **Not a decision already taken.** An ADR records a decision about how the
  system is built, made once. A playbook's Selection is a rule for deciding
  each time. It may name a default, but the choice is made when the scenario
  arises, against the evidence gathered then.
* **Not only a list of steps.** An SOP or a runbook says what to do. A
  playbook also says what to do at once, what to avoid, what to find out
  first, what the options are, how to choose between them, and how to tell
  whether it worked.

## The sections

The order follows how the problem will be handled. Each section builds on the
one before.

| Section                       | What it holds                                                                                                                                              |
| ----------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Scenario**                  | The triggering situation, in one to three sentences, with any standing conditions that frame it.                                                           |
| **Problem**                   | The exact decision to make, as one clear question, and the primary risks if it is decided wrongly. No solutions.                                           |
| **Immediate action**          | The action to take at once, before gathering information, to make the situation safe or stop it getting worse, and who takes it, or the SOP for it. Or "No immediate action is indicated." |
| **Counter-indicated actions** | Optional. The actions known to make the situation worse, or to waste time, that someone under pressure may be tempted to take, and why each is wrong.     |
| **Information**               | The facts to collect before deciding: what, where each is found, and who owns it, or the SOP for collecting it. Known facts apart from assumptions. Critical unknowns, and how to settle them quickly. |
| **Options**                   | The realistic courses of action, distinct and actionable, with at least one conservative fallback. For each: when it applies, and one or two tradeoffs.  |
| **Selection**                 | How to choose: explicit criteria, which findings from Information lead to which option, the default, the confidence in the rule, and its key assumptions. |
| **Execution**                 | The steps, in order, with owners and dependencies, referring to SOPs for routine tasks. What to record and communicate. The rollback path. Steps to prevent a recurrence. |
| **Evaluation**                | How success is measured once the playbook has been used: measurable outcomes, the SOPs for routine checks, when and by whom they are reviewed, and follow-up actions if targets are missed. |

A playbook may add sections that help the person following it, such as
**Communication Templates** with the message to send for each outcome, or
**Reference Links** to the documentation of the systems involved.

### The immediate action

The immediate action is taken before anyone has looked at the evidence, so
it must be right whatever the evidence turns out to be. It is usually small,
safe and reversible: turn off a feature flag, pause a job, stop sending
emails, put up a status notice. An action that is right only in some cases
is an option, chosen by the selection rules, not an immediate action.

Many problems need nothing at once: nothing gets worse while the information
is gathered. Then the playbook says so, in those words: "No immediate action
is indicated." That is a finding, not a gap, and it tells the reader they may
take their time to investigate.

### Counter-indicated actions

These are the actions that seem right under pressure, and are not: a restart
that destroys the evidence, a retry that charges the customer again, a
rollback that loses data written since the release, a refund issued before a
dispute is settled. Each names the action and why it is wrong. The section
is optional: write it only when such actions are known, from the
documentation, the code, an earlier incident, or the user.

## What makes a good one

* **Actionable.** It is written to be followed, not read. The first thing it
  tells the reader to do, they can do at once.
* **Concrete.** It names the actual log source, storage path, record field,
  dashboard page and menu path, or the SOP that does. An Information section
  that says "check the logs" is no help at two in the morning. One that
  names the service whose logs to read, and the events to look for, is.
  Concrete is not prescriptive: it names a command where the command is how
  to find a fact, and leaves everyday actions, such as replying to an issue,
  to whatever tools the reader uses.
* **Self-contained.** Someone who was not in the conversation can follow it,
  with the SOPs it links to, without asking the author what they meant.
* **Decidable.** The Selection rules map what the investigation finds onto an
  option, with no gap and no overlap. Where the evidence will not settle it in
  time, a fallback says what to do.
* **Conservative where it is unsure.** There is always a safe option, such as
  escalating, pausing or asking for approval, for when the evidence is
  ambiguous.
* **Auditable.** Execution says what to record: who acted, when (in UTC), the
  identifiers involved, and a reason code where the team uses them. The
  decision can then be checked afterwards.
* **Closed.** Evaluation says how the team will know the playbook works, with
  measures and targets, a review, and what to do if the targets are missed.
* **Current.** It is kept up to date as the systems and policies it names
  change. An out-of-date playbook misleads the person who trusts it. Routine
  steps kept in SOPs are kept current in one place.

## How the sections read

The Scenario is in the present tense: "A customer contacts support and says
they cannot cancel their subscription." The Problem is a question, or a
statement of the decision: "What decision process should support use to
cancel a subscription safely on the customer's behalf?" The Immediate action
is an instruction, or "No immediate action is indicated." The rest is
written as instructions and rules: "Do not retry the payment", "Collect these
facts before deciding", "Select Option 1 if…", "Roll back the release by
[SOP 3](../sop/SOP-0003-roll-back-a-release.md)". These are followed by the
measures and reviews of Evaluation.
