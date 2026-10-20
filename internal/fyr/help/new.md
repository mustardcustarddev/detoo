Creates a new, numbered fire drill, in draft.  The title is a single argument,
so it must be quoted if it contains spaces.  The fire drill is opened for
editing in the editor specified by the VISUAL or EDITOR environment variable
(VISUAL is preferred; EDITOR is used if VISUAL is not set).  After editing,
the file name of the fire drill is output to stdout, so the command can be
used in scripts.

If the fire drill directory contains a file `templates/template.md`, this is
used as the template for the new fire drill.  Otherwise a built-in template is
used: a title, a date, a status, and a section for the playbook it practises,
how to play, the ground rules, the setup, and the drill itself.

Beside the fire drill, new creates a directory named as it is without `.md`,
which holds two directories, each with a README that says what goes in it:

 * `problem`, for the pages the drill sends the reader to: the incident, the
   evidence, and what happens after each choice;
 * `solution`, for the expected response and the debrief, read when the
   drill is over.

Multiple --replace and --update options can be given, so that the new fire
drill can replace or update multiple existing fire drills.  A replaced fire
drill loses its Draft status; an updated one keeps its status.
