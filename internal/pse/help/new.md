Creates a new, numbered PIOSEE playbook, in draft.  The title is a single
argument, so it must be quoted if it contains spaces.  The playbook is opened
for editing in the editor specified by the VISUAL or EDITOR environment
variable (VISUAL is preferred; EDITOR is used if VISUAL is not set).  After
editing, the file name of the playbook is output to stdout, so the command can
be used in scripts.

If the playbook directory contains a file `templates/template.md`, this is
used as the template for the new playbook.  Otherwise a built-in template is
used, with a section for the scenario and for each step of PIOSEE: problem,
information, options, selection, execution and evaluation.

Multiple --replace and --update options can be given, so that the new
playbook can replace or update multiple existing playbooks.  A replaced
playbook loses its Draft status; an updated one keeps its status.
