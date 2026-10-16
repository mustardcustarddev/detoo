# rfc: a second record tool

Date: 2026-09-10

## Goal

Add an `rfc` command that keeps a log of requests for comments the way `adr`
keeps a log of architecture decision records. Reorganise the code so that the
two tools share one engine and differ only in data.

## Behaviour

`rfc` has the same commands as `adr`: `init`, `new`, `link`, `list` (with
`--toc`), `graph` and `config`, plus `help` and `completion`.

| | adr | rfc |
| --- | --- | --- |
| Directory and marker file | `docs/adr`, `.adr-dir` | `docs/rfc`, `.rfc-dir` |
| Status of a new record | Accepted | Draft |
| Relations | Supercedes / Superceded by, Updates / Updated by | Obsoletes / Obsoleted by, Updates / Updated by |
| Flags on `new` and `link` | `-s, --supercede`, `-u, --update` | `-o, --obsolete`, `-u, --update` |
| Environment variables | `ADR_TEMPLATE`, `ADR_DATE`, `ADR_PAGER` | `RFC_TEMPLATE`, `RFC_DATE`, `RFC_PAGER` |
| `config` output | `adr_dir`, `adr_template` | `rfc_dir`, `rfc_template` |
| Table of contents heading | Architecture Decision Records | Requests for Comments |
| Graph node IDs | `adr1` | `rfc1` |
| First record | 1. Record architecture decisions | 1. Use requests for comments |

Obsoleting an RFC removes its `Draft` line, as superceding an ADR removes
`Accepted`. An RFC whose status has been edited to something else keeps that
status beside "Obsoleted by", as a Proposed ADR does when it is superceded.

A new RFC is created from this template, unless the project or `RFC_TEMPLATE`
provides one:

```markdown
# NUMBER. TITLE

Date: DATE

## Status

STATUS

## Summary

The proposal, in a paragraph.

## Motivation

Why the change should be made: the problem it solves, and who it matters to.

## Design

The change, in enough detail to comment on and to build.

## Drawbacks

Why the change might not be made.

## Alternatives

What else was considered, and why this design was chosen over it.

## Unresolved questions

What should be settled before the RFC is accepted, and what can wait.
```

`rfc init` writes `0001-use-requests-for-comments.md`, an accepted RFC in the
same sections, which records the decision to propose significant changes as
RFCs.

## Design

### Tool definitions

`internal/adr` and `internal/rfc` each export:

* `Kind`, a `record.Kind` holding the tool's vocabulary, and
* `Files`, an embedded file system holding its `templates/` and `help/`.

`Kind.Name` names everything that belongs to a tool, by one convention applied
where each thing is used: the command, the `NAME_TEMPLATE`, `NAME_DATE` and
`NAME_PAGER` variables, the `.NAME-dir` marker, the `docs/NAME` directory, the
`NAME_dir` and `NAME_template` lines of `config`, and the `NAMEn` graph nodes.

### Shared engine

* `internal/domain/adr` becomes `internal/domain/record`, which also defines
  `Kind`. `Relation` becomes `Replaces | Updates`, and the words for each come
  from the Kind.
* The slices keep their shape and their ports. Those that write the tool's
  words (create, initialise, link, toc and graph) are given a Kind.
* The store, templates and clock adapters take the Kind, and derive their
  names from it.
* `app.New` and `cli.NewRootCommand` take the Kind, and the CLI takes the
  tool's `Files` for its help.
* `cmd/adr/main.go` and `cmd/rfc/main.go` each wire their own tool. The
  Makefile already builds everything under `cmd/`.

### Tests

* The refactor comes first and must not change adr. Its golden files move to
  `test/testdata/adr/`. The only content changes are two flag descriptions in
  `adr help graph` and `adr help list`, reworded to fit both tools.
* The slice tests gain cases using RFC wording, to show the words come from
  the Kind.
* rfc's end-to-end tests, in `test/rfc_test.go` with golden files in
  `test/testdata/rfc/`, cover what differs: init, the Draft status, obsoleting
  and updating, the contents heading, the graph, `config`, the `RFC_*`
  variables, help and completion. One test checks that adr and rfc logs in the
  same project do not see each other. Behaviour the tools share (numbering,
  slugs, directory search, the editor) is not tested again through rfc.

## Out of scope

* Starting an RFC log in this repository.
* Commands that change a record's status, such as accepting an RFC.
