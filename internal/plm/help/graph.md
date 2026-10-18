Draws the problem playbooks and the links between them, in Mermaid format, or
in Graphviz format with --graphviz.  The graph can be piped into mmdc or dot to
make an image.

Each node in the graph is a playbook, and links to the playbook's document.
Dotted arrows chain the playbooks in order, and each link between two playbooks
is drawn once, from the playbook that made it.

When there are no playbooks, nothing is written.
