---
title: create-rfc
description: Draft a request for comments with rfc from a one-line request, reading the code for the facts and asking you what you intend and why.
section: skills
order: 7
---

`create-rfc` drafts a request for comments with [`rfc`](../../rfc/). It turns
a short request into a draft others can comment on, by reading the project's
RFC template afresh, gathering facts from the code and documents, and
interviewing you section by section for your intent and reasoning.

**Needs:** `rfc` on the `PATH`.

## Ask for it

* "Draft an RFC to migrate the system to hexagonal architecture."
* "Write an RFC proposing we adopt Connect for internal services."
* "Put the move to monorepo out for comment."

## How it works

1. Finds the log, and reads its template every time.
2. Checks that the change needs an RFC. A small, local change may not, and a
   change already decided is a decision to record with
   [`create-adr`](../create-adr/) instead.
3. Gathers the facts from the code and the documents.
4. Interviews you, section by section, for what only you know.
5. Shows you a draft for approval, then creates it with `rfc new`.

## What you get

A new RFC with the status Draft, ready for others to comment on.

## Source

The skill is [`skills/create-rfc/SKILL.md`](https://github.com/mustardcustarddev/detoo/blob/main/skills/create-rfc/SKILL.md), with its references beside it.
