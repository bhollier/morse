package morse

import "io"

// Writer is an interface for writing Morse signals.
// While it generally functions the same as an io.Writer,
// prefer using this where possible
type Writer interface {
	Write(p []Signal) (n int, err error)
}

// ToByteWriter is a simple wrapper around a Reader,
// so it acts like an io.Reader. This works because
// Signal is a byte under the hood. Beware when using
// this, as it can be confusing to treat a Signal as a byte
type ToByteWriter struct {
	Writer
}

func (r ToByteWriter) Write(p []byte) (n int, err error) {
	return r.Writer.Write(byteArrayToSignalArray(p))
}

// FromByteWriter is a simple wrapper around an io.Writer,
// so it acts like a Writer. This works because Signal is
// a byte under the hood. Beware when using this, as it
// can be confusing to treat a Signal as a byte
type FromByteWriter struct {
	io.Writer
}

func (r FromByteWriter) Write(p []Signal) (n int, err error) {
	return r.Writer.Write(signalArrayToByteArray(p))
}

// BlockingChannelWriter writes signals to a Go channel,
// blocking the running thread until the signals are all written
type BlockingChannelWriter chan<- Signal

// NonBlockingChannelWriter writes signals to a Go channel.
// If no more signals can be written, the writer immediately returns
type NonBlockingChannelWriter chan<- Signal

// WriterFromChan creates a Writer that writes signals to a Go channel.
// The blocking argument determines whether Writer.Write blocks until
// all signals are written, see BlockingChannelReader and
// NonBlockingChannelReader for more info
func WriterFromChan(c chan<- Signal, blocking bool) Writer {
	if blocking {
		return BlockingChannelWriter(c)
	} else {
		return NonBlockingChannelWriter(c)
	}
}

func (r BlockingChannelWriter) Write(p []Signal) (n int, err error) {
	for _, s := range p {
		r <- s
		n++
	}
	return
}

func (r NonBlockingChannelWriter) Write(p []Signal) (n int, err error) {
	for _, s := range p {
		select {
		case r <- s:
			n++
		default:
			return
		}
	}
	return
}
