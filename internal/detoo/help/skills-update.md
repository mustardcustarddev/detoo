Updates the detoo skills the lock file records from detoo's repository on
GitHub, mustardcustarddev/detoo, at the branch, tag or commit each was
installed from: every detoo skill, or those named.  Each branch, tag or commit
is fetched once.

A detoo skill is one installed from detoo's repository.  The lock file may
record skills from other sources too, installed by the skills CLI, or by
`detoo skills install --source` from a clone or a fork.  Those are left as
they are, and nothing is updated if one of them is named.  Update them with
the skills CLI, or by installing them again.

A skill whose files differ from the repository's is installed again, and so
is one missing from `.agents/skills`.  The links to it need no change.  A
skill the repository no longer offers is left as it is, and reported.

When every skill could be checked, update exits 0, whether or not any was
updated; otherwise it exits 1.
