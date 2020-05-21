---
title: "Ferret v0.11.0"
subtitle: "More features"
draft: false
author: "Vladimir Fetisov"
authorLink: "https://github.com/3timeslazy"
date: "2020-05-19"
---
# COVID-19
Today, there are almost 5 million infected COVID-19 in the world. This disease kills thousands of people every day.

We hope you are healthy and safe.

But just in case: stay home 😷

---

Hello fellow miners, **[Ferret v0.11.0](https://github.com/MontFerret/ferret/releases/tag/v0.11.0)** has been **finally** released!

In this release, we mostly focused on our full-time jobs and families.

But despite this, we have done a lot of cool features and projects.

Let’s take a look!

---

# What's added
## USE statement
There is no need to write a full function path anymore! Just add ``USE PATH::TO::PACKAGE`` at the top of your query.
 
{{< code fql >}}
USE IO::FS

LET favicon = DOWNLOAD('https://www.google.com/favicon.ico')

// 'WRITE' the same as 'IO::FS::WRITE'.
// RETURN IO::FS::WRITE('google.favicon.ico', favicon) also valid.
RETURN WRITE('google.favicon.ico', favicon)
{{</ code >}}

## PATH functions
New functions for working with file paths.

{{< code fql >}}
LET filename = 'main.go'

RETURN PATH::EXT(filename) == '.go'
{{</ code >}}

---

# What's changed
## Scroll options
Scroll document's window and elements according to **[MDN](https://developer.mozilla.org/en-US/docs/Web/API/Element/scrollIntoView)**

{{< code fql >}}
LET doc = DOCUMENT(@url, { driver: 'cdp' })

SCROLL_BOTTOM(doc, {
    behavior: 'auto',
    block: 'end',
    inline: 'end'
})

RETURN 1
{{</ code >}}

---

# What's fixed
## RANDOM_TOKEN
``RANDOM_TOKEN`` return random string now! (even on Windows)

## IO::FS::WRITE
Now files are created with read access automatically.

---
# Around the Ferret
## OSS-Fuzz platform
We have tried to connect Ferret with [OSS-Fuzz](https://github.com/google/oss-fuzz) platform. But Google has [rejected](https://github.com/google/oss-fuzz/pull/3782) Ferret :(
