# Skills

[Agent Skills](https://agentskills.io/specification) for keeping an
architecture decision log with the `adr` tool. Each directory is one skill,
and is self-contained, so it can be installed on its own.

| Skill | Use it to |
| --- | --- |
| [adr-create](adr-create/SKILL.md) | Record a decision with `adr`, filling in the template only from evidence, and interviewing the user for the rest. |
| [adr-review](adr-review/SKILL.md) | Judge each existing record by whether it is truly architectural, show the reasoning, and offer to deprecate, supersede or delete the rest. |
| [adr-missing](adr-missing/SKILL.md) | Survey the whole codebase once, and propose records for the architectural decisions it embodies but the log lacks. |
| [adr-check](adr-check/SKILL.md) | Check that the code follows one record, or every record in force, and advise on remediation. |

All four share two principles. An agent writes nothing it cannot evidence, and
asks the user rather than inventing context or reasons. And decisions are made
at the right level, when they're needed: only architectural decisions belong in
the log.

## Installing

Copy or link a skill's directory into wherever your agent looks for skills,
such as `.claude/skills/` in a project. Each needs `adr` on the `PATH`; see the
[README](../README.md#install).

## Maintaining

`references/significance.md`, which defines what counts as architectural, is
copied into `adr-create`, `adr-review` and `adr-missing`, so that each skill
works on its own. Keep the copies identical:

```bash
diff skills/adr-review/references/significance.md skills/adr-create/references/significance.md
diff skills/adr-review/references/significance.md skills/adr-missing/references/significance.md
```
