Initialises the directory of standard operating procedures (SOPs):

 * creates a subdirectory of the current working directory
 * creates the first SOP in that subdirectory, which sets out how to write
   an SOP.

If the DIRECTORY is not given, the SOPs are stored in the directory
`docs/sop`.

If the log already has SOPs, init does nothing, and giving it a different
DIRECTORY is an error.

A DIRECTORY whose name begins with a hyphen must be preceded by `--`, so that
it is not read as an option.
