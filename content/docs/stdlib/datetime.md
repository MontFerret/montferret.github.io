---
title: "datetime"
weight: 1
draft: false
---


## DATE_DAYOFWEEK
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/dayofweek.go#L13)

DateDayOfWeek returns number of the weekday from the date. Sunday is the 0th day of week.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`date` | `DateTime` | Source datetime.


**Returns** `Int` Return number of the weekday.
- - - -

## DATE_YEAR
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/year.go#L13)

DateYear returns the year extracted from the given date.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`date` | `DateTime` | Source datetime.


**Returns** `Int` A year number.
- - - -

## DATE_MILLISECOND
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/millisecond.go#L13)

DateMillisecond returns the millisecond of date as a number.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`date` | `DateTime` | Source datetime.


**Returns** `Int` A millisecond number.
- - - -

## DATE_DIFF
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/diff.go#L16)

DateDiff returns the difference between two dates in given time unit.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`date1` | `DateTime` | First datetime.
`date2` | `DateTime` | Second datetime.
`unit` | `String` | Time unit to return the difference in.
`asFloat` | `Boolean, optional` | If true amount of unit will be as float.


**Returns** `Int, Float` Difference between date1 and date2.
- - - -

## DATE_HOUR
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/hour.go#L13)

DateHour returns the hour of date as a number.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`date` | `DateTime` | Source datetime.


**Returns** `Int` A hour number.
- - - -

## DATE_SECOND
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/second.go#L13)

DateSecond returns the second of date as a number.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`date` | `DateTime` | Source datetime.


**Returns** `Int` A second number.
- - - -

## DateCompare
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/compare.go#L17)

DateCompare check if two partial dates match.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`date1, date2` | `DateTime` | Comparable dates.
`unitRangeStart` | `String` | Unit to start from.
`unitRangeEnd` | `String, Optional` | Unit to end with. error will be returned if unitrangestart unit less that unitrangeend.


**Returns** `Boolean` True if the dates match, else false.
- - - -

## DATE_ADD
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/add_subtract.go#L30)

DateAdd add amount given in unit to date.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`date` | `DateTime` | Source date.
`amount` | `Int` | Amount of units
`unit` | `String` | Unit.


**Returns** `DateTime` Calculated date. the following units are available: * y, year, year * m, month, months * w, week, weeks * d, day, days * h, hour, hours * i, minute, minutes * s, second, seconds * f, millisecond, milliseconds
- - - -

## DATE_SUBTRACT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/add_subtract.go#L60)

DateSubtract subtract amount given in unit to date.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`date` | `DateTime` | Source date.
`amount` | `Int` | Amount of units
`unit` | `String` | Unit.


**Returns** `DateTime` Calculated date. the following units are available: * y, year, year * m, month, months * w, week, weeks * d, day, days * h, hour, hours * i, minute, minutes * s, second, seconds * f, millisecond, milliseconds
- - - -

## DATE_FORMAT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/format.go#L13)

DateFormat format date according to the given format string.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`date` | `DateTime` | Source datetime object.


**Returns** `String` Formatted date.
- - - -

## DATE_MONTH
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/month.go#L13)

DateMonth returns the month of date as a number.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`date` | `DateTime` | Source datetime.


**Returns** `Int` A month number.
- - - -

## DATE_MINUTE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/minute.go#L13)

DateMinute returns the minute of date as a number.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`date` | `DateTime` | Source datetime.


**Returns** `Int` A minute number.
- - - -

## DATE_DAYS_IN_MONTH
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/daysinmonth.go#L29)

DateDaysInMonth returns the number of days in the month of date.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`date` | `DateTime` | Source datetime.


**Returns** `Int` Number of the days.
- - - -

## DATE_DAY
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/day.go#L13)

DateDay returns the day of date as a number.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`date` | `DateTime` | Source datetime.


**Returns** `Int` A day number.
- - - -

## DATE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/date.go#L14)

Date convert RFC3339 date time string to DateTime object.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`timeString` | `String` | String in rfc3339 format.


**Returns** `DateTime` New datetime object derived from timestring.
- - - -

## DATE_LEAPYEAR
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/leapyear.go#L13)

DateLeapYear returns true if date is in a leap year else false.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`date` | `DateTime` | Source datetime.


**Returns** `Boolean` Date is in a leap year.
- - - -

## DATE_QUARTER
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/quarter.go#L14)

DateQuarter returns which quarter date belongs to.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`date` | `DateTime` | Source datetime.


**Returns** `Int` A quarter number.
- - - -

## DATE_DAYOFYEAR
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/dayofyear.go#L14)

DateDayOfYear returns the day of year number of date. The return value range from 1 to 365 (366 in a leap year).

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`date` | `DateTime` | Source datetime.


**Returns** `Int` A day of year number.
- - - -

## NOW
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/datetime/now.go#L12)

Now returns new DateTime object with Time equal to time.Now().

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description


**Returns** `DateTime` New datetime object.
- - - -
