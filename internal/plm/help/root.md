Command-line tools for working with a log of problem playbooks.  Each playbook
plans, with PIOSEE, how to respond to a problem, such as an incident, before
it arises: the Scenario, the Problem to decide, the Immediate action to take,
the Counter-indicated actions to avoid, the Information to gather, the
Options, how to make the Selection, the steps of Execution, and the
Evaluation of the result.  A playbook refers to standard operating procedures,
kept with the sop tool, for the routine steps.

Problem playbooks are stored as Markdown files in a subdirectory of the
project, by default `docs/problem`.  Each playbook's file name begins with
`PROBLEM-` and its number, such as
`PROBLEM-0001-plan-for-incidents-with-piosee.md`.  Run "plm init" to start a
log, and "plm new" to add to it.
