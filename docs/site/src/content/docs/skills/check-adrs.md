---
title: check-adrs
description: Check that the code follows one ADR or every ADR in force, rule by rule, with evidence, and advise on fixing the code or changing the decision.
section: skills
order: 6
---

`check-adrs` checks that the code follows one ADR, or the whole log. It turns
each decision into concrete rules, searches the code for evidence of
compliance and of violation, and gives a verdict on each rule with the
evidence it found.

**Needs:** `adr` on the `PATH`, and read access to the whole repository.

## Ask for it

* "Check that the code follows ADR 7."
* "Are we still following our architecture decisions?"
* "Find any architecture drift before I rely on the ADRs."

## How it works

Each rule gets a verdict:

| Verdict | Meaning |
| --- | --- |
| Adheres | Searched thoroughly, and no violation found. |
| Violates | The code breaks the rule, with every instance listed. |
| Partial | Followed in some places and broken in others, as in an unfinished migration. |
| Undetermined | The evidence is inconclusive, with what would settle it. |
| Not checkable | The rule is not about code. |
| Obsolete | The record concerns code that no longer exists. |

A record's overall verdict is the worst of its rules, leaving out any that
are not checkable.

## What you get

A report, and for every rule that is broken, two ways forward: bring the code
into line, or record a new decision that supersedes the old one, with
[`create-adr`](../create-adr/). The record itself is never edited to match
the code.

## Source

The skill is [`skills/check-adrs/SKILL.md`](https://github.com/mustardcustarddev/detoo/blob/main/skills/check-adrs/SKILL.md), with its references beside it.
