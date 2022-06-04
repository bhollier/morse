package morse

import (
	"bufio"
	"github.com/bhollier/morse/internal/buffer"
	"io"
	"strings"
)

// CodeStringReader converts Morse code in text form from an io.Reader into Morse signals.
// It expects the code to be made up of the following characters:
//
// - '・' (or '.'), representing a Dit
//
// - '－' (or '-'), representing a Dah
//
// - ' ', representing a RuneSpace (space between the letters of a word)
//
// - '/', representing a WordSpace
//
// All other characters are ignored, as well as extra whitespace
type CodeStringReader struct {
	overflow    buffer.Overflow[Signal]
	codeScanner *bufio.Scanner
	wordStarted bool
}

// ReaderFromCodeString creates a CodeStringReader that retrieves code strings
// from the given io.Reader and converts it into Code
func ReaderFromCodeString(r io.Reader) *CodeStringReader {
	return ReaderFromCodeStringScanner(bufio.NewScanner(r))
}

// ReaderFromCodeStringScanner creates a CodeStringReader that retrieves code strings
// from the given bufio.Scanner and converts it into Code
func ReaderFromCodeStringScanner(s *bufio.Scanner) *CodeStringReader {
	s.Split(bufio.ScanWords)
	return &CodeStringReader{codeScanner: s}
}

func (r *CodeStringReader) Read(p []Signal) (n int, err error) {
	// First, try to empty the overflow from the last read
	n = r.overflow.Empty(p)
	p = p[n:]

	// While there is space in p and there is code to read
	for len(p) > 0 {
		if !r.codeScanner.Scan() {
			return n, io.EOF
		}

		wordStarted := r.wordStarted
		if !wordStarted {
			r.wordStarted = true
		}

		// Read a rune
		codeRunes := []rune(r.codeScanner.Text())

		code := make(Code, 0, (len(codeRunes)*2)-1)

		for i, codeRune := range codeRunes {
			switch codeRune {
			case '・', '.':
				code = append(code, Dit)
			case '－', '-':
				code = append(code, Dah)
			case '/':
				code = append(code, WordSpace)
				r.wordStarted = false
				// Continue, because we don't want to add a signal space
				continue
			default:
				continue
			}
			// If there is another code rune, and it isn't a word space
			if i+1 < len(codeRunes) && codeRunes[i+1] != '/' {
				code = append(code, SignalSpace)
			}
		}

		// Check if we need to add a preceding space
		// to separate from the previous rune
		if wordStarted {
			// Only add if the first signal was audible
			if len(code) > 0 && code[0].Audible() {
				code = append(Code{RuneSpace}, code...)
			}
		}

		// Copy the code into p (with the remaining going into the buffer)
		signalsCopied := r.overflow.Copy(p, code)
		p = p[signalsCopied:]
		n += signalsCopied
	}

	return
}

// FromCodeString converts the given string of morse code into Code.
// See CodeStringReader for more info
func FromCodeString(code string) Code {
	e := ReaderFromCodeString(strings.NewReader(code))
	c, err := ReadAll(e)
	if err != nil {
		// panic on error as neither strings.Reader or CodeReader should ever error
		panic(err)
	}
	return c
}
