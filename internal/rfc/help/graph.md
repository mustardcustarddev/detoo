Draws the RFCs and the links between them, in Mermaid format, or in Graphviz
format with --graphviz.  The graph can be piped into mmdc or dot to make an
image.

Each node in the graph is an RFC, and links to the RFC's document.  Dotted
arrows chain the RFCs in order, and each link between two RFCs is drawn once,
from the RFC that made it.

When there are no RFCs, nothing is written.
