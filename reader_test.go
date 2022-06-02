package morse

import (
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestFromCode(t *testing.T) {
	a := assert.New(t)

	codeStr := "・・・ －－－ ・・・"
	code := Code{
		Dit, SignalSpace, Dit, SignalSpace, Dit, RuneSpace,
		Dah, SignalSpace, Dah, SignalSpace, Dah, RuneSpace,
		Dit, SignalSpace, Dit, SignalSpace, Dit}

	a.Equal(codeStr, FromCodeString(codeStr).String())

	// Read individual signals (to test the overflow)
	r := ReaderFromCodeString(strings.NewReader(codeStr))
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

	codeStr = ".../---/..."
	code = Code{
		Dit, SignalSpace, Dit, SignalSpace, Dit, WordSpace,
		Dah, SignalSpace, Dah, SignalSpace, Dah, WordSpace,
		Dit, SignalSpace, Dit, SignalSpace, Dit}
	a.Equal(code.String(), FromCodeString(codeStr).String())

	codeStr = ".../ --- /..."
	a.Equal(code.String(), FromCodeString(codeStr).String())

	codeStr = ".... . .-.. .-.. --- / .-- --- .-. .-.. -.."
	code = Code{
		Dit, SignalSpace, Dit, SignalSpace, Dit, SignalSpace, Dit, RuneSpace,
		Dit, RuneSpace,
		Dit, SignalSpace, Dah, SignalSpace, Dit, SignalSpace, Dit, RuneSpace,
		Dit, SignalSpace, Dah, SignalSpace, Dit, SignalSpace, Dit, RuneSpace,
		Dah, SignalSpace, Dah, SignalSpace, Dah,
		WordSpace,
		Dit, SignalSpace, Dah, SignalSpace, Dah, RuneSpace,
		Dah, SignalSpace, Dah, SignalSpace, Dah, RuneSpace,
		Dit, SignalSpace, Dah, SignalSpace, Dit, RuneSpace,
		Dit, SignalSpace, Dah, SignalSpace, Dit, SignalSpace, Dit, RuneSpace,
		Dah, SignalSpace, Dit, SignalSpace, Dit, SignalSpace,
	}
	a.Equal(code.String(), FromCodeString(codeStr).String())

	codeStr = "..i -. n ...- v .-alid .-.. t -..ex t"
	codeStrClean := ".. -. ...- .- .-.. -.."
	a.Equal(FromCodeString(codeStrClean).String(), FromCodeString(codeStr).String())
}
