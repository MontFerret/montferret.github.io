---
title: "LIMIT"
weight: 5
draft: false
----

{{< header href="limit" >}}
LIMIT
{{</ header >}}

The ``LIMIT`` statement allows slicing the result array using an offset and a count. It reduces the number of elements in the result to at most the specified number. Two general forms of ``LIMIT`` are followed:

{{< code lang="fql" height="100px" >}}
LIMIT count
LIMIT offset, count
{{</ code >}}

The first form allows specifying only the count value whereas the second form allows specifying both offset and count. The first form is identical using the second form with an offset value of 0.

{{< editor height="160px" orientation="horizontal" >}}
LET users = [{"id":1,"firstName":"Johny","lastName":"Purdie"},{"id":2,"firstName":"Wayland","lastName":"Bewshaw"},{"id":3,"firstName":"Julius","lastName":"Taplin"},{"id":4,"firstName":"Jarrad","lastName":"Dollman"},{"id":5,"firstName":"Leia","lastName":"Meechan"}]

FOR u IN users
  LIMIT 2
  RETURN u
{{</ editor >}}

Above query returns the first five elements of the array of user objects. It could also be written as ``LIMIT 0, 5`` for the same result. 

The ``offset`` value specifies how many elements from the result shall be skipped. It must be 0 or greater. The ``count`` value specifies how many elements should be at most included in the result.

{{< editor height="160px" orientation="horizontal" >}}
LET users = [{"id":1,"firstName":"Johny","lastName":"Purdie"},{"id":2,"firstName":"Wayland","lastName":"Bewshaw"},{"id":3,"firstName":"Julius","lastName":"Taplin"},{"id":4,"firstName":"Jarrad","lastName":"Dollman"},{"id":5,"firstName":"Leia","lastName":"Meechan"}]

FOR u IN users
  SORT u.firstName
  LIMIT 3, 2
  RETURN u
{{</ editor >}}

In above example, the objects of users are sorted, the first three results get skipped and it returns the next two user objects.

<div class="notification is-info">
    Note that variables, expressions and subqueries can not be used for offset and count. The values for offset and count must be known at query compile time, which means that you can only use number literals, bind parameters or expressions that can be resolved at query compile time.
</div>

Where a ``LIMIT`` is used in relation to other operations in a query has meaning. ``LIMIT`` operations before ``FILTER``s in particular can change the result significantly, because the operations are executed in the order in which they are written in the query. See [``FILTER``]({{< ref "filter" >}} "FILTER") for a detailed example.