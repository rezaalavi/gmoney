package money

import "github.com/shopspring/decimal"

type calculator struct{}

func (c *calculator) add(a, b Amount) Amount {
	return a.Add(b)
}

func (c *calculator) subtract(a, b Amount) Amount {
	return a.Sub(b)
}

func (c *calculator) multiply(a Amount, m Amount) Amount {
	return a.Mul(m)
}

func (c *calculator) divide(a Amount, d Amount, precision int32) Amount {
	return a.Div(d).Truncate(precision)
}

func (c *calculator) modulus(a Amount, b Amount, precision int32) Amount {
	_, rem := a.QuoRem(b, precision)
	return rem
}

func (c *calculator) allocate(a, r, s Amount, precision int32) Amount {
	if a.IsZero() || s.IsZero() {
		return decimal.NewFromInt(0)
	}

	return a.Mul(r).DivRound(s, precision)
}

func (c *calculator) absolute(a Amount) Amount {
	return a.Abs()
}

func (c *calculator) negative(a Amount) Amount {

	return a.Neg()
}

func (c *calculator) round(a Amount, e int32) Amount {

	return a.Round(e)
}
