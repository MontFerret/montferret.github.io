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
LET response = IO::NET::HTTP::GET("https://jsonplaceholder.typicode.com/posts/1")
LET data = JSON_PARSE(TO_STRING(response))
RETURN data
'
{{< /terminal >}}
{{< /tab >}}

{{< tab title="Try in browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET response = IO::NET::HTTP::GET("https://jsonplaceholder.typicode.com/posts/1")
LET data = JSON_PARSE(TO_STRING(response))
RETURN data
{{< /editor >}}
{{< /tab >}}
{{< /tabs >}}

`IO::NET::HTTP::GET` returns raw bytes. Use `TO_STRING` to convert to a string, then `JSON_PARSE` to decode JSON.

## Make a POST request

Use `IO::NET::HTTP::POST` with a body and headers:

{{< code lang="fql" >}}
LET response = IO::NET::HTTP::POST({
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

RETURN JSON_PARSE(TO_STRING(response))
{{</ code >}}

## Iterate over API results

Fetch a list and process it with `FOR`:

{{< tabs >}}
{{< tab title="Terminal" >}}
{{< terminal command="true" >}}
ferret run -e '
LET response = IO::NET::HTTP::GET("https://jsonplaceholder.typicode.com/posts")
LET posts = JSON_PARSE(TO_STRING(response))

FOR post IN posts
    LIMIT 5
    RETURN {
        id: post.id,
        title: post.title
    }
'
{{< /terminal >}}
{{< /tab >}}

{{< tab title="Try in browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET response = IO::NET::HTTP::GET("https://jsonplaceholder.typicode.com/posts")
LET posts = JSON_PARSE(TO_STRING(response))

FOR post IN posts
    LIMIT 5
    RETURN {
        id: post.id,
        title: post.title
    }
{{< /editor >}}
{{< /tab >}}
{{< /tabs >}}

## Paginate an API

Many APIs use offset or page-based pagination:

{{< code lang="fql" >}}
LET baseURL = "https://jsonplaceholder.typicode.com/posts?_start="
LET pageSize = 10

LET result = (
    FOR pageNum IN 0..2
        LET offset = pageNum * pageSize
        LET url = baseURL + TO_STRING(offset) + "&_limit=" + TO_STRING(pageSize)
        LET response = IO::NET::HTTP::GET(url)
        LET posts = JSON_PARSE(TO_STRING(response))

        FOR post IN posts
            RETURN {
                id: post.id,
                title: post.title
            }
)

RETURN result
{{</ code >}}

## Add headers and authentication

Pass custom headers for APIs that require authentication:

{{< code lang="fql" >}}
LET response = IO::NET::HTTP::GET({
    url: "https://api.example.com/data",
    headers: {
        "Authorization": "Bearer " + @token,
        "Accept": "application/json"
    }
})

RETURN JSON_PARSE(TO_STRING(response))
{{</ code >}}

Use a bind parameter (`@token`) so the secret is not hardcoded in the script:

{{< terminal command="true" >}}
ferret run script.fql --param token=your-api-key
{{< /terminal >}}

## Combine API and HTML data

A powerful pattern: fetch structured data from an API and enrich it with data from HTML pages:

{{< code lang="fql" >}}
LET response = IO::NET::HTTP::GET("https://jsonplaceholder.typicode.com/posts")
LET posts = JSON_PARSE(TO_STRING(response))

FOR post IN posts
    LIMIT 3

    LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
        ON ERROR RETURN NONE

    RETURN {
        id: post.id,
        title: post.title,
        pageTitle: page?.title
    }
{{</ code >}}

## Handle API errors

Use `ON ERROR RETURN` and `ON ERROR RETRY` for network failures:

{{< code lang="fql" >}}
LET response = IO::NET::HTTP::GET("https://api.example.com/data")
    ON ERROR RETRY 3 DELAY 1s BACKOFF EXPONENTIAL
    OR RETURN NONE

RETURN response != NONE
    ? JSON_PARSE(TO_STRING(response))
    : { error: "API unavailable" }
{{</ code >}}

## Next steps

{{< docs-related tiles="guide-static-pages,guide-error-handling,stdlib,language-parameters" >}}
