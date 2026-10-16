Initialises the directory of PIOSEE records:

 * creates a subdirectory of the current working directory
 * creates the first record in that subdirectory, recording the decision to
   make decisions with PIOSEE.

If the DIRECTORY is not given, the records are stored in the directory
`docs/piosee`.

If the log already has records, init does nothing, and giving it a different
DIRECTORY is an error.

A DIRECTORY whose name begins with a hyphen must be preceded by `--`, so that
it is not read as an option.
