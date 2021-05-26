package tc_test

import (
	"fmt"
	"github.com/opencinemac/vtc-go/pkg/rate"
	"github.com/opencinemac/vtc-go/pkg/tc"
)

// Basic comparison
func ExampleTimecode_Cmp() {
	tc1, _ := tc.FromTimecode("01:00:00:00", rate.F23_98)
	tc2, _ := tc.FromTimecode("02:00:00:00", rate.F23_98)

	// The returned type is a new type of int, but implements fmt.Stringer
	fmt.Println(tc1.Cmp(tc2))

	// Output:
	// LT
}

// Comparisons are done based on the real-world seconds elapsed for a timecode.
//
// Both the timecodes below have a representation of 01:00:00:00, but the timecode
// running at an NTSC rate will take slightly longer to reach that value than the
// timecode running at true 24 fps.
//
// We can see this by checking their runtimes.
func ExampleTimecode_Cmp_mixedRates() {
	tc1, _ := tc.FromTimecode("01:00:00:00", rate.F23_98)
	tc2, _ := tc.FromTimecode("01:00:00:00", rate.F24)

	fmt.Println(tc1.Cmp(tc2))
	fmt.Println("RUNTIME 1:", tc1.Runtime(3))
	fmt.Println("RUNTIME 2:", tc2.Runtime(3))

	// Output:
	// GT
	// RUNTIME 1: 01:00:03.6
	// RUNTIME 2: 01:00:00.0
}

func ExampleFromTimecode_basic() {
	timecode, _ := tc.FromTimecode("01:00:00:00", rate.F23_98)

	fmt.Println(timecode)

	// Output:
	// 01:00:00:00 @ 23.98 NTSC NDF
}

// We can parse partial timecodes with no problems.
func ExampleFromTimecode_partial() {
	timecode, _ := tc.FromTimecode("1:12", rate.F23_98)

	fmt.Println(timecode)

	// Output:
	// 00:00:01:12 @ 23.98 NTSC NDF
}

// We can also parse timecodes where a value has overflowed it's place.
func ExampleFromTimecode_overflow() {
	// The max value that SHOULD appear in the frames place is '23'.
	timecode, _ := tc.FromTimecode("00:00:00:48", rate.F23_98)

	// But that's okay
	fmt.Println(timecode)

	// Output:
	// 00:00:02:00 @ 23.98 NTSC NDF
}
