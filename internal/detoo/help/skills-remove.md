Removes the skills named, or with `--all`, every skill the lock file records:
from `.agents/skills`, from each agent's directory that has a link to it or a
copy of it, and from the lock file.  A skills directory that removing them
leaves empty, such as `.claude/skills`, is removed too.

Nothing is removed unless every skill named is installed.
