package morse

import (
	"bytes"
	"github.com/bhollier/morse/internal/buffer"
	"io"
	"unicode"
	"unicode/utf8"
)

// Decoder converts Morse Code into human-readable text
type Decoder struct {
	morseWordScanner Scanner
	overflow         buffer.Overflow[byte]
	started          bool
}

// NewDecoder creates a Morse Decoder from the given Reader
func NewDecoder(r Reader) *Decoder {
	return NewDecoderFromScanner(NewScanner(r))
}

// NewDecoderFromScanner creates a Morse Decoder from the given Scanner
func NewDecoderFromScanner(s Scanner) *Decoder {
	s.Split(ScanWords)
	return &Decoder{morseWordScanner: s}
}

func (d *Decoder) Read(b []byte) (n int, err error) {
	// First, try to empty the overflow from the last read
	n = d.overflow.Empty(b)
	b = b[n:]

	for len(b) > 0 {
		if !d.morseWordScanner.Scan() {
			err = d.morseWordScanner.Err()
			if err == nil {
				err = io.EOF
			}
			return
		}

		word := bytes.Buffer{}

		// First, check if we need to add a preceding space
		// to separate from the previous word
		if d.started {
			word.WriteRune(' ')
		} else {
			d.started = true
		}

		wordCode := d.morseWordScanner.Code()
		codeRuneBuf := make(Code, 0, len(StandardWordCode))
		for i, s := range wordCode {
			// If the signal isn't a rune
			if s != RuneSpace {
				// Add the signal to the buffer
				codeRuneBuf = append(codeRuneBuf, s)
			}

			// If the signal is a rune space or this is the end of the word
			if s == RuneSpace || i+1 == len(wordCode) {
				r := Dictionary.FromCode(codeRuneBuf)
				if r == utf8.RuneError {
					r = '?'
				} else {
					r = unicode.ToUpper(r)
				}
				word.WriteRune(r)
				codeRuneBuf = codeRuneBuf[:0]
			}
		}

		// Copy it
		bytesCopied := d.overflow.Copy(b, word.Bytes())
		b = b[bytesCopied:]
		n += bytesCopied
	}

	return
}

// Decode returns the human-readable text of the given code
func Decode(code Code) string {
	d := NewDecoder(NewReader(code))
	b, err := io.ReadAll(d)
	if err != nil {
		// panic on error as neither CodeReader or Decoder should ever error
		panic(err)
	}
	return string(b)
}
