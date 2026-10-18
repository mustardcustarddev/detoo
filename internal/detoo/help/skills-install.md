Installs agent skills from a source: by default, every skill in detoo's
repository on GitHub, mustardcustarddev/detoo, as it is on its default branch.
Name skills to install only those.

`--source` names another source:

* a repository on GitHub, as `OWNER/REPO`, or `OWNER/REPO#REF` for a branch,
  tag or commit, such as `mustardcustarddev/detoo#v1.2.0`;
* a directory on this machine, as a path beginning with `.`, `/` or `~`, such
  as `../detoo`.

A source's skills are the directories under its `skills` directory that hold
a `SKILL.md`.  Each is named by the `name` in its `SKILL.md`.  A source with
no `skills` directory offers its own directories that hold a `SKILL.md`.

Each skill is installed in `.agents/skills`, replacing any skill of that name
already there.  It is linked into the directory of Claude Code, and of each
other agent whose configuration directory exists in your home directory, when
that agent does not read `.agents/skills` itself, but for Kiro CLI and Junie,
which are installed for only when `--agent` names them.  `--agent` installs
for the agents named instead.  Where a link cannot be made, the skill is copied.

The lock file records each skill installed, with its source.  Nothing is
installed unless every skill and agent named is one there is.

The agents are: claude-code, codex, cursor, gemini-cli, github-copilot,
opencode, amp, cline, windsurf, roo, kiro-cli, goose, continue, junie,
qwen-code, augment and trae.

detoo downloads from GitHub's API.  For a private repository, or if GitHub
limits your requests, set `GITHUB_TOKEN` or `GH_TOKEN` to a token.
