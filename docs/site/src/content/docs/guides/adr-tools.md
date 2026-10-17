---
title: Migrating from adr-tools
description: Move a decision log kept with adr-tools to where adr looks for it, renamed with the ADR- prefix, with its links updated.
section: guides
order: 2
---

A project that keeps its decisions with [adr-tools](https://github.com/npryce/adr-tools)
has them in `doc/adr`, numbered without a prefix, such as
`0001-record-architecture-decisions.md`. `adr` looks in `docs/adr`, for files
such as `ADR-0001-record-architecture-decisions.md`. `detoo doctor` reports
what differs, and `detoo doctor --fix` fixes it:

```bash
detoo doctor                # report what is out of place
detoo doctor --fix --git    # git mv doc/adr to docs/adr, rename to ADR-0001-…, update links
adr list                    # adr finds the decisions again
```

The links between decisions, such as `Superseded by`, follow the renamed
files, and so do links to them from other Markdown documents in the
project's root, or under `docs/` or `doc/`, such as a `README.md`. Links from
anywhere else are not changed.

If you started adr-tools with a directory of your own, `.adr-dir` names it.
If that is `doc/adr`, delete `.adr-dir` before fixing, or `doctor` will ask
you to fix it by hand.

Commit before fixing, so you can review the changes with `git status` and
`git diff`. With `--git`, git shows each file as renamed.
