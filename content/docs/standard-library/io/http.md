---
title: "IO / HTTP"
sidebarTitle: "HTTP"
weight: 20
draft: false
description: "HTTP functions in the Ferret standard library."
aliases:
  - /docs/stdlib/io-net-http/
menuTitle: io::net::http
menu: [delete,do,get,post,put,]
---



{{< header href="delete" >}}

io::net::http::delete

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/io/net/http/delete.go#L15)

delete makes a HTTP delete request.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`params` | `Object`  |  | Request parameters.
`params.url` | `String`  |  | Target url
`params.body` | `Binary`  |  | Request data
`params.headers` | `Object`  |  | Http headers


**Returns** `Binary` Response in binary format
- - - -


{{< header href="do" >}}

io::net::http::do

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/io/net/http/request.go#L27)

REQUEST makes a HTTP request.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`params` | `Object`  |  | Request parameters.
`params.method` | `String`  |  | Http method
`params.url` | `String`  |  | Target url
`params.body` | `Binary`  |  | Request data
`params.headers` | `Object`  |  | Http headers


**Returns** `Binary` Response in binary format
- - - -


{{< header href="get" >}}

io::net::http::get

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/io/net/http/get.go#L16)

get makes a HTTP get request.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`urlOrParam` | `Object` `String`  |  | Target url or parameters.
`param.url` | `String`  |  | Target url or parameters.
`param.headers` | `Object`  |  | Http headers


**Returns** `Binary` Response in binary format
- - - -


{{< header href="post" >}}

io::net::http::post

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/io/net/http/post.go#L15)

post makes a post request.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`params` | `Object`  |  | Request parameters.
`params.url` | `String`  |  | Target url
`params.body` | `Binary`  |  | Request data
`params.headers` | `Object`  |  | Http headers


**Returns** `Binary` Response in binary format
- - - -


{{< header href="put" >}}

io::net::http::put

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/io/net/http/put.go#L15)

put makes a put HTTP request.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`params` | `Object`  |  | Request parameters.
`params.url` | `String`  |  | Target url
`params.body` | `Binary`  |  | Request data
`params.headers` | `Object`  |  | Http headers


**Returns** `Binary` Response in binary format
- - - -

## Next steps

{{< docs-related tiles="web-extraction,tools-worker,stdlib" >}}
