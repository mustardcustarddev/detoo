Initialises the directory of problem playbooks, `docs/problem`:

 * creates `docs/problem` under the current working directory, unless the
   current working directory or one of its parents already has one
 * creates the first playbook in it, which plans for deciding when to write
   a playbook, and how.

If the log already has playbooks, init does nothing.
