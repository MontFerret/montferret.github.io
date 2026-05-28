---
title: "Quick Start"
sidebarTitle: "Quick Start"
weight: 30
draft: false
description: "Run your first Ferret query against static and dynamic pages."
aliases:
  - /docs/getting-started/
---

# Quick Start
In this guide, you will run your first Ferret query and learn the basic shape of an FQL script.

## 1. Run Ferret without a browser
Start with a simple expression:
{{< tabs >}}

{{< tab title="Terminal" >}}
{{< terminal >}}
ferret run -e 'RETURN 1 + 1'
{{< /terminal >}}
{{< /tab >}}

{{< tab title="Browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
RETURN 1 + 1
{{< /editor >}}
{{< /tab >}}

{{< /tabs >}}


## 2. Work with structured data
{{< tabs >}}
{{< tab title="Terminal" >}}
{{< terminal command="true" >}}
ferret run -e '
LET user = {
    name: "Ada",
    roles: ["admin", "editor"]
}

RETURN {
    name: user.name,
    isAdmin: CONTAINS(user.roles, "admin")
}
'
{{< /terminal >}}
{{< /tab >}}

{{< tab title="Browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET user = {
    name: "Ada",
    roles: ["admin", "editor"]
}

RETURN {
    name: user.name,
    isAdmin: CONTAINS(user.roles, "admin")
}
{{< /editor >}}
{{< /tab >}}
{{< /tabs >}}

## 3. Query HTML

{{< tabs >}}
{{< tab title="Terminal" >}}
{{< terminal command="true" >}}
ferret run -e '
LET doc = DOCUMENT("https://mockery.montferret.dev")
RETURN doc[~ css`article`]
'
{{< /terminal >}}
{{< /tab >}}

{{< tab title="Browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev")
RETURN doc[~ css`article`]
{{< /editor >}}
{{< /tab >}}
{{< /tabs >}}

## 4. Use a browser

{{< tabs >}}
{{< tab title="Terminal" >}}
{{< terminal command="true" >}}
ferret run -e '
LET page = DOCUMENT("https://mockery.montferret.dev", { driver: "cdp" })
RETURN page.title
'
{{< /terminal >}}
{{< /tab >}}

{{< tab title="Browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET page = DOCUMENT("https://mockery.montferret.dev", { driver: "cdp" })
RETURN page.title
{{< /editor >}}
{{< /tab >}}
{{< /tabs >}}

## 5. Save a script
{{< terminal command="true" >}}
echo 'LET name = @name ?: "Ferret"
RETURN "Hello, " + name' > hello.fql
{{< /terminal >}}

{{< terminal command="true" >}}
ferret run hello.fql --param name=Steve
{{< /terminal >}}

## Where to go next

And this is where your tiles are perfect:

* CLI Usage
* Embedding Ferret in Go
* HTML Querying
* Browser Automation
* JavaScript Rendering
* Custom Modules
* Lab / Testing

So the distinction becomes:

Overview - what Ferret is, philosophy, use cases
Installation - how to get it
Quick Start - shortest successful hands-on path
CLI Usage - all CLI commands and flags
Language Guide - FQL syntax and semantics
Modules / Drivers - HTML, browser, JS rendering, custom integrations
