package adr

import "testing"

func TestRenderReplacesEachPlaceholderOnceALine(t *testing.T) {
	template := "# NUMBER. TITLE\n\nDate: DATE\n\n## Status\n\nSTATUS\n"
	want := "# 7. Use MySQL\n\nDate: 1992-01-12\n\n## Status\n\nAccepted\n"
	if got := render(template, 7, "Use MySQL", "1992-01-12", "Accepted"); got != want {
		t.Errorf("render() =\n%s\nwant\n%s", got, want)
	}
}

// A title that happens to contain a placeholder is inserted, not rescanned.
func TestRenderDoesNotRescanTheTitle(t *testing.T) {
	want := "# 1. The NUMBER of things\n"
	if got := render("# NUMBER. TITLE\n", 1, "The NUMBER of things", "", ""); got != want {
		t.Errorf("render() = %q, want %q", got, want)
	}
}
