// check-links finds the internal links and assets in the built site that lead
// nowhere.
//
// The site is served from a subdirectory on GitHub Pages, such as /detoo/, so
// a link written from the root, such as /blog, works with `astro dev` and
// breaks once deployed.  This walks every page in dist, resolves each link and
// asset against the base path, and fails if any is missing, or is written
// from the root without the base.
//
// Usage: node scripts/check-links.mjs  (after `astro build`)

import { readFileSync, readdirSync, statSync, existsSync } from 'node:fs';
import { join, relative, dirname, posix } from 'node:path';

const dist = new URL('../dist/', import.meta.url).pathname;
const trimmed = (process.env.BASE_PATH ?? '/detoo').replace(/^\/+|\/+$/g, '');
const base = trimmed ? `/${trimmed}/` : '/';
// The site's own origin: absolute URLs to it, such as a page's og:image, are
// checked as internal links.
const origin = new URL(process.env.SITE_URL ?? 'https://mustardcustarddev.github.io').origin;

function* files(dir, extensions) {
  for (const name of readdirSync(dir)) {
    const path = join(dir, name);
    if (statSync(path).isDirectory()) yield* files(path, extensions);
    else if (extensions.some((e) => name.endsWith(e))) yield path;
  }
}

// exists reports whether a path under dist is served: as a file, or as a
// directory with an index.html.
function exists(path) {
  const file = join(dist, decodeURIComponent(path));
  if (existsSync(file) && statSync(file).isFile()) return true;
  return existsSync(join(file, 'index.html')) || existsSync(file.replace(/\/$/, '') + '.html');
}

const problems = [];
for (const page of files(dist, ['.html'])) {
  const html = readFileSync(page, 'utf8');
  // The page's own URL path, under dist, such as docs/adr/new/
  const here = '/' + relative(dist, dirname(page)).split('\\').join('/') + '/';
  const links = [
    ...[...html.matchAll(/\s(href|src)="([^"]*)"/g)].map(([, attr, link]) => [attr, link]),
    // Social cards and search engines read these, and nothing on the page links to them.
    ...[...html.matchAll(/<meta (?:property|name)="((?:og|twitter):(?:image|url))" content="([^"]*)"/g)].map(
      ([, attr, link]) => [attr, link.startsWith(origin + '/') ? link.slice(origin.length) : link],
    ),
  ];
  for (const [attr, link] of links) {
    if (/^([a-z]+:|\/\/|#|$)/i.test(link)) continue; // external, fragment or empty
    const path = link.split(/[?#]/)[0];
    let target;
    if (path.startsWith('/')) {
      if (!(path + '/').startsWith(base)) {
        problems.push(`${relative(dist, page)}: ${attr}="${link}" is not under ${base}`);
        continue;
      }
      target = path.slice(base.length);
    } else {
      target = posix.normalize(posix.join(here, path)).replace(/^\//, '');
    }
    if (!exists(target)) problems.push(`${relative(dist, page)}: ${attr}="${link}" leads nowhere`);
  }
}

// The Markdown versions of the pages are read out of context, pasted into a
// chat or fetched by a model, so every link in them must be absolute, and
// those to this site must lead somewhere.
for (const file of files(dist, ['.md', '.txt'])) {
  const text = readFileSync(file, 'utf8');
  let fence = '';
  for (const line of text.split('\n')) {
    const opening = line.match(/^\s*(`{3,}|~{3,})/);
    if (fence) {
      if (opening && opening[1][0] === fence[0] && opening[1].length >= fence.length) fence = '';
      continue;
    }
    if (opening) {
      fence = opening[1];
      continue;
    }
    for (const [, link] of line.matchAll(/\]\(([^)\s]+)/g)) {
      if (link.startsWith('#') || (/^[a-z]+:/i.test(link) && !link.startsWith(origin + '/'))) continue;
      if (!link.startsWith(origin + '/')) {
        problems.push(`${relative(dist, file)}: (${link}) is not absolute`);
        continue;
      }
      const path = link.slice(origin.length).split(/[?#]/)[0];
      if (!(path + '/').startsWith(base) || !exists(path.slice(base.length))) {
        problems.push(`${relative(dist, file)}: (${link}) leads nowhere`);
      }
    }
  }
}

if (problems.length) {
  const unique = [...new Set(problems)];
  console.error(unique.join('\n'));
  console.error(`\n${unique.length} broken links`);
  process.exit(1);
}
console.log('All internal links resolve.');
