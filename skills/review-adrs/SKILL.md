---
name: review-adrs
description: Reviews an existing log of Architecture Decision Records (ADRs) for value, judging each record against the principle of decisions made at the right level, when they are needed - whether it is truly architectural, a software design decision that could be refactored without architectural change, trivial, or premature - showing the reasoning for every record, and offering to deprecate, supersede or delete low-value ones. Use when the user asks to review, audit, prune, clean up or tidy their ADRs, or asks which ADRs are worth keeping.
license: Apache-2.0
compatibility: Requires the adr command-line tool on PATH.
---

# Review ADRs

Judge every record in the ADR log by whether it belongs there, show your
reasoning for each one, and offer to retire those that do not. The guiding
principle is **decisions made at the right level, when they're needed**.

A log full of local design choices and trivia buries the decisions that shape
the system. Pruning it makes the remaining records more likely to be read and
kept.

## Ground rules

* **Evidence over impressions.** Base each verdict on the record's text and on
  the code as it stands. Where a test in the significance guide can be answered
  from the code, such as how far a dependency reaches, measure it and cite what
  you found.
* **Show your thinking for every record**, including those you would keep. The
  user must be able to disagree with a specific step.
* **Change nothing without confirmation.** This skill proposes; the user
  decides, record by record.
* **Do not invent history.** If a verdict depends on something only the user
  knows, such as whether a decision was contested at the time, ask them,
  through your question tool, as [asking.md](references/asking.md) describes.

## Workflow

### 1. Read the log

```bash
adr config        # where the log is
adr list          # every record, in order
adr graph         # how the records link to each other
```

Read every record in full. Note each one's status, and its links: a record
that is already `Superseded by` another is history, and is reviewed only for
whether it has anything left to say.

Read [significance.md](references/significance.md) before judging anything.

### 2. Judge each record

For each record, apply every test in the significance guide, then give one
verdict:

| Verdict           | Meaning                                                                               |
| ----------------- | ------------------------------------------------------------------------------------- |
| **Architectural** | Belongs in the log.                                                                   |
| **Borderline**    | The tests disagree. Say which point which way.                                        |
| **Design**        | A real decision, but local; it could be refactored away without architectural change. |
| **Trivial**       | No real alternatives, a convention, or taste.                                         |
| **Premature**     | Decided before it was needed, about something that does not exist yet.                |

Also note, where you see them, problems that affect value but are not about
level. Mention them; do not let them drive the verdict.

* **Not a decision**: it describes a fact, a task or a plan, and decides
  nothing.
* **Several decisions** in one record.
* **Duplicated** by, or contradicting, another record.
* **Stale**: the code no longer matches it. Suggest the `check-adrs` skill, if
  it is available, rather than checking in depth here.

Leave record 1, written by `adr init` to record the decision to keep ADRs,
alone unless the user asks about it.

### 3. Report

Start with a summary table, then a section for every record:

```markdown
| ADR | Title | Verdict | Recommendation |
| --- | --- | --- | --- |
| 4 | Use a builder for Config | Design | Delete, after moving its rationale to a doc comment |

### ADR 4: Use a builder for Config

**Verdict:** Design

**Reasoning:**
- Reversal cost: confined to `internal/config`; the builder is not exported
  beyond it (`grep` finds 3 call sites, all in that package).
- Refactoring: replacing it with a struct literal changes no other package.
- Real alternatives: …
- Future reader: …
- Timing: …

**Recommendation:** Delete. The reason for the builder is worth keeping, as a
doc comment on `config.Builder`.
```

Keep each piece of reasoning short and specific to the record. "This is a
local concern" is not reasoning; "only `internal/config` imports it" is.

### 4. Offer actions

For each record you would not keep as it is, ask the user what to do. Ask
about several records at once if your question tool allows it, offering these
options, with your recommendation first:

* **Keep** it as it is.
* **Deprecate** it: change its status line to `Deprecated`. The record stays
  in the log as history, and the number is never reused. This is the gentlest
  option, and usually the right one for a record that was once acted on.
* **Supersede** it, when another record, existing or new, replaces it:
  `adr link -s NEW OLD`. The tool rewrites both records' status sections. A new
  replacement record should be written with the `create-adr` skill, if
  available.
* **Delete** it. Before offering this, check what deletion would involve:
  * other records that link to it, found by searching the log for its file
    name; each link line must be removed;
  * whether it is the highest-numbered record: if so, the next `adr new` will
    reuse its number;
  * any generated table of contents, which must be regenerated with
    `adr list --toc`.

  Version control keeps the file's history, but readers of the log will no
  longer see it.
* **Move the knowledge, then retire the record.** For a design decision whose
  reasoning is still useful, offer to put that reasoning where the code is: a
  doc comment, a package README. Write it only from what the record says.

### 5. Carry out what was approved

Do exactly what the user approved, record by record:

* Deprecate: edit the line in the Status section that reads `Accepted` to
  read `Deprecated`. Leave any link lines as they are. Add a reason only if
  the user gives one.
* Supersede: run `adr link -s NEW OLD`. Never write link lines by hand.
* Delete: remove the file, and remove the link lines to it from other records.

Then run `adr list` and `adr graph` to confirm the log is consistent, offer to
regenerate any table of contents, and summarise what changed. Leave committing
to the user.
