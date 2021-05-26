package tc_test

import (
	"fmt"
	"github.com/opencinemac/vtc-go/pkg/internal/testdata"
	"github.com/opencinemac/vtc-go/pkg/rate"
	"github.com/opencinemac/vtc-go/pkg/tc"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestParseTableTests(t *testing.T) {
	t.Run("sequence stat time", func(t *testing.T) {
		testParseTimecodeInfo(t, testdata.ManyBasicEditsData.StartTime)
	})

	for i, event := range testdata.ManyBasicEditsData.Events {
		t.Run(fmt.Sprintf("Event %03d", i), func(t *testing.T) {
			t.Run("Record In", func(t *testing.T) {
				testParseTimecodeInfo(t, event.RecordIn)
			})
			t.Run("Record Out", func(t *testing.T) {
				testParseTimecodeInfo(t, event.RecordOut)
			})
			t.Run("Source In", func(t *testing.T) {
				testParseTimecodeInfo(t, event.SourceIn)
			})
			t.Run("Source Out", func(t *testing.T) {
				testParseTimecodeInfo(t, event.SourceOut)
			})
		})
	}
}

// getTimecodeRate gets the Framerate from a testdata.TimecodeData value.
func getTimecodeRate(t *testing.T, data testdata.TimecodeData) rate.Framerate {
	assert := assert.New(t)
	ntsc := rate.NTSCNone
	if data.Ntsc {
		ntsc = rate.NTSCNonDrop
	}
	if data.DropFrame {
		ntsc = rate.NTSCDrop
	}

	framerate, err := rate.FromInt(data.Timebase, ntsc)
	if !assert.NoError(err, "parse timebase to framerate") {
		t.FailNow()
	}

	return framerate
}

func testParseTimecodeInfo(t *testing.T, data testdata.TimecodeData) {
	assert := assert.New(t)

	framerate := getTimecodeRate(t, data)

	t.Run("From Timecode", func(t *testing.T) {
		parsed, err := tc.FromTimecode(data.Timecode, framerate)
		if !assert.NoError(err, "parse from timecode") {
			t.FailNow()
		}

		checkTimecodeParse(t, parsed, data)
	})

	t.Run("From Frames", func(t *testing.T) {
		parsed := tc.FromFrames(data.Frame, framerate)
		checkTimecodeParse(t, parsed, data)
	})

	t.Run("From Seconds", func(t *testing.T) {
		seconds := big.Rat(*data.SecondsRational)
		parsed := tc.FromSeconds(&seconds, framerate)
		checkTimecodeParse(t, parsed, data)
	})

	t.Run("From Runtime", func(t *testing.T) {
		parsed, err := tc.FromRuntime(data.Runtime, framerate)
		if !assert.NoError(err, "parse from runtime") {
			t.FailNow()
		}
		checkTimecodeParse(t, parsed, data)
	})

	t.Run("From PPro Ticks", func(t *testing.T) {
		parsed := tc.FromPremiereTicks(data.PProTicks, framerate)
		checkTimecodeParse(t, parsed, data)
	})

	t.Run("From Feet And Frames", func(t *testing.T) {
		parsed, err := tc.FromFeetAndFrames(data.FeetAndFrames, framerate)
		if !assert.NoError(err, "parse from feet and frames") {
			t.FailNow()
		}
		checkTimecodeParse(t, parsed, data)
		checkTimecodeParse(t, parsed, data)
	})
}

func checkTimecodeParse(t *testing.T, parsed tc.Timecode, data testdata.TimecodeData) {
	assert := assert.New(t)

	expectedSeconds := big.Rat(*data.SecondsRational)

	assert.Equal(data.Timecode, parsed.Timecode(), "timecode")
	assert.Equal(data.Frame, parsed.Frames(), "frames")
	assert.Equal(&expectedSeconds, parsed.Seconds(), "seconds")
	assert.Equal(data.Runtime, parsed.Runtime(9), "runtime")
	assert.Equal(data.PProTicks, parsed.PremiereTicks(), "ppro ticks")
	assert.Equal(data.FeetAndFrames, parsed.FeetAndFrames(), "feet and frames")
}

// TestTallySequenceTimecode is going to keep a running tally of the current record time
// and frame count. For each event, then length of the record and source time will be
// checked to see if it is correct against the raw value from the FCP7XML.
//
// We will also check that by adding the length of the event to a running timecode
// counter, we get the current record out time of each event.
//
// This test ensures that our addition and subtraction works as intended and does not
// drift over time.
func TestTallySequenceTimecode(t *testing.T) {
	seq := testdata.ManyBasicEditsData
	seqRate := getTimecodeRate(t, seq.StartTime)

	// Current total will start at 0.
	currentTotal := tc.FromFrames(0, seqRate)

	// currentTc will start at the sequence start time. When we add the length of each
	// event, we should get the record out of the event as a result.
	currentTc := tc.FromFrames(seq.StartTime.Frame, seqRate)

	for i, event := range seq.Events {
		failNow := false

		t.Run(fmt.Sprintf("Event %03d", i), func(t *testing.T) {
			// We're going to add the source length to the running tallies. We use the
			// source so it's more removed from the record timecode we are hoping to
			// match.
			var srcLength tc.Timecode

			// Check that we can calculate the length of the source timecode correctly.
			t.Run("Source Length", func(t *testing.T) {
				t.Cleanup(func() {
					if t.Failed() {
						failNow = true
					}
				})

				sourceIn := getTimecodeFromData(t, event.SourceIn)
				sourceOut := getTimecodeFromData(t, event.SourceOut)

				srcLength = sourceOut.Sub(sourceIn)

				assert.Equal(t, event.DurationFrames, srcLength.Frames(), "record length expected")
			})

			var recordOut tc.Timecode

			// Check that we can calculate the length of the record timecode correctly.
			t.Run("Record Length", func(t *testing.T) {
				t.Cleanup(func() {
					if t.Failed() {
						failNow = true
					}
				})

				recordIn := getTimecodeFromData(t, event.RecordIn)
				recordOut = getTimecodeFromData(t, event.RecordOut)

				length := recordOut.Sub(recordIn)

				assert.Equal(t, event.DurationFrames, length.Frames(), "record length expected")
			})

			// Check that our running tallies are correct.
			t.Run("Running Record Time", func(t *testing.T) {
				t.Cleanup(func() {
					if t.Failed() {
						failNow = true
					}
				})

				currentTotal = currentTotal.Add(srcLength)
				currentTc = currentTc.Add(srcLength)

				assert.Equal(
					t,
					tc.CmpEq,
					currentTc.Cmp(recordOut),
					"current running tc (%v) is record out (%v)",
					currentTc,
					recordOut,
				)
			})

		})

		if failNow {
			t.FailNow()
		}
	}

	// Check that our running frame total is correct.
	t.Run("Frame Total", func(t *testing.T) {
		assert.Equal(
			t,
			seq.TotalDurationFrames,
			currentTotal.Frames(),
			"frame total matches",
		)
	})
}

// getTimecodeFromData extracts a tc.Timecode from testdata.TimecodeData using the
// timecode string and framerate.
func getTimecodeFromData(t *testing.T, data testdata.TimecodeData) tc.Timecode {
	t.Helper()

	framerate := getTimecodeRate(t, data)

	timecode, err := tc.FromTimecode(data.Timecode, framerate)
	if !assert.NoError(t, err, "parse timecode string") {
		t.FailNow()
	}

	return timecode
}
