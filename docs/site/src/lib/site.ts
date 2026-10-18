/** What the site says about itself, wherever it says it. */
export const SITE = {
  name: 'detoo',
  description:
    'Command-line tools and agent skills for keeping architecture decision records, requests for comments, user stories, problem playbooks and standard operating procedures.',
  repo: 'https://github.com/mustardcustarddev/detoo',
  /** The tools, in the order the site lists them. */
  tools: ['adr', 'rfc', 'ust', 'plm', 'sop', 'detoo'] as const,
};

export type Tool = (typeof SITE.tools)[number];
