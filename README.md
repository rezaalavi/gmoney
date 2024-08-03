# Money


![alt text](http://i.imgur.com/c3XmCC6.jpg "Money")

This is a fork of [go-money](github.com/rhymond/go-money) that uses (github.com/shopspring/decimal) instead of int64 to store the amount. This allows for more functionality when working with money, including but not limited to multiply and divide operations by float values.


[![Go Report Card](https://goreportcard.com/badge/github.com/rezaalavi/gmoney)](https://goreportcard.com/report/github.com/rezaalavi/gmoney)
[![Coverage Status](https://coveralls.io/repos/github/rezaalavi/gmoney/badge.svg?branch=master)](https://coveralls.io/github/rezaalavi/gmoney?branch=master)
[![GoDoc](https://godoc.org/github.com/rezaalavi/gmoney?status.svg)](https://godoc.org/github.com/rezaalavi/gmoney)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**GoMoney** provides ability to work with [monetary value using a currency's smallest unit](https://martinfowler.com/eaaCatalog/money.html).
This package provides basic and precise Money operations such as rounding, splitting and allocating.  Monetary values should not be stored as floats due to small rounding differences.

```go
package main

import (
  "log"

  "github.com/rezaalavi/gmoney"
)

func main() {
    pound := money.New(1.00, money.GBP)
    twoPounds, err := pound.Add(pound)

    if err != nil {
        log.Fatal(err)
    }

    parties, err := twoPounds.Split(3)

    if err != nil {
        log.Fatal(err)
    }

    parties[0].Display() // £0.67
    parties[1].Display() // £0.67
    parties[2].Display() // £0.66
}

```
Quick start
-
Get the package:

``` bash
$ go get github.com/rezaalavi/gmoney
```

## Features
* Provides a Money struct which stores information about an Money amount value and its currency.
* Provides a ```Money.Amount``` struct which encapsulates all information about a monetary unit.
* Represents monetary values as decimals. This avoids floating point rounding errors.
* Represents currency as ```Money.Currency``` instances providing a high level of flexibility.

Usage
-
### Initialization
Initialize Money by using smallest unit value (e.g 100 represents 1 pound). Use ISO 4217 Currency Code to set money Currency. Note that constants are also provided for all ISO 4217 currency codes.
```go
pound := money.New(1.00, money.GBP)
```
Or initialize Money using the any other numerical values.

Comparison
-
**Gmoney** provides base compare operations like:

* Equals
* GreaterThan
* GreaterThanOrEqual
* LessThan
* LessThanOrEqual
* Compare

Comparisons must be made between the same currency units.

```go
pound := money.New(1.00, money.GBP)
twoPounds := money.New(2.00, money.GBP)
twoEuros := money.New(2.00, money.EUR)

pound.GreaterThan(twoPounds) // false, nil
pound.LessThan(twoPounds) // true, nil
twoPounds.Equals(twoEuros) // false, error: Currencies don't match
twoPounds.Compare(pound) // 1, nil
pound.Compare(twoPounds) // -1, nil
pound.Compare(pound) // 0, nil
pound.Compare(twoEuros) // pound.amount, ErrCurrencyMismatch
```
Asserts
-
* IsZero
* IsNegative
* IsPositive

#### Zero value

To assert if Money value is equal to zero use `IsZero()`

```go
pound := money.New(1.00, money.GBP)
result := pound.IsZero() // false
```

#### Positive value

To assert if Money value is more than zero use `IsPositive()`

```go
pound := money.New(1.00, money.GBP)
pound.IsPositive() // true
```

#### Negative value

To assert if Money value is less than zero use `IsNegative()`

```go
pound := money.New(1.00, money.GBP)
pound.IsNegative() // false
```

Operations
-
* Add
* Subtract
* Multiply
* Absolute
* Negative

Comparisons must be made between the same currency units.

#### Addition

Additions can be performed using `Add()`.

```go
pound := money.New(1.00, money.GBP)
twoPounds := money.New(2.00, money.GBP)

result, err := pound.Add(twoPounds) // £3.00, nil
```

#### Subtraction

Subtraction can be performed using `Subtract()`.

```go
pound := money.New(1.00, money.GBP)
twoPounds := money.New(2.00, money.GBP)

result, err := pound.Subtract(twoPounds) // -£1.00, nil
```

#### Multiplication

Multiplication can be performed using `Multiply()`.

```go
pound := money.New(1.00, money.GBP)

result := pound.Multiply(decimal.NewFromFloat(2.5)) // £2.50
```

#### Absolute

Return `absolute` value of Money structure

```go
pound := money.New(-1.00, money.GBP)

result := pound.Absolute() // £1.00
```

#### Negative

Return `negative` value of Money structure

```go
pound := money.New(1.00, money.GBP)

result := pound.Negative() // -£1.00
```

Allocation
-

* Split
* Allocate

#### Splitting

In order to split Money for parties without losing any pennies due to rounding differences, use `Split()`.

After division leftover pennies will be distributed round-robin amongst the parties. This means that parties listed first will likely receive more pennies than ones that are listed later.

```go
pound := money.New(100, money.GBP)
parties, err := pound.Split(3)

if err != nil {
    log.Fatal(err)
}

parties[0].Display() // £0.34
parties[1].Display() // £0.33
parties[2].Display() // £0.33
```

Format
-

To format and return Money as a string use `Display()`.

```go
money.New(1234567.89, money.EUR).Display() // €1,234,567.89
```
To format and return Money as a float64 representing the amount value in the currency's subunit use `AsMajorUnits()`.

```go
money.New(1234567.89, money.EUR).AsMajorUnits() // 1234567.89
```

Contributing
-
Thank you for considering contributing!
Please use GitHub issues and Pull Requests for contributing.

License
-
The MIT License (MIT). Please see License File for more information.



[![forthebadge](http://forthebadge.com/images/badges/built-with-love.svg)](https://github.com/rezaalavi/gmoney)
