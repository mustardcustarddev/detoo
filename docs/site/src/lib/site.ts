/** What the site says about itself, wherever it says it. */
export const SITE = {
  name: 'detoo',
  description:
    'Command-line tools and agent skills for keeping architecture decision records, requests for comments, user stories, problem playbooks and standard operating procedures.',
  repo: 'https://github.com/mustardcustarddev/detoo',
  /** The tools, in the order the site lists them. */
  tools: ['adr', 'rfc', 'sty', 'plm', 'sop', 'detoo'] as const,
  /**
   * The tools that were renamed, by their former names, so that a post
   * written before the rename still points to the tool.  astro.config.mjs
   * redirects their pages.
   */
  renamed: { pse: 'plm' } as const,
};

export type Tool = (typeof SITE.tools)[number];
export type FormerTool = keyof typeof SITE.renamed;

/** current returns the name a tool goes by now. */
export function current(tool: Tool | FormerTool): Tool {
  return tool in SITE.renamed ? SITE.renamed[tool as FormerTool] : (tool as Tool);
}
