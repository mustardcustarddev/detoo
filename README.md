# detoo

Command-line tools for keeping logs of numbered, linked Markdown records in
your repository, so they version alongside the code they describe:

* `adr` keeps [Architecture Decision Records][ADRs].
* `rfc` keeps requests for comments: proposals for significant changes, written
  down for comment before the change is made.

## Install

Build from source, which also installs shell completion for each tool:

```bash
make build                  # compile into build/
make deploy                 # install into /usr/local (override with PREFIX=)
```

## Usage

```text
adr <command> [options] [arguments]
rfc <command> [options] [arguments]
```

Both tools have the same commands. Options come before operands. Run `adr help`
or `rfc help <command>` for the built-in documentation.

### Commands

| Command | Description |
| --- | --- |
| `init [DIRECTORY]` | Start a log, and write its first record |
| `new TITLE` | Create the next numbered record and open it in your editor |
| `link (-s \| -u) SOURCE TARGET` | Record that one ADR supercedes or updates another |
| `link (-o \| -u) SOURCE TARGET` | Record that one RFC obsoletes or updates another |
| `list [--toc]` | List the records, or write a table of contents in Markdown |
| `graph [--mermaid \| --graphviz]` | Draw the links between records in Mermaid or Graphviz format |
| `config` | Print the configuration in use |
| `help [COMMAND]` | Explain a command, or list them all |
| `completion SHELL` | Write a shell completion script |

### Options

| Flag | Command | Description | Default |
| --- | --- | --- | --- |
| `-s, --supercede REFERENCE` | `adr new` | A record this one supercedes, by number or partial filename; repeatable | — |
| `-o, --obsolete REFERENCE` | `rfc new` | A record this one obsoletes, by number or partial filename; repeatable | — |
| `-u, --update REFERENCE` | `new` | A record this one updates, by number or partial filename; repeatable | — |
| `-s, --supercede` | `adr link` | SOURCE supercedes TARGET, which loses its Accepted status | — |
| `-o, --obsolete` | `rfc link` | SOURCE obsoletes TARGET, which loses its Draft status | — |
| `-u, --update` | `link` | SOURCE makes substantive changes to TARGET | — |
| `--toc` | `list` | Write a table of contents in Markdown | — |
| `-i, --intro INTRO` | `list --toc` | Markdown file to precede the contents | — |
| `-o, --outro OUTRO` | `list --toc` | Markdown file to follow the contents | — |
| `-p, --prefix LINK_PREFIX` | `list --toc`, `graph` | Prefix each link to a record | — |
| `--mermaid` | `graph` | Draw the graph in Mermaid format, the default | — |
| `--graphviz` | `graph` | Draw the graph in Graphviz format | — |
| `-e, --extension EXTENSION` | `graph` | File extension the generated links refer to | `.html` |
| `-h, --help` | all | Show help for the command | `false` |

### How the tools differ

| | adr | rfc |
| --- | --- | --- |
| Directory and marker file | `docs/adr`, `.adr-dir` | `docs/rfc`, `.rfc-dir` |
| Status of a new record | Accepted | Draft |
| Relations | Supercedes / Superceded by, Updates / Updated by | Obsoletes / Obsoleted by, Updates / Updated by |
| Built-in template | Context, Decision, Consequences | Summary, Motivation, Design, Drawbacks, Alternatives, Unresolved questions |
| Table of contents heading | Architecture Decision Records | Requests for Comments |

Superceding an ADR removes its `Accepted` line, and obsoleting an RFC removes
its `Draft` line. A record whose status has been edited to something else keeps
that status beside the link.

### Configuration

Each tool keeps its log in `docs/adr` or `docs/rfc` by default. To use another
directory, either pass it to `init`, or put its path in a `.adr-dir` or
`.rfc-dir` file in the root of the project. `config` prints the directory and
template in use.

A new record is created from `templates/template.md` inside the record
directory if the project provides one, and otherwise from the tool's built-in
template. adr's follows the style described by Michael Nygard. The placeholders
`NUMBER`, `TITLE`, `DATE` and `STATUS` are each replaced once per line.

Each tool reads its own variables, named after it: rfc reads `RFC_DATE`, and
ignores `ADR_DATE`.

| Variable | Description | Default |
| --- | --- | --- |
| `VISUAL` | Editor opened on a new record; preferred over `EDITOR` | — |
| `EDITOR` | Editor opened on a new record | — |
| `ADR_TEMPLATE`, `RFC_TEMPLATE` | Template a new record is created from | project, else built-in |
| `ADR_DATE`, `RFC_DATE` | Date written into a new record | today |
| `ADR_PAGER`, `RFC_PAGER` | Pager help is shown in; preferred over `PAGER` | — |
| `PAGER` | Pager help is shown in | `more` |

## Examples

Start a log of decisions in a directory of your choosing, and record the first
decision:

```bash
adr init doc/architecture/decisions
adr new "Implement as Go commands"
```

Supersede an earlier decision, updating the status of the old record:

```bash
adr new -s 9 "Use Rust for performance-critical functionality"
```

Record that a later decision makes substantive changes to an earlier one, which
stays accepted:

```bash
adr link --update 12 9
```

Propose a change for comment, and later replace the proposal with a better one:

```bash
rfc new "Adopt gRPC for internal services"
rfc new -o 4 "Adopt Connect for internal services"
```

Publish a table of contents, and draw the records as a diagram, with Mermaid
or with Graphviz:

```bash
adr list --toc > docs/adr/README.md
rfc graph | mmdc -i - -o docs/rfc/graph.svg
adr graph --graphviz | dot -Tpng -o docs/adr/graph.png
```

## Exit Codes

| Code | Meaning |
| --- | --- |
| `0` | Success |
| `1` | Any failure, reported on stderr and prefixed with the tool's name, such as `adr:` |

## Architecture

The application is organised by capability, in vertical slices, and knows
nothing about the command line. Cobra is an adapter that translates what the
user typed into a call on a slice and renders what comes back, so the same
behaviour could be driven by an HTTP API or an MCP server with no change to
anything under `internal/features`. Dependencies run one way, `cmd → cli → app
→ features → domain`, which makes that property checkable rather than intended:

```bash
go list -deps ./internal/features/... | grep spf13/cobra    # finds nothing
```

The tools share one engine. Each tool is a definition, in
[internal/adr](internal/adr) or [internal/rfc](internal/rfc), of the vocabulary
of its log, with the templates and help it ships with. The slices, adapters and
commands are written once and take that definition, and each `cmd` package puts
one tool together. [ADR 3](docs/adr/0003-share-one-engine-between-the-record-tools.md)
records why.

The layout is described in full by the `feature-driven-vertical-slices` skill
in [.claude/skills](.claude/skills). Run `make` on its own to list the build
targets. Tests follow the architecture, one layer at a time; the end-to-end
tests in [test](test) compare whole terminal transcripts against golden files,
one directory for each tool, which `go test ./test -update` rewrites.

With thanks to Nat Pryce for [adr-tools][], which inspired this project.

## License

Apache License, Version 2.0. See [LICENSE.txt](LICENSE.txt).

[ADRs]: http://thinkrelevance.com/blog/2011/11/15/documenting-architecture-decisions
[adr-tools]: https://github.com/npryce/adr-tools
