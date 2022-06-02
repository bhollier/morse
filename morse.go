package morse

import (
	"fmt"
	"strings"
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
		standardWordRuneDitDuration := standardWordDuration - standardWordFarnsworthDuration
		farnsworthDitDuration := ((time.Minute / time.Duration(farnsworthWPM)) -
			(time.Duration(standardWordRuneDitDuration) * ditDuration)) /
			time.Duration(standardWordFarnsworthDuration)
		return time.Duration(s.duration) * farnsworthDitDuration
	}
}

// String returns the signal's string representation, e.g. '・', '－', ' '
func (s Signal) String() string {
	return s.str
}

// Code is a sequence of Morse signals
type Code []Signal

func (c Code) Equal(o Code) bool {
	if len(c) != len(o) {
		return false
	}
	for i := range c {
		if c[i] != o[i] {
			return false
		}
	}
	return true
}

func (c Code) duration() (d uint) {
	for _, s := range c {
		d += s.duration
	}
	return
}

// Duration returns the duration of the code, at the given WPM.
// The word the WPM is based on is "PARIS".
//
// If farnsworthWPM is non-zero, the duration uses Farnsworth timing,
// where the speed of the characters is determined by wpm, but the
// actual words per minute is determined by farnsworthWPM. This is
// achieved by elongating the duration between letters and words
func (c Code) Duration(wpm, farnsworthWPM uint) (d time.Duration) {
	for _, s := range c {
		d += s.Duration(wpm, farnsworthWPM)
	}
	return d
}

func (c Code) String() string {
	sb := strings.Builder{}
	for _, s := range c {
		sb.WriteString(s.String())
	}
	return sb.String()
}
