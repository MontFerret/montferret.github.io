---
title: "io/net/http"
weight: 1
draft: false
menuTitle: IO::NET::HTTP
menu: [DELETE,DO,GET,POST,PUT,]
---



{{< header href="delete" >}}

IO::NET::HTTP::DELETE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/io/net/http/delete.go#L15)

DELETE makes a HTTP DELETE request.

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

IO::NET::HTTP::DO

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

IO::NET::HTTP::GET

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/io/net/http/get.go#L16)

GET makes a HTTP GET request.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`urlOrParam` | `Object` `String`  |  | Target url or parameters.
`param.url` | `String`  |  | Target url or parameters.
`param.headers` | `Object`  |  | Http headers


**Returns** `Binary` Response in binary format
- - - -


{{< header href="post" >}}

IO::NET::HTTP::POST

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/io/net/http/post.go#L15)

POST makes a POST request.

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

IO::NET::HTTP::PUT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/io/net/http/put.go#L15)

PUT makes a PUT HTTP request.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`params` | `Object`  |  | Request parameters.
`params.url` | `String`  |  | Target url
`params.body` | `Binary`  |  | Request data
`params.headers` | `Object`  |  | Http headers


**Returns** `Binary` Response in binary format
- - - -
