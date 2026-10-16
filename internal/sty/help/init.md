Initialises the directory of user stories:

 * creates a subdirectory of the current working directory
 * creates the first story in that subdirectory, describing the practice of
   keeping the work in user stories.

If the DIRECTORY is not given, the stories are stored in the directory
`docs/sty`.

If the log already has stories, init does nothing, and giving it a different
DIRECTORY is an error.

A DIRECTORY whose name begins with a hyphen must be preceded by `--`, so that
it is not read as an option.
