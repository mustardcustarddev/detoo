package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/mustardcustard/detoo/internal/adapters/clock"
	"github.com/mustardcustard/detoo/internal/adapters/editor"
	"github.com/mustardcustard/detoo/internal/adapters/memlog"
	"github.com/mustardcustard/detoo/internal/adapters/project"
	"github.com/mustardcustard/detoo/internal/adapters/templates"
	"github.com/mustardcustard/detoo/internal/adr"
	"github.com/mustardcustard/detoo/internal/app"
	"github.com/mustardcustard/detoo/internal/cli"
	"github.com/mustardcustard/detoo/internal/cli/output"
	"github.com/mustardcustard/detoo/internal/detoo"
	"github.com/mustardcustard/detoo/internal/domain/record"
	"github.com/mustardcustard/detoo/internal/features/agents/skills"
	"github.com/mustardcustard/detoo/internal/fyr"
	"github.com/mustardcustard/detoo/internal/plm"
	"github.com/mustardcustard/detoo/internal/rfc"
	"github.com/mustardcustard/detoo/internal/sop"
	"github.com/mustardcustard/detoo/internal/ust"
)

// Page is one page of the site's command reference: what a tool, or one of
// its commands, says of itself in its help, as Markdown.
type Page struct {
	Tool        string // the tool, such as "adr"
	Command     string // the command, such as "new", or "" for the tool itself
	Title       string // "adr new", or "adr"
	Description string // one line, for the page's meta description
	Order       int    // where the page goes among the tool's pages
	Body        string // the help as Markdown: its prose, usage, examples, commands and options
}

// workflow is the order a reader meets the commands in, which is the order
// the sidebar lists them.  A command not named here goes after these, in
// alphabetical order, so that a new command is never left out.
var workflow = []string{"init", "new", "link", "list", "graph", "config", "doctor", "skills", "install", "update", "remove"}

// tool is a command-line tool the site documents, and a way to build its
// command tree, writing to stdout.
type tool struct {
	name string
	root func(stdout io.Writer) *cobra.Command
}

// tools returns the tools, in the order the site lists them.
//
// Each command tree is built as package tool builds it, but over adapters that
// touch nothing: help never reaches the application, so the adapters are there
// only because the constructors ask for them.  getenv finds nothing, so no
// pager is run.
func tools() []tool {
	none := func(string) string { return "" }
	recordTool := func(kind record.Kind, files fs.FS) tool {
		return tool{kind.Name, func(stdout io.Writer) *cobra.Command {
			application := app.New(
				kind,
				memlog.Of(kind),
				templates.New(os.TempDir(), none, kind, files),
				clock.New(none, kind),
				editor.New(os.TempDir(), none, strings.NewReader(""), io.Discard, io.Discard),
			)
			return cli.NewRootCommand(kind, files, application, stdout, io.Discard, none)
		}}
	}
	doctor := tool{"detoo", func(stdout io.Writer) *cobra.Command {
		p := project.New(os.TempDir())
		application := app.NewDetoo(detoo.Kinds, p, project.NewGit(p), nil, skills.Installations{})
		terminal := func() output.Terminal { return output.Terminal{} }
		return cli.NewDetooCommand(detoo.Files, application, terminal, nil, stdout, io.Discard, none)
	}}
	return []tool{
		recordTool(adr.Kind, adr.Files),
		recordTool(rfc.Kind, rfc.Files),
		recordTool(ust.Kind, ust.Files),
		recordTool(plm.Kind, plm.Files),
		recordTool(sop.Kind, sop.Files),
		recordTool(fyr.Kind, fyr.Files),
		doctor,
	}
}

// documented reports whether a command gets a page.  Cobra adds help and
// completion to every tool, and they say nothing about the tool itself.
func documented(cmd *cobra.Command) bool {
	return !cmd.Hidden && cmd.Name() != "help" && cmd.Name() != "completion"
}

// Pages returns a page for each tool, and for each of its commands, in the
// order the site lists them.
func Pages() ([]Page, error) {
	var pages []Page
	for _, t := range tools() {
		root := t.root(io.Discard)
		description := firstSentence(root.Long)
		pages = append(pages, Page{
			Tool:        t.name,
			Title:       t.name,
			Description: description,
			Body:        withoutLead(body(root), description),
		})

		for i, cmd := range commandsUnder(root) {
			// The command's path below the tool, such as "skills install".
			command := strings.TrimPrefix(cmd.CommandPath(), root.Name()+" ")
			pages = append(pages, Page{
				Tool:        t.name,
				Command:     command,
				Title:       t.name + " " + command,
				Description: cmd.Short,
				Order:       i + 1,
				Body:        body(cmd),
			})
		}
	}
	return pages, nil
}

// commandsUnder returns the documented commands beneath cmd, in the order of
// the workflow, each followed by those beneath it.
func commandsUnder(cmd *cobra.Command) []*cobra.Command {
	var commands []*cobra.Command
	for _, sub := range cmd.Commands() {
		if documented(sub) {
			commands = append(commands, sub)
		}
	}
	slices.SortStableFunc(commands, func(a, b *cobra.Command) int {
		return rank(a.Name()) - rank(b.Name())
	})
	var all []*cobra.Command
	for _, sub := range commands {
		all = append(all, sub)
		all = append(all, commandsUnder(sub)...)
	}
	return all
}

// rank places a command in the workflow; Cobra has already sorted the
// commands alphabetically, so a stable sort keeps the others in that order.
func rank(name string) int {
	if i := slices.Index(workflow, name); i >= 0 {
		return i
	}
	return len(workflow)
}

// firstSentence returns the first sentence of text, on one line.
func firstSentence(text string) string {
	text = strings.Join(strings.Fields(text), " ")
	if i := strings.Index(text, ". "); i >= 0 {
		return text[:i+1]
	}
	return text
}

// withoutLead returns body without its first paragraph, if that paragraph is
// lead: the page shows its description beneath its title, so it need not say
// it twice.
func withoutLead(body, lead string) string {
	first, rest, ok := strings.Cut(body, "\n\n")
	if ok && strings.Join(strings.Fields(first), " ") == lead {
		return rest
	}
	return body
}

// Path returns where the page is written, relative to the reference
// collection: the tool's overview is its index, and a command beneath
// another is in a directory named after it, as detoo/skills/install.md.
func (p Page) Path() string {
	if p.Command == "" {
		return p.Tool + "/index.md"
	}
	return p.Tool + "/" + strings.ReplaceAll(p.Command, " ", "/") + ".md"
}

// Markdown returns the page as the site reads it: frontmatter, then the body.
//
// The frontmatter's values are quoted, since a description may hold a colon.
// Go's quoting is valid YAML.
func (p Page) Markdown() string {
	var b strings.Builder
	fmt.Fprintf(&b, "---\n")
	fmt.Fprintf(&b, "tool: %s\n", strconv.Quote(p.Tool))
	fmt.Fprintf(&b, "command: %s\n", strconv.Quote(p.Command))
	fmt.Fprintf(&b, "title: %s\n", strconv.Quote(p.Title))
	fmt.Fprintf(&b, "description: %s\n", strconv.Quote(p.Description))
	fmt.Fprintf(&b, "order: %d\n", p.Order)
	fmt.Fprintf(&b, "---\n\n")
	b.WriteString(p.Body)
	return b.String()
}

// body returns what cmd's help says, as Markdown: the same parts, from the
// same fields of the command, that Cobra prints in a terminal.  The prose in
// the help's ".md" files is Markdown already, so it is used as it is; the
// usage and examples are commands, so they are code; and the commands and
// options are tables.  The --help option, and the line saying how to get help
// on a command, are left out: a page is the help.
func body(cmd *cobra.Command) string {
	var b strings.Builder
	if long := strings.TrimSpace(cmd.Long); long != "" {
		fmt.Fprintf(&b, "%s\n\n", long)
	}

	var options []string
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if f.Hidden || f.Name == "help" {
			return
		}
		options = append(options, fmt.Sprintf("| %s | %s |", cell(option(f)), cell(optionUsage(f))))
	})

	// A command with commands beneath it, and no options but --help, runs only
	// to print its help, so it is used only with one of its commands.
	var usage []string
	if cmd.Runnable() && (!cmd.HasAvailableSubCommands() || len(options) > 0) {
		usage = append(usage, cmd.UseLine())
	}
	if cmd.HasAvailableSubCommands() {
		usage = append(usage, cmd.CommandPath()+" [command]")
	}
	fmt.Fprintf(&b, "## Usage\n\n%s\n\n", codeBlock(strings.Join(usage, "\n")))

	if cmd.HasExample() {
		fmt.Fprintf(&b, "## Examples\n\n%s\n\n", codeBlock(dedent(cmd.Example)))
	}

	var commands []string
	for _, sub := range cmd.Commands() {
		if documented(sub) {
			// Each command's page is in a directory named after its parent's
			// page, so the link is relative to this one.
			commands = append(commands, fmt.Sprintf("| [`%s`](%s/) | %s |", sub.Name(), sub.Name(), cell(sub.Short)))
		}
	}
	if len(commands) > 0 {
		fmt.Fprintf(&b, "## Commands\n\n| Command | Description |\n| --- | --- |\n%s\n\n", strings.Join(commands, "\n"))
	}

	if len(options) > 0 {
		fmt.Fprintf(&b, "## Options\n\n| Option | Description |\n| --- | --- |\n%s\n\n", strings.Join(options, "\n"))
	}

	return strings.TrimRight(b.String(), "\n") + "\n"
}

// option returns how an option is written, long form first, such as
// "`--supersede REFERENCE`, `-s`".
func option(f *pflag.Flag) string {
	name, _ := pflag.UnquoteUsage(f)
	long := "--" + f.Name
	if name != "" {
		long += " " + name
	}
	s := "`" + long + "`"
	if f.Shorthand != "" {
		s += ", `-" + f.Shorthand + "`"
	}
	return s
}

// optionUsage returns what an option does, as Cobra describes it, with its
// default when it has one worth saying.
func optionUsage(f *pflag.Flag) string {
	_, usage := pflag.UnquoteUsage(f)
	switch f.DefValue {
	case "", "false", "[]", "0":
	default:
		usage += fmt.Sprintf(" (default `%s`)", f.DefValue)
	}
	return usage
}

// cell makes text safe to put in a cell of a Markdown table.
func cell(text string) string {
	return strings.ReplaceAll(strings.Join(strings.Fields(text), " "), "|", `\|`)
}

// codeBlock fences text as shell commands.  Help is prose with code in
// backticks, so the fence is made longer than any run of backticks in it.
func codeBlock(text string) string {
	fence := strings.Repeat("`", max(3, longestRun(text, '`')+1))
	return fence + "bash\n" + strings.TrimRight(text, "\n") + "\n" + fence
}

// dedent removes the indentation that every line of text shares, which the
// ".example" files carry so that the examples are indented in a terminal.
func dedent(text string) string {
	lines := strings.Split(strings.Trim(text, "\n"), "\n")
	indent := -1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		n := len(line) - len(strings.TrimLeft(line, " \t"))
		if indent < 0 || n < indent {
			indent = n
		}
	}
	for i, line := range lines {
		if len(line) >= indent {
			lines[i] = line[max(indent, 0):]
		} else {
			lines[i] = strings.TrimLeft(line, " \t")
		}
	}
	return strings.Join(lines, "\n")
}

// longestRun returns the length of the longest run of c in s.
func longestRun(s string, c byte) int {
	longest, run := 0, 0
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			run++
			longest = max(longest, run)
		} else {
			run = 0
		}
	}
	return longest
}
