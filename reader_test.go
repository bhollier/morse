package morse

import (
	"io"
	"testing"
)

// todo test and benchmark CodeReader

type genericReader[T any] interface {
	Read(b []T) (n int, err error)
}

func benchmarkReader[T any](b *testing.B, bufferSize int, readerGenerator func() genericReader[T]) {
	buf := make([]T, bufferSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Pause the timer just in case readerGenerator is slow
		b.StopTimer()
		r := readerGenerator()
		b.StartTimer()
		for true {
			_, err := r.Read(buf)
			if err != nil {
				if err != io.EOF {
					b.Fatal(err)
				}
				break
			}
		}
	}
}
