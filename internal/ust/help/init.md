Initialises the directory of user stories, `docs/userstory`:

 * creates `docs/userstory` under the current working directory, unless the
   current working directory or one of its parents already has one
 * creates a `README.md` in it, which sets out what a user story is, and how
   to write one. It is not a story: the first story is number 1.

If the log already has stories, or a `README.md`, init does nothing.
