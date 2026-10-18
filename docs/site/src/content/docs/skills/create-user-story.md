---
title: create-user-story
description: Write a user story with ust from a one-line subject, asking who it is for, what they need, why, and how you will know it is done.
section: skills
order: 8
---

`create-user-story` writes user stories with [`ust`](../../ust/). It turns a
short subject into a draft story with acceptance criteria, by reading the
project's story template afresh, gathering facts from the code, and
interviewing you for who the story is for, what they need and why.

**Needs:** `ust` on the `PATH`.

## Ask for it

* "Write a user story for exporting invoices."
* "Split the checkout epic into stories."
* "Add acceptance criteria for password reset."

## How it works

1. Finds the log, and reads its template every time.
2. Checks that the subject is a story. A large or vague subject is an epic,
   which it plans to split.
3. Gathers the facts from the code.
4. Interviews you for who wants it, what, why, and how you will know it is
   done.
5. Shows you a draft for approval, then creates it with `ust new`.

## What you get

A new story, or several from an epic, each with the status Draft. A story
that replaces or updates another is linked to it in both directions.

## Source

The skill is [`skills/create-user-story/SKILL.md`](https://github.com/mustardcustarddev/detoo/blob/main/skills/create-user-story/SKILL.md), with its references beside it.
