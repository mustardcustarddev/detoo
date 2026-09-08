package adr

import (
	"sort"
	"strings"
)

// Command is one of the commands of the adr tool.  A command's Help is the text
// shown by "adr help COMMAND", and starts with its usage line.
type Command struct {
	Name   string
	Help   string
	Hidden bool
	Run    func(*Env, []string) error
	Subs   []*Command
}

// commands are the commands of the adr tool, in the order they are listed.
// It is populated by init to break the initialisation cycle between a command
// and the help that lists every command.
var commands []*Command

func init() {
	commands = []*Command{
		{Name: "config", Help: configHelp, Run: runConfig},
		{Name: "generate", Help: generateHelp, Run: runGenerate, Subs: generators},
		{Name: "help", Help: helpHelp, Run: runHelp},
		{Name: "init", Help: initHelp, Run: runInit},
		{Name: "link", Help: linkHelp, Run: runLink},
		{Name: "list", Help: listHelp, Run: runList},
		{Name: "new", Help: newHelp, Run: runNew},
		{Name: "upgrade-repository", Help: upgradeHelp, Run: runUpgrade},
		{Name: "autocomplete", Hidden: true, Run: runAutocomplete},
	}
	sort.Slice(commands, func(i, j int) bool { return commands[i].Name < commands[j].Name })
}

// Run executes a command line.
func Run(e *Env, args []string) error {
	var cmd *Command
	if len(args) > 0 {
		cmd = lookup(commands, args[0])
	}
	if cmd == nil {
		printCommands(e)
		return exitWith(1)
	}
	return cmd.Run(e, args[1:])
}

func lookup(in []*Command, name string) *Command {
	for _, cmd := range in {
		if cmd.Name == name {
			return cmd
		}
	}
	return nil
}

// visible returns the commands users are told about.
func visible(in []*Command) []*Command {
	var shown []*Command
	for _, cmd := range in {
		if !cmd.Hidden {
			shown = append(shown, cmd)
		}
	}
	return shown
}

func names(in []*Command) []string {
	var named []string
	for _, cmd := range in {
		named = append(named, cmd.Name)
	}
	return named
}

// matching returns the names that start with prefix.
func matching(candidates []string, prefix string) []string {
	var matches []string
	for _, candidate := range candidates {
		if strings.HasPrefix(candidate, prefix) {
			matches = append(matches, candidate)
		}
	}
	return matches
}
