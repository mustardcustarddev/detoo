package skillstore

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"slices"
	"strings"

	"github.com/mustardcustard/detoo/internal/features/agents/skills"
)

// folderHash is the hash a project's lock file records of a skill, as the
// skills CLI computes it: SHA-256 over the path and content of each file, in
// the order JavaScript's localeCompare sorts the paths.
func folderHash(files []skills.File) string {
	sorted := slices.Clone(files)
	slices.SortFunc(sorted, func(a, b skills.File) int { return collate(a.Path, b.Path) })
	h := sha256.New()
	for _, f := range sorted {
		h.Write([]byte(f.Path))
		h.Write(f.Content)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// collation is the order of the Unicode Collation Algorithm's root collation
// for the characters a file name is likely to hold: punctuation and symbols,
// then digits, then letters, whatever their case.  A space or a control
// character comes before all of them, and a character not here after.
const collation = "_-,;:!?.'\"()[]{}@*/\\&#%`^+<=>|~$0123456789abcdefghijklmnopqrstuvwxyz"

// collate compares two paths as localeCompare does, for the characters in
// collation: by their characters, ignoring case, and then, when they differ
// only in case, by it, lower case first.
func collate(a, b string) int {
	ka, kb := []rune(a), []rune(b)
	for i := 0; i < min(len(ka), len(kb)); i++ {
		if d := primary(ka[i]) - primary(kb[i]); d != 0 {
			return d
		}
	}
	if d := len(ka) - len(kb); d != 0 {
		return d
	}
	for i := range ka {
		if d := tertiary(ka[i]) - tertiary(kb[i]); d != 0 {
			return d
		}
	}
	return strings.Compare(a, b)
}

func primary(r rune) int {
	if r <= ' ' {
		return 0
	}
	if 'A' <= r && r <= 'Z' {
		r += 'a' - 'A'
	}
	if i := strings.IndexRune(collation, r); i >= 0 {
		return i + 1
	}
	return len(collation) + int(r)
}

func tertiary(r rune) int {
	if 'A' <= r && r <= 'Z' {
		return 1
	}
	return 0
}

// treeHash is the hash the user's lock file records of a skill: the hash git
// gives the tree of its files, which the skills CLI reads from GitHub.
func treeHash(files []skills.File) string {
	return hex.EncodeToString(tree(files))
}

// tree returns the hash of the git tree object that holds files, by their
// paths beneath it.
func tree(files []skills.File) []byte {
	type item struct {
		name, mode string
		hash       []byte
	}
	var items []item
	subdirs := map[string][]skills.File{}
	for _, f := range files {
		if dir, rest, ok := strings.Cut(f.Path, "/"); ok {
			if _, seen := subdirs[dir]; !seen {
				items = append(items, item{name: dir, mode: "40000"})
			}
			subdirs[dir] = append(subdirs[dir], skills.File{Path: rest, Executable: f.Executable, Content: f.Content})
			continue
		}
		mode := "100644"
		if f.Executable {
			mode = "100755"
		}
		items = append(items, item{name: f.Path, mode: mode, hash: gitObject("blob", f.Content)})
	}
	for i := range items {
		if items[i].mode == "40000" {
			items[i].hash = tree(subdirs[items[i].name])
		}
	}
	// git orders a tree's entries by name, a directory's as if it ended in /.
	key := func(it item) string {
		if it.mode == "40000" {
			return it.name + "/"
		}
		return it.name
	}
	slices.SortFunc(items, func(a, b item) int { return strings.Compare(key(a), key(b)) })
	var body bytes.Buffer
	for _, it := range items {
		fmt.Fprintf(&body, "%s %s\x00", it.mode, it.name)
		body.Write(it.hash)
	}
	return gitObject("tree", body.Bytes())
}

// gitObject returns the hash git gives an object of a kind, such as a blob.
func gitObject(kind string, content []byte) []byte {
	h := sha1.New()
	fmt.Fprintf(h, "%s %d\x00", kind, len(content))
	h.Write(content)
	return h.Sum(nil)
}
