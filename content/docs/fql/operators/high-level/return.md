---
title: "RETURN"
weight: 3
draft: false
----

{{< header href="return" >}}
RETURN
{{</ header >}}

The ``RETURN`` statement can be used to produce the result of a query. It is mandatory to specify a ``RETURN`` statement at the end of each block in a data-selection query, otherwise the query will be invalid. Using ``RETURN`` on the main level in data-modification queries is optional.

The general syntax for ``RETURN`` is:

{{< code lang="fql" height="80px" >}}
RETURN expression
{{</ code >}}

The expression returned by ``RETURN`` is produced for each iteration in the block the ``RETURN`` statement is placed in. That means the result of a ``RETURN`` statement is always an array when inside [``FOR``]({{< ref "for" >}} "FOR") loop.

To return all elements from the currently iterated array without modification, the following simple form can be used:

{{< code lang="fql" height="100px" >}}
FOR variableName IN expression
  RETURN variableName
{{</ code >}}

As ``RETURN`` allows specifying an expression, arbitrary computations can be performed to calculate the result elements. Any of the variables valid in the scope the ``RETURN`` is placed in can be used for the computations.

To iterate over all elements of an array of objects and return the full objects, you can write:

{{< editor height="100px" orientation="horizontal" >}}
FOR i IN [{name: 'Mike', age: 30}, {name: 'James', age: 35}]
  RETURN i
{{</ editor >}}

In each iteration of the for-loop, an element of the array is assigned to a variable ``i`` and returned unmodified in this example. To return only one attribute of each element, you could use a different return expression:

{{< editor height="100px" orientation="horizontal" >}}
FOR i IN [{name: 'Mike', age: 30}, {name: 'James', age: 35}]
  RETURN i.name
{{</ editor >}}

Or to return multiple attributes, an object can be constructed like this:

{{< editor height="100px" orientation="horizontal" >}}
FOR i IN [{firstName: 'Mike', lastName: 'Wazowski', age: 30}, {firstName: 'James', lastName: 'Sullivan', age: 35}]
  RETURN { name: i.firstName + ' ' + i.lastName, age: i.age }
{{</ editor >}}

Dynamic attribute names are supported as well:

{{< editor height="100px" orientation="horizontal" >}}
FOR i IN [{name: 'Mike', age: 30}, {name: 'James', age: 35}]
  RETURN { [i.name]: i.age }
{{</ editor >}}