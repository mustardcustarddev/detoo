Draws the fire drills and the links between them, in Mermaid format, or in
Graphviz format with --graphviz.  The graph can be piped into mmdc or dot to
make an image.

Each node in the graph is a fire drill, and links to the fire drill's
document.  Dotted arrows chain the fire drills in order, and each link between
two fire drills is drawn once, from the fire drill that made it.

When there are no fire drills, nothing is written.
