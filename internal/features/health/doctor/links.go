package doctor

import (
	"path"
	"regexp"
	"strings"
)

var (
	// inlineLink matches a link or an image written inline, [text](target)
	// or [text](target "title"), capturing its text and its target.
	inlineLink = regexp.MustCompile(`!?\[([^\]]*)\]\(([^()\s]+)(?:\s+(?:"[^"]*"|'[^']*'))?\)`)

	// referenceLink matches the definition of a reference link,
	// [label]: target, capturing its target.
	referenceLink = regexp.MustCompile(`^ {0,3}\[[^\]]+\]:[ \t]*(\S+)`)

	// scheme begins a target that is a URL rather than a path: https:,
	// mailto:.
	scheme = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+.-]*:`)

	// fence opens or closes a block of code, whose links are left alone.
	fence = regexp.MustCompile("^ {0,3}(```|~~~)")
)

// Relink returns content, a Markdown document that was at from and is now at
// to, with each link that would break changed so that it leads where it did,
// and how many it changed.  where returns the path after the fix of a path
// before it.
//
// A relative link is resolved from where the document was, and is changed
// only when it no longer resolves, from where the document is, to where its
// target is now.  A link from the project's root, /docs/adr, is changed when
// its target moved.  Links to other sites, to a heading in the same document,
// and in blocks of code, are left as they are.
func Relink(content, from, to string, where func(string) string) (string, int) {
	return relink(content, from, to, where, nil)
}

// relink is Relink, which also changes the text of a link when retext, given
// the path its target is at after the fix, returns new text for it.
func relink(content, from, to string, where func(string) string, retext func(target, text string) (string, bool)) (string, int) {
	lines := strings.SplitAfter(content, "\n")
	changed := 0
	inCode := false
	for i, line := range lines {
		if fence.MatchString(line) {
			inCode = !inCode
			continue
		}
		if inCode {
			continue
		}
		var n int
		line, n = replaceGroups(line, inlineLink, func(groups []string) ([]string, bool) {
			text, target := groups[0], groups[1]
			target, moved := relinkTarget(target, from, to, where)
			retexted := false
			if retext != nil {
				if at, ok := resolve(target, to); ok {
					text, retexted = retext(at, text)
				}
			}
			return []string{text, target}, moved || retexted
		})
		changed += n
		line, n = replaceGroups(line, referenceLink, func(groups []string) ([]string, bool) {
			target, ok := relinkTarget(groups[0], from, to, where)
			return []string{target}, ok
		})
		changed += n
		lines[i] = line
	}
	return strings.Join(lines, ""), changed
}

// replaceGroups replaces the groups of each match of pattern in line by what
// replace returns for them, and counts the matches it changed.
func replaceGroups(line string, pattern *regexp.Regexp, replace func([]string) ([]string, bool)) (string, int) {
	var b strings.Builder
	n, last := 0, 0
	for _, m := range pattern.FindAllStringSubmatchIndex(line, -1) {
		groups := make([]string, 0, len(m)/2-1)
		for g := 2; g < len(m); g += 2 {
			groups = append(groups, line[m[g]:m[g+1]])
		}
		replaced, ok := replace(groups)
		if !ok {
			continue
		}
		for g, group := range replaced {
			b.WriteString(line[last:m[2+2*g]])
			b.WriteString(group)
			last = m[3+2*g]
		}
		n++
	}
	if n == 0 {
		return line, 0
	}
	b.WriteString(line[last:])
	return b.String(), n
}

// resolve returns the path in the project that a link in the document at
// name leads to, and whether it leads to one: not to another site, nor to a
// heading in the same document.
func resolve(target, name string) (string, bool) {
	if i := strings.IndexAny(target, "#?"); i >= 0 {
		target = target[:i]
	}
	if target == "" || scheme.MatchString(target) {
		return "", false
	}
	if rooted, ok := strings.CutPrefix(target, "/"); ok {
		return path.Clean(rooted), true
	}
	return path.Join(path.Dir(name), target), true
}

// relinkTarget returns the target a link in the document should have now, and
// whether it differs from target.
func relinkTarget(target, from, to string, where func(string) string) (string, bool) {
	if target == "" || strings.HasPrefix(target, "#") || scheme.MatchString(target) {
		return target, false
	}
	p, suffix := target, ""
	if i := strings.IndexAny(target, "#?"); i >= 0 {
		p, suffix = target[:i], target[i:]
	}
	if p == "" {
		return target, false
	}
	slash := ""
	if strings.HasSuffix(p, "/") {
		slash = "/"
	}

	if rooted, ok := strings.CutPrefix(p, "/"); ok {
		was := path.Clean(rooted)
		now := where(was)
		if now == was {
			return target, false
		}
		return "/" + now + trailing(now, slash) + suffix, true
	}

	now := where(path.Join(path.Dir(from), p))
	if path.Join(path.Dir(to), p) == now {
		return target, false
	}
	rel, ok := relative(path.Dir(to), now)
	if !ok {
		return target, false
	}
	return rel + trailing(rel, slash) + suffix, true
}

// trailing returns the slash a link to a directory ended in, unless the new
// target ends in one already.
func trailing(target, slash string) string {
	if strings.HasSuffix(target, "/") {
		return ""
	}
	return slash
}

// relative returns the path that leads from the directory dir to target, and
// whether there is one: both must be relative, or both absolute.
func relative(dir, target string) (string, bool) {
	if path.IsAbs(dir) != path.IsAbs(target) {
		return "", false
	}
	segments := func(p string) []string {
		p = strings.TrimPrefix(path.Clean(p), "/")
		if p == "." || p == "" {
			return nil
		}
		return strings.Split(p, "/")
	}
	d, t := segments(dir), segments(target)
	common := 0
	for common < len(d) && common < len(t) && d[common] == t[common] {
		common++
	}
	up := d[common:]
	for _, s := range up {
		if s == ".." {
			return "", false
		}
	}
	parts := make([]string, 0, len(up)+len(t)-common)
	for range up {
		parts = append(parts, "..")
	}
	parts = append(parts, t[common:]...)
	if len(parts) == 0 {
		return ".", true
	}
	return strings.Join(parts, "/"), true
}
