package money

import (
	"strings"

	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

// Formatter stores Money formatting information.
type Formatter struct {
	Fraction int32
	Decimal  string
	Thousand string
	Grapheme string
	Template string
}

// NewFormatter creates new Formatter instance.
func NewFormatter(fraction int32, decimal, thousand, grapheme, template string) *Formatter {
	return &Formatter{
		Fraction: fraction,
		Decimal:  decimal,
		Thousand: thousand,
		Grapheme: grapheme,
		Template: template,
	}
}

// Format returns string of formatted integer using given currency template.
func (f *Formatter) Format(amount decimal.Decimal) string {
	// Work with absolute amount value
	sa := cast.ToString(amount.IntPart())
	sa = strings.TrimPrefix(sa, "-")

	if f.Thousand != "" {
		for i := len(sa) - 3; i > 0; i -= 3 {
			sa = sa[:i] + f.Thousand + sa[i:]
		}
	}

	if f.Fraction > 0 {
		dg := cast.ToString(amount.Sub(decimal.NewFromInt(amount.IntPart())).InexactFloat64())
		dg = dg[strings.Index(dg, ".")+1:]

		if len(dg) > int(f.Fraction) {
			dg = dg[:f.Fraction]
		} else {
			for i := len(dg); i < int(f.Fraction); i++ {
				dg = dg + "0"
			}
		}

		sa = sa + f.Decimal + dg
	}
	sa = strings.Replace(f.Template, "1", sa, 1)
	sa = strings.Replace(sa, "$", f.Grapheme, 1)

	// Add minus sign for negative amount.
	if amount.IsNegative() {
		sa = "-" + sa
	}

	return sa
}

// abs return absolute value of given integer.
func (f Formatter) abs(amount int64) int64 {
	if amount < 0 {
		return -amount
	}

	return amount
}
