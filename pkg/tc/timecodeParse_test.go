package tc_test

import (
	"fmt"
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
	FeetAndFrames string
}

func TestParseTimecode(t *testing.T) {
	cases := []ParseCase{
		// 23.98 NTSC
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
		// 24 fps
		{
			Name:          "01:00:00:00 24 fps",
			Rate:          rate.F24,
			Seconds:       big.NewRat(3600, 1),
			Frames:        86400,
			Timecode:      "01:00:00:00",
			Runtime:       "01:00:00.0",
			PremiereTicks: 914457600000000,
			FeetAndFrames: "5400+00",
		},
		{
			Name:          "00:40:00:00 24 fps",
			Rate:          rate.F24,
			Seconds:       big.NewRat(2400, 1),
			Frames:        57600,
			Timecode:      "00:40:00:00",
			Runtime:       "00:40:00.0",
			PremiereTicks: 609638400000000,
			FeetAndFrames: "3600+00",
		},
		// 29.97 Drop-frame
		{
			Name:          "00:00:00;00 29.97 Drop-Frame",
			Rate:          rate.F29_97Df,
			Seconds:       big.NewRat(0, 1),
			Frames:        0,
			Timecode:      "00:00:00;00",
			Runtime:       "00:00:00.0",
			PremiereTicks: 0,
			FeetAndFrames: "0+00",
		},
		{
			Name:          "00:00:02;02 29.97 Drop-Frame",
			Rate:          rate.F29_97Df,
			Seconds:       big.NewRat(31031, 15000),
			Frames:        62,
			Timecode:      "00:00:02;02",
			Runtime:       "00:00:02.068733333",
			PremiereTicks: 525491366400,
			FeetAndFrames: "3+14",
		},
		{
			Name:          "00:01:00;02 29.97 Drop-Frame",
			Rate:          rate.F29_97Df,
			Seconds:       big.NewRat(3003, 50),
			Frames:        1800,
			Timecode:      "00:01:00;02",
			Runtime:       "00:01:00.06",
			PremiereTicks: 15256200960000,
			FeetAndFrames: "112+08",
		},
		{
			Name:          "00:10:00;00 29.97 Drop-Frame",
			Rate:          rate.F29_97Df,
			Seconds:       big.NewRat(2999997, 5000),
			Frames:        17982,
			Timecode:      "00:10:00;00",
			Runtime:       "00:09:59.9994",
			PremiereTicks: 152409447590400,
			FeetAndFrames: "1123+14",
		},
		{
			Name:          "00:11:00;02 29.97 Drop-Frame",
			Rate:          rate.F29_97Df,
			Seconds:       big.NewRat(3300297, 5000),
			Frames:        19782,
			Timecode:      "00:11:00;02",
			Runtime:       "00:11:00.0594",
			PremiereTicks: 167665648550400,
			FeetAndFrames: "1236+06",
		},
		{
			Name:          "01:00:00;00 29.97 Drop-Frame",
			Rate:          rate.F29_97Df,
			Seconds:       big.NewRat(8999991, 2500),
			Frames:        107892,
			Timecode:      "01:00:00;00",
			Runtime:       "00:59:59.9964",
			PremiereTicks: 914456685542400,
			FeetAndFrames: "6743+04",
		},
		// 59.94 DF Cases
		{
			Name:          "00:00:00;00 59.94 DF",
			Rate:          rate.F59_94Df,
			Seconds:       big.NewRat(0, 1),
			Frames:        0,
			Timecode:      "00:00:00;00",
			Runtime:       "00:00:00.0",
			PremiereTicks: 0,
			FeetAndFrames: "0+00",
		},
		{
			Name:          "00:00:01;01 59.94 DF",
			Rate:          rate.F59_94Df,
			Seconds:       big.NewRat(61061, 60000),
			Frames:        61,
			Timecode:      "00:00:01;01",
			Runtime:       "00:00:01.017683333",
			PremiereTicks: 258507849600,
			FeetAndFrames: "3+13",
		},
		{
			Name:          "00:00:01;03 59.94 DF",
			Rate:          rate.F59_94Df,
			Seconds:       big.NewRat(21021, 20000),
			Frames:        63,
			Timecode:      "00:00:01;03",
			Runtime:       "00:00:01.05105",
			PremiereTicks: 266983516800,
			FeetAndFrames: "3+15",
		},
		{
			Name:          "00:01:00;04 59.94 DF",
			Rate:          rate.F59_94Df,
			Seconds:       big.NewRat(3003, 50),
			Frames:        3600,
			Timecode:      "00:01:00;04",
			Runtime:       "00:01:00.06",
			PremiereTicks: 15256200960000,
			FeetAndFrames: "225+00",
		},
		// 239.76 NDF CASES ---------------------
		// We're going to use this to test very large values beyond what you would
		// normally see in the wild to put pressure on possible integer overflow points.
		//
		// This value represents a timecode of over 123 hours running at 240 fps. In the
		// real world, one would be VERY unlikely to see a timecode like this. We are
		// using an NTSC timebase as NTSC bases are far more likely to create large
		// numerators / denominators.
		{
			Name: "123:17:34:217 239.9476 NTSC",
			Rate: func() rate.Framerate {
				framerate, err := rate.FromFloat(239.76, rate.NTSCNonDrop)
				if err != nil {
					panic(fmt.Errorf("error making framerata: %w", err))
				}
				return framerate
			}(),
			Seconds:       big.NewRat(106631702177, 240000),
			Frames:        106525177,
			Timecode:      "123:17:34:217",
			Runtime:       "123:24:58.759070833",
			PremiereTicks: 112858993584136800,
			FeetAndFrames: "6657823+09",
		},
	}

	for _, thisCase := range cases {
		t.Run(thisCase.Name, func(t *testing.T) {
			t.Run("Positive", func(t *testing.T) {
				runParseCase(t, thisCase)
			})

			// Don't need to test negative on 0 values.
			if thisCase.Frames == 0 {
				return
			}

			t.Run("Negative", func(t *testing.T) {
				thisCase.Timecode = "-" + thisCase.Timecode
				thisCase.Frames = -thisCase.Frames
				thisCase.Seconds = thisCase.Seconds.Neg(thisCase.Seconds)
				thisCase.Runtime = "-" + thisCase.Runtime
				thisCase.PremiereTicks = -thisCase.PremiereTicks
				thisCase.FeetAndFrames = "-" + thisCase.FeetAndFrames

				runParseCase(t, thisCase)
			})
		})
	}
}

