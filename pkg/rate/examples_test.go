package rate_test

import (
	"fmt"
	"github.com/opencinemac/vtc-go/pkg/rate"
	"math/big"
)

// Parse a Framerate from a *big.Rat NTSC playback speed
func ExampleFromRat_ntscPlayback() {
	framerate, err := rate.FromRat(big.NewRat(24000, 1001), rate.NTSCNonDrop)
	if err != nil {
		panic(err)
	}

	fmt.Println("RESULT:", framerate)

	// Output:
	// RESULT: 23.98 NTSC NDF
}

// Parse a Framerate from a *big.Rat NTSC timebase
func ExampleFromRat_ntscTimebase() {
	framerate, err := rate.FromRat(big.NewRat(24, 1), rate.NTSCNonDrop)
	if err != nil {
		panic(err)
	}

	fmt.Println("RESULT:", framerate)

	// Output:
	// RESULT: 23.98 NTSC NDF
}

// Parse a Framerate from a *big.Rat non-NTSC playback speed
func ExampleFromRat_fps() {
	framerate, err := rate.FromRat(big.NewRat(24, 1), rate.NTSCNone)
	if err != nil {
		panic(err)
	}

	fmt.Println("RESULT:", framerate)

	// Output:
	// RESULT: 24 fps
}

// Parse a Framerate from a int64 as an ntsc timebase
func ExampleFromInt_ntsc() {
	framerate, err := rate.FromInt(24, rate.NTSCNonDrop)
	if err != nil {
		panic(err)
	}

	fmt.Println("RESULT:", framerate)

	// Output:
	// RESULT: 23.98 NTSC NDF
}

// Parse a Framerate from a int64 as a non-ntsc framerate
func ExampleFromInt_fps() {
	framerate, err := rate.FromInt(24, rate.NTSCNone)
	if err != nil {
		panic(err)
	}

	fmt.Println("RESULT:", framerate)

	// Output:
	// RESULT: 24 fps
}

// Parse a Framerate from a float64 as a NTSC playback speed.
func ExampleFromFloat_ntscPlayback() {
	framerate, err := rate.FromFloat(23.98, rate.NTSCNonDrop)
	if err != nil {
		panic(err)
	}

	fmt.Println("RESULT:", framerate)

	// Output:
	// RESULT: 23.98 NTSC NDF
}

// Parse a Framerate from a float64 as a NTSC timebase.
func ExampleFromFloat_ntscTimebase() {
	framerate, err := rate.FromFloat(24.0, rate.NTSCNonDrop)
	if err != nil {
		panic(err)
	}

	fmt.Println("RESULT:", framerate)

	// Output:
	// RESULT: 23.98 NTSC NDF
}

// Parse a Framerate from a float64 as a non-ntsc framerate
func ExampleFromFloat_nonNTSCError() {
	_, err := rate.FromFloat(24.0, rate.NTSCNone)
	fmt.Println("ERROR:", err)

	// Output:
	// ERROR: could not parse Framerate: non-ntsc framerates cannot be parsed from floats due to imprecision
}

// Parse a rational string as an NTSC playback speed.
func ExampleFromString_rationalNTSCPlayback() {
	framerate, err := rate.FromString("24000/1001", rate.NTSCNonDrop)
	if err != nil {
		panic(err)
	}

	fmt.Println("RESULT:", framerate)

	// Output:
	// RESULT: 23.98 NTSC NDF
}

// Parse a rational string as an NTSC playback speed.
func ExampleFromString_rationalNTSCTimebase() {
	framerate, err := rate.FromString("24/1", rate.NTSCNonDrop)
	if err != nil {
		panic(err)
	}

	fmt.Println("RESULT:", framerate)

	// Output:
	// RESULT: 23.98 NTSC NDF
}

// Parse a rational string as an NTSC playback speed.
//
// Many applications report the timebase as the seconds-per-frame (1/24) rather than the
// frames-per-seconds (24/1). When parsing a rational string, FromString will
// automatically invert the timebase if the numerator is larger than the denominator.
func ExampleFromString_rationalNTSCFlipped() {
	framerate, err := rate.FromString("1/24", rate.NTSCNonDrop)
	if err != nil {
		panic(err)
	}

	fmt.Println("RESULT:", framerate)

	// Output:
	// RESULT: 23.98 NTSC NDF
}

// Parse a float string as an NTSC playback speed / timebase.
//
// Parsing float strings for non-ntsc Framerates will result in an error. Floats are
// not precise enough to parse without the ability to coerce to an NTSC rate.
func ExampleFromString_float() {
	framerate, err := rate.FromString("23.98", rate.NTSCNonDrop)
	if err != nil {
		panic(err)
	}

	fmt.Println("RESULT:", framerate)

	// Output:
	// RESULT: 23.98 NTSC NDF
}

// Parse an integer string as an NTSC playback timebase.
func ExampleFromString_int() {
	framerate, err := rate.FromString("24", rate.NTSCNonDrop)
	if err != nil {
		panic(err)
	}

	fmt.Println("RESULT:", framerate)

	// Output:
	// RESULT: 23.98 NTSC NDF
}
