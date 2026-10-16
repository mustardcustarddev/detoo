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

// recordName matches the file name of a record, after any prefix: a number, a
// hyphen and a slug.
var recordName = regexp.MustCompile(`^[0-9]+-[^/]*\.md`)

// leadingNumber matches the number a record's file name starts with, after any
// prefix.
var leadingNumber = regexp.MustCompile(`^[0-9]+`)

// digits matches a reference that is a record's number.
var digits = regexp.MustCompile(`^[0-9]+$`)

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

// IsRecordName reports whether a file name is that of a record of this kind,
// so that a README or a template directory is not mistaken for one.
func (k Kind) IsRecordName(name string) bool {
	rest, ok := k.unprefixed(name)
	return ok && recordName.MatchString(rest)
}

// FileName returns the path of the record numbered number, with the given
// title, in dir.
func (k Kind) FileName(dir string, number int, title string) string {
	name := fmt.Sprintf("%04d-%s.md", number, Slug(title))
	if k.Prefix != "" {
		name = k.Prefix + "-" + name
	}
	return path.Join(dir, name)
}

// Number returns the number a record's file name starts with, after its
// kind's prefix, or 0 if the file is not named as a record of this kind.
func (k Kind) Number(file string) int {
	rest, ok := k.unprefixed(path.Base(file))
	if !ok {
		return 0
	}
	number, err := strconv.Atoi(leadingNumber.FindString(rest))
	if err != nil {
		return 0
	}
	return number
}

// unprefixed returns a file name without this kind's prefix and the hyphen
// after it, and whether the name had them.
func (k Kind) unprefixed(name string) (string, bool) {
	if k.Prefix == "" {
		return name, true
	}
	return strings.CutPrefix(name, k.Prefix+"-")
}

// Find returns the record that ref refers to, and whether there is one.  A
// reference is a record's number, its path, or part of its file name, and the
// first record that matches wins.
//
// A reference in digits alone is a number, and finds that record or none,
// rather than the first whose name happens to contain the digits.  Only the
// file name is searched, never the directory, so that a reference such as
// "adr" does not match every record in docs/adr.
func Find(records []Record, ref string) (Record, bool) {
	if digits.MatchString(ref) {
		number, err := strconv.Atoi(ref)
		for _, r := range records {
			if err == nil && r.Number == number {
				return r, true
			}
		}
		return Record{}, false
	}
	for _, r := range records {
		if r.File == path.Clean(ref) || strings.Contains(path.Base(r.File), ref) {
			return r, true
		}
	}
	return Record{}, false
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
