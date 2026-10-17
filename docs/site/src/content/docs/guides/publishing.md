---
title: Publishing a table of contents and graphs
description: Write a Markdown table of contents for a log, and draw the links between its documents with Mermaid or Graphviz.
section: guides
order: 4
---

Every record tool can publish its log as a table of contents, and draw the
links between documents as a graph. The examples use `adr`; the others are
the same.

## A table of contents

```bash
adr list --toc > docs/adr/README.md
```

GitHub shows a directory's `README.md` beneath its files, so this puts an
index of the decisions beside them.

To publish one elsewhere, prefix each link so that it still reaches the
records, and add your own introduction and conclusion:

```bash
adr list --toc --intro intro.md --outro outro.md --prefix ../docs/adr/ > site/decisions.md
```

## A graph

`graph` writes the documents and the links between them, such as
"supersedes", in Mermaid by default:

```bash
adr graph > docs/adr/graph.mmd
```

GitHub draws a Mermaid diagram in Markdown inside a `mermaid` code block. To
make an image, use the [Mermaid CLI](https://github.com/mermaid-js/mermaid-cli):

```bash
sty graph | mmdc --input - --output docs/userstory/graph.svg
```

Or draw it with [Graphviz](https://graphviz.org/):

```bash
adr graph --graphviz | dot -Tsvg > graph.svg
```

Each node links to its document. `--extension` sets the extension of those
links, `.html` by default, and `--prefix` puts a path before them:

```bash
adr graph --graphviz --extension .pdf | dot -Tpdf > graph.pdf
```
