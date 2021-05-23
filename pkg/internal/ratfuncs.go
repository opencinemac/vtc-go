package internal

import (
	"github.com/wadey/go-rounding"
	"math/big"
)

// RoundRat rounds a rational value to it's nearest integer value in-place.
func RoundRat(value *big.Rat) *big.Rat {
	// If this value is an integer already, we don't need to round it.
	return rounding.Round(value, 0, rounding.HalfUp)
}

// DivModRat returns the floor division and remainder of a rational value.
//
// The returned values are newly created.
//
// x is set to remainder and then returned. y is not modified.
func DivModRat(x *big.Rat, y *big.Rat) (dividend *big.Rat, remainder *big.Rat) {
	divisor := new(big.Rat).Inv(y)
	dividend = rounding.Round(new(big.Rat).Mul(x, divisor), 0, rounding.Down)

	dividendMultiplied := new(big.Rat).Mul(dividend, y)
	remainder = x.Sub(x, dividendMultiplied)
	return dividend, remainder
}
