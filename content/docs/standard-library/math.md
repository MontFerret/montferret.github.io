---
title: "Math"
weight: 50
draft: false
description: "Math functions in the Ferret standard library."
aliases:
  - /docs/stdlib/math/
menuTitle: 
menu: [abs,acos,asin,atan,atan2,average,ceil,cos,degrees,exp,exp2,floor,log,log10,log2,max,median,min,percentile,pi,pow,radians,rand,range,round,sin,sqrt,stddev_population,stddev_sample,sum,tan,variance_population,variance_sample,]
---



{{< header href="abs" >}}

abs

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/abs.go#L14)

abs returns the absolute value of a given number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Float` The absolute value of a given number.
- - - -


{{< header href="acos" >}}

acos

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/acos.go#L14)

acos returns the arccosine, in radians, of a given number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Float` The arccosine, in radians, of a given number.
- - - -


{{< header href="asin" >}}

asin

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/asin.go#L14)

asin returns the arcsine, in radians, of a given number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Float` The arcsine, in radians, of a given number.
- - - -


{{< header href="atan" >}}

atan

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/atan.go#L14)

atan returns the arctangent, in radians, of a given number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Float` The arctangent, in radians, of a given number.
- - - -


{{< header href="atan2" >}}

atan2

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/atan2.go#L15)

atan2 returns the arc tangent of y/x, using the signs of the two to determine the quadrant of the return value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number1` | `Int` `Float`  |  | Input number.
`number2` | `Int` `Float`  |  | Input number.


**Returns** `Float` The arc tangent of y/x, using the signs of the two to determine the quadrant of the return value.
- - - -


{{< header href="average" >}}

average

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/average.go#L13)

average Returns the average (arithmetic mean) of the values in array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Int[]` `Float[]`  |  | Array of numbers.


**Returns** `Float` The average of the values in array.
- - - -


{{< header href="ceil" >}}

ceil

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/ceil.go#L14)

ceil returns the least integer value greater than or equal to a given value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Int` The least integer value greater than or equal to a given value.
- - - -


{{< header href="cos" >}}

cos

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/cos.go#L14)

cos returns the cosine of a given number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Float` The cosine of a given number.
- - - -


{{< header href="degrees" >}}

degrees

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/degrees.go#L13)

degrees returns the angle converted from radians to degrees.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | The input number.


**Returns** `Float` The angle in degrees
- - - -


{{< header href="exp" >}}

exp

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/exp.go#L14)

exp returns Euler's constant (2.71828...) raised to the power of value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Float` Euler's constant raised to the power of value.
- - - -


{{< header href="exp2" >}}

exp2

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/exp2.go#L14)

exp2 returns 2 raised to the power of value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Float` 2 raised to the power of value.
- - - -


{{< header href="floor" >}}

floor

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/floor.go#L14)

floor returns the greatest integer value less than or equal to a given value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Int` The greatest integer value less than or equal to a given value.
- - - -


{{< header href="log" >}}

log

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/log.go#L14)

log returns the natural logarithm of a given value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Float` The natural logarithm of a given value.
- - - -


{{< header href="log10" >}}

log10

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/log10.go#L14)

log10 returns the decimal logarithm of a given value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Float` The decimal logarithm of a given value.
- - - -


{{< header href="log2" >}}

log2

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/log2.go#L14)

log2 returns the binary logarithm of a given value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Float` The binary logarithm of a given value.
- - - -


{{< header href="max" >}}

max

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/max.go#L13)

max returns the greatest (arithmetic mean) of the values in array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Int[]` `Float[]`  |  | Array of numbers.


**Returns** `Float` The greatest of the values in array.
- - - -


{{< header href="median" >}}

median

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/median.go#L14)

median returns the median of the values in array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Int[]` `Float[]`  |  | Array of numbers.


**Returns** `Float` The median of the values in array.
- - - -


{{< header href="min" >}}

min

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/min.go#L13)

min returns the smallest (arithmetic mean) of the values in array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Int[]` `Float[]`  |  | Array of numbers.


