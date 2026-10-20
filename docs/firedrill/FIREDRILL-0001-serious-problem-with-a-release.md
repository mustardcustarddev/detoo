# 1. Serious problem with a release

Date: 2026-10-20

## Status

Draft

## Playbook

[1. Serious problem with a release](../problem/PROBLEM-0001-serious-problem-with-a-release.md).
Have it open throughout, with the SOPs it links to.

## How to play

This drill is for one person. Set it up, then read the drill a step at a
time, and do what each step says. Where a step asks you to choose, write down
your choice and why, then read the page it names, in the `problem` directory
beside this document. When a page says the drill is over, read the
`solution` directory.

A user opens an issue to say that detoo destroyed their work. You are the
detoo maintainer. Allow about 45 minutes; you will need to build the current
release and try it - most of the time is for familiarisation and decision-making.

## Ground rules

* Act only on the local system. A step that would act on anything else is
  written down, not taken. Nothing is posted on the real issues of
  `mustardcustarddev/detoo`, and no release, tag or formula is changed.
  Reading the real issues, releases and documentation is not acting on them,
  and is allowed.
* Send nothing. Write any message the drill has you send, and keep it for the
  debrief.
* Run `detoo doctor --fix` only in `~/firedrill-scratch/architecture`, the
  folder the drill has you make. It changes the nearest directory above where
  it runs that holds a `docs` or `doc` directory (`detoo help doctor`), so run
  anywhere else, it could change the detoo working copy, or another project.
* Stop at any time, for a real incident or any other reason.

## Setup

1. **The local system.** You need `git`, `make`, the GitHub CLI `gh`, and the
   Go version that `go.mod` names. From the detoo working copy, find the
   current release's tag:

   ```bash
   gh release view --repo mustardcustarddev/detoo --json tagName --jq .tagName
   ```

   Build that release, from its tag, in a scratch folder outside the working
   copy, so that nothing in the working copy changes. In place of `<tag>`,
   use the tag the command printed:

   ```bash
   git fetch --tags
   git worktree add ~/firedrill-scratch/detoo-release <tag>
   make -C ~/firedrill-scratch/detoo-release build
   ```

   `make build` vets, tests and compiles every command into `build/`
   ([README](../../README.md#for-developers)). Built from a tag with nothing
   changed, each command reports the tag's version. Do not install it, with
   `make deploy` or Homebrew: run it by its path,
   `~/firedrill-scratch/detoo-release/build/detoo`.

   No fault is added. The release behaves as it always does.
2. **What to have to hand**: the playbook,
   [SOP 2](../sop/SOP-0002-revoke-a-detoo-release.md) and
   [SOP 1](../sop/SOP-0001-release-detoo.md), and somewhere to write down
   each choice, and each comment you would post.
3. **Afterwards**, remove the worktree with
   `git worktree remove ~/firedrill-scratch/detoo-release`, and then remove
   `~/firedrill-scratch`.

## The drill

1. **The issue.** It is a weekday morning. A new issue is opened on
   `mustardcustarddev/detoo`. Read
   [the issue](FIREDRILL-0001-serious-problem-with-a-release/problem/incident.md).
2. **Your first move.** The reporter says detoo deleted their decision
   records, and asks you to pull the release. What do you do first?
   * Revoke the release now, by SOP 2, before it does the same to anyone
     else. Read
     [revoke-the-release.md](FIREDRILL-0001-serious-problem-with-a-release/problem/revoke-the-release.md).
   * Ask the reporter, on the issue, for what you need to know. Read
     [ask-the-reporter.md](FIREDRILL-0001-serious-problem-with-a-release/problem/ask-the-reporter.md).
   * Try `detoo doctor --fix` yourself, on a folder like theirs, before you
     reply. Read
     [try-it-first.md](FIREDRILL-0001-serious-problem-with-a-release/problem/try-it-first.md).
   * Close the issue: detoo has no command that deletes records. Read
     [close-at-once.md](FIREDRILL-0001-serious-problem-with-a-release/problem/close-at-once.md).
3. **The reporter's reply.** From
   [ask-the-reporter.md](FIREDRILL-0001-serious-problem-with-a-release/problem/ask-the-reporter.md)
   only. Read
   [the reply](FIREDRILL-0001-serious-problem-with-a-release/problem/reporter-reply.md).
   What do you do now?
   * Revoke the release: the reporter has confirmed it, with the output.
     Read
     [revoke-the-release.md](FIREDRILL-0001-serious-problem-with-a-release/problem/revoke-the-release.md).
   * Run what they ran, with the release they have. Go on to step 4.
   * Reply with what their output shows, and close the issue. Read
     [close-from-the-output.md](FIREDRILL-0001-serious-problem-with-a-release/problem/close-from-the-output.md).
4. **The reproduction.** Make a folder like the reporter's, with records like
   theirs, as adr-tools writes them, and run what they ran, with the release
   you built. From the detoo working copy:

   ```bash
   mkdir -p ~/firedrill-scratch
   cp -R docs/firedrill/FIREDRILL-0001-serious-problem-with-a-release/problem/architecture ~/firedrill-scratch/
   cd ~/firedrill-scratch/architecture
   ~/firedrill-scratch/detoo-release/build/detoo doctor
   ~/firedrill-scratch/detoo-release/build/detoo doctor --fix
   ```

   Look at what is in the folder now, and what is in each record. Compare
   what doctor printed with what the reporter pasted.

   Then look for other reports of the same problem with the current release,
   as Information 3 says, among the real issues of `mustardcustarddev/detoo`,
   only reading them. The drill's issue is not among them. If you find a real
   report of lost data, stop the drill, and treat it as real.
5. **The decision.** Which option does the playbook's Selection give?
   * Option 4, **Ask, and wait**. Read
     [ask-again.md](FIREDRILL-0001-serious-problem-with-a-release/problem/ask-again.md).
   * Option 2, **Keep, and fix forward**. Read
     [fix-forward.md](FIREDRILL-0001-serious-problem-with-a-release/problem/fix-forward.md).
   * Option 1, **Revoke, then fix**. Read
     [revoke-the-release.md](FIREDRILL-0001-serious-problem-with-a-release/problem/revoke-the-release.md).
   * Option 3, **Not detoo: close**. Read
     [reply-and-close.md](FIREDRILL-0001-serious-problem-with-a-release/problem/reply-and-close.md).
