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
* Which item of Information 1 told you most? Did the reporter's reply leave
  out anything you needed?
* How long did building the release by SOP 3 take, against the rest of the
  drill? Did each of its steps work as written?
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

* Did the Selection rules lead you to one option, clearly?
* Which choices did you make that the expected response does not, and what
  led you to them?
* Is there anything detoo could change, such as what doctor prints, so that
  the next user who runs `--fix` on adr-tools records does not open this
  issue? That would be a new issue of its own, not a step of this playbook.

## What the drill found in the playbook

None, as the playbook and SOP 3 stand. Write down anything you found missing
or wrong in either, for revising the playbook with the `create-playbook`
skill, and SOP 3 with the `create-sop` skill.
