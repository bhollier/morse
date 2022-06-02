package morse

import (
	"github.com/stretchr/testify/assert"
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
