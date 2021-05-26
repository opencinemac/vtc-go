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

func testParseTimecodeInfo(t *testing.T, data testdata.TimecodeData) {
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
