package morse

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/loremipsum.v1"
	"io"
	"strings"
	"testing"
)

func TestFromText(t *testing.T) {
	a := assert.New(t)

	text := "SOS"
	code := FromCodeString("... --- ...")
	a.Equal(code.String(), FromText(text).String())

	text = "sos"
	a.Equal(code.String(), FromText(text).String())

	text = "Hello World"
	code = FromCodeString(".... . .-.. .-.. --- / .-- --- .-. .-.. -..")
	a.Equal(code.String(), FromText(text).String())

	text = "PARIS"
	code = FromCodeString(".--. .- .-. .. ...")
	a.Equal(code.String(), FromText(text).String())

	// Read individual signals (to test the overflow)
	r := ReaderFromText(strings.NewReader(text))
	buf := make([]Signal, 1)
	for i := range code {
		n, err := r.Read(buf)
		a.NoError(err)
		a.Equal(1, n)
		a.Equal(code[i], buf[0])
	}
	n, err := r.Read(buf)
	a.Equal(0, n)
	a.Equal(io.EOF, err)
}

const benchmarkTextEncoderSeed = 42
const benchmarkTextEncoderBufferSize = 512

func benchmarkTextEncoder(b *testing.B, paragraphs int) {
	loremIpsumGenerator := loremipsum.NewWithSeed(benchmarkTextEncoderSeed)
	text := loremIpsumGenerator.Paragraphs(paragraphs)

	benchmarkReader(b, benchmarkTextEncoderBufferSize, func() genericReader[Signal] {
		return ReaderFromText(strings.NewReader(text))
	})
}

func BenchmarkTextEncoder1(b *testing.B) {
	benchmarkTextEncoder(b, 1)
}

func BenchmarkTextEncoder2(b *testing.B) {
	benchmarkTextEncoder(b, 2)
}

func BenchmarkTextEncoder3(b *testing.B) {
	benchmarkTextEncoder(b, 3)
}

func BenchmarkTextEncoder4(b *testing.B) {
	benchmarkTextEncoder(b, 4)
}

func BenchmarkTextEncoder5(b *testing.B) {
	benchmarkTextEncoder(b, 5)
}

func BenchmarkTextEncoder10(b *testing.B) {
	benchmarkTextEncoder(b, 10)
}
