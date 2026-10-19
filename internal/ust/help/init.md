Initialises the directory of user stories, `docs/userstory`:

 * creates `docs/userstory` under the current working directory, unless the
   current working directory or one of its parents already has one
 * creates the first story in it, describing the practice of keeping the
   work in user stories.

If the log already has stories, init does nothing.
