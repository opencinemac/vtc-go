package tc

import (
	"fmt"
	"github.com/opencinemac/vtc-go/pkg/rate"
	"math"
)

// dropFrameNumAdjustment returns the adjustment to apply to the frame number in
// drop-frame timecode calculations.
//
// algorithm adapted from:
// https://www.davidheidelberger.com/2010/06/10/drop-frame-timecode/
func dropFrameNumAdjustment(frameNumber int64, framerate rate.Framerate) int64 {
	// The timebase of a drop-frame timecode will always be a whole-number, so we can
	// just get the numerator as our timebase.
	timebase := framerate.Timebase().Num().Int64()

	// Get the number of frames we need to drop each time we drop frames
	// (ex: 2 for 29.97).
	dropFrames := dropFramesForTimebase(timebase)

	framesPerMinute := timebase * 60

	// Get the number of frames are in a minute where we have dropped frames at the
	// beginning.
	framesPerMinuteDrop := timebase*60 - dropFrames

	// Get the number of actual frames in a 10-minute span for drop frame timecode.
	// Since we drop 9 times in 10 minutes, it will be 9 drop-minute frame counts + 1
	// whole-minute frame count.
	framesPer10MinuteDrop := framesPerMinuteDrop*9 + framesPerMinute

	tensOfMinutes := frameNumber / framesPer10MinuteDrop
	frames := frameNumber % framesPer10MinuteDrop

	// Create an adjustment for the number of 10s of minutes. It will be 9 times the
	// drop value (we drop for the first 9 minutes, then leave the 10th alone).
	adjustment := 9 * dropFrames * tensOfMinutes

	// If our remaining frames are less than a whole minute, we aren't going to drop
	// again.
	if frames < framesPerMinute {
		return adjustment
	}

	// Remove the first full minute (we don't drop until the next minute) and add the
	// drop-rate to the adjustment.
	frames -= timebase
	adjustment += dropFrames

	// Get the number of remaining drop-minutes present, and add a drop adjustment for
	// each.
	minutesDrop := frames / framesPerMinuteDrop
	adjustment += minutesDrop * dropFrames

	return adjustment
}

// dropFrameParseAdjustment creates a frame number adjustment for parsing drop-frame
// timecode.
func dropFrameParseAdjustment(sections TimecodeSections, framerate rate.Framerate) (int64, error) {
	// Drop-frame timebases are always whole-numbers.
	timebase := framerate.Timebase().Num().Int64()
	dropFrames := dropFramesForTimebase(timebase)

	hasBadFrames := sections.Frames < dropFrames
	isTenthMinute := sections.Minutes%10 == 0

	if hasBadFrames && !isTenthMinute {
		return 0, fmt.Errorf(
			"%w: found frame value of '%v', should be < '%v'",
			ErrBadDropFrameValue,
			sections.Frames,
			dropFrames,
		)
	}

	totalMinutes := 60*sections.Hours + sections.Minutes
	// calculate the adjustment, we need to remove two frames for each minute except for
	// every 10th minute.
	adjustment := dropFrames * (totalMinutes - totalMinutes/10)

	// We need the adjustment to remove frames, so return a negative.
	return -adjustment, nil
}

// dropFramesForTimebase returns the number of frames to skip when we drop frames
// based on the timebase.
func dropFramesForTimebase(timebase int64) int64 {
	return int64(math.Round(float64(timebase) * 0.066666))
}
