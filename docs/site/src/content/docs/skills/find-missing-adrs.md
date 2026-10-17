---
title: find-missing-adrs
description: Survey a whole codebase for the architecture decisions it embodies that have no ADR, and propose records for them, backed by evidence.
section: skills
order: 5
---

`find-missing-adrs` surveys a whole codebase, thoroughly and in one pass, for
architecture decisions that are built into the code but have no ADR. It
proposes new records based on the code as it really is.

**Needs:** `adr` on the `PATH`, and read access to the whole repository.

## Ask for it

* "Find the architecture decisions in this codebase that have no ADR."
* "We've never kept ADRs. Which decisions should we record first?"
* "What architectural decisions are undocumented?"

## How it works

1. Reads what the log already records.
2. Maps the repository, and surveys every area of it.
3. Judges each candidate: only genuinely architectural decisions are
   proposed, not local design choices or trivia.
4. Reports each proposal with the evidence for it, separating what the code
   shows from the reasons only you can supply.

## What you get

A list of proposed records. You choose which to write, and for each, the
skill hands over to [`create-adr`](../create-adr/), which asks you for the
reasons. Every record starts as a Draft.

## Source

The skill is [`skills/find-missing-adrs/SKILL.md`](https://github.com/mustardcustarddev/detoo/blob/main/skills/find-missing-adrs/SKILL.md), with its references beside it.
