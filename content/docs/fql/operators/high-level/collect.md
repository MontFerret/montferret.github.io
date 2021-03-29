---
title: "COLLECT"
weight: 7
draft: false
----

{{< header href="collect" >}}
COLLECT
{{</ header >}}

The ``COLLECT`` keyword can be used to group an array by one or multiple group criteria.

The ``COLLECT`` statement will eliminate all local variables in the current scope. After ``COLLECT`` only the variables introduced by ``COLLECT`` itself are available.

There are several syntax variants for ``COLLECT`` operations:

{{< code lang="fql" height="230px" >}}
COLLECT variableName = expression
COLLECT variableName = expression INTO groupsVariable
COLLECT variableName = expression INTO groupsVariable = projectionExpression
COLLECT variableName = expression INTO groupsVariable KEEP keepVariable
COLLECT variableName = expression WITH COUNT INTO countVariable
COLLECT variableName = expression AGGREGATE variableName = aggregateExpression
COLLECT variableName = expression AGGREGATE variableName = aggregateExpression INTO groupsVariable
COLLECT AGGREGATE variableName = aggregateExpression
COLLECT AGGREGATE variableName = aggregateExpression INTO groupsVariable
COLLECT WITH COUNT INTO countVariable
{{</ code >}}

{{< header size="3" href="collect-grouping" >}}
Grouping syntaxes
{{</ header >}}

The first syntax form of ``COLLECT`` only groups the result by the defined group criteria specified in expression. In order to further process the results produced by ``COLLECT``, a new variable (specified by ``variableName``) is introduced. This variable contains the group value.

Here’s an example query that find the distinct values in u.city and makes them available in variable city:

{{< editor height="200px" orientation="horizontal" >}}
LET users = [{"id":1,"name":"Thorn De Banke","city":"New York City"},{"id":2,"name":"Adolph MacAindreis","city":"Moscow"},{"id":3,"name":"Flory Maxworthy","city":"London"},{"id":4,"name":"Merola Blundell","city":"New York City"},{"id":5,"name":"Benson Manoelli","city":"Moscow"},{"id":6,"name":"Cordy McVie","city":"London"},{"id":7,"name":"Reggy Tuck","city":"New York City"},{"id":8,"name":"Courtney Toogood","city":"Berlin"},{"id":9,"name":"Berkley Fields","city":"Berlin"},{"id":10,"name":"Findley Gauvain","city":"Stockholm"}]

FOR u IN users
  COLLECT city = u.city
  
  RETURN { 
    "city" : city 
  }
{{</ editor >}}

The second form does the same as the first form, but additionally introduces a variable (specified by ``groupsVariable``) that contains all elements that fell into the group. This works as follows: The ``groupsVariable`` variable is an array containing as many elements as there are in the group. Each member of that array is a JSON object in which the value of every variable that is defined in the FQL query is bound to the corresponding property. Note that this considers all variables that are defined before the ``COLLECT`` statement, but not those on the top level (outside of any FOR), unless the ``COLLECT`` statement is itself on the top level, in which case all variables are taken. 

{{< editor height="200px" orientation="horizontal" >}}
LET users = [{"id":1,"name":"Thorn De Banke","city":"New York City"},{"id":2,"name":"Adolph MacAindreis","city":"Moscow"},{"id":3,"name":"Flory Maxworthy","city":"London"},{"id":4,"name":"Merola Blundell","city":"New York City"},{"id":5,"name":"Benson Manoelli","city":"Moscow"},{"id":6,"name":"Cordy McVie","city":"London"},{"id":7,"name":"Reggy Tuck","city":"New York City"},{"id":8,"name":"Courtney Toogood","city":"Berlin"},{"id":9,"name":"Berkley Fields","city":"Berlin"},{"id":10,"name":"Findley Gauvain","city":"Stockholm"}]

FOR u IN users
  COLLECT city = u.city INTO groups
  
  RETURN { 
    "city" : city, 
    "usersInCity" : groups 
  }
{{</ editor >}}

