# Skills

[Agent Skills](https://agentskills.io/specification) for keeping an
architecture decision log with the `adr` tool, and a log of requests for
comments with the `rfc` tool. Each directory is one skill, and is
self-contained, so it can be installed on its own.

## Catalogue

### Architecture decision records

These need the `adr` tool.

* [**create-adr**](create-adr/SKILL.md): record a decision with `adr`, filling
  in the template only from evidence, and interviewing the user for the rest.
* [**review-adrs**](review-adrs/SKILL.md): judge each existing record by
  whether it is truly architectural, show the reasoning, and offer to
  deprecate, supersede or delete the rest.
* [**find-missing-adrs**](find-missing-adrs/SKILL.md): survey the whole
  codebase once, and propose records for the architectural decisions it
  embodies but the log lacks.
* [**check-adrs**](check-adrs/SKILL.md): check that the code follows one
  record, or every record in force, and advise on remediation.

### Requests for comments

This needs the `rfc` tool.

* [**create-rfc**](create-rfc/SKILL.md): draft an RFC with `rfc` from a
  one-line request, by reading the template, gathering facts, and interviewing
  the user for their intent and reasons.

### What they share

All of them share two principles. An agent writes nothing it cannot evidence,
and asks the user rather than inventing context or reasons. And decisions are
made at the right level, when they're needed: only architectural decisions
belong in the ADR log, and only changes others would want a say in need an RFC.

Every ADR and RFC an agent creates has the status `Draft`. Accepting one is for
people to do.

## Installing

### 1. Install the tools

Each skill drives `adr` or `rfc`, which must be on the `PATH` of the machine
the agent runs on. Build and install them as described in the
[README](../README.md#install).

### 2. Make sure you can reach the repository

detoo is a private repository, so you need read access to
`mustardcustarddev/detoo` on GitHub, and Git set up to use it. Any one of
these will do:

* a Git credential helper that already has your GitHub credentials;
* the GitHub CLI, signed in with `gh auth login`;
* an SSH key added to your GitHub account.

If you can `git clone https://github.com/mustardcustarddev/detoo.git`, or
`git clone git@github.com:mustardcustarddev/detoo.git`, you are ready.

### 3. Install the skills

Install with the [skills](https://github.com/vercel-labs/skills) CLI. It uses
the same command for private repositories as for public ones, and the Git or
GitHub CLI authentication you already have:

```bash
# See which skills the repository offers, without installing
npx skills add mustardcustarddev/detoo --list

# Install some of them into the current project
npx skills add mustardcustarddev/detoo --skill create-adr --skill create-rfc

# Install all of them
npx skills add mustardcustarddev/detoo --skill create-adr --skill review-adrs \
  --skill find-missing-adrs --skill check-adrs --skill create-rfc
```

Name the skills you want, rather than using `--skill '*'`: the CLI also finds
`feature-driven-vertical-slices`, in `.claude/skills/`, which is for working
on detoo itself, not for keeping a log.

If HTTPS does not work for you, give the repository's SSH address instead:

```bash
npx skills add git@github.com:mustardcustarddev/detoo.git --skill create-rfc
```

Some useful options:

| Option           | Effect                                                                                            |
| ---------------- | ------------------------------------------------------------------------------------------------- |
| `-g`, `--global` | Install for your user, in every project, rather than in the current project.                      |
| `-a AGENT`       | Install for a particular agent, such as `-a claude-code` or `-a codex`. Repeat for more than one. |
| `-y`, `--yes`    | Skip the confirmation prompts.                                                                    |

Later, `npx skills list` shows what is installed, `npx skills update` brings
it up to date, and `npx skills remove` removes it.

For GitHub, the CLI tries your Git credentials first, then the GitHub CLI,
then SSH. It never reads the GitHub CLI's stored token. Setting `GITHUB_TOKEN`
or `GH_TOKEN` is optional, and only needed if none of those is set up; see the
CLI's notes on
[private repositories](https://github.com/vercel-labs/skills#private-repositories).

### Without the CLI

Each skill is a self-contained directory, so you can also copy or link it
into wherever your agent looks for skills, such as `.claude/skills/` in a
project or `~/.claude/skills/` for your user.

## Maintaining

`references/significance.md`, which defines what counts as architectural, is
copied into `create-adr`, `review-adrs` and `find-missing-adrs`, so that each
skill works on its own. Keep the copies identical:

```bash
diff skills/review-adrs/references/significance.md skills/create-adr/references/significance.md
diff skills/review-adrs/references/significance.md skills/find-missing-adrs/references/significance.md
```

A skill's directory and the `name` in its `SKILL.md` must be the same verb
phrase, such as `create-rfc`.
