package money

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

// Injection points for backward compatibility.
// If you need to keep your JSON marshal/unmarshal way, overwrite them like below.
//
//	money.UnmarshalJSON = func (m *Money, b []byte) error { ... }
//	money.MarshalJSON = func (m Money) ([]byte, error) { ... }
var (
	// UnmarshalJSON is injection point of json.Unmarshaller for money.Money
	UnmarshalJSON = defaultUnmarshalJSON
	// MarshalJSON is injection point of json.Marshaller for money.Money
	MarshalJSON = defaultMarshalJSON

	// ErrCurrencyMismatch happens when two compared Money don't have the same currency.
	ErrCurrencyMismatch = errors.New("currencies don't match")

	// ErrInvalidJSONUnmarshal happens when the default money.UnmarshalJSON fails to unmarshal Money because of invalid data.
	ErrInvalidJSONUnmarshal = errors.New("invalid json unmarshal")

	one = decimal.NewFromInt(1)
)

var Zero = New(0, "")

func defaultUnmarshalJSON(m *Money, b []byte) error {
	data := make(map[string]interface{})
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	var amount float64
	if amountRaw, ok := data["amount"]; ok {
		amount, ok = amountRaw.(float64)
		if !ok {
			return ErrInvalidJSONUnmarshal
		}
	}

	var currency string
	if currencyRaw, ok := data["currency"]; ok {
		currency, ok = currencyRaw.(string)
		if !ok {
			return ErrInvalidJSONUnmarshal
		}
	}

	var ref *Money
	if amount == 0 && currency == "" {
		ref = &Money{}
	} else {
		ref = New(decimal.NewFromFloat(amount), currency)
	}

	*m = *ref
	return nil
}

func defaultMarshalJSON(m Money) ([]byte, error) {
	if m == (Money{}) {
		m = *New(0, "")
	}

	buff := bytes.NewBufferString(fmt.Sprintf(`{"amount": %.`+cast.ToString(m.currency.Fraction)+`f, "currency": "%s"}`, m.Amount(), m.Currency().Code))
	return buff.Bytes(), nil
}

// Amount is a data structure that stores the amount being used for calculations.
type Amount = decimal.Decimal

// Money represents monetary value information, stores
// currency and amount value.
type Money struct {
	currency *Currency `db:"currency"`
	amount   Amount    `db:"amount"`
}

// New creates and returns new instance of Money.
func New(amount any, code string) *Money {
	currency := newCurrency(code).get()
	return &Money{
		amount:   ConvertToDecimal(amount),
		currency: currency,
	}
}

// NewFromFloat creates and returns new instance of Money from a float64.
// Always rounding trailing decimals down.
func NewFromFloat(_amount float64, code string) *Money {
	amount := decimal.NewFromFloat(_amount)
	return New(amount, code)
}

// Currency returns the currency used by Money.
func (m *Money) Currency() *Currency {
	return m.currency
}

// Amount returns a copy of the internal monetary value as an int64.
func (m *Money) Amount() float64 {
	val, _ := m.amount.Truncate(m.currency.Fraction).Float64()
	return val
}

// SameCurrency check if given Money is equals by currency.
func (m *Money) SameCurrency(om *Money) bool {
	return m.currency.equals(om.currency)
}

func (m *Money) assertSameCurrency(om *Money) error {
	if !m.SameCurrency(om) {
		return ErrCurrencyMismatch
	}

	return nil
}

func (m *Money) compare(om *Money) int {
	return m.amount.Compare(om.amount)

}

// Equals checks equality between two Money types.
func (m *Money) Equals(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == 0, nil
}

// GreaterThan checks whether the value of Money is greater than the other.
func (m *Money) GreaterThan(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == 1, nil
}

// GreaterThanOrEqual checks whether the value of Money is greater or equal than the other.
func (m *Money) GreaterThanOrEqual(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) >= 0, nil
}

// LessThan checks whether the value of Money is less than the other.
func (m *Money) LessThan(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == -1, nil
}

// LessThanOrEqual checks whether the value of Money is less or equal than the other.
func (m *Money) LessThanOrEqual(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) <= 0, nil
}

// IsZero returns boolean of whether the value of Money is equals to zero.
func (m *Money) IsZero() bool {
	return m.amount.IsZero()
}

// IsPositive returns boolean of whether the value of Money is positive.
func (m *Money) IsPositive() bool {
	return m.amount.GreaterThan(decimal.NewFromInt(0))
}

