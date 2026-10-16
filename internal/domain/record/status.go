package record

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

// statusHeading is the heading that opens a record's status section.
const statusHeading = "## Status"

// linkPattern matches a link in a status section, capturing its description and
// the file of the record it points at.
var linkPattern = regexp.MustCompile(`^(.+) \[.*\]\((.*)\)`)

// Link is a reference from one record to another.
type Link struct {
	Type   string // how the source record describes the link, e.g. "Supersedes"
	Target int    // the number of the record linked to
}

// Status returns the lines of a record's status section, ignoring headings and
// blank lines.
func Status(lines []string) []string {
	var status []string
	inStatus := false
	for _, line := range lines {
		if !inStatus {
			if !strings.HasPrefix(line, statusHeading) {
				continue
			}
			inStatus = true
		} else if strings.HasPrefix(line, "#") {
			break
		}
		if !strings.HasPrefix(line, "#") && strings.TrimSpace(line) != "" {
			status = append(status, line)
		}
	}
	return status
}

// Links returns the links recorded in a record's status section.  A link's
// target is numbered by its file name, as this kind names its records, and a
// link to any other file is not one.
func (k Kind) Links(lines []string) []Link {
	var links []Link
	for _, line := range Status(lines) {
		match := linkPattern.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		if target := k.Number(match[2]); target > 0 {
			links = append(links, Link{Type: match[1], Target: target})
		}
	}
	return links
}

// AddLink returns lines with a link to the target record appended to the status
// section, described by linkType and labelled with the target's title.
//
// The link is inserted before the heading that closes the status section, which
// means a record whose status section runs to the end of the file is returned
// unchanged.  Every record created from a template has a Context section after
// its status, so this does not arise in a log built by this tool; the behaviour
// is preserved as it stands rather than corrected, because the records it would
// alter are already written.
func AddLink(lines []string, linkType, targetTitle, targetFile string) []string {
	link := fmt.Sprintf("%s [%s](%s)", linkType, targetTitle, path.Base(targetFile))

	linked := make([]string, 0, len(lines)+2)
	inStatus := false
	for _, line := range lines {
		if strings.HasPrefix(line, "##") {
			if inStatus {
				linked = append(linked, link, "")
			}
			inStatus = false
		}
		if line == statusHeading {
			inStatus = true
		}
		linked = append(linked, line)
	}
	return linked
}

// RemoveStatus returns lines with a status line deleted from the status
// section, collapsing the blank line it leaves behind.  Only the status section
// is searched, so the same text elsewhere in the record is left alone.
func RemoveStatus(lines []string, status string) []string {
	remaining := make([]string, 0, len(lines))
	inStatus, afterBlank := false, false
	for _, line := range lines {
		if strings.HasPrefix(line, "##") {
			inStatus = false
		}
		if line == statusHeading {
			inStatus = true
		}
		if inStatus {
			if strings.TrimSpace(line) == "" {
				if !afterBlank {
					remaining = append(remaining, line)
				}
				afterBlank = true
				continue
			}
			if line == status {
				continue
			}
			afterBlank = false
		}
		remaining = append(remaining, line)
	}
	return remaining
}
