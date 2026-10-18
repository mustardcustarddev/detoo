package skills

import "slices"

// Agent is an agent that reads skills, and where it looks for them.
type Agent struct {
	// Name is how --agent names it, as the skills CLI does: claude-code.
	Name string

	// Title is what it is called: Claude Code.
	Title string

	// Dir is the directory it reads in a project, relative to the project's
	// directory.
	Dir string

	// Global is the directory it reads for the user, in every project.
	Global Place

	// Config is the directory that says it is installed.
	Config Place
}

// Place is a directory under one of the user's directories.
type Place struct {
	Base Base
	Path string
}

// Base is a directory of the user's that others are found under.
type Base int

const (
	// Home is the user's home directory.
	Home Base = iota

	// ConfigHome is $XDG_CONFIG_HOME, or else ~/.config.
	ConfigHome

	// ClaudeHome is $CLAUDE_CONFIG_DIR, or else ~/.claude.
	ClaudeHome

	// CodexHome is $CODEX_HOME, or else ~/.codex.
	CodexHome
)

// Canonical is where each skill is installed in a project, and Global where
// each is installed for the user: the directory of the agents that follow the
// convention, which every other agent is given a link into.
const Canonical = ".agents/skills"

// Global is where each skill is installed for the user.
var Global = Place{Home, ".agents/skills"}

// DefaultAgent is installed for whether or not it seems to be installed.
const DefaultAgent = "claude-code"

// Agents are the agents skills are installed for: those that read
// .agents/skills in a project, and the best known of those that look
// elsewhere.  Their names and directories are those of the skills CLI.
var Agents = []Agent{
	{"claude-code", "Claude Code", ".claude/skills", Place{ClaudeHome, "skills"}, Place{ClaudeHome, ""}},
	{"codex", "Codex", Canonical, Global, Place{CodexHome, ""}},
	{"cursor", "Cursor", Canonical, Place{Home, ".cursor/skills"}, Place{Home, ".cursor"}},
	{"gemini-cli", "Gemini CLI", Canonical, Place{Home, ".gemini/skills"}, Place{Home, ".gemini"}},
	{"github-copilot", "GitHub Copilot", Canonical, Place{Home, ".copilot/skills"}, Place{Home, ".copilot"}},
	{"opencode", "OpenCode", Canonical, Place{ConfigHome, "opencode/skills"}, Place{ConfigHome, "opencode"}},
	{"amp", "Amp", Canonical, Place{ConfigHome, "agents/skills"}, Place{ConfigHome, "amp"}},
	{"cline", "Cline", Canonical, Global, Place{Home, ".cline"}},
	{"windsurf", "Windsurf", ".windsurf/skills", Place{Home, ".codeium/windsurf/skills"}, Place{Home, ".codeium/windsurf"}},
	{"roo", "Roo Code", ".roo/skills", Place{Home, ".roo/skills"}, Place{Home, ".roo"}},
	{"kiro-cli", "Kiro CLI", ".kiro/skills", Place{Home, ".kiro/skills"}, Place{Home, ".kiro"}},
	{"goose", "Goose", ".goose/skills", Place{ConfigHome, "goose/skills"}, Place{ConfigHome, "goose"}},
	{"continue", "Continue", ".continue/skills", Place{Home, ".continue/skills"}, Place{Home, ".continue"}},
	{"junie", "Junie", ".junie/skills", Place{Home, ".junie/skills"}, Place{Home, ".junie"}},
	{"qwen-code", "Qwen Code", ".qwen/skills", Place{Home, ".qwen/skills"}, Place{Home, ".qwen"}},
	{"augment", "Augment", ".augment/skills", Place{Home, ".augment/skills"}, Place{Home, ".augment"}},
	{"trae", "Trae", ".trae/skills", Place{Home, ".trae/skills"}, Place{Home, ".trae"}},
}

// namedOnly are the agents installed for only when --agent names them, even
// when they seem to be installed.
var namedOnly = []string{"kiro-cli", "junie"}

// Automatic reports whether skills are installed for the agent when it seems
// to be installed, without --agent naming it.
func (a Agent) Automatic() bool {
	return !slices.Contains(namedOnly, a.Name)
}

// AgentNamed returns the agent --agent calls name.
func AgentNamed(name string) (Agent, bool) {
	i := slices.IndexFunc(Agents, func(a Agent) bool { return a.Name == name })
	if i < 0 {
		return Agent{}, false
	}
	return Agents[i], true
}

// AgentNames returns the name of each agent, as --agent takes it.
func AgentNames() []string {
	names := make([]string, len(Agents))
	for i, a := range Agents {
		names[i] = a.Name
	}
	return names
}
