package rate_test

import (
	"fmt"
	"github.com/opencinemac/vtc-go/pkg/rate"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

// ExpectedFramerate is used by framerate parse tests to declare the expected values.
type ExpectedFramerate struct {
	// Ntsc is the expected Framerate.NTSC return.
	Ntsc rate.NTSC
	// Playback is the expected Framerate.Playback return.
	Playback *big.Rat
	// Timebase is the expected Framerate.Timebase return.
	Timebase *big.Rat
	// Err is any expected error. Other values will not be checked if Err is not nil.
	Err error
}

// TestFromRat tests our logic parsing from a *big.Rat value.
func TestFromRat(t *testing.T) {
	cases := []struct {
		Name     string
		Input    []*big.Rat
		Expected ExpectedFramerate
	}{
		// 23.98 NTSC
		{
			Name: "23.98 NTSC Non-Drop",
			Input: []*big.Rat{
				big.NewRat(24000, 1001),
				big.NewRat(24, 1),
			},
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNonDrop,
				Playback: big.NewRat(24000, 1001),
				Timebase: big.NewRat(24, 1),
				Err:      nil,
			},
		},
		// 29.97 NTSC
		{
			Name: "29.97 NTSC Non-Drop",
			Input: []*big.Rat{
				big.NewRat(30000, 1001),
				big.NewRat(30, 1),
			},
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNonDrop,
				Playback: big.NewRat(30000, 1001),
				Timebase: big.NewRat(30, 1),
				Err:      nil,
			},
		},
		{
			Name: "29.97 NTSC Drop",
			Input: []*big.Rat{
				big.NewRat(30000, 1001),
				big.NewRat(30, 1),
			},
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCDrop,
				Playback: big.NewRat(30000, 1001),
				Timebase: big.NewRat(30, 1),
				Err:      nil,
			},
		},
		// 59.94 NTSC
		{
			Name: "59.94 NTSC Non-Drop",
			Input: []*big.Rat{
				big.NewRat(60000, 1001),
				big.NewRat(60, 1),
			},
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNonDrop,
				Playback: big.NewRat(60000, 1001),
				Timebase: big.NewRat(60, 1),
				Err:      nil,
			},
		},
		{
			Name: "59.94 NTSC Drop",
			Input: []*big.Rat{
				big.NewRat(60000, 1001),
				big.NewRat(60, 1),
			},
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCDrop,
				Playback: big.NewRat(60000, 1001),
				Timebase: big.NewRat(60, 1),
				Err:      nil,
			},
		},
		// Non-NTSC
		{
			Name:  "24 fps",
			Input: []*big.Rat{big.NewRat(24, 1)},
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNone,
				Playback: big.NewRat(24, 1),
				Timebase: big.NewRat(24, 1),
				Err:      nil,
			},
		},
		{
			Name:  "24000/1001 fps",
			Input: []*big.Rat{big.NewRat(24000, 1001)},
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNone,
				Playback: big.NewRat(24000, 1001),
				Timebase: big.NewRat(24000, 1001),
				Err:      nil,
			},
		},
		// ERRORS
		{
			Name:  "Error Negative",
			Input: []*big.Rat{big.NewRat(-24000, 1001)},
			Expected: ExpectedFramerate{
				Ntsc: rate.NTSCNonDrop,
				Err:  rate.ErrNegative,
			},
		},
		{
			Name:  "Error Bad Drop-Frame",
			Input: []*big.Rat{big.NewRat(24000, 1001)},
			Expected: ExpectedFramerate{
				Ntsc: rate.NTSCDrop,
				Err:  rate.ErrBadDropFrameRate,
			},
		},
		{
			Name:  "Bad Bad Ntsc",
			Input: []*big.Rat{big.NewRat(24000, 1001)},
			Expected: ExpectedFramerate{
				Ntsc: 100,
				Err:  rate.ErrBadNtsc,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			for _, source := range tc.Input {
				t.Run(fmt.Sprint(source), func(t *testing.T) {
					framerate, err := rate.FromRat(source, tc.Expected.Ntsc)
					checkParse(t, framerate, err, tc.Expected)
				})
			}
		})
	}
}

