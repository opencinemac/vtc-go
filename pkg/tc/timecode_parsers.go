package tc

import (
	"github.com/opencinemac/vtc-go/pkg/internal"
	"github.com/opencinemac/vtc-go/pkg/rate"
	"math/big"
	"regexp"
	"strconv"
)

// FromSecondsRat creates a new Timecode based on a rational representation of the seconds count.
//
// seconds will be rounded to the nearest whole-frame based on rate.
func FromSecondsRat(seconds *big.Rat, framerate rate.Framerate) Timecode {
	playback := framerate.Playback()
	playbackDivisor := new(big.Rat).Inv(playback)
	if !new(big.Rat).Mul(seconds, playbackDivisor).IsInt() {
		frames := new(big.Rat).Mul(seconds, playback)
		seconds = new(big.Rat).Mul(frames, playbackDivisor)
	} else {
		// We need to make a copy of the seconds value so it doesn't get changed out from under us by the caller.
		seconds = new(big.Rat).Set(seconds)
	}

	return Timecode{
		seconds: seconds,
		rate:    framerate,
	}
}

// timecodeRegex is the regex we are going to use to parse timecode.
var timecodeRegex = regexp.MustCompile(
	`^(?P<negative>-)?` +
		`((?P<section1>[0-9]+)[:|;])?` +
		`((?P<section2>[0-9]+)[:|;])?` +
		`((?P<section3>[0-9]+)[:|;])?` +
		`(?P<frames>[0-9]+)$`,
)

// Indexes of our submatch groups.
const (
	timecodeRegexNegative = 1
	timecodeRegexSection1 = 3
	timecodeRegexSection2 = 5
	timecodeRegexSection3 = 7
	timecodeRegexFrames   = 8
)

// FromFrames converts a frame count / number to a Timecode value.
func FromFrames(frames int64, framerate rate.Framerate) Timecode {
	playback := framerate.Playback()
	playbackDivisor := playback.Inv(playback)

	framesRat := new(big.Rat).SetInt64(frames)
	seconds := framesRat.Mul(framesRat, playbackDivisor)

	return FromSecondsRat(seconds, framerate)
}

// FromTimecode parses a new timecode value from a string.
func FromTimecode(tc string, framerate rate.Framerate) (Timecode, error) {
	// See if our regex gets a match
	match := timecodeRegex.FindStringSubmatch(tc)
	if match == nil {
		return Timecode{}, ErrFormatNotRecognized
	}

	// The hours, minutes, and seconds place are only optionally present, and annoyingly with the way regex works,
	// will shift what group they match two depending on which ones are present. We need to put them into a
	// slice, and pop them off the end.
	sections := make([]string, 0, 3)
	for _, sectionIndex := range []int{timecodeRegexSection1, timecodeRegexSection2, timecodeRegexSection3} {
		if section := match[sectionIndex]; section != "" {
			sections = append(sections, section)
		}
	}

	var frames int64
	var seconds int64
	var minutes int64
	var hours int64
	if len(sections) >= 1 {
		secondsStr := sections[len(sections)-1]
		seconds, _ = strconv.ParseInt(secondsStr, 10, 64)
	}
	if len(sections) >= 2 {
		minutesStr := sections[len(sections)-2]
		minutes, _ = strconv.ParseInt(minutesStr, 10, 64)
	}
	if len(sections) >= 3 {
		hoursStr := sections[len(sections)-3]
		hours, _ = strconv.ParseInt(hoursStr, 10, 64)
	}

	framesStr := match[timecodeRegexFrames]
	if framesStr != "" {
		frames, _ = strconv.ParseInt(framesStr, 10, 64)
	}

	// Now we need to get the seconds as a rational value so we can multiply it by our timebase.
	seconds = minutes*secondsPerMinute + hours*secondsPerHour + seconds
	secondsRat := big.NewRat(seconds, 1)

	// We are going to calculate our frames as a rational. We multiply our seconds by our timebase then add the
	// frames as a rational value to it.
	framesRat := big.NewRat(frames, 1)
	framesRat = secondsRat.Mul(secondsRat, framerate.Timebase()).Add(secondsRat, framesRat)

	// Then round the result and extract the numerator to get the actual frame count.
	frames = internal.RoundRat(framesRat).Num().Int64()

	// If this was a negative value, we need to make the frames negative.
	isNegative := match[timecodeRegexNegative] != ""
	if isNegative {
		frames = -frames
	}

	// Now we can use our FromFrames conversion.
	return FromFrames(frames, framerate), nil
}