// IsNegative returns boolean of whether the value of Money is negative.
func (m *Money) IsNegative() bool {
	return m.amount.LessThan(decimal.NewFromInt(0))
}

// Absolute returns new Money struct from given Money using absolute monetary value.
func (m *Money) Absolute() *Money {
	return &Money{amount: mutate.calc.absolute(m.amount), currency: m.currency}
}

// Negative returns new Money struct from given Money using negative monetary value.
func (m *Money) Negative() *Money {
	return &Money{amount: mutate.calc.negative(m.amount), currency: m.currency}
}

func (m *Money) ToDecimal() decimal.Decimal {
	return m.amount
}

// Add returns new Money struct with value representing sum of Self and Other Money.
func (m *Money) Add(ms ...*Money) (*Money, error) {
	if len(ms) == 0 {
		return m, nil
	}
	var k *Money
	if m.currency != nil && m.currency.Code != "" {
		k = New(0, m.currency.Code)
	} else {
		fnzv, found := lo.Find(ms, func(m *Money) bool {
			if m.currency != nil && m.currency.Code != "" {
				return true
			}
			return false
		})
		if !found {
			return nil, errors.New("no currency found")
		}
		k = New(0, fnzv.currency.Code)
	}

	for _, m2 := range ms {
		if m2.IsZero() {
			continue
		}
		if err := m.assertSameCurrency(m2); err != nil {
			return nil, err
		}

		k.amount = mutate.calc.add(k.amount, m2.amount)
	}

	return &Money{amount: mutate.calc.add(m.amount, k.amount), currency: m.currency}, nil
}

// Subtract returns new Money struct with value representing difference of Self and Other Money.
func (m *Money) Subtract(ms ...*Money) (*Money, error) {
	if len(ms) == 0 {
		return m, nil
	}

	var k *Money
	if m.currency != nil && m.currency.Code != "" {
		k = New(0, m.currency.Code)
	} else {
		fnzv, ok := lo.Find(ms, func(m *Money) bool {
			if m.currency != nil && m.currency.Code != "" {
				return true
			}
			return false
		})
		if !ok {
			return nil, errors.New("no currency found")
		}
		k = New(0, fnzv.currency.Code)
	}

	for _, m2 := range ms {
		if m2.IsZero() {
			continue
		}
		if err := m.assertSameCurrency(m2); err != nil {
			return nil, err
		}

		k.amount = mutate.calc.add(k.amount, m2.amount)
	}

	return &Money{amount: mutate.calc.subtract(m.amount, k.amount), currency: m.currency}, nil
}

// Multiply returns new Money struct with value representing Self multiplied value by multiplier.
func (m *Money) Multiply(muls ...any) *Money {
	if len(muls) == 0 {
		panic("At least one multiplier is required to multiply")
	}

	k := New(1, m.currency.Code)

	for _, m2 := range muls {

		k.amount = mutate.calc.multiply(k.amount, ConvertToDecimal(m2))
	}

	return &Money{amount: mutate.calc.multiply(m.amount, k.amount), currency: m.currency}
}

// Round returns new Money struct with value rounded to nearest zero.
func (m *Money) Round() *Money {
	return &Money{amount: mutate.calc.round(m.amount, 0), currency: m.currency}
}
func (m *Money) Divide(amount any) *Money {
	return &Money{amount: mutate.calc.divide(m.amount, ConvertToDecimal(amount), m.currency.Fraction), currency: m.currency}
}

func (m *Money) setCurrency(code string) *Money {
	m.currency = newCurrency(code).get()
	return m
}

// Split returns slice of Money structs with split Self value in given number.
// After division leftover pennies will be distributed round-robin amongst the parties.
// This means that parties listed first will likely receive more pennies than ones that are listed later.
func (m *Money) Split(n int) ([]*Money, error) {
	if n <= 0 {
		return nil, errors.New("split must be higher than zero")
	}

	a := mutate.calc.divide(m.amount, decimal.NewFromInt(int64(n)), m.currency.Fraction)
	ms := make([]*Money, n)

	for i := 0; i < n; i++ {
		ms[i] = &Money{amount: a, currency: m.currency}
	}

	r := mutate.calc.modulus(m.amount, decimal.NewFromInt(int64(n)), m.currency.Fraction)
	// l := mutate.calc.absolute(r)
	// Add leftovers to the first parties.
	v := one.Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt32(m.currency.Fraction)))
	// v := one
	if m.amount.IsNegative() {
		v = v.Neg()
	}
	for p := 0; !r.IsZero(); p++ {
		ms[p].amount = mutate.calc.add(ms[p].amount, v)
		r = r.Sub(v)
		if p == n-1 {
			p = 0
		}
	}

	return ms, nil
}

