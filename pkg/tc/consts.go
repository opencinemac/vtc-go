package tc

import "math/big"

// cache some constants for seconds per interval
const (
	secondsPerMinute int64 = 60
	secondsPerHour         = secondsPerMinute * 60
)

var (
	secondsPerMinuteRat = big.NewRat(secondsPerMinute, 1)
	secondsPerHourRat   = big.NewRat(secondsPerHour, 1)
)

// premiereTicksPerSecond is the number of ticks Adobe Premiere Pro tracks within a
// seconds.
const premiereTicksPerSecond int64 = 254016000000

// framesPerFoot is the number of frames in a foot of 35mm, 4-perf film.
const framesPerFoot int64 = 16

// premiereTicksPerSecondsRat is the rational version of premiereTicksPerSecond.
var premiereTicksPerSecondsRat = big.NewRat(premiereTicksPerSecond, 1)

// zeroRat will be used to check if rational values are negative.
var zeroRat = big.NewRat(0, 1)
