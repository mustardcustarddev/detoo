# PIOSEE playbooks

## What a playbook is

A PIOSEE playbook is a plan, written in advance, for deciding and acting when
a foreseen situation arises. It is named for its steps: **P**roblem,
**I**nformation, **O**ptions, **S**election, **E**xecution and
**E**valuation. They are prefaced by the **Scenario** the playbook is for.

The situation may be one that recurs, such as a customer who cannot cancel a
subscription, or a charge reported after an account was deleted. It may be
one that is scheduled or foreseen, such as crossing a tax registration
threshold on a known date. It may be one that has happened once and could
happen again. What these have in common is that someone will one day have to
decide what to do, often under pressure. The playbook means the thinking has
already been done, calmly, by the people who know the system and the policy.

When the situation arises, the person on hand does the following:

1. Recognises the scenario.
2. Gathers the information it lists.
3. Applies its selection rules to choose an option.
4. Follows the execution steps.
5. Later, judges the result by its evaluation measures.

They decide in the moment, but by rules agreed beforehand. So similar cases
are decided alike, and each step leaves an audit trail.

## What a playbook is not

* **Not a record of an incident.** It does not say what happened on a given
  day, who decided what, or how it turned out. A postmortem or an incident
  ticket does that. A playbook is written for every occurrence of a situation.
* **Not a decision already taken.** An ADR records a decision about how the
  system is built, made once. A playbook's Selection is a rule for deciding
  each time. It may name a default, but the choice is made when the scenario
  arises, against the evidence gathered then.
* **Not only a list of steps.** A runbook says what to do. A playbook also
  says what to find out first, what the options are, how to choose between
  them, and how to tell whether it worked.

## The sections

The order follows how the situation will be handled. Each section builds on
the one before.

| Section         | What it holds                                                                                                                                              |
| --------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Scenario**    | The triggering situation, in one to three sentences, with any standing conditions that frame it.                                                           |
| **Problem**     | The exact decision to make, as one clear question, and the primary risks if it is decided wrongly. No solutions.                                           |
| **Information** | The facts to collect before deciding: what, where each is found, and who owns it. Known facts apart from assumptions. Critical unknowns, and how to settle them quickly. |
| **Options**     | The realistic courses of action, distinct and actionable, with at least one conservative fallback. For each: when it applies, and one or two tradeoffs.  |
| **Selection**   | How to choose: explicit criteria, which findings from Information lead to which option, the default, the confidence in the rule, and its key assumptions. |
| **Execution**   | The steps, in order, with owners and dependencies. What to record and communicate. The rollback path. Steps to prevent a recurrence.                     |
| **Evaluation**  | How success is measured once the playbook has been used: measurable outcomes, when and by whom they are reviewed, and follow-up actions if targets are missed. |

A playbook may add sections that help the person following it, such as
**Communication Templates** with the message to send for each outcome, or
**Reference Links** to the documentation of the systems involved.

## What makes a good one

* **Concrete.** It names the actual log source, storage path, record field,
  dashboard page and menu path. An Information section that says "check the
  logs" is no help at two in the morning. One that names the service whose
  logs to read, and the events to look for, is.
* **Self-contained.** Someone who was not in the conversation can follow it,
  without asking the author what they meant.
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
  change. An out-of-date playbook misleads the person who trusts it.

## How the sections read

The Scenario is in the present tense: "A customer contacts support and says
they cannot cancel their subscription." The Problem is a question, or a
statement of the decision: "What decision process should support use to
cancel a subscription safely on the customer's behalf?" The rest is written
as instructions and rules: "Collect these facts before deciding", "Select
Option 1 if…", "Verify the customer's identity". These are followed by the
measures and reviews of Evaluation.
