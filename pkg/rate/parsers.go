package rate

import (
	"fmt"
	"github.com/opencinemac/vtc-go/pkg/internal"
	"math/big"
	"strconv"
)

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
// • Integer: '24'
//
// • Float: '23.98'
//
// • Rational: '24/1'
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
