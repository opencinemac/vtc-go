package rate

import (
	"fmt"
	"github.com/opencinemac/vtc-go/pkg/internal"
	"math/big"
)

// Framerate is the rate at which a video file frames are played back.
//
// Framerate is measured in frames-per-second (24000/1001 = 23.98 frames-per-second).
type Framerate struct {
	playback *big.Rat
	ntsc     NTSC
}

// String implements fmt.Stringer.
func (rate Framerate) String() string {
	rateFloat, _ := rate.playback.Float64()
	var floatString string
	// If this playback is an int, we don't need to to show any places after the 0, and can just truncate the
	// float.
	if rate.playback.IsInt() {
		floatString = fmt.Sprintf("%.0f", rateFloat)
	} else {
		// Otherwise round it to 2 places.
		floatString = fmt.Sprintf("%.2f", rateFloat)
	}

	return fmt.Sprintf("%v %v", floatString, rate.ntsc)
}

// NTSC returns if and which type of NTSC standard this framerate adheres to.
func (rate Framerate) NTSC() NTSC {
	return rate.ntsc
}

// Playback returns the real-world playback speed of the Framerate in frames-per-second.
func (rate Framerate) Playback() *big.Rat {
	// we need to return a copy of the inner value so it does not get modified
	return new(big.Rat).Set(rate.playback)
}

// Timebase returns the speed at which timecode is interpreted at in frames-per-second.
func (rate Framerate) Timebase() *big.Rat {
	if rate.ntsc == NTSCNone {
		return rate.Playback()
	}

	// big.Rat does not offer a rounding method, so we need to convert it to a float,
	// round it, convert it to an int, then feed it to a new big.Rat.
	return internal.RoundRat(rate.Playback())
}

// mustNew panics if a new Framerate could not be created. Used for creating our framerate constants.
func mustNew(timebase int64, ntsc NTSC) Framerate {
	rate, err := FromInt(timebase, ntsc)
	if err != nil {
		panic(fmt.Errorf("error creating Framerate from %v: %w", timebase, err))
	}
	return rate
}

// revive:disable

// vtc comes with a number of pre-defined, common Framerate values.
var (
	// F23_98 is a 23.98 NTSC, Non-Drop framerate.
	F23_98 = mustNew(24, NTSCNonDrop)

	// F24 is a 24 fps, Non-NTSC framerate.
	F24 = mustNew(24, NTSCNone)

	// F29_97Ndf is a 29.97 NTSC, Non-Drop framerate.
	F29_97Ndf = mustNew(30, NTSCNonDrop)

	// F29_97Df is a 29.97 NTSC, Drop-Frame framerate.
	F29_97Df = mustNew(30, NTSCDrop)

	// F30 is a 30 fps, Non-NTSC framerate.
	F30 = mustNew(30, NTSCNone)

	// F47_95 is a 57.95 NTSC, Non-Drop framerate.
	F47_95 = mustNew(48, NTSCNonDrop)

	// F48 is a 48 fps, Non-NTSC framerate.
	F48 = mustNew(48, NTSCNone)

	// F59_94Ndf is a 59.94 NTSC, Non-Drop framerate.
	F59_94Ndf = mustNew(60, NTSCNonDrop)

	// F59_94Df is a 59.94 NTSC, Drop-Frame framerate.
	F59_94Df = mustNew(60, NTSCDrop)

	// F60 is a 60 fps, Non-NTSC framerate.
	F60 = mustNew(60, NTSCNone)
)

// revive:enable
