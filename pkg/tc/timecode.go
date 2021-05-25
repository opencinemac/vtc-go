package tc

import (
	"fmt"
	"github.com/opencinemac/vtc-go/pkg/internal"
	"github.com/opencinemac/vtc-go/pkg/rate"
	"github.com/wadey/go-rounding"
	"math/big"
	"strings"
)

// TimecodeSections holds the individual sections of a timecode for formatting /
// manipulation.
type TimecodeSections struct {
	// IsNegative is whether the timecode is less than 0.
	IsNegative bool
	// Hours is the value of the hours place.
	Hours int64
	// Minutes is the value of the hours place.
	Minutes int64
	// Seconds is the value of the hours place.
	Seconds int64
	// Frames is the value of the hours place.
	Frames int64
}

// Timecode represents the frame at a particular time in a video.
type Timecode struct {
	// seconds holds the rational representation of the real-world number of seconds.
	seconds *big.Rat
	// rate holds information about our framerate.
	rate rate.Framerate
}

// Rate returns the rate.Framerate of the timecode.
func (tc Timecode) Rate() rate.Framerate {
	return tc.rate
}

// IsNegative returns true if the value is negative.
func (tc Timecode) IsNegative() bool {
	return tc.seconds.Cmp(zeroRat) == -1
}

/*
Seconds returns the rational representation of the real-world seconds that would have
elapsed between 00:00:00:00 and this timecode.

What it is

The number of real-world seconds that have elapsed between 00:00:00:00 and the timecode
value. With NTSC timecode, the timecode drifts from the real-world elapsed time.

Where you see it

- Anywhere real-world time needs to be calculated.
- In code that needs to do lossless calculations of playback time not rely on frame
  count, like adding two timecodes together with different framerates.
*/
func (tc Timecode) Seconds() *big.Rat {
	return new(big.Rat).Set(tc.seconds)
}

// Sections returns the individual sections of a timecode string as int values.
//
// Note: this method will panic on framerates where the timebase is not a whole integer.
func (tc Timecode) Sections() TimecodeSections {
	timebase := tc.Rate().Timebase()

	framesInt := tc.Frames()
	isNegative := tc.IsNegative()
	if isNegative {
		framesInt = -framesInt
	}
	if tc.rate.NTSC() == rate.NTSCDrop {
		framesInt += dropFrameNumAdjustment(framesInt, tc.rate)
	}

	frames := big.NewRat(framesInt, 1)

	framesPerMinute := new(big.Rat).Mul(secondsPerMinuteRat, timebase)
	framesPerHour := new(big.Rat).Mul(secondsPerHourRat, timebase)

	hours, frames := internal.DivModRat(frames, framesPerHour)
	minutes, frames := internal.DivModRat(frames, framesPerMinute)

	seconds, frames := internal.DivModRat(frames, timebase)
	frames = internal.RoundRat(frames)

	return TimecodeSections{
		IsNegative: isNegative,
		Hours:      hours.Num().Int64(),
		Minutes:    minutes.Num().Int64(),
		Seconds:    seconds.Num().Int64(),
		Frames:     frames.Num().Int64(),
	}
}

/*
Timecode returns the the formatted SMPTE timecode: (ex: 01:00:00:00).

What it is

Timecode is used as a human-readable way to represent the id of a given frame. It is
formatted to give a rough sense of where to find a frame:
{HOURS}:{MINUTES}:{SECONDS}:{FRAME}. For more on timecode, see Frame.io's
[excellent post](https://blog.frame.io/2017/07/17/timecode-and-frame-rates/) on the
subject.

Where you see it

Timecode is ubiquitous in video editing, a small sample of places you might see
timecode:

- Source and Playback monitors in your favorite NLE.
- Burned into the footage for dailies.
- Cut lists like an EDL.

Warning

Currently, this method will panic on framerates where the timebase is not a whole
integer.
*/
func (tc Timecode) Timecode() string {
	sections := tc.Sections()

	// We'll add a negative sign if the timecode is negative.
	sign := ""
	if sections.IsNegative {
		sign = "-"
	}

	// If this is a drop-frame timecode, we need to use a ';' to separate the frames
	// from the seconds.
	frameSep := ":"
	if tc.Rate().NTSC() == rate.NTSCDrop {
		frameSep = ";"
	}

	return fmt.Sprintf(
		"%v%02d:%02d:%02d%v%02d",
		sign,
		sections.Hours,
		sections.Minutes,
		sections.Seconds,
		frameSep,
		sections.Frames,
	)
}

