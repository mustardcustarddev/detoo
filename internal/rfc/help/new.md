Creates a new, numbered RFC, in draft.  The title is a single argument, so it
must be quoted if it contains spaces.  The RFC is opened for editing in the
editor specified by the VISUAL or EDITOR environment variable (VISUAL is
preferred; EDITOR is used if VISUAL is not set).  After editing, the file name
of the RFC is output to stdout, so the command can be used in scripts.

If the RFC directory contains a file `templates/template.md`, this is used as
the template for the new RFC.  Otherwise a built-in template is used, with
sections for the summary, motivation, design, drawbacks, alternatives and
unresolved questions.

Multiple --obsolete and --update options can be given, so that the new RFC can
obsolete or update multiple existing RFCs.  An obsoleted RFC loses its Draft
status; an updated one keeps its status.
