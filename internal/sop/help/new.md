Creates a new, numbered standard operating procedure (SOP), in draft.  The
title is a single argument, so it must be quoted if it contains spaces.  The
SOP is opened for editing in the editor specified by the VISUAL or EDITOR
environment variable (VISUAL is preferred; EDITOR is used if VISUAL is not
set).  After editing, the file name of the SOP is output to stdout, so the
command can be used in scripts.

If the SOP directory contains a file `templates/template.md`, this is used as
the template for the new SOP.  Otherwise a built-in template is used: a
title, a date, a status, and the steps, as an ordered list.

Multiple --supersede and --update options can be given, so that the new SOP
can supersede or update multiple existing SOPs.  A superseded SOP loses its
Draft or Approved status; an updated one keeps its status.
