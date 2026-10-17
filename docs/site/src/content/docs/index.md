---
title: Introduction
description: detoo is a set of command-line tools and agent skills for keeping decision records, requests for comments, user stories and PIOSEE playbooks.
section: getting-started
order: 1
---

detoo is a set of command-line tools and agent skills for working with
decision records. Each tool keeps one kind of document as a numbered log of
Markdown files in your repository, next to the code it is about.

| Tool | Keeps | Where, by default |
| --- | --- | --- |
| [`adr`](adr/) | Architecture decision records: decisions, and why they were made. | `docs/adr` |
| [`rfc`](rfc/) | Requests for comments: proposed changes, written down for comment before they are made. | `docs/rfc` |
| [`sty`](sty/) | User stories: who wants something of the software, what, and why, with the acceptance criteria that say when it is done. | `docs/userstory` |
| [`pse`](pse/) | PIOSEE playbooks, each a plan for a situation before it arises. | `docs/piosee` |

A fifth tool, [`detoo`](detoo/), looks after the logs as a whole. `detoo
doctor` checks that each log is where its tool looks for it, and that its
documents are named and numbered as the tool expects.

The four record tools work the same way. `init` starts a log, `new` adds the
next numbered document from a template, `link` records that one document
replaces or updates another, `list` lists them or writes a table of contents,
`graph` draws the links between them, and `config` prints where the tool is
looking.

The [agent skills](skills/) use the tools to write the documents with you.
They read the code and the existing records, ask you for the reasons they
cannot find, and leave every document as a Draft for your team to review.

## Where next

* [Install](install/) the tools.
* Follow the [quick start](quick-start/) to write your first record.
* Look up a command, such as [`adr new`](adr/new/). Each command's page shows
  exactly what the tool prints for `help`.
