Initialises the directory of requests for comments, `docs/rfc`:

 * creates `docs/rfc` under the current working directory, unless the
   current working directory or one of its parents already has one
 * creates the first RFC in it, recording the decision to propose
   significant changes in RFCs.

If the log already has RFCs, init does nothing.
