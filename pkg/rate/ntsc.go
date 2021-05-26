package rate

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
