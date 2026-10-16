Creates a new, numbered user story, in draft.  The title is a single argument,
so it must be quoted if it contains spaces.  The story is opened for editing in
the editor specified by the VISUAL or EDITOR environment variable (VISUAL is
preferred; EDITOR is used if VISUAL is not set).  After editing, the file name
of the story is output to stdout, so the command can be used in scripts.

If the story directory contains a file `templates/template.md`, this is used
as the template for the new story.  Otherwise a built-in template is used, with
sections for the story, its acceptance criteria and notes, in the form
described by Mike Cohn in this guide:

  https://www.mountaingoatsoftware.com/agile/user-stories

Multiple --replace and --update options can be given, so that the new story
can replace or update multiple existing stories.  A replaced story loses its
Draft status; an updated one keeps its status.