// TestFromInt tests our logic parsing from an int64 value.
func TestFromInt(t *testing.T) {
	cases := []struct {
		Name     string
		Input    int64
		Expected ExpectedFramerate
	}{
		// 23.98 NTSC
		{
			Name:  "24 NTSC Non-Drop",
			Input: 24,
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNonDrop,
				Playback: big.NewRat(24000, 1001),
				Timebase: big.NewRat(24, 1),
				Err:      nil,
			},
		},
		// 29.97 NTSC
		{
			Name:  "30 NTSC Non-Drop",
			Input: 30,
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNonDrop,
				Playback: big.NewRat(30000, 1001),
				Timebase: big.NewRat(30, 1),
				Err:      nil,
			},
		},
		{
			Name:  "30 NTSC Drop",
			Input: 30,
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCDrop,
				Playback: big.NewRat(30000, 1001),
				Timebase: big.NewRat(30, 1),
				Err:      nil,
			},
		},
		// 59.94 NTSC
		{
			Name:  "60 NTSC Non-Drop",
			Input: 60,
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNonDrop,
				Playback: big.NewRat(60000, 1001),
				Timebase: big.NewRat(60, 1),
				Err:      nil,
			},
		},
		{
			Name:  "60 NTSC Drop",
			Input: 60,
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCDrop,
				Playback: big.NewRat(60000, 1001),
				Timebase: big.NewRat(60, 1),
				Err:      nil,
			},
		},
		// Non-NTSC
		{
			Name:  "24 fps",
			Input: 24,
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNone,
				Playback: big.NewRat(24, 1),
				Timebase: big.NewRat(24, 1),
				Err:      nil,
			},
		},
		// ERRORS
		{
			Name:  "Error Negative",
			Input: -24,
			Expected: ExpectedFramerate{
				Ntsc: rate.NTSCNonDrop,
				Err:  rate.ErrNegative,
			},
		},
		{
			Name:  "Bad Bad Ntsc",
			Input: 24,
			Expected: ExpectedFramerate{
				Ntsc: 100,
				Err:  rate.ErrBadNtsc,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			framerate, err := rate.FromInt(tc.Input, tc.Expected.Ntsc)
			checkParse(t, framerate, err, tc.Expected)
		})
	}
}

// TestFromFloat tests our logic parsing from a float64 value.
func TestFromFloat(t *testing.T) {
	cases := []struct {
		Name     string
		Input    []float64
		Expected ExpectedFramerate
	}{
		// 23.98 NTSC
		{
			Name:  "23.98 NTSC Non-Drop",
			Input: []float64{23.98, 23.976, 24, 23.5},
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNonDrop,
				Playback: big.NewRat(24000, 1001),
				Timebase: big.NewRat(24, 1),
				Err:      nil,
			},
		},
		// 29.97 NTSC
		{
			Name:  "29.97 NTSC Non-Drop",
			Input: []float64{29.97, 30},
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNonDrop,
				Playback: big.NewRat(30000, 1001),
				Timebase: big.NewRat(30, 1),
				Err:      nil,
			},
		},
		{
			Name:  "29.97 NTSC Drop",
			Input: []float64{29.97, 30},
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCDrop,
				Playback: big.NewRat(30000, 1001),
				Timebase: big.NewRat(30, 1),
				Err:      nil,
			},
		},
		// 59.94 NTSC
		{
			Name:  "59.94 NTSC Non-Drop",
			Input: []float64{59.94, 60},
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNonDrop,
				Playback: big.NewRat(60000, 1001),
				Timebase: big.NewRat(60, 1),
				Err:      nil,
			},
		},
		{
			Name:  "59.94 NTSC Drop",
			Input: []float64{59.94, 60},
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCDrop,
				Playback: big.NewRat(60000, 1001),
				Timebase: big.NewRat(60, 1),
				Err:      nil,
			},
		},
		// ERRORS
		{
			Name:  "Error Negative",
			Input: []float64{-23.98},
			Expected: ExpectedFramerate{
				Ntsc: rate.NTSCNonDrop,
				Err:  rate.ErrNegative,
			},
		},
		{
			Name:  "Error Bad Drop-Frame",
			Input: []float64{23.98},
			Expected: ExpectedFramerate{
				Ntsc: rate.NTSCDrop,
				Err:  rate.ErrBadDropFrameRate,
			},
		},
		{
			Name:  "Error Imprecise",
			Input: []float64{29.97, 59.94},
			Expected: ExpectedFramerate{
				Ntsc: rate.NTSCNone,
				Err:  rate.ErrImprecise,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			for _, source := range tc.Input {
				t.Run(fmt.Sprint(source), func(t *testing.T) {
					framerate, err := rate.FromFloat(source, tc.Expected.Ntsc)
					checkParse(t, framerate, err, tc.Expected)
				})
			}
		})
	}
}

