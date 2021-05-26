package tc_test

import (
	"fmt"
	"github.com/opencinemac/vtc-go/pkg/rate"
	"github.com/opencinemac/vtc-go/pkg/tc"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestTimecode_Cmp(t *testing.T) {
	cases := []struct {
		Tc1      tc.Timecode
		Tc2      tc.Timecode
		Expected tc.Cmp
	}{
		// 24 FPS CASES ----------
		// -----------------------
		{
			Tc1:      mustTC("01:00:00:00", rate.F24),
			Tc2:      mustTC("01:00:00:00", rate.F24),
			Expected: tc.CmpEq,
		},
		{
			Tc1:      mustTC("01:00:00:00", rate.F24),
			Tc2:      mustTC("00:59:59:24", rate.F24),
			Expected: tc.CmpEq,
		},
		{
			Tc1:      mustTC("01:00:00:00", rate.F24),
			Tc2:      mustTC("02:00:00:00", rate.F24),
			Expected: tc.CmpLt,
		},
		{
			Tc1:      mustTC("01:00:00:00", rate.F24),
			Tc2:      mustTC("01:00:00:01", rate.F24),
			Expected: tc.CmpLt,
		},
		{
			Tc1:      mustTC("01:00:00:00", rate.F24),
			Tc2:      mustTC("00:59:59:23", rate.F24),
			Expected: tc.CmpGt,
		},
		// 23.98 NTSC CASES ------
		// -----------------------
		{
			Tc1:      mustTC("01:00:00:00", rate.F23_98),
			Tc2:      mustTC("01:00:00:01", rate.F23_98),
			Expected: tc.CmpLt,
		},
		{
			Tc1:      mustTC("00:00:00:00", rate.F23_98),
			Tc2:      mustTC("02:00:00:01", rate.F23_98),
			Expected: tc.CmpLt,
		},
		// MIXED FPS CASES ------
		// ----------------------
		{
			Tc1:      mustTC("01:00:00:00", rate.F23_98),
			Tc2:      mustTC("01:00:00:00", rate.F24),
			Expected: tc.CmpGt,
		},
	}

	for _, thisCase := range cases {
		name := fmt.Sprintf("%v %v %v", thisCase.Tc1, thisCase.Expected, thisCase.Tc2)
		t.Run(name, func(t *testing.T) {
			t.Run("Regular", func(t *testing.T) {
				assert.Equal(t, thisCase.Expected, thisCase.Tc1.Cmp(thisCase.Tc2))
			})

			t.Run("Flipped", func(t *testing.T) {
				// Invert our expected result.
				var newExpected tc.Cmp
				switch thisCase.Expected {
				case tc.CmpLt:
					newExpected = tc.CmpGt
				case tc.CmpEq:
					newExpected = tc.CmpEq
				case tc.CmpGt:
					newExpected = tc.CmpLt
				}

				assert.Equal(t, newExpected, thisCase.Tc2.Cmp(thisCase.Tc1))
			})
		})
	}
}

func mustTC(timecode string, framerate rate.Framerate) tc.Timecode {
	parsed, err := tc.FromTimecode(timecode, framerate)
	if err != nil {
		panic(fmt.Errorf("error parsing timecode '%v': %w", timecode, err))
	}

	return parsed
}

func TestTimecode_Add(t *testing.T) {
	cases := []struct {
		Tc1      tc.Timecode
		Tc2      tc.Timecode
		Expected tc.Timecode
	}{
		{
			Tc1:      mustTC("01:00:00:00", rate.F24),
			Tc2:      mustTC("01:00:00:00", rate.F24),
			Expected: mustTC("02:00:00:00", rate.F24),
		},
		{
			Tc1:      mustTC("01:00:00:00", rate.F24),
			Tc2:      mustTC("00:00:00:01", rate.F24),
			Expected: mustTC("01:00:00:01", rate.F24),
		},
		{
			Tc1:      mustTC("01:00:00:00", rate.F24),
			Tc2:      mustTC("-00:30:00:00", rate.F24),
			Expected: mustTC("00:30:00:00", rate.F24),
		},
	}

	for _, testCase := range cases {
		name := fmt.Sprintf("%v + %v = %v", testCase.Tc1, testCase.Tc2, testCase.Expected)
		t.Run(name, func(t *testing.T) {
			t.Run("Regular", func(t *testing.T) {
				result := testCase.Tc1.Add(testCase.Tc2)
				assert.Equal(t, tc.CmpEq, testCase.Expected.Cmp(result), "result equals expected")
			})

			t.Run("Flipped", func(t *testing.T) {
				result := testCase.Tc2.Add(testCase.Tc1)
				assert.Equal(t, tc.CmpEq, testCase.Expected.Cmp(result), "result equals expected")
			})
		})
	}
}

func TestTimecode_Sub(t *testing.T) {
	cases := []struct {
		Tc1      tc.Timecode
		Tc2      tc.Timecode
		Expected tc.Timecode
	}{
		{
			Tc1:      mustTC("01:00:00:00", rate.F24),
			Tc2:      mustTC("01:00:00:00", rate.F24),
			Expected: mustTC("00:00:00:00", rate.F24),
		},
		{
			Tc1:      mustTC("01:00:00:00", rate.F24),
			Tc2:      mustTC("00:00:00:01", rate.F24),
			Expected: mustTC("00:59:59:23", rate.F24),
		},
		{
			Tc1:      mustTC("01:00:00:00", rate.F24),
			Tc2:      mustTC("-00:30:00:00", rate.F24),
			Expected: mustTC("01:30:00:00", rate.F24),
		},
	}

	for _, testCase := range cases {
		name := fmt.Sprintf("%v - %v = %v", testCase.Tc1, testCase.Tc2, testCase.Expected)
		t.Run(name, func(t *testing.T) {
			result := testCase.Tc1.Sub(testCase.Tc2)
			assert.Equal(t, tc.CmpEq, testCase.Expected.Cmp(result), "result equals expected")
		})
	}
}

func TestTimecode_Mul(t *testing.T) {
	cases := []struct {
		Tc         tc.Timecode
		Multiplier *big.Rat
		Expected   tc.Timecode
	}{
		{
			Tc:         mustTC("01:00:00:00", rate.F24),
			Multiplier: big.NewRat(2, 1),
			Expected:   mustTC("02:00:00:00", rate.F24),
		},
		{
			Tc:         mustTC("01:00:00:00", rate.F24),
			Multiplier: big.NewRat(3, 2),
			Expected:   mustTC("01:30:00:00", rate.F24),
		},
		{
			Tc:         mustTC("01:00:00:00", rate.F24),
			Multiplier: big.NewRat(0, 1),
			Expected:   mustTC("00:00:00:00", rate.F24),
		},
		{
			Tc:         mustTC("00:00:00:00", rate.F24),
			Multiplier: big.NewRat(10, 1),
			Expected:   mustTC("00:00:00:00", rate.F24),
		},
	}

	for _, testCase := range cases {
		name := fmt.Sprintf("%v * %v = %v", testCase.Tc, testCase.Multiplier, testCase.Expected)
		t.Run(name, func(t *testing.T) {
			result := testCase.Tc.Mul(testCase.Multiplier)
			assert.Equal(t, tc.CmpEq, testCase.Expected.Cmp(result), "result equals expected")
		})
	}
}

func TestTimecode_DivMod(t *testing.T) {
	cases := []struct {
		Tc                tc.Timecode
		Divisor           *big.Rat
		ExpectedDividend  tc.Timecode
		ExpectedRemainder tc.Timecode
	}{
		{
			Tc:                mustTC("01:00:00:00", rate.F24),
			Divisor:           big.NewRat(2, 1),
			ExpectedDividend:  mustTC("00:30:00:00", rate.F24),
			ExpectedRemainder: mustTC("00:00:00:00", rate.F24),
		},
		{
			Tc:                mustTC("01:00:00:01", rate.F24),
			Divisor:           big.NewRat(2, 1),
			ExpectedDividend:  mustTC("00:30:00:00", rate.F24),
			ExpectedRemainder: mustTC("00:00:00:01", rate.F24),
		},
		{
			Tc:                mustTC("01:00:00:00", rate.F24),
			Divisor:           big.NewRat(4, 1),
			ExpectedDividend:  mustTC("00:15:00:00", rate.F24),
			ExpectedRemainder: mustTC("00:00:00:00", rate.F24),
		},
		{
			Tc:                mustTC("01:00:00:03", rate.F24),
			Divisor:           big.NewRat(4, 1),
			ExpectedDividend:  mustTC("00:15:00:00", rate.F24),
			ExpectedRemainder: mustTC("00:00:00:03", rate.F24),
		},
		{
			Tc:                mustTC("01:00:00:04", rate.F24),
			Divisor:           big.NewRat(3, 2),
			ExpectedDividend:  mustTC("00:40:00:02", rate.F24),
			ExpectedRemainder: mustTC("00:00:00:01", rate.F24),
		},
	}

	for _, testCase := range cases {
		name := fmt.Sprintf("%v /%% %v", testCase.Tc, testCase.Divisor)
		t.Run(name, func(t *testing.T) {
			t.Run("Div", func(t *testing.T) {
				result := testCase.Tc.Div(testCase.Divisor)
				assert.Equal(
					t, tc.CmpEq, testCase.ExpectedDividend.Cmp(result), "dividend equals expected",
				)
			})

			t.Run("Mod", func(t *testing.T) {
				result := testCase.Tc.Mod(testCase.Divisor)
				assert.Equal(
					t, tc.CmpEq, testCase.ExpectedRemainder.Cmp(result), "remainder equals expected",
				)
			})

			t.Run("DivMod", func(t *testing.T) {
				dividend, remainder := testCase.Tc.DivMod(testCase.Divisor)
				assert.Equal(
					t,
					tc.CmpEq,
					testCase.ExpectedDividend.Cmp(dividend),
					"dividend (%v) equals expected (%v)",
					dividend,
					testCase.ExpectedDividend,
				)

				assert.Equal(
					t,
					tc.CmpEq,
					testCase.ExpectedRemainder.Cmp(remainder),
					"remainder (%v) equals expected (%v)",
					dividend,
					testCase.ExpectedDividend,
					testCase.ExpectedRemainder,
				)
			})
		})
	}
}

func TestTimecode_Neg(t *testing.T) {
	cases := []struct {
		Tc       tc.Timecode
		Expected tc.Timecode
	}{
		{
			Tc:       mustTC("01:00:00:00", rate.F24),
			Expected: mustTC("-01:00:00:00", rate.F24),
		},
		{
			Tc:       mustTC("-01:00:00:00", rate.F24),
			Expected: mustTC("01:00:00:00", rate.F24),
		},
	}

	for _, testCase := range cases {
		t.Run(fmt.Sprint(testCase.Tc), func(t *testing.T) {
			result := testCase.Tc.Neg()
			assert.Equal(
				t, tc.CmpEq, testCase.Expected.Cmp(result), "negative equals expected",
			)
		})
	}
}

func TestTimecode_Abs(t *testing.T) {
	cases := []struct {
		Tc       tc.Timecode
		Expected tc.Timecode
	}{
		{
			Tc:       mustTC("01:00:00:00", rate.F24),
			Expected: mustTC("01:00:00:00", rate.F24),
		},
		{
			Tc:       mustTC("-01:00:00:00", rate.F24),
			Expected: mustTC("01:00:00:00", rate.F24),
		},
	}

	for _, testCase := range cases {
		t.Run(fmt.Sprint(testCase.Tc), func(t *testing.T) {
			result := testCase.Tc.Abs()
			assert.Equal(
				t, tc.CmpEq, testCase.Expected.Cmp(result), "negative equals expected",
			)
		})
	}
}

func TestTimecode_Rebase(t *testing.T) {
	cases := []struct {
		Tc       tc.Timecode
		Rate     rate.Framerate
		Expected tc.Timecode
	}{
		{
			Tc:       mustTC("01:00:00:00", rate.F24),
			Rate:     rate.F48,
			Expected: mustTC("00:30:00:00", rate.F48),
		},
	}

	for _, testCase := range cases {
		name := fmt.Sprintf("%v @ %v", testCase.Tc, testCase.Rate)
		t.Run(name, func(t *testing.T) {
			result := testCase.Tc.Rebase(testCase.Rate)
			assert.Equal(
				t, tc.CmpEq, testCase.Expected.Cmp(result), "rebase equals expected",
			)
			assert.Equal(
				t, testCase.Rate.NTSC(), result.Rate().NTSC(), "rate NTSC expected",
			)
		})
	}
}
