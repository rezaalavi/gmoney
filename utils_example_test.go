package money_test

import (
	"fmt"
	"log"

	money "github.com/rezaalavi/gmoney"
	"github.com/shopspring/decimal"
)

func ExampleSum() {
	pound := money.New(decimal.NewFromInt(1), "GBP")
	twoPounds := money.New(decimal.NewFromInt(2), "GBP")
	threePounds := money.New(decimal.NewFromInt(3), "GBP")

	sum, err := money.Sum(pound, twoPounds, threePounds)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(sum.Display())

	// Output:
	// £6.00
}

func ExampleAverage() {
	pound := money.New(decimal.NewFromInt(1), "GBP")
	twoPounds := money.New(decimal.NewFromInt(2), "GBP")
	fourPounds := money.New(decimal.NewFromInt(4), "GBP")

	average, err := money.Average(pound, twoPounds, fourPounds)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(average.Display())

	// Output:
	// £2.33
}

func ExampleMin() {
	pound := money.New(decimal.NewFromInt(1), "GBP")
	twoPounds := money.New(decimal.NewFromInt(2), "GBP")
	fourPounds := money.New(decimal.NewFromInt(4), "GBP")

	min, err := money.Min(pound, twoPounds, fourPounds)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(min.Display())

	// Output:
	// £1.00
}

func ExampleMax() {
	pound := money.New(decimal.NewFromInt(1), "GBP")
	twoPounds := money.New(decimal.NewFromInt(2), "GBP")
	fourPounds := money.New(decimal.NewFromInt(4), "GBP")

	max, err := money.Max(pound, twoPounds, fourPounds)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(max.Display())

	// Output:
	// £4.00
}

func ExampleSort() {
	pound := money.New(decimal.NewFromInt(1), "GBP")
	twoPounds := money.New(decimal.NewFromInt(2), "GBP")
	fourPounds := money.New(decimal.NewFromInt(4), "GBP")

	sorted, err := money.Sort([]*money.Money{fourPounds, pound, twoPounds})

	if err != nil {
		log.Fatal(err)
	}

	for _, m := range sorted {
		fmt.Println(m.Display())
	}

	// Output:
	// £1.00
	// £2.00
	// £4.00
}

func ExampleMedian() {
	pound := money.New(decimal.NewFromInt(1), "GBP")
	twoPounds := money.New(decimal.NewFromInt(2), "GBP")
	threePounds := money.New(decimal.NewFromInt(3), "GBP")
	fourPounds := money.New(decimal.NewFromInt(4), "GBP")

	median, err := money.Median(threePounds, pound, twoPounds, fourPounds)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(median.Display())

	// Output:
	// £2.50
}
