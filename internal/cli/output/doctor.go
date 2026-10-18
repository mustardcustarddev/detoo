package output

import (
	"fmt"
	"image/color"
	"io"
	"path"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"

	"github.com/mustardcustard/detoo/internal/domain/record"
	"github.com/mustardcustard/detoo/internal/features/health/doctor"
)

// Terminal is what a report knows of the terminal it is written to.  The
// caller finds it out; nothing here looks at the terminal itself.
type Terminal struct {
	// Width is the number of columns, or 0 when it is not known.
	Width int

	// Dark reports whether the background is dark, which chooses between the
	// dark and light palettes.
	Dark bool

	// Profile is the colours the terminal shows.  Output to a file or a pipe
	// has none, and is written as plain text.
	Profile colorprofile.Profile

	// Home is the user's home directory, which a path inside is shown under
	// as ~.
	Home string

	// Version is the version of detoo.
	Version string
}

// widest is the most columns a report spreads across, however wide the
// terminal.
const widest = 120

// Doctor writes the health of a project's logs: a header, a box for each log
// found when any has problems (or one box listing every log when none has),
// and a status line ending in the exit status the command returns.
func Doctor(w io.Writer, t Terminal, result doctor.Result) error {
	return doctorReport(w, t, result, false)
}

// doctorReport writes the report Doctor writes, after a fix that changed
// nothing when fixed is true, which hints at what is left to fix.
func doctorReport(w io.Writer, t Terminal, result doctor.Result, fixed bool) error {
	r := newReport(t, result)
	if result.Healthy() {
		return r.write(w, r.healthyBox(), r.statusBar())
	}
	boxes := make([]string, 0, len(result.Logs))
	var missing []string
	for _, l := range result.Logs {
		if l.Found() {
			boxes = append(boxes, r.logBox(l))
		} else {
			missing = append(missing, l.Kind.Name)
		}
	}
	boxes = append(boxes, r.missing(missing)...)
	hint := r.fixHint()
	if fixed {
		hint = r.renumberHint()
	}
	if hint != "" {
		boxes = append(boxes, hint)
	}
	return r.write(w, strings.Join(boxes, "\n\n"), r.statusBar())
}

// renumberHint says that doctor --fix --renumber renumbers documents that
// share a number, when a fix has left some, or is empty.
func (r *report) renumberHint() string {
	if r.result.Renumberable() == 0 {
		return ""
	}
	return r.fg(r.th.primary).Render("→") + " Run " + bold("detoo doctor --fix --renumber") + " " +
		r.fg(r.th.subtext).Render("to renumber the documents that share a number.")
}

// fixHint says that doctor --fix fixes the problems that do not need fixing
// by hand, and how many of them that is, or is empty when every problem does.
func (r *report) fixHint() string {
	problems, byHand, fixable := r.result.Problems(), r.result.ByHand(), r.result.Fixable()
	if fixable == 0 {
		return ""
	}
	text := "to fix them."
	switch {
	case problems == 1:
		text = "to fix it."
	case byHand > 0:
		verb := "needs"
		if byHand > 1 {
			verb = "need"
		}
		text = fmt.Sprintf("to fix %d of the %d problems; %d %s fixing by hand.", fixable, problems, byHand, verb)
	}
	return r.fg(r.th.primary).Render("→") + " Run " + bold("detoo doctor --fix") + " " + r.fg(r.th.subtext).Render(text)
}

// DoctorFix writes what a fix changed: a header, a box for each log found,
// saying what was fixed in it and what is left to fix by hand, a box of the
// documents whose links were changed, and a status line ending in the exit
// status the command returns.  When nothing was changed, it is the report
// doctor writes without fixing.
func DoctorFix(w io.Writer, t Terminal, fixed doctor.Fixed) error {
	if fixed.Changes() == 0 {
		return doctorReport(w, t, fixed.After, true)
	}
	r := newReport(t, fixed.After)
	boxes := make([]string, 0, len(fixed.After.Logs)+1)
	var missing []string
	for i, l := range fixed.After.Logs {
		f := fixed.Logs[i]
		switch {
		case f.Count() > 0 || l.Problems() > 0:
			boxes = append(boxes, r.fixBox(l, f, removedBy(fixed, i)))
		case l.Found():
			boxes = append(boxes, r.logBox(l))
		default:
			missing = append(missing, l.Kind.Name)
		}
	}
	if len(fixed.Relinked) > 0 {
		boxes = append(boxes, r.linksBox(fixed.Relinked))
	}
	boxes = append(boxes, r.missing(missing)...)
	if fixed.Git {
		boxes = append(boxes, r.tick()+" Used "+bold("git mv")+r.fg(r.th.subtext).Render(", so git records the moves and renames. Review with git status."))
	}
	if hint := r.renumberHint(); hint != "" {
		boxes = append(boxes, hint)
	}
	return r.write(w, strings.Join(boxes, "\n\n"), r.fixStatusBar(fixed))
}

