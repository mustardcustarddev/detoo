# Skills

[Agent Skills](https://agentskills.io/specification) for keeping logs with
the detoo tools: architecture decision records with `adr`, requests for
comments with `rfc`, user stories with `ust`, problem playbooks with `plm`,
and standard operating procedures with `sop`. Each directory is one skill,
and is self-contained, so it can be installed on its own.

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

This needs the `ust` tool.

* [**create-user-story**](create-user-story/SKILL.md): write a user story with
  `ust` from a one-line subject, asking the user who it is for, what they
  need, why, and how they will know it is done. Splits epics into smaller
  stories.

### Problem response

These need the `plm` and `sop` tools.

* [**create-playbook**](create-playbook/SKILL.md): write a problem playbook
  with `plm`, by the PIOSEE method: a plan, made in advance, for responding
  to a problem the team expects to face, such as an incident. It covers the
  scenario, the decision it calls for, the immediate action to take (or that
  none is indicated), the actions that would make things worse, the
  information to gather, the options, the rules for selecting one, the steps
  to execute it, and how to evaluate the result, and refers to SOPs for the
  routine steps. It asks the user about the problem, then researches the
  code, the SOPs and the documentation of the systems involved, such as a
  payment provider's, to propose the rest for the user to confirm. It
  interviews the user step by step only where its research finds nothing. It
  also revises a playbook after it has been used.
* [**create-sop**](create-sop/SKILL.md): write a standard operating procedure
  with `sop`: the steps of a routine task, such as rolling back a release,
  as an ordered list. Unlike the other skills, it works from research rather
  than interview: it reads the scripts, the CI and deployment configuration
  and the documentation of the tools involved, drafts the steps with their
  sources, and asks the user only what research cannot settle.

### What they have in common

Each skill writes only what it has a source for: the repository, the
documentation it has read, or the user. It asks the user for anything else.
Each also keeps its log to what belongs there. Only architectural decisions go in the ADR log, and only changes that
others would want a say in need an RFC. A user story describes an outcome
someone wants. A problem playbook plans the response to a problem before it
arises, by rules the team sets. An SOP holds the steps of a routine task,
with no decision in it.

Every ADR, RFC, user story, playbook and SOP an agent creates has the status
`Draft`. Accepting one is for people to do.

## Installing

### 1. Install the tools

Each skill drives `adr`, `rfc`, `ust`, `plm` or `sop`, which must be on the
`PATH` of the machine the agent runs on. Build and install them as described in the
[README](../README.md#install).

### 2. Install the skills

Install them with `detoo skills`, which comes with the tools:

```bash
# See which skills there are, without installing any
detoo skills install --list

# Install all of them into the current project
detoo skills install

# Install some of them
detoo skills install create-adr create-rfc
```

`detoo skills install` installs each skill in `.agents/skills`, and links it into
the directory of each other agent it finds, such as Claude Code's
`.claude/skills`. It links it for Kiro CLI and Junie only when you name them
with `--agent`. It records what it installed in `skills-lock.json`.

Some useful options:

| Option            | Effect                                                                                                      |
| ----------------- | ----------------------------------------------------------------------------------------------------------- |
| `--global`        | Install for your user, in every project.                                                                    |
| `--agent AGENT`   | Install for a particular agent, such as `--agent claude-code` or `--agent codex`. Repeat for more than one. |
| `--source SOURCE` | Install from another source, such as a tag (`mustardcustarddev/detoo#v1.2.0`) or a clone (`../detoo`).    |

Later, `detoo skills update` brings the skills up to date, and
`detoo skills remove` removes them. Run `detoo help skills install` for every
option.

`detoo skills` keeps the same files as the
[skills](https://github.com/vercel-labs/skills) CLI, so either can update or
remove what the other installed.

### Without detoo skills

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
