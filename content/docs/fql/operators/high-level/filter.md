---
title: "FILTER"
weight: 4
draft: false
----

{{< header href="filter" >}}
FILTER
{{</ header >}}

The ``FILTER`` statement can be used to restrict the results to elements that match an arbitrary logical condition.

The general syntax for ``FILTER`` is:

{{< code lang="fql" height="80px" >}}
FILTER expression
{{</ code >}}

``expression`` must be a condition that evaluates to either false or true. If the condition result is false, the current element is skipped, so it will not be processed further and not be part of the result. If the condition is true, the current element is not skipped and can be further processed. See [``Operators``]({{< ref "operators" >}} "Operators") for a list of comparison operators, logical operators etc. that you can use in conditions.

{{< editor height="100px" orientation="horizontal" >}}
FOR country IN JSON_PARSE(IO::NET::HTTP::GET("http://country.io/names.json"))
  FILTER country =~ 'A'
  RETURN country
{{</ code >}}

It is allowed to specify multiple ``FILTER`` statements in a query, even in the same block. If multiple ``FILTER`` statements are used, their results will be combined with a logical ``AND``, meaning all filter conditions must be true to include an element.

{{< editor height="100px" orientation="horizontal" >}}
FOR country IN JSON_PARSE(IO::NET::HTTP::GET("http://country.io/names.json"))
  FILTER country LIKE "A*"
  FILTER country LIKE "*a"
  RETURN country
{{</ code >}}

{{< header href="filter-order" size="2" >}}
Order of operations
{{</ header >}}

Note that the positions of ``FILTER`` statements can influence the result of a query. There are 250 countries in the test data for instance:

{{< editor height="100px" orientation="horizontal" >}}
FOR country IN JSON_PARSE(IO::NET::HTTP::GET("http://country.io/names.json"))
  FILTER country LIKE "A*"
  RETURN country
{{</ code >}}

We can limit the result set to 5 countries at most:

{{< editor height="100px" orientation="horizontal" >}}
FOR country IN JSON_PARSE(IO::NET::HTTP::GET("http://country.io/names.json"))
  FILTER LIKE(country, "A%")
  LIMIT 5
  RETURN country
{{</ code >}}

This may return strings Argentina, Anguilla, Armenia, Antarctica, American Samoa. Which ones are returned is undefined, since there is no SORT statement to ensure a particular order. If we add a second FILTER statement to only return women…