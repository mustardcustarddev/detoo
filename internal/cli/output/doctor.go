package output

import (
	"fmt"
	"image/color"
	"io"
	"path"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"

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
	r := newReport(t, result)
	var b strings.Builder
	b.WriteString(r.header())
	b.WriteString("\n\n")
	if result.Healthy() {
		b.WriteString(r.healthyBox())
	} else {
		boxes := make([]string, 0, len(result.Logs))
		var missing []string
		for _, l := range result.Logs {
			if l.Found() {
				boxes = append(boxes, r.logBox(l))
			} else {
				missing = append(missing, l.Kind.Name)
			}
		}
		if len(missing) > 0 {
			boxes = append(boxes, r.dash()+" "+r.fg(r.th.subtext).Render("No "+series(missing)+" log found."))
		}
		b.WriteString(strings.Join(boxes, "\n\n"))
	}
	b.WriteString("\n\n")
	b.WriteString(r.statusBar())
	b.WriteString("\n")

	_, err := io.WriteString(&colorprofile.Writer{Forward: w, Profile: t.Profile}, b.String())
	return err
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
	if markSize+len(gap)+lipgloss.Width(info) > r.width {
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

// markSize is the width and height of the mark, in pixels.  Two pixels share
// a character, one above the other, so it is markSize columns by half as many
// lines.
const markSize = 14

// markPalette is the mark's bands, from the amber of the avatar's set square
// at its corner to the orange of its shoulder.
var markPalette = []string{"#fdd08a", "#fcbf5a", "#fba200", "#f59024", "#e57929", "#d66f24", "#b85a1a"}

// markPixel returns the colour of the mark's pixel at row and col, or "" for
// none: a set square, a right triangle with its centre cut out.
func markPixel(row, col int) string {
	if col > row || (row <= 10 && col >= 3 && col <= row-3) {
		return ""
	}
	return markPalette[min(len(markPalette)-1, (row+col)*len(markPalette)/27)]
}

// mark draws the detoo mark in half blocks.  Without colour, a character
// holding two pixels of different colours is drawn as a full block, so that
// the shape survives.
func (r *report) mark() string {
	colour := r.t.Profile > colorprofile.ASCII
	lines := make([]string, 0, markSize/2)
	for row := 0; row < markSize; row += 2 {
		var line strings.Builder
		for col := range markSize {
			line.WriteString(halfBlock(markPixel(row, col), markPixel(row+1, col), colour))
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

	aside := count(l.Documents(), l.Kind.Noun, l.Kind.Plural)
	if l.Problems() > 0 {
		aside += " · " + count(l.Problems(), "problem", "problems")
	}
	return r.box(title, r.fg(r.th.faint).Render(aside), r.logBody(l), border)
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
	for _, m := range l.Misplaced {
		lines = append(lines, r.misplaced(l, m)...)
	}
	if len(l.Misnamed) > 0 {
		lines = append(lines, r.misnamed(l)...)
	}
	return strings.Join(lines, "\n")
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
			rows = append(rows, []string{from, r.bang(), r.fg(r.th.yellow).Render(rn.To + " is there already: fix by hand")})
			continue
		}
		rows = append(rows, []string{from, r.arrow(), r.fg(r.th.green).Bold(true).Render(prefix) + r.fg(r.th.subtext).Render(rest)})
	}
	return append(lines, r.mapping(rows)...)
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
	id := func(n int) string {
		if l.Kind.Prefix == "" {
			return fmt.Sprintf("%04d", n)
		}
		return fmt.Sprintf("%s-%04d", l.Kind.Prefix, n)
	}
	if first == last {
		return id(first)
	}
	return id(first) + " … " + id(last)
}

// statusBar sums up across the width of the report, ending in the exit status.
func (r *report) statusBar() string {
	bar := lipgloss.NewStyle().Background(r.th.bar)
	on := func(c color.Color) lipgloss.Style { return bar.Foreground(c) }

	left := r.badge("DOCTOR", r.th.primary) + bar.Render(" ")
	exit := r.badge("EXIT 0", r.th.green)
	if r.result.Healthy() {
		left += on(r.th.green).Render("✓") + bar.Render(" "+r.healthySummary()+"  ")
	} else {
		exit = r.badge("EXIT 1", r.th.red)
		troubled := 0
		for _, l := range r.result.Logs {
			if l.Problems() > 0 {
				troubled++
			}
		}
		left += on(r.th.red).Render("✗") + bar.Render(fmt.Sprintf(" %s in %s  ",
			count(r.result.Problems(), "problem", "problems"), count(troubled, "log", "logs")))
		if n := r.result.ByHand(); n > 0 {
			verb := "needs"
			if n > 1 {
				verb = "need"
			}
			left += on(r.th.yellow).Bold(true).Render("!") + on(r.th.yellow).Render(fmt.Sprintf(" %d %s fixing by hand  ", n, verb))
		}
	}

	fill := r.width - lipgloss.Width(left) - lipgloss.Width(exit)
	note := "Nothing was changed."
	if fill < len(note) {
		note = ""
	}
	return left + on(r.th.subtext).Width(max(fill, 0)).Render(note) + exit
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
// pse".
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