// write writes the header, the body and the status bar, in the colours the
// terminal shows.
func (r *report) write(w io.Writer, body, status string) error {
	text := r.header() + "\n\n" + body + "\n\n" + status + "\n"
	_, err := io.WriteString(&colorprofile.Writer{Forward: w, Profile: r.t.Profile}, text)
	return err
}

// missing says which logs the project does not keep, if any.
func (r *report) missing(names []string) []string {
	if len(names) == 0 {
		return nil
	}
	return []string{r.dash() + " " + r.fg(r.th.subtext).Render("No "+series(names)+" log found.")}
}

// removedBy returns the directories removed after the moves out of them that
// the log of index i made, each given to the last log that moved out of it.
func removedBy(fixed doctor.Fixed, i int) []string {
	var dirs []string
	for _, dir := range fixed.Removed {
		last := -1
		for j, f := range fixed.Logs {
			for _, m := range f.Moved {
				if path.Dir(m.From) == dir {
					last = j
				}
			}
		}
		if last == i {
			dirs = append(dirs, dir)
		}
	}
	return dirs
}

// theme is the palette the report is drawn in, taken from detoo's avatar: its
// orange shoulder, amber set square, teal joints, cream body and charcoal
// visor.  Each colour has a shade for a light background and one for a dark
// background.  Plain text is left in the terminal's own colour.
type theme struct {
	subtext, faint, border, bar, base color.Color
	primary, red, yellow, green       color.Color
}

func newTheme(dark bool) theme {
	pick := lipgloss.LightDark(dark)
	c := func(light, dark string) color.Color {
		return pick(lipgloss.Color(light), lipgloss.Color(dark))
	}
	return theme{
		subtext: c("#6b6157", "#cbbfa9"),
		faint:   c("#8f857a", "#958b7d"),
		border:  c("#8aa39e", "#4d6662"),
		bar:     c("#d6e0dc", "#314442"),
		base:    c("#faf2e2", "#2f2e2d"),
		primary: c("#cf6418", "#ee8f49"),
		red:     c("#c4372f", "#f07c6c"),
		yellow:  c("#a87500", "#fbc04d"),
		green:   c("#3d8350", "#9ccf8f"),
	}
}

// report draws one result for one terminal.
type report struct {
	t      Terminal
	th     theme
	width  int
	result doctor.Result
}

func newReport(t Terminal, result doctor.Result) *report {
	width := t.Width
	if width <= 0 {
		width = 80
	}
	return &report{t: t, th: newTheme(t.Dark), width: min(width, widest), result: result}
}

func (r *report) fg(c color.Color) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(c)
}

func bold(s string) string {
	return lipgloss.NewStyle().Bold(true).Render(s)
}

func (r *report) badge(text string, bg color.Color) string {
	return lipgloss.NewStyle().Background(bg).Foreground(r.th.base).Bold(true).Padding(0, 1).Render(text)
}

func (r *report) cross() string { return r.fg(r.th.red).Render("✗") }
func (r *report) bang() string  { return r.fg(r.th.yellow).Bold(true).Render("!") }
func (r *report) tick() string  { return r.fg(r.th.green).Render("✓") }
func (r *report) dash() string  { return r.fg(r.th.faint).Render("–") }

// header is the mark beside what was checked.
func (r *report) header() string {
	label := lipgloss.NewStyle().Foreground(r.th.subtext).Width(10)
	names := make([]string, 0, len(r.result.Logs))
	for _, l := range r.result.Logs {
		names = append(names, l.Kind.Name)
	}
	info := strings.Join([]string{
		r.badge("detoo", r.th.primary) + " " + r.fg(r.th.faint).Render(r.t.Version),
		"",
		bold("doctor") + " " + r.fg(r.th.subtext).Render("helps maintain and troubleshoot decision record logs"),
		"",
		label.Render("project") + r.home(r.result.Root),
		label.Render("logs") + strings.Join(names, "  "),
		label.Render("checked") + fmt.Sprintf("%s in %s",
			count(r.result.Documents(), "document", "documents"),
			count(r.result.Directories(), "directory", "directories")),
	}, "\n")

	const gap = "    "
	if markWidth+len(gap)+lipgloss.Width(info) > r.width {
		return info
	}
	// Joining pads every line to the longest; the padding is dropped again.
	lines := strings.Split(lipgloss.JoinHorizontal(lipgloss.Bottom, r.mark(), gap, info), "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " ")
	}
	return strings.Join(lines, "\n")
}

// home shows a path inside the user's home directory under ~.
func (r *report) home(dir string) string {
	if r.t.Home != "" {
		if rest, ok := strings.CutPrefix(dir, r.t.Home); ok && (rest == "" || strings.HasPrefix(rest, "/")) {
			return "~" + rest
		}
	}
	return dir
}

// Two vertical pixels share a character, so the portrait occupies markWidth
// columns and markHeight/2 lines.
const (
	markWidth  = 17
	markHeight = 14
)

