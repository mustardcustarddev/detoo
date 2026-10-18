Records, in both playbooks, that SOURCE replaces or updates TARGET.

SOURCE and TARGET are both a reference (number or partial filename) to a
problem playbook.  Exactly one of --replace and --update must be given.

With --replace, SOURCE is a new plan for the problem TARGET planned for, and
TARGET is marked "Replaced by" SOURCE and loses its Draft status.

With --update, SOURCE changes part of TARGET, such as a case it left out, and
TARGET is marked "Updated by" SOURCE and keeps its status.
