Draws the SOPs and the links between them, in Mermaid format, or in Graphviz
format with --graphviz.  The graph can be piped into mmdc or dot to make an
image.

Each node in the graph is an SOP, and links to the SOP's document.  Dotted
arrows chain the SOPs in order, and each link between two SOPs is drawn once,
from the SOP that made it.

When there are no SOPs, nothing is written.
