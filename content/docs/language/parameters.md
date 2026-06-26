---
title: "Bind Parameters"
sidebarTitle: "Parameters"
weight: 100
draft: false
description: "Pass external values into Ferret queries with bind parameters."
aliases:
    - /docs/fql/bind-parameters/
---

# Bind Parameters

Bind parameters allow a query to receive values from the execution context instead of embedding those values directly in the query text.

A bind parameter is written as `@name`, where name is the parameter identifier.

<div class="is-info notification">
Click on <b>PARAMS</b> to open the parameter editor and provide values for the parameters used in the query.
</div>

{{< editor lang="fql" params=`{"active": true}` >}}
LET users = [
    { name: "Ada", active: true },
    { name: "Grace", active: false },
    { name: "Linus", active: true }
]

FOR user IN users
    FILTER user.active == @active
    RETURN user.name
{{</ editor >}}

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

In this example, the query expects a parameter named `name`.

## Parameters are values

Bind parameters represent values directly. They should not be wrapped in quotes.

{{< editor lang="fql" params=`{"name": "Grace"}` >}}
LET users = [
    { name: "Ada" },
    { name: "Grace" },
    { name: "Linus" }
]

FOR user IN users
    FILTER user.name == @name
    RETURN user.name
{{</ editor >}}

Quoted text is always treated as a string literal. An unquoted @name expression refers to the parameter value.

## Using parameters in expressions

Bind parameters can be used anywhere a value expression is expected.

{{< editor lang="fql" params=`{"id": "F03D4CB9"}` >}}
RETURN CONCAT("prefix-", @id, "-suffix")
{{</ editor >}}

The parameter value participates in the expression the same way as any other value.

## Dynamic property access

Bind parameters can also be used with bracket notation to access object properties dynamically.

{{< editor lang="fql" params=`{"attr": "foo", "subattr": "bar"}` >}}
LET doc = {
    foo: {
        bar: "baz"
    }
}

RETURN doc[@attr][@subattr]
{{</ editor >}}

## Next steps

{{< docs-related tiles="language-functions,embedding-parameters,language-control-flow" >}}