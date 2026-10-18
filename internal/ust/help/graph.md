Draws the user stories and the links between them, in Mermaid format, or in
Graphviz format with --graphviz.  The graph can be piped into mmdc or dot to
make an image.

Each node in the graph is a story, and links to the story's document.  Dotted
arrows chain the stories in order, and each link between two stories is drawn
once, from the story that made it.

When there are no stories, nothing is written.