// markPixels is a front-facing portrait of the mascot: cream helmet, dark
// visor and orange shoulders. Each character is one pixel.
var markPixels = [markHeight]string{
	"                 ",
	"  ccccccccccccc  ",
	" cwvvvvvvvvvvvwc ",
	" cvvvvvvvvvvvvvc ",
	" cvgvvvvvvvvvvvc ",
	" cvgvvvvvvvvvvvc ",
	" cvvvvvvvvvvvvvc ",
	" cvvvvvvvvvvvvvc ",
	" cvvvvvvvvvvvvvc ",
	" ccvvvvvvvvvvvcc ",
	"  ccccccccccccc  ",
	"      ttttt      ",
	"  ooocccccccooo  ",
	" oooo ccccc oooo ",
}

var markPalette = map[byte]string{
	'c': "#d6cbb0", // cream shell in shadow
	'w': "#f5eedb", // lit shell
	'v': "#303436", // visor
	'g': "#5a6262", // visor reflection
	't': "#386363", // teal neck joint
	'o': "#e98729", // orange shoulders
}

// markPixel returns the colour of a portrait pixel, or "" for empty space.
func markPixel(row, col int) string {
	return markPalette[markPixels[row][col]]
}

// mark draws the detoo mark in half blocks.  Without colour, a character
// holding two pixels of different colours is drawn as a full block, so that
// the shape survives.
func (r *report) mark() string {
	colour := r.t.Profile > colorprofile.ASCII
	lines := make([]string, 0, markHeight/2)
	for row := 0; row < markHeight; row += 2 {
		var line strings.Builder
		for col := range markWidth {
			top, bottom := markPixel(row, col), markPixel(row+1, col)
			// Leave the visor open without colour so its outline remains visible.
			if !colour {
				if p := markPixels[row][col]; p == 'v' || p == 'g' {
					top = ""
				}
				if p := markPixels[row+1][col]; p == 'v' || p == 'g' {
					bottom = ""
				}
			}
			line.WriteString(halfBlock(top, bottom, colour))
		}
		lines = append(lines, line.String())
	}
	return strings.Join(lines, "\n")
}

func halfBlock(top, bottom string, colour bool) string {
	style := lipgloss.NewStyle()
	switch {
	case top == "" && bottom == "":
		return " "
	case bottom == "":
		return style.Foreground(lipgloss.Color(top)).Render("▀")
	case top == "":
		return style.Foreground(lipgloss.Color(bottom)).Render("▄")
	case top == bottom || !colour:
		return style.Foreground(lipgloss.Color(top)).Render("█")
	default:
		return style.Foreground(lipgloss.Color(top)).Background(lipgloss.Color(bottom)).Render("▀")
	}
}

// box draws body in a rounded border the width of the report, with title set
// into the top edge on the left and aside on the right.
func (r *report) box(title, aside, body string, border color.Color) string {
	edge := r.fg(border)
	b := lipgloss.RoundedBorder()
	sides := lipgloss.NewStyle().Border(b, false, true, true, true).BorderForeground(border).Padding(0, 1).Width(r.width)

	head := edge.Render(b.TopLeft+b.Top) + " " + title + " "
	tail := " " + aside + " " + edge.Render(b.Top+b.TopRight)
	fill := r.width - lipgloss.Width(head) - lipgloss.Width(tail)
	if fill < 1 {
		tail = edge.Render(b.TopRight)
		fill = r.width - lipgloss.Width(head) - lipgloss.Width(tail)
	}
	head += edge.Render(strings.Repeat(b.Top, max(fill, 1))) + tail
	return head + "\n" + sides.Render(body)
}

// logBox draws the problems with one log.
func (r *report) logBox(l doctor.Log) string {
	border, title := r.logTitle(l)
	return r.box(title, r.logAside(l, 0), r.logBody(l), border)
}

// logTitle is the border and the title of a log's box, which say whether
// anything is wrong with it, and whether it needs fixing by hand.
func (r *report) logTitle(l doctor.Log) (color.Color, string) {
	border, icon := r.th.border, r.tick()
	switch {
	case l.ByHand() > 0:
		border, icon = r.th.yellow, r.bang()
	case l.Problems() > 0:
		icon = r.cross()
	}

	dir := l.Dir
	if len(l.Misplaced) > 0 {
		dir = l.Misplaced[0].From
	}
	title := icon + " " + bold(l.Kind.Name) + "  " + r.fg(r.th.subtext).Render(dir)
	if l.ByHand() > 0 {
		title += "  " + r.badge("BY HAND", r.th.yellow)
	}
	return border, title
}

// logAside counts a log's documents, the problems fixed in it, and the
// problems with it.
func (r *report) logAside(l doctor.Log, fixed int) string {
	aside := count(l.Documents(), l.Kind.Noun, l.Kind.Plural)
	if fixed > 0 {
		aside += fmt.Sprintf(" · %d fixed", fixed)
	}
	if l.Problems() > 0 {
		aside += " · " + count(l.Problems(), "problem", "problems")
	}
	return r.fg(r.th.faint).Render(aside)
}

