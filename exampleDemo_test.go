package vtc_test

import (
	"fmt"
	"github.com/opencinemac/vtc-go/pkg/rate"
	"github.com/opencinemac/vtc-go/pkg/tc"
	"math/big"
)

func Example() {
	// It's easy to make a new 23.98 NTSC Timecode.
	timecode, err := tc.FromTimecode("01:00:00:00", rate.F23_98)
	if err != nil {
		panic(fmt.Errorf("error parsing timecode: %w", err))
	}

	// 01:00:00:00 @ 23.98 NTSC NDF
	fmt.Println(timecode)

	// We can get all sorts of ways to represent the timecode.

	// 01:00:00:00
	fmt.Println(timecode.Timecode())

	// 86400
	fmt.Println(timecode.Frames())

	// 18018/5
	fmt.Println(timecode.Seconds())

	// 01:00:03.6
	fmt.Println(timecode.Runtime(9))

	// 915372057600000
	fmt.Println(timecode.PremiereTicks())

	// 5400+00
	fmt.Println(timecode.FeetAndFrames())

	// We can inspect the framerate:

	// 24000/1001
	fmt.Println(timecode.Rate().Playback())

	// 24/1
	fmt.Println(timecode.Rate().Timebase())

	// NTSC NDF
	fmt.Println(timecode.Rate().NTSC())

	// Parsing is flexible:

	// Partial timecodes:
	timecode, _ = tc.FromTimecode("3:12", rate.F23_98)

	// 00:00:03:12 @ 23.98 NTSC NDF
	fmt.Println(timecode)

	// Frames:
	timecode = tc.FromFrames(24, rate.F23_98)

	// 00:00:01:00 @ 23.98 NTSC NDF
	fmt.Println(timecode)

	// Seconds:
	timecode = tc.FromSeconds(big.NewRat(3, 2), rate.F23_98)

	// 00:00:01:12 @ 23.98 NTSC NDF
	fmt.Println(timecode)

	// Premiere Ticks:
	timecode = tc.FromPremiereTicks(254016000000, rate.F23_98)

	// 00:00:01:00 @ 23.98 NTSC NDF
	fmt.Println(timecode)

	// Runtime:
	timecode, _ = tc.FromRuntime("01:00:00.5", rate.F23_98)

	// 00:59:56:22 @ 23.98 NTSC NDF
	fmt.Println(timecode)

	// Feet + Frames:
	timecode, _ = tc.FromFeetAndFrames("213+07", rate.F23_98)

	// 00:02:22:07 @ 23.98 NTSC NDF
	fmt.Println(timecode)

	// We can add two timecodes:
	timecode, _ = tc.FromTimecode("17:23:13:02", rate.F23_98)
	other, _ := tc.FromTimecode("01:00:00:00", rate.F23_98)

	timecode = timecode.Add(other)

	// 18:23:13:02 @ 23.98 NTSC NDF
	fmt.Println(timecode)

	// We can subtract too.
	timecode = timecode.Sub(other)

	// 17:23:13:02 @ 23.98 NTSC NDF
	fmt.Println(timecode)

	// It's easy to compare two timecodes:

	// GT (1)
	fmt.Println(timecode.Cmp(other))

	// We can multiply:
	timecode = timecode.Mul(big.NewRat(2, 1))

	// 34:46:26:04 @ 23.98 NTSC NDF
	fmt.Println(timecode)

	// ... divide...
	timecode = timecode.Div(big.NewRat(2, 1))

	// 17:23:13:02 @ 23.98 NTSC NDF
	fmt.Println(timecode)

	// and even get the remainder of division
	dividend, remainder := timecode.DivMod(big.NewRat(3, 2))

	// DIVIDEND: 11:35:28:17 @ 23.98 NTSC NDF
	// REMAINDER: 00:00:00:01 @ 23.98 NTSC NDF
	fmt.Println("DIVIDEND:", dividend)
	fmt.Println("REMAINDER:", remainder)

	// We can make a timecode negative:
	timecode = timecode.Neg()

	// -17:23:13:02 @ 23.98 NTSC NDF
	fmt.Println(timecode)

	// Or get it's absolute value:
	timecode = timecode.Abs()

	// 17:23:13:02 @ 23.98 NTSC NDF
	fmt.Println(timecode)

	// We can make dropframe timecode for 29.97 or 59.94 using one of the pre-set
	// framerates.
	dropFrame := tc.FromFrames(15000, rate.F29_97Df)

	// 00:08:20;18 @ 29.97 NTSC DF
	// NTSC DF
	fmt.Println(dropFrame)
	fmt.Println(dropFrame.Rate().NTSC())

	// We can make new timecodes with arbitrary framerates if we want:
	framerate, _ := rate.FromInt(137, rate.NTSCNone)
	timecode, _ = tc.FromTimecode("01:00:00:00", framerate)

	// 493200
	fmt.Println(timecode.Frames())

	// We can make NTSC values for timebases and playback speeds that do not ship with
	// this crate:
	framerate, _ = rate.FromInt(120, rate.NTSCNonDrop)
	timecode, _ = tc.FromTimecode("01:00:00:00", framerate)

	// 01:00:00:00 @ 119.88 NTSC NDF
	fmt.Println(timecode)

	// We can also rebase them using another framerate:
	timecode = timecode.Rebase(rate.F59_94Ndf)

	// 02:00:00:00 @ 59.94 NTSC NDF
	fmt.Println(timecode)

	// Output:
	// 01:00:00:00 @ 23.98 NTSC NDF
	// 01:00:00:00
	// 86400
	// 18018/5
	// 01:00:03.6
	// 915372057600000
	// 5400+00
	// 24000/1001
	// 24/1
	// NTSC NDF
	// 00:00:03:12 @ 23.98 NTSC NDF
	// 00:00:01:00 @ 23.98 NTSC NDF
	// 00:00:01:12 @ 23.98 NTSC NDF
	// 00:00:01:00 @ 23.98 NTSC NDF
	// 00:59:56:22 @ 23.98 NTSC NDF
	// 00:02:22:07 @ 23.98 NTSC NDF
	// 18:23:13:02 @ 23.98 NTSC NDF
	// 17:23:13:02 @ 23.98 NTSC NDF
	// GT
	// 34:46:26:04 @ 23.98 NTSC NDF
	// 17:23:13:02 @ 23.98 NTSC NDF
	// DIVIDEND: 11:35:28:17 @ 23.98 NTSC NDF
	// REMAINDER: 00:00:00:01 @ 23.98 NTSC NDF
	// -17:23:13:02 @ 23.98 NTSC NDF
	// 17:23:13:02 @ 23.98 NTSC NDF
	// 00:08:20;18 @ 29.97 NTSC DF
	// NTSC DF
	// 493200
	// 01:00:00:00 @ 119.88 NTSC NDF
	// 02:00:00:00 @ 59.94 NTSC NDF
}
