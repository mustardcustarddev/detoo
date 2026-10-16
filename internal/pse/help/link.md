Records, in both records, that SOURCE reverses or updates TARGET.

SOURCE and TARGET are both a reference (number or partial filename) to a
PIOSEE record.  Exactly one of --reverse and --update must be given.

With --reverse, SOURCE undoes or abandons the decision in TARGET, which is
marked "Reversed by" SOURCE and loses its Draft status.

With --update, SOURCE builds on the decision in TARGET, such as the next
cycle of PIOSEE after an evaluation, and TARGET is marked "Updated by" SOURCE
and keeps its status.