func (r *report) logBody(l doctor.Log) string {
	subtext := r.fg(r.th.subtext)
	if l.Problems() == 0 {
		text := "In place."
		if l.Kind.Prefix != "" {
			text += " File names begin with " + l.Kind.Prefix + "-."
		}
		return r.tick() + " " + subtext.Render(text)
	}

	var lines []string
	for _, m := range l.FormerMarkers {
		lines = append(lines, r.formerMarker(l, m)...)
	}
	for _, m := range l.Misplaced {
		lines = append(lines, r.misplaced(l, m)...)
	}
	if len(l.Misnamed) > 0 {
		lines = append(lines, r.misnamed(l)...)
	}
	if len(l.Duplicates) > 0 {
		lines = append(lines, r.duplicates(l)...)
	}
	if len(l.Mistitled) > 0 {
		lines = append(lines, r.mistitled(l)...)
	}
	return strings.Join(lines, "\n")
}

// formerMarker describes a marker file named after a former name of the tool,
// which it no longer reads, and how to fix it.
func (r *report) formerMarker(l doctor.Log, m doctor.MarkerRename) []string {
	subtext := r.fg(r.th.subtext)
	name := l.Kind.Name
	lines := []string{r.cross() + " Marker file has a former name: " + name + " reads " + bold(m.To)}
	if !m.Clash {
		lines = append(lines, "  "+subtext.Render(fmt.Sprintf("%s no longer reads %s, which names %s.", name, m.From, m.Dir)))
		return append(lines, r.mapping([][]string{{m.From, r.arrow(), r.fg(r.th.green).Bold(true).Render(m.To)}})...)
	}
	lines = append(lines, r.bang()+" Needs fixing by hand: both marker files exist, so renaming would clash.")
	return append(lines, grid([][]string{
		{m.From, subtext.Render("names " + m.Dir), subtext.Render("· former name")},
		{m.To, subtext.Render("names " + l.Dir), subtext.Render("· current name, which " + name + " reads")},
	})...)
}

// misplaced describes a log kept in the wrong directory, and how to fix it.
func (r *report) misplaced(l doctor.Log, m doctor.Move) []string {
	subtext := r.fg(r.th.subtext)
	name := l.Kind.Name
	var lines []string
	switch m.Reason {
	case doctor.UnderDoc:
		lines = append(lines, r.cross()+" Kept under "+bold("doc/")+" instead of "+bold("docs/"))
	case doctor.Former:
		lines = append(lines, r.cross()+" Kept in a former directory: "+l.Kind.Plural+" now live in "+bold(m.To))
		if m.Clash != doctor.MarkerNamesIt && m.Documents > 0 {
			lines = append(lines, "  "+subtext.Render(fmt.Sprintf("%s no longer finds %s in %s.",
				name, the(m.Documents, l.Kind.Noun, l.Kind.Plural), m.From)))
		}
	}

	switch m.Clash {
	case doctor.NoClash:
		lines = append(lines, r.mapping([][]string{{r.from(m), r.arrow(), r.to(m)}})...)
	case doctor.TargetExists:
		role := "former directory"
		if m.Reason == doctor.UnderDoc {
			role = "under doc/"
		}
		lines = append(lines, r.bang()+" Needs fixing by hand: both directories exist, so moving one would clash.")
		lines = append(lines, grid([][]string{
			{m.From, subtext.Render(count(m.Documents, l.Kind.Noun, l.Kind.Plural)), subtext.Render("· " + role)},
			{m.To, subtext.Render(count(l.DocumentsIn(m.To), l.Kind.Noun, l.Kind.Plural)), subtext.Render("· current directory, where " + name + " looks")},
		})...)
	case doctor.MarkerNamesIt:
		lines = append(lines, r.bang()+" Needs fixing by hand: "+bold(l.Marker)+" names "+m.From+", so "+name+" still looks there.")
	case doctor.SharedTarget:
		lines = append(lines, r.bang()+" Needs fixing by hand: "+m.ClashWith+" moves to "+m.To+" too, so moving both would clash.")
	}
	return lines
}

func (r *report) arrow() string { return r.fg(r.th.faint).Render("→") }

// from shows a misplaced directory, with what is wrong with it marked.
func (r *report) from(m doctor.Move) string {
	if rest, ok := strings.CutPrefix(m.From, "doc/"); ok && m.Reason == doctor.UnderDoc {
		return r.fg(r.th.red).Render("doc/") + rest
	}
	return m.From
}

// to shows where a misplaced directory belongs, with what changes marked.
func (r *report) to(m doctor.Move) string {
	green := r.fg(r.th.green).Bold(true)
	subtext := r.fg(r.th.subtext)
	if rest, ok := strings.CutPrefix(m.To, "docs/"); ok && m.Reason == doctor.UnderDoc {
		return green.Render("docs/") + subtext.Render(rest)
	}
	parent, base := path.Split(m.To)
	return subtext.Render(parent) + green.Render(base)
}

