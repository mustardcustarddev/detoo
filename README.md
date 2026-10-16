# detoo

<img src="assets/detoo-avatar.png" alt="Detoo robot drawing with a set square" width="320">

Command-line tools that keep logs of numbered, linked Markdown documents in
your repository, so they are versioned alongside the code they describe:

* `adr` keeps [architecture decision records][ADRs]: decisions, and why they
  were made.
* `rfc` keeps requests for comments: proposed changes, written down for
  comment before they are made.
* `sty` keeps [user stories][stories]: who wants something of the software,
  what, and why, with the acceptance criteria that say when it is done.
* `pse` keeps decisions made with [PIOSEE][piosee]: the Problem, the
  Information, the Options, the one Selected, how it was Executed, and an
  Evaluation of whether it worked.

## Install

Build from source with Go, which also builds shell completion for each tool:

```bash
make build                  # vet, test and compile into build/
make deploy                 # install into /usr/local (override with PREFIX=)
```

## Usage

Each tool works the same way. Start a log, which by default lives in
`docs/adr`, `docs/rfc`, `docs/userstory` or `docs/piosee`, then add to it:

```bash
adr init
adr new "Use PostgreSQL for persistence"
```

`new` creates the next numbered document from a template, opens it in
`$VISUAL` or `$EDITOR`, and prints its path. Its file name is the number and a
slug of the title, such as `0002-use-postgresql-for-persistence.md`; user
stories and PIOSEE records begin with a prefix as well, as in
`story-0012-reorder-a-previous-purchase.md` and
`piosee-0004-roll-back-release-2-3.md`. A project can provide its own
template in `templates/template.md` inside the log directory.

A new document can replace or update earlier ones, and the link is recorded
in both:

```bash
adr new --supercede 9 "Use Rust for performance-critical functionality"
rfc new --obsolete 4 "Adopt Connect for internal services"
sty new --replace 12 "Reorder a previous purchase"
pse new --reverse 4 "Restore release 2.3 with the fix"
```

Publish a table of contents, or draw the links between documents with Mermaid
or Graphviz:

```bash
adr list --toc > docs/adr/README.md
sty graph | mmdc --input - --output docs/userstory/graph.svg
```

Run `adr help`, or `sty help new`, for each tool's commands and options.

## Agent skills

[skills](skills/README.md) holds agent skills that drive the tools for you:
writing ADRs, RFCs and user stories, and working decisions through PIOSEE, by
interviewing you for the reasons rather than guessing them, and reviewing and checking an ADR log. The catalogue there
says how to install them.

## Thanks

With thanks to Nat Pryce for [adr-tools][], which inspired this project.

## License

Apache License, Version 2.0. See [LICENSE.md](LICENSE.md).

[ADRs]: http://thinkrelevance.com/blog/2011/11/15/documenting-architecture-decisions
[stories]: https://www.mountaingoatsoftware.com/agile/user-stories
[piosee]: https://en.wikipedia.org/wiki/Pilot_decision_making
[adr-tools]: https://github.com/npryce/adr-tools
