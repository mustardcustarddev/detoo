---
name: find-missing-adrs
description: Surveys a whole codebase, thoroughly and in one pass, for architecture decisions that are embodied in the code but have no Architecture Decision Record (ADR), and proposes new ADRs based on the code as it really is - each one judged to be genuinely architectural rather than a local design choice or trivia, backed by cited evidence, and separating what the code shows from the rationale only the user can supply. Use when the user asks which ADRs are missing, wants to backfill or bootstrap an ADR log for an existing codebase, or asks what architectural decisions are undocumented.
license: Apache-2.0
compatibility: Requires the adr command-line tool on PATH, and read access to the whole repository.
---

# Find missing ADRs

Survey the codebase for the architecture decisions it already embodies, find
those the ADR log does not record, and propose records for the ones that are
genuinely architectural.

This is meant to be run once, so **be exhaustive**. Do not stop at the first
few obvious candidates, and do not sample at random. Cover every area in the
checklist, across the whole repository, and show the user what you covered.

## Ground rules

* **Describe the code as it really is**, not as the README or the directory
  names say it is. Where they disagree, report both.
* **Cite evidence for every candidate**: file paths, line numbers, and the
  commands you ran with their results. Measure rather than estimate.
* **Code shows what, never why.** Record what the code shows as fact. Do not
  infer the reasons for a decision, the alternatives considered, or who made
  it. List those as open questions for the user.
* **Only architectural decisions.** Judge every candidate with
  [significance.md](references/significance.md). Design choices and trivia are
  rejected, but listed, so the user can see they were considered.
* **Change nothing during the survey.** Use read-only commands. Records are
  created afterwards, only those the user chooses.

## Workflow

### 1. Read what is already recorded

```bash
adr config
adr list
```

Read every record, including superseded ones, so that you neither propose a
decision that is already recorded nor miss one that has been reversed. If
there is no log, the survey is still useful; creating the log comes later.

Read [significance.md](references/significance.md) and
[survey-checklist.md](references/survey-checklist.md).

### 2. Map the repository

List every tracked file (`git ls-files`, or the equivalent), and note the
top-level structure, the languages, the manifests, and the documentation.
This is the territory the survey must cover; plan to visit all of it.

If your environment can run subagents in parallel, you may divide the
checklist between them, as long as each reports evidence in the same form and
you merge and judge the results yourself.

### 3. Survey every area

Work through all twelve areas of
[survey-checklist.md](references/survey-checklist.md). For each, note what you
examined and what you found. Read the code itself: imports, call sites,
configuration, build files, tests. Do not judge from file names alone.

Record each potential decision as a candidate, with:

* the decision as practised, stated in one sentence from the code: "All
  persistence goes through `internal/store`, which wraps SQLite";
* the evidence, with paths, lines and measurements;
* how far it reaches, measured: how many packages or files depend on it.

Watch especially for:

* **Inconsistency**: two approaches to the same concern, such as two HTTP
  clients or two error styles. This is either a missing decision or an
  unfinished migration. Report it as a finding; do not pick a side.
* **Decisions written down outside the log**: a README section, a design
  document, a comment that explains why. These are the best candidates, since
  someone has already stated the decision.
* **Records that the code contradicts.** Note them for the `check-adrs` skill,
  if available, rather than proposing a duplicate.

### 4. Judge every candidate

Apply every test in the significance guide. Keep a candidate only if it is
architectural, or borderline with a clear explanation. Drop any already
covered by an existing record.

Rank what remains by reversal cost and by the risk that someone undoes it
without knowing why:

* **High**: shapes the whole system, or a contract with users; expensive to
  reverse; easy to break by accident.
* **Medium**: architectural, but narrower in reach, or well protected by
  tests or tooling.

### 5. Report

Write the report in this order.

**Coverage.** A table of all twelve checklist areas, saying what was examined
and the result. An area with nothing found says so, with what was checked.

**Proposed ADRs**, highest priority first. For each:

```markdown
### Proposed: Keep persistence behind internal/store

**Priority:** High
**What the code shows:** Every read and write of application data goes
through `internal/store` (`store.go:12-80`). No other package imports
`database/sql` (`grep -rl '"database/sql"' --include=*.go` finds only
`internal/store/store.go`).
**Why it is architectural:** reversal cost …; blast radius: 14 packages
import `internal/store` …; refactoring …; real alternatives …
**Related records:** none.
**Open questions for you:**
- What led to this decision, and when?
- Were other approaches considered?
- What has it made harder?
```

Titles should follow the style of the log's existing titles.

**Rejected candidates.** A short table of what you considered and set aside,
with the verdict (design or trivial) and the reason in one line each.

**Other findings.** Inconsistencies, records that the code appears to
contradict, and decisions documented elsewhere that disagree with the code.

If the report is long, offer to save it to a file as well.

### 6. Offer to write the records

Ask the user which proposals to record, using your environment's question
tool if it has one (such as `ask_question` or `AskUserQuestion`), allowing
several choices; otherwise ask in the conversation.

For each one chosen, follow the `create-adr` skill if it is available. If it
is not:

* interview the user about the open questions; do not answer them yourself;
* ask whether to date the record today, or on the date the decision was made,
  using `ADR_DATE`; history can show when code arrived, not when the decision
  was made;
* draft it, show it with its sources, and get approval;
* create it with `VISUAL= EDITOR= adr new "Title"`, and fill in the template.

Where the user cannot remember why a decision was made, the record can still
say what was decided and what follows from it, and state plainly that the
original reasons were not recorded.
