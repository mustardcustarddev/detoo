Initialises the directory of standard operating procedures (SOPs),
`docs/sop`:

 * creates `docs/sop` under the current working directory, unless the
   current working directory or one of its parents already has one
 * creates the first SOP in it, which sets out how to write an SOP.

If the log already has SOPs, init does nothing.
