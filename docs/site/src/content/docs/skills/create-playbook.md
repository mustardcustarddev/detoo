---
title: create-playbook
description: Write a problem playbook with plm, by the PIOSEE method - a plan made in advance for responding to a problem such as an incident.
section: skills
order: 9
---

`create-playbook` writes a problem playbook with [`plm`](../../plm/). A
playbook is a plan, made in advance, for responding to a problem the team
expects to face, such as an incident, an outage, a customer reporting an
unexpected charge, or a failure that has happened before.

It follows PIOSEE, a decision-making method from aviation. It sets out the
Scenario, the Problem to decide, the Immediate action to take before anything
else, or that none is indicated, the Counter-indicated actions that would
make things worse, the Information to gather, the Options, the rules for
Selection, the steps of Execution, and how to Evaluate the result. PIOSEE is
a way of acting, so the playbook is written to be followed: the first thing
it tells the reader to do, they can do at once.

Where a playbook needs a routine task done, such as rolling back a release,
it refers to a standard operating procedure (SOP), kept with
[`sop`](../../sop/), rather than repeating the steps. If there is no SOP for
it, the skill offers to draft one with [create-sop](../create-sop/).

**Needs:** `plm` on the `PATH`, and `sop` to find the project's SOPs.
Reading a third party's documentation, such as a payment provider's, needs a
web fetch or search tool.

## Ask for it

* "Write a playbook for when a customer says they've been charged twice."
* "Plan what we do when checkout errors spike after a release."
* "Review the double-charge playbook now that we've used it."

## How it works

1. Interviews you about the problem: the situation, the decision it calls
   for, and whether it gets worse while it waits.
2. Researches the rest: the code that handles the situation, what stops the
   harm growing and what makes it worse, where the evidence lives, the SOPs,
   and the documentation of the systems involved.
3. Proposes the immediate action, the investigation, the options and the
   steps, for you to confirm.
4. Interviews you step by step only where its research found nothing.
5. Refers to an SOP for each routine step, or offers to draft one.
6. Shows you a draft for approval, then creates it with `plm new`.

Policy is the team's: authority limits, owners, deadlines and targets come
from you, never from the skill.

## What you get

A new playbook with the status Draft. After a playbook has been used, the
skill can revise it against its own Evaluation section. It also brings a
playbook written for `pse`, the tool `plm` replaced, up to date with an
Immediate action section.

## Source

The skill is [`skills/create-playbook/SKILL.md`](https://github.com/mustardcustarddev/detoo/blob/main/skills/create-playbook/SKILL.md), with its references beside it.