func runParseCase(t *testing.T, thisCase ParseCase) {
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
}

// checkParse checks that was parsed correctly.
func checkParse(t *testing.T, thisCase ParseCase, parsed tc.Timecode, err error) {
	assert := assert.New(t)

	if !assert.NoError(err, "parse error") {
		t.FailNow()
	}

	assert.Equal(
		thisCase.Seconds,
		parsed.Seconds(),
		"seconds. expected: %v, got: %v",
		thisCase.Seconds,
		parsed.Seconds(),
	)
	assert.Equal(thisCase.Frames, parsed.Frames(), "frames")
	assert.Equal(thisCase.Timecode, parsed.Timecode(), "timecode")
	assert.Equal(thisCase.Runtime, parsed.Runtime(9), "runtime")
	assert.Equal(thisCase.PremiereTicks, parsed.PremiereTicks(), "Premiere Ticks")
	assert.Equal(thisCase.FeetAndFrames, parsed.FeetAndFrames(), "Feet And Frames")
}

// TesTcOverflowParsing tests that tc strings with overflowed values are parsed
// correctly.
func TestTcOverflowParsing(t *testing.T) {
	cases := []struct {
		In       string
		Expected string
	}{
		{
			In:       "00:59:59:24",
			Expected: "01:00:00:00",
		},
		{
			In:       "00:59:59:28",
			Expected: "01:00:00:04",
		},
		{
			In:       "00:00:62:04",
			Expected: "00:01:02:04",
		},
		{
			In:       "00:62:01:04",
			Expected: "01:02:01:04",
		},
		{
			In:       "00:62:62:04",
			Expected: "01:03:02:04",
		},
		{
			In:       "123:00:00:00",
			Expected: "123:00:00:00",
		},
		{
			In:       "01:00:00:48",
			Expected: "01:00:02:00",
		},
		{
			In:       "01:00:120:00",
			Expected: "01:02:00:00",
		},
		{
			In:       "01:120:00:00",
			Expected: "03:00:00:00",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.In, func(t *testing.T) {
			assert := assert.New(t)

			timecode, err := tc.FromTimecode(testCase.In, rate.F24)
			if !assert.NoError(err, "parse In") {
				t.FailNow()
			}

			assert.Equal(testCase.Expected, timecode.Timecode(), "fixed tc")
		})
	}
}

