package testdata

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
)

//go:embed "test-timelines/PPRO/Many Basic Edits/Many Basic Edits.json"
var manyBasiEditsJSON []byte

// ManyBasicEditsData contains 215 events with no blank spaces, respeeds, or other
// effects
var ManyBasicEditsData = mustLoadSequence(manyBasiEditsJSON)

// TimecodeData holds data about a specific timecode value combined from XML and EDL
// data.
type TimecodeData struct {
	// Timebase is the Timecode's timebase.
	Timebase int64
	// NTSC is whether the Timecode is displayed using the NTSC convention.
	Ntsc bool
	// DropFrame is whether the Timecode is displayed using the NTSC drop-frame
	// convention/
	DropFrame bool `json:"drop_frame"`
	// FrameRateFrac is the framerate represented as a fraction.
	FrameRateFrac *JsonRat `json:"frame_rate_frac"`
	// Timecode is the timecode value pulled from the EDL.
	Timecode string
	// Frame is the frame number pulled from the XML.
	Frame int64
	// FrameXMLRaw is the raw frame number pulled from the XML without the start time
	// of the media/sequence added.
	FrameXMLRaw int64 `json:"frame_xml_raw"`
	// SecondsRational is the seconds representation of the timecode as a fraction.
	SecondsRational *JsonRat `json:"seconds_rational"`
	// SecondsDecimal is the seconds representation of the timecode as a decimal string
	// with a precision of 9 places.
	SecondsDecimal string `json:"seconds_decimal"`
	// PProTicks is the number of Adobe Premiere Pro ticks for this timecode pulled from
	// the XML.
	PProTicks int64 `json:"ppro_ticks"`
	// PProTicksRaw is the number of Adobe Premiere Pro ticks for this timecode pulled
	// from the XML without the start time of the media / sequence added.
	PProTicksRaw int64 `json:"ppro_ticks_raw"`
	// FeetAndFrames is the feet-and-frames representation of this timecode.
	FeetAndFrames string `json:"feet_and_frames"`
	// Runtime is the real-world runtime representation of this timecode with a
	// precision of up to 9 places.
	Runtime string
}

// EventData holds the timecode data of a single sequence event using data from a
// combined EDL and FCP7XML
type EventData struct {
	// DurationFrames is the length of thee event in frames.
	DurationFrames int64 `json:"duration_frames"`
	// RecordIn is the sequence timecode info for the start of this event.
	RecordIn TimecodeData `json:"record_in"`
	// RecordOut is the sequence timecode info for the end of this event.
	RecordOut TimecodeData `json:"record_out"`
	// SourceIn is the media timecode info for the start of this event.
	SourceIn TimecodeData `json:"source_in"`
	// SourceOut is the media timecode info for the end of this event.
	SourceOut TimecodeData `json:"source_out"`
}

// SequenceData holds all the timecode events of a sequence using data from a combined
// EDL and FCP7XML
type SequenceData struct {
	// StartTime is the timecode info for the first frame of this sequence.
	StartTime TimecodeData `json:"start_time"`
	// TotalDurationFrames is the number of frames this sequence contains.
	TotalDurationFrames int64 `json:"total_duration_frames"`
	// Events is a slice of timecode information for all the edit events for this
	// sequence.
	Events []*EventData
}

// JsonRat implements json.Unmarshaler for *big.Rat.
type JsonRat big.Rat

func (rat *JsonRat) UnmarshalJSON(data []byte) error {
	dataStr := string(data)
	dataStr = strings.Trim(dataStr, "\"")
	ratValue, ok := new(big.Rat).SetString(dataStr)
	if !ok {
		return fmt.Errorf("could not decode '%v' as fraction", dataStr)
	}

	*rat = JsonRat(*ratValue)
	return nil
}

func mustLoadSequence(data []byte) *SequenceData {
	deserialized := new(SequenceData)

	buffer := bytes.NewBuffer(data)
	decoder := json.NewDecoder(buffer)
	err := decoder.Decode(deserialized)
	if err != nil {
		panic(fmt.Errorf("error loading test sequence data: %w", err))
	}

	return deserialized
}
