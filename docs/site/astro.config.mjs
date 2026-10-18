// @ts-check
import { defineConfig } from 'astro/config';

import tailwindcss from '@tailwindcss/vite';
import react from '@astrojs/react';
import mdx from '@astrojs/mdx';
import sitemap from '@astrojs/sitemap';

// The site is served by GitHub Pages from a subdirectory named after the
// repository.  SITE_URL and BASE_PATH override them, for a custom domain or a
// fork.  Every internal link goes through url() in src/lib/url.ts so that it
// carries the base.
const site = process.env.SITE_URL ?? 'https://mustardcustarddev.github.io';
const base = process.env.BASE_PATH ?? '/detoo';

// https://astro.build/config
export default defineConfig({
  site,
  base,
  trailingSlash: 'always',
  output: 'static',

  vite: {
    plugins: [tailwindcss()],
  },

  integrations: [react(), mdx(), sitemap()],

  // Pages that moved, so that links to them, from elsewhere and from posts
  // written before the move, still lead somewhere.  pse became plm, and its
  // skill, create-piosee, became create-playbook.
  redirects: {
    ...Object.fromEntries(
      ['', 'init/', 'new/', 'link/', 'list/', 'graph/', 'config/'].map((page) => [
        `/docs/pse/${page}`,
        `${base}/docs/plm/${page}`,
      ])
    ),
    '/docs/skills/create-piosee/': `${base}/docs/skills/create-playbook/`,
  },

  prefetch: {
    prefetchAll: true,
  },

  build: {
    inlineStylesheets: 'always',
  },

  markdown: {
    shikiConfig: {
      themes: { light: 'github-light', dark: 'github-dark' },
    },
  },
});
