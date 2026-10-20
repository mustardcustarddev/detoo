---
title: "Senior to staff engineer: don't be your team's single point of failure"
pubDate: 2026-10-18
author: "tom-regan"
tags: [ "staff engineering", "incident management", "decision-making", "SRE" ]
categories: [ "Technology" ]
tool: plm
description: "Pilots use PIOSEE, a six-step process anyone can follow under pressure. Put your incident judgement into PIOSEE playbooks and the whole rota can use it."
heroImage: ../../assets/senior-to-staff-single-point-of-failure.jpg
---

Checkout starts failing on a Tuesday evening. The engineer on call acknowledges the page, opens three dashboards and
posts in the incident channel. Two more engineers join and start looking at logs. Then everyone waits for one person.

When the senior engineer arrives, they read the graphs for a minute, say the connection pool is exhausted, tell someone
to roll back the afternoon's release and someone else to warn the payments team. Twenty minutes later checkout is
back. Everyone remembers who fixed it.

On a lot of teams, that's how you get promoted to senior. I think it's also what stops you getting to staff. The team's
incident response works as well as that one engineer does, for as long as they're awake, reachable and still working
there. Anywhere else in the system, we'd call that a single point of failure.

## When the checklist runs out

An airliner loses an engine just after take-off. The crew has a checklist for that, and they follow it: keep the
aircraft flying, secure the dead engine, climb away. The checklist deals with the engine. It doesn't tell them where to
land on the one that's left, with a longer landing distance and perhaps not enough fuel for the airport they'd planned
on.

For that, pilots have decision-making models, and one of them is PIOSEE: Problem, Information, Options, Select, Execute,
Evaluate. Define the problem you actually have to solve, gather what bears on it, lay out what you could do, choose
one, carry it out, and check whether it's working. If something has changed, go back to the start.

PIOSEE is a structure, so it works whoever is in the seat. Pressure narrows people's attention, and a trained crew
working through six steps in order will catch things a brilliant but frightened one misses. Airlines teach a process
like it to every pilot, because they can't count on each cockpit having someone gifted at emergencies in it.

## What the senior engineer is actually doing

Watch the senior engineer on that incident call and you'll see PIOSEE, done silently. While everyone else stared at the
error count, they worked out that the problem was customers unable to pay. They knew which dashboard to trust. They had
four possible fixes in mind and knew which one could be undone. After the rollback, they checked the success rate to
see whether it had worked.

None of that is visible to the rest of the team, so none of it is passed on. The engineer on call learns that the
answer is to wait for the senior engineer.

Tanya Reilly, in *The Staff Engineer's Path* (2022), makes levelling up the engineers around you one of the three
pillars of the staff role. For incident response, I think that means writing down the judgement and teaching the team
the process.

## Write the judgement down

Most teams already have runbooks: how to roll back a release, how to fail over the database, how to issue a refund.
They're checklists, and they're useful. But each one assumes someone has already decided it's the right thing to do,
and the deciding is what the senior engineer was carrying in their head.

A PIOSEE playbook is where the deciding goes. It's written for one problem, before it happens, and it links to the
runbooks for the routine steps. Take a customer who says they've been charged twice:

- **Scenario.** A customer reports two charges for one order.
- **Problem.** Has this customer actually paid twice, and is it still happening to other people?
- **Immediate action.** Like a pilot's memory items, anything safe to do before the cause is known. For one report,
  none. If the duplicate-capture alert is firing too, pause the retry job.
- **Counter-indicated actions.** What makes it worse. Don't refund before checking which charge was captured, because
  the other may be an authorisation that will lapse, leaving the customer with nothing paid at all.
- **Information.** Where to look: the provider's transactions, the internal ledger, the idempotency keys, recent
  deployments.
- **Options.** Wait for the authorisation to drop off, cancel it, refund a confirmed duplicate, pause the retry job, or
  roll back. For each, what it costs, whether it can be undone and who's allowed to do it.
- **Selection.** The rules for choosing, agreed beforehand with whoever owns the trade-offs. One confirmed duplicate
  gets a refund. Several in an hour means pausing the retry job before anyone refunds anything.
- **Execution.** Who does what, linking to the runbook for each step.
- **Evaluation.** What should change and by when, and what to do if it hasn't.

With that open, the engineer on call can make the first decisions before anyone senior has joined.

## Teach the process on the call

When the next incident comes, join the call and ask the questions. What problem are we solving? What do we know, and
what don't we? What could we do? Which of those can we undo? How will we know it's worked, and by when? Leave the
answers to the engineer on call, with the playbook open.

It's slower the first few times, and you'll often see the answer before anyone else does. Unless the harm is growing
faster than the team can work, keep it to yourself and note it down. Anything you wanted to say that wasn't in the
playbook belongs in it.

## How I'd start

[detoo](../../) includes [`plm`](../../docs/plm/), a command-line tool that keeps PIOSEE playbooks as Markdown files in
your repository, and an agent skill, [create-playbook](../../docs/skills/create-playbook/), that writes them with you.

1. Look back over your last few incidents and count the ones you joined without being on call. What you were asked on
   those calls is your list of playbooks.
2. [Install `plm` and the skill](../../docs/skills/install/), and ask your agent for the first playbook in a sentence,
   such as "Write a playbook for when checkout starts failing after a release."
3. Answer its questions about the problem, with the engineer who'll follow the playbook sitting next to you. It reads
   your code, your runbooks and the documentation of the systems involved, then proposes the rest for you to confirm.
4. Agree the selection rules with the team. The skill won't invent them, and every playbook it creates is a draft until
   the team agrees it.
5. Practise it with someone junior on call and you asking the questions, and see where it falls apart.
6. At the next real incident, ask the questions and let the engineer on call answer them.

You wouldn't approve a design with a single point of failure in it. Don't be one.

---

### Further reading

- Petter Hörnfeldt (Mentour Pilot), "HOW to deal with CRISIS.", YouTube (2020).
- Tanya Reilly, *The Staff Engineer's Path: A Guide for Individual Contributors Navigating Growth and Change* (2022).
