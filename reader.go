package morse

import (
	"io"
)

// Reader is an interface for reading Morse signals.
// While it generally functions the same as an io.Reader,
// prefer using this where possible
type Reader interface {
	Read(p []Signal) (n int, err error)
}

// ToByteReader is a simple wrapper around a Reader,
// so it acts like an io.Reader. This works because
// Signal is a byte under the hood. Beware when using
// this, as it can be confusing to treat a Signal as a byte
type ToByteReader struct {
	Reader
}

func (r ToByteReader) Read(p []byte) (n int, err error) {
	return r.Reader.Read(byteArrayToSignalArray(p))
}

// FromByteReader is a simple wrapper around an io.Reader,
// so it acts like a Reader. This works because Signal is
// a byte under the hood. Beware when using this, as it
// can be confusing to treat a Signal as a byte
type FromByteReader struct {
	io.Reader
}

func (r FromByteReader) Read(p []Signal) (n int, err error) {
	return r.Reader.Read(signalArrayToByteArray(p))
}

// ReadAll is the Signal equivalent of io.ReadAll
func ReadAll(r Reader) ([]Signal, error) {
	// Use io.ReadAll and ToByteReader since Signals are bytes
	bytes, err := io.ReadAll(ToByteReader{r})
	return byteArrayToSignalArray(bytes), err
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
