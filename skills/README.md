# Skills

[Agent Skills](https://agentskills.io/specification) for keeping logs with
the detoo tools: architecture decision records with `adr`, requests for
comments with `rfc`, user stories with `sty`, and PIOSEE playbooks with
`pse`. Each directory is one skill, and is self-contained, so it can be
installed on its own.

## Catalogue

### Architecture decision records

These need the `adr` tool.

* [**create-adr**](create-adr/SKILL.md): record a decision with `adr`, filling
  in the template only from evidence, and interviewing the user for the rest.
* [**review-adrs**](review-adrs/SKILL.md): judge whether each existing record
  is architectural, show the reasoning, and offer to deprecate, supersede or
  delete those that are not.
* [**find-missing-adrs**](find-missing-adrs/SKILL.md): read the whole
  codebase once, and propose records for the architectural decisions in the
  code that the log does not record.
* [**check-adrs**](check-adrs/SKILL.md): check that the code follows one
  record, or every record in force, and suggest fixes where it does not.

### Requests for comments

This needs the `rfc` tool.

* [**create-rfc**](create-rfc/SKILL.md): draft an RFC with `rfc` from a
  one-line request, reading the code for the facts, and asking the user what
  they intend and why.

### User stories

This needs the `sty` tool.

* [**create-user-story**](create-user-story/SKILL.md): write a user story with
  `sty` from a one-line subject, asking the user who it is for, what they
  need, why, and how they will know it is done. Splits epics into smaller
  stories.

### PIOSEE playbooks

This needs the `pse` tool.

* [**create-piosee**](create-piosee/SKILL.md): write a PIOSEE playbook with
  `pse`: a plan, made in advance, for a situation the team expects to face.
  It covers the scenario, the decision it calls for, the information to
  gather, the options, the rules for selecting one, the steps to execute it,
  and how to evaluate the result. It asks the user about the problem, then
  researches the code and the documentation of the systems involved, such as
  a payment provider's, to propose the investigation, the options and the
  steps for the user to confirm. It interviews the user step by step only
  where its research finds nothing. It also revises a playbook after it has
  been used.

### What they have in common

Each skill writes only what it has a source for: the repository, the
documentation it has read, or the user. It asks the user for anything else. Each also keeps its log to what belongs
there. Only architectural decisions go in the ADR log, and only changes that
others would want a say in need an RFC. A user story describes an outcome
someone wants. A PIOSEE playbook plans for a situation before it arises, by
rules the team sets.

Every ADR, RFC, user story and PIOSEE playbook an agent creates has the status
`Draft`. Accepting one is for people to do.

## Installing

### 1. Install the tools

Each skill drives `adr`, `rfc`, `sty` or `pse`, which must be on the `PATH` of
the machine the agent runs on. Build and install them as described in the
[README](../README.md#install).

### 2. Install the skills

Install with the [skills](https://github.com/vercel-labs/skills) CLI:

```bash
# See which skills the repository offers, without installing
npx skills add mustardcustarddev/detoo --list

# Install some of them into the current project
npx skills add mustardcustarddev/detoo --skill create-adr --skill create-rfc

# Install all of them
npx skills add mustardcustarddev/detoo --skill create-adr --skill review-adrs \
  --skill find-missing-adrs --skill check-adrs --skill create-rfc \
  --skill create-user-story --skill create-piosee
```

Name the skills you want. `--skill '*'` would also install
`feature-driven-vertical-slices`, from `.claude/skills/`, which is for working
on detoo itself.

Some useful options:

| Option           | Effect                                                                                            |
| ---------------- | ------------------------------------------------------------------------------------------------- |
| `-g`, `--global` | Install for your user, in every project.                                                          |
| `-a AGENT`       | Install for a particular agent, such as `-a claude-code` or `-a codex`. Repeat for more than one. |
| `-y`, `--yes`    | Skip the confirmation prompts.                                                                    |

Later, `npx skills list` shows what is installed, `npx skills update` brings
it up to date, and `npx skills remove` removes it.

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

`references/asking.md`, which says how to put every question to the user
through the agent's question tool, is copied into every skill. Keep the copies
identical:

```bash
for f in skills/*/references/asking.md; do
  diff -q skills/create-user-story/references/asking.md "$f"
done
```

A skill's directory and the `name` in its `SKILL.md` must be the same verb
phrase, such as `create-rfc`.
