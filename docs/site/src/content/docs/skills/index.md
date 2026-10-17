---
title: Agent skills
description: Agent skills that use the detoo tools to write ADRs, RFCs, user stories and PIOSEE playbooks with you, and to review and check an ADR log.
section: skills
order: 1
---

The skills are [Agent Skills](https://agentskills.io/specification) for
coding agents such as Claude Code and Codex. They drive the detoo tools to
write ADRs, RFCs, user stories and PIOSEE playbooks with you, asking you for
the reasons behind each one. Others review an ADR log, and check the code
against it.

Let your agent write the first draft. The skills read the code, the existing
records and the documentation of the systems involved, ask you about what
they cannot find, and fill in your project's template.

## The skills

| Skill | What it does | Needs |
| --- | --- | --- |
| [create-adr](create-adr/) | Records a decision, filling in the template only from evidence, and interviewing you for the rest. | `adr` |
| [review-adrs](review-adrs/) | Judges whether each record is architectural, shows the reasoning, and offers to deprecate, supersede or delete those that are not. | `adr` |
| [find-missing-adrs](find-missing-adrs/) | Reads the whole codebase once, and proposes records for the architectural decisions in the code that the log does not record. | `adr` |
| [check-adrs](check-adrs/) | Checks that the code follows one record, or every record in force, and suggests fixes where it does not. | `adr` |
| [create-rfc](create-rfc/) | Drafts an RFC from a one-line request, reading the code for the facts, and asking you what you intend and why. | `rfc` |
| [create-user-story](create-user-story/) | Writes a user story from a one-line subject, asking who it is for, what they need, why, and how you will know it is done. | `sty` |
| [create-piosee](create-piosee/) | Writes a PIOSEE playbook: a plan, made in advance, for a situation your team expects to face. | `pse` |

## Ask in your own words

You don't name the skill. Ask your agent for what you want, and it picks the
skill:

| Prompt | Skill |
| --- | --- |
| "Write an ADR for moving our persistence to PostgreSQL." | `create-adr` |
| "Draft an RFC to migrate the system to hexagonal architecture." | `create-rfc` |
| "Write a user story for exporting invoices." | `create-user-story` |
| "Write a playbook for when a customer says they've been charged twice." | `create-piosee` |
| "Review our ADRs, and tell me which aren't really architectural." | `review-adrs` |
| "Find the architecture decisions in this codebase that have no ADR." | `find-missing-adrs` |
| "Check that the code follows ADR 7." | `check-adrs` |

## What they have in common

Each skill writes only what it has a source for: the repository, the
documentation it has read, or you. It asks you for anything else. Each also
keeps its log to what belongs there. Only architectural decisions go in the
ADR log, and only changes that others would want a say in need an RFC. A user
story describes an outcome someone wants. A PIOSEE playbook plans for a
situation before it arises, by rules the team sets.

Every ADR, RFC, user story and PIOSEE playbook an agent creates has the
status Draft. Accepting one is for people to do.

Next, [install the skills](install/).
