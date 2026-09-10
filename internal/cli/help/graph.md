Draws the decision records and the links between them, in Mermaid format, or
in Graphviz format with --graphviz.  The graph can be piped into mmdc or dot to
make an image.

Each node in the graph is a decision record, and links to the record's
document.  Dotted arrows chain the records in order, and each link between two
records is drawn once, from the record that made it.

When there are no records, nothing is written.
