Records, in both stories, that SOURCE replaces or updates TARGET.

SOURCE and TARGET are both a reference (number or partial filename) to a user
story.  Exactly one of --replace and --update must be given.

With --replace, SOURCE replaces TARGET, which is marked "Replaced by" SOURCE
and loses its Draft status.  When a story is split, each of the smaller
stories replaces it.

With --update, SOURCE makes substantive changes to TARGET, which is marked
"Updated by" SOURCE and keeps its status.
