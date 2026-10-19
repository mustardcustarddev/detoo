Initialises the directory of architecture decision records, `docs/adr`:

 * creates `docs/adr` under the current working directory, unless the
   current working directory or one of its parents already has one
 * creates the first ADR in it, recording the decision to record
   architectural decisions with ADRs.

If the log already has records, init does nothing.