// runtimeRegex is the regex we are going to use to parse runtimes.
var runtimeRegex = regexp.MustCompile(
	`^(?P<negative>-)?((?P<section1>[0-9]+)[:|;])?((?P<section2>[0-9]+)[:|;])?(?P<seconds>[0-9]+(\.[0-9]+)?)$`,
)

// Indexes of our submatch groups.
const (
	runtimeRegexNegative = 1
	runtimeRegexSection1 = 3
	runtimeRegexSection2 = 5
	runtimeRegexSeconds  = 6
)

// FromRuntime parses a new timecode from a runtime string like "01:12:34.342".
func FromRuntime(runtime string, framerate rate.Framerate) (Timecode, error) {
	// See if our regex gets a match
	match := runtimeRegex.FindStringSubmatch(runtime)
	if match == nil {
		return Timecode{}, ErrFormatNotRecognized
	}

	// The hours, and minutes, places are only optionally present, and annoyingly with
	// the way regex works, will shift what group they match two depending on which ones
	// are present. We need to put them into a slice, and pop them off the end.
	sections := make([]string, 0, 2)
	for _, sectionIndex := range []int{runtimeRegexSection1, runtimeRegexSection2} {
		if section := match[sectionIndex]; section != "" {
			sections = append(sections, section)
		}
	}

	var minutes int64
	var hours int64
	if len(sections) >= 1 {
		minutesStr := sections[len(sections)-1]
		minutes, _ = strconv.ParseInt(minutesStr, 10, 64)
	}
	if len(sections) >= 2 {
		hoursStr := sections[len(sections)-2]
		hours, _ = strconv.ParseInt(hoursStr, 10, 64)
	}

	secondsInt := hours*secondsPerHour + minutes*secondsPerMinute

	// This value will always be here if the regex matches, we don't need to check.
	secondsStr := match[runtimeRegexSeconds]
	seconds, _ := new(big.Rat).SetString(secondsStr)

	seconds.Add(seconds, big.NewRat(secondsInt, 1))
	// If this was a negative value, we need to make the frames negative.
	isNegative := match[runtimeRegexNegative] != ""
	if isNegative {
		seconds = seconds.Inv(seconds)
	}

	return FromSecondsRat(seconds, framerate), nil
}

// feetAndFramesRegex will be used to parse our feet and frames value.
var feetAndFramesRegex = regexp.MustCompile(`^(?P<negative>-)?(?P<feet>[0-9]+)\+(?P<frames>[0-9]+)$`)

// Indexes of our submatch groups.
const (
	fafRegexNegative = 1
	fafRegexFeet     = 2
	fafRegexFrames   = 3
)

// FromFeetAndFrames parses a timecode from a feet+frames string like
func FromFeetAndFrames(faf string, framerate rate.Framerate) (Timecode, error) {
	// See if our regex gets a match
	match := feetAndFramesRegex.FindStringSubmatch(faf)
	if match == nil {
		return Timecode{}, ErrFormatNotRecognized
	}

	feet, _ := strconv.ParseInt(match[fafRegexFeet], 10, 64)
	frames, _ := strconv.ParseInt(match[fafRegexFrames], 10, 64)

	frames += feet * framesPerFoot
	// If this was a negative value, we need to make the frames negative.
	isNegative := match[fafRegexNegative] != ""
	if isNegative {
		frames = -frames
	}

	return FromFrames(frames, framerate), nil
}

// ticksDivisor holds an inverted version of premiereTicksPerSecondsRat for dividing.
var ticksDivisor = new(big.Rat).Inv(premiereTicksPerSecondsRat)

// FromPremiereTicks returns a new Timecode from a number of Adobe Premiere Pro Ticks.
//
// The resulting timecode will be rounded to the nearest whole-frame, given framerate.
func FromPremiereTicks(ticks int64, framerate rate.Framerate) Timecode {
	ticksRat := big.NewRat(ticks, 1)
	secondsRat := ticksRat.Mul(ticksRat, ticksDivisor)
	return FromSecondsRat(secondsRat, framerate)
}
