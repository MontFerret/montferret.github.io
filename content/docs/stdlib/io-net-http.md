---
title: "io/net/http"
weight: 1
draft: false
---


## POST
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/io/net/http/post.go#L14)

POST makes a POST request.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`params (Object) - request parameters. * url (String) - Target url * body` | `Binary` | Post data * headers (object) optional - http headers


**Returns** `None`
- - - -

## DELETE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/io/net/http/delete.go#L14)

DELETE makes a HTTP DELETE request.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`params (Object) - request parameters. * url (String) - Target url * body` | `Binary` | Post data * headers (object) optional - http headers


**Returns** `None`
- - - -

## PUT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/io/net/http/put.go#L14)

PUT makes a PUT HTTP request.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`params (Object) - request parameters. * url (String) - Target url. * body` | `Binary` | Post data. * headers (object) optional - http headers.


**Returns** `None`
- - - -

## GET
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/io/net/http/get.go#L15)

GET makes a HTTP GET request.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`url or  (String) - target url or parameters. * url` | `String` | Target url * headers (object) optional - http headers


**Returns** `None`
- - - -

## DO
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/io/net/http/request.go#L26)

REQUEST makes a HTTP request.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`params (Object) - request parameters. * method (String) - HTTP method. * url (String) - Target url. * body` | `Binary` | Post data. * headers (object) optional - http headers.


**Returns** `None`
- - - -
