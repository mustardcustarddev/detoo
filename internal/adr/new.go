package adr

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// NewOptions varies how New creates a record.
type NewOptions struct {
	// Superceded holds references to records that the new record supercedes.
	Superceded []string

	// Links holds "TARGET:LINK:REVERSE-LINK" specifications of links from the
	// new record to existing ones.
	Links []string

	// Template is the template to create the record from.  When empty the
	// project's own template is used, or the default template if it has none.
	Template string

	// NoEdit suppresses opening the new record in an editor.
	NoEdit bool
}

// New creates the next numbered decision record and returns its path, relative
// to the working directory.
func (e *Env) New(title string, opts NewOptions) (string, error) {
	if title == "" {
		return "", errors.New("no title given")
	}
	dir, err := e.ADRDir()
	if err != nil {
		return "", err
	}
	template := opts.Template
	if template == "" {
		if template, err = e.template(dir); err != nil {
			return "", err
		}
	}
	number, err := e.nextNumber(dir)
	if err != nil {
		return "", err
	}
	date := e.getenv("ADR_DATE")
	if date == "" {
		date = time.Now().Format(time.DateOnly)
	}

	file := path.Join(dir, fmt.Sprintf("%04d-%s.md", number, Slug(title)))
	if err := os.MkdirAll(e.path(dir), 0o755); err != nil {
		return "", err
	}
	record := render(template, number, title, date, "Accepted")
	if err := os.WriteFile(e.path(file), []byte(record), 0o644); err != nil {
		return "", err
	}

	for _, target := range opts.Superceded {
		if err := e.AddLink(target, "Superceded by", file); err != nil {
			return "", err
		}
		if err := e.RemoveStatus("Accepted", target); err != nil {
			return "", err
		}
		if err := e.AddLink(file, "Supercedes", target); err != nil {
			return "", err
		}
	}
	for _, spec := range opts.Links {
		target, forward, reverse, err := parseLink(spec)
		if err != nil {
			return "", err
		}
		if err := e.AddLink(file, forward, target); err != nil {
			return "", err
		}
		if err := e.AddLink(target, reverse, file); err != nil {
			return "", err
		}
	}

	if !opts.NoEdit {
		if err := e.edit(file); err != nil {
			return "", err
		}
	}
	return file, nil
}

// template returns the template a new record in dir is created from: the
// project's own, if it has one, and otherwise the default.
func (e *Env) template(dir string) (string, error) {
	if file := e.getenv("ADR_TEMPLATE"); file != "" {
		return readFile(e.path(file))
	}
	project, err := readFile(e.path(path.Join(dir, ProjectTemplate)))
	if errors.Is(err, fs.ErrNotExist) {
		return DefaultTemplate, nil
	}
	return project, err
}

// nextNumber returns the number for the next record in dir, one higher than the
// highest already there.
func (e *Env) nextNumber(dir string) (int, error) {
	entries, err := os.ReadDir(e.path(dir))
	if errors.Is(err, fs.ErrNotExist) {
		return 1, nil
	} else if err != nil {
		return 0, err
	}
	highest := 0
	for _, entry := range entries {
		if number := Number(entry.Name()); number > highest {
			highest = number
		}
	}
	return highest + 1, nil
}

// parseLink splits a "TARGET:LINK:REVERSE-LINK" link specification.
func parseLink(spec string) (target, forward, reverse string, err error) {
	parts := strings.Split(spec, ":")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("link %q is not of the form TARGET:LINK:REVERSE-LINK", spec)
	}
	return parts[0], parts[1], parts[2], nil
}

// edit opens a record in the editor named by VISUAL, or by EDITOR if VISUAL is
// not set.  With neither set, the record is left for the user to open.
func (e *Env) edit(file string) error {
	editor := e.getenv("VISUAL")
	if editor == "" {
		editor = e.getenv("EDITOR")
	}
	words := strings.Fields(editor)
	if len(words) == 0 {
		return nil
	}
	cmd := exec.Command(words[0], append(words[1:], filepath.FromSlash(file))...)
	cmd.Dir = e.Dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = e.Stdout
	cmd.Stderr = e.Stderr
	return cmd.Run()
}

func readFile(file string) (string, error) {
	data, err := os.ReadFile(file)
	return string(data), err
}
