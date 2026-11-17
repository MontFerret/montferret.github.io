import assert from "node:assert/strict";
import { glob, readFile } from "node:fs/promises";
import path from "node:path";

const output = path.join(process.cwd(), "public");
const sitemap = await readFile(path.join(output, "sitemap.xml"), "utf8");
const sitemapLinks = [...sitemap.matchAll(/<loc>([^<]+)<\/loc>/g)].map((match) => match[1]);
assert.ok(sitemapLinks.length > 0, "Sitemap contains no routes");
const origin = new URL(sitemapLinks[0]).origin;
const pages = new Map();

// These expressions read the quoted attributes emitted by the theme and Hugo.
function attributes(tag) {
    return Object.fromEntries([...tag.matchAll(/([\w:-]+)\s*=\s*(?:"([^"]*)"|'([^']*)')/g)]
        .map((match) => [match[1], match[2] ?? match[3]]));
}

function route(href, source) {
    const url = new URL(href.replaceAll("&amp;", "&"), new URL(source, origin));
    assert.equal(url.origin, origin, `${source} links outside the site: ${href}`);
    assert.ok(!url.pathname.includes("//"), `${source} has a malformed route: ${href}`);
    return url.pathname;
}

async function readPage(href, source) {
    const pathname = route(href, source);
    if (!pages.has(pathname)) {
        const file = path.join(output, decodeURIComponent(pathname), "index.html");
        try {
            pages.set(pathname, await readFile(file, "utf8"));
        } catch (error) {
            throw new Error(`${source} links to a missing page: ${href}`, { cause: error });
        }
    }
    return pages.get(pathname);
}

const sitemapRoutes = new Set();
for (const href of sitemapLinks) {
    const pathname = route(href, "/sitemap.xml");
    assert.ok(!sitemapRoutes.has(pathname), `Duplicate sitemap route: ${pathname}`);
    sitemapRoutes.add(pathname);
    await readPage(href, "/sitemap.xml");
}
for (const pathname of ["/blog/", "/docs/", "/try/", "/docs/standard-library/"]) {
    assert.ok(sitemapRoutes.has(pathname), `Sitemap is missing ${pathname}`);
}

const articleRoutes = new Set();
const listedArticles = new Set();
let archiveCount = 0;
let docsLinkCount = 0;

for await (const file of glob(["blog/**/index.html", "docs/**/index.html"], { cwd: output })) {
    const pathname = `/${file.replaceAll(path.sep, "/").replace(/index\.html$/, "")}`;
    const html = await readPage(pathname, pathname);
    const isArchive = /<ul\b[^>]*class="[^"]*\bblog-list\b/.test(html);
    if (isArchive) {
        archiveCount++;
    }
    if (/<article\b[^>]*class="blog-post"/.test(html)) {
        articleRoutes.add(pathname);
        const content = html.match(/<div class="blog-post-content">([\s\S]*?)<nav class="blog-post-nav"/);
        assert.ok(content?.[1].replace(/<[^>]*>/g, "").trim(), `Empty article body: ${pathname}`);
    }

    let listedCount = 0;
    for (const [tag] of html.matchAll(/<a\b[^>]*>/g)) {
        const attrs = attributes(tag);
        const classes = new Set((attrs.class ?? "").split(/\s+/));
        if (!attrs.href) {
            continue;
        }
        if (isArchive && classes.has("read-more")) {
            const target = route(attrs.href, pathname);
            assert.ok(target.startsWith("/blog/"), `Article link escaped /blog/: ${attrs.href}`);
            assert.ok(!listedArticles.has(target), `Article listed more than once: ${target}`);
            listedArticles.add(target);
            listedCount++;
            await readPage(attrs.href, pathname);
        }
        if (isArchive && (classes.has("pagination-previous") || classes.has("pagination-next"))) {
            const target = route(attrs.href, pathname);
            assert.match(target, /^\/blog\/(?:page\/\d+\/)?$/, `Invalid blog pagination: ${attrs.href}`);
            await readPage(attrs.href, pathname);
        }
        if (classes.has("docs-nav-link") || classes.has("docs-section-card") || classes.has("docs-card")) {
            // Cards may link to other site sections or external documentation.
            if (classes.has("docs-card") && /^https?:\/\//.test(attrs.href) && new URL(attrs.href).origin !== origin) {
                continue;
            }
            const target = route(attrs.href, pathname);
            if (!classes.has("docs-card")) {
                assert.ok(target.startsWith("/docs/"), `Docs link escaped /docs/: ${attrs.href}`);
            }
            await readPage(attrs.href, pathname);
            docsLinkCount++;
        }
    }
    if (isArchive) {
        assert.ok(listedCount > 0, `Empty blog archive: ${pathname}`);
    }
}

assert.ok(archiveCount > 0, "Blog archive was not rendered");
assert.ok(articleRoutes.size > 0, "Blog articles were not rendered");
assert.deepEqual(listedArticles, articleRoutes, "Blog archives must link to every rendered article exactly once");
assert.ok(docsLinkCount > 0, "Docs navigation was not rendered");

const firstPage = await readPage("/blog/page/1/", "/blog/");
const refreshTag = [...firstPage.matchAll(/<meta\b[^>]*>/g)]
    .map(([tag]) => attributes(tag))
    .find((attrs) => attrs["http-equiv"]?.toLowerCase() === "refresh");
const redirect = refreshTag?.content?.match(/\burl=(.+)$/i)?.[1];
assert.ok(redirect, "First blog page has no redirect");
assert.equal(route(redirect, "/blog/page/1/"), "/blog/", "First blog page redirects outside /blog/");

for (const file of ["index.xml", "blog/index.xml"]) {
    const feed = await readFile(path.join(output, file), "utf8");
    const links = [...feed.matchAll(/<link>([^<]+)<\/link>/g)].map((match) => match[1]);
    assert.ok(links.length > 1, `${file} contains no feed entries`);
    assert.equal(route(links[0], `/${file}`), file === "index.xml" ? "/" : "/blog/", `Incorrect feed URL in ${file}`);
    for (const href of links) {
        if (file === "blog/index.xml") {
            assert.ok(route(href, `/${file}`).startsWith("/blog/"), `Blog RSS link escaped /blog/: ${href}`);
        }
        await readPage(href, `/${file}`);
    }
}

console.log(`Verified ${sitemapRoutes.size} sitemap routes, ${articleRoutes.size} blog articles, ${archiveCount} archives, and ${docsLinkCount} docs links.`);
