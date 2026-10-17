package output

import (
	"errors"
	"io"
	"path"
	"strings"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"

	"github.com/mustardcustard/detoo/internal/features/health/doctor"
)

// Question is a question answered yes or no.
type Question struct {
	// Title is the question itself, and Description what the user needs to
	// know to answer it.
	Title, Description string

	// Yes is the answer chosen to begin with.
	Yes bool
}

// ErrStopped is returned when the user stops a question without answering.
var ErrStopped = errors.New("stopped without an answer; nothing was changed")

// Ask asks q at the terminal t, reading keys from in and drawing on out, as a
// choice between Yes and No, and returns the answer.
func Ask(in io.Reader, out io.Writer, t Terminal, q Question) (bool, error) {
	answer := q.Yes
	field := huh.NewConfirm().
		Title(q.Title).
		Description(q.Description).
		Affirmative("Yes").
		Negative("No").
		Value(&answer)
	width := t.Width
	if width <= 0 {
		width = 80
	}
	form := huh.NewForm(huh.NewGroup(field)).
		WithTheme(confirmTheme(t)).
		WithWidth(min(width, widest)).
		WithShowHelp(false).
		WithInput(in).
		WithOutput(out)
	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return false, ErrStopped
		}
		return false, err
	}
	return answer, nil
}

// confirmTheme draws a question in the report's palette: its title in the
// primary colour, and the button chosen on it.
func confirmTheme(t Terminal) huh.Theme {
	return huh.ThemeFunc(func(dark bool) *huh.Styles {
		th := newTheme(t.Dark)
		s := huh.ThemeBase(t.Dark)
		s.Focused.Title = s.Focused.Title.Foreground(th.primary).Bold(true)
		s.Focused.Description = s.Focused.Description.Foreground(th.subtext)
		s.Focused.FocusedButton = s.Focused.FocusedButton.Foreground(th.base).Background(th.primary).Bold(true)
		s.Focused.BlurredButton = s.Focused.BlurredButton.Foreground(th.subtext).Background(th.bar)
		s.Focused.Base = s.Focused.Base.BorderForeground(th.border)
		s.Blurred = s.Focused
		s.Blurred.Base = s.Focused.Base.BorderStyle(lipgloss.HiddenBorder())
		return s
	})
}

// Plan writes what a fix will do that it does not ask about: each log it
// moves, and each document it renames to begin with its prefix.  It writes
// nothing when there is none.
func Plan(w io.Writer, t Terminal, result doctor.Result) error {
	lines := PlanItems(result)
	if len(lines) == 0 {
		return nil
	}
	th := newTheme(t.Dark)
	bullet := lipgloss.NewStyle().Foreground(th.primary).Render("•")
	text := bold("detoo doctor --fix") + " will:\n"
	for _, line := range lines {
		text += "  " + bullet + " " + line + "\n"
	}
	_, err := io.WriteString(&colorprofile.Writer{Forward: w, Profile: t.Profile}, text+"\n")
	return err
}

// PlanItems lists what a fix will do that it does not ask about, a line for
// each log it moves and each document it renames.
func PlanItems(result doctor.Result) []string {
	var lines []string
	for _, l := range result.Logs {
		for _, m := range l.Misplaced {
			if m.Clash == doctor.NoClash {
				lines = append(lines, "Move "+m.From+"/ → "+m.To+"/")
			}
		}
		for _, rn := range l.Misnamed {
			if !rn.Clash {
				lines = append(lines, "Rename "+path.Join(rn.Dir, rn.From)+" → "+rn.To)
			}
		}
	}
	return lines
}

// RenumberQuestion asks whether to renumber the documents in result that
// share a number, listing each document it would rename, and warning of what
// renaming them could break.  It is No to begin with.
func RenumberQuestion(result doctor.Result) Question {
	var b strings.Builder
	b.WriteString("Each document gets its own number, in the order git added them:\n")
	for _, l := range result.Logs {
		for _, r := range l.Renumbers {
			if r.Clash == doctor.NoClash {
				b.WriteString("  • " + r.From + " → " + r.To + "\n")
			}
		}
	}
	b.WriteString("Links to them in this project's Markdown files will be updated.\n")
	b.WriteString("Other references to them may break: other branches, links from outside the project, " +
		"and mentions in text, tickets or commit messages.")
	return Question{Title: "Renumber documents that share a number?", Description: b.String()}
}

// GitQuestion asks whether to move and rename with git mv.  It is Yes to
// begin with.
var GitQuestion = Question{
	Title:       "Use git mv?",
	Description: "This project is in a git repository. With git mv, git records each move and rename as a rename, not as a deleted file and a new one.",
	Yes:         true,
}