// Allocate returns slice of Money structs with split Self value in given ratios.
// It lets split money by given ratios without losing pennies and as Split operations distributes
// leftover pennies amongst the parties with round-robin principle.
// func (m *Money) Allocate(rs ...int) ([]*Money, error) {
// 	if len(rs) == 0 {
// 		return nil, errors.New("no ratios specified")
// 	}

// 	// Calculate sum of ratios.
// 	var sum decimal.Decimal
// 	for _, r := range rs {
// 		if r < 0 {
// 			return nil, errors.New("negative ratios not allowed")
// 		}
// 		if int64(r) > (math.MaxInt64 - sum) {
// 			return nil, errors.New("sum of given ratios exceeds max int")
// 		}
// 		sum += int64(r)
// 	}

// 	var total int64
// 	ms := make([]*Money, 0, len(rs))
// 	for _, r := range rs {
// 		party := &Money{
// 			amount:   mutate.calc.allocate(m.amount, r, sum, m.currency.Fraction),
// 			currency: m.currency,
// 		}

// 		ms = append(ms, party)
// 		total += party.amount
// 	}

// 	// if the sum of all ratios is zero, then we just returns zeros and don't do anything
// 	// with the leftover
// 	if sum == 0 {
// 		return ms, nil
// 	}

// 	// Calculate leftover value and divide to first parties.
// 	lo := m.amount - total
// 	sub := int64(1)
// 	if lo < 0 {
// 		sub = -sub
// 	}

// 	for p := 0; lo != 0; p++ {
// 		ms[p].amount = mutate.calc.add(ms[p].amount, sub)
// 		lo -= sub
// 	}

// 	return ms, nil
// }

// Display lets represent Money struct as string in given Currency value.
func (m *Money) Display() string {
	c := m.currency.get()
	return c.Formatter().Format(m.amount)
}

// Similar to Display but without the currency symbol
func (m *Money) Simple() string {
	c := *m.currency.get()
	c.Formatter().Grapheme = ""
	return c.Formatter().Format(m.amount)
}

// AsMajorUnits lets represent Money struct as subunits (float64) in given Currency value
func (m *Money) AsMajorUnits() float64 {
	c := m.currency.get()
	if c.Fraction == 0 {
		return float64(m.amount.Round(0).IntPart())
	}
	return m.Amount()
}

// UnmarshalJSON is implementation of json.Unmarshaller
func (m *Money) UnmarshalJSON(b []byte) error {
	return UnmarshalJSON(m, b)
}

// MarshalJSON is implementation of json.Marshaller
func (m Money) MarshalJSON() ([]byte, error) {
	return MarshalJSON(m)
}

// Compare function compares two money of the same type
//
//	if m.amount > om.amount returns (1, nil)
//	if m.amount == om.amount returns (0, nil
//	if m.amount < om.amount returns (-1, nil)
//
// If compare moneys from distinct currency, return (m.amount, ErrCurrencyMismatch)
func (m *Money) Compare(om *Money) (int, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return m.compare(om), err
	}

	return m.compare(om), nil
}

func ConvertToDecimal[T any](value T) decimal.Decimal {
	switch v := any(value).(type) {
	case int:
		return decimal.NewFromInt(int64(v))
	case int8:
		return decimal.NewFromInt(int64(v))
	case int16:
		return decimal.NewFromInt(int64(v))
	case int32:
		return decimal.NewFromInt(int64(v))
	case int64:
		return decimal.NewFromInt(v)
	case uint:
		return decimal.NewFromInt(int64(v))
	case uint8:
		return decimal.NewFromInt(int64(v))
	case uint16:
		return decimal.NewFromInt(int64(v))
	case uint32:
		return decimal.NewFromInt(int64(v))
	case uint64:
		return decimal.NewFromInt(int64(v))
	case float32:
		return decimal.NewFromFloat(float64(v))
	case float64:
		return decimal.NewFromFloat(float64(v))
	case string:
		dec, err := decimal.NewFromString(v)
		if err != nil {
			panic(err)
		}
		return dec
	case decimal.Decimal:
		return v
	case *decimal.Decimal:
		return *v

	default:
		panic("unsupported type")
	}
}
