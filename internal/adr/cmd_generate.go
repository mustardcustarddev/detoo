package adr

import (
	"path"
	"strings"
)

const generateHelp = `usage: adr generate [REPORT [OPTION ...]]

Generates summary documentation about the architecture decision records
that is typically fed into the tool chain for publishing a project's
documentatation.

To list the report types that can be generated, run:

    adr generate

To get help on a specific report generator, run:

    adr help generate <report-type>

For example, the following command will run the "toc" (table of contents)
generator:

    adr generate toc

The following command will display help for the "toc" generator:

    adr help generate toc
`

// generators are the reports that "adr generate" can produce.
var generators = []*Command{
	{Name: "graph", Help: graphHelp, Run: runGraph},
	{Name: "toc", Help: tocHelp, Run: runTOC},
}

func runGenerate(e *Env, args []string) error {
	if len(args) == 0 {
		for _, generator := range visible(generators) {
			e.print(generator.Name)
		}
		return nil
	}
	generator := lookup(generators, args[0])
	if generator == nil {
		return &UsageError{Usage: "adr generate [" + strings.Join(names(visible(generators)), "|") + "] [OPTION ...]"}
	}
	return generator.Run(e, args[1:])
}

const graphHelp = `usage: adr generate graph [-p LINK_PREFIX] [-e LINK-EXTENSION]

Generates a visualisation of the links between decision records in
Graphviz format.  This can be piped into the graphviz tools to
generate a an image file.

Each node in the graph represents a decision record and is linked to
the decision record document.

Options:

-e LINK-EXTENSION
        the file extension of the documents to which generated links refer.
        Defaults to ` + "`.html`" + `.

-p LINK_PREFIX
        prefix each decision file link with LINK_PREFIX.

E.g. to generate a graph visualisation of decision records in SVG format:

    adr generate graph | dot -Tsvg > graph.svg

E.g. to generate a graph visualisation in PDF format, in which all links
are to .pdf files:

   adr generate graph -e .pdf | dot -Tpdf > graph.pdf
`

func runGraph(e *Env, args []string) error {
	linkPrefix, linkExtension := "", ".html"
	if _, err := parseOptions(args, "ep", func(option byte, value string) error {
		switch option {
		case 'e':
			linkExtension = value
		case 'p':
			linkPrefix = value
		}
		return nil
	}); err != nil {
		return err
	}
	records, err := e.List()
	if err != nil {
		return err
	}

	e.print("digraph {")
	e.print("  node [shape=plaintext];")
	e.print("  subgraph {")
	for _, record := range records {
		number := Number(record)
		title, err := e.Title(record)
		if err != nil {
			return err
		}
		name := strings.TrimSuffix(path.Base(record), ".md")
		e.printf("    _%d [label=%q; URL=%q];\n", number, title, linkPrefix+name+linkExtension)
		if number > 1 {
			e.printf("    _%d -> _%d [style=\"dotted\", weight=1];\n", number-1, number)
		}
	}
	e.print("  }")
	for _, record := range records {
		links, err := e.Links(record)
		if err != nil {
			return err
		}
		for _, link := range links {
			// Reverse links duplicate a link recorded by the other record.
			if strings.HasSuffix(link.Type, " by") {
				continue
			}
			e.printf("  _%d -> _%d [label=%q, weight=0]\n", Number(record), link.Target, link.Type)
		}
	}
	e.print("}")
	return nil
}

const tocHelp = `usage: adr generate toc [-i INTRO] [-o OUTRO] [-p LINK_PREFIX]

Generates a table of contents in Markdown format to stdout.

Options:

-i INTRO  precede the table of contents with the given INTRO text.
-o OUTRO  follow the table of contents with the given OUTRO text.
-p LINK_PREFIX
          prefix each decision file link with LINK_PREFIX.

Both INTRO and OUTRO must be in Markdown format.
`

func runTOC(e *Env, args []string) error {
	intro, outro, linkPrefix := "", "", ""
	if _, err := parseOptions(args, "iop", func(option byte, value string) error {
		switch option {
		case 'i':
			intro = value
		case 'o':
			outro = value
		case 'p':
			linkPrefix = value
		}
		return nil
	}); err != nil {
		return err
	}
	records, err := e.List()
	if err != nil {
		return err
	}

	e.print("# Architecture Decision Records")
	e.print("")
	if intro != "" {
		text, err := readFile(e.path(intro))
		if err != nil {
			return err
		}
		e.printf("%s\n", strings.TrimSuffix(text, "\n"))
		e.print("")
	}
	for _, record := range records {
		title, err := e.Title(record)
		if err != nil {
			return err
		}
		e.printf("* [%s](%s)\n", title, linkPrefix+path.Base(record))
	}
	if outro != "" {
		text, err := readFile(e.path(outro))
		if err != nil {
			return err
		}
		e.print("")
		e.printf("%s\n", strings.TrimSuffix(text, "\n"))
	}
	return nil
}
