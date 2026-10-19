---
title: create-firedrill
description: Write a fire drill from a problem playbook - an exercise on a local copy of the system, with realistic assets, that causes no real problem.
section: skills
order: 11
---

`create-firedrill` writes a fire drill from a problem playbook, written with
[create-playbook](../create-playbook/). In a fire drill the team meets a
simulated incident and responds as the playbook says, so that it is ready
for a real one, and finds out where the playbook falls short. See
[Fire Drills](https://openpracticelibrary.com/practice/fire-drills/) in the
Open Practice Library for the practice.

The drill runs on a local copy of the system, run the way the project runs
it locally. Nothing is broken to make the incident: the fire drill document
describes the symptoms, and its assets stand in for the evidence. Those
assets are based on the real system. The log lines are the ones the code
writes, in its format, and the events and errors of the other systems
involved, such as a payment provider, are the ones their documentation
gives.

Like [create-sop](../create-sop/), it works from research rather than
interview: the playbook gives the scenario and the response, the code gives
the setup and the logs, and the documentation of the other systems gives
what they report.

**Needs:** `plm` on the `PATH` to find the playbook, and `sop` to find the
SOPs it refers to. Reading a third party's documentation needs a web fetch or
search tool.

## Ask for it

* "Write a fire drill for the double-charge playbook."
* "Prepare a game day for playbook 4."
* "We want to practise what we do when a release breaks something."

## How it works

1. Reads the playbook, and the SOPs it refers to.
2. Chooses a case: the facts the evidence will show, and the option the
   playbook's Selection rules lead to.
3. Reads how the project runs the system locally.
4. Reads what the code logs on the problem's path, and the documentation of
   each other system involved, for the version the project uses.
5. Drafts the fire drill document and its assets, and asks you once whether
   to write them, with up to three questions it could not settle.
6. Writes them to a `firedrill` directory beside the problem log.

It never causes a real problem, and runs nothing to make the drill. Where
the playbook's response would act outside the local system, such as on a
release or a customer, the drill has the team describe the step rather than
take it.

## What you get

A directory with:

* `FIREDRILL.md`, for whoever runs the drill: the setup, the ground rules,
  the symptoms, the expected response, a debrief, and the sources;
* `assets/`, for the team: the incident description that starts the drill,
  the logs, and any other evidence the playbook says to check.

The document has the status Draft, and holds the expected response, so keep
it from the team until the debrief. Afterwards,
[create-playbook](../create-playbook/) can revise the playbook from what the
drill found.

## Source

The skill is [`skills/create-firedrill/SKILL.md`](https://github.com/mustardcustarddev/detoo/blob/main/skills/create-firedrill/SKILL.md), with its references beside it.
