package tc

import (
	"errors"
	"fmt"
)

// tc comes with sentinel errors for catching various parsing errors.
var (
	// ErrParseTimecode is the sentinel error returned when a timecode could not be
	// parsed. All other parsing errors wrap this error.
	ErrParseTimecode = errors.New("could not parse Timecode")

	// ErrFormatNotRecognized is the sentinel error returned when a timecode string's
	// format is not recognized.
	ErrFormatNotRecognized = fmt.Errorf("%w: string format not recognized", ErrParseTimecode)

	// ErrBadDropFrameValue is returned when a timecode to be parsed includes a
	// disallowed frame value. Ex: ('00:01:00:01', since this frame should be dropped.)
	ErrBadDropFrameValue = fmt.Errorf(
		"%w: frames value not allowed in Drop-Frame timecode", ErrParseTimecode,
	)
)