/*
Frames returns the number of frames that would have elapsed between 00:00:00:00 and this
timecode.

What it is

Frame number / frames count is the number of a frame if the timecode started at
00:00:00:00 and had been running until the current value. A timecode of '00:00:00:10'
has a frame number of 10. A timecode of '01:00:00:00' has a frame number of 86400.

Where you see it

- Frame-sequence files: 'my_vfx_shot.0086400.exr'
- FCP7XML cut lists:

	```xml
	<timecode>
		<rate>
			<timebase>24</timebase>
			<ntsc>TRUE</ntsc>
		</rate>
		<string>01:00:00:00</string>
		<frame>86400</frame>  <!-- <====THIS LINE-->
		<displayformat>NDF</displayformat>
	</timecode>
	```
*/
func (tc Timecode) Frames() int64 {
	playback := tc.rate.Playback()

	// Get the frames by multiplying our seconds by the playback speed, then rounding
	// the result.
	frames := playback.Mul(tc.seconds, playback)
	frames = internal.RoundRat(frames)

	// Once the rational value is rounded, return the numerator.
	return frames.Num().Int64()
}

// rat10 is used to check if we need to add a leading zero to the seconds place of the
// runtime by checking whether the seconds value is less than 10.
var rat10 = big.NewRat(10, 1)

/*
Runtime Returns the true, real-world runtime of the timecode in HH:MM:SS.FFFFFFFFF
format.

What it is

The formatted version of seconds. It looks like timecode, but with a decimal seconds
value instead of a frame number place.

Where you see it

- Anywhere real-world time is used.
- FFMPEG commands:

   ```shell
   ffmpeg -ss 00:00:30.5 -i input.mov -t 00:00:10.25 output.mp4
   ```
*/
func (tc Timecode) Runtime(precision int) string {
	seconds := tc.Seconds()
	// If this is a negative value, make it positive for the purposes of parsing the
	// value.
	isNegative := tc.IsNegative()
	if isNegative {
		seconds.Neg(seconds)
	}

	hours, seconds := internal.DivModRat(seconds, secondsPerHourRat)
	minutes, seconds := internal.DivModRat(seconds, secondsPerMinuteRat)

	seconds = rounding.Round(seconds, precision, rounding.HalfUp)

	var secondsStr string
	if seconds.IsInt() {
		secondsStr = fmt.Sprintf("%02d.0", seconds.Num().Int64())
	} else {
		secondsStr = seconds.FloatString(precision)
		// Trim any trailing zeros.
		secondsStr = strings.TrimRight(secondsStr, "0")
		// If the seconds is less than 10, we need to pad a leading 0.
		if seconds.Cmp(rat10) == -1 {
			secondsStr = "0" + secondsStr
		}
	}

	sign := ""
	if isNegative {
		sign = "-"
	}

	return fmt.Sprintf("%v%02d:%02d:%v", sign, hours.Num().Int64(), minutes.Num().Int64(), secondsStr)
}

/*
FeetAndFrames returns the number of feet and frames this timecode represents if it were
shot on 35mm 4-perf film (16 frames per foot). ex: '5400+13'.

What it is

On physical film, each foot contains a certain number of frames. For 35mm, 4-perf film
(the most common type on Hollywood movies), this number is 16 frames per foot.
Feet-And-Frames was often used in place of Keycode to quickly reference a frame in the
edit.

Where you see it

For the most part, feet + frames has died out as a reference, because digital media is
not measured in feet. The most common place it is still used is Studio Sound
Departments. Many Sound Mixers and Designers intuitively think in feet + frames, and it
is often burned into the reference picture for them.

- Telecine.
- Sound turnover reference picture.
- Sound turnover change lists.
*/
func (tc Timecode) FeetAndFrames() string {
	frames := tc.Frames()
	// If this is a negative value, make it positive.
	isNegative := tc.IsNegative()
	if isNegative {
		frames = -frames
	}

	feet := frames / framesPerFoot
	frames = frames % framesPerFoot

	sign := ""
	if isNegative {
		sign = "-"
	}

	return fmt.Sprintf("%v%v+%02d", sign, feet, frames)
}

/*
PremiereTicks returns the number of elapsed ticks this timecode represents in Adobe
Premiere Pro.

What it is

Internally, Adobe Premiere Pro uses ticks to divide up a second, and keep track of how
far into that second we are. There are 254016000000 ticks in a second, regardless of
framerate in Premiere.

Where you see it

- Premiere Pro Panel functions and scripts
- FCP7XML cutlists generated from Premiere:

	```xml
	<clipitem id="clipitem-1">
		...
		<in>158</in>
		<out>1102</out>
		<pproTicksIn>1673944272000</pproTicksIn>
		<pproTicksOut>11675231568000</pproTicksOut>
		...
	</clipitem>
	```
*/
func (tc Timecode) PremiereTicks() int64 {
	seconds := tc.Seconds()
	ticks := seconds.Mul(seconds, premiereTicksPerSecondsRat)
	return internal.RoundRat(ticks).Num().Int64()
}
