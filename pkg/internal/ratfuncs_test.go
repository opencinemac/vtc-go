package internal_test

import (
	"fmt"
	"github.com/opencinemac/vtc-go/pkg/internal"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestDivModRat(t *testing.T) {
	cases := []struct {
		a         *big.Rat
		b         *big.Rat
		dividend  *big.Rat
		remainder *big.Rat
	}{
		{
			a:         big.NewRat(5, 1),
			b:         big.NewRat(2, 1),
			dividend:  big.NewRat(2, 1),
			remainder: big.NewRat(1, 1),
		},
		{
			a:         big.NewRat(24, 1),
			b:         big.NewRat(5, 1),
			dividend:  big.NewRat(4, 1),
			remainder: big.NewRat(4, 1),
		},
		{
			a:         big.NewRat(9, 2),
			b:         big.NewRat(2, 1),
			dividend:  big.NewRat(2, 1),
			remainder: big.NewRat(1, 2),
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v /%% %v = %v, %v", tc.a, tc.b, tc.dividend, tc.remainder), func(t *testing.T) {
			assert := assert.New(t)

			dividend, remainder := internal.DivModRat(tc.a, tc.b)

			assert.Equal(
				tc.dividend, dividend, "dividend expected: %v, got %v", tc.dividend, dividend,
			)
			assert.Equal(
				tc.remainder, remainder, "remainder expected: %v, got %v", tc.remainder, remainder,
			)
		})
	}
}
