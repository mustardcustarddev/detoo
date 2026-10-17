import rss from '@astrojs/rss';
import { getCollection } from 'astro:content';
import { SITE } from '../lib/site';
import { url } from '../lib/url';

export async function GET(context) {
  const posts = (await getCollection('blog', (post) => !post.data.draft)).sort(
    (a, b) => b.data.pubDate.valueOf() - a.data.pubDate.valueOf(),
  );

  return rss({
    title: `${SITE.name} blog`,
    description: 'Why and how software teams record their decisions, plan for incidents, and agree on change before they make it.',
    // The channel links to the blog, under the base path, not to the domain.
    site: new URL(url('blog/'), context.site).toString(),
    items: posts.map((post) => ({
      title: post.data.title,
      pubDate: post.data.pubDate,
      description: post.data.description,
      link: url(`blog/${post.id}/`),
    })),
    customData: '<language>en-gb</language>',
  });
}
