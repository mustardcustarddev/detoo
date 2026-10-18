package output

import (
	"fmt"
	"image/color"
	"io"
	"path"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"

	"github.com/mustardcustard/detoo/internal/features/agents/skills"
)

// SkillsProgress writes what the skills commands do as they do it, a line
// for each step: plain lines, with a mark in colour for each skill, and none
// of a report's boxes.
type SkillsProgress struct {
	w  io.Writer
	t  Terminal
	th theme

	// Updating says a skill installed was updated, rather than added.
	Updating bool
}

// NewSkillsProgress returns a SkillsProgress that writes to w, for the
// terminal t.
func NewSkillsProgress(w io.Writer, t Terminal) *SkillsProgress {
	return &SkillsProgress{w: &colorprofile.Writer{Forward: w, Profile: t.Profile}, t: t, th: newTheme(t.Dark)}
}

func (p *SkillsProgress) line(format string, a ...any) {
	fmt.Fprintf(p.w, format+"\n", a...)
}

func (p *SkillsProgress) fg(c color.Color, s string) string {
	return lipgloss.NewStyle().Foreground(c).Render(s)
}

func (p *SkillsProgress) tick() string  { return p.fg(p.th.green, "✓") }
func (p *SkillsProgress) cross() string { return p.fg(p.th.red, "✗") }
func (p *SkillsProgress) bang() string  { return p.fg(p.th.yellow, "!") }
func (p *SkillsProgress) dash() string  { return p.fg(p.th.faint, "–") }
func (p *SkillsProgress) faint(s string) string {
	return p.fg(p.th.faint, s)
}

// Fetching says where skills are being fetched from.
func (p *SkillsProgress) Fetching(src skills.Source) {
	switch src.Kind {
	case skills.GitHub:
		p.line("Fetching %s from GitHub…", bold(src.String()))
	default:
		p.line("Reading skills from %s…", bold(p.home(src.String())))
	}
}

// Fetched says how many skills the source offers.
func (p *SkillsProgress) Fetched(_ skills.Source, offered []skills.Skill) {
	p.line("%s", p.faint("Found "+count(len(offered), "skill", "skills")+"."))
}

// Installing says where the skills go, and which agents they are for.
func (p *SkillsProgress) Installing(t skills.Targets) {
	where := "Installing into " + bold(p.home(t.Dir))
	if len(t.Readers) > 0 {
		where += ", for " + and(titles(t.Readers))
	}
	p.line("%s", where)
	for _, l := range t.Links {
		p.line("Linking from %s, for %s", bold(p.home(l.Dir)), l.Agent.Title)
	}
}

// Installed marks a skill installed or updated.
func (p *SkillsProgress) Installed(name string) {
	if p.Updating {
		p.line("  %s %s %s", p.tick(), name, p.faint("updated"))
		return
	}
	p.line("  %s %s", p.tick(), name)
}

// UpToDate marks a skill that needed no update.
func (p *SkillsProgress) UpToDate(name string) {
	p.line("  %s %s %s", p.dash(), name, p.faint("up to date"))
}

// Gone marks a skill its source no longer offers.
func (p *SkillsProgress) Gone(name string, src skills.Source) {
	p.line("  %s %s %s", p.bang(), name, p.faint("is no longer in "+p.home(src.String())+", so is left as it is"))
}

// Failed marks a skill that could not be updated, and says why.
func (p *SkillsProgress) Failed(name string, err error) {
	p.line("  %s %s %s", p.cross(), name, p.faint(err.Error()))
}

// Removed marks a skill removed, and says from where.
func (p *SkillsProgress) Removed(name string, paths []string) {
	var dirs []string
	for _, at := range paths {
		dirs = append(dirs, p.home(path.Dir(at)))
	}
	p.line("  %s %s %s", p.tick(), name, p.faint("removed from "+strings.Join(dirs, ", ")))
}

// home shows a path inside the user's home directory under ~.
func (p *SkillsProgress) home(dir string) string {
	if p.t.Home != "" {
		if rest, ok := strings.CutPrefix(dir, p.t.Home); ok && (rest == "" || strings.HasPrefix(rest, "/")) {
			return "~" + rest
		}
	}
	return dir
}

