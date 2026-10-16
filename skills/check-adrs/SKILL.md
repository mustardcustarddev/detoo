---
name: check-adrs
description: Checks that the codebase adheres to one Architecture Decision Record (ADR) or to the whole ADR log - turning each decision into concrete, checkable rules, searching the code for evidence of compliance and of violation, and reporting a verdict per rule with cited evidence and remediation advice, either to bring the code back in line or to record a new decision that supersedes the old. Use when the user asks whether the code follows the ADRs, to verify, audit or enforce an ADR, to find architecture drift or violations, or before relying on an ADR as a description of the system.
license: Apache-2.0
compatibility: Requires the adr command-line tool on PATH, and read access to the whole repository.
---

# Check adherence to ADRs

Check whether the code does what the ADR log says was decided: for one record,
or for every record in force. Report each rule with evidence, and advise what
to do where code and record disagree.

## Ground rules

* **Every verdict needs evidence.** Cite file paths, line numbers, and the
  commands you ran with their results. Never report a violation, or compliance,
  that you have not found in the code.
* **Look for violations, not just confirmations.** Finding one file that
  follows a rule proves little. Search the whole repository for code that
  breaks it: tests, scripts, build files and configuration included.
* **Read-only.** Do not change code or records while checking. Remediation is
  advice, carried out only if the user asks.
* **Do not decide which side is right.** When code and record disagree, either
  may be wrong. Lay out both remedies; the user chooses. If you need facts only
  they have, such as whether a change was deliberate, ask them, using your
  environment's question tool if it has one (such as `ask_question` or
  `AskUserQuestion`), or in the conversation.

## Workflow

### 1. Choose the records

```bash
adr config
adr list
```

If the user named a record, by number or part of its title, check that one.
Otherwise check every record.

Read each record's Status section:

* **Accepted** records are in force: check them.
* A record that is **Updated by** another is in force as amended: read both,
  and check them together.
* **Superceded** (the tool's spelling), **Deprecated** or **Rejected** records
  are not in force: skip them, but check that the record replacing a
  superseded one is followed.
* **Draft** or **Proposed** records are not yet in force: check them only if
  the user asks, and report them separately.

### 2. Turn each decision into rules

Read the Decision section, and the Consequences where they state something
that must hold. Write each obligation as a rule that can be tested, with how
you will test it:

| Kind of rule | Example                                     | How to check                                    |
| ------------ | ------------------------------------------- | ----------------------------------------------- |
| Structure    | `features` must not import `cli`            | List each package's imports                     |
| Dependency   | Cobra appears only in `cli` and `cmd`       | Find every importer of the module               |
| Technology   | All persistence uses PostgreSQL             | Search for other database drivers and clients   |
| Convention   | Every record type lives in `domain/record`  | Search for the types defined elsewhere          |
| Contract     | The API is versioned under `/v1`            | Read the route definitions                      |
| Technique    | End-to-end tests compare golden transcripts | Read the tests, and look for others that do not |

Stay within what the record says. Do not add rules it does not state, however
sensible. If the record is too vague to test, report that as a finding rather
than guessing what it meant.

Some decisions cannot be checked from code, such as a review process or
who approves releases. Mark them **Not checkable** and say why.

### 3. Gather the evidence

For each rule, search the whole repository. Prefer exact, repeatable checks:
the language's dependency tools (for example `go list -deps`,
`npm ls`, `pipdeptree`), searches for imports and identifiers, and reading the
code at every match. Record each command and its result, so the user can run
it again.

Note what the record refers to that no longer exists: packages, files,
technologies. A record about code that is gone may be **Obsolete**.

Look also for records that contradict each other: two records in force that
cannot both be followed.

### 4. Give a verdict per rule

| Verdict           | Meaning                                                                                 |
| ----------------- | --------------------------------------------------------------------------------------- |
| **Adheres**       | Searched thoroughly; no violation found. Say what was searched.                         |
| **Violates**      | Code breaks the rule. List every instance, or the count and a sample if there are many. |
| **Partial**       | Followed in some places, broken in others, as in an unfinished migration.               |
| **Undetermined**  | The evidence is inconclusive. Say what would settle it.                                 |
| **Not checkable** | The rule is not about code.                                                             |
| **Obsolete**      | The record concerns code that no longer exists.                                         |

Give a record an overall verdict from its rules: the worst of them, apart from
Not checkable.

### 5. Advise on remediation

For every rule that is violated, partial or obsolete, give the user both
directions, with what each would involve:

* **Bring the code into line.** Which files change, what the change is, and
  roughly how much: a handful of call sites, or a migration. Point out any risk,
  such as a contract other code depends on.
* **Change the decision.** If the code reflects a deliberate new direction,
  the record should not be edited to match. A new record should supersede it,
  or update it if the change is partial. Suggest the `create-adr` skill, if
  available, which will ask the user why the direction changed.

Where you can, also suggest a way to stop the drift recurring, such as a test,
a lint rule, or a CI step that runs the check you used. Offer it; do not add
it unasked.

For a record too vague to test, suggest clarifying it with a new record that
updates it. For an obsolete record, suggest deprecating or superseding it; the
`review-adrs` skill, if available, covers this.

### 6. Report

```markdown
## Summary

| ADR | Title | Verdict | Rules: adheres / violates / other |
| --- | --- | --- | --- |
| 3 | Share one engine between the record tools | Adheres | 4 / 0 / 1 |
| 5 | Keep persistence behind internal/store | Violates | 1 / 1 / 0 |

## ADR 5: Keep persistence behind internal/store

### Rule 5.1: Only internal/store imports database/sql — Violates

**Checked with:** `grep -rln '"database/sql"' --include='*.go' .`
**Found:** `internal/store/store.go` (allowed), and
`internal/report/export.go:9` (violation), which queries the database
directly at lines 40-58.

**Remediation:**
- *Bring the code into line:* move the query into `internal/store` as a new
  method, and call it from `export.go`. One file, one call site.
- *Change the decision:* if reports are meant to read the database directly,
  record that in a new ADR that updates ADR 5.
- *Prevent recurrence:* a test that fails if any package but `internal/store`
  imports `database/sql`.
```

List every rule of every record checked, including those that adhere, so the
user can see what was covered. End with the violations in order of risk, and
ask which, if any, the user wants to act on.
