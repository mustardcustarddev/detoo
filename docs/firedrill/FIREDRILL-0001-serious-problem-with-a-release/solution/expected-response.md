# Expected response

The response
[1. Serious problem with a release](../../../problem/PROBLEM-0001-serious-problem-with-a-release.md)
leads to in this drill's case: a user who keeps decisions with adr-tools runs
`detoo doctor --fix`, finds `doc/adr` gone, and reports that detoo deleted
their records. Doctor moved them. Nothing was lost.

## Immediate action

None. The playbook says "No immediate action is indicated." The first thing
to do is to ask the reporter for the information, on the issue.

## Counter-indicated actions

The playbook names both of the actions the drill tempted you to take:

* **Revoking the release on the report alone**, at step 2 or step 3. The
  reporter is sure, and asks for it, but what they report is what they
  believe happened. Revoking cannot be undone, and here it would have warned
  every user off a release that does what the one before it does.
* **Closing the issue before the information is collected**, because detoo
  has no command that deletes records. If the report were right, the
  release would stay published, and more people would install it.

## Information

1. **From the reporter.** Ask on the issue for each item of Information 1:
   the version, the exact commands, what detoo printed, what was lost, and
   whether the folder is in git. The reply gives all five: the latest
   release, by `brew list --versions detoo`; `detoo doctor`,
   `detoo doctor --fix` and `ls doc/adr`, in `~/work/architecture`; what
   `--fix` printed; the three adr-tools records; and no git, so no copy.
2. **Confirmation that detoo caused it.** Run the same commands, with the
   same release, on records like theirs, by SOP 3: step 4 of the drill.
   `detoo doctor`
   reports the log as kept under `doc/` instead of `docs/`, and the file
   names as not beginning with `ADR-`. `detoo doctor --fix` prints what the
   reporter pasted, apart from its header and its line about the other
   logs, which they left out. Then, from the detoo working copy:

   ```bash
   ls -R ~/firedrill-scratch/architecture
   for f in docs/firedrill/FIREDRILL-0001-serious-problem-with-a-release/problem/architecture/doc/adr/*.md; do
     diff "$f" ~/firedrill-scratch/architecture/docs/adr/ADR-"$(basename "$f")"
   done
   ```

   `ls` shows `doc/` gone, and the three records in `docs/adr`, as
   `ADR-0001-record-architecture-decisions.md`, `ADR-0002-use-mysql.md` and
   `ADR-0003-use-postgresql.md`. `diff` compares each record with the one
   you copied: the only lines changed are the two links between records 2
   and 3, which now lead to the new names. Run in the folder, the `adr` you
   built, `~/firedrill-scratch/detoo-release/build/adr list`, lists all
   three. Nothing was lost: `ls doc/adr` failed because `doc/adr` had been
   moved, and `Removed doc/, which is empty now` is doctor removing the
   folder the move left empty. This is what `detoo help doctor` and the
   README's "Migrating from adr-tools" say `--fix` does.
3. **Other reports.** The real issues show no other report of records lost
   to `doctor --fix` with the current release. Had there been one, the drill
   would have stopped, for a real incident.

## Selection

**Option 3, No fault in detoo: reply and close**, by the rule "Select
**Option 3** when the information shows that detoo is not at fault:
something else caused the problem, or nothing was lost or damaged and detoo
did what its documentation says, such as moving files the reporter took to
be deleted." The problem reported is the loss of three records. The
information shows that they were not lost, but moved and renamed, as
documented, with their contents and links intact.

* Option 1 and Option 2 need detoo "confirmed to cause the problem". The
  reproduction confirms the opposite.
* Option 4 is for when the problem "can be confirmed neither way". After
  step 4, it can be.

## Execution

The maintainer takes each step of Option 3.

1. **Reply on the issue with what was found.** Write it down; do not post
   it. The playbook says to say first whether anything was lost, and where
   what detoo moved is now. A good reply:
   * says at once that the records are safe, since the reporter believes
     they are gone;
   * says where they are, `docs/adr`, and under what names, `ADR-0001-…`;
   * says that `--fix` moved the folder from `doc/adr` to `docs/adr`, where
     detoo's `adr` keeps decisions, renamed the files, and changed the links
     between them to follow, and that "Removed doc/" means the empty folder
     the move left behind;
   * gives `adr list` to check, and links the README's "Migrating from
     adr-tools";
   * suggests keeping a copy, or committing, before running `--fix` on
     records that are not in git, as the README advises.
2. **Close the issue.** Write it down; do not close anything.

## Evaluation

The playbook's evaluation is the maintainer's review of whether it led to a
sound decision in time, and what it lacked: that is the
[debrief](debrief.md). The checks in SOP 2 are for Option 1 only, and do not
apply. The drill shows the result on
[reply-and-close.md](../problem/reply-and-close.md): the reporter finds
their records.

## The path

* **Step 2, your first move: ask the reporter**
  ([ask-the-reporter.md](../problem/ask-the-reporter.md)). Information 1 is
  asked for on the issue, and nothing else can be decided without it.
  * [Revoke the release now](../problem/revoke-the-release.md): a
    counter-indicated action, and the second primary risk, taken on the
    report alone.
  * [Try it first](../problem/try-it-first.md): Information 2 confirms "from
    what the reporter sent". Without their version, commands and output, a
    reproduction shows only what doctor does to your folder, and the
    reporter waits unanswered.
  * [Close the issue](../problem/close-at-once.md): a counter-indicated
    action. Option 3 needs the information to show that detoo is not at
    fault, and nothing has been looked at.
* **Step 3, the reply: run what they ran** (step 4). Information 2.
  * [Revoke](../problem/revoke-the-release.md): the reporter's word is not
    confirmation, and their own output says the log was moved, not deleted.
  * [Reply from the output](../problem/close-from-the-output.md): the right
    answer, as it turns out, but reached by reading a paste, not by
    Information 2's reproduction, and without Information 3.
* **Step 5, the decision: Option 3**
  ([reply-and-close.md](../problem/reply-and-close.md)), as Selection says
  above.
  * [Option 4](../problem/ask-again.md): the information is complete, and
    asking again leaves a frightened reporter waiting for an answer you
    have.
  * [Option 2](../problem/fix-forward.md): there is no fault to fix, and the
    reporter is not told their records are safe.
  * [Option 1](../problem/revoke-the-release.md): nothing was lost.
