---
title: create-sop
description: Write a standard operating procedure with sop - the steps of a routine task - researched from the code and the tools' documentation, with few questions.
section: skills
order: 10
---

`create-sop` writes a standard operating procedure (SOP) with
[`sop`](../../sop/): the steps of a routine task, such as rolling back a
release, rotating a key or restoring a backup, as an ordered list. An SOP has
no decision in it, so it is done the same way each time. Problem playbooks,
written with [create-playbook](../create-playbook/), refer to SOPs for their
routine steps.

Unlike the other skills, it works from research rather than interview. How a
project does a routine task is usually in its Makefile, its scripts, its CI
workflows and its deployment configuration, and the rest is in the
documentation of the tools involved. So the skill reads those, drafts the
steps with their sources, and asks you only what it cannot find.

**Needs:** `sop` on the `PATH`. Reading a tool's documentation needs a web
fetch or search tool.

## Ask for it

* "Write an SOP for rolling back a release."
* "Document how we rotate the webhook signing secret."
* "Update the backup restore SOP now that we've moved to the new cluster."

## How it works

1. Reads the scripts, build files, CI workflows, deployment configuration,
   runbooks and documents that touch the task.
2. Reads the official documentation of each tool involved, for the version
   the project uses.
3. Drafts the steps, each with its source, and leaves any gap as an open
   question in the draft rather than guessing.
4. Asks you once: whether to create the draft, and up to three questions it
   could not settle, such as which of two sources holds.
5. Creates it with `sop new`.

It never runs the procedure to test it: it reads the steps, it does not
carry them out.

## What you get

A new SOP with the status Draft, with any open questions marked for the
team. Have someone who does the task follow it once, and correct it.

## Source

The skill is [`skills/create-sop/SKILL.md`](https://github.com/mustardcustarddev/detoo/blob/main/skills/create-sop/SKILL.md), with its references beside it.
