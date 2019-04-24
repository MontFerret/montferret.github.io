---
title: "Data types"
weight: 3
draft: false
---

# Data types

FQL supports both primitive data types consisting of exactly one value and compound data types comprised of multiple values. The following types are available:

<table class="table">
    <thead>
        <tr>
            <th>Data type</th>
            <th>Description</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>none</td>
            <td>An empty value, also: the absence of a value</td>
        </tr>
        <tr>
            <td>boolean</td>
            <td>Boolean truth value with possible values false and true</td>
        </tr>
        <tr>
            <td>integer</td>
            <td>Signed whole number</td>
        </tr>
        <tr>
            <td>float</td>
            <td>Signed floating point number</td>
        </tr>
        <tr>
            <td>string</td>
            <td>UTF-8 encoded text value</td>
        </tr>
        <tr>
            <td>array</td>
            <td>Sequence of values, referred to by their positions</td>
        </tr>
        <tr>
            <td>object</td>
            <td>Sequence of values, referred to by their names</td>
        </tr>
        <tr>
            <td>binary</td>
            <td>Sequence of binary values</td>
        </tr>
        <tr>
            <td>custom</td>
            <td>User defined types</td>
        </tr>
    </tbody>
</table>

## Primitive types

### None value
A none value can be used to represent an empty or absent value. It is different from a numerical value of zero (null != 0) and other falsy values (false, zero-length string, empty array or object). It is also known as nil or null in other languages.

### Boolean data type
The Boolean data type has two possible values, ``true`` and ``false``. They represent the two truth values in logic and mathematics.

### Numeric literals
Numeric literals can be integers or floating-point numbers. They can optionally be signed with the + or - symbols. A decimal point . is used as separator for the optional fractional part. 

{{< highlight javascript >}}
1
+1
42
-1
-42
0.5
1.23
-99.99
{{</ highlight >}}

### String literals
String literals must be enclosed in single or double quotes. If the used quote character is to be used itself within the string literal, it must be escaped using the backslash symbol. A literal backslash also needs to be escaped with a backslash.

{{< highlight javascript >}}
"yikes!"
"don't know"
"this is a \"quoted\" word"
"this is a longer string."
"the path separator on Windows is \\"

'yikes!'
'don\'t know'
'this is a "quoted" word'
'this is a longer string.'
'the path separator on Windows is \\'
{{</ highlight >}}

## Compound types
FQL supports two compound types:

- array: A composition of unnamed values, each accessible by their positions.
- object: A composition of named values, each accessible by their names.
- binary: An array of binary values.

### Arrays
The first supported compound type is the array type. Arrays are effectively sequences of (unnamed / anonymous) values. Individual array elements can be accessed by their positions. The order of elements in an array is important.

An array declaration starts with a left square bracket ``[`` and ends with a right square bracket ``]``. The declaration contains zero, one or more expressions, separated from each other with the comma ``,`` symbol. Whitespace around elements is ignored in the declaration, thus line breaks, tab stops and blanks can be used for formatting.

In the easiest case, an array is empty and thus looks like:

{{< highlight javascript >}}
[ ]
{{</ highlight >}}

Array elements can be any legal expression values. Nesting of arrays is supported.

{{< highlight javascript >}}
[ true ]
[ 1, 2, 3 ]
[ -99, "yikes!", [ false, ["no"], [] ], 1 ]
[ [ "fox", "marshal" ] ]
{{</ highlight >}}

Individual array values can later be accessed by their positions using the [] accessor. The position of the accessed element must be a numeric value. Positions start at 0.

{{< highlight javascript >}}
// access 1st array element (elements start at index 0)
u.friends[0]

// access 3rd array element
u.friends[2]
{{</ highlight >}}

### Objects

The other supported compound type is the object type. Objects are a composition of zero to many attributes. Each attribute is a name/value pair. Object attributes can be accessed individually by their names. This data type is also known as dictionary, map, associative array and other names.

Object declarations start with a left curly bracket ``{`` and end with a right curly bracket ``}``. An object contains zero to many attribute declarations, separated from each other with the ``,`` symbol. Whitespace around elements is ignored in the declaration, thus line breaks, tab stops and blanks can be used for formatting.

In the simplest case, an object is empty. Its declaration would then be:

{{< highlight javascript >}}
{ }
{{</ highlight >}}

Each attribute in an object is a name/value pair. Name and value of an attribute are separated using the colon ``:`` symbol. The name is always a string, whereas the value can be of any type including sub-objects.

The attribute name is mandatory - there can't be anonymous values in an object. It can be specified as a quoted or unquoted string:

{{< highlight javascript >}}
{ name: … }    // unquoted
{ 'name': … }  // quoted (apostrophe / "single quote mark")
{ "name": … }  // quoted (quotation mark / "double quote mark")
{{</ highlight >}}

