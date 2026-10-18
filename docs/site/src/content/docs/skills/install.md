---
title: Installing skills
description: Install the detoo agent skills into a project, or for every project, with the skills CLI or by copying them.
section: skills
order: 2
---

## 1. Install the tools

Each skill drives `adr`, `rfc`, `sty`, `plm` or `sop`, which must be on the
`PATH` of the machine the agent runs on. [Install](../../install/) them first.

## 2. Install the skills

Install with the [skills](https://github.com/vercel-labs/skills) CLI:

```bash
# See which skills the repository offers, without installing
npx skills add mustardcustarddev/detoo --list

# Install some of them into the current project
npx skills add mustardcustarddev/detoo --skill create-adr --skill create-rfc

# Install all of them
npx skills add mustardcustarddev/detoo --skill '*'
```

Some useful options:

| Option | Effect |
| --- | --- |
| `--global` | Install for your user, in every project. |
| `--agent AGENT` | Install for a particular agent, such as `--agent claude-code` or `--agent codex`. Repeat for more than one. |
| `--yes` | Skip the confirmation prompts. |

Later, `npx skills list` shows what is installed, `npx skills update` brings
it up to date, and `npx skills remove` removes it.

## Without the CLI

Each skill is a self-contained directory under
[`skills/`](https://github.com/mustardcustarddev/detoo/tree/main/skills) in
the repository, so you can also copy or link it into wherever your agent
looks for skills, such as `.claude/skills/` in a project or
`~/.claude/skills/` for your user.
