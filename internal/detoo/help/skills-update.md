Updates the skills the lock file records, from the sources they were installed
from: every skill, or those named.  Each source is fetched once.

A skill whose files differ from its source's is installed again, and so is
one missing from `.agents/skills`.  The links to it need no change.  A skill
its source no longer offers is left as it is, and reported.  A skill
installed from a source detoo cannot fetch, such as a git server other than
GitHub, is left for the skills CLI to update.

When every skill could be checked, update exits 0, whether or not any was
updated; otherwise it exits 1.
