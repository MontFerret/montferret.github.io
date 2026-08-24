---
title: "Construct FQL Safely"
sidebarTitle: "Safe FQL Construction"
weight: 35
draft: false
description: "Keep FQL source separate from runtime data when embedding Ferret in Go."
---

# Construct FQL safely

An embedding application provides two different inputs to Ferret: FQL source to compile and runtime values for that source to consume. Keep that boundary explicit.

{{% notification type="warning" %}}
**FQL source is code.** Do not insert runtime data with `fmt.Sprintf` or string concatenation. URLs, CSS selectors, HTML, credentials, identifiers, request bodies, and configuration values should normally be passed as engine or session parameters.
{{% /notification %}}

## Parameterize runtime values

Interpolating a value into FQL makes the value part of the source text. Quotes or other FQL syntax in that value can change the program or make it fail to compile.

Do not construct a query like this:

{{< code lang="go" >}}
query := fmt.Sprintf(`return web::html::open("%s").title`, url)
{{</ code >}}

Compile stable FQL that references `@url`, then provide the URL when creating a session:

{{< code lang="go" >}}
htmlmod, err := html.New(
    html.WithDefaultDriver(memory.New()),
)
if err != nil {
    return err
}

engine, err := ferret.New(
    ferret.WithModules(htmlmod),
)
if err != nil {
    return err
}
defer engine.Close()

plan, err := engine.Compile(ctx, source.New("page-title", `
    return web::html::open(@url).title
`))
if err != nil {
    return err
}
defer plan.Close()

session, err := plan.NewSession(ctx,
    ferret.WithSessionParam("url", url),
)
if err != nil {
    return err
}
defer session.Close()

output, err := session.Run(ctx)
if err != nil {
    return err
}

fmt.Println(string(output.Content))
{{</ code >}}

The engine owns the configured HTML module. The plan contains compiled FQL code, and each session supplies the data for one execution.

Use engine parameters for stable application-wide values and session parameters for per-execution values. Both keep values out of the source text. See [Parameters]({{< ref "/docs/embedding/go/parameters" >}}) for the complete API.

## Pass selectors as parameters

A CSS selector is runtime data when it comes from a request, configuration file, database, or another host input. Pass it separately even when the FQL program uses it as a query selector:

{{< code lang="go" >}}
plan, err := engine.Compile(ctx, source.New("selected-text", `
    let page = web::html::parse(@html)
    return web::html::inner_text(page, @selector)
`))
if err != nil {
    return err
}
defer plan.Close()

session, err := plan.NewSession(ctx,
    ferret.WithSessionParam("html", htmlSource),
    ferret.WithSessionParam("selector", selector),
)
if err != nil {
    return err
}
defer session.Close()

output, err := session.Run(ctx)
if err != nil {
    return err
}

fmt.Println(string(output.Content))
{{</ code >}}

The selector remains a runtime string. It cannot add FQL statements or change the structure of the compiled plan.

## Pass existing HTML as a parameter

When the host already has HTML, pass the content to `web::html::parse` as a parameter instead of embedding it in a string literal:

{{< code lang="go" >}}
plan, err := engine.Compile(ctx, source.New("parsed-title", `
    let page = web::html::parse(@html)
    return page.title
`))
if err != nil {
    return err
}
defer plan.Close()

session, err := plan.NewSession(ctx,
    ferret.WithSessionParam("html", htmlSource),
)
if err != nil {
    return err
}
defer session.Close()

output, err := session.Run(ctx)
if err != nil {
    return err
}

fmt.Println(string(output.Content))
{{</ code >}}

This preserves the HTML exactly as data, including quotes, backticks, and text that resembles FQL syntax.

## Compose source only from authored FQL

An intentionally authored FQL fragment is code, so an application may compose it into a larger FQL source. Runtime data consumed by that fragment must still be passed separately.

{{< code lang="go" >}}
const predicate = `item.price <= @max_price`

query := `
    return for item in @items {
        filter ` + predicate + `
        return item
    }
`

plan, err := engine.Compile(ctx, source.New("affordable-products", query))
if err != nil {
    return err
}
defer plan.Close()

session, err := plan.NewSession(ctx,
    ferret.WithSessionParam("items", items),
    ferret.WithSessionParam("max_price", maxPrice),
)
if err != nil {
    return err
}
defer session.Close()

output, err := session.Run(ctx)
if err != nil {
    return err
}

fmt.Println(string(output.Content))
{{</ code >}}

Here `predicate` is deliberately authored FQL code. `items` and `maxPrice` are runtime data and remain parameters.

## Next steps

{{< docs-related tiles="embedding-go-parameters,embedding-go-executing,embedding-go-application-guide" >}}
