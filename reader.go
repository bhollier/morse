package morse

import (
	"io"
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
