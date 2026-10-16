Creates a new, numbered PIOSEE record, in draft.  The title is a single
argument, so it must be quoted if it contains spaces.  The record is opened
for editing in the editor specified by the VISUAL or EDITOR environment
variable (VISUAL is preferred; EDITOR is used if VISUAL is not set).  After
editing, the file name of the record is output to stdout, so the command can
be used in scripts.

If the record directory contains a file `templates/template.md`, this is used
as the template for the new record.  Otherwise a built-in template is used,
with a section for each step of PIOSEE: problem, information, options,
select, execute and evaluate.

Multiple --reverse and --update options can be given, so that the new record
can reverse or update multiple existing records.  A reversed record loses its
Draft status; an updated one keeps its status.
