Installs the agent skills that drive adr, rfc, ust, plm and sop, keeps them up
to date, and removes them.  The skills come from detoo's repository on
GitHub, mustardcustarddev/detoo, unless you name another source.

detoo installs skills where the skills CLI (`npx skills`) does, and records
them in the same lock files, so either tool can update or remove what the
other installed:

* each skill is installed once, in `.agents/skills`, which Codex, Cursor,
  Gemini CLI, GitHub Copilot, OpenCode, Amp and Cline read;
* the directory of each agent that looks elsewhere, such as Claude Code's
  `.claude/skills`, gets a link to it;
* `skills-lock.json`, in the project's directory, records where each skill
  came from, so that `detoo skills update` can update it.  Commit it with the
  skills.

The project is the nearest of the current directory and its parents that
holds `skills-lock.json`, `.agents` or `.git`, below your home directory.  If
none does, it is the current directory.

With `--global`, skills are installed for your user, in every project:
in `~/.agents/skills`, linked into each agent's own directory, such as
`~/.claude/skills`, and recorded in `~/.agents/.skill-lock.json`.

Each command says what it is doing as it goes, on stderr.
