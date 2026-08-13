---
title: "Datetime"
weight: 60
draft: false
description: "Date and time functions in the Ferret standard library."
aliases:
  - /docs/stdlib/datetime/
menuTitle: 
menu: [date,date_add,date_compare,date_day,date_dayofweek,date_dayofyear,date_days_in_month,date_diff,date_format,date_hour,date_leapyear,date_millisecond,date_minute,date_month,date_quarter,date_second,date_subtract,date_year,now,]
---

DateTime values also support checked `+` and `-` operators with native Duration values. Subtracting two DateTime values returns a Duration. Convert other inputs explicitly with `to_duration` before applying an operator. These operators use canonical instants and raise range errors on overflow; the functions below remain available for calendar-unit operations.



{{< header href="date" >}}

date

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/date.go#L14)

date converts RFC3339 date time string to DateTime object.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`time` | `String`  |  | String in rfc3339 format.


**Returns** `DateTime` New datetime object derived from timestring.
- - - -


{{< header href="date_add" >}}

date_add

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/add_subtract.go#L30)

date_add adds amount given in unit to date. The following units are available: * y, year, year * m, month, months * w, week, weeks * d, day, days * h, hour, hours * i, minute, minutes * s, second, seconds * f, millisecond, milliseconds

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source date.
`amount` | `Int`  |  | Amount of units
`unit` | `String`  |  | Unit.


**Returns** `DateTime` Calculated date.
- - - -


{{< header href="date_compare" >}}

date_compare

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/compare.go#L17)

date_compare checks if two partial dates match.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date1` | `DateTime`  |  | First date.
`date2` | `DateTime`  |  | Second date.
`unitRangeStart` | `String`  |  | Unit to start from.
`unitRangeEnd` | `String`  | `"millisecond"` | Unit to end with. error will be returned if unitrangestart unit less that unitrangeend.


**Returns** `Boolean` True if the dates match, else false.
- - - -


{{< header href="date_day" >}}

date_day

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/day.go#L13)

date_day returns the day of date as a number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Int` A day number.
- - - -


{{< header href="date_dayofweek" >}}

date_dayofweek

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/dayofweek.go#L13)

date_dayofweek returns number of the weekday from the date. Sunday is the 0th day of week.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Int` Number of the weekday.
- - - -


{{< header href="date_dayofyear" >}}

date_dayofyear

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/dayofyear.go#L14)

date_dayofyear returns the day of year number of date. The return value range from 1 to 365 (366 in a leap year).

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Int` A day of year number.
- - - -


{{< header href="date_days_in_month" >}}

date_days_in_month

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/daysinmonth.go#L29)

date_days_in_month returns the number of days in the month of date.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Int` Number of the days.
- - - -


{{< header href="date_diff" >}}

date_diff

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/diff.go#L16)

date_diff returns the difference between two dates in given time unit.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date1` | `DateTime`  |  | First date.
`date2` | `DateTime`  |  | Second date.
`unit` | `String`  |  | Time unit to return the difference in.
`asFloat` | `Boolean`  | `False` | If true amount of unit will be as float.


**Returns** `Int` `Float` Difference between date1 and date2.
- - - -


{{< header href="date_format" >}}

date_format

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/format.go#L13)

date_format format date according to the given format string.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime object.


**Returns** `String` Formatted date.
- - - -


{{< header href="date_hour" >}}

date_hour

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/hour.go#L13)

date_hour returns the hour of date as a number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Int` An hour number.
- - - -


{{< header href="date_leapyear" >}}

date_leapyear

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/leapyear.go#L13)

date_leapyear returns true if date is in a leap year else false.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Boolean` Date is in a leap year.
- - - -


{{< header href="date_millisecond" >}}

date_millisecond

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/millisecond.go#L13)

date_millisecond returns the millisecond of date as a number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Int` A millisecond number.
- - - -


{{< header href="date_minute" >}}

date_minute

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/minute.go#L13)

date_minute returns the minute of date as a number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Int` A minute number.
- - - -


{{< header href="date_month" >}}

date_month

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/month.go#L13)

date_month returns the month of date as a number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Int` A month number.
- - - -


{{< header href="date_quarter" >}}

date_quarter

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/quarter.go#L14)

date_quarter returns which quarter date belongs to.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Int` A quarter number.
- - - -


{{< header href="date_second" >}}

date_second

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/second.go#L13)

date_second returns the second of date as a number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Int` A second number.
- - - -


{{< header href="date_subtract" >}}

date_subtract

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/add_subtract.go#L60)

date_subtract subtract amount given in unit to date. The following units are available: * y, year, year * m, month, months * w, week, weeks * d, day, days * h, hour, hours * i, minute, minutes * s, second, seconds * f, millisecond, milliseconds

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source date.
`amount` | `Int`  |  | Amount of units
`unit` | `String`  |  | Unit.


**Returns** `DateTime` Calculated date.
- - - -


{{< header href="date_year" >}}

date_year

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/year.go#L13)

date_year returns the year extracted from the given date.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`date` | `DateTime`  |  | Source datetime.


**Returns** `Int` A year number.
- - - -


{{< header href="now" >}}

now

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/now.go#L12)

now returns new DateTime object with Time equal to time.Now().

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description


**Returns** `DateTime` New datetime object.
- - - -

## Next steps

{{< docs-related tiles="stdlib,language-types-basic,language-operators-comparison" >}}
