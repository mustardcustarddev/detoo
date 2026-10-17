---
title: review-adrs
description: Judge each record in an ADR log for whether it is really architectural, show the reasoning, and offer to deprecate, supersede or delete the rest.
section: skills
order: 4
---

`review-adrs` reviews an existing ADR log for value. It judges each record by
whether the decision was made at the right level, when it was needed, and
shows its reasoning for every record.

**Needs:** `adr` on the `PATH`.

## Ask for it

* "Review our ADRs, and tell me which aren't really architectural."
* "Which of our decision records are worth keeping?"
* "Tidy up the ADR log."

## How it works

It reads the whole log, then gives each record one verdict:

| Verdict | Meaning |
| --- | --- |
| Architectural | Belongs in the log. |
| Borderline | The tests disagree, and it says which point which way. |
| Design | A real decision, but local: it could be refactored away without architectural change. |
| Trivial | No real alternatives, a convention, or taste. |
| Premature | Decided before it was needed, about something that does not exist yet. |

It also notes records that decide nothing, hold several decisions, duplicate
or contradict another, or no longer match the code.

## What you get

A report, then a question for each record it would not keep as it is: keep
it, deprecate it, supersede it with `adr link`, or delete it. It changes
nothing until you choose.

## Source

The skill is [`skills/review-adrs/SKILL.md`](https://github.com/mustardcustarddev/detoo/blob/main/skills/review-adrs/SKILL.md), with its references beside it.
