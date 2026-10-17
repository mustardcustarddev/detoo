---
title: create-adr
description: Record an architecture decision with adr, filling in the template only from evidence, and interviewing you for the rest.
section: skills
order: 3
---

`create-adr` records an architecture decision with [`adr`](../../adr/). It
fills in the log's template only from evidence: what you have said, the code,
and documents in the repository. It interviews you for anything it cannot
evidence, rather than inventing context, alternatives or justifications.

**Needs:** `adr` on the `PATH`.

## Ask for it

* "Write an ADR for moving our persistence to PostgreSQL."
* "Record that we've decided to use UTC timestamps throughout."
* "Supersede ADR 4: we're moving off Redis."
* "Start an ADR log for this project."

## How it works

1. Finds the log, or offers to start one with `adr init`.
2. Checks that the decision is architectural, and so belongs in the log.
3. Gathers the facts from the code and the documents.
4. Interviews you for what it cannot find, such as why you chose this option
   over the others.
5. Shows you a draft, with where each statement came from, for approval.
6. Creates the record with `adr new`, and fills in the template.

An accepted ADR is history, so the skill does not rewrite what it decided. To
change or reverse a decision, it writes a new record that supersedes or
updates the old one.

## What you get

A new record in your ADR log, with the status Draft. Accepting it is for your
team to do.

## Source

The skill is [`skills/create-adr/SKILL.md`](https://github.com/mustardcustarddev/detoo/blob/main/skills/create-adr/SKILL.md), with its references beside it.
