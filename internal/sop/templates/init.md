# Standard operating procedures

A standard operating procedure (SOP) is the steps of a routine task, such as
rolling back a release, as an ordered list, so that it is done the same way
each time. It has no decision in it: a task that needs one belongs in a
problem playbook, which refers to SOPs for its routine steps.

## Following an SOP

Do the steps in order, as written. Each names the command, page or setting to
use, and how to check that it worked. Where a step is wrong or out of date,
say so, so that the SOP is corrected.

## Writing an SOP

Ask your agent, with the `create-sop` skill, which `detoo skills install`
installs: "Write an SOP for rolling back a release." It researches the code
and the documentation of the tools involved, and asks you only what research
cannot settle.

To write one yourself, run `sop new` with the task, as an instruction, such
as `sop new "Roll back a release"`. Have someone who does the task follow it,
and correct it where it differs from what they do. An SOP is a draft until
the people who own the task approve it, when its status becomes Approved.
