---
title: "FOR"
weight: 1
draft: false
----

{{< header href="for" >}}
FOR
{{</ header >}}

The versatile FOR keyword can be used to iterate over a different collections like array, objects and set of HTML elements.

{{< header size="3" href="for-in" >}}
FOR-IN
{{</ header >}}

{{< code lang="fql" height="80px" >}}
FOR variableName IN expression
{{</ code >}}

A ``FOR-IN`` statement iterates through all entries of an array, object, or values that implement [``Iterable``]("https://github.com/MontFerret/ferret/blob/master/pkg/runtime/core/value.go#L22") interface. 
Each element returned by expression is visited exactly once. The current element is made available for further processing in the variable specified by ``variableName``.

{{< editor height="100px" orientation="horizontal" >}}
FOR i IN JSON_PARSE(IO::NET::HTTP::GET('http://country.io/continent.json'))
  RETURN i
{{</ editor >}}

This will iterate over all elements from the array and make the current array element available in variable ``i``. ``i`` is not modified in this example but simply pushed into the result using the ``RETURN`` keyword.

The variable introduced by ``FOR`` is available until the scope the ``FOR`` is placed in is closed.

Another example that uses a statically declared array of values to iterate over:

{{< editor height="100px" orientation="horizontal" >}}
FOR year IN [ 2011, 2012, 2013 ]
  RETURN { "year" : year, "isLeapYear" : year % 4 == 0 && (year % 100 != 0 || year % 400 == 0) }
{{</ editor >}}

Nesting of multiple ``FOR`` statements is allowed, too. When ``FOR`` statements are nested, a cross product of the array elements returned by the individual ``FOR`` statements will be created.

{{< editor height="100px" orientation="horizontal" >}}
FOR m IN [ 2019, 2020, 2021 ]
  FOR y IN [ "jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sept", "oct", "nov", "dec" ]
    RETURN { "year" : y, "month" : m }
{{</ editor >}}

In this example, there are two array iterations: an outer iteration over the array of months plus an inner iteration over the array of years. The inner array is traversed as many times as there are elements in the outer array. For each iteration, the current values of months and years are made available for further processing in the variable ``m`` and ``y``.

{{< header size="3" href="for-in" >}}
FOR-WHILE
{{</ header >}}

{{< code lang="fql" height="80px" >}}
FOR variableName [DO] WHILE expression
{{</ code >}}

``FOR-WHILE`` variant specifies the repeated iteration as long as a boolean condition evaluates to true. The condition is evaluated before each iteration. ``variableName`` variable holds a number of each iteration.