// SkillsInstalled writes what an install installed.
func (p *SkillsProgress) SkillsInstalled(installed skills.Installed, global bool) {
	where := "in " + p.home(installed.Targets.Root)
	if global {
		where = "for every project"
	}
	n := len(installed.Installed)
	p.line("%s Installed %s %s. %s records %s.", p.tick(), count(n, "skill", "skills"), where, p.home(installed.Targets.LockFile), them(n))
}

// SkillsUpdated writes what an update found, and reports whether every skill
// it checked is up to date now.
func (p *SkillsProgress) SkillsUpdated(updated skills.Updated) bool {
	checked := len(updated.Updated) + len(updated.UpToDate) + len(updated.Gone) + len(updated.Failed)
	switch {
	case checked == 0:
		p.line("No skills are recorded in %s. Install them with detoo skills install.", p.home(updated.LockFile))
	case len(updated.Failed) > 0:
		p.line("%s %s could not be updated.", p.cross(), count(len(updated.Failed), "skill", "skills"))
	case len(updated.Updated) == 0 && len(updated.Gone) == 0:
		p.line("%s %s up to date.", p.tick(), all(len(updated.UpToDate)))
	default:
		summary := "Updated " + count(len(updated.Updated), "skill", "skills")
		if n := len(updated.UpToDate); n > 0 {
			summary += fmt.Sprintf("; %d %s up to date", n, was(n))
		}
		p.line("%s %s.", p.tick(), summary)
	}
	if n := len(updated.Gone); n > 0 {
		p.line("%s %s no longer in %s source. Remove %s with detoo skills remove.", p.bang(), are(n), its(n), them(n))
	}
	return len(updated.Failed) == 0
}

// SkillsRemoved writes what a remove removed.
func (p *SkillsProgress) SkillsRemoved(removed skills.Removed) {
	if len(removed.Removed) == 0 {
		p.line("No skills are recorded in %s, so none was removed.", p.home(removed.LockFile))
		return
	}
	p.line("%s Removed %s.", p.tick(), count(len(removed.Removed), "skill", "skills"))
	if n := len(removed.Emptied); n > 0 {
		var dirs []string
		for _, dir := range removed.Emptied {
			dirs = append(dirs, p.home(dir))
		}
		p.line("%s Removed %s, which %s empty now.", p.tick(), and(dirs), isAre(n))
	}
}

func isAre(n int) string {
	if n == 1 {
		return "is"
	}
	return "are"
}

// SkillsAvailable writes the skills a source offers, a line each: its name,
// and as much of its description as fits.
func SkillsAvailable(w io.Writer, t Terminal, offered []skills.Skill) error {
	width := t.Width
	if width <= 0 {
		width = 80
	}
	widest := 0
	for _, s := range offered {
		widest = max(widest, len(s.Name))
	}
	out := &colorprofile.Writer{Forward: w, Profile: t.Profile}
	for _, s := range offered {
		line := bold(s.Name)
		if s.Description != "" {
			room := width - widest - 2
			line += strings.Repeat(" ", widest-len(s.Name)+2) + truncate(s.Description, room)
		}
		if _, err := fmt.Fprintln(out, strings.TrimRight(line, " ")); err != nil {
			return err
		}
	}
	return nil
}

// truncate shortens text to n characters, ending in … when it is cut.
func truncate(text string, n int) string {
	runes := []rune(text)
	if n <= 0 {
		return ""
	}
	if len(runes) <= n {
		return text
	}
	return strings.TrimRight(string(runes[:n-1]), " ") + "…"
}

func titles(agents []skills.Agent) []string {
	var names []string
	for _, a := range agents {
		names = append(names, a.Title)
	}
	return names
}

// and joins words as a list in prose: "Codex", "Codex and Cursor", "Codex,
// Cursor and Amp".
func and(words []string) string {
	if len(words) < 2 {
		return strings.Join(words, "")
	}
	return strings.Join(words[:len(words)-1], ", ") + " and " + words[len(words)-1]
}

// all says how many skills: "The skill is", "All 8 skills are".
func all(n int) string {
	if n == 1 {
		return "The skill is"
	}
	return fmt.Sprintf("All %d skills are", n)
}

func was(n int) string {
	if n == 1 {
		return "was"
	}
	return "were"
}

func are(n int) string {
	if n == 1 {
		return "1 skill is"
	}
	return fmt.Sprintf("%d skills are", n)
}

func its(n int) string {
	if n == 1 {
		return "its"
	}
	return "their"
}
