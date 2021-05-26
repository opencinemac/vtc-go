package rate

import (
	"errors"
	"fmt"
)

// rate comes with a number of sentinel errors for catching framerate parsing problems.
var (
	// ErrParseFramerate is returned when there is an error parsing a Framerate, and
	// is wrapped by all other errors.
	ErrParseFramerate = errors.New("could not parse Framerate")

	// ErrBadDropFrameRate is returned when a drop-frame playback rate is not cleanly
	// divisible by 30000/1001 (29.97 NTSC) For more information on why drop-frame
	// framerates must be a clean multiple of 29.97 NTSC, see:
	// https://www.davidheidelberger.com/2010/06/10/drop-frame-timecode/
	ErrBadDropFrameRate = fmt.Errorf(
		"%w: drop-frame Framerate values must have a playback cleanly divisible by "+
			"30000/1001",
		ErrParseFramerate,
	)

	// ErrImprecise is returned when a float is being parsed as part of a non-NTSC
	// framerate. Without knowing how to coerce the float to a sane framerate, the values
	// are not precise enough to yield a meaningful value.
	ErrImprecise = fmt.Errorf(
		"%w: non-ntsc framerates cannot be parsed from floats due to imprecision",
		ErrParseFramerate,
	)

	// ErrNegative is returned when a negative value is passed to a Framerate parser.
	// Negative framerates are not supported.
	ErrNegative = fmt.Errorf("%w: Framerate cannot be negative", ErrParseFramerate)

	// ErrBadNtsc is returned when an enum value outside the predefined NTSC constant values
	// is passed into a Framerate parser.
	ErrBadNtsc = fmt.Errorf("%w: NTSC value not recognized", ErrParseFramerate)
)
