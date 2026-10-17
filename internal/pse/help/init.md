Initialises the directory of PIOSEE playbooks:

 * creates a subdirectory of the current working directory
 * creates the first playbook in that subdirectory, which plans for deciding
   when to write a playbook, and how.

If the DIRECTORY is not given, the playbooks are stored in the directory
`docs/piosee`.

If the log already has playbooks, init does nothing, and giving it a
different DIRECTORY is an error.

A DIRECTORY whose name begins with a hyphen must be preceded by `--`, so that
it is not read as an option.
