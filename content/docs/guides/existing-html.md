---
title: "Process existing HTML"
sidebarTitle: "Existing HTML"
weight: 15
draft: false
description: "Pass HTML from Go into Ferret, extract structured data, and return modified markup."
---

# Process existing HTML

Use `web::html::parse` when your Go application already has the HTML. The content might come from `os.ReadFile`, a database, an HTTP response fetched by the host, or another service. Ferret can parse, query, and modify that content without fetching it again.

This workflow does not require CDP. Configure the HTML module with the in-process `memory` driver and pass the HTML through a session parameter.

## Parse supplied HTML

Keep the FQL source stable and reference the supplied content as `@html`:

{{< code lang="fql" >}}
let page = web::html::parse(@html)
let products = query ".product" in page using css

return for product in products
    return {
        name: (query one ".name" in product using css).textContent,
        url: (query one "a" in product using css).attributes.href
    }
{{</ code >}}

`web::html::parse` accepts an FQL string or binary value. In Go, `WithSessionParam` converts a `string` to an FQL string and `[]byte` to an FQL binary value. Do not interpolate the HTML into the FQL source.

## Run the workflow from Go

Create a module and install Ferret, the HTML module, and the small HTML parser used at the end of the example:

{{< terminal command="true" >}}
mkdir ferret-existing-html && cd ferret-existing-html
go mod init example.com/ferret-existing-html
go get github.com/MontFerret/ferret/v2 \
  github.com/MontFerret/contrib/modules/web/html \
  golang.org/x/net/html
{{</ terminal >}}

The following program runs two plans against the same document. The first passes a Go string and decodes structured output. The second passes `[]byte`, modifies the document, and decodes the returned HTML string.

{{< code lang="go" >}}
package main

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "log"
    "strings"

    htmlmodule "github.com/MontFerret/contrib/modules/web/html"
    "github.com/MontFerret/contrib/modules/web/html/drivers/memory"
    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/source"
    gohtml "golang.org/x/net/html"
)

const htmlSource = `<!doctype html>
<html>
<head><title>Catalog</title></head>
<body>
  <article class="product">
    <a href="/keyboard"><span class="name">Keyboard</span></a>
  </article>
  <article class="product">
    <a href="/mouse"><span class="name">Mouse</span></a>
  </article>
</body>
</html>`

type Product struct {
    Name string `json:"name"`
    URL  string `json:"url"`
}

func main() {
    if err := run(); err != nil {
        log.Fatal(err)
    }
}

func run() (err error) {
    htmlmod, err := htmlmodule.New(
        htmlmodule.WithDefaultDriver(memory.New()),
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
    defer func() {
        err = errors.Join(err, engine.Close())
    }()

    ctx := context.Background()

    productsPlan, err := engine.Compile(ctx, source.New("products", `
        let page = web::html::parse(@html)
        let products = query ".product" in page using css

        return for product in products
            return {
                name: (query one ".name" in product using css).textContent,
                url: (query one "a" in product using css).attributes.href
            }
    `))
    if err != nil {
        return err
    }
    defer func() {
        err = errors.Join(err, productsPlan.Close())
    }()

    productsOutput, err := execute(ctx, productsPlan, htmlSource)
    if err != nil {
        return err
    }
    if productsOutput.ContentType != "application/json" {
        return fmt.Errorf("unexpected products content type %q", productsOutput.ContentType)
    }

    var products []Product
    if err := json.Unmarshal(productsOutput.Content, &products); err != nil {
        return fmt.Errorf("decode products: %w", err)
    }
    fmt.Printf("products (%s): %+v\n", productsOutput.ContentType, products)

    htmlPlan, err := engine.Compile(ctx, source.New("modified-html", `
        let page = web::html::parse(@html)
        web::html::attr_set(page.body, "data-processed", "true")
        return page.innerHTML
    `))
    if err != nil {
        return err
    }
    defer func() {
        err = errors.Join(err, htmlPlan.Close())
    }()

    htmlOutput, err := execute(ctx, htmlPlan, []byte(htmlSource))
    if err != nil {
        return err
    }
    if htmlOutput.ContentType != "application/json" {
        return fmt.Errorf("unexpected HTML content type %q", htmlOutput.ContentType)
    }

    var renderedHTML string
    if err := json.Unmarshal(htmlOutput.Content, &renderedHTML); err != nil {
        return fmt.Errorf("decode HTML: %w", err)
    }
    if !strings.Contains(renderedHTML, `data-processed="true"`) {
        return errors.New("modified HTML is missing the data-processed attribute")
    }

    // Pass the decoded string to another HTML-processing library.
    if _, err := gohtml.Parse(strings.NewReader(renderedHTML)); err != nil {
        return fmt.Errorf("parse returned HTML: %w", err)
    }

    fmt.Printf("html (%s): %s\n", htmlOutput.ContentType, renderedHTML)
    return nil
}

func execute(
    ctx context.Context,
    plan *ferret.Plan,
    htmlInput any,
) (*ferret.Output, error) {
    session, err := plan.NewSession(ctx,
        ferret.WithSessionParam("html", htmlInput),
    )
    if err != nil {
        return nil, err
    }

    output, runErr := session.Run(ctx)
    closeErr := session.Close()
    if err := errors.Join(runErr, closeErr); err != nil {
        return nil, err
    }

    return output, nil
}
{{</ code >}}

Resolve the complete module graph and run the program:

{{< terminal command="true" >}}
go mod tidy
go run .
{{</ terminal >}}

The engine owns runtime-wide configuration, each plan owns compiled code and its VM pool, and each session owns one execution. The example closes all three levels. A real application can reuse the engine and plans while creating a new session for each input document.

## Decode the output format

`Session.Run` returns an `Output` with two fields:

- `Output.Content` contains encoded bytes, not a raw Ferret runtime value.
- `Output.ContentType` identifies the codec that produced those bytes.

JSON is the default output format. Objects and arrays therefore decode into Go structs, slices, maps, or other JSON targets with `json.Unmarshal`.

A returned FQL string is also JSON-encoded. For example, returned HTML begins and ends with JSON quotes and may contain escaped characters in `Output.Content`. Decode it into a Go `string` before writing it to a file, returning it as HTML, or passing it to another library.

The example checks `Output.ContentType` before using `json.Unmarshal`. Keep that check when the application allows [configurable output codecs]({{< ref "/docs/embedding/go/value-encoders" >}}); do not assume every result is JSON.

## Work with the returned markup

`page.innerHTML` serializes the current in-memory document, including mutations such as the `data-processed` attribute in the example. HTML parsing and serialization can normalize the markup, including tag structure, whitespace, quoting, and character escaping. Treat the result as equivalent parsed HTML rather than a byte-for-byte copy of the input.

Once decoded into a Go string, the markup is ordinary host data. The example passes it to `golang.org/x/net/html`, but the same boundary works with another parser, sanitizer, renderer, or storage layer.

## When a browser-backed document is needed

The memory driver is the default for supplied HTML because parsing, querying, mutation, and serialization all run in process. Configure the CDP driver only when the supplied document must support browser-backed operations. CDP requires a running browser endpoint and is not needed for the workflow in this guide.

## Next steps

{{< docs-related tiles="guide-static-pages,embedding-go-parameters,embedding-go-safe-fql-construction,embedding-go-value-encoders" >}}
