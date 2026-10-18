package doctor

import (
	"context"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/mustardcustard/detoo/internal/domain/record"
)

// heading matches a document's title, its first line, as the tools write it:
// "# 3. Use PostgreSQL", capturing the number, if there is one, and the rest.
var heading = regexp.MustCompile(`^# (?:([0-9]+)\. )?(.*)$`)

// mistitled returns a mistitle for each document in dir, of kind, whose title
// does not match its file name: numbered otherwise, or not at all, or in
// other words.  Words match when they make the same slug, as the tools name a
// document from its title, so that case and punctuation do not count.  A
// document whose first line is not a title, such as an empty one, has no
// title to match.
func mistitled(ctx context.Context, project Project, kind record.Kind, kinds []record.Kind, dir Dir) ([]Mistitle, error) {
	var found []Mistitle
	for _, name := range dir.Documents {
		content, err := project.Read(ctx, path.Join(dir.Path, name))
		if err != nil {
			return nil, err
		}
		first, _, _ := strings.Cut(content, "\n")
		m := heading.FindStringSubmatch(strings.TrimSuffix(first, "\r"))
		if m == nil {
			continue
		}
		number, slug := parts(kind, kinds, name)
		titled, err := strconv.Atoi(m[1]) // fails when the title has no number
		mistitle := Mistitle{
			Dir:    dir.Path,
			File:   name,
			Title:  strings.TrimPrefix(m[0], "# "),
			Number: err != nil || titled != number,
			Words:  record.Slug(m[2]) != record.Slug(slug),
		}
		if mistitle.Number || mistitle.Words {
			found = append(found, mistitle)
		}
	}
	return found, nil
}

// parts returns the number and the slug of a document of kind, from the
// name it should have.
func parts(kind record.Kind, kinds []record.Kind, name string) (int, string) {
	to, _ := canonical(kind, kinds, name)
	if kind.Prefix != "" {
		to = strings.TrimPrefix(to, kind.Prefix+"-")
	}
	digits := leadingDigits.FindString(to)
	number, _ := strconv.Atoi(digits)
	slug := strings.TrimSuffix(strings.TrimPrefix(to[len(digits):], "-"), ".md")
	return number, slug
}
