---
title: Install
description: Build and install the detoo tools from source with Go, with shell completion, and check that they work.
section: getting-started
order: 2
---

The tools are built from source, so you need [Go](https://go.dev/dl/) 1.26 or
later, `make` and `git`.

## Build and install

```bash
git clone https://github.com/mustardcustarddev/detoo.git
cd detoo
make build                  # vet, test and compile into build/
make deploy                 # install into /usr/local
```

`make deploy` builds first, so it is enough on its own. It installs `adr`,
`rfc`, `sty`, `pse` and `detoo` into `/usr/local/bin`. To install somewhere
else, such as your home directory, set `PREFIX`:

```bash
make deploy PREFIX="$HOME/.local"
```

Make sure the `bin` directory under `PREFIX` is on your `PATH`.

## Shell completion

Each tool generates its own bash completion. `make deploy` installs it into
`$PREFIX/etc/bash_completion.d`, where bash-completion picks it up. Each
tool also writes completion for zsh, fish and PowerShell, and explains how to
load it:

```bash
adr completion zsh --help
```

## Check that it works

```bash
adr --help
detoo doctor
```

`adr --help` lists the commands. `detoo doctor`, run in a project, reports on
its logs.

## Settings

The tools are configured through environment variables, named after each
tool. These are `adr`'s, and the others are the same with their own names,
such as `RFC_PAGER`.

| Variable | Effect |
| --- | --- |
| `VISUAL`, `EDITOR` | The editor `new` opens a new document in. `VISUAL` is preferred. Set both empty to create documents without opening an editor. |
| `ADR_PAGER` | A pager for the help, such as `less`. Help is not paged unless this is set. |
| `ADR_DATE` | The date written into a new record, instead of today's. |
| `ADR_TEMPLATE` | A template file for new records, instead of the log's own or the built-in one. |

## Next

Follow the [quick start](../quick-start/).