In the above example, the array users will be grouped by the property city. The result is a new array of objects, with one object per distinct ``u.city`` value. The elements from the original array (here: ``users``) per city are made available in the variable groups. This is due to the ``INTO`` clause.

``COLLECT`` also allows specifying multiple group criteria. Individual group criteria can be separated by commas:

{{< editor height="200px" orientation="horizontal" >}}
LET users = [{"id":1,"name":"Thorn De Banke","city":"New York City","country":"USA"},{"id":2,"name":"Adolph MacAindreis","city":"Moscow","country":"Russia"},{"id":3,"name":"Flory Maxworthy","city":"London","country":"UK"},{"id":4,"name":"Merola Blundell","city":"New York City","country":"USA"},{"id":5,"name":"Benson Manoelli","city":"Moscow","country":"Russia"},{"id":6,"name":"Cordy McVie","city":"London","country":"UK"},{"id":7,"name":"Reggy Tuck","city":"New York City","country":"USA"},{"id":8,"name":"Courtney Toogood","city":"Berlin","country":"Germany"},{"id":9,"name":"Berkley Fields","city":"Berlin","country":"Germany"},{"id":10,"name":"Findley Gauvain","city":"Stockholm","country":"Sweden"}]

FOR u IN users
  COLLECT country = u.country, city = u.city INTO groups
  
  RETURN { 
    "country" : country, 
    "city" : city, 
    "usersInCity" : groups 
  }
{{</ editor >}}

In the above example, the array users is grouped by country first and then by city, and for each distinct combination of country and city, the users will be returned.

{{< header size="3" href="collect-obsolete-variables" >}}
Discarding obsolete variables
{{</ header >}}

The third form of ``COLLECT`` allows rewriting the contents of the ``groupsVariable`` using an arbitrary ``projectionExpression``:

{{< editor height="200px" orientation="horizontal" >}}
LET users = [{"id":1,"name":"Thorn De Banke","city":"New York City","country":"USA"},{"id":2,"name":"Adolph MacAindreis","city":"Moscow","country":"Russia"},{"id":3,"name":"Flory Maxworthy","city":"London","country":"UK"},{"id":4,"name":"Merola Blundell","city":"New York City","country":"USA"},{"id":5,"name":"Benson Manoelli","city":"Moscow","country":"Russia"},{"id":6,"name":"Cordy McVie","city":"London","country":"UK"},{"id":7,"name":"Reggy Tuck","city":"New York City","country":"USA"},{"id":8,"name":"Courtney Toogood","city":"Berlin","country":"Germany"},{"id":9,"name":"Berkley Fields","city":"Berlin","country":"Germany"},{"id":10,"name":"Findley Gauvain","city":"Stockholm","country":"Sweden"}]

FOR u IN users
  COLLECT country = u.country, city = u.city INTO groups = u.name
  
  RETURN { 
    "country" : country, 
    "city" : city, 
    "userNames" : groups 
  }
{{</ editor >}}

In the above example, only the ``projectionExpression`` is ``u.name``. Therefore, only this property is copied into the ``groupsVariable`` for each object. This is probably much more efficient than copying all variables from the scope into the ``groupsVariable`` as it would happen without a ``projectionExpression``.

The expression following ``INTO`` can also be used for arbitrary computations:

{{< editor height="250px" orientation="horizontal" >}}
LET users = [{"id":1,"name":"Thorn De Banke","city":"New York City","country":"USA","status":"active"},{"id":2,"name":"Adolph MacAindreis","city":"Moscow","country":"Russia","status":"active"},{"id":3,"name":"Flory Maxworthy","city":"London","country":"UK","status":"active"},{"id":4,"name":"Merola Blundell","city":"New York City","country":"USA","status":"inactive"},{"id":5,"name":"Benson Manoelli","city":"Moscow","country":"Russia","status":"inactive"},{"id":6,"name":"Cordy McVie","city":"London","country":"UK","status":"active"},{"id":7,"name":"Reggy Tuck","city":"New York City","country":"USA","status":"active"},{"id":8,"name":"Courtney Toogood","city":"Berlin","country":"Germany","status":"active"},{"id":9,"name":"Berkley Fields","city":"Berlin","country":"Germany","status":"inactive"},{"id":10,"name":"Findley Gauvain","city":"Stockholm","country":"Sweden","status":"inactive"}]

