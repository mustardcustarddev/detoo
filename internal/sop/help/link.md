Records, in both SOPs, that SOURCE supersedes or updates TARGET.

SOURCE and TARGET are both a reference (number or partial filename) to an
SOP.  Exactly one of --supersede and --update must be given.

With --supersede, SOURCE is a new procedure for the task TARGET set out, and
TARGET is marked "Superseded by" SOURCE and loses its Draft or Approved
status.

With --update, SOURCE changes part of TARGET, such as a step for a case it
left out, and TARGET is marked "Updated by" SOURCE and keeps its status.
