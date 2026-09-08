package adr

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// recordName matches the file name of a decision record: a number, a hyphen and
// a slug.
var recordName = regexp.MustCompile(`^[0-9]+-[^/]*\.md`)

// leadingNumber matches the number a record's file name starts with.
var leadingNumber = regexp.MustCompile(`^[0-9]+`)

// MissingDirError reports that the decision record directory does not exist.
type MissingDirError struct {
	Dir string
}

func (e *MissingDirError) Error() string {
	return fmt.Sprintf("the %s directory does not exist", e.Dir)
}

// List returns the decision records, in order, as paths relative to the working
// directory.
func (e *Env) List() ([]string, error) {
	dir, err := e.ADRDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(e.path(dir))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, &MissingDirError{Dir: dir}
	} else if err != nil {
		return nil, err
	}
	var records []string
	for _, entry := range entries {
		if !entry.IsDir() && recordName.MatchString(entry.Name()) {
			records = append(records, path.Join(dir, entry.Name()))
		}
	}
	sort.Strings(records)
	return records, nil
}

// Resolve returns the record referred to by ref, which is either a record's
// number or any part of its file name.  The first matching record wins.
func (e *Env) Resolve(ref string) (string, error) {
	records, err := e.List()
	if err != nil {
		return "", err
	}
	for _, record := range records {
		if strings.Contains(record, ref) {
			return record, nil
		}
	}
	return "", fmt.Errorf("no decision record matches %q", ref)
}

// Title returns the title of the record referred to by ref, taken from its
// first line.
func (e *Env) Title(ref string) (string, error) {
	file, err := e.Resolve(ref)
	if err != nil {
		return "", err
	}
	lines, err := readLines(e.path(file))
	if err != nil {
		return "", err
	}
	if len(lines) == 0 {
		return "", nil
	}
	return strings.TrimPrefix(lines[0], "# "), nil
}

// Number returns the number a record's file name starts with.
func Number(file string) int {
	digits := leadingNumber.FindString(path.Base(file))
	number, err := strconv.Atoi(digits)
	if err != nil {
		return 0
	}
	return number
}

// Slug converts a title into the part of a record's file name that follows its
// number: lower case, with each run of other characters collapsed to a hyphen.
func Slug(title string) string {
	var slug strings.Builder
	hyphenated := false
	for _, r := range title {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			slug.WriteRune(unicode.ToLower(r))
			hyphenated = false
		} else if !hyphenated {
			slug.WriteByte('-')
			hyphenated = true
		}
	}
	return strings.Trim(slug.String(), "-")
}

// readLines reads a file as lines, dropping the final line terminator so that
// writeLines round-trips it.
func readLines(file string) ([]string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	text := strings.TrimSuffix(string(data), "\n")
	if text == "" {
		return nil, nil
	}
	return strings.Split(text, "\n"), nil
}

// writeLines writes lines to a file, each terminated by a newline.
func writeLines(file string, lines []string) error {
	var text strings.Builder
	for _, line := range lines {
		text.WriteString(line)
		text.WriteByte('\n')
	}
	return os.WriteFile(file, []byte(text.String()), 0o644)
}
