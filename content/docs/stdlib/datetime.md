---
title: "datetime"
weight: 1
draft: false
menuTitle: 
menu: [DATE,DATE_ADD,DATE_COMPARE,DATE_DAY,DATE_DAYOFWEEK,DATE_DAYOFYEAR,DATE_DAYS_IN_MONTH,DATE_DIFF,DATE_FORMAT,DATE_HOUR,DATE_LEAPYEAR,DATE_MILLISECOND,DATE_MINUTE,DATE_MONTH,DATE_QUARTER,DATE_SECOND,DATE_SUBTRACT,DATE_YEAR,NOW,]
---



{{< header >}}

DATE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/date.go#L14)

DATE converts RFC3339 date time string to DateTime object.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`time` | `String`  |  | String in rfc3339 format.


**Returns** `DateTime` New datetime object derived from timestring.
- - - -


{{< header >}}

DATE_ADD

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/add_subtract.go#L30)

DATE_ADD adds amount given in unit to date. The following units are available: * y, year, year * m, month, months * w, week, weeks * d, day, days * h, hour, hours * i, minute, minutes * s, second, seconds * f, millisecond, milliseconds

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source date.
`amount` | `Int`  |  | Amount of units
`unit` | `String`  |  | Unit.


**Returns** `DateTime` Calculated date.
- - - -


{{< header >}}

DATE_COMPARE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/compare.go#L17)

DATE_COMPARE checks if two partial dates match.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date1` | `DateTime`  |  | First date.
`date2` | `DateTime`  |  | Second date.
`unitRangeStart` | `String`  |  | Unit to start from.
`unitRangeEnd` | `String`  | `"millisecond"` | Unit to end with. error will be returned if unitrangestart unit less that unitrangeend.


**Returns** `Boolean` True if the dates match, else false.
- - - -


{{< header >}}

DATE_DAY

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/day.go#L13)

DATE_DAY returns the day of date as a number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Int` A day number.
- - - -


{{< header >}}

DATE_DAYOFWEEK

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/dayofweek.go#L13)

DATE_DAYOFWEEK returns number of the weekday from the date. Sunday is the 0th day of week.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Int` Number of the weekday.
- - - -


{{< header >}}

DATE_DAYOFYEAR

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/dayofyear.go#L14)

DATE_DAYOFYEAR returns the day of year number of date. The return value range from 1 to 365 (366 in a leap year).

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Int` A day of year number.
- - - -


{{< header >}}

DATE_DAYS_IN_MONTH

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/daysinmonth.go#L29)

DATE_DAYS_IN_MONTH returns the number of days in the month of date.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Int` Number of the days.
- - - -


{{< header >}}

DATE_DIFF

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/diff.go#L16)

DATE_DIFF returns the difference between two dates in given time unit.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date1` | `DateTime`  |  | First date.
`date2` | `DateTime`  |  | Second date.
`unit` | `String`  |  | Time unit to return the difference in.
`asFloat` | `Boolean`  | `False` | If true amount of unit will be as float.


**Returns** `Int` `Float` Difference between date1 and date2.
- - - -


{{< header >}}

DATE_FORMAT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/format.go#L13)

DATE_FORMAT format date according to the given format string.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime object.


**Returns** `String` Formatted date.
- - - -


{{< header >}}

DATE_HOUR

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/hour.go#L13)

DATE_HOUR returns the hour of date as a number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Int` An hour number.
- - - -


{{< header >}}

DATE_LEAPYEAR

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/leapyear.go#L13)

DATE_LEAPYEAR returns true if date is in a leap year else false.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Boolean` Date is in a leap year.
- - - -


{{< header >}}

DATE_MILLISECOND

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/millisecond.go#L13)

DATE_MILLISECOND returns the millisecond of date as a number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Int` A millisecond number.
- - - -


{{< header >}}

DATE_MINUTE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/minute.go#L13)

DATE_MINUTE returns the minute of date as a number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Int` A minute number.
- - - -


{{< header >}}

DATE_MONTH

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/month.go#L13)

DATE_MONTH returns the month of date as a number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Int` A month number.
- - - -


{{< header >}}

DATE_QUARTER

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/quarter.go#L14)

DATE_QUARTER returns which quarter date belongs to.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Int` A quarter number.
- - - -


{{< header >}}

DATE_SECOND

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/second.go#L13)

DATE_SECOND returns the second of date as a number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Int` A second number.
- - - -


{{< header >}}

DATE_SUBTRACT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/add_subtract.go#L60)

DATE_SUBTRACT subtract amount given in unit to date. The following units are available: * y, year, year * m, month, months * w, week, weeks * d, day, days * h, hour, hours * i, minute, minutes * s, second, seconds * f, millisecond, milliseconds

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source date.
`amount` | `Int`  |  | Amount of units
`unit` | `String`  |  | Unit.


**Returns** `DateTime` Calculated date.
- - - -


{{< header >}}

DATE_YEAR

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/year.go#L13)

DATE_YEAR returns the year extracted from the given date.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Int` A year number.
- - - -


{{< header >}}

NOW

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/now.go#L12)

NOW returns new DateTime object with Time equal to time.Now().

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description


**Returns** `DateTime` New datetime object.
- - - -
