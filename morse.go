// Package morse contains definitions for representing morse code programmatically,
// as well as some general readers and writers similar to io.Reader/io.Writer
package morse

import (
	"strings"
	"time"
)

// Code is a sequence of Morse signals
type Code []Signal

// Join concatenates the given slice of Code, using the given delimiter.
// Similar to strings.Join
func Join(elems []Code, sep Code) (code Code) {
	switch len(elems) {
	case 0:
		return Code{}
	case 1:
		return elems[0]
	}
	n := len(sep) * (len(elems) - 1)
	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}

	code = make(Code, 0, n)
	code = append(code, elems[0]...)
	for _, c := range elems[1:] {
		code = append(code, sep...)
		code = append(code, c...)
	}
	return
}

// JoinSignals joins the given codes with a SignalSpace between code sequences.
// Useful for creating prosigns
func JoinSignals(elems ...Code) Code {
	return Join(elems, Code{SignalSpace})
}

// JoinLetters joins the given codes with a RuneSpace between the code sequences
func JoinLetters(letters ...Code) Code {
	return Join(letters, Code{RuneSpace})
}

// JoinWords joins the given codes with a WordSpace between the code sequences
func JoinWords(words ...Code) Code {
	return Join(words, Code{WordSpace})
}

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

// DitDuration returns the duration of the code, relative to a Dit.
// See Signal.DitDuration for more info
func (c Code) DitDuration() (d uint) {
	for _, s := range c {
		d += s.DitDuration()
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
