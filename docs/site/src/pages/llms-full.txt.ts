import type { APIRoute } from 'astro';
import { pages } from '../lib/sidebar';
import { pageMarkdown } from '../lib/markdown';
import { absolute } from '../lib/url';
import { SITE } from '../lib/site';

/**
 * llms-full.txt is the whole of the documentation as one Markdown file, in the
 * order of the contents, for a language model to read in one go.
 */
export const GET: APIRoute = async ({ site }) => {
  const all = await pages();
  const intro = [
    `# ${SITE.name} documentation`,
    '',
    `> ${SITE.description}`,
    '',
    `Every page of the documentation at ${absolute('docs/', site)}, in the order of its contents.`,
  ].join('\n');
  // A blank line before each rule, so that it is not read as a heading's underline.
  const body = [intro, ...all.map((page) => pageMarkdown(page, site).trim())].join('\n\n---\n\n');
  return new Response(body + '\n', { headers: { 'Content-Type': 'text/plain; charset=utf-8' } });
};
