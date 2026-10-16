Initialises the directory of requests for comments:

 * creates a subdirectory of the current working directory
 * creates the first RFC in that subdirectory, recording the decision to
   propose significant changes in RFCs.

If the DIRECTORY is not given, the RFCs are stored in the directory `docs/rfc`.

If the log already has RFCs, init does nothing, and giving it a different
DIRECTORY is an error.

A DIRECTORY whose name begins with a hyphen must be preceded by `--`, so that
it is not read as an option.
