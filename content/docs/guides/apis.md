---
title: "Work with APIs"
sidebarTitle: "APIs"
weight: 60
draft: false
description: "Fetch data from REST APIs and combine it with HTML extraction."
---

# Work with APIs

Ferret is not limited to HTML. It can fetch data from REST APIs, parse JSON responses, and combine API data with HTML extraction in a single script.

## Make a GET request

Use `IO::NET::HTTP::GET` to fetch data from an API:

{{< tabs >}}
{{< tab title="Terminal" >}}
{{< terminal command="true" >}}
ferret run -e '
let response = IO::NET::HTTP::GET("https://jsonplaceholder.typicode.com/posts/1")
let data = JSON_PARSE(TO_STRING(response))
return data
'
{{< /terminal >}}
{{< /tab >}}

{{< tab title="Try in browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
let response = IO::NET::HTTP::GET("https://jsonplaceholder.typicode.com/posts/1")
let data = JSON_PARSE(TO_STRING(response))
return data
{{< /editor >}}
{{< /tab >}}
{{< /tabs >}}

`IO::NET::HTTP::GET` returns raw bytes. Use `TO_STRING` to convert to a string, then `JSON_PARSE` to decode JSON.

## Make a POST request

Use `IO::NET::HTTP::POST` with a body and headers:

{{< code lang="fql" >}}
let response = IO::NET::HTTP::POST({
    url: "https://jsonplaceholder.typicode.com/posts",
    body: TO_BINARY(JSON_STRINGIFY({
        title: "Ferret",
        body: "Data extraction",
        userId: 1
    })),
    headers: {
        "Content-Type": "application/json"
    }
})

return JSON_PARSE(TO_STRING(response))
{{</ code >}}

## Iterate over API results

Fetch a list and process it with `for`:

{{< tabs >}}
{{< tab title="Terminal" >}}
{{< terminal command="true" >}}
ferret run -e '
let response = IO::NET::HTTP::GET("https://jsonplaceholder.typicode.com/posts")
let posts = JSON_PARSE(TO_STRING(response))

return for post in posts
    limit 5
    return {
        id: post.id,
        title: post.title
    }
'
{{< /terminal >}}
{{< /tab >}}

{{< tab title="Try in browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
let response = IO::NET::HTTP::GET("https://jsonplaceholder.typicode.com/posts")
let posts = JSON_PARSE(TO_STRING(response))

return for post in posts
    limit 5
    return {
        id: post.id,
        title: post.title
    }
{{< /editor >}}
{{< /tab >}}
{{< /tabs >}}

## Paginate an API

Many APIs use offset or page-based pagination:

{{< code lang="fql" >}}
let baseURL = "https://jsonplaceholder.typicode.com/posts?_start="
let pageSize = 10

let result = (
    for pageNum in 0..2
        let offset = pageNum * pageSize
        let url = baseURL + TO_STRING(offset) + "&_limit=" + TO_STRING(pageSize)
        let response = IO::NET::HTTP::GET(url)
        let posts = JSON_PARSE(TO_STRING(response))

        for post in posts
            return {
                id: post.id,
                title: post.title
            }
)

return result
{{</ code >}}

## Add headers and authentication

Pass custom headers for APIs that require authentication:

{{< code lang="fql" >}}
let response = IO::NET::HTTP::GET({
    url: "https://api.example.com/data",
    headers: {
        "Authorization": "Bearer " + @token,
        "Accept": "application/json"
    }
})

return JSON_PARSE(TO_STRING(response))
{{</ code >}}

Use a bind parameter (`@token`) so the secret is not hardcoded in the script:

{{< terminal command="true" >}}
ferret run script.fql --param token=your-api-key
{{< /terminal >}}

## Combine API and HTML data

A powerful pattern: fetch structured data from an API and enrich it with data from HTML pages:

{{< code lang="fql" >}}
let response = IO::NET::HTTP::GET("https://jsonplaceholder.typicode.com/posts")
let posts = JSON_PARSE(TO_STRING(response))

return for post in posts
    limit 3

    let page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
        on error return none

    return {
        id: post.id,
        title: post.title,
        pageTitle: page?.title
    }
{{</ code >}}

## Handle API errors

Use `on error return` and `on error retry` for network failures:

{{< code lang="fql" >}}
let response = IO::NET::HTTP::GET("https://api.example.com/data")
    on error retry 3 delay 1s backoff EXPONENTIAL
    or return none

return response != none
    ? JSON_PARSE(TO_STRING(response))
    : { error: "API unavailable" }
{{</ code >}}

## Next steps

{{< docs-related tiles="guide-static-pages,guide-error-handling,stdlib,language-parameters" >}}
