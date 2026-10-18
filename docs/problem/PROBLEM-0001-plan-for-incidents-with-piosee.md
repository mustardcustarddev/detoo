# 1. Plan for incidents with PIOSEE

Date: 2026-10-16

## Status

Draft

## Scenario

We foresee a situation that will call for a decision under pressure, such as
a customer reporting an unexpected charge, a deadline set by a regulator, or a
failure we have seen before, and we have no plan for it.

## Problem

Should we write a playbook for this situation now, before it arises, and if
so, how?

Primary risks:

* Deciding in a hurry when it arises, on part of the facts, by whoever is
  there, and differently each time.
* Spending effort on plans for situations that never arise, or that are
  better handled as they come.

## Immediate action

No immediate action is indicated.

## Counter-indicated actions

* Writing the playbook while the situation is being handled. Deal with it
  first, then write the playbook from what was learned.

## Information

Collect these before deciding:

* How often the situation arises, or is likely to, and when it last did.
* What it costs to decide it wrongly, or late: to customers, to money, to
  compliance.
* Whether a playbook, standard operating procedure (SOP), runbook or policy
  already covers it, or a situation close to it, in this directory or
  elsewhere.
* Who will have to act when it arises, and what access and authority they
  will need.

## Options

1. **Decide when it arises.** Suits a situation that is rare and cheap to get
   wrong. Costs nothing now, but the decision is made under pressure.
2. **Write a playbook.** Suits a situation that recurs, or is costly to get
   wrong. Takes effort now, and must be kept up to date.
3. **Extend an existing playbook.** Suits a situation close to one that a
   playbook already covers. Keeps related plans together, but can make that
   playbook harder to follow.

## Selection

* Extend an existing playbook (option 3) when one covers a situation close
  to this one, and the same decision, options and selection rules apply.
* Otherwise, write a playbook (option 2) when the situation recurs, or a
  wrong or late decision would be costly. One known option is enough, with a
  conservative fallback for when it does not apply.
* Otherwise, decide when it arises (option 1).

## Execution

1. Create the playbook in this directory with `plm new`. Each playbook is a
   Markdown file, numbered in sequence. Numbers are not reused.
2. Fill in each section: the scenario it is for, the problem to decide, the
   action to take at once, if any, and the actions to avoid, the information
   to gather, the options, the rules for selecting one, the steps to execute
   it, and how to evaluate the result.
3. Where a step is a routine procedure, such as rolling back a release,
   refer to the standard operating procedure (SOP) for it, kept with
   `sop new`, rather than repeating its steps. Write the SOP first if there
   is none.
4. Have the people who will act on it review it. A new playbook is a draft
   until they agree it.
5. When the situation arises, follow the playbook, and record what was
   decided and done where the team records its incidents.
6. A playbook rewritten for the same situation replaces the one it rewrites,
   and both say so. A playbook that changes part of another, such as one that
   extends it to a situation close to it, updates it, and both say so.

## Evaluation

* Review each playbook after it is used: whether it led to a sound decision
  in time, and what it lacked.
* Review the directory regularly: which playbooks were used, which are out
  of date, and which situations still lack one.
* Rewrite or retire a playbook that did not help.
