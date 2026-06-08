---
title: "Bind Parameters"
sidebarTitle: "Parameters"
weight: 40
draft: false
description: "Pass external values into Ferret queries with bind parameters."
aliases:
    - /docs/fql/bind-parameters/
---

# Bind Parameters

Bind parameters allow a query to receive values from the execution context instead of embedding those values directly in the query text.

A bind parameter is written as `@name`, where name is the parameter identifier.

{{< code lang="fql" >}}
LET users = [
    { name: “Ada”, active: true },
    { name: “Grace”, active: false },
    { name: “Linus”, active: true }
]

FOR user IN users
    FILTER user.active == @active
    RETURN user.name
{{</ code >}}

In this example, `@active` is not a boolean literal written in the query. It is a parameter whose value is provided when the query is executed.

## Parameter names

Parameter names must start with a letter or underscore. They may contain letters, digits, and underscores.

{{< code lang="fql" >}}
RETURN @name
RETURN @user_id
RETURN @value1
RETURN @_value
{{</ code >}}

The leading `@` is part of the parameter syntax, but not part of the parameter name itself.

## Passing parameter values

Parameter values are passed together with the query by the host application, runtime, or tool that executes the query.

They are not written as part of the query text.

A query fails if it references a bind parameter that was not provided.

{{< editor lang="fql" >}}
RETURN @name
{{</ editor >}}

In this example, the query expects a parameter named name.

Depending on the execution environment, passing extra parameters that are not referenced by the query may also be rejected.

## Parameters are values

Bind parameters represent values directly. They should not be wrapped in quotes.

{{< code lang="fql" height="100px" >}}
FILTER user.name == "@name" // Compares with the literal string "@name"
FILTER user.name == @name   // Compares with the parameter value
{{</ code >}}

Quoted text is always treated as a string literal. An unquoted @name expression refers to the parameter value.

## Using parameters in expressions

Bind parameters can be used anywhere a value expression is expected.

{{< code lang="fql" height="90px" >}}
RETURN CONCAT("prefix", @id, "suffix")
{{</ code >}}

The parameter value participates in the expression the same way as any other value.

## Dynamic property access

Bind parameters can also be used with bracket notation to access object properties dynamically.

{{< code lang="fql" height="120px" >}}
LET doc = {
foo: {
bar: "baz"
}
}

RETURN doc[@attr][@subattr]
{{</ code >}}

In this example, @attr and @subattr provide the property names used during evaluation.