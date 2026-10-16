// Package record holds the concepts the features of every record tool
// exchange: what a record is, how its file is named, how its Markdown is read
// and written, and the vocabulary each kind of log describes it in.
//
// Everything here is pure.  The package reads no files, writes no output, and
// knows nothing of Cobra, so a slice built on it can be driven by a CLI, an
// HTTP API or a test with equal ease.  Anything that touches the disk belongs
// in an adapter; anything that belongs to one use case belongs in that slice.
package record

import (
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// recordName matches the file name of a record: a number, a hyphen and a slug.
var recordName = regexp.MustCompile(`^[0-9]+-[^/]*\.md`)

// leadingNumber matches the number a record's file name starts with.
var leadingNumber = regexp.MustCompile(`^[0-9]+`)

// Record identifies one record in a log.
//
// It carries no content: the title and status live in the file, and a slice
// that needs them reads its lines through a store and parses them here.  That
// keeps listing a log cheap, and keeps this type free of I/O.
type Record struct {
	// Number is the record's position in the log, taken from its file name.
	Number int

	// File is the record's path, relative to the working directory and always
	// slash separated, as every path reported by these features is.
	File string
}

// MissingDirError reports that the record directory does not exist.
//
// It is a domain condition rather than a filesystem one — this project has no
// log yet — so a store adapter translates fs.ErrNotExist into it, and
// every client is free to decide what that should mean.
type MissingDirError struct {
	Dir string
}

func (e *MissingDirError) Error() string {
	return fmt.Sprintf("the %s directory does not exist", e.Dir)
}

// IsRecordName reports whether a file name is that of a record, so
// that a README or a template directory is not mistaken for one.
func IsRecordName(name string) bool {
	return recordName.MatchString(name)
}

// FileName returns the path of the record numbered number, with the given
// title, in dir.
func FileName(dir string, number int, title string) string {
	return path.Join(dir, fmt.Sprintf("%04d-%s.md", number, Slug(title)))
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

// Title returns a record's title, taken from the first of its lines.
func Title(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	return strings.TrimPrefix(lines[0], "# ")
}
