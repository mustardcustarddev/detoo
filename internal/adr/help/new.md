Creates a new, numbered ADR, with the status Draft: change it to Accepted once
the decision is agreed.  The title is a single argument, so it must be quoted
if it contains spaces.  The ADR is opened for editing in the editor
specified by the VISUAL or EDITOR environment variable (VISUAL is preferred;
EDITOR is used if VISUAL is not set).  After editing, the file name of the ADR
is output to stdout, so the command can be used in scripts.

If the ADR directory contains a file `templates/template.md`, this is used as
the template for the new ADR.  Otherwise a built-in template is used that
follows the style described by Michael Nygard in this article:

  http://thinkrelevance.com/blog/2011/11/15/documenting-architecture-decisions

Multiple --supersede and --update options can be given, so that the new ADR can
supersede or update multiple existing ADRs.  A superseded ADR loses its Draft
or Accepted status; an updated one keeps its status.
