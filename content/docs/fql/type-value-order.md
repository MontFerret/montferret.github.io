---
title: "Type and value order"
weight: 6
draft: false
---

# Type and value order

When checking for equality or inequality or when determining the sort order of values, FQL uses a deterministic algorithm that takes both the data types and the actual values into account.

The compared operands are first compared by their data types, and only by their data values if the operands have the same data types.

The following type order is used when comparing data types:

{{< code lang="fql" height="90px" >}}
none < bool < number < string < array/list < object
{{</ code >}}

This means none is the smallest type in FQL and object is the type with the highest order. If the compared operands have a different type, then the comparison result is determined and the comparison is finished.

For example, the boolean true value will always be less than any numeric or string value, any array (even an empty array) or any object. Additionally, any string value (even an empty string) will always be greater than any numeric value, a boolean value, true or false

{{< code lang="fql" height="730px" >}}
none < false
none < true
none < 0
none < ''
none < ' '
none < '0'
none < 'abc'
none < [ ]
none < { }

false < true
false < 0
false < ''
false < ' '
false < '0'
false < 'abc'
false < [ ]
false < { }

true < 0
true < ''
true < ' '
true < '0'
true < 'abc'
true < [ ]
true < { }

0 < ''
0 < ' '
0 < '0'
0 < 'abc'
0 < [ ]
0 < { }

'' < ' '
'' < '0'
'' < 'abc'
'' < [ ]
'' < { }

[ ] < { }
{{</ code >}}

If the two compared operands have the same data types, then the operands values are compared. For the primitive types (none, boolean, number, and string), the result is defined as follows:

- none: none is equal to none
- boolean: false is less than true
- number: numeric values are ordered by their cardinal value
- string: string values are ordered using a localized comparison, using the configured server language for sorting according to the alphabetical order rules of that language

<div class="notification is-info">
  Unlike in SQL, none can be compared to any value, including none itself, without the result being converted into none automatically.
</div>

For compound, types the following special rules are applied:

Two array values are compared by comparing their individual elements position by position, starting at the first element. For each position, the element types are compared first. If the types are not equal, the comparison result is determined, and the comparison is finished. If the types are equal, then the values of the two elements are compared.

If an array element is itself a compound value (an array or an object), then the comparison algorithm will check the element's sub values recursively. The element's sub-elements are compared recursively.

{{< code lang="fql" height="120px" >}}
[ ] < [ 0 ]
[ 1 ] < [ 2 ]
[ false ] < [ true ]
{{</ code >}}

Two object operands are compared by checking attribute names and value. The attribute names are compared first. Before attribute names are compared, a combined array of all attribute names from both operands is created and sorted lexicographically. This means that the order in which attributes are declared in an object is not relevant when comparing two objects.

The combined and sorted array of attribute names is then traversed, and the respective attributes from the two compared operands are then looked up. If one of the objects does not have an attribute with the sought name, its attribute value is considered to be none. Finally, the attribute value of both objects is compared using the before mentioned data type and value comparison. The comparisons are performed for all object / document attributes until there is an unambiguous comparison result. If an unambiguous comparison result is found, the comparison is finished. If there is no unambiguous comparison result, the two compared objects are considered equal.

{{< code lang="fql" height="200px" >}}
{ } < { "a" : 1 }
{ } < { "a" : none }
{ "a" : 1 } < { "a" : 2 }
{ "b" : 1 } < { "a" : 0 }
{ "a" : { "c" : true } } < { "a" : { "c" : 0 } }
{ "a" : { "c" : true, "a" : 0 } } < { "a" : { "c" : false, "a" : 1 } }

{ "a" : 1, "b" : 2 } == { "b" : 2, "a" : 1 }
{{</ code >}}