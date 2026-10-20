# Debrief

Look back on the drill with these questions, and with what you wrote down at
each choice. Then read what the drill found in the playbook.

## Detect

* The drill started with an issue, and GitHub's notification of it. How
  soon would you have seen a real one? Is a new issue on detoo noticed the
  same day?
* The reporter wrote "GONE" and "Please pull this release". Did that change
  how quickly you wanted to act, and how?

## Identify

* Did the reporter's paste tell you what had happened before you reproduced
  it? Which line: `Moved from doc/ to docs/`, `doc/adr  →  docs/adr`, or
  `Removed doc/, which is empty now`? Which line misled the reporter?
* Did Information 1 ask for everything you needed? The reply had what
  `--fix` printed only because the reporter chose to paste it.
* How long did building the release from its tag take, against the rest of
  the drill? Would you have known how to run the reporter's version without
  installing it, without the setup's steps?
* Did you compare the records' contents after the fix, or only their names?
  What would have shown you that a record had been changed?

## Communicate

* Your first comment, at step 2: did it ask for each item of Information 1,
  in words the reporter could act on, such as the command that gives the
  version?
* Your final reply: did it say first that the records were safe? Did it say
  where they were, what "Removed doc/" meant, and how to check, as the
  expected response sketches?
* Did anything you wrote blame the reporter for not reading the output, or
  for not keeping a copy?

## Resolve

* Did the Selection rules lead you to one option, clearly? Did "Not detoo"
  feel like the right name for a case where detoo moved the files, as asked?
* Which choices did you make that the expected response does not, and what
  led you to them?
* Is there anything detoo could change, such as what doctor prints, so that
  the next user who runs `--fix` on adr-tools records does not open this
  issue? That would be a new issue of its own, not a step of this playbook.

## What the drill found in the playbook

These are for revising the playbook with the `create-playbook` skill.

* **Option 3 does not quite fit this case.** Its name, "Not detoo: close",
  and its rule, "the information shows that detoo did not cause it", are
  written for a problem with another cause. Here detoo did what the reporter
  says it did, removing `doc/`, but nothing was lost: detoo worked as
  documented, and the report is a misreading. The rule fits only if "it" is
  read as the harm. A rule, or an option, for "no problem: working as
  documented" would make the selection clear.
* **There is no Counter-indicated actions section.** Revoking on the report
  alone, before Information 2 confirms it, is the action a frightened
  reporter asks for, and the one that cannot be undone. It is only implied
  by the primary risks.
* **Information 1 does not ask for the output.** What `--fix` printed held
  the answer, and lets Information 2 compare. The reporter gave it
  unasked.
* **Information 2 does not say how to run "the same version".** Running the
  reporter's release without installing it, by building it from its tag in
  a worktree outside the working copy, had to come from the drill's setup.
  It could be a step of the playbook, or an SOP of its own.
* **Option 3's reply has no guidance on recovery.** A reporter who believes
  their data is lost needs to be told first that it is not, and where it
  is. Execution says only "Reply on the issue with what was found."