// misnamed lists the documents not named with the log's prefix, and the names
// they should have.  When the log is kept in more than one directory, each
// document is shown with its directory.
func (r *report) misnamed(l doctor.Log) []string {
	prefix := l.Kind.Prefix + "-"
	lines := []string{r.cross() + " File names do not begin with " + bold(prefix)}
	rows := make([][]string, 0, len(l.Misnamed))
	for _, rn := range l.Misnamed {
		rest := strings.TrimPrefix(rn.To, prefix)
		lead := strings.TrimSuffix(rn.From, rest)
		from := r.fg(r.th.red).Render(lead) + rest
		if len(l.Dirs) > 1 {
			from = r.fg(r.th.faint).Render(rn.Dir+"/") + from
		}
		if rn.Clash {
			why := rn.To + " is there already"
			if rn.ClashWith != rn.To {
				why = rn.ClashWith + " would have the same name"
			}
			rows = append(rows, []string{from, r.bang(), r.fg(r.th.yellow).Render(why + ": fix by hand")})
			continue
		}
		rows = append(rows, []string{from, r.arrow(), r.fg(r.th.green).Bold(true).Render(prefix) + r.fg(r.th.subtext).Render(rest)})
	}
	return append(lines, r.mapping(rows)...)
}

// duplicates lists the numbers documents of the log share, and how to give
// each its own: by when git added them, or by hand.
func (r *report) duplicates(l doctor.Log) []string {
	subtext := r.fg(r.th.subtext)
	var lines []string
	clash := doctor.NoClash
	for _, d := range l.Duplicates {
		clash = max(clash, d.Clash)
		lines = append(lines, r.cross()+" "+fmt.Sprintf("%d %s share the number ", len(d.Documents), l.Kind.Plural)+bold(strconv.Itoa(d.Number)))
		rows := make([][]string, 0, len(d.Documents))
		for i, name := range d.Documents {
			if len(l.Dirs) > 1 {
				name = r.fg(r.th.faint).Render(d.Dir+"/") + name
			}
			note := ""
			if d.Clash == doctor.NoClash {
				note = "· renumbered"
				if i == 0 {
					note = fmt.Sprintf("· added first, keeps %d", d.Number)
				}
			}
			rows = append(rows, []string{name, subtext.Render(note)})
		}
		lines = append(lines, grid(rows)...)
	}

	switch clash {
	case doctor.NoClash:
		lines = append(lines, r.hang(r.fg(r.th.primary).Render("→"), bold("--fix")+" can renumber them by when git added each:"))
	case doctor.NoHistory:
		lines = append(lines,
			r.hang(r.bang(), "Needs fixing by hand: this is not a git repository, so which came first is not known."),
			r.hang(" ", subtext.Render("Keep the first of each, and renumber the others, such as:")))
	case doctor.NamesClash:
		lines = append(lines,
			r.hang(r.bang(), "Needs fixing by hand: fix the file names that clash first, and doctor can renumber these."),
			r.hang(" ", subtext.Render("Or keep the first of each, and renumber the others, such as:")))
	}
	lines = append(lines, r.renumbers(l, l.Renumbers)...)
	n := len(l.Renumbers)
	warning := "Renumbering renames " + count(n, l.Kind.Noun, l.Kind.Plural) + ". Links to " + them(n) +
		" in the project's Markdown are updated, but other references to " + them(n) + " could be affected."
	if clash != doctor.NoClash {
		warning = "Renumbering renames " + count(n, l.Kind.Noun, l.Kind.Plural) + ": update the title of each, the links to " +
			them(n) + ", and anything else that refers to " + them(n) + " by name or number."
	}
	return append(lines, r.hang(r.bang(), r.fg(r.th.yellow).Render(warning)))
}

// mistitled lists the documents whose titles do not match their file names,
// each with what differs and, under it, its title, which may be long.  They
// are fixed by hand: which of the two is wrong is not known.
func (r *report) mistitled(l doctor.Log) []string {
	subtext := r.fg(r.th.subtext)
	lines := []string{r.cross() + " Titles do not match file names"}
	rows := make([][]string, 0, len(l.Mistitled))
	for _, m := range l.Mistitled {
		name := m.File
		if len(l.Dirs) > 1 {
			name = r.fg(r.th.faint).Render(m.Dir+"/") + name
		}
		var differs string
		switch {
		case m.Number && m.Words:
			differs = "· number and words differ"
		case m.Number:
			differs = "· number differs"
		default:
			differs = "· words differ"
		}
		rows = append(rows, []string{name, subtext.Render(differs)})
	}
	for i, row := range grid(rows) {
		lines = append(lines, row, "    "+r.fg(r.th.red).Render("# "+l.Mistitled[i].Title))
	}
	return append(lines, r.hang(r.bang(), "Needs fixing by hand: correct the title, or rename the "+l.Kind.Noun+
		" and the links to it, whichever is wrong."))
}

