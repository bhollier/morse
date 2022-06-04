package morse

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/loremipsum.v1"
	"io"
	"testing"
)

func TestDecode(t *testing.T) {
	a := assert.New(t)

	code := FromCodeString("... --- ...")
	text := "SOS"
	a.Equal(text, Decode(code))

	code = FromCodeString(".... . .-.. .-.. --- / .-- --- .-. .-.. -..")
	text = "HELLO WORLD"
	a.Equal(text, Decode(code))

	code = FromCodeString(".--. .- .-. .. ...")
	text = "PARIS"
	a.Equal(text, Decode(code))

	// Read individual signals (to test the overflow)
	r := NewDecoder(NewReader(code))
	buf := make([]byte, 1)
	for i := range text {
		n, err := r.Read(buf)
		a.NoError(err)
		a.Equal(1, n)
		a.Equal(text[i], buf[0])
	}
	n, err := r.Read(buf)
	a.Equal(0, n)
	a.Equal(io.EOF, err)
}

const benchmarkDecoderSeed = 42
const benchmarkDecoderBufferSize = 512

func benchmarkDecoder(b *testing.B, paragraphs int) {
	loremIpsumGenerator := loremipsum.NewWithSeed(benchmarkDecoderSeed)
	text := loremIpsumGenerator.Paragraphs(paragraphs)
	// Convert the text into morse code
	code := FromText(text)

	benchmarkReader(b, benchmarkDecoderBufferSize, func() genericReader[byte] {
		return NewDecoder(NewReader(code))
	})
}

func BenchmarkDecoder1(b *testing.B) {
	benchmarkDecoder(b, 1)
}

func BenchmarkDecoder2(b *testing.B) {
	benchmarkDecoder(b, 2)
}

func BenchmarkDecoder3(b *testing.B) {
	benchmarkDecoder(b, 3)
}

func BenchmarkDecoder4(b *testing.B) {
	benchmarkDecoder(b, 4)
}

func BenchmarkDecoder5(b *testing.B) {
	benchmarkDecoder(b, 5)
}

func BenchmarkDecoder10(b *testing.B) {
	benchmarkDecoder(b, 10)
}
