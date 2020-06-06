---
title: "Syntax"
weight: 2
draft: false
---

# Syntax

### Query types
An FQL query must return a result indicated by usage either of the ``RETURN`` or ``FOR IN`` keywords. The FQL parser will return an error if it cannot find any of these two statements.

### Whitespace
Whitespaces (blanks, carriage returns, line feeds, and tab stops) can be used in the query text to increase its readability. Tokens have to be separated by any number of whitespaces. Whitespace within strings or names must be enclosed in quotes in order to be preserved.

### Comments
Comments can be embedded at any position in a query. The text contained in the comment is ignored by the FQL parser.

Multi-line comments cannot be nested, which means subsequent comment starts within comments are ignored, comment ends will end the comment.

FQL supports two types of comments:

- Single line comments: These start with a double forward slash and end at the end of the line, or the end of the query string (whichever is first).
- Multi line comments: These start with a forward slash and asterisk, and end with an asterisk and a following forward slash. They can span as many lines as necessary.

{{< code lang="fql" height="170px" >}}
/* this is a comment */ RETURN 1
/* these */ RETURN /* are */ 1 /* multiple */ + /* comments */ 1
/* this is
   a multi line
   comment */
// a single line comment
{{</ code >}}

### Keywords
On the top level, FQL offers the following operations:

- ``USE``: allows the use of types in a namespace without typing full type path.
- ``FOR``: array iteration
- ``RETURN``: results projection
- ``FILTER``: non-view results filtering
- ``SEARCH``: view results filtering
- ``SORT``: result sorting
- ``LIMIT``: result slicing
- ``LET``: variable assignment
- ``COLLECT``: result grouping

Each of the above operations can be initiated in a query by using a keyword of the same name. An FQL query can (and typically does) consist of multiple of the above operations.

An example FQL query may look like this:

{{< code lang="fql" height="180px" >}}
LET page = DOCUMENT("https://github.com/trending")

FOR row IN ELEMENTS(page, "ol.repo-list li")
    LET name = INNER_TEXT(row, "div:nth-child(1)")
    LET description = INNER_TEXT(row, "div:nth-child(3)")
    
    RETURN { name, description }
{{</ code >}}

In this example query, the terms ``FOR``, ``FILTER``, and ``RETURN`` initiate the higher-level operation according to their name. These terms are also keywords, meaning that they have a special meaning in the language.

For example, the query parser will use the keywords to find out which high-level operations to execute. That also means keywords can only be used at certain locations in a query. This also makes all keywords reserved words that must not be used for other purposes than they are intended for.

At this moment, keywords are not case-insensitive, meaning they cannot be specified in lower or mixed case in queries.

There are a few more keywords in addition to the higher-level operation keywords. Additional keywords may be added in future versions of Ferret. The complete list of keywords is currently:

<div class="columns">
    <div class="column">
    <ul>
        <li>AGGREGATE</li>
        <li>ALL</li>
        <li>AND</li>
        <li>ANY</li>
        <li>ASC</li>
        <li>COLLECT</li>
        <li>DESC</li>
        <li>DISTINCT</li>
        <li>FALSE</li>
        <li>FILTER</li>
        <li>FOR</li>
    </ul>
    </div>
    <div class="column">
        <ul>
            <li>IN</li>
            <li>LET</li>
            <li>LIMIT</li>
            <li>NONE</li>
            <li>NOT</li>
            <li>NULL</li>
            <li>OR</li>
            <li>RETURN</li>
            <li>SORT</li>
            <li>TRUE</li>
        </ul>
    </div>
</div>

### Names
In general, names are used to identify objects (properties, variables, and functions) in FQL queries.

The maximum supported length of any name is 64 bytes. Names in FQL are always case-sensitive.

Keywords must not be used as names. If a reserved keyword should be used as a name, the name must be enclosed in single or double quotes. Enclosing a name in quotes makes it possible to use otherwise reserved keywords as names. An example for this is:

{{< code lang="fql" height="100px" >}}
FOR i IN [{ "RETURN": "foobar" }]
    RETURN i."RETURN"
{{</ code >}}

#### Variable names
FQL allows the user to assign values to additional variables in a query. All variables that are assigned a value must have a name that is unique within the context of the query.

{{< code lang="fql" height="120px" >}}
LET users = [{ name: "Steve" }]
FOR u IN users
  RETURN { name : u.name }
{{</ code >}}

Allowed characters in variable names are the letters a to z (both in lower and upper case), the numbers 0 to 9, the underscore (_) symbol and the dollar ($) sign. A variable name must not start with a number. If a variable name starts with the underscore character, the underscore must be followed by least one letter (a-z or A-Z) or digit (0-9).