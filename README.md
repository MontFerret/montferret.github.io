# Ferret Website

Built with Hugo

## Prerequisites
- [Go](https://golang.org/doc/install)
- [Hugo](https://gohugo.io/getting-started/installing/)
- [Node.js](https://nodejs.org/en/download/)
- [Mage](https://magefile.org) for building the website.
- [frep](https://github.com/subchen/frep/releases)(optional) is the tool used for templating the documentation

## Getting Started
### Installing dependencies
```bash
go mod tidy && mage install
```

### Dev server
```bash
mage serve
```

The standard development server does not generate the Pagefind search bundle.

### Production build
Build the static site and its client-side search index:

```bash
mage build
```

The generated site, including Pagefind assets under `public/pagefind/`, is written to `public/`.

### Search preview
Build the site and search index, then serve the generated static output with Pagefind's preview server:

```bash
mage serveSearch
```

### Generating API docs
Generating `stdlib` documentation requires doc rep YAML.

```bash
mage generate
```

### Deployment

GitHub Actions builds, tests, and deploys the site to GitHub Pages when changes
are pushed to `main`. Pull requests targeting `main` run the same checks without
deploying, and maintainers can also start the workflow manually from `main` in
the Actions tab.

The Pages source must be set to **GitHub Actions** in the repository settings.
The custom domain remains `ferretlang.org`; configure it and HTTPS in the Pages
settings rather than relying on the generated `CNAME` file alone.
