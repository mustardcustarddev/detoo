Command-line tools for working with a log of standard operating procedures
(SOPs).  Each SOP is the steps of a routine task, such as rolling back a
release, as an ordered list, so that it is done the same way each time.
Problem playbooks, kept with the plm tool, refer to SOPs for their routine
steps.

SOPs are stored as Markdown files in the project's `docs/sop` directory.
Each SOP's file name begins with `SOP-` and its number, such as
`SOP-0003-roll-back-a-release.md`.  Run "sop init" to start a log, and
"sop new" to add to it.