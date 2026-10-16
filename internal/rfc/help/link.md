Records, in both RFCs, that SOURCE obsoletes or updates TARGET.

SOURCE and TARGET are both a reference (number or partial filename) to an RFC.
Exactly one of --obsolete and --update must be given.

With --obsolete, SOURCE replaces TARGET, which is marked "Obsoleted by" SOURCE
and loses its Draft status.

With --update, SOURCE makes substantive changes to TARGET, which is marked
"Updated by" SOURCE and keeps its status.
