package morse

import (
	"fmt"
	"time"
	"unsafe"
)

// MaxSignalDuration is the maximum duration of a signal,
// relative to the length of a dit. This value is 2^7 (128)
// as 7 out of 8 of the bits in a signal byte are for the
// duration (the other is the audible flag), and not 2^7-1
// as the duration is 1-based (not 0-based) because 0 is
// not a valid signal duration
const MaxSignalDuration = 128

// Signal represents a single (possible audible) Morse signal,
// e.g. a dit, dah, single space, word space, etc.
//
// Internally, signals are stored as bytes for memory
// efficiency, where the first bit is if the signal is
// audible, and the other 7 are duration (relative to the
// length of a dit). The duration is 1-based (not 0-based),
// as 0 is not a valid duration
type Signal byte

// NewSignal creates a new audible or inaudible signal.
// duration specifies the duration of the signal relative to
// a Dit or SignalSpace. So a Dah signal has a duration of 3,
// a SignalSpace 1, RuneSpace 3, etc.
//
// For example, the following creates a Dit signal:
//  NewSignal(true, 1)
// Or the following to create a WordSpace:
//  NewSignal(false, 7)
//
// This method shouldn't generally need to be used, as the
// predefined Signal constants exist for standard Morse
// signals (Dit, Dah, SignalSpace, RuneSpace, WordSpace).
//
// Panics if duration is 0 or larger than
// MaxSignalDuration (128)
func NewSignal(audible bool, duration uint8) Signal {
	if duration == 0 {
		panic("signal duration cannot be 0")
	}
	if duration > MaxSignalDuration {
		panic(fmt.Errorf("signal duration %d > %d", duration, MaxSignalDuration))
	}

	var b byte

	// The first bit is if the signal is audible
	if audible {
		b |= 1 << 0
	}

	// The final 7 bits are the duration
	b |= (duration - 1) << 1

	return Signal(b)
}

// Audible returns whether the signal is audible,
// aka whether it is a beep or silent
func (s Signal) Audible() bool {
	return s&1 != 0
}

// DitDuration is the signal duration relative to a
// Dit (e.g. the dit duration of Dah is 3 as it is
// x3 as long as a Dit). This will never return 0.
//
// For most circumstances Duration should be used as
// it gives a more usable time.Duration
func (s Signal) DitDuration() uint {
	return uint(s>>1) + 1
}

// String returns the signal's string representation,
// e.g. '・', '－', ' '. If the signal is unknown,
// returns "?"
func (s Signal) String() string {
	switch s {
	case Dit:
		return "・"
	case Dah:
		return "－"
	case SignalSpace:
		return ""
	case RuneSpace:
		return " "
	case WordSpace:
		return "  "
	default:
		return "?"
	}
}

// Converts the given signal array to a byte array,
// utilising the fact that Signal's underlying type is byte.
// This is relatively risky as we're subverting Go's type
// system with unsafe.Pointer, but is much faster than
// creating a new slice and copying the individual Signal
// bytes
func signalArrayToByteArray(s []Signal) []byte {
	// unsafe, yuck!
	return *(*[]byte)(unsafe.Pointer(&s))
}

// Reverse of signalArrayToByteArray
func byteArrayToSignalArray(b []byte) []Signal {
	return *(*[]Signal)(unsafe.Pointer(&b))
}

const (
	// Dit (aka dot) is an audible beep signal. Equivalent to
	//  NewSignal(true, 1)
	Dit = Signal(0b00000001)

	// Dah (aka dash) is an audible beep that is 3x as long
	// as a Dit. Equivalent to
	//  NewSignal(true, 3)
	Dah = Signal(0b00000101)

	// SignalSpace represents the space between audible signals,
	// and is as long as a Dit. Equivalent to
	//  NewSignal(false, 1)
	SignalSpace = Signal(0b00000000)

	// RuneSpace represents a space between runes (aka letters),
	// so is 3x as long as a SignalSpace. Equivalent to
	//  NewSignal(false, 3)
	RuneSpace = Signal(0b00000100)

	// WordSpace represents a space between words, so is 7x as
	// long as a SignalSpace. Equivalent to
	//  NewSignal(false, 7)
	WordSpace = Signal(0b00001100)
)

// StandardWord is the standard word for WPM calculations
const StandardWord = "PARIS"

// StandardWordCode is the Morse Code of StandardWord, including the WordSpace on the end
var StandardWordCode = append(Join([]Code{P, A, R, I, S}, Code{RuneSpace}), WordSpace)

// The duration of the standard word, where 1 = the length of a dit
var standardWordDuration = StandardWordCode.DitDuration()

// The duration of the dits, dahs and rune spaces in the standard word.
// Used to calculate Farnsworth timing
var standardWordRuneDuration uint

// The duration of the spaces in the standard word (between letters and words).
// Used to calculate Farnsworth timing
var standardWordSpaceDuration uint

func init() {
	for _, s := range StandardWordCode {
		d := s.DitDuration()
		if s == Dit || s == Dah || s == SignalSpace {
			standardWordRuneDuration += d
		} else if s == RuneSpace || s == WordSpace {
			standardWordSpaceDuration += d
		} else {
			panic("this shouldn't be possible")
		}
	}
}

// Duration returns the duration of the signal, at the given WPM.
// The standard word for the WPM is "PARIS".
//
// If farnsworthWPM is non-zero, the duration is based on Farnsworth timing,
// see Code.Duration for more info
func (s Signal) Duration(wpm, farnsworthWPM uint) time.Duration {
	if farnsworthWPM > wpm {
		panic(fmt.Errorf("farnswordWPM (%d) > wpm (%d)", farnsworthWPM, wpm))
	}

	ditDuration := time.Minute / time.Duration(standardWordDuration*wpm)
	if farnsworthWPM == 0 || (s != RuneSpace && s != WordSpace) {
		return time.Duration(s.DitDuration()) * ditDuration
	} else {
		farnsworthDitDuration := ((time.Minute / time.Duration(farnsworthWPM)) -
			(time.Duration(standardWordRuneDuration) * ditDuration)) /
			time.Duration(standardWordSpaceDuration)
		return time.Duration(s.DitDuration()) * farnsworthDitDuration
	}
}
