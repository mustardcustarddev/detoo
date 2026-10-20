Initialises the directory of standard operating procedures (SOPs),
`docs/sop`:

 * creates `docs/sop` under the current working directory, unless the
   current working directory or one of its parents already has one
 * creates a `README.md` in it, which sets out what an SOP is, and how to
   follow and write one. It is not an SOP: the first SOP is number 1.

If the log already has SOPs, or a `README.md`, init does nothing.
