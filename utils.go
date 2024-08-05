package money

import "github.com/shopspring/decimal"

func Sum(moneys ...*Money) (sum *Money, err error) {
	for _, money := range moneys {
		if sum == nil {
			sum = money
			continue
		}
		sum, err = sum.Add(money)
		if err != nil {
			return
		}
	}
	return
}

func Average(moneys ...*Money) (average *Money, err error) {
	sum, err := Sum(moneys...)
	if err != nil {
		return
	}
	average = sum.Divide(decimal.NewFromInt(int64(len(moneys))))
	return
}

func Min(moneys ...*Money) (min *Money, err error) {
	for _, money := range moneys {
		if min == nil {
			min = money
			continue
		}
		lt, err := money.LessThan(min)
		if err != nil {
			return nil, err
		}
		if lt {
			min = money
		}
	}
	return
}

func Max(moneys ...*Money) (max *Money, err error) {
	for _, money := range moneys {
		if max == nil {
			max = money
			continue
		}
		gt, err := money.GreaterThan(max)
		if err != nil {
			return nil, err
		}
		if gt {
			max = money
		}
	}
	return
}

func Sort(moneys []*Money) (sorted []*Money, err error) {
	sorted = make([]*Money, len(moneys))
	copy(sorted, moneys)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			lt, err := sorted[j].LessThan(sorted[i])
			if err != nil {
				return nil, err
			}
			if lt {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	return
}

func Median(moneys ...*Money) (mean *Money, err error) {
	sorted, err := Sort(moneys)
	if err != nil {
		return
	}
	if len(sorted)%2 == 0 {
		mean, err = Sum(sorted[len(sorted)/2-1], sorted[len(sorted)/2])
		if err != nil {
			return
		}
		mean = mean.Divide(decimal.NewFromInt(2))
	} else {
		mean = sorted[len(sorted)/2]
	}
	return
}