FOR u IN users
  COLLECT country = u.country, city = u.city INTO groups = { 
    "name" : u.name, 
    "isActive" : u.status == "active"
  }
  RETURN { 
    "country" : country, 
    "city" : city, 
    "usersInCity" : groups 
  }
{{</ editor >}}

<!-- ``COLLECT`` also provides an optional ``KEEP`` clause that can be used to control which variables will be copied into the variable created by ``INTO``. If no ``KEEP`` clause is specified, all variables from the scope will be copied as sub-attributes into the ``groupsVariable``. This is safe but can have a negative impact on performance if there are many variables in scope or the variables contain massive amounts of data.

The following example limits the variables that are copied into the ``groupsVariable`` to just name. The variables ``u`` and ``someCalculation`` also present in the scope will not be copied into ``groupsVariable`` because they are not listed in the ``KEEP`` clause:

{{< editor height="250px" orientation="horizontal" >}}
LET users = [{"id":1,"name":"Thorn De Banke","city":"New York City","country":"USA","status":"active"},{"id":2,"name":"Adolph MacAindreis","city":"Moscow","country":"Russia","status":"active"},{"id":3,"name":"Flory Maxworthy","city":"London","country":"UK","status":"active"},{"id":4,"name":"Merola Blundell","city":"New York City","country":"USA","status":"inactive"},{"id":5,"name":"Benson Manoelli","city":"Moscow","country":"Russia","status":"inactive"},{"id":6,"name":"Cordy McVie","city":"London","country":"UK","status":"active"},{"id":7,"name":"Reggy Tuck","city":"New York City","country":"USA","status":"active"},{"id":8,"name":"Courtney Toogood","city":"Berlin","country":"Germany","status":"active"},{"id":9,"name":"Berkley Fields","city":"Berlin","country":"Germany","status":"inactive"},{"id":10,"name":"Findley Gauvain","city":"Stockholm","country":"Sweden","status":"inactive"}]

FOR u IN users
  LET name = u.name
  LET someCalculation = u.value1 + u.value2
  COLLECT city = u.city INTO groups KEEP name 
  
  RETURN { 
    "city" : city, 
    "userNames" : groups[*].name 
  }
{{</ editor >}} -->

{{< header size="3" href="collect-group-length" >}}
Group length calculation
{{</ header >}}

``COLLECT`` also provides a special ``WITH COUNT`` clause that can be used to determine the number of group members efficiently.

The simplest form just returns the number of items that made it into the ``COLLECT``:

{{< editor height="180px" orientation="horizontal" >}}
LET users = [{"id":1,"name":"Thorn De Banke","city":"New York City","country":"USA","status":"active"},{"id":2,"name":"Adolph MacAindreis","city":"Moscow","country":"Russia","status":"active"},{"id":3,"name":"Flory Maxworthy","city":"London","country":"UK","status":"active"},{"id":4,"name":"Merola Blundell","city":"New York City","country":"USA","status":"inactive"},{"id":5,"name":"Benson Manoelli","city":"Moscow","country":"Russia","status":"inactive"},{"id":6,"name":"Cordy McVie","city":"London","country":"UK","status":"active"},{"id":7,"name":"Reggy Tuck","city":"New York City","country":"USA","status":"active"},{"id":8,"name":"Courtney Toogood","city":"Berlin","country":"Germany","status":"active"},{"id":9,"name":"Berkley Fields","city":"Berlin","country":"Germany","status":"inactive"},{"id":10,"name":"Findley Gauvain","city":"Stockholm","country":"Sweden","status":"inactive"}]

FOR u IN users
  COLLECT WITH COUNT INTO length
  
  RETURN length
{{</ editor >}}

