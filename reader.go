package morse

import (
	"bufio"
	"github.com/bhollier/morse/internal/buffer"
	"io"
	"strings"
)

// Reader is an interface for reading Morse signals.
// Functions the same as an io.Reader, but for reading Signal
type Reader interface {
	Read(p []Signal) (n int, err error)
}

// ReadAll reads from r until an error or io.EOF and returns the data it read.
// A successful call returns err == nil, not err == io.EOF. Because ReadAll is
// defined to read from src until io.EOF, it does not treat an io.EOF from Read
// as an error to be reported.
func ReadAll(r Reader) ([]Signal, error) {
	b := make([]Signal, 0, 512)
	for {
		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, Signal{})[:len(b)]
		}
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return b, err
		}
	}
}

// CodeReader is a reader for Code slice
type CodeReader struct {
	c Code
	i int64
}

// NewReader creates a Reader for the given Code
func NewReader(p Code) *CodeReader {
	return &CodeReader{p, 0}
}

func (r *CodeReader) Read(p []Signal) (n int, err error) {
	if r.i >= int64(len(r.c)) {
		return 0, io.EOF
	}
	n = copy(p, r.c[r.i:])
	r.i += int64(n)
	return
}

// BlockingChannelReader reads signals from a Go channel,
// blocking the running thread until the channel is closed
type BlockingChannelReader <-chan Signal

// NonBlockingChannelReader reads signals from a Go channel.
// If no signals are available, the reader immediately returns
type NonBlockingChannelReader <-chan Signal

// ReaderFromChan creates a Reader that retrieves
// signals from a Go channel. The blocking argument determines whether
// Reader.Read blocks until the slice is full, see BlockingChannelReader
// and NonBlockingChannelReader for more info
func ReaderFromChan(c <-chan Signal, blocking bool) Reader {
	if blocking {
		return BlockingChannelReader(c)
	} else {
		return NonBlockingChannelReader(c)
	}
}

func (r BlockingChannelReader) Read(p []Signal) (n int, err error) {
	for i := range p {
		s, more := <-r
		if !more {
			return n, io.EOF
		}
		p[i] = s
		n++
	}
	return
}

func (r NonBlockingChannelReader) Read(p []Signal) (n int, err error) {
	for i := range p {
		select {
		case s, more := <-r:
			if !more {
				return n, io.EOF
			}
			p[i] = s
			n++
		default:
			return
		}
	}
	return
}

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
