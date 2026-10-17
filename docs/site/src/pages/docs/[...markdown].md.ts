import type { APIRoute, GetStaticPaths } from 'astro';
import { pages, type Page } from '../../lib/sidebar';
import { markdownPath, pageMarkdown } from '../../lib/markdown';

/**
 * Each page of the documentation as Markdown, beside the page: docs/adr/new.md.
 *
 * The .md is in the file's name, not the parameter, so that Astro knows the
 * route ends in an extension and serves it without the trailing slash it
 * requires of pages; otherwise `astro dev` answers only docs/adr/new.md/.
 */
export const getStaticPaths = (async () => {
  return (await pages()).map((page) => ({
    params: { markdown: markdownPath(page.slug).replace(/^docs\//, '').replace(/\.md$/, '') },
    props: { page },
  }));
}) satisfies GetStaticPaths;

export const GET: APIRoute<{ page: Page }> = ({ props, site }) =>
  new Response(pageMarkdown(props.page, site), {
    headers: { 'Content-Type': 'text/markdown; charset=utf-8' },
  });
