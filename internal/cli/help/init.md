Initialises the directory of architecture decision records:

 * creates a subdirectory of the current working directory
 * creates the first ADR in that subdirectory, recording the decision to
   record architectural decisions with ADRs.

If the DIRECTORY is not given, the ADRs are stored in the directory `docs/adr`.

If the log already has records, init does nothing, and giving it a different
DIRECTORY is an error.

A DIRECTORY whose name begins with a hyphen must be preceded by `--`, so that
it is not read as an option.
