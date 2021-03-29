---
title: "SORT"
weight: 4
draft: false
----

{{< header href="sort" >}}
SORT
{{</ header >}}

The ``SORT`` statement will force a sort of the array of already produced intermediate results in the current block. ``SORT`` allows specifying one or multiple sort criteria and directions. The general syntax is:

{{< code lang="fql" height="80px" >}}
SORT expression direction
{{</ code >}}

Example query that is sorting by lastName (in ascending order), then firstName (in ascending order), then by id (in descending order):

{{< editor height="160px" orientation="horizontal" >}}
LET users = [{"id":1,"firstName":"Johny","lastName":"Purdie"},{"id":2,"firstName":"Wayland","lastName":"Bewshaw"},{"id":3,"firstName":"Julius","lastName":"Taplin"},{"id":4,"firstName":"Jarrad","lastName":"Dollman"},{"id":5,"firstName":"Leia","lastName":"Meechan"}]

FOR u IN users
  SORT u.lastName, u.firstName, u.id DESC
  RETURN u
{{</ editor >}}

Specifying the direction is optional. The default (implicit) direction for a sort expression is the ascending order. To explicitly specify the sort direction, the keywords ``ASC`` (ascending) and ``DESC`` can be used. Multiple sort criteria can be separated using commas. In this case the direction is specified for each expression separately. For example:

{{< code lang="fql" height="80px" >}}
SORT u.lastName, u.firstName
{{</ code >}}

will first sort elements by lastName in ascending order and then by firstName in ascending order.

{{< code lang="fql" height="80px" >}}
SORT u.lastName DESC, u.firstName
{{</ code >}}

will first sort documents by lastName in descending order and then by firstName in ascending order.

{{< code lang="fql" height="80px" >}}
SORT u.lastName, u.firstName DESC
{{</ code >}}

will first sort elements by lastName in ascending order and then by firstName in descending order.
