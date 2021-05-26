<h1 align="center">vtc-go</h1>
<p align="center">
    <img height=150 class="heightSet" align="center" src="https://raw.githubusercontent.com/opencinemac/vtc-py/master/zdocs/source/_static/logo1.svg"/>
</p>
<p align="center">A SMPTE Timecode Library for Go</p>
<p align="center">
    <a href="https://dev.azure.com/peake100/Open%20Cinema%20Collective/_build?definitionId=17"><img src="https://dev.azure.com/peake100/Open%20Cinema%20Collective/_apis/build/status/vtc-go?branchName=dev" alt="click to see build pipeline"></a>
    <a href="https://dev.azure.com/peake100/Open%20Cinema%20Collective/_build?definitionId=17"><img src="https://img.shields.io/azure-devops/tests/peake100/Open%20Cinema%20Collective/17/dev?compact_message" alt="click to see build pipeline"></a>
    <a href="https://dev.azure.com/peake100/Open%20Cinema%20Collective/_build?definitionId=17"><img src="https://img.shields.io/azure-devops/coverage/peake100/Open%20Cinema%20Collective/17/dev?compact_message" alt="click to see build pipeline"></a>
</p>
<p align="center">
    <a href="https://goreportcard.com/report/github.com/opencinemac/vtc-go"><img src="https://goreportcard.com/badge/github.com/opencinemac/vtc-go" alt="click to see report card"></a>
    <a href="https://codeclimate.com/github/opencinemac/vtc-go/maintainability"><img src="https://api.codeclimate.com/v1/badges/72bffc76c41f12c9ab71/maintainability" alt="click here to see report"/></a>
</p>
<p align="center">
    <a href="https://github.com/opencinemac/vtc-go"><img src="https://img.shields.io/github/go-mod/go-version/opencinemac/vtc-go" alt="Repo"></a>
    <a href="https://pkg.go.dev/github.com/opencinemac/vtc-go?readme=expanded#section-documentation"><img src="https://pkg.go.dev/badge/github.com/opencinemac/vtc-go?readme=expanded#section-documentation.svg" alt="Go Reference"></a>
</p>

# Overview

``vtc-go`` is inspired by years of scripting workflow solutions in a Hollywood cutting
room. It aims to capture all the ways in which timecode is used throughout the industry 
so  users can spend more time on their workflow logic, and less time handling the
corner-cases of parsing and calculating timecode.

## Demo

Let's take a quick high-level look at what you can do with vtc-rs:

```go
package main

import (
	"fmt"
	"github.com/opencinemac/vtc-go/pkg/rate"
	"github.com/opencinemac/vtc-go/pkg/tc"
	"math/big"
)

func main() {
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
}
```

## Features

- SMPTE Conventions:
    - [X] NTSC
    - [X] Drop-Frame
    - [ ] Interlaced timecode
- Timecode Representations:
    - Timecode    | '01:00:00:00'
    - Frames      | 86400
    - Seconds     | 3600.0
    - Runtime     | '01:00:00.0'
    - Rational    | 18018/5
    - Feet+Frames | '5400+00'
        - [X] 35mm, 4-perf
        - [ ] 35mm, 3-perf
        - [ ] 35mm, 2-perf
        - [ ] 16mm
    - Premiere Ticks | 15240960000000
- Operations:
    - Comparisons (==, <, <=, >, >=)
    - Add
    - Subtract
    - Scale (multiply and divide)
    - Div/Rem
    - Modulo
    - Negative
    - Absolute
    - Rebase (recalculate frame count at new framerate)
- Flexible Parsing:
    - Partial timecodes      | '1:12'
    - Partial runtimes       | '1.5'
    - Negative string values | '-1:12', '-3+00'
    - Poorly formatted tc    | '1:13:4'
- Built-in consts for common framerates.

## Goals

- Parse and fetch all Timecode representations.
- A clean, idiomatic API.
- Support all operations that make sense for timecode.

## Non-Goals

- Real-time timecode generators.

## Attributions

<div>Drop-frame calculations adapted from <a href="https://www.davidheidelberger.com/2010/06/10/drop-frame-timecode/">David Heidelberger's blog.</a></div>
<div>Logo made by <a href="" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a></div>
