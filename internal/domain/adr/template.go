package adr

import (
	"strconv"
	"strings"
)

// ProjectTemplate is the template a project can provide, relative to its
// decision record directory, to override the built-in one.
const ProjectTemplate = "templates/template.md"

// Render fills in a template, replacing the first occurrence of each
// placeholder on each line.
//
// Replacing only the first occurrence per line, rather than every one, is what
// stops a title that itself contains NUMBER or DATE from being rescanned and
// substituted a second time.
func Render(template string, number int, title, date, status string) string {
	placeholders := [][2]string{
		{"NUMBER", strconv.Itoa(number)},
		{"TITLE", title},
		{"DATE", date},
		{"STATUS", status},
	}
	lines := strings.Split(template, "\n")
	for i, line := range lines {
		for _, placeholder := range placeholders {
			line = replaceFirst(line, placeholder[0], placeholder[1])
		}
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}

func replaceFirst(s, old, new string) string {
	before, after, found := strings.Cut(s, old)
	if !found {
		return s
	}
	return before + new + after
}
