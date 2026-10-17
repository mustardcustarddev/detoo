import type { APIRoute } from 'astro';
import { absolute } from '../lib/url';
import { SITE } from '../lib/site';
import { sidebar } from '../lib/sidebar';
import { markdownPath } from '../lib/markdown';

/**
 * llms.txt describes the site for language models, following llmstxt.org: a
 * summary, then links to the Markdown version of every page, grouped as the
 * contents are, and to llms-full.txt, which holds them all.  It is generated so
 * that its links follow SITE_URL and BASE_PATH, and every page added to the
 * docs.
 */
export const GET: APIRoute = async ({ site }) => {
  const groups = await sidebar();
  const body = [
    `# ${SITE.name}`,
    '',
    `> ${SITE.description}`,
    '',
    `The whole of the documentation is in one file: ${absolute('llms-full.txt', site)}`,
    '',
    ...groups.flatMap((group) => [
      `## ${group.label}`,
      '',
      ...group.pages.map(
        (page) => `- [${page.entry.data.title}](${absolute(markdownPath(page.slug), site)}): ${page.entry.data.description}`,
      ),
      '',
    ]),
    '## Optional',
    '',
    `- [Blog](${absolute('blog/', site)}): why and how software teams record their decisions`,
    `- [Source](${SITE.repo}): the tools, the agent skills and this site`,
    '',
  ].join('\n');
  return new Response(body, { headers: { 'Content-Type': 'text/plain; charset=utf-8' } });
};
