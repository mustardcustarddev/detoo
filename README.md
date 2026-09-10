# adr

A command-line tool for keeping a log of [Architecture Decision Records][ADRs].

Records are plain Markdown files in your repository, numbered and linked, so
they version alongside the code they describe.

## Install

```bash
go install github.com/mustardcustard/decision-records/cmd/adr@latest
```

Or build from source, which also installs shell completion:

```bash
make build                  # compile into build/
make deploy                 # install into /usr/local (override with PREFIX=)
```

## Usage

```text
adr <command> [options] [arguments]
```

Options come before operands. Run `adr help` or `adr help <command>` for the
built-in documentation.

### Commands

| Command | Description |
| --- | --- |
| `adr init [DIRECTORY]` | Start a log of decision records, and write the first one |
| `adr new TITLE` | Create the next numbered record and open it in your editor |
| `adr link (-s \| -u) SOURCE TARGET` | Record that one record supercedes or updates another |
| `adr list` | List the decision records |
| `adr generate toc` | Write a table of contents in Markdown |
| `adr generate graph` | Write the links between records in Graphviz format |
| `adr config` | Print the configuration in use |
| `adr help [COMMAND]` | Explain a command, or list them all |
| `adr completion SHELL` | Write a shell completion script |

### Options

| Flag | Command | Description | Default |
| --- | --- | --- | --- |
| `-s, --supercede REFERENCE` | `new` | A decision this one supercedes, by number or partial filename; repeatable | — |
| `-u, --update REFERENCE` | `new` | A decision this one updates, by number or partial filename; repeatable | — |
| `-s, --supercede` | `link` | SOURCE supercedes TARGET, which is no longer accepted | — |
| `-u, --update` | `link` | SOURCE makes substantive changes to TARGET | — |
| `-p, --prefix LINK_PREFIX` | `generate toc`, `generate graph` | Prefix each record link | — |
| `-i, --intro INTRO` | `generate toc` | Markdown file to precede the contents | — |
| `-o, --outro OUTRO` | `generate toc` | Markdown file to follow the contents | — |
| `-e, --extension EXTENSION` | `generate graph` | File extension the generated links refer to | `.html` |
| `-h, --help` | all | Show help for the command | `false` |

### Configuration

Records live in `docs/adr` by default. To use another directory, either pass it
to `adr init`, or put its path in a `.adr-dir` file in the root of the project.
`adr config` prints the directory and template in use.

A new record is created from `templates/template.md` inside the record
directory if the project provides one, and otherwise from a built-in template
in the style described by Michael Nygard. The placeholders `NUMBER`, `TITLE`,
`DATE` and `STATUS` are each replaced once per line.

| Variable | Description | Default |
| --- | --- | --- |
| `VISUAL` | Editor opened on a new record; preferred over `EDITOR` | — |
| `EDITOR` | Editor opened on a new record | — |
| `ADR_TEMPLATE` | Template a new record is created from | project, else built-in |
| `ADR_DATE` | Date written into a new record | today |
| `ADR_PAGER` | Pager help is shown in; preferred over `PAGER` | — |
| `PAGER` | Pager help is shown in | `more` |

## Examples

Start a log in a directory of your choosing, and record the first decision:

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

Publish a table of contents, and render the decisions as a diagram:

```bash
adr generate toc > docs/adr/README.md
adr generate graph | dot -Tpng -o docs/adr/graph.png
```

## Exit Codes

| Code | Meaning |
| --- | --- |
| `0` | Success |
| `1` | Any failure, reported on stderr and prefixed `adr:` |

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

The layout is described in full by the `feature-driven-vertical-slices` skill
in [.claude/skills](.claude/skills). Run `make` on its own to list the build
targets. Tests follow the architecture, one layer at a time; the end-to-end
tests in [test](test) compare whole terminal transcripts against golden files,
which `go test ./test -update` rewrites.

With thanks to Nat Pryce for [adr-tools][], which inspired this project.

## License

Apache License, Version 2.0. See [LICENSE.txt](LICENSE.txt).

[ADRs]: http://thinkrelevance.com/blog/2011/11/15/documenting-architecture-decisions
[adr-tools]: https://github.com/npryce/adr-tools
