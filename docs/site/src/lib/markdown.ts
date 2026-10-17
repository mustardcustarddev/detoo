/**
 * The documentation as Markdown, for language models and anyone who would
 * rather read it as plain text: each page on its own, and all of them in one
 * file, llms-full.txt.
 *
 * A page's Markdown is its source, the same text the site renders, with its
 * title and description above it.  Its links are made absolute, since a link
 * such as ../install/ means nothing once the text is pasted somewhere else, and
 * its images become their descriptions.
 */
import type { Page } from './sidebar';
import { absolute } from './url';

/** markdownPath returns the path under the base of a page's Markdown. */
export function markdownPath(slug: string): string {
  return slug ? `docs/${slug}.md` : 'docs/index.md';
}

/** pageMarkdown returns a page of the documentation as Markdown. */
export function pageMarkdown(page: Page, site: URL | undefined): string {
  const here = absolute(page.slug ? `docs/${page.slug}/` : 'docs/', site);
  const { title, description } = page.entry.data;
  const body = absoluteLinks(page.entry.body ?? '', here).trim();
  return `# ${title}\n\n> ${description}\n\nURL: ${here}\n\n${body}\n`;
}

/**
 * absoluteLinks resolves the relative links in markdown against base, and
 * replaces each image by its description.  Code blocks are left as they are,
 * since what looks like a link there is code.
 */
export function absoluteLinks(markdown: string, base: string): string {
  let fence = '';
  return markdown
    .split('\n')
    .map((line) => {
      const opening = line.match(/^\s*(`{3,}|~{3,})/);
      if (fence) {
        if (opening && opening[1][0] === fence[0] && opening[1].length >= fence.length) fence = '';
        return line;
      }
      if (opening) {
        fence = opening[1];
        return line;
      }
      return line
        .replace(/!\[([^\]]*)\]\([^)]*\)/g, (_, alt: string) => `(Image: ${alt})`)
        .replace(/\]\(([^)\s]+)((?:\s+"[^"]*")?)\)/g, (match, target: string, title: string) =>
          /^([a-z]+:|#)/i.test(target) ? match : `](${new URL(target, base)}${title})`,
        );
    })
    .join('\n');
}