The above is equivalent to, but less efficient than:

{{< editor height="130px" orientation="horizontal" >}}
LET users = [{"id":1,"name":"Thorn De Banke","city":"New York City","country":"USA","status":"active"},{"id":2,"name":"Adolph MacAindreis","city":"Moscow","country":"Russia","status":"active"},{"id":3,"name":"Flory Maxworthy","city":"London","country":"UK","status":"active"},{"id":4,"name":"Merola Blundell","city":"New York City","country":"USA","status":"inactive"},{"id":5,"name":"Benson Manoelli","city":"Moscow","country":"Russia","status":"inactive"},{"id":6,"name":"Cordy McVie","city":"London","country":"UK","status":"active"},{"id":7,"name":"Reggy Tuck","city":"New York City","country":"USA","status":"active"},{"id":8,"name":"Courtney Toogood","city":"Berlin","country":"Germany","status":"active"},{"id":9,"name":"Berkley Fields","city":"Berlin","country":"Germany","status":"inactive"},{"id":10,"name":"Findley Gauvain","city":"Stockholm","country":"Sweden","status":"inactive"}]

RETURN LENGTH(users)
{{</ editor >}}

The ``WITH COUNT`` clause can also be used to efficiently count the number of items in each group:

{{< editor height="220px" orientation="horizontal" >}}
LET users = [{"id":1,"name":"Thorn De Banke","city":"New York City","country":"USA","status":"active","age":30},{"id":2,"name":"Adolph MacAindreis","city":"Moscow","country":"Russia","status":"active","age":25},{"id":3,"name":"Flory Maxworthy","city":"London","country":"UK","status":"active","age":40},{"id":4,"name":"Merola Blundell","city":"New York City","country":"USA","status":"inactive","age":25},{"id":5,"name":"Benson Manoelli","city":"Moscow","country":"Russia","status":"inactive","age":30},{"id":6,"name":"Cordy McVie","city":"London","country":"UK","status":"active","age":25},{"id":7,"name":"Reggy Tuck","city":"New York City","country":"USA","status":"active","age":18},{"id":8,"name":"Courtney Toogood","city":"Berlin","country":"Germany","status":"active","age":30},{"id":9,"name":"Berkley Fields","city":"Berlin","country":"Germany","status":"inactive","age":30},{"id":10,"name":"Findley Gauvain","city":"Stockholm","country":"Sweden","status":"inactive","age":30}]

FOR u IN users
  COLLECT age = u.age WITH COUNT INTO length
  
  RETURN { 
    "age" : age, 
    "count" : length 
  }
{{</ editor >}}

<div class="notification is-info">
    The WITH COUNT clause can only be used together with an INTO clause.
</div>

{{< header size="3" href="collect-aggregation" >}}
Aggregation
{{</ header >}}

A ``COLLECT`` statement can be used to perform aggregation of data per group. To only determine group lengths, the ``WITH COUNT INTO`` variant of ``COLLECT`` can be used as described before.

For other aggregations, it is possible to run aggregate functions on the ``COLLECT`` results:

{{< editor height="310px" orientation="horizontal" >}}
LET users = [{"id":1,"name":"Thorn De Banke","city":"New York City","country":"USA","status":"active","age":30},{"id":2,"name":"Adolph MacAindreis","city":"Moscow","country":"Russia","status":"active","age":25},{"id":3,"name":"Flory Maxworthy","city":"London","country":"UK","status":"active","age":40},{"id":4,"name":"Merola Blundell","city":"New York City","country":"USA","status":"inactive","age":25},{"id":5,"name":"Benson Manoelli","city":"Moscow","country":"Russia","status":"inactive","age":30},{"id":6,"name":"Cordy McVie","city":"London","country":"UK","status":"active","age":25},{"id":7,"name":"Reggy Tuck","city":"New York City","country":"USA","status":"active","age":18},{"id":8,"name":"Courtney Toogood","city":"Berlin","country":"Germany","status":"active","age":30},{"id":9,"name":"Berkley Fields","city":"Berlin","country":"Germany","status":"inactive","age":30},{"id":10,"name":"Findley Gauvain","city":"Stockholm","country":"Sweden","status":"inactive","age":30}]