func TestParsePartialTc(t *testing.T) {
	cases := []struct {
		In       string
		Expected string
	}{
		{
			In:       "1:02:03:04",
			Expected: "01:02:03:04",
		},
		{
			In:       "02:03:04",
			Expected: "00:02:03:04",
		},
		{
			In:       "2:03:04",
			Expected: "00:02:03:04",
		},
		{
			In:       "03:04",
			Expected: "00:00:03:04",
		},
		{
			In:       "3:04",
			Expected: "00:00:03:04",
		},
		{
			In:       "04",
			Expected: "00:00:00:04",
		},
		{
			In:       "4",
			Expected: "00:00:00:04",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.In, func(t *testing.T) {
			assert := assert.New(t)

			timecode, err := tc.FromTimecode(testCase.In, rate.F24)
			if !assert.NoError(err, "parse In") {
				t.FailNow()
			}

			assert.Equal(testCase.Expected, timecode.Timecode(), "fixed tc")
		})
	}
}

func TestParseRuntimePartial(t *testing.T) {
	cases := []struct {
		In       string
		Expected string
	}{
		{
			In:       "1:02:03.5",
			Expected: "01:02:03.5",
		},
		{
			In:       "02:03.5",
			Expected: "00:02:03.5",
		},
		{
			In:       "2:03.5",
			Expected: "00:02:03.5",
		},
		{
			In:       "03.5",
			Expected: "00:00:03.5",
		},
		{
			In:       "3.5",
			Expected: "00:00:03.5",
		},
		{
			In:       "0.5",
			Expected: "00:00:00.5",
		},
		{
			In:       "1:2:3.5",
			Expected: "01:02:03.5",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.In, func(t *testing.T) {
			assert := assert.New(t)

			timecode, err := tc.FromRuntime(testCase.In, rate.F24)
			if !assert.NoError(err, "parse In") {
				t.FailNow()
			}

			assert.Equal(testCase.Expected, timecode.Runtime(9), "fixed runtime")
		})
	}
}

func TestFromTimecode_ErrFormat(t *testing.T) {
	assert := assert.New(t)

	_, err := tc.FromTimecode("not a timecode", rate.F24)

	if !assert.Error(err, "error occurred") {
		t.FailNow()
	}
	assert.ErrorIs(err, tc.ErrParseTimecode, "is parse err")
	assert.ErrorIs(err, tc.ErrFormatNotRecognized, "is correct sub err")
}

func TestFromTimecode_BadDropFrames(t *testing.T) {
	assert := assert.New(t)

	_, err := tc.FromTimecode("00:01:00:01", rate.F29_97Df)

	if !assert.Error(err, "error occurred") {
		t.FailNow()
	}
	assert.ErrorIs(err, tc.ErrParseTimecode, "is parse err")
	assert.ErrorIs(err, tc.ErrBadDropFrameValue, "is correct sub err")
}

func TestFromRuntime_ErrFormat(t *testing.T) {
	assert := assert.New(t)

	_, err := tc.FromRuntime("not a timecode", rate.F24)

	if !assert.Error(err, "error occurred") {
		t.FailNow()
	}
	assert.ErrorIs(err, tc.ErrParseTimecode, "is parse err")
	assert.ErrorIs(err, tc.ErrFormatNotRecognized, "is correct sub err")
}

func TestFromFeetAndFrames_ErrFormat(t *testing.T) {
	assert := assert.New(t)

	_, err := tc.FromFeetAndFrames("not a timecode", rate.F24)

	if !assert.Error(err, "error occurred") {
		t.FailNow()
	}
	assert.ErrorIs(err, tc.ErrParseTimecode, "is parse err")
	assert.ErrorIs(err, tc.ErrFormatNotRecognized, "is correct sub err")
}
