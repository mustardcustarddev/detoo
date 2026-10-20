---
title: create-firedrill
description: Write a self-led, Choose Your Own Adventure fire drill from a problem playbook, on a local copy of the system, that causes no real problem.
section: skills
order: 11
---

`create-firedrill` writes a fire drill with `fyr`, from a problem playbook
written with [create-playbook](../create-playbook/). In a fire drill you meet
a simulated incident and respond as the playbook says, so that you are ready
for a real one, and find out where the playbook falls short. See
[Fire Drills](https://openpracticelibrary.com/practice/fire-drills/) in the
Open Practice Library for the practice.

Each drill is self-led: one person runs it alone. It reads like a Choose
Your Own Adventure book for software engineers. The fire drill document sets
out the setup, then what happens, step by step. Where the playbook asks for
a decision, it asks you to choose, and each option leads to a page of what
happens next. Every path ends at the solution: the response the playbook
expected, and a debrief.

The drill runs on a local copy of the system, run the way the project runs
it locally. Nothing is broken to make the incident: the pages describe the
symptoms, and stand in for the evidence. They are based on the real system.
The log lines are the ones the code writes, in its format, and the events
and errors of the other systems involved, such as a payment provider, are
the ones their documentation gives.

A drill is one of three kinds, chosen from the playbook:

* **Reproduced**: you run the real system locally on the reporter's inputs,
  such as a command-line tool on their files, and see what it does.
* **Log-driven**: you search log files the drill holds, from each service,
  with ordinary traffic around the incident, and follow one identifier
  through them.
* **Alert-driven**: an alert starts it, from the project's own rules, and
  the metrics, events and logs arrive as time passes, and change with what
  you do.

Like [create-sop](../create-sop/), it works from research rather than
interview: the playbook gives the scenario, the decisions and the response,
the code gives the setup and the logs, and the documentation of the other
systems gives what they report.

**Needs:** [`fyr`](../../fyr/) on the `PATH` to write the drill, `plm` to
find the playbook, and `sop` to find the SOPs it refers to. Reading a third
party's documentation needs a web fetch or search tool.

## Ask for it

* "Write a fire drill for the double-charge playbook."
* "Prepare a game day for playbook 4."
* "We want to practise what we do when a release breaks something."

## How it works

1. Reads the playbook, and the SOPs it refers to.
2. Chooses the kind of drill, and a case: the facts the evidence will show,
   and the option the playbook's Selection rules lead to.
3. Reads how the project runs the system locally.
4. Reads what the code logs on the problem's path, and the documentation of
   each other system involved, for the version the project uses.
5. Plots the drill: the steps, the choices at each decision the playbook
   makes, and the page each choice leads to.
6. Drafts the drill, and asks you once whether to write it, with up to three
   questions it could not settle.
7. Creates the drill with `fyr new`, fills it in, and writes its pages.

It never causes a real problem, and runs nothing to make the drill. Where a
choice would act outside the local system, such as on a release or a
customer, the drill has you write the step down rather than take it.

## What you get

In `docs/firedrill`, the fire drill document, such as
`FIREDRILL-0003-customer-charged-twice.md`, and a directory beside it with
the same name:

* the document, where you start: the playbook, how to play, the ground
  rules, the setup, and the drill, step by step;
* `problem`: the incident that starts the drill, the logs and other evidence,
  and a page for each choice;
* `solution`: `expected-response.md`, the response the playbook leads to, and
  `debrief.md`, questions to look back with, and what the drill found in the
  playbook.

The document has the status Draft. Read the solution only when the drill is
over. Afterwards, [create-playbook](../create-playbook/) can revise the
playbook from what the drill found.

## Source

The skill is [`skills/create-firedrill/SKILL.md`](https://github.com/mustardcustarddev/detoo/blob/main/skills/create-firedrill/SKILL.md), with its references beside it.
