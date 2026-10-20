Command-line tools for working with a log of fire drills.  Each fire drill is
an exercise in which one person practises a problem playbook, kept with the
plm tool, on a local copy of the system.  It is read like a book in which the
reader chooses what happens: the drill sets out the setup and what happens, a
step at a time, and each choice leads to a page of what happens next.
Nothing is broken to make the incident.

Fire drills are stored as Markdown files in the project's `docs/firedrill`
directory.  Each drill's file name begins with `FIREDRILL-` and its number,
such as `FIREDRILL-0002-customer-charged-twice.md`, and beside it is a
directory named as it is, which holds its `problem` and `solution` pages.
Run "fyr init" to start a log, and "fyr new" to add to it.