---
title: "Getting started"
weight: 3
draft: false
---

# Quick start
### Browserless mode
If you want to play with FQL and check its syntax, you can run CLI with the following commands to run Ferret CLI in REPL mode:

```bash
$ ferret
```

```bash
Welcome to Ferret REPL
Please use `exit` or `Ctrl-D` to exit this program.
>%
>LET doc = DOCUMENT('https://news.ycombinator.com/')
>FOR post IN ELEMENTS(doc, '.storylink')
>RETURN post.attributes.href
>%
```

**NOTE**: symbol % is used to start and end multi-line queries. You also can use the heredoc format.

If you want to execute a query stored in a file, just pass a file name:

```bash
$ ferret ./docs/examples/static-page.fql
```

```bash
$ cat ./docs/examples/static-page.fql | ferret
```

```bash
$ ferret < ./docs/examples/static-page.fql
```

### Browser mode
By default, Ferret loads HTML pages via HTTP protocol, because it's faster.
But nowadays, there are more and more websites rendered with JavaScript, and therefore, this 'old school' approach does not really work.
For such cases, you may fetch documents using Chrome or Chromium via Chrome DevTools protocol (aka CDP).
First, you need to make sure that you launched a Chrome instance with open debugging port.

You can use your local installation:

```bash
$ chrome.exe --remote-debugging-port=9222 --headless
```

Or use [Docker](https://www.docker.com/):

```bash
$ docker pull montferret/chromium
$ docker run -d -p 9222:9222 montferret/chromium
```

Second, you need to pass the address to Ferret CLI.

```bash
$ ferret --cdp http://127.0.0.1:9222
```

**NOTE**: By default, Ferret will try to use this local address as a default one, so it makes sense to explicitly pass the parameter only in case of either different port number or remote address.

Alternatively, you can tell CLI to launch Chrome for you.

```bash
$ ferret --cdp-launch
```

Once Ferret knows how to communicate with Chrome, you can use ``DOCUMENT(url, opts)`` function with true boolean value for dynamic pages:

```bash
Welcome to Ferret REPL
Please use `exit` or `Ctrl-D` to exit this program.
>%
>LET doc = DOCUMENT('https://soundcloud.com/charts/top', { driver: "cdp" })
>WAIT_ELEMENT(doc, '.chartTrack__details', 5000)
>LET tracks = ELEMENTS(doc, '.chartTrack__details')
>FOR track IN tracks
>    LET username = ELEMENT(track, '.chartTrack__username')
>    LET title = ELEMENT(track, '.chartTrack__title')
>    RETURN {
>       artist: username.innerText,
>        track: title.innerText
>    }
>%
```