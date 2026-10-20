Initialises the directory of problem playbooks, `docs/problem`:

 * creates `docs/problem` under the current working directory, unless the
   current working directory or one of its parents already has one
 * creates a `README.md` in it, which sets out what a problem playbook is,
   and how to use and write one. It is not a playbook: the first playbook is
   number 1.

If the log already has playbooks, or a `README.md`, init does nothing.