// them is the pronoun for n things.
func them(n int) string {
	if n == 1 {
		return "it"
	}
	return "them"
}

// hang writes text after icon, wrapped to fit a box, with the lines after the
// first indented to line up with the first.
func (r *report) hang(icon, text string) string {
	inner := r.width - 4 // the box's borders and padding
	wrapped := lipgloss.NewStyle().Width(inner - 2).Render(text)
	lines := strings.Split(wrapped, "\n")
	for i, line := range lines {
		line = strings.TrimRight(line, " ")
		if i == 0 {
			lines[i] = icon + " " + line
		} else {
			lines[i] = "  " + line
		}
	}
	return strings.Join(lines, "\n")
}

// renumbers lays out documents renumbered, or to be, with the number that
// changes marked, and where each that shared a number goes in the log.
func (r *report) renumbers(l doctor.Log, renumbers []doctor.Renumber) []string {
	subtext := r.fg(r.th.subtext)
	dirs := map[string]bool{}
	for _, rn := range renumbers {
		dirs[rn.Dir] = true
	}
	rows := make([][]string, 0, len(renumbers))
	var places []string
	for _, rn := range renumbers {
		from := r.number(l, rn.From, r.th.red, false)
		if len(dirs) > 1 || len(l.Dirs) > 1 {
			from = r.fg(r.th.faint).Render(rn.Dir+"/") + from
		}
		rows = append(rows, []string{from, r.arrow(), r.number(l, rn.To, r.th.green, true)})
		id := func(name string) string { return ident(l.Kind, l.Kind.Number(name)) }
		switch {
		case rn.After == "":
		case rn.Before == "":
			places = append(places, r.hang(" ", subtext.Render(id(rn.To)+" goes last, after "+id(rn.After)+", the last added before it.")))
		default:
			places = append(places, r.hang(" ", subtext.Render(id(rn.To)+" goes between "+id(rn.After)+", the last added before it, and "+id(rn.Before)+".")))
		}
	}
	return append(r.mapping(rows), places...)
}

// number shows a document's name with its number in colour c, and the rest
// subdued when the name is what a document should have.
func (r *report) number(l doctor.Log, name string, c color.Color, should bool) string {
	prefix := ""
	if l.Kind.Prefix != "" {
		prefix = l.Kind.Prefix + "-"
	}
	rest := strings.TrimPrefix(name, prefix)
	i := strings.IndexFunc(rest, func(r rune) bool { return r < '0' || r > '9' })
	if i < 0 {
		return name
	}
	style := lipgloss.NewStyle()
	if should {
		style = r.fg(r.th.subtext)
	}
	return style.Render(prefix) + r.fg(c).Bold(should).Render(rest[:i]) + style.Render(rest[i:])
}

// mapping lays out rows of what is there, an arrow or a warning, and what
// should be, in a grid.  When the grid is too wide for a box, each row is
// stacked instead, with what should be beneath what is there.
func (r *report) mapping(rows [][]string) []string {
	lines := grid(rows)
	inner := r.width - 4 // the box's borders and padding
	for _, line := range lines {
		if lipgloss.Width(line) > inner {
			lines = lines[:0]
			for _, row := range rows {
				lines = append(lines, "  "+row[0], "    "+row[1]+" "+row[2])
			}
			break
		}
	}
	return lines
}

// fixBox draws what a fix changed in one log, and what is left wrong with it:
// the marker files renamed, the moves, the directories they left empty and
// that were removed, the renames, and the problems that need fixing by hand.
func (r *report) fixBox(l doctor.Log, f doctor.Fixes, removed []string) string {
	subtext := r.fg(r.th.subtext)
	var lines []string
	for _, m := range f.Markers {
		lines = append(lines, r.tick()+" Renamed from a former name to "+bold(m.To))
		lines = append(lines, r.mapping([][]string{{m.From, r.arrow(), r.fg(r.th.green).Bold(true).Render(m.To)}})...)
	}
	for _, m := range l.FormerMarkers {
		lines = append(lines, r.formerMarker(l, m)...)
	}
	for _, m := range f.Moved {
		switch m.Reason {
		case doctor.UnderDoc:
			lines = append(lines, r.tick()+" Moved from "+bold("doc/")+" to "+bold("docs/"))
		case doctor.Former:
			lines = append(lines, r.tick()+" Moved from a former directory to "+bold(m.To))
		}
		lines = append(lines, r.mapping([][]string{{m.From, r.arrow(), r.to(m)}})...)
	}
	for _, dir := range removed {
		lines = append(lines, r.tick()+" Removed "+bold(dir+"/")+subtext.Render(", which is empty now"))
	}
	if len(f.Renamed) > 0 {
		lines = append(lines, r.renamed(l, f.Renamed)...)
	}
	for _, m := range l.Misplaced {
		lines = append(lines, r.misplaced(l, m)...)
	}
	if len(f.Renumbered) > 0 {
		lines = append(lines, r.hang(r.tick(), "Renumbered by when git added each:"))
		lines = append(lines, r.renumbers(l, f.Renumbered)...)
		lines = append(lines, r.hang(r.bang(), r.fg(r.th.yellow).Render(
			"Links to "+them(len(f.Renumbered))+" in the project's Markdown were updated. Check other references to "+
				them(len(f.Renumbered))+" by name or number, such as links from elsewhere, other branches, and mentions in text.")))
	}
	if len(l.Misnamed) > 0 {
		lines = append(lines, r.misnamed(l)...)
	}
	if len(l.Duplicates) > 0 {
		lines = append(lines, r.duplicates(l)...)
	}

	border, title := r.logTitle(l)
	return r.box(title, r.logAside(l, f.Count()), strings.Join(lines, "\n"), border)
}

