---
title: "Bind parameters"
weight: 6
draft: false
---

{{< header >}}
Bind parameters
{{</ header >}}

FQL supports the usage of bind parameters, thus allowing to separate the query text from literal values used in the query. It is good practice to separate the query text from the literal values because it allows to reuse same query in different scenarios.

Using bind parameters, the meaning of an existing query cannot be changed. Bind parameters can be used everywhere in a query where literals can be used.

The syntax for bind parameters is ``@name`` where ``@`` signifies that this is a bind parameter and name is the actual parameter name. Parameter names must start with any of the letters a to z (upper or lower case) or a digit (0 to 9), and can be followed by any letter, digit or the underscore symbol.

{{< code lang="fql" height="190px" >}}
LET google = DOCUMENT("https://www.google.com/", true)

INPUT(google, 'input[name="q"]', @criteria)

WAIT_NAVIGATION(google)

RETURN ELEMENTS(google, '.g')
{{</ code >}}

The bind parameter values need to be passed along with the query when it is executed, but not as part of the query text itself.

Bind parameters that are declared in the query must also be passed a parameter value, or the query will fail. Specifying parameters that are not declared in the query will result in an error too.

Bind variables represent a value like a string, and must not be put in quotes in the FQL code:

{{< code lang="fql" height="100px" >}}
FILTER u.name == "@name" // wrong
FILTER u.name == @name   // correct
{{</ code >}}

If you need to do string processing in the query, you need to use string functions to do so:

{{< code lang="fql" height="90px" >}}
RETURN CONCAT('prefix', @id, 'suffix')
{{</ code >}}

Bind parameters can be used for square bracket notation for sub-attribute access. They can also be chained:

{{< code lang="fql" height="120px" >}}
LET doc = { foo: { bar: "baz" } }

RETURN doc[@attr][@subattr]
{{</ code >}}