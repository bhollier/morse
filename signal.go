package morse

import (
	"fmt"
	"time"
)

// Signal represents a single (possible audible) Morse signal,
// e.g. a dit, dah, single space, word space, etc.
type Signal struct {
	audible  bool
	duration uint
	str      string
}

// Audible returns whether the signal is audible,
// aka whether it is a beep or silent
func (s Signal) Audible() bool {
	return s.audible
}

// String returns the signal's string representation, e.g. '・', '－', ' '
func (s Signal) String() string {
	return s.str
}

var (
	Dit         = Signal{true, 1, "・"}
	Dah         = Signal{true, 3, "－"}
	SignalSpace = Signal{false, 1, ""}
	RuneSpace   = Signal{false, 3, " "}
	WordSpace   = Signal{false, 7, "  "}
)

// StandardWord is the standard word for WPM calculations
const StandardWord = "PARIS"

// StandardWordCode is the Morse Code of StandardWord, including the WordSpace on the end
var StandardWordCode = append(Join([]Code{P, A, R, I, S}, Code{RuneSpace}), WordSpace)

// The duration of the standard word, where 1 = the length of a dit
var standardWordDuration = StandardWordCode.duration()

// The duration of the dits, dahs and rune spaces in the standard word.
// Used to calculate Farnsworth timing
var standardWordRuneDuration uint

// The duration of the spaces in the standard word (between letters and words).
// Used to calculate Farnsworth timing
var standardWordSpaceDuration uint

func init() {
	for _, s := range StandardWordCode {
		d := s.duration
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
		return time.Duration(s.duration) * ditDuration
	} else {
		farnsworthDitDuration := ((time.Minute / time.Duration(farnsworthWPM)) -
			(time.Duration(standardWordRuneDuration) * ditDuration)) /
			time.Duration(standardWordSpaceDuration)
		return time.Duration(s.duration) * farnsworthDitDuration
	}
}
