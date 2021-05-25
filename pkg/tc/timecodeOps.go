package tc

import (
	"github.com/opencinemac/vtc-go/pkg/internal"
	"github.com/opencinemac/vtc-go/pkg/rate"
	"github.com/wadey/go-rounding"
	"math/big"
)

// Add adds two timecodes together.
func (tc Timecode) Add(other Timecode) Timecode {
	seconds := tc.Seconds()
	seconds = seconds.Add(seconds, other.seconds)

	return Timecode{
		seconds: seconds,
		rate:    tc.rate,
	}
}

// Sub subtracts a timecode from the caller.
func (tc Timecode) Sub(other Timecode) Timecode {
	seconds := tc.Seconds()
	seconds = seconds.Sub(seconds, other.seconds)

	return Timecode{
		seconds: seconds,
		rate:    tc.rate,
	}
}

// Mul multiplies a timecode by a scalar.
func (tc Timecode) Mul(multiplier *big.Rat) Timecode {
	seconds := tc.Seconds()
	seconds = seconds.Mul(seconds, multiplier)

	return Timecode{
		seconds: seconds,
		rate:    tc.rate,
	}
}

// Div divides a timecode by a scalar. Divide returns a result as if floor division had
// been done to the frame count.
func (tc Timecode) Div(divisor *big.Rat) Timecode {
	divisor = new(big.Rat).Inv(divisor)

	frames := big.NewRat(tc.Frames(), 1)
	frames = frames.Mul(frames, divisor)

	frames = rounding.Round(frames, 0, rounding.Down)
	return FromFrames(frames.Num().Int64(), tc.Rate())
}

// Mod divides a timecode by a scalar and returns the dividend and remainder. Mod
// returns a result as if floor division had been done to the frame count.
func (tc Timecode) Mod(divisor *big.Rat) Timecode {
	_, mod := tc.DivMod(divisor)
	return mod
}

// DivMod divides a timecode by a scalar and returns the dividend and remainder.
// DivMod returns a result as if floor division had been done to the frame count.
func (tc Timecode) DivMod(divisor *big.Rat) (dividend Timecode, remainder Timecode) {
	divisor = new(big.Rat).Inv(divisor)

	frames := big.NewRat(tc.Frames(), 1)
	frames = frames.Mul(frames, divisor)

	dividendRat, remainderRat := internal.DivModRat(frames, divisor)
	remainderRat = internal.RoundRat(remainderRat)
	return FromFrames(dividendRat.Num().Int64(), tc.rate), FromFrames(remainderRat.Num().Int64(), tc.rate)
}

// Neg returns the negative version of the timecode (will be positive if current value
// is negative).
func (tc Timecode) Neg() Timecode {
	seconds := tc.Seconds()
	seconds = seconds.Neg(seconds)

	return Timecode{
		seconds: seconds,
		rate:    tc.rate,
	}
}

// Abs returns the absolute value of the Timecode.
func (tc Timecode) Abs() Timecode {
	if !tc.IsNegative() {
		return tc
	}

	return tc.Neg()
}

// Rebase returns a Timecode with the same number of frames running at a different
// Framerate.
func (tc Timecode) Rebase(framerate rate.Framerate) Timecode {
	// Just use the FromFrames parser to parse the frame count at the new framerate.
	return FromFrames(tc.Frames(), framerate)
}