It must be quoted if it contains whitespace, escape sequences or characters other than ASCII letters (a-z, A-Z), digits (0-9), underscores (_) and dollar signs ($). The first character has to be a letter, underscore or dollar sign.

If a keyword is used as an attribute name then the attribute name must be quoted or escaped by backticks:

{{< highlight javascript >}}
{ return: … }    // error, return is a keyword!
{ 'return': … }  // quoted
{ "return": … }  // quoted
{ `return`: … }  // escaped (backticks)
{ ´return´: … }  // escaped (ticks)
{{</ highlight >}}

Attribute names can be computed using dynamic expressions, too. To disambiguate regular attribute names from attribute name expressions, computed attribute names must be enclosed in square brackets ``[ … ]``:

{{< highlight javascript >}}
{ [ CONCAT("test/", "bar") ] : "someValue" }
{{</ highlight >}}

There is also shorthand notation for attributes which is handy for returning existing variables easily:

{{< highlight javascript >}}
LET name = "Peter"
LET age = 42
RETURN { name, age }
{{</ highlight >}}

The above is the shorthand equivalent for the generic form:

{{< highlight javascript >}}
LET name = "Peter"
LET age = 42
RETURN { name: name, age: age }
{{</ highlight >}}

Any valid expression can be used as an attribute value. That also means nested objects can be used as attribute values:

{{< highlight javascript >}}
{ name : "Peter" }
{ "name" : "Vanessa", "age" : 15 }
{ "name" : "John", likes : [ "Swimming", "Skiing" ], "address" : { "street" : "Cucumber lane", "zip" : "94242" } }
{{</ highlight >}}

Individual object attributes can later be accessed by their names using the dot ``.`` accessor:

{{< highlight javascript >}}
u.address.city.name
u.friends[0].name.first
{{</ highlight >}}

Attributes can also be accessed using the square bracket ``[]`` accessor:

{{< highlight javascript >}}
u["address"]["city"]["name"]
u["friends"][0]["name"]["first"]
{{</ highlight >}}

In contrast to the dot accessor, the square brackets allow for expressions:

{{< highlight javascript >}}
LET attr1 = "friends"
LET attr2 = "name"
u[attr1][0][attr2][ CONCAT("fir", "st") ]
{{</ highlight >}}

## Custom types
Custom types are the values defined by a user or author of a 3rd party library that extends FQL functionality.
The values can represent either specific primitives or compund types. For example, tt could be http cookies, database records, images and etc. 

<div class="notification is-warning">
  The feature is a runtime specific.
  Currently, it's the Go runtime.
</div>

<div class="notification is-warning">
  At this moment, the feature is supported in embedded mode only and not available via CLI.
</div>

In order to define a custom value, it needs to have 2 Go types.
One should implement ``core.Type`` and another ``core.Value``.

{{< highlight go >}}
type Type interface {
    ID() int64
    String() string
    Equals(other Type) bool
}
{{</ highlight >}}

{{< highlight go >}}
type Value interface {
    json.Marshaler
    Type() Type
    String() string
    Compare(other Value) int64
    Unwrap() interface{}
    Hash() uint64
    Copy() Value
}
{{</ highlight >}}


For simplicity, there is a helper function in ``core`` package, that allows to create a type:

{{< highlight go >}}
var MyValueType = core.NewType("MyValueType")
{{</ highlight >}}

``core.Value`` interface provides a basic set of functionality which allows the FQL runtime to use a custom value. 

{{< highlight javascript >}}
LET myType = SOME_FUNC()

RETURN [myType]
{{</ highlight >}}

### Reading data
By default, implementing ``core.Value`` gives the FQL runtime just a basic set of functionality to work with values. This means that a type which implemented only ``core.Value`` interfaces is treated as a primitive data type.

If the underlying data type is compound and there is a need for providing access to underlying data, the type should implement ``core.Getter`` interface:

{{< highlight go >}}
type Getter interface {
    GetIn(ctx context.Context, path []Value) (Value, error)
}
{{</ highlight >}}

After that, it will be possible to use ``.`` or ``[]`` accessors:

{{< highlight javascript >}}
LET myType = SOME_FUNC()

RETURN myType.someProperty
{{</ highlight >}}

### Iterations
In order to use a custom type as source for ``FOR IN`` statement, it must implement ``core.Iterable`` interface:

{{< highlight go >}}
type Iterable interface {
    Iterate(ctx context.Context) (Iterator, error)
}
{{</ highlight >}}

{{< highlight go >}}
type Iterator interface {
    Next(ctx context.Context) (value Value, key Value, err error)
}
{{</ highlight >}}

### Clean up
If a custom type needs to do some clean up after a query execution, it needs to implement Go's ``io.Closer`` interface.