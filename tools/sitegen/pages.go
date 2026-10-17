package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

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
	"github.com/mustardcustard/detoo/internal/pse"
	"github.com/mustardcustard/detoo/internal/rfc"
	"github.com/mustardcustard/detoo/internal/sty"
)

// Page is one page of the site's command reference: what a tool, or one of
// its commands, prints as its help.
type Page struct {
	Tool        string // the tool, such as "adr"
	Command     string // the command, such as "new", or "" for the tool itself
	Title       string // "adr new", or "adr"
	Description string // one line, for the page's meta description
	Order       int    // where the page goes among the tool's pages
	Help        string // exactly what `<tool> help <command>` prints
}

// workflow is the order a reader meets the commands in, which is the order
// the sidebar lists them.  A command not named here goes after these, in
// alphabetical order, so that a new command is never left out.
var workflow = []string{"init", "new", "link", "list", "graph", "config", "doctor"}

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
		p := project.New(os.TempDir(), detoo.Kinds)
		application := app.NewDetoo(detoo.Kinds, p, project.NewGit(p))
		terminal := func() output.Terminal { return output.Terminal{} }
		return cli.NewDetooCommand(detoo.Files, application, terminal, nil, stdout, io.Discard, none)
	}}
	return []tool{
		recordTool(adr.Kind, adr.Files),
		recordTool(rfc.Kind, rfc.Files),
		recordTool(sty.Kind, sty.Files),
		recordTool(pse.Kind, pse.Files),
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
		help, err := t.help()
		if err != nil {
			return nil, err
		}
		pages = append(pages, Page{
			Tool:        t.name,
			Title:       t.name,
			Description: firstSentence(root.Long),
			Help:        help,
		})

		var commands []*cobra.Command
		for _, cmd := range root.Commands() {
			if documented(cmd) {
				commands = append(commands, cmd)
			}
		}
		slices.SortStableFunc(commands, func(a, b *cobra.Command) int {
			return rank(a.Name()) - rank(b.Name())
		})
		for i, cmd := range commands {
			help, err := t.help(cmd.Name())
			if err != nil {
				return nil, err
			}
			pages = append(pages, Page{
				Tool:        t.name,
				Command:     cmd.Name(),
				Title:       t.name + " " + cmd.Name(),
				Description: cmd.Short,
				Order:       i + 1,
				Help:        help,
			})
		}
	}
	return pages, nil
}

// help runs `<tool> help <command>` on a fresh command tree, the same code a
// user's terminal runs, and returns what it prints.
func (t tool) help(command ...string) (string, error) {
	var out bytes.Buffer
	root := t.root(&out)
	root.SetArgs(append([]string{"help"}, command...))
	if err := root.Execute(); err != nil {
		return "", fmt.Errorf("%s help %s: %w", t.name, strings.Join(command, " "), err)
	}
	return out.String(), nil
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

// Path returns where the page is written, relative to the reference
// collection: the tool's overview is its index.
func (p Page) Path() string {
	if p.Command == "" {
		return p.Tool + "/index.md"
	}
	return p.Tool + "/" + p.Command + ".md"
}

// Markdown returns the page as the site reads it: frontmatter, then the help
// in a fenced block, so that it shows exactly as it does in a terminal.
//
// Help is prose with code in backticks, so the fence is made longer than any
// run of backticks in it, and the frontmatter's values are quoted, since a
// description may hold a colon.  Go's quoting is valid YAML.
func (p Page) Markdown() string {
	fence := strings.Repeat("`", max(3, longestRun(p.Help, '`')+1))
	var b strings.Builder
	fmt.Fprintf(&b, "---\n")
	fmt.Fprintf(&b, "tool: %s\n", strconv.Quote(p.Tool))
	fmt.Fprintf(&b, "command: %s\n", strconv.Quote(p.Command))
	fmt.Fprintf(&b, "title: %s\n", strconv.Quote(p.Title))
	fmt.Fprintf(&b, "description: %s\n", strconv.Quote(p.Description))
	fmt.Fprintf(&b, "order: %d\n", p.Order)
	fmt.Fprintf(&b, "---\n\n")
	fmt.Fprintf(&b, "%stext\n%s\n%s\n", fence, strings.TrimRight(p.Help, "\n"), fence)
	return b.String()
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
