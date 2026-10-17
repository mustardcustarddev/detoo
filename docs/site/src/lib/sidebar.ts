/**
 * The documentation, as the sidebar lists it: hand-written pages and the
 * generated command reference, in one order.  The sidebar, the routes and
 * the previous and next links all come from here, so they cannot disagree,
 * and a page added to either collection appears in all three.
 */
import { getCollection, type CollectionEntry } from 'astro:content';
import { SITE } from './site';
import { url } from './url';

export type Entry = CollectionEntry<'docs'> | CollectionEntry<'reference'>;

export interface Page {
  /** The path under docs/, such as "adr/new" or "install"; "" for the introduction. */
  slug: string;
  href: string;
  /** What the sidebar calls it. */
  label: string;
  /** The label of the group it is in. */
  group: string;
  entry: Entry;
}

export interface Group {
  label: string;
  /** Whether the label is a command, and so set in monospace. */
  code: boolean;
  pages: Page[];
}

const sections = {
  'getting-started': 'Getting started',
  guides: 'Guides',
  skills: 'Agent skills',
} as const;

/** slugOf turns an entry's id into its path under docs/: an index is its directory. */
export function slugOf(id: string): string {
  return id.replace(/(^|\/)index$/, '');
}

function page(entry: Entry, group: string, label: string): Page {
  const slug = slugOf(entry.id);
  return { slug, href: url(slug ? `docs/${slug}/` : 'docs/'), label, group, entry };
}

/** sidebar returns the groups of pages, in order. */
export async function sidebar(): Promise<Group[]> {
  const docs = await getCollection('docs');
  const reference = await getCollection('reference');
  const byOrder = (a: Entry, b: Entry) => a.data.order - b.data.order;

  const section = (key: keyof typeof sections): Group => ({
    label: sections[key],
    code: false,
    pages: docs
      .filter((d) => d.data.section === key)
      .sort(byOrder)
      .map((d) => page(d, sections[key], d.data.title)),
  });

  const tool = (name: string): Group => ({
    label: name,
    code: true,
    pages: reference
      .filter((r) => r.data.tool === name)
      .sort(byOrder)
      .map((r) => page(r, name, r.data.command || 'Overview')),
  });

  return [section('getting-started'), ...SITE.tools.map(tool), section('guides'), section('skills')].filter(
    (g) => g.pages.length > 0,
  );
}

/** pages returns every page of the documentation, in the sidebar's order. */
export async function pages(): Promise<Page[]> {
  return (await sidebar()).flatMap((g) => g.pages);
}
