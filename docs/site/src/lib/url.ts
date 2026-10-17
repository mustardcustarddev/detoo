/**
 * url returns the path to a page or file of this site, under the base path
 * the site is served from, such as /detoo/ on GitHub Pages.
 *
 * Every internal link goes through it: a link written from the root, such as
 * "/blog", works with `astro dev` but leads nowhere once deployed.  A URL with
 * a scheme, or a fragment, is returned as it is.
 */
export function url(path = ''): string {
  if (/^([a-z]+:|\/\/|#)/i.test(path)) return path;
  const base = import.meta.env.BASE_URL.replace(/\/+$/, '');
  return `${base}/${path.replace(/^\/+/, '')}`;
}

/** absolute returns the full URL of a path of this site, for metadata. */
export function absolute(path: string, site: URL | undefined): string {
  return new URL(url(path), site ?? 'https://mustardcustarddev.github.io').toString();
}
