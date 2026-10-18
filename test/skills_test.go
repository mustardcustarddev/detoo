package clitest

import "testing"

// detoo skills installs agent skills, updates them and removes them.  These
// tests install from a directory, ~/src, laid out as detoo's repository is,
// so that they need no network; the project is ~/app, so that it is found
// below the home directory, where an agent's own directories are.

// source writes a repository of two skills, one with a reference.
func source(w *workspace) {
	w.write("src/skills/alpha/SKILL.md", "---\nname: alpha\ndescription: Does the first thing.\n---\n# Alpha\n")
	w.write("src/skills/alpha/references/notes.md", "Notes\n")
	w.write("src/skills/beta/SKILL.md", "---\nname: beta\ndescription: Does the second thing.\n---\n# Beta\n")
}

// Every skill is installed in .agents/skills, linked for Claude Code, and
// recorded; an update installs again what changed in the source, and reports
// what has gone from it; and remove removes a skill from everywhere, and the
// skills directories it leaves empty.
func TestSkillsInstallUpdateAndRemove(t *testing.T) {
	w := newWorkspace(t, "detoo")
	source(w)
	w.mkdir("app")
	w.cd("app")
	w.detoo("skills", "install", "--source", "~/src", "--list")
	w.detoo("skills", "install", "--source", "~/src")
	w.ls(".agents/skills")
	w.cat(".claude/skills/alpha/references/notes.md")
	w.cat("skills-lock.json")

	w.detoo("skills", "update")
	w.write("../src/skills/alpha/SKILL.md", "---\nname: alpha\ndescription: Does the first thing better.\n---\n# Alpha 2\n")
	w.rm("../src/skills/beta")
	w.detoo("skills", "update")
	w.cat(".claude/skills/alpha/SKILL.md")

	w.detoo("skills", "remove", "alpha")
	w.ls(".agents/skills")
	w.ls(".claude/skills")
	w.detoo("skills", "remove", "--all")
	w.ls(".")
	w.ls(".claude")
	w.cat("skills-lock.json")
	w.detoo("skills", "update")
	w.check()
}

// Skills are linked for each agent found in the home directory that does not
// read .agents/skills, but for Kiro CLI and Junie, or for the agents named.
func TestSkillsInstallForAgents(t *testing.T) {
	w := newWorkspace(t, "detoo")
	source(w)
	w.mkdir(".cursor")
	w.mkdir(".codeium/windsurf")
	w.mkdir(".kiro")
	w.mkdir(".junie")
	w.mkdir("app")
	w.cd("app")
	w.detoo("skills", "install", "--source", "~/src", "alpha")
	w.ls(".windsurf/skills")
	w.detoo("skills", "install", "--source", "~/src", "--agent", "codex", "--agent", "junie", "beta")
	w.ls(".junie/skills")
	w.check()
}

// --global installs for the user, in ~/.agents/skills, linked into the
// directories the agents read for the user.
func TestSkillsGlobal(t *testing.T) {
	w := newWorkspace(t, "detoo")
	source(w)
	w.mkdir("app")
	w.cd("app")
	w.detoo("skills", "install", "--global", "--source", "~/src", "beta")
	w.ls("../.agents/skills")
	w.cat("../.claude/skills/beta/SKILL.md")
	w.detoo("skills", "update", "--global")
	w.detoo("skills", "remove", "--global", "beta")
	w.ls("../.claude")
	w.check()
}

// Nothing is installed or removed unless what is named is there, and the
// options are checked.
func TestSkillsRefuseWhatIsNotThere(t *testing.T) {
	w := newWorkspace(t, "detoo")
	source(w)
	w.mkdir("app")
	w.cd("app")
	w.detoo("skills")
	w.detoo("skills", "install", "--source", "~/src", "alpha", "gamma")
	w.detoo("skills", "install", "--source", "~/src", "--agent", "vim")
	w.detoo("skills", "install", "--source", "nowhere")
	w.detoo("skills", "install", "--source", "./missing")
	w.detoo("skills", "install", "--source", "~/src", "--list", "alpha")
	// Nothing has been installed.
	w.detoo("skills", "update")
	w.detoo("skills", "update", "alpha")
	w.detoo("skills", "remove")
	w.detoo("skills", "remove", "--all", "alpha")
	w.detoo("skills", "remove", "alpha")
	w.check()
}
