package buffer

// Overflow is a buffer for copying into slices and
// storing the extra elements that there wasn't space for.
// Mainly designed to be used in a Reader
type Overflow[T any] struct {
	buf []T
}

// Empty as many of the elements in the overflow buffer into dst as possible,
// returning the number of elements that were copied
func (b *Overflow[T]) Empty(dst []T) (n int) {
	n = copy(dst, b.buf)
	b.buf = b.buf[n:]
	return
}

// Copy as many elements from src into dst as possible.
// Any remaining elements are added onto the end of the overflow
func (b *Overflow[T]) Copy(dst, src []T) (n int) {
	n = copy(dst, src)
	if n < len(src) {
		b.buf = append(b.buf, src[n:]...)
	}
	return
}
