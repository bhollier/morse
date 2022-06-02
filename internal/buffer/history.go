package buffer

// History is a buffer that stores N number of previous values,
// using a ring buffer
type History[T any] struct {
	buf  []T
	curr int
}

// NewHistory creates a history buffer with the given buffer size
func NewHistory[T any](bufSize int) History[T] {
	return History[T]{buf: make([]T, bufSize)}
}

// Add an element to the buffer
func (b *History[T]) Add(t T) {
	b.curr = mod(b.curr+1, len(b.buf))
	b.buf[b.curr] = t
}

// LastN retrieves the last n elements that were added,
// in reverse chronological order. If n is larger than the
// buffer size, the returned slice will be equal to the
// size of the buffer and the entire history will be returned
func (b *History[T]) LastN(n int) (history []T) {
	history = make([]T, intMin(len(b.buf), n))
	curr := b.curr
	for i := range history {
		history[i] = b.buf[curr]
		curr = mod(curr-1, len(b.buf))
	}
	return
}
