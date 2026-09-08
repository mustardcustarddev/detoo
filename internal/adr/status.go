package adr

import (
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"
)

// statusHeading is the heading that opens a record's status section.
const statusHeading = "## Status"

// linkPattern matches a link in a status section, capturing its description and
// the number of the record it points at.
var linkPattern = regexp.MustCompile(`^(.+) \[.*\]\(0*([1-9][0-9]*).*\)`)

// Link is a reference from one decision record to another.
type Link struct {
	Type   string // how the source record describes the link, e.g. "Supercedes"
	Target int    // the number of the record linked to
}

// Status returns the lines of the status section of the record referred to by
// ref, ignoring headings and blank lines.
func (e *Env) Status(ref string) ([]string, error) {
	file, err := e.Resolve(ref)
	if err != nil {
		return nil, err
	}
	lines, err := readLines(e.path(file))
	if err != nil {
		return nil, err
	}
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
	return status, nil
}

// Links returns the links recorded in the status section of the record referred
// to by ref.
func (e *Env) Links(ref string) ([]Link, error) {
	status, err := e.Status(ref)
	if err != nil {
		return nil, err
	}
	var links []Link
	for _, line := range status {
		match := linkPattern.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		target, err := strconv.Atoi(match[2])
		if err != nil {
			continue
		}
		links = append(links, Link{Type: match[1], Target: target})
	}
	return links, nil
}

// AddLink records a link from the source record to the target record, at the end
// of the source's status section.
func (e *Env) AddLink(sourceRef, linkType, targetRef string) error {
	source, err := e.Resolve(sourceRef)
	if err != nil {
		return err
	}
	target, err := e.Resolve(targetRef)
	if err != nil {
		return err
	}
	title, err := e.Title(target)
	if err != nil {
		return err
	}
	lines, err := readLines(e.path(source))
	if err != nil {
		return err
	}
	link := fmt.Sprintf("%s [%s](%s)", linkType, title, path.Base(target))

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
	return writeLines(e.path(source), linked)
}

// RemoveStatus deletes a status line from the record referred to by ref,
// collapsing the blank line it leaves behind.
func (e *Env) RemoveStatus(status, ref string) error {
	file, err := e.Resolve(ref)
	if err != nil {
		return err
	}
	lines, err := readLines(e.path(file))
	if err != nil {
		return err
	}
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
	return writeLines(e.path(file), remaining)
}