**Returns** `Float` The smallest of the values in array.
- - - -


{{< header href="percentile" >}}

percentile

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/percentile.go#L17)

percentile returns the nth percentile of the values in a given array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Int[]` `Float[]`  |  | Array of numbers.
`number` | `Int`  |  | A number which must be between 0 (excluded) and 100 (included).
`method` | `String`  | `"rank"` | "rank" or "interpolation".


**Returns** `Float` The nth percentile, or null if the array is empty or only null values are contained in it or the percentile cannot be calculated.
- - - -


{{< header href="pi" >}}

pi

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/pi.go#L12)

pi returns Pi value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description


**Returns** `Float` Pi value.
- - - -


{{< header href="pow" >}}

pow

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/pow.go#L15)

pow returns the base to the exponent value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`base` | `Int` `Float`  |  | The base value.
`exp` | `Int` `Float`  |  | The exponent value.


**Returns** `Float` The exponentiated value.
- - - -


{{< header href="radians" >}}

radians

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/radians.go#L13)

radians returns the angle converted from degrees to radians.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | The input number.


**Returns** `Float` The angle in radians.
- - - -


{{< header href="rand" >}}

rand

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/rand.go#L13)

rand return a pseudo-random number between 0 and 1.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`max` | `Int` `Float`  |  | Upper limit.
`min` | `Int` `Float`  |  | Lower limit.


**Returns** `Float` A number greater than 0 and less than 1.
- - - -


{{< header href="range" >}}

range

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/range.go#L15)

range returns an array of numbers in the specified range, optionally with increments other than 1.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`start` | `Int` `Float`  |  | The value to start the range at (inclusive).
`end` | `Int` `Float`  |  | The value to end the range with (inclusive).
`step` | `Int` `Float`  | `1.0` | How much to increment in every step.


**Returns** `Int[]` `Float[]` Array of numbers in the specified range, optionally with increments other than 1.
- - - -


{{< header href="round" >}}

round

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/round.go#L14)

round returns the nearest integer, rounding half away from zero.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Int` The nearest integer, rounding half away from zero.
- - - -


{{< header href="sin" >}}

sin

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/sin.go#L14)

sin returns the sine of the radian argument.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Float` The sin, in radians, of a given number.
- - - -


{{< header href="sqrt" >}}

sqrt

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/sqrt.go#L14)

sqrt returns the square root of a given number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Int` `Float`  |  | A number.


**Returns** `Float` The square root.
- - - -


{{< header href="stddev_population" >}}

stddev_population

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/stddev_population.go#L14)

stddev_population returns the population standard deviation of the values in a given array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`numbers` | `Int[]` `Float[]`  |  | Array of numbers.


**Returns** `Float` The population standard deviation.
- - - -


{{< header href="stddev_sample" >}}

stddev_sample

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/stddev_sample.go#L14)

stddev_sample returns the sample standard deviation of the values in a given array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`numbers` | `Int[]` `Float[]`  |  | Array of numbers.


**Returns** `Float` The sample standard deviation.
- - - -


{{< header href="sum" >}}

sum

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/sum.go#L13)

sum returns the sum of the values in a given array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`numbers` | `Int[]` `Float[]`  |  | Array of numbers.


**Returns** `Float` The sum of the values.
- - - -


{{< header href="tan" >}}

tan

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/tan.go#L14)

tan returns the tangent of a given number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | A number.


**Returns** `Float` The tangent.
- - - -


{{< header href="variance_population" >}}

variance_population

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/variance_population.go#L14)

variance_population returns the population variance of the values in a given array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`numbers` | `Int[]` `Float[]`  |  | Array of numbers.


**Returns** `Float` The population variance.
- - - -


{{< header href="variance_sample" >}}

variance_sample

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/variance_sample.go#L14)

variance_sample returns the sample variance of the values in a given array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`numbers` | `Int[]` `Float[]`  |  | Array of numbers.


**Returns** `Float` The sample variance.
- - - -

## Next steps

{{< docs-related tiles="stdlib,language-types-basic,language-operators-arithmetic" >}}
