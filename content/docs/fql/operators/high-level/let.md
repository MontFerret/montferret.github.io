---
title: "LET"
weight: 6
draft: false
----

{{< header href="let" >}}
LET
{{</ header >}}

The ``LET`` statement can be used to assign an arbitrary value to a variable. The variable is then introduced in the scope the ``LET`` statement is placed in.

The general syntax is:

{{< code lang="fql" height="80px" >}}
LET variableName = expression
{{</ code >}}

Variables are immutable in FQL, which means they can not be re-assigned:

{{< code lang="fql" height="150px" >}}
LET a = [1, 2, 3]  // initial assignment

a = PUSH(a, 4)     // syntax error, unexpected identifier
LET a = PUSH(a, 4) // parsing error, variable 'a' is assigned multiple times
LET b = PUSH(a, 4) // allowed, result: [1, 2, 3, 4]
{{</ code >}}

``LET`` statements are mostly used to declare complex computations and to avoid repeated computations of the same value at multiple parts of a query.

{{< editor height="240px" orientation="horizontal" >}}
LET users = [{"id":1,"firstName":"Kikelia","lastName":"Coper","cart":[{"Name":"Garlic","Price":"$7.46"},{"Name":"Flower - Commercial Spider","Price":"$6.59"}]},{"id":2,"firstName":"Toni","lastName":"MacTeggart","cart":[{"Name":"Spice - Paprika","Price":"$6.31"},{"Name":"Extract - Vanilla,artificial","Price":"$4.74"},{"Name":"Wine - White, Cooking","Price":"$1.50"},{"Name":"Nutmeg - Ground","Price":"$1.30"}]},{"id":3,"firstName":"Neile","lastName":"Saice","cart":[{"Name":"Mustard Prepared","Price":"$2.28"},{"Name":"Flower - Commercial Bronze","Price":"$4.80"}]}]

FOR u IN users
  LET numProducts = LENGTH(u.cart)
  
  RETURN { 
    "user" : u, 
    "numProducts" : numProducts, 
    "discount" : numProducts >= 3
  } 
{{</ editor >}}

In the above example, the computation of the number of cart items is factored out using a ``LET`` statement, thus avoiding computing the value twice in the ``RETURN`` statement.

Another use case for ``LET`` is to declare a complex computation in a subquery, making the whole query more readable.

{{< editor height="300px" orientation="horizontal" >}}
LET users = [{"id":1,"name":"Moises Grisewood"},{"id":2,"name":"Dell Marnes"},{"id":3,"name":"Tobin Bilbery"},{"id":4,"name":"Lorianne Posten"},{"id":5,"name":"Drucill Cryer"}]
LET friends = [{"id":1,"name":"Maximo Massard","userId":1},{"id":2,"name":"Delainey Sancho","userId":1},{"id":3,"name":"Lindon Beale","userId":1},{"id":4,"name":"Gus Sprey","userId":3},{"id":5,"name":"Virgil Dallander","userId":3},{"id":6,"name":"Agretha Mackerel","userId":4},{"id":7,"name":"Christalle Aldins","userId":4},{"id":8,"name":"Karalynn Margery","userId":5},{"id":9,"name":"Rodolph Ladd","userId":5},{"id":10,"name":"Babette Brassill","userId":1}]

FOR u IN users
  LET friends = (
    FOR f IN friends 
        FILTER u.id == f.userId
        RETURN f
  )

  RETURN { 
    "user" : u, 
    "friends" : friends, 
    "numFriends" : LENGTH(friends)
  }
{{</ editor >}}