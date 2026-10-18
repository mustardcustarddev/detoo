---
title: Installing skills
description: Install the detoo agent skills into a project, or for every project, with detoo skills, and keep them up to date.
section: skills
order: 2
---

## 1. Install the tools

Each skill drives `adr`, `rfc`, `ust`, `plm` or `sop`, which must be on the
`PATH` of the machine the agent runs on. [Install](../../install/) them first.

## 2. Install the skills

Install them with `detoo skills`, which comes with the tools:

```bash
# See which skills there are, without installing any
detoo skills install --list

# Install all of them into the current project
detoo skills install

# Install some of them
detoo skills install create-adr create-rfc
```

`detoo skills install` downloads the skills from the
[repository](https://github.com/mustardcustarddev/detoo/tree/main/skills) on
GitHub. It installs each one in `.agents/skills`, which Codex, Cursor, Gemini
CLI, GitHub Copilot, OpenCode, Amp and Cline read. Then it links it into the
directory of each other agent it finds, such as Claude Code's
`.claude/skills`. It links it for Kiro CLI and Junie only when you name them
with `--agent`. It records what it installed in `skills-lock.json`. Commit
that file with the skills.

Some useful options:

| Option | Effect |
| --- | --- |
| `--global` | Install for your user, in every project. |
| `--agent AGENT` | Install for a particular agent, such as `--agent claude-code` or `--agent codex`. Repeat for more than one. |
| `--source SOURCE` | Install from another source: a branch, tag or commit, such as `mustardcustarddev/detoo#v1.2.0`, a fork on GitHub, or a directory, such as `../detoo`. |

Later, `detoo skills update` brings the skills up to date, and
`detoo skills remove` removes them. The
[command reference](../../detoo/skills/) has every option.

`detoo skills` keeps the same files as the
[skills](https://github.com/vercel-labs/skills) CLI, so if you installed the
skills with `npx skills add`, `detoo skills` can update and remove them, and
the other way round.

## Without detoo skills

Each skill is a self-contained directory under
[`skills/`](https://github.com/mustardcustarddev/detoo/tree/main/skills) in
the repository, so you can also copy or link it into wherever your agent
looks for skills, such as `.claude/skills/` in a project or
`~/.claude/skills/` for your user.