// renamed lists the documents a fix renamed, with their old and new names.
// When they were in more than one directory, each is shown with the
// directory it was in.
func (r *report) renamed(l doctor.Log, renames []doctor.Rename) []string {
	prefix := l.Kind.Prefix + "-"
	lines := []string{r.tick() + " Renamed to begin with " + bold(prefix)}
	dirs := map[string]bool{}
	for _, rn := range renames {
		dirs[rn.Dir] = true
	}
	rows := make([][]string, 0, len(renames))
	for _, rn := range renames {
		rest := strings.TrimPrefix(rn.To, prefix)
		from := r.fg(r.th.faint).Render(strings.TrimSuffix(rn.From, rest)) + rest
		if len(dirs) > 1 {
			from = r.fg(r.th.faint).Render(rn.Dir+"/") + from
		}
		rows = append(rows, []string{from, r.arrow(), r.fg(r.th.green).Bold(true).Render(prefix) + r.fg(r.th.subtext).Render(rest)})
	}
	return append(lines, r.mapping(rows)...)
}

// linksBox lists the documents whose links a fix changed, and those whose
// titles it renumbered.
func (r *report) linksBox(relinked []doctor.Relinked) string {
	subtext := r.fg(r.th.subtext)
	links, titles := 0, 0
	rows := make([][]string, 0, len(relinked))
	for _, d := range relinked {
		links += d.Links
		var changed []string
		if d.Links > 0 {
			changed = append(changed, count(d.Links, "link", "links"))
		}
		if d.Retitled {
			titles++
			changed = append(changed, "title")
		}
		rows = append(rows, []string{d.Path, subtext.Render(strings.Join(changed, ", "))})
	}
	text := "Changed to lead to the documents where they are now."
	if titles > 0 {
		text = "Changed to lead to the documents where they are now, and by the numbers they have now."
	}
	lines := append([]string{r.hang(r.tick(), subtext.Render(text))}, grid(rows)...)
	title := r.tick() + " " + bold("links")
	var aside []string
	if links > 0 {
		aside = append(aside, count(links, "link", "links"))
	}
	if titles > 0 {
		aside = append(aside, count(titles, "title", "titles"))
	}
	return r.box(title, r.fg(r.th.faint).Render(strings.Join(aside, " and ")+" in "+count(len(relinked), "document", "documents")),
		strings.Join(lines, "\n"), r.th.border)
}

// healthyBox lists every log, each in its place.
func (r *report) healthyBox() string {
	subtext := r.fg(r.th.subtext)
	rows := make([][]string, 0, len(r.result.Logs))
	for _, l := range r.result.Logs {
		if !l.Found() {
			rows = append(rows, []string{r.dash(), bold(l.Kind.Name), subtext.Render(l.Dir), subtext.Render("no log"), ""})
			continue
		}
		rows = append(rows, []string{
			r.tick(), bold(l.Kind.Name), l.Dir,
			subtext.Render(count(l.Documents(), l.Kind.Noun, l.Kind.Plural)),
			subtext.Render(span(l)),
		})
	}
	var lines []string
	for _, line := range grid(rows) {
		lines = append(lines, strings.TrimPrefix(line, "  "))
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).BorderForeground(r.th.border).
		Padding(0, 1).Width(r.width).
		Render(strings.Join(lines, "\n"))
}

// span names the first and last documents of a log that is in its place:
// ADR-0001 … ADR-0003.
func span(l doctor.Log) string {
	if len(l.Dirs) == 0 || len(l.Dirs[0].Documents) == 0 {
		return ""
	}
	first, last := 0, 0
	for _, name := range l.Dirs[0].Documents {
		n := l.Kind.Number(name)
		if first == 0 || n < first {
			first = n
		}
		last = max(last, n)
	}
	if first == last {
		return ident(l.Kind, first)
	}
	return ident(l.Kind, first) + " … " + ident(l.Kind, last)
}

// ident names a document by its kind's prefix and its number: ADR-0003.
func ident(kind record.Kind, n int) string {
	if kind.Prefix == "" {
		return fmt.Sprintf("%04d", n)
	}
	return fmt.Sprintf("%s-%04d", kind.Prefix, n)
}

