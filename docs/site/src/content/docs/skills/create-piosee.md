---
title: create-piosee
description: Write a PIOSEE playbook with pse, a plan made in advance for a situation your team expects to face, from research and a short interview.
section: skills
order: 9
---

`create-piosee` writes a PIOSEE playbook with [`pse`](../../pse/). A playbook
is a plan, made in advance, for a situation the team expects to face, such as
a customer reporting an unexpected charge, a tax threshold being crossed, or
a failure that has happened before. It sets out the Scenario, the Problem to
decide, the Information to gather, the Options, the rules for Selection, the
steps of Execution, and how to Evaluate the result.

**Needs:** `pse` on the `PATH`. Reading a third party's documentation, such
as a payment provider's, needs a web fetch or search tool.

## Ask for it

* "Write a playbook for when a customer says they've been charged twice."
* "Plan what we do when a Stripe payment is handled incorrectly."
* "Review the double-charge playbook now that we've used it."

## How it works

1. Interviews you about the problem: the situation and the decision it calls
   for.
2. Researches the rest: the code that handles the situation, where the
   evidence lives, and the documentation of the systems involved.
3. Proposes the investigation, the options and the steps, for you to
   confirm.
4. Interviews you step by step only where its research found nothing.
5. Shows you a draft for approval, then creates it with `pse new`.

Policy is the team's: authority limits, owners, deadlines and targets come
from you, never from the skill.

## What you get

A new playbook with the status Draft. After a playbook has been used, the
skill can revise it against its own Evaluation section.

## Source

The skill is [`skills/create-piosee/SKILL.md`](https://github.com/mustardcustarddev/detoo/blob/main/skills/create-piosee/SKILL.md), with its references beside it.
