---
title: "Quick Start"
sidebarTitle: "Quick Start"
weight: 40
draft: false
description: "Run your first Ferret query against static and dynamic pages."
---

# Quick Start

In this guide, you will run your first Ferret queries and learn the basic shape of an FQL script.

Ferret scripts are written in FQL (Ferret Query Language). A script can define values, work with structured data, query documents, interact with web pages, and return a result. The examples below start with the smallest possible script and gradually introduce the pieces you will use in real-world extraction workflows.

## Hello world: a simple expression

Start with a simple expression:

{{< tabs >}}

{{< tab title="Terminal" >}}
{{< terminal >}}
ferret run -e 'return 1 + 1'
{{< /terminal >}}
{{< /tab >}}

{{< tab title="Try in browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
return 1 + 1
{{< /editor >}}
{{< /tab >}}

{{< /tabs >}}

This script does not open a page or use a browser. It evaluates an expression and returns the result.

`return` produces the script's final value. In this case, the returned value is the result of `1 + 1`; without an explicit return, a non-empty script completes with `none`.

## Work with structured data

Ferret can also work with objects, arrays, and function calls:

{{< tabs >}}
{{< tab title="Terminal" >}}
{{< terminal command="true" >}}
ferret run -e '
let user = {
    name: "Ada",
    roles: ["admin", "editor"]
}

return {
    name: user.name,
    isAdmin: contains(user.roles, "admin")
}
'
{{< /terminal >}}
{{< /tab >}}

{{< tab title="Try in browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
let user = {
    name: "Ada",
    roles: ["admin", "editor"]
}

return {
    name: user.name,
    isAdmin: contains(user.roles, "admin")
}
{{< /editor >}}
{{< /tab >}}
{{< /tabs >}}

This example defines a variable with `let`, creates an object, reads fields with dot notation, and returns a new object.

`let` creates a local binding. Here, user is an object with a name field and a roles array. The `return` expression builds a new object from that data and uses `contains` to check whether the user has the admin role.

## Query HTML

Now let’s load a page and query HTML elements from it:

{{< tabs >}}
{{< tab title="Terminal" >}}
{{< terminal command="true" >}}
ferret run -e '
let page = web::html::open("https://mockery.ferretlang.org")
return page[~ css`article`]
'
{{< /terminal >}}
{{< /tab >}}

{{< tab title="Try in browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
let page = web::html::open("https://mockery.ferretlang.org")
return page[~ css`article`]
{{< /editor >}}
{{< /tab >}}
{{< /tabs >}}

`web::html::open` loads the URL and returns an HTML page value. The expression `page[~ css'article']` queries that page using the CSS dialect and returns matching elements.

The `~` operator is FQL’s shorthand query operator. The full query syntax is covered later in the language reference.

HTML support is available in the CLI. When embedding Ferret in a Go application, HTML querying is provided through a module. JavaScript-defined modules cannot supply the Go HTML host values and drivers these examples require, so the HTML and browser examples are not available through `@montferret/ferret`.

## Drive a browser  

Some pages need JavaScript to render their content. For those cases, Ferret can use a browser-backed driver:  

{{< tabs >}} 
{{< tab title="Terminal" >}} 
{{< terminal command="true" >}}
ferret run -e '
let page = web::html::open("https://mockery.ferretlang.org", { driver: "cdp" })
return page.title
'
{{< /terminal >}}
{{< /tab >}}
{{< tab title="Try in browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
let page = web::html::open("https://mockery.ferretlang.org", { driver: "cdp" })
return page.title
{{< /editor >}}
{{< /tab >}}
{{< /tabs >}}

This example uses the `cdp` driver, which talks to a browser through the Chrome DevTools Protocol.  
Use the `cdp` driver when the page depends on JavaScript, client-side rendering, delayed content, user interaction, or browser APIs. 

For static HTML pages, the non-browser mode is usually simpler and faster.  

## Save a script

For anything longer than a small example, save the script to a `.fql` file:

{{< terminal command="true" >}}
echo 'return "Hello, Ferret"' > hello.fql
{{< /terminal >}}

Then run it:

{{< terminal command="true" >}}
ferret run hello.fql
{{< /terminal >}}

## Pass parameters

Scripts can also read parameters passed from the CLI. Parameters are available through the `@` prefix:

{{< terminal command="true" >}}
echo 'return "Hello, " + @name' > hello.fql
{{< /terminal >}}

Run it with a parameter:

{{< terminal command="true" >}}
ferret run hello.fql --param name=Steve
{{< /terminal >}}

In this example, `@name` reads the `name` parameter passed from the CLI.

## Next steps

{{< docs-related tiles="language,web-extraction,tools-cli-browser,tools-cli" >}}
