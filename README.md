# detoo

<img src="https://github.com/user-attachments/assets/43a08c20-05b9-4365-a7d7-ff14c3360a1d" alt="Detoo robot drawing with a set square" width="320">

[![Apache License 2.0](https://img.shields.io/:license-Apache%20License%202.0-blue.svg?style=shield)](LICENSE.md)
[![Build](https://github.com/mustardcustarddev/detoo/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/mustardcustarddev/detoo/actions/workflows/build.yml)
[![Documentation](https://img.shields.io/badge/docs-GitHub%20Pages-ea580c.svg)](https://mustardcustarddev.github.io/detoo/)

Command-line tools and agentic skills for working with decision records:

* `adr` keeps [architecture decision records][ADRs]: decisions, and why they
  were made.
* `rfc` keeps requests for comments: proposed changes, written down for
  comment before they are made.
* `ust` keeps [user stories][stories]: who wants something of the software,
  what, and why, with the acceptance criteria that say when it is done.
* `plm` keeps problem playbooks, each a plan, made in advance, for responding
  to a problem such as an incident. It follows the [PIOSEE][piosee] method:
  the Problem to decide, the Immediate action to take, the Counter-indicated
  actions to avoid, the Information to gather, the Options, how to make the
  Selection, the steps of Execution, and how to Evaluate the result.
* `sop` keeps standard operating procedures: the steps of a routine task,
  such as rolling back a release, which playbooks refer to.

## Install

On an Apple silicon Mac, install the tools with [Homebrew](https://brew.sh).
Add the Mustard Custard tap, once, then install detoo from it:

```bash
brew tap mustardcustarddev/tap
brew install detoo
```

This installs `adr`, `rfc`, `ust`, `plm`, `sop` and `detoo`, with shell
completion for bash, zsh and fish. On other machines, build them from source,
as [For developers](#for-developers) describes.

## Usage

Each tool works the same way. Start a log, which by default lives in
`docs/adr`, `docs/rfc`, `docs/userstory`, `docs/problem` or `docs/sop`, then
add to it:

```bash
adr init
adr new "Use PostgreSQL for persistence"
```

`new` creates the next numbered document from a template, opens it in
`$VISUAL` or `$EDITOR`, and prints its path. Its file name is the log's prefix,
the number and a slug of the title, such as
`ADR-0002-use-postgresql-for-persistence.md`,
`RFC-0004-adopt-connect-for-internal-services.md`,
`STORY-0012-reorder-a-previous-purchase.md`,
`PROBLEM-0004-unexpected-charge-after-deletion.md` or
`SOP-0003-roll-back-a-release.md`. A project can provide its
own template in `templates/template.md` inside the log directory.

A new document can replace or update earlier ones, and the link is recorded
in both:

```bash
adr new --supersede 9 "Use Rust for performance-critical functionality"
rfc new --obsolete 4 "Adopt Connect for internal services"
ust new --replace 12 "Reorder a previous purchase"
plm new --replace 4 "Customer cannot cancel a subscription"
sop new --supersede 3 "Roll back a release"
```

Publish a table of contents, or draw the links between documents with Mermaid
or Graphviz:

```bash
adr list --toc > docs/adr/README.md
ust graph | mmdc --input - --output docs/userstory/graph.svg
```

Run `adr help`, or `ust help new`, for each tool's commands and options, or
read them all on the [documentation site](https://mustardcustarddev.github.io/detoo/).

## Agent skills

Install the skills into your project with the
[skills](https://github.com/vercel-labs/skills) CLI:

```bash
npx skills add mustardcustarddev/detoo
```

The skills use the tools to write ADRs, RFCs, user stories, problem
playbooks and SOPs, asking you for the reasons behind each one, or, for an
SOP, researching the steps for you. Others review an ADR
log, and check the code against it. The [catalogue](skills/README.md)
describes each skill, and other ways to install them.

> [!TIP]
> Let your agent write the first draft. The skills read the code, the
> existing records and the documentation of the systems involved, ask you
> about what they cannot find, and fill in your project's template. They
> write only what you told them or what they read, and every document starts
> as a Draft for your team to review.

Ask your agent in your own words, for example:

| Prompt                                                                    | Skill               |
| ------------------------------------------------------------------------- | ------------------- |
| "Write an ADR for moving our persistence to PostgreSQL."                  | `create-adr`        |
| "Draft an RFC to migrate the system to hexagonal architecture."           | `create-rfc`        |
| "Write a user story for exporting invoices."                              | `create-user-story` |
| "Write a playbook for when a customer says they've been charged twice."   | `create-playbook`   |
| "Write an SOP for rolling back a release."                                | `create-sop`        |
| "Review our ADRs, and tell me which aren't really architectural."         | `review-adrs`       |
| "Find the architecture decisions in this codebase that have no ADR."      | `find-missing-adrs` |
| "Check that the code follows ADR 7."                                      | `check-adrs`        |

## Troubleshooting

<div style="display: flex; justify-content: center;">
<img src="https://mustardcustarddev.github.io/detoo/_astro/detoo-doctor.C33HMgiX_Z24PRTX.webp" alt="detoo doctor in a terminal, reporting that the adr, rfc, sty and pse logs are all in place and named with their prefix" width="1024">
</div>

```bash
detoo doctor                # check every log directory
detoo doctor --fix          # put right what it finds
```

`detoo doctor` checks that each log directory is where its tool looks for it,
that each document is named the way its tool expects, that no two share a
number, and that each title matches its file name. It changes nothing, and exits 1 when it finds a problem, so it can run
in CI.

`detoo doctor --fix` moves and renames what is out of place, and updates the
links to it. In a git repository it offers to use `git mv`, so that git
records the renames, and to renumber documents that share a number, in the
order git added them. `--git` and `--renumber` say yes without asking.
Renumbering renames documents, so check anything else that refers to them.

## Migrating from adr-tools

A project that keeps its decisions with [adr-tools][] has them in `doc/adr`,
numbered without a prefix, such as `0001-record-architecture-decisions.md`.
`detoo doctor` reports what differs, and `detoo doctor --fix` fixes it:

```bash
detoo doctor                # report what is out of place
detoo doctor --fix          # move doc/adr to docs/adr, rename to ADR-0001-…, update links
adr list                    # adr finds the decisions again
```

The links between decisions, such as `Superseded by`, follow the renamed
files, and so do links to them from a `README.md` or other documents. If you
started adr-tools with a directory of your own, `.adr-dir` names it: if that
is `doc/adr`, delete `.adr-dir` before fixing, or `doctor` will ask you to fix
it by hand. Commit before fixing, so you can review the changes with
`git status` and `git diff`. In a git repository, `--fix` offers to use
`git mv`, so that git shows each file as renamed.

## For developers

Building needs [Go](https://go.dev/dl/) 1.26 or later, `make` and `git`. It
also builds shell completion for each tool:

```bash
make build                  # vet, test and compile into build/
make deploy                 # install into /usr/local (override with PREFIX=)
```

To release the tools, see [CONTRIBUTING.md](CONTRIBUTING.md).

## Thanks

With thanks to Nat Pryce for [adr-tools][], which inspired this project.

## License

Apache License, Version 2.0. See [LICENSE.md](LICENSE.md).

[ADRs]: https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions
[stories]: https://www.mountaingoatsoftware.com/agile/user-stories
[piosee]: https://en.wikipedia.org/wiki/Pilot_decision_making
[adr-tools]: https://github.com/npryce/adr-tools
