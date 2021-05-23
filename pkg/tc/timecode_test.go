package tc_test

import (
	"github.com/opencinemac/vtc-go/pkg/rate"
	"github.com/opencinemac/vtc-go/pkg/tc"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

// ParseCase hold a timecode parsing test case.
type ParseCase struct {
	Name          string
	Rate          rate.Framerate
	Seconds       *big.Rat
	Frames        int64
	Timecode      string
	Runtime       string
	PremiereTicks int64
	ErrExpected   error
	FeetAndFrames string
}

func TestParseTimecode(t *testing.T) {
	cases := []ParseCase{
		{
			Name:          "01:00:00:00 23.98 NTSC",
			Rate:          rate.F23_98,
			Seconds:       big.NewRat(18018, 5),
			Frames:        86400,
			Timecode:      "01:00:00:00",
			Runtime:       "01:00:03.6",
			PremiereTicks: 915372057600000,
			FeetAndFrames: "5400+00",
		},
		{
			Name:          "00:40:00:00 23.98 NTSC",
			Rate:          rate.F23_98,
			Seconds:       big.NewRat(12012, 5),
			Frames:        57600,
			Timecode:      "00:40:00:00",
			Runtime:       "00:40:02.4",
			PremiereTicks: 610248038400000,
			FeetAndFrames: "3600+00",
		},
	}

	for _, thisCase := range cases {
		t.Run(thisCase.Name, func(t *testing.T) {
			t.Run("From Seconds", func(t *testing.T) {
				parsed := tc.FromSecondsRat(thisCase.Seconds, thisCase.Rate)
				checkParse(t, thisCase, parsed, nil)
			})

			t.Run("From Frames", func(t *testing.T) {
				parsed := tc.FromFrames(thisCase.Frames, thisCase.Rate)
				checkParse(t, thisCase, parsed, nil)
			})

			t.Run("From Timecode", func(t *testing.T) {
				parsed, err := tc.FromTimecode(thisCase.Timecode, thisCase.Rate)
				checkParse(t, thisCase, parsed, err)
			})

			t.Run("From Runtime", func(t *testing.T) {
				parsed, err := tc.FromRuntime(thisCase.Runtime, thisCase.Rate)
				checkParse(t, thisCase, parsed, err)
			})

			t.Run("From Feet and Frames", func(t *testing.T) {
				parsed, err := tc.FromFeetAndFrames(
					thisCase.FeetAndFrames, thisCase.Rate,
				)
				checkParse(t, thisCase, parsed, err)
			})

			t.Run("From Premiere Ticks", func(t *testing.T) {
				parsed := tc.FromPremiereTicks(thisCase.PremiereTicks, thisCase.Rate)
				checkParse(t, thisCase, parsed, nil)
			})
		})
	}
}

// checkParse checks that was parsed correctly.
func checkParse(t *testing.T, thisCase ParseCase, parsed tc.Timecode, err error) {
	assert := assert.New(t)

	if thisCase.ErrExpected != nil {
		assert.ErrorIs(err, thisCase.ErrExpected, "error is expected sentinel")
		assert.ErrorIs(err, tc.ErrParseTimecode, "error is expected sentinel")
		return
	}

	if !assert.NoError(err, "parse error") {
		t.FailNow()
	}

	assert.Equal(thisCase.Seconds, parsed.Seconds(), "seconds")
	assert.Equal(thisCase.Frames, parsed.Frames(), "frames")
	assert.Equal(thisCase.Timecode, parsed.Timecode(), "timecode")
	assert.Equal(thisCase.Runtime, parsed.Runtime(9), "runtime")
	assert.Equal(thisCase.PremiereTicks, parsed.PremiereTicks(), "Premiere Ticks")
	assert.Equal(thisCase.FeetAndFrames, parsed.FeetAndFrames(), "Feet And Frames")
}