// barStyle is the background of the status bar, with text in colour c, or in
// the terminal's own colour when c is nil.
func (r *report) barStyle(c color.Color) lipgloss.Style {
	style := lipgloss.NewStyle().Background(r.th.bar)
	if c != nil {
		style = style.Foreground(c)
	}
	return style
}

// statusBar sums up across the width of the report, ending in the exit status.
func (r *report) statusBar() string {
	if r.result.Healthy() {
		return r.bar(r.barStyle(r.th.green).Render("✓")+r.barStyle(nil).Render(" "+r.healthySummary()+"  "), "Nothing was changed.")
	}
	troubled := 0
	for _, l := range r.result.Logs {
		if l.Problems() > 0 {
			troubled++
		}
	}
	summary := r.barStyle(r.th.red).Render("✗") + r.barStyle(nil).Render(fmt.Sprintf(" %s in %s  ",
		count(r.result.Problems(), "problem", "problems"), count(troubled, "log", "logs")))
	return r.bar(summary+r.byHandSummary(), "Nothing was changed.")
}

// fixStatusBar sums up what a fix changed, and what it left.
func (r *report) fixStatusBar(fixed doctor.Fixed) string {
	tick := r.barStyle(r.th.green).Render("✓")
	plain := r.barStyle(nil)
	summary := tick + plain.Render(" Fixed "+count(fixed.Count(), "problem", "problems")+".")
	if r.result.Healthy() {
		summary += plain.Render(" All is well.  ")
	} else {
		summary += plain.Render("  ")
		if left := r.result.Problems() - r.result.ByHand(); left > 0 {
			summary += r.barStyle(r.th.red).Render("✗") + plain.Render(" "+count(left, "problem", "problems")+" left  ")
		}
		summary += r.byHandSummary()
	}
	return r.bar(summary, "Made "+count(fixed.Changes(), "change", "changes")+".")
}

// byHandSummary counts the problems that need fixing by hand, if any.
func (r *report) byHandSummary() string {
	n := r.result.ByHand()
	if n == 0 {
		return ""
	}
	verb := "needs"
	if n > 1 {
		verb = "need"
	}
	yellow := r.barStyle(r.th.yellow)
	return yellow.Bold(true).Render("!") + yellow.Render(fmt.Sprintf(" %d %s fixing by hand  ", n, verb))
}

// bar lays out the status bar: the badge, the summary, a note when there is
// room for it, and the exit status.
func (r *report) bar(summary, note string) string {
	exit := r.badge("EXIT 0", r.th.green)
	if !r.result.Healthy() {
		exit = r.badge("EXIT 1", r.th.red)
	}
	left := r.badge("DOCTOR", r.th.primary) + r.barStyle(nil).Render(" ") + summary
	fill := r.width - lipgloss.Width(left) - lipgloss.Width(exit)
	if fill < len(note) {
		note = ""
	}
	return left + r.barStyle(r.th.subtext).Width(max(fill, 0)).Render(note) + exit
}

func (r *report) healthySummary() string {
	found := 0
	var only string
	for _, l := range r.result.Logs {
		if l.Found() {
			found++
			only = l.Kind.Name
		}
	}
	switch {
	case found == 0:
		return "No logs found."
	case found == 1:
		return "The " + only + " log is in place and named with its prefix."
	case found == len(r.result.Logs):
		return fmt.Sprintf("All %d logs are in place and named with their prefix.", found)
	default:
		return fmt.Sprintf("The %d logs found are in place and named with their prefix.", found)
	}
}

// grid lays rows out in columns, two spaces apart, indented by two.  Every
// column but the last is padded to its widest cell.
func grid(rows [][]string) []string {
	var widths []int
	for _, row := range rows {
		for i, cell := range row {
			if i == len(widths) {
				widths = append(widths, 0)
			}
			widths[i] = max(widths[i], lipgloss.Width(cell))
		}
	}
	lines := make([]string, 0, len(rows))
	for _, row := range rows {
		var line strings.Builder
		line.WriteString("  ")
		for i, cell := range row {
			line.WriteString(cell)
			if i < len(row)-1 {
				line.WriteString(strings.Repeat(" ", widths[i]-lipgloss.Width(cell)+2))
			}
		}
		lines = append(lines, strings.TrimRight(line.String(), " "))
	}
	return lines
}

// count writes n with the singular or plural noun: "1 document", "3
// documents".
func count(n int, singular, plural string) string {
	if n == 1 {
		return "1 " + singular
	}
	return fmt.Sprintf("%d %s", n, plural)
}

// series joins words as a list in prose: "adr", "adr or rfc", "adr, rfc or
// plm".
func series(words []string) string {
	if len(words) < 2 {
		return strings.Join(words, "")
	}
	return strings.Join(words[:len(words)-1], ", ") + " or " + words[len(words)-1]
}

// the writes n of a noun as the object of a sentence: "the user story", "the
// 2 user stories".
func the(n int, singular, plural string) string {
	if n == 1 {
		return "the " + singular
	}
	return fmt.Sprintf("the %d %s", n, plural)
}