FOR u IN users
  COLLECT ageGroup = FLOOR(u.age / 5) * 5 INTO g

  LET ages = (
      FOR i IN g
        RETURN i.u.age
  )
  
  RETURN { 
    "ageGroup" : ageGroup,
    "minAge" : MIN(ages),
    "maxAge" : MAX(ages)
  }
{{</ editor >}}

The above however requires storing all group values during the collect operation for all groups, which can be inefficient.

The special ``AGGREGATE`` variant of ``COLLECT`` allows building the aggregate values incrementally during the collect operation, and is therefore often more efficient.

With the ``AGGREGATE`` variant the above query becomes:

{{< editor height="250px" orientation="horizontal" >}}
LET users = [{"id":1,"name":"Thorn De Banke","city":"New York City","country":"USA","status":"active","age":30},{"id":2,"name":"Adolph MacAindreis","city":"Moscow","country":"Russia","status":"active","age":25},{"id":3,"name":"Flory Maxworthy","city":"London","country":"UK","status":"active","age":40},{"id":4,"name":"Merola Blundell","city":"New York City","country":"USA","status":"inactive","age":25},{"id":5,"name":"Benson Manoelli","city":"Moscow","country":"Russia","status":"inactive","age":30},{"id":6,"name":"Cordy McVie","city":"London","country":"UK","status":"active","age":25},{"id":7,"name":"Reggy Tuck","city":"New York City","country":"USA","status":"active","age":18},{"id":8,"name":"Courtney Toogood","city":"Berlin","country":"Germany","status":"active","age":30},{"id":9,"name":"Berkley Fields","city":"Berlin","country":"Germany","status":"inactive","age":30},{"id":10,"name":"Findley Gauvain","city":"Stockholm","country":"Sweden","status":"inactive","age":30}]

FOR u IN users
  COLLECT ageGroup = FLOOR(u.age / 5) * 5 
  AGGREGATE minAge = MIN(u.age), maxAge = MAX(u.age)
  
  RETURN {
    ageGroup, 
    minAge, 
    maxAge 
  }
{{</ editor >}}

The ``AGGREGATE`` keyword can only be used after the ``COLLECT`` keyword. If used, it must directly follow the declaration of the grouping keys. If no grouping keys are used, it must follow the ``COLLECT`` keyword directly:

{{< editor height="250px" orientation="horizontal" >}}
LET users = [{"id":1,"name":"Thorn De Banke","city":"New York City","country":"USA","status":"active","age":30},{"id":2,"name":"Adolph MacAindreis","city":"Moscow","country":"Russia","status":"active","age":25},{"id":3,"name":"Flory Maxworthy","city":"London","country":"UK","status":"active","age":40},{"id":4,"name":"Merola Blundell","city":"New York City","country":"USA","status":"inactive","age":25},{"id":5,"name":"Benson Manoelli","city":"Moscow","country":"Russia","status":"inactive","age":30},{"id":6,"name":"Cordy McVie","city":"London","country":"UK","status":"active","age":25},{"id":7,"name":"Reggy Tuck","city":"New York City","country":"USA","status":"active","age":18},{"id":8,"name":"Courtney Toogood","city":"Berlin","country":"Germany","status":"active","age":30},{"id":9,"name":"Berkley Fields","city":"Berlin","country":"Germany","status":"inactive","age":30},{"id":10,"name":"Findley Gauvain","city":"Stockholm","country":"Sweden","status":"inactive","age":30}]

FOR u IN users
  COLLECT AGGREGATE minAge = MIN(u.age), maxAge = MAX(u.age)
  
  RETURN {
    minAge, 
    maxAge 
  }
{{</ editor >}}

Only specific expressions are allowed on the right-hand side of each AGGREGATE assignment:

- an aggregate expression must not refer to variables introduced by the COLLECT itself