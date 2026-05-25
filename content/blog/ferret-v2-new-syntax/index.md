---
title: "Inside Ferret v2: New Language Capabilities"
subtitle: "New Statements, Capabilities, and Syntax"
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2026-05-04"
---

Hello friends,

Today I want to talk about one of the most exciting parts of Ferret v2: new language capabilities that are now available directly in the syntax.

Ferret v2 is not meant to be a completely different language. The goal is to preserve the original spirit of Ferret - a practical language for extracting structured data from messy websites - while making the language more expressive, more consistent, and easier to extend.

The familiar Ferret style remains: small scripts, readable data flow, structured values, and extraction-focused operations. The new syntax is there to make those scripts clearer, not to make Ferret feel like a large general-purpose programming language.

That means adding new keywords, statements, and expressions where they help describe common extraction workflows more clearly: querying data, dispatching events, updating values, handling failures, waiting for conditions, and defining reusable functions.

So this post is not about replacing Ferret’s identity. It is about giving the language a stronger foundation for the next stage of the project.

## Why add new language capabilities?

Ferret v1 proved that a query-oriented language can be useful for web automation and data extraction. It allowed developers to describe browser interaction and data shaping in a way that felt more focused than stitching together low-level automation calls.

That idea is still at the center of Ferret v2.

But as Ferret grew, it became clear that some concepts deserved to be represented more directly in the language. Querying is one example. Waiting for something to appear is another. Dispatching an event, updating a value, or choosing an error strategy are all common parts of extraction workflows.

In Ferret v1, some of these ideas had to be expressed indirectly through functions or domain-specific helpers. That works, but it has limits. It can make scripts harder to read, harder to optimize, and harder to explain through good diagnostics.

Ferret v2 introduces syntax-level support for these concepts not because the language needs to look different, but because these operations are important enough to be first-class.

A function call can perform the same work, but it cannot always describe the same intent.

For example, a compiler can understand that `QUERY ".product-item" IN document USING css` is a query operation. It can attach better diagnostics to the selector, apply query-specific policies, validate dialect support, and eventually optimize or trace the operation as part of the execution model.

With a plain function call, most of that meaning is hidden behind an arbitrary host function boundary.

This also matters for error reporting. When the compiler understands that a piece of code is a query, a dispatch operation, a wait condition, or an assignment path, it can point to the part of the script that actually caused the problem. Better syntax gives Ferret better structure, and better structure leads to better diagnostics.

The goals are simple:

- keep common extraction workflows readable.
- make behavior explicit where it matters.
- make the language easier to compile and optimize.
- improve diagnostics by giving the compiler clearer structure.
- make Ferret extensible beyond a single hardcoded domain.

In other words, Ferret v2 is still Ferret. It just has a more capable language foundation.

## Query as a language capability

The most important example is querying.

In Ferret v1, querying was strongly associated with documents, elements, and selectors, but was implemented via host functions like `ELEMENTS` or `XPATH`.

In Ferret v2, querying becomes a general language capability. A value can support queries if it implements the queryable capability.

{{< editor lang="fql" height="132px" apiVersion="2" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev")
LET title = QUERY "h1" IN doc USING css
RETURN title
{{</ editor >}}

The important part is `USING css`. The query string does not define the meaning by itself. The dialect does.

That means the same language construct can work with different query dialects:

{{< code lang="fql" height="96px" >}}
LET products = QUERY ".product" IN doc USING css
LET rows = QUERY "SELECT * FROM products WHERE price > 100" IN db USING sql
{{</ code >}}

The language does not need separate statements for CSS selectors, XPath, SQL, JSONPath, or future query systems.

Instead, the target value and the selected dialect decide how the query is interpreted.

This is one of the main design directions in Ferret v2: keep the syntax stable, but allow capabilities to define behavior.

### Query modifiers

Queries often need different result shapes. Sometimes you want all matches. Sometimes you want one match. Sometimes you only care whether something exists or how many matches there are.

Ferret v2 introduces query modifiers for those cases:

{{< code lang="fql" height="128px" >}}
LET products = QUERY ".product" IN doc USING css
LET hasLogin = QUERY EXISTS "form.login" IN doc USING css
LET count = QUERY COUNT ".product" IN doc USING css
LET first = QUERY ONE ".product" IN doc USING css
{{</ code >}}

By default, `QUERY` returns the normal result shape for the selected dialect and target value. For CSS-style document queries, that usually means a collection of matching elements.

`QUERY EXISTS` returns a boolean value:

{{< editor lang="fql" height="84px" apiVersion="2" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev")
RETURN QUERY EXISTS "h1" IN doc USING css
{{</ editor >}}

`QUERY COUNT` returns the number of matches:

{{< editor lang="fql" height="84px" apiVersion="2" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev/scenarios/ecommerce/categories/laptops/")
RETURN QUERY COUNT ".product-card" IN doc USING css
{{</ editor >}}

`QUERY ONE` returns a single matching result:

{{< editor lang="fql" height="84px" apiVersion="2" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev/scenarios/ecommerce/categories/laptops/")
RETURN QUERY ONE ".product-card" IN doc USING css
{{</ editor >}}

This is useful for the common case where a script expects one element and does not want to query a collection and then index into it manually.

Instead of writing:

{{< code lang="fql" height="84px" >}}
LET products = QUERY ".product-card" IN doc USING css
LET product = products[0]
{{</ code >}}

A script can express the intent directly:

{{< code lang="fql" height="84px" >}}
LET product = QUERY ONE ".product" IN doc USING css
{{</ code >}}

The modifier makes cardinality visible. A query that returns all matches is different from a query that returns one match. A query that checks existence is different from a query that counts matches.

By making that intent visible in the syntax, Ferret can make scripts easier to read and potentially easier to optimize.

It also gives the runtime and diagnostics a clearer model. If a script asks for one result, Ferret can treat that as a distinct operation instead of a collection query followed by an index access.

### Query shorthand

For common cases, the long form can be too verbose. Ferret v2 also supports shorthand query expressions.

The regular shorthand uses `~`:

{{< editor lang="fql" height="112px" apiVersion="2" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev")
LET headings = doc[~ css`h1`]
RETURN headings
{{</ editor >}}

This is equivalent to a regular query:

{{< code lang="fql" height="84px" >}}
LET headings = QUERY "h1" IN doc USING css
{{</ code >}}

When a script expects a single result, Ferret also supports the `~?` shorthand:

{{< editor lang="fql" height="112px" apiVersion="2" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev")
LET title = doc[~? css`h1`]
RETURN title
{{</ editor >}}

This is equivalent to `QUERY ONE`:

{{< code lang="fql" height="84px" >}}
LET title = QUERY ONE "h3" IN doc USING css
{{</ code >}}

The distinction is small, but important.

`~` asks for the normal query result shape. For CSS-style document queries, that usually means a collection of matching elements.

`~?` asks for one matching result.

That makes the common “give me the first matching thing” case more readable without forcing scripts to query a collection and then index into it manually:

{{< code lang="fql" height="128px" >}}
LET title = doc[~? css`h1`]

// instead of
LET title = doc[~ css`h1`][0]
{{</ code >}}

More dynamic query expressions should still use the long form:

{{< code lang="fql" height="96px" >}}
LET selector = ".product[data-id='" + id + "']"
LET product = QUERY ONE selector IN doc USING css WITH { timeout: 5000 }
{{</ code >}}

{{< code lang="fql" height="96px" >}}
LET product = QUERY ONE "SELECT * FROM products WHERE id = ?" IN doc USING css WITH { params: [@id] }
{{</ code >}}

This gives Ferret both sides: concise syntax for everyday cases, and explicit syntax for cases that need more control.

The shorthand forms are intentionally limited. They should make common scripts pleasant to write, but they should not become a second full query language hidden inside brackets.

## Array operators for data shaping

Web extraction rarely returns one perfectly shaped value.

More often, a script gets a list of elements, rows, links, products, or records, and then needs to transform that list into clean structured output.

Ferret v2 brings array operators inspired by ArangoDB’s AQL into the language to make this kind of data shaping easier to express.

Instead of writing loops for every small transformation, scripts can map, filter, slice, and project arrays directly.

{{< editor lang="fql" height="356px" apiVersion="2" >}}
LET products = [
    { name: "Widget", price: 19.99 },
    { name: "Gadget", price: 149.99 },
    { name: "Thingamajig", price: 49.99 },
    { name: "Doodad", price: 9.99 },
    { name: "Doohickey", price: 199.99 }
]

LET names = products[*].name
LET expensive = products[? FILTER .price > 100]
LET firstThree = products[* LIMIT 2]

RETURN {
    names: names,
    expensive: expensive,
    firstThree: firstThree
}
{{</ editor >}}

The basic operators cover the most common cases:

- `[n]` accesses an array element by index.
- `[*]` expands an array and allows projecting fields from each item.
- `[**]`, `[***]`, and deeper forms flatten nested arrays.
- `[* FILTER ...]` filters an array while expanding it.
- `[* LIMIT ...]` limits the expanded result.
- `[* RETURN ...]` projects each expanded item into a new shape.
- `[? ...]` checks whether an array contains values matching a condition.

The `[*]` operator is especially useful after queries. For example, extracting all link targets from a page can stay compact:

{{< editor lang="fql" height="128px" apiVersion="2" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev/")
RETURN doc[~ css`a`][*].attributes.href
{{</ editor >}}

Inline expressions make filtering and projection more explicit. They use `.` to refer to the current array item:

{{< editor lang="fql" height="328px" apiVersion="2" >}}
LET products = [
    { name: "Widget", price: 19.99 },
    { name: "Gadget", price: 149.99 },
    { name: "Thingamajig", price: 49.99 },
    { name: "Doodad", price: 9.99 },
    { name: "Doohickey", price: 199.99 }
]

RETURN products[*
    FILTER .price > 100
    RETURN {
        title: .name,
        price: .price
    }
]
{{</ editor >}}

`FILTER`, `LIMIT`, and `RETURN` can be combined in that order:

{{< editor lang="fql" height="376px" apiVersion="2">}}
LET products = [
    { name: "Widget", price: 19.99 },
    { name: "Gadget", price: 149.99 },
    { name: "Thingamajig", price: 49.99 },
    { name: "Doodad", price: 9.99 },
    { name: "Doohickey", price: 199.99 },
    { name: "Whatchamacallit", price: 129.99 },
    { name: "Contraption", price: 89.99 },
    { name: "Gizmo", price: 179.99 }
]

RETURN products[*
    FILTER .price > 100
    LIMIT 3
    RETURN {
        title: .name,
        price: TO_FLOAT(.price)
    }
]
{{</ editor >}}

The question-mark operator is different. It does not return the filtered items. It answers whether matching items exist.

{{< editor lang="fql" height="274px" apiVersion="2">}}
LET products = [
    { name: "Widget", price: 19.99 },
    { name: "Gadget", price: 149.99 },
    { name: "Thingamajig", price: 49.99 },
    { name: "Doodad", price: 9.99 },
    { name: "Doohickey", price: 199.99 },
    { name: "Whatchamacallit", price: 129.99 },
    { name: "Contraption", price: 89.99 },
    { name: "Gizmo", price: 179.99 }
]

RETURN products[? ANY FILTER .price > 100]
{{</ editor >}}

This returns a boolean value.

That distinction matters: use `[* FILTER ...]` when you want the matching values, and `[? ...]` when you want to test whether matching values exist.

Array contraction is useful when querying nested collections:

{{< editor lang="fql" height="152px" apiVersion="2" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev/scenarios/ecommerce/")
LET sections = QUERY "section" IN doc USING css
LET linksBySection = sections[* RETURN  .[~ css`a`]]

RETURN linksBySection[**].attributes.href
{{</ editor >}}

This is especially useful after a query. `QUERY` gives the script a collection of values, and array operators help turn that collection into the shape the caller actually needs.

The goal is the same as with the rest of Ferret v2 syntax: make common extraction workflows readable without forcing every small data transformation into a manual loop.

## Dispatch as a language capability

Many values are not only queryable. They can also receive events, commands, or signals.

In browser automation, clicking an element is the obvious example. But dispatch is not limited to DOM events. A queue, actor, stream, workflow, or custom host object could also expose dispatch behavior.

Ferret v2 represents this with `DISPATCH`:

{{< code lang="fql" height="84px" >}}
DISPATCH "click" IN button
{{</ code >}}

With payload and options, the same statement stays explicit:

{{< code lang="fql" height="128px" >}}
DISPATCH "input" IN searchBox WITH {
    value: "ferret"
} OPTIONS {
    bubbles: true
}
{{</ code >}}

For simple payload-less signals, Ferret v2 also has a concise shorthand:

{{< code lang="fql" height="84px" >}}
button <- "click"  
{{</ code >}}

This reads as: send the click signal to button.

The shorthand is intentionally narrow. Once payloads, options, or more explicit behavior are needed, the long form is clearer.

Again, the important part is not only the browser use case. The important part is that dispatch becomes a language-level operation over values that support the dispatchable capability.

The value itself decides how to interpret the signal, what to do with the payload, and how to handle options.

Values with dispatch capabilities are provided by registered host modules.

## Match for structured control flow

Ferret v2 also improves control flow with `MATCH`.

The goal is not to make simple `IF` expressions obsolete. The goal is to provide a better structure when branching logic grows beyond one or two conditions.

Guard-style matching can express condition-based branching:

{{< editor lang="fql" height="160px" apiVersion="2" >}}
LET status = 501
RETURN MATCH (
    WHEN status == 404 => "not_found",
    WHEN status == 403 => "forbidden",
    WHEN status >= 500 => "server_error",
    _ => "ok"
)
{{</ editor >}}

Scrutinee-style matching can inspect a value directly:

{{< editor lang="fql" height="160px" apiVersion="2" >}}
LET status = 404
RETURN MATCH status (
    200 => "ok",
    404 => "not_found",
    _ => "unknown"
)
{{</ editor >}}

Ferret v2 also supports object pattern matching. This is useful when a script needs to branch based on the shape or selected fields of a value:

{{< editor lang="fql" height="256px" apiVersion="2" >}}
LET response = {
    status: 500,
    body: "Internal Server Error"
}

RETURN MATCH response (
    { status: 200, body: body } => body,
    { status: 404 } => NONE,
    { status: status } WHEN status >= 500 => "server_error",
    _ => "unknown"
)
{{</ editor >}}

This makes common extraction and normalization logic easier to express. Instead of pulling fields out first and then writing a chain of conditions, the match arm can describe the shape it expects and bind the values it needs.

The pattern system will continue to evolve, but object pattern matching is already supported and is part of Ferret’s current direction. `MATCH` should become the primary way to express structured branching when a script has multiple cases to handle.

This is especially useful in extraction workflows, where scripts often need to classify responses, handle missing values, normalize data, or branch based on page state.

## String templates

Ferret v2 also adds string templates for cases where scripts need to build readable strings from values.

Extraction scripts often need to create URLs, format labels, build messages, or normalize output fields. Plain string concatenation works, but it quickly becomes noisy.

{{< code lang="fql" height="96px" >}}
LET name = @user
LET message = "Hello " + name
RETURN message
{{</ code >}}

With string templates, the same expression is easier to read:

{{< code lang="fql" height="96px" >}}
LET message = `Hello ${@user}`
RETURN message
{{</ code >}}

Expressions inside `${...}` are evaluated and inserted into the final string:

{{< editor lang="fql" height="160px" apiVersion="2" >}}
LET product = {
    title: "Keyboard",
    price: 99.99
}

RETURN `${product.title}: $${product.price}`
{{</ editor >}}

String templates are especially useful when building dynamic query strings or URLs:

{{< editor lang="fql" height="128px" apiVersion="2" >}}
LET category = "phones"
LET page = 2

LET url = `https://mockery.montferret.dev/scenarios/ecommerce/categories/${category}/page/${page}`
RETURN WEB::ARTICLE::EXTRACT(DOCUMENT(url))
{{</ editor >}}

They also work well for shaping final output:

{{< editor lang="fql" height="224px" apiVersion="2" >}}
LET product = {
    title: "Keyboard",
    price: 99.99,
    available: true
}

RETURN {
    title: product.title,
    label: `${product.title} - $${product.price}`,
    status: `available: ${product.available}`
}
{{</ editor >}}

The goal is simple: keep common string-building cases readable without forcing scripts into long chains of concatenation.

## Mutable values and assignment

Ferret has traditionally favored query-style expression flow, but some tasks are simply easier with local mutable state.

Ferret v2 separates immutable and mutable bindings:

{{< editor lang="fql" height="352px" apiVersion="2" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev/scenarios/ecommerce/categories/laptops/", {
    driver: "cdp"
})
LET pages = doc[~ css`[data-testid="page-link"]`][* RETURN TO_INT(.innerText)]
VAR current = 1

LET aggregated = (FOR WHILE current < LENGTH(pages)
    current += 1
    LET products = doc[~ css`.product-card`][* RETURN {
        title: `${.[~? css`.product-brand`].textContent} ${.[~? css`.product-title`].textContent}`,
        price: TO_FLOAT(.[~? css`.product-price`].attributes["data-price"]) ON ERROR RETURN NONE
    }]

    DISPATCH "click" IN doc[~? css`:nth(${current}, [data-testid])`]
    WAITFOR EVENT "navigation" IN doc TIMEOUT 10s

    RETURN products
)

RETURN aggregated[**]
{{</ editor >}}

`LET` remains immutable. `VAR` is explicit. Reassignment is allowed only when the nearest binding is mutable.

{{< editor lang="fql" height="192px" apiVersion="2" >}}
LET baseUrl = "https://example.com"
LET attempts = 0

FOR WHILE attempts < 3
    attempts += 1

RETURN attempts
{{</ editor >}}

The second example is invalid because `attempts` was declared with `LET`.

This keeps mutation available without making every binding mutable by default.

One important distinction is that `LET` prevents rebinding the variable itself. It does not necessarily make the underlying value deeply immutable. If a value supports mutation, its fields may still be updated.

For objects and other mutable values, Ferret uses familiar assignment syntax:

{{< editor lang="fql" height="312px" apiVersion="2" >}}
LET user = {
    name: "Bob",
    profile: {
        active: false
    },
    stats: {
        visits: 0
    }
}

user.name = "Alice"
user.profile.active = true
user.stats.visits += 1

RETURN user
{{</ editor >}}

Safe access also applies naturally to mutation paths:

{{< editor lang="fql" height="128px" apiVersion="2" >}}
LET user = NONE
user?.profile?.active = true
RETURN user
{{</ editor >}}

The goal is not to make Ferret more imperative than necessary. The goal is to support the cases where mutation describes the workflow more naturally, while preserving immutability as the default style.

### Deleting properties

Assignment lets scripts create or update values, but extraction workflows also often need to remove data.

A script may need to clean up intermediate fields, remove deprecated metadata, drop optional values, or normalize an object before returning it.

Ferret v2 adds `DELETE` for removing properties from mutable values:

{{< editor lang="fql" height="240px" apiVersion="2" >}}
LET user = {
    name: "Alice",
    profile: {
        active: true,
        deprecated: true
    }
}

DELETE user.profile.deprecated

RETURN user
{{</ editor >}}

This removes the final property in the path. It does not delete the whole object. In this example, `profile` remains, but `profile.deprecated` is removed.

Bracket access is supported as well:

{{< editor lang="fql" height="240px" apiVersion="2" >}}
LET user = {
    name: "Alice",
    metadata: {
        source: "import",
        temp: true
    }
}

DELETE user.metadata["temp"]

RETURN user
{{</ editor >}}

Safe access can be used when intermediate values may be missing:

{{< editor lang="fql" height="152px" apiVersion="2" >}}
LET user = NONE

DELETE user?.profile?.deprecated

RETURN user
{{</ editor >}}

This makes deletion safe and explicit. If the path cannot be reached because a safe segment evaluates to `NONE`, the operation does nothing.

Like assignment, `DELETE` works on values that support mutation. The statement describes the operation, while the value decides how that operation is applied.

The goal is to make cleanup and normalization code readable without introducing helper functions for basic property removal.

## Waiting as part of extraction

Waiting is another operation that deserves first-class treatment in web extraction.

Pages are dynamic. Data may appear after a network request, a DOM update, an animation, or a client-side route change. In Ferret v1, this kind of behavior often had to be expressed through helper functions or custom retry logic.

Ferret v2 makes waiting explicit:

{{< editor lang="fql" height="224px" apiVersion="2" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev/scenarios/delayed-rendering/", {
    driver: "cdp"
})

RETURN WAITFOR VALUE doc[~ css`[data-testid="delayed-long"]`]
    WHEN LENGTH(.) > 0 AND FIRST(.).attributes["data-state"] == "ready"
    TIMEOUT 10s
    EVERY 250ms
    ON TIMEOUT RETURN NONE
{{</ editor >}}

This describes the operation directly: evaluate the value repeatedly, use a polling interval, stop after a timeout, and choose a fallback if the condition is not met.

The result is easier to read than hand-written retry logic, and easier for the runtime to trace, optimize, and explain.

### Waiting for network activity

Some extraction workflows depend not only on the DOM, but also on network activity.

Modern pages often load data lazily through background requests. A product list, search result, price, availability status, or recommendation block may appear only after one or more API calls finish.

The CDP driver exposes network lifecycle events that can be observed from Ferret scripts:

- `network.request_started`
- `network.response_received`
- `network.request_finished`
- `network.request_failed`
- `network.idle`

This makes it possible to wait for network behavior directly instead of guessing with fixed delays.

For example, a script can wait until the page becomes network-idle before querying the DOM:

{{< editor lang="fql" height="224px" apiVersion="2" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev/scenarios/network/network-idle/", {
    driver: "cdp"
})

WAITFOR EVENT "network.idle" IN doc TIMEOUT 10s

RETURN doc[~ css`#network-log li`][*].textContent
{{</ editor >}}

A script can also wait for a specific request to finish before reading the updated page state:

{{< editor lang="fql" height="224px" apiVersion="2" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev/scenarios/network/delayed-requests/", {
    driver: "cdp"
})

RETURN WAITFOR EVENT "network.request_finished" IN doc
    WHEN CONTAINS(.url, "slow-1.json")
    TIMEOUT 10s
{{</ editor >}}

Network events are also useful for debugging or collecting metadata from a page session:

{{< editor lang="fql" height="384px" apiVersion="2" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev/scenarios/network/polling/", {
    driver: "cdp"
})

LET response = WAITFOR EVENT "network.response_received" IN doc
    WHEN CONTAINS(.url, ".json")
    TIMEOUT 10s
    ON TIMEOUT RETURN NONE

RETURN MATCH response (
    NONE => {
        ok: false,
        reason: "products_api_not_seen"
    },
    _ => {
        ok: true,
        url: response.url,
        status: response.status
    }
)
{{</ editor >}}

The important part is that network activity becomes observable through the same language-level waiting model. Ferret does not need a separate `WAITFOR NETWORK` construct. The CDP driver can expose network activity as events, and `WAITFOR EVENT` can observe them.

This keeps the language general while still supporting browser-specific workflows.

## Error and timeout policies close to the operation

Web data extraction often fails for normal reasons: a page is slow, an element is missing, a network request times out, or a site returns an unexpected response.

In Ferret v2, failure policy can live close to the operation that may fail.

{{< editor lang="fql" height="128px" apiVersion="2" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev")
RETURN QUERY VALUE "#price" IN doc USING css 
    ON ERROR RETURN NONE
{{</ editor >}}

Timeout behavior can be expressed in a similar way where supported:

{{< editor lang="fql" height="146px" apiVersion="2" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev")
RETURN WAITFOR VALUE doc[~ css`.loaded`]
    TIMEOUT 5s
    EVERY 250ms
    ON TIMEOUT RETURN NONE
{{</ editor >}}

This keeps the happy path readable while making fallback behavior explicit.

It also gives Ferret a clearer execution model. A timeout policy or fallback value is not hidden inside arbitrary user code. It is part of the operation itself.

{{< editor lang="fql" height="120px" apiVersion="2" >}}
LET parsed = TO_FLOAT("not a number") ON ERROR RETURN 42
RETURN parsed
{{</ editor >}}

## User-defined functions

Ferret v2 is also moving toward user-defined functions.

The goal is not to turn Ferret into a general-purpose application language. The goal is to let scripts define reusable extraction and normalization logic directly where it belongs.

Functions are especially useful for normalization logic: parsing prices, cleaning text, mapping statuses, extracting IDs, and turning inconsistent page data into stable output shapes.

For larger functions, the block form gives enough structure without relying on indentation-sensitive syntax or `END` markers.

{{< editor lang="fql" height="192px" apiVersion="2" >}}
FUNC normalizePrice(input) (
    LET cleaned = TRIM(input)
    LET numeric = SUBSTITUTE(cleaned, "$", "")
    RETURN TO_FLOAT(numeric)
)

RETURN normalizePrice("$19.99")
{{</ editor >}}

For smaller functions, the body can stay compact:

{{< editor lang="fql" height="96px" apiVersion="2" >}}
FUNC add(a, b) => a + b
RETURN add(2, 3)
{{</ editor >}}

This should make Ferret scripts easier to organize while keeping the language focused.

### Control flow using pattern matching and user-defined functions

With user-defined functions and pattern matching, Ferret can express more complex logic without relying on host functions or external code.

This is an important step toward making Ferret a more self-contained language for data extraction and processing.

{{< editor lang="fql" height="196px" apiVersion="2" >}}
FUNC fib(n) (
    RETURN MATCH n (
        0 => 0,
        1 => 1,
        _ => fib(n - 1) + fib(n - 2)
    )
)

RETURN fib(10)
{{</ editor >}}

## Where values with capabilities come from

One important question is where these values with capabilities come from.

A Ferret script does not manually attach capabilities to a value. Capabilities come from the runtime, modules, and host applications that embed or extend Ferret.

For example, an HTML module can expose a document value that supports CSS queries:

{{< code lang="fql" height="96px" >}}
LET document = PARSE(page)
LET title = QUERY VALUE "h1" IN document USING css
{{</ code >}}

A DOM element can be both queryable and dispatchable:

{{< code lang="fql" height="96px" >}}
LET button = QUERY ONE "button.submit" IN document USING css
DISPATCH "click" IN button
{{</ code >}}

A database module could expose a connection or table value that supports SQL queries:

{{< code lang="fql" height="96px" >}}
LET rows = QUERY "SELECT * FROM products" IN db USING sql
{{</ code >}}

A host application embedding Ferret can also provide its own values. That could be a queue, cache, workflow object, browser session, API client, or any other domain-specific object.

From the script’s point of view, these values participate in the same language constructs. The script does not need to know whether a value came from the standard library, a contrib module, or the embedding application. It only matters which capabilities the value exposes.

As the compiler and runtime mature, this model also gives Ferret a clearer way to report capability errors. The problem is not just that a function failed, but that a value does not support the operation being requested.

For example, `QUERY ... IN value USING css` requires the target value to support querying with the selected dialect. If it does not, Ferret can report that directly.

This is especially important for the Ferret ecosystem. Core Ferret can stay small, while modules can add support for new document types, protocols, storage systems, APIs, or browser/runtime integrations without requiring new syntax for each one.

This is what allows Ferret to stay small at the language level while still being extensible at the runtime and module level.

## Module aliases with USE

Ferret v2 also improves how scripts work with modules.

Modules can expose host functions, but fully-qualified names can become noisy when a script uses the same module repeatedly.

`USE` lets a script create a local alias for a module namespace:

{{< editor lang="fql" height="160px" apiVersion="2" >}}
USE IO::NET::HTTP::GET AS GET

LET out = GET("https://mockery.montferret.dev/api/products/index.json")

RETURN JSON_PARSE(out)
{{</ editor >}}

This keeps module-based scripts readable without pulling individual functions directly into local scope.

The important part is that `USE` does not hide where functionality comes from. The function call is still namespaced, but the namespace can be made shorter and more convenient for the script.

This keeps Ferret explicit while making module-heavy scripts easier to read.

## Putting it together

The individual features are useful on their own, but they are designed to work together.

{{< editor lang="fql" height="520px" apiVersion="2" >}}
FUNC normalizePrice(input) (
    LET cleaned = TRIM(input)
    LET numeric = SUBSTITUTE(cleaned, "$", "")
    RETURN TO_FLOAT(numeric)
)

FUNC processItem(product) (
    RETURN {
        title: product[~ css`product-title`][0]?.textContent,
        price: normalizePrice(product[~ css`.product-price`][0]?.textContent)
    }
)

LET doc = DOCUMENT("https://mockery.montferret.dev/scenarios/dynamic-products/load-more/", { driver: "cdp" })

DISPATCH "scroll" IN doc WITH { to: "bottom", behavior: "smooth" }
LET btn = FIRST(doc[~ css`#load-more-products`])

LET _ = (
    FOR WHILE btn.attributes["data-complete"] != "true"
        DISPATCH "click" IN btn

        RETURN WAITFOR EVENT "network.idle" IN doc TIMEOUT 40s 
)

LET products = doc[~ css`.product-card`]

RETURN products[* LIMIT 5 RETURN processItem(.)]
{{</ editor >}}

This is the kind of script Ferret v2 is designed to make easier to write: query the page, wait for dynamic content, handle missing data explicitly, normalize the result, and return a clean structured value.

## The bigger idea: capability-oriented syntax

The common thread behind these additions is capability-oriented design.

Ferret v2 does not need every domain to become a new language feature. Instead, the language provides a small set of operations that can work across different kinds of values:

- `QUERY` for values that can be queried.
- `DISPATCH` for values that can receive events or signals.
- `WAITFOR` for values or expressions that can be observed over time.
- `MATCH` for structured branching.
- assignment for local and object mutation.
- operation-level policies for errors and timeouts.

The goal is not to turn every useful library operation into syntax. Only operations that shape the structure of extraction workflows should become language constructs.

The syntax stays focused, while capabilities allow modules and host applications to define what those operations mean for their own values.

That is what makes Ferret more than a browser scripting language. Browser automation remains an important use case, but the language is being shaped as a programmable data extraction engine.

Ferret should help developers turn messy websites and data sources into clean structured data. These language-level capabilities are designed around that goal.

## Closing thoughts

Ferret v2 is still evolving, but the direction is intentional.

The language should stay concise for common extraction tasks, explicit when behavior matters, and structured enough for better diagnostics, optimization, and tooling.

The new language capabilities are not about changing Ferret for the sake of change. They are about giving important extraction concepts a clear place in the language.

Ferret v2 is still Ferret: practical, focused, and built for developers who need reliable structured data from messy sources.

It just has a stronger foundation now.