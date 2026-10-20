Initialises the directory of fire drills, `docs/firedrill`:

 * creates `docs/firedrill` under the current working directory, unless the
   current working directory or one of its parents already has one
 * creates a `README.md` in it, which sets out how to practise a playbook in
   a fire drill. It is not a drill: the first drill is number 1.

If the log already has fire drills, or a `README.md`, init does nothing.
