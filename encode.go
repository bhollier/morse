package morse

import (
	"bufio"
	"github.com/bhollier/morse/internal/buffer"
	"io"
	"strings"
)

// TextEncoder converts human-readable byte text from an io.Reader into Morse signals
type TextEncoder struct {
	wordScanner *bufio.Scanner
	overflow    buffer.Overflow[Signal]
	started     bool
}

// ReaderFromText creates a TextEncoder that retrieves human-readable text
// from the given io.Reader and converts it into Morse Code
func ReaderFromText(r io.Reader) *TextEncoder {
	return ReaderFromTextScanner(bufio.NewScanner(r))
}

// ReaderFromTextScanner creates a TextEncoder that retrieves human-readable text
// from the given bufio.Scanner and converts it into Morse Code
func ReaderFromTextScanner(s *bufio.Scanner) *TextEncoder {
	s.Split(bufio.ScanWords) // todo what if the input has no spaces?
	// todo also doesn't handle newlines
	return &TextEncoder{wordScanner: s}
}

func (e *TextEncoder) Read(p []Signal) (n int, err error) {
	// First, try to empty the overflow from the last read
	n = e.overflow.Empty(p)
	p = p[n:]

	// While there is space in p and there are words to encode
	for len(p) > 0 {
		if !e.wordScanner.Scan() {
			return n, io.EOF
		}

		// Read a word
		word := e.wordScanner.Text()

		wordCode := make(Code, 0)

		// First, check if we need to add a preceding space
		// to separate from the previous word
		if e.started {
			wordCode = append(wordCode, WordSpace)
		} else {
			e.started = true
		}

		// Iterate over the runes of the word
		for i, r := range word {
			runeCode := Dictionary.FromRune(r)
			if runeCode == nil {
				runeCode = Dictionary.FromRune('?')
			}
			wordCode = append(wordCode, runeCode...)
			if i+1 < len(word) {
				wordCode = append(wordCode, RuneSpace)
			}
		}

		// Copy the code into p (with the remaining going into the buffer)
		signalsCopied := e.overflow.Copy(p, wordCode)
		p = p[signalsCopied:]
		n += signalsCopied
	}

	return
}

// FromText returns the Morse code of the given human-readable string
func FromText(text string) Code {
	e := ReaderFromText(strings.NewReader(text))
	c, err := ReadAll(e)
	if err != nil {
		// panic on error as neither strings.Reader or TextEncoder should ever error
		panic(err)
	}
	return c
}