// TestFromRat tests our logic parsing from a string value.
func TestFromString(t *testing.T) {
	cases := []struct {
		Name     string
		Input    []string
		Expected ExpectedFramerate
	}{
		// 23.98 NTSC
		{
			Name: "23.98 NTSC Non-Drop",
			Input: []string{
				"24/1",
				"1/24",
				"24000/1001",
				"1001/24000",
				"24",
				"23.98",
				"23.976",
				"23.5",
			},
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNonDrop,
				Playback: big.NewRat(24000, 1001),
				Timebase: big.NewRat(24, 1),
				Err:      nil,
			},
		},
		// 29.97 NTSC
		{
			Name: "29.97 NTSC Non-Drop",
			Input: []string{
				"30/1",
				"1/30",
				"30000/1001",
				"1001/30000",
				"30",
				"29.97",
			},
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNonDrop,
				Playback: big.NewRat(30000, 1001),
				Timebase: big.NewRat(30, 1),
				Err:      nil,
			},
		},
		{
			Name: "29.97 NTSC Drop",
			Input: []string{
				"30/1",
				"1/30",
				"30000/1001",
				"1001/30000",
				"30",
				"29.97",
			},
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCDrop,
				Playback: big.NewRat(30000, 1001),
				Timebase: big.NewRat(30, 1),
				Err:      nil,
			},
		},
		// 59.94 NTSC
		{
			Name: "59.94 NTSC Non-Drop",
			Input: []string{
				"60/1",
				"1/60",
				"60000/1001",
				"1001/60000",
				"60",
				"59.94",
			},
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNonDrop,
				Playback: big.NewRat(60000, 1001),
				Timebase: big.NewRat(60, 1),
				Err:      nil,
			},
		},
		{
			Name: "59.94 NTSC Drop",
			Input: []string{
				"60/1",
				"1/60",
				"60000/1001",
				"1001/60000",
				"60",
				"59.94",
			},
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCDrop,
				Playback: big.NewRat(60000, 1001),
				Timebase: big.NewRat(60, 1),
				Err:      nil,
			},
		},
		// Non-NTSC
		{
			Name: "24 fps",
			Input: []string{
				"24/1",
				"1/24",
				"24",
			},
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNone,
				Playback: big.NewRat(24, 1),
				Timebase: big.NewRat(24, 1),
				Err:      nil,
			},
		},
		{
			Name: "24000/1001 fps",
			Input: []string{
				"24000/1001",
				"1001/24000",
			},
			Expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNone,
				Playback: big.NewRat(24000, 1001),
				Timebase: big.NewRat(24000, 1001),
				Err:      nil,
			},
		},
		// ERRORS
		{
			Name: "Error Negative",
			Input: []string{
				"-24/1",
				"-1/24",
				"-24000/1001",
				"-1001/24000",
				"-24",
				"-23.98",
				"-23.976",
				"-23.5",
			},
			Expected: ExpectedFramerate{
				Ntsc: rate.NTSCNonDrop,
				Err:  rate.ErrNegative,
			},
		},
		{
			Name: "Error Bad Drop-Frame",
			Input: []string{
				"24/1",
				"1/24",
				"24000/1001",
				"1001/24000",
				"24",
				"23.98",
				"23.976",
				"23.5",
			},
			Expected: ExpectedFramerate{
				Ntsc: rate.NTSCDrop,
				Err:  rate.ErrBadDropFrameRate,
			},
		},
		{
			Name: "Error Bad Ntsc",
			Input: []string{
				"24/1",
				"1/24",
				"24000/1001",
				"1001/24000",
				"24",
			},
			Expected: ExpectedFramerate{
				Ntsc: 100,
				Err:  rate.ErrBadNtsc,
			},
		},
		{
			Name: "Error Imprecise",
			Input: []string{
				"23.98",
				"23.976",
			},
			Expected: ExpectedFramerate{
				Ntsc: rate.NTSCNone,
				Err:  rate.ErrImprecise,
			},
		},
		{
			Name: "Not Recognized",
			Input: []string{
				"Not a Framerate",
			},
			Expected: ExpectedFramerate{
				Ntsc: rate.NTSCNone,
				Err:  rate.ErrParseFramerate,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			for _, source := range tc.Input {
				t.Run(fmt.Sprint(source), func(t *testing.T) {
					framerate, err := rate.FromString(source, tc.Expected.Ntsc)
					checkParse(t, framerate, err, tc.Expected)
				})
			}
		})
	}
}

