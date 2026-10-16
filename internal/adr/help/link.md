Records, in both ADRs, that SOURCE supercedes or updates TARGET.

SOURCE and TARGET are both a reference (number or partial filename) to an ADR.
Exactly one of --supercede and --update must be given.

With --supercede, SOURCE replaces TARGET, which is marked "Superceded by"
SOURCE and is no longer accepted.

With --update, SOURCE makes substantive changes to TARGET, which is marked
"Updated by" SOURCE and keeps its status.
