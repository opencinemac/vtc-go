package rate

import (
	"fmt"
	"github.com/opencinemac/vtc-go/pkg/internal"
	"math/big"
	"strconv"
)

// NTSC is an enum-like type for specifying whether a framerate adheres to the NTSC standard.
type NTSC int

const (
	// NTSCNone means the framerate is not an NTSC framerate.
	NTSCNone NTSC = iota
	// NTSCNonDrop means this in an NTSC, Non-drop-frame framerate.
	NTSCNonDrop
	// NTSCDrop means this in an NTSC, drop-frame framerate.
	NTSCDrop
)

// String implements fmt.Stringer.
func (ntsc NTSC) String() string {
	switch ntsc {
	case NTSCNone:
		return "fps"
	case NTSCNonDrop:
		return "NTSC NDF"
	case NTSCDrop:
		return "NTSC DF"
	default:
		return "[INVALID NTSC VALUE]"
	}
}

// IsNTSC returns whether this value represents an NTSC standard.
//
// Returns false for NTSCNone
//
// Returns true for NTSCNonDrop and NTSCDrop
func (ntsc NTSC) IsNTSC() bool {
	return ntsc == NTSCNonDrop || ntsc == NTSCDrop
}

// Validate returns ErrBadNtsc if this value is not one of the pre-defined NTSC enum constants that ships with
// this library.
func (ntsc NTSC) Validate() error {
	if ntsc < 0 || ntsc > NTSCDrop {
		return ErrBadNtsc
	}
	return nil
}

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

// FromRat creates a Framerate from a *big.Rat value.
//
// If ntsc != NTSCNone, value will be coerced to the nearest valid NTSC framerate by
// rounding to the nearest whole number, and putting over a denominator of 1001.
//
// IsNegative values are currently disallowed and will return an error.
//
// If ntsc == NTSCDrop and the value is not cleanly divisible by 30000/1001 after the
// normal NTSC coercion, an error will be returned. Drop-frame values must be a clean
// multiple of a 29.97 NTSC framerate. For more information on why this is, see:
// https://www.davidheidelberger.com/2010/06/10/drop-frame-timecode/
func FromRat(value *big.Rat, ntsc NTSC) (Framerate, error) {
	// Return an error if this is not a valid NTSC value.
	if err := ntsc.Validate(); err != nil {
		return Framerate{}, err
	}

	// Return an error if this value is negative. We do not allow negative framerates at this time.
	if value.Cmp(zeroRat) == -1 {
		return Framerate{}, ErrNegative
	}

	// If this is an NTSC value, but the denominator is not 1001, then return an error.
	if ntsc.IsNTSC() && value.Denom().Int64() != 1001 {
		timebase := internal.RoundRat(value)
		value = new(big.Rat).SetFrac64(timebase.Num().Int64()*1000, 1001)
	} else {
		// We want to make a copy so this rate does not get swapped out from under us by the caller by mistake.
		value = new(big.Rat).Set(value)
	}

	// If this is a drop-frame value, but is not cleanly divisible by 30000/1001 (29.97 NTSC), return an error. For
	// more info on why drop-frame must divisible by 29.97, see:
	// https://www.davidheidelberger.com/2010/06/10/drop-frame-timecode/
	if ntsc == NTSCDrop && !new(big.Rat).Mul(value, dropFrameDivisor).IsInt() {
		return Framerate{}, ErrBadDropFrameRate
	}

	return Framerate{
		playback: value,
		ntsc:     ntsc,
	}, nil
}

// FromFloat creates a new Framerate from a float64 value.
//
// If ntsc == NTSCNone, this method will fail. Floating point values are not precise enough to be yield expected
// results if not being coerced to the nearest valid NTSC value.
func FromFloat(value float64, ntsc NTSC) (Framerate, error) {
	if !ntsc.IsNTSC() {
		return Framerate{}, ErrImprecise
	}

	return FromRat(new(big.Rat).SetFloat64(value), ntsc)
}

// FromInt converts an integer value into a Framerate.
//
// If ntsc != NTSCNone, the value wil be assumed to be a timebase, and will be multiplied by 1000 and put over a
// denominator of 1001.
func FromInt(value int64, ntsc NTSC) (Framerate, error) {
	return FromRat(new(big.Rat).SetInt64(value), ntsc)
}

// FromString creates a new Framerate from a string value, which may have the following forms:
//
// - Integer: '24'
// - Float: '23.98'
// - Rational: '24/1'
//
// Float representations are only allowed if ntsc != NTSCNone.
//
// It's common for some metadata reporting tools to represent the Timebase as '1/24' instead of '24/1'. If the
// numerator of a parsed rational is found to be smaller than the denominator, the rational will be inverted.
func FromString(value string, ntsc NTSC) (Framerate, error) {
	intVal, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		return FromInt(intVal, ntsc)
	}

	floatVal, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return FromFloat(floatVal, ntsc)
	}

	ratValue, ok := new(big.Rat).SetString(value)
	if !ok {
		return Framerate{}, fmt.Errorf(
			"%w: string format not recognized. must be int, float, or rational", ErrParseFramerate,
		)
	}

	// If the numerator is less than the denominator, this is likely a string where the values were flipped. Some
	// programs print 1/24 instead of 24/1. We will invert the value if this is the case.
	if ratValue.Num().Cmp(ratValue.Denom()) == -1 {
		ratValue.Inv(ratValue)
	}

	return FromRat(ratValue, ntsc)
}

// dropFrameDivisor is used to test whether a playback value is a valid drop-frame playback rate. Drop frame must
// be divisible by 30000/1001 (29.97 NTSC). If the result of multiplying an incoming playback value by this value
// is not an integer, then we should return an error
//
// For more information on why drop-frame value must be divisible by 29.97, see:
// https://www.davidheidelberger.com/2010/06/10/drop-frame-timecode/
var dropFrameDivisor = big.NewRat(1001, 30000)

// zeroRat will be used to check whether an incoming framerate value is negative.
var zeroRat = big.NewRat(0, 1)

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
