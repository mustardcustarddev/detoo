package adr

import (
	_ "embed"
	"strconv"
	"strings"
)

// DefaultTemplate is the template used for a new record when the project does
// not provide one of its own.  It follows the style described by Michael Nygard
// in http://thinkrelevance.com/blog/2011/11/15/documenting-architecture-decisions
//
//go:embed templates/template.md
var DefaultTemplate string

// InitTemplate is the template for the first record of a repository, which
// records the decision to record decisions.
//
//go:embed templates/init.md
var InitTemplate string

// ProjectTemplate is the template a project can provide, relative to its
// decision record directory, to override DefaultTemplate.
const ProjectTemplate = "templates/template.md"

// render fills in a template, replacing the first occurrence of each placeholder
// on each line.
func render(template string, number int, title, date, status string) string {
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
	i := strings.Index(s, old)
	if i < 0 {
		return s
	}
	return s[:i] + new + s[i+len(old):]
}
