package money_test

import (
	"fmt"
	"log"

	money "github.com/rezaalavi/gmoney"
	"github.com/shopspring/decimal"
)

func ExampleMoney() {
	pound := money.New(decimal.NewFromInt(1), "GBP")
	twoPounds, err := pound.Add(pound)

	if err != nil {
		log.Fatal(err)
	}

	parties, err := twoPounds.Split(3)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(parties[0].Display())
	fmt.Println(parties[1].Display())
	fmt.Println(parties[2].Display())

	// Output:
	// £0.67
	// £0.67
	// £0.66
}

func ExampleNew() {
	pound := money.New(decimal.NewFromInt(1), "GBP")

	fmt.Println(pound.Display())

	// Output:
	// £1.00
}

func ExampleMoney_comparisons() {
	pound := money.New(decimal.NewFromInt(1), "GBP")
	twoPounds := money.New(decimal.NewFromInt(2), "GBP")
	twoEuros := money.New(decimal.NewFromInt(2), "EUR")

	gt, err := pound.GreaterThan(twoPounds)
	fmt.Println(gt, err)

	lt, err := pound.LessThan(twoPounds)
	fmt.Println(lt, err)

	eq, err := twoPounds.Equals(twoEuros)
	fmt.Println(eq, err)

	// Output:
	// false <nil>
	// true <nil>
	// false currencies don't match
}

func ExampleMoney_IsZero() {
	pound := money.New(decimal.NewFromInt(1), "GBP")
	fmt.Println(pound.IsZero())

	// Output:
	// false
}

func ExampleMoney_IsPositive() {
	pound := money.New(decimal.NewFromInt(1), "GBP")
	fmt.Println(pound.IsPositive())

	// Output:
	// true
}

func ExampleMoney_IsNegative() {
	pound := money.New(decimal.NewFromInt(1), "GBP")
	fmt.Println(pound.IsNegative())

	// Output:
	// false
}

func ExampleMoney_Add() {
	pound := money.New(decimal.NewFromInt(1), "GBP")
	twoPounds := money.New(decimal.NewFromInt(2), "GBP")

	result, err := pound.Add(twoPounds)
	fmt.Println(result.Display(), err)

	// Output:
	// £3.00 <nil>
}

func ExampleMoney_Subtract() {
	pound := money.New(decimal.NewFromInt(1), "GBP")
	twoPounds := money.New(decimal.NewFromInt(2), "GBP")

	result, err := pound.Subtract(twoPounds)
	fmt.Println(result.Display(), err)

	// Output:
	// -£1.00 <nil>
}

func ExampleMoney_Multiply() {
	pound := money.New(decimal.NewFromInt(1), "GBP")

	result := pound.Multiply(decimal.NewFromInt(2))
	fmt.Println(result.Display())

	// Output:
	// £2.00
}

func ExampleMoney_Absolute() {
	pound := money.New(decimal.NewFromInt(-1), "GBP")

	result := pound.Absolute()
	fmt.Println(result.Display())

	// Output:
	// £1.00
}

func ExampleMoney_Split() {
	pound := money.New(decimal.NewFromInt(1), "GBP")
	parties, err := pound.Split(3)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(parties[0].Display())
	fmt.Println(parties[1].Display())
	fmt.Println(parties[2].Display())

	// Output:
	// £0.34
	// £0.33
	// £0.33
}

// func ExampleMoney_Allocate() {
// 	pound := money.New(decimal.NewFromInt(1), "GBP")
// 	// Allocate is variadic function which can receive ratios as
// 	// slice (int[]{33, 33, 33}...) or separated by a comma integers
// 	parties, err := pound.Allocate(33, 33, 33)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println(parties[0].Display())
// 	fmt.Println(parties[1].Display())
// 	fmt.Println(parties[2].Display())

// 	// Output:
// 	// £0.34
// 	// £0.33
// 	// £0.33
// }

func ExampleMoney_Display() {
	fmt.Println(money.New(decimal.NewFromFloat(1234567.89), "EUR").Display())

	// Output:
	// €1,234,567.89
}

func ExampleMoney_AsMajorUnits() {
	result := money.New(decimal.NewFromFloat(1234567.89), "EUR").AsMajorUnits()
	fmt.Printf("%.2f", result)

	// Output:
	// 1234567.89
}
