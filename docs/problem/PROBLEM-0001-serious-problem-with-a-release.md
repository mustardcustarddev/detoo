# 1. Serious problem with a release

Date: 2026-10-18

## Status

Accepted

## Scenario

A GitHub issue on `mustardcustarddev/detoo` reports that a published release
of detoo causes a serious problem, such as losing data.

## Problem

Whether the release must be revoked, by
[SOP 3. Revoke a detoo release](../sop/SOP-0003-revoke-a-detoo-release.md),
or the report is dealt with otherwise.

Primary risks:

* Leaving a harmful release published, so that more people install it.
* Revoking a release for a problem detoo did not cause. If it is the only
  release, detoo cannot be installed with Homebrew until the next one.

## Immediate action

No immediate action is indicated.

## Information

1. **From the reporter**, asked for on the issue:
   * the version, from `brew list --versions detoo`, or from the version
     beside the detoo badge in the header of `detoo doctor`;
   * the exact commands they ran, in order, and in which project;
   * what was lost or damaged: which files, and what they held before;
   * whether the project is a git repository, and whether what was lost had
     been committed.
2. **Confirmation that detoo caused it**, from what the reporter sent: run
   the same commands, with the same version, in a scratch project. Start with
   the commands that change existing files:
   * `new` with `--supersede` or `--update`, which changes the records it
     names as well as writing the new one;
   * `detoo doctor --fix`, which moves and renames logs and documents.
3. **Other reports** of the same problem with the same version, among the
   repository's issues. Several identical reports confirm it together, even
   when each has less information than one good one.

## Options

1. **Revoke, then fix.** Revoke the release, and release the fix.
2. **Keep, and fix forward.** Leave the release published, and release the
   fix.
3. **Not detoo: close.** Reply on the issue with what was found, and close it.
4. **Ask, and wait.** The conservative fallback. Ask the reporter for what is
   missing, label the issue `question`, and decide again when more comes in.

## Selection

Evaluate every case. There is no default.

* Select **Option 1** when detoo is confirmed to cause the problem, by the
  reporter's information or by several identical reports, and the problem is
  serious, such as data loss.
* Select **Option 2** when detoo is confirmed to cause the problem, but it is
  not serious.
* Select **Option 3** when the information shows that detoo did not cause
  it.
* Select **Option 4** when it can be confirmed neither way.

Serious is a question of judgment, besides data loss.

## Execution

The detoo maintainer takes each step.

* **Option 1**
  1. Revoke the release by
     [SOP 3. Revoke a detoo release](../sop/SOP-0003-revoke-a-detoo-release.md).
  2. Reply on the issue, with a link to the notice SOP 3 posts on the
     documentation site.
  3. Fix the problem, and release the fix by
     [SOP 2. Release detoo](../sop/SOP-0002-release-detoo.md).
  4. Reply on the issue with the version that fixes it.
* **Option 2**
  1. Fix the problem, and release the fix by
     [SOP 2. Release detoo](../sop/SOP-0002-release-detoo.md).
  2. Reply on the issue with the version that fixes it.
* **Option 3**
  1. Reply on the issue with what was found.
  2. Close the issue.
* **Option 4**
  1. Ask on the issue for what is missing from the information in
     Information 1.
  2. Label the issue `question`.
  3. Return to Information when the reporter replies, or another report
     arrives.

## Evaluation

* After the playbook is used, the maintainer reviews whether it led to a
  sound decision in time, and what it lacked.
* After Option 1, the checks in SOP 3 show whether the release was revoked:
  that the release is gone and its tag kept, that the notice and the README
  warning are published, and that Homebrew installs the earlier release.
