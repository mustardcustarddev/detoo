---
title: Quick start
description: Start a decision log, write and supersede a record, and publish a table of contents and a graph, in a few commands.
section: getting-started
order: 3
---

This walks through `adr`, the architecture decision record tool. The other
record tools work the same way.

## Start a log

In the root of your project:

```bash
adr init
```

This creates `docs/adr`, with a first record that explains the decision to
keep ADRs: `docs/adr/ADR-0001-record-architecture-decisions.md`.

## Record a decision

```bash
adr new "Use PostgreSQL for persistence"
```

`new` creates the next numbered record from the template, opens it in your
editor, and prints its path: `docs/adr/ADR-0002-use-postgresql-for-persistence.md`.
The title is one argument, so quote it. A new record has the status Draft:
change it to Accepted once the decision is agreed.

## Change your mind

Decisions are not edited away. A new decision supersedes the old one, and
both records say so:

```bash
adr new --supersede 2 "Use SQLite for persistence"
```

The superseded record loses its Draft or Accepted status, and points to the
record that replaced it. A record can also update another, with `--update`,
which leaves the other's status as it is.

## See the log

```bash
adr list                                # the records' paths, one a line
adr list --toc > docs/adr/README.md     # a table of contents in Markdown
adr graph                               # the links between records, in Mermaid
```

See [Publishing a table of contents and graphs](../guides/publishing/) for
more.

## The other tools

```bash
rfc init && rfc new "Adopt Connect for internal services"
ust init && ust new "Reorder a previous purchase"
plm init && plm new "Customer charged twice"
sop init && sop new "Roll back a release"
fyr init && fyr new "Customer charged twice"
```

`rfc init`, like `adr init`, writes a first document that explains the
practice. `ust init`, `plm init`, `sop init` and `fyr init` write a
`README.md` instead, which says what the log is for and how to add to it, so
the first document you write is number 1.

`rfc new --obsolete`, `ust new --replace`, `plm new --replace`,
`sop new --supersede` and `fyr new --replace` replace an earlier document, as
`adr new --supersede` does.

## Check the logs

```bash
detoo doctor
```

`detoo doctor` checks that every log is where its tool looks for it, and that
its documents are named and numbered as the tool expects. See
[Troubleshooting with detoo doctor](../guides/doctor/).

## Let your agent help

The [agent skills](../skills/) write these documents with you, asking for the
reasons behind each one.
