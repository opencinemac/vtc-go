package tc

import (
	"errors"
	"fmt"
)

// ErrParseTimecode is the sentinel error returned when a timecode could not be parsed.
var ErrParseTimecode = errors.New("could not parse Timecode")

// ErrFormatNotRecognized is the sentinel error returned when a timecode string's format is not recognized.
var ErrFormatNotRecognized = fmt.Errorf("%w: string format not recognized", ErrParseTimecode)
