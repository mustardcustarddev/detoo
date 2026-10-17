Records, in both ADRs, that SOURCE supersedes or updates TARGET.

SOURCE and TARGET are both a reference (number or partial filename) to an ADR.
Exactly one of --supersede and --update must be given.

With --supersede, SOURCE replaces TARGET, which is marked "Superseded by"
SOURCE and loses its Draft or Accepted status.

With --update, SOURCE makes substantive changes to TARGET, which is marked
"Updated by" SOURCE and keeps its status.
