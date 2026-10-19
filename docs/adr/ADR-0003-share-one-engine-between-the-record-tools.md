# 3. Share one engine between the record tools

Date: 2026-09-10

## Status

Accepted

## Context

This project is meant to provide several command-line tools, each keeping a
log of numbered Markdown records. The first, `adr`, keeps architecture decision
records. The second, `rfc`, keeps requests for comments.

The two tools do the same things. Each starts a log, creates a record from a
template, links one record to another in both records, lists the log, and draws
it as a graph. They differ in words and in files, not in behaviour. An RFC is
created as a draft rather than accepted, and obsoletes a record rather than
superceding it. It lives in `docs/rfc` rather than `docs/adr`, and is
configured by `RFC_` variables rather than `ADR_` ones. Each tool also has its
own templates and help.

We considered three ways to organise the code:

* **One engine, two tools.** The domain, slices, adapters and command tree are
  shared, and each is given a definition of the tool it serves.
* **A shared domain, with a stack for each tool.** The domain and adapters are
  shared, and each tool has its own slices, application and command tree. This
  duplicates about 1,300 lines.
* **A full copy.** Each tool is a complete stack of its own, duplicating about
  2,100 lines.

## Decision

We will share one engine between the record tools.

The domain package `record` defines `Kind`, the vocabulary of one kind of log:
the tool's name, what its records are called, the status a new record is given,
and the words for replacing and updating a record. Each tool is a package,
`internal/adr` or `internal/rfc`, that holds its Kind, embeds its templates and
help, and holds nothing else. The slices, adapters and commands are written once
and take a Kind. Each tool has its own `cmd` package, which puts its definition
and the engine together.

The tool's name names everything else that belongs to it, by one convention:
the `NAME_DATE`, `NAME_TEMPLATE` and `NAME_PAGER` variables, and the
`docs/NAME` directory.

## Consequences

A new record tool is a Kind, its templates and help, and a `main` package of
about fifty lines. A fix to the engine fixes every tool.

The tools cannot diverge without the engine knowing. A command that only one
tool has, or a relation that only one kind of log uses, has to become an option
in the engine, rather than code in one tool. If the tools grow apart, that cost
is the signal to revisit this decision.

The end-to-end tests of `adr` cover the behaviour the tools share. The tests of
`rfc` cover only what its definition changes.
