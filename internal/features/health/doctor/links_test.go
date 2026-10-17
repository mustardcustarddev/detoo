package doctor

import "testing"

// where is a fix that moved doc/rfc to docs/rfc, and renamed two documents.
func where(name string) string {
	files := map[string]string{
		"docs/adr/0001-first.md":   "docs/adr/ADR-0001-first.md",
		"doc/rfc/0001-proposal.md": "docs/rfc/RFC-0001-proposal.md",
	}
	if to, ok := files[name]; ok {
		return to
	}
	return swapPrefix(name, map[string]string{"doc/rfc": "docs/rfc"})
}

func TestRelink(t *testing.T) {
	for _, test := range []struct {
		name, from, to, content, want string
		changed                       int
	}{
		{
			name: "a link to a renamed document in the same directory",
			from: "docs/adr/0002-second.md", to: "docs/adr/ADR-0002-second.md",
			content: "Supersedes [1. First](0001-first.md)\n",
			want:    "Supersedes [1. First](ADR-0001-first.md)\n", changed: 1,
		},
		{
			name: "a link from the root, with a heading and a title",
			from: "README.md", to: "README.md",
			content: "See [the first](docs/adr/0001-first.md#context \"First\") and [an RFC](doc/rfc/0001-proposal.md).\n",
			want:    "See [the first](docs/adr/ADR-0001-first.md#context \"First\") and [an RFC](docs/rfc/RFC-0001-proposal.md).\n", changed: 2,
		},
		{
			name: "a link from a document that moved, to one that did not",
			from: "doc/rfc/0001-proposal.md", to: "docs/rfc/RFC-0001-proposal.md",
			content: "Builds on [ADR 2](../adr/0002-second.md) and [ADR 3](../../docs/adr/0003-third.md).\n",
			want:    "Builds on [ADR 2](../../doc/adr/0002-second.md) and [ADR 3](../../docs/adr/0003-third.md).\n", changed: 1,
		},
		{
			name: "a reference link, a link from the project's root, and a directory",
			from: "README.md", to: "README.md",
			content: "[first]: docs/adr/0001-first.md\nThe [RFCs](/doc/rfc/) and [the list](doc/rfc/).\n",
			want:    "[first]: docs/adr/ADR-0001-first.md\nThe [RFCs](/docs/rfc/) and [the list](docs/rfc/).\n", changed: 3,
		},
		{
			name: "links to other sites, headings, and in code, are left alone",
			from: "README.md", to: "README.md",
			content: "[site](https://example.com/doc/rfc) [here](#doc) [mail](mailto:a@b.c)\n```\n[code](doc/rfc/0001-proposal.md)\n```\n",
			want:    "[site](https://example.com/doc/rfc) [here](#doc) [mail](mailto:a@b.c)\n```\n[code](doc/rfc/0001-proposal.md)\n```\n",
		},
		{
			name: "an image in a directory that moved",
			from: "docs/adr/ADR-0002-second.md", to: "docs/adr/ADR-0002-second.md",
			content: "![Diagram](../../doc/rfc/diagram.png)\n",
			want:    "![Diagram](../rfc/diagram.png)\n", changed: 1,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, changed := Relink(test.content, test.from, test.to, where)
			if got != test.want || changed != test.changed {
				t.Errorf("Relink() = %q, %d\nwant         %q, %d", got, changed, test.want, test.changed)
			}
		})
	}
}

func TestRelative(t *testing.T) {
	for _, test := range []struct{ dir, target, want string }{
		{"docs/adr", "docs/adr/ADR-0001-first.md", "ADR-0001-first.md"},
		{"docs/rfc", "docs/adr/ADR-0001-first.md", "../adr/ADR-0001-first.md"},
		{".", "docs/rfc", "docs/rfc"},
		{"docs/rfc", ".", "../.."},
		{"/srv/adr", "/srv/rfc/RFC-0001-x.md", "../rfc/RFC-0001-x.md"},
	} {
		if got, ok := relative(test.dir, test.target); !ok || got != test.want {
			t.Errorf("relative(%q, %q) = %q, %v, want %q", test.dir, test.target, got, ok, test.want)
		}
	}
	if _, ok := relative("docs", "/srv/adr"); ok {
		t.Errorf("relative() leads from a relative path to an absolute one")
	}
}
