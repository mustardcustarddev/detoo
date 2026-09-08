Decision Records
================

Command-line tools for working with a log of [Architecture Decision Records][ADRs]
(ADRs), written in Go.

This is a port of Nat Pryce's [adr-tools](https://github.com/npryce/adr-tools),
which was implemented as a set of Bash scripts.  The commands, the file formats
and the behaviour are the same; the whole of the original test suite is ported
along with the tool.  The repository is laid out to hold more than the one
command: each tool is a package under `cmd/`, and the logic they share lives
under `internal/`.

Quick Start
-----------

Build and install the tools:

    make build      # compile everything into build/
    make deploy     # install into /usr/local (override with PREFIX=)

Use the `adr` command to manage ADRs.  Try running `adr help`.

ADRs are stored in a subdirectory of your project as Markdown files.  The
default directory is `doc/adr`, but you can specify the directory when you
initialise the ADR log.

1.  Create an ADR directory in the root of your project:

        adr init doc/architecture/decisions

    This will create a directory named `doc/architecture/decisions` containing
    the first ADR, which records that you are using ADRs to record architectural
    decisions and links to [Michael Nygard's article on the subject][ADRs].

2.  Create Architecture Decision Records:

        adr new Implement as Go commands

    This will create a new, numbered ADR file and open it in your editor of
    choice (as specified by the `VISUAL` or `EDITOR` environment variable).

    To create a new ADR that supercedes a previous one (ADR 9, for example), use
    the `-s` option:

        adr new -s 9 Use Rust for performance-critical functionality

    This will create a new ADR, and change the status of the superceded ADR to
    record that it has been superceded by the new one.

3.  For further information, use the built-in help:

        adr help

The decision records are just Markdown files, so you can edit them by hand, and
manage them with your version control system alongside the code they describe.

Commands
--------

| Command                   | What it does                                            |
| ------------------------- | ------------------------------------------------------- |
| `adr init [DIRECTORY]`    | starts a log of decision records                          |
| `adr new TITLE...`        | creates the next numbered record                          |
| `adr link ...`            | links two records to each other                           |
| `adr list`                | lists the records                                         |
| `adr generate toc`        | writes a table of contents in Markdown                    |
| `adr generate graph`      | writes the links between records in Graphviz format       |
| `adr upgrade-repository`  | brings older records up to the current format             |
| `adr config`              | prints the configuration in use                           |
| `adr help [COMMAND]`      | explains a command, or lists them all                     |

Templates
---------

A new record is created from `templates/template.md` inside the decision record
directory, if the project provides one, and otherwise from a built-in template
in the style described by Michael Nygard.  In a template, the placeholders
`NUMBER`, `TITLE`, `DATE` and `STATUS` are replaced, once each per line, when the
record is created.

Development
-----------

`make` on its own lists the targets:

| Target    | What it does                                                     |
| --------- | ---------------------------------------------------------------- |
| `build`   | vets, tests, and compiles every command in `cmd/` into `build/`   |
| `deploy`  | installs the commands and their shell completion under `$PREFIX`  |
| `release` | tags the current commit and publishes a release                   |
| `clean`   | removes the build output                                          |

The tests live in two places:

*   `internal/adr` holds unit tests for the library the commands are built from.
*   `test` holds the ported end-to-end tests.  Each one drives the real binary
    through a scenario and compares the whole transcript against a golden file
    in `test/testdata`.  After a deliberate change to what a command prints,
    rewrite the golden files with:

        go test ./test -update

    and read the diff before committing it.

Differences from adr-tools
--------------------------

The commands behave as they did in the shell implementation, with three
exceptions:

*   Errors are reported on stderr, prefixed with `adr:`, in the Go convention.
    The shell version wrote them to stdout.
*   `adr config` prints the resolved decision record directory and template
    rather than the location of the scripts, which a single binary does not have.
*   `adr autocomplete` is a subcommand of `adr` rather than a separate
    `_adr_autocomplete` script.  `completion/adr.bash` is the bash completion
    that calls it.

Environment variables
---------------------

| Variable       | Effect                                                       |
| -------------- | ------------------------------------------------------------ |
| `VISUAL`       | the editor opened on a new record; preferred over `EDITOR`    |
| `EDITOR`       | the editor opened on a new record                             |
| `ADR_DATE`     | the date written into a new record; defaults to today         |
| `ADR_TEMPLATE` | the template a new record is created from                     |
| `ADR_PAGER`    | the pager help is shown in; preferred over `PAGER`            |
| `PAGER`        | the pager help is shown in; defaults to `more`                |

License
-------

Copyright © Nat Pryce, and the contributors to this port.

This software is licensed under the GNU General Public License, version 3 or
later.  See [LICENSE.txt](LICENSE.txt).

[ADRs]: http://thinkrelevance.com/blog/2011/11/15/documenting-architecture-decisions