func TestConstRates(t *testing.T) {
	cases := []struct {
		Name     string
		constant rate.Framerate
		expected ExpectedFramerate
	}{
		{
			Name:     "F23_98",
			constant: rate.F23_98,
			expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNonDrop,
				Playback: big.NewRat(24000, 1001),
				Timebase: big.NewRat(24, 1),
				Err:      nil,
			},
		},
		{
			Name:     "F24",
			constant: rate.F24,
			expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNone,
				Playback: big.NewRat(24, 1),
				Timebase: big.NewRat(24, 1),
				Err:      nil,
			},
		},
		{
			Name:     "F29_97Ndf",
			constant: rate.F29_97Ndf,
			expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNonDrop,
				Playback: big.NewRat(30000, 1001),
				Timebase: big.NewRat(30, 1),
				Err:      nil,
			},
		},
		{
			Name:     "F29_97Df",
			constant: rate.F29_97Df,
			expected: ExpectedFramerate{
				Ntsc:     rate.NTSCDrop,
				Playback: big.NewRat(30000, 1001),
				Timebase: big.NewRat(30, 1),
				Err:      nil,
			},
		},
		{
			Name:     "F30",
			constant: rate.F30,
			expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNone,
				Playback: big.NewRat(30, 1),
				Timebase: big.NewRat(30, 1),
				Err:      nil,
			},
		},
		{
			Name:     "F47_95",
			constant: rate.F47_95,
			expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNonDrop,
				Playback: big.NewRat(48000, 1001),
				Timebase: big.NewRat(48, 1),
				Err:      nil,
			},
		},
		{
			Name:     "F48",
			constant: rate.F48,
			expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNone,
				Playback: big.NewRat(48, 1),
				Timebase: big.NewRat(48, 1),
				Err:      nil,
			},
		},
		{
			Name:     "F59_94Ndf",
			constant: rate.F59_94Ndf,
			expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNonDrop,
				Playback: big.NewRat(60000, 1001),
				Timebase: big.NewRat(60, 1),
				Err:      nil,
			},
		},
		{
			Name:     "F59_94Df",
			constant: rate.F59_94Df,
			expected: ExpectedFramerate{
				Ntsc:     rate.NTSCDrop,
				Playback: big.NewRat(60000, 1001),
				Timebase: big.NewRat(60, 1),
				Err:      nil,
			},
		},
		{
			Name:     "F60",
			constant: rate.F60,
			expected: ExpectedFramerate{
				Ntsc:     rate.NTSCNone,
				Playback: big.NewRat(60, 1),
				Timebase: big.NewRat(60, 1),
				Err:      nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprint(tc.Name), func(t *testing.T) {
			checkParse(t, tc.constant, nil, tc.expected)
		})
	}
}

// checkParse checks that framerate and err adheres to the expected values.
func checkParse(
	t *testing.T,
	framerate rate.Framerate,
	err error,
	expected ExpectedFramerate,
) {
	assert := assert.New(t)

	if expected.Err != nil {
		assert.ErrorIs(err, expected.Err, "error expected.")
		assert.ErrorIs(err, rate.ErrParseFramerate, "is framerate parse error.")
		return
	}

	if !assert.NoError(err, "parse framerate") {
		return
	}

	assert.Equal(expected.Ntsc, framerate.NTSC(), "ntsc value")

	assert.True(
		expected.Playback.Cmp(framerate.Playback()) == 0,
		"playback: expected %v, got %v",
		expected.Playback,
		framerate.Playback(),
	)

	assert.True(
		expected.Timebase.Cmp(framerate.Timebase()) == 0,
		"timebase: expected %v, got %v",
		expected.Timebase,
		framerate.Timebase(),
	)
}

func TestFramerate_String(t *testing.T) {
	cases := []struct {
		Name     string
		Rate     rate.Framerate
		Expected string
	}{
		{
			Name:     "",
			Rate:     rate.F23_98,
			Expected: "23.98 NTSC NDF",
		},
		{
			Name:     "",
			Rate:     rate.F24,
			Expected: "24 fps",
		},
		{
			Name:     "",
			Rate:     rate.F29_97Ndf,
			Expected: "29.97 NTSC NDF",
		},
		{
			Name:     "",
			Rate:     rate.F29_97Df,
			Expected: "29.97 NTSC DF",
		},
		{
			Name:     "",
			Rate:     rate.F30,
			Expected: "30 fps",
		},
		{
			Name:     "",
			Rate:     rate.F47_95,
			Expected: "47.95 NTSC NDF",
		},
		{
			Name:     "",
			Rate:     rate.F48,
			Expected: "48 fps",
		},
		{
			Name:     "",
			Rate:     rate.F59_94Ndf,
			Expected: "59.94 NTSC NDF",
		},
		{
			Name:     "",
			Rate:     rate.F59_94Df,
			Expected: "59.94 NTSC DF",
		},
		{
			Name:     "",
			Rate:     rate.F60,
			Expected: "60 fps",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Log("STRING:", tc.Rate.String())
			assert.Equal(t, tc.Expected, tc.Rate.String())
		})
	}
}

func TestNTSC_String_Invalid(t *testing.T) {
	assert.Equal(t, "[INVALID NTSC VALUE]", rate.NTSC(100).String())
}
