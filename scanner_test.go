package morse

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// Here we mostly test the split functions (rather than the Scanner itself)
// because bufio.Scanner is already tested pretty thoroughly in the standard lib

func TestScanner_ScanWords(t *testing.T) {
	a := assert.New(t)

	codeString := "- . ... - / ... - .-. .. -. --."
	codeReader := ReaderFromCodeString(strings.NewReader(codeString))
	s := NewScanner(codeReader)
	s.Split(ScanWords)
	a.True(s.Scan())
	a.NoError(s.Err())
	a.Equal("－ ・ ・・・ －", s.Code().String())

	a.True(s.Scan())
	a.NoError(s.Err())
	a.Equal("・・・ － ・－・ ・・ －・ －－・", s.Code().String())

	a.False(s.Scan())
	a.NoError(s.Err())
}

func TestScanner_ScanRunes(t *testing.T) {
	a := assert.New(t)

	codeString := "- . ... - / ... - .-. .. -. --."
	codeReader := ReaderFromCodeString(strings.NewReader(codeString))
	s := NewScanner(codeReader)
	s.Split(ScanRunes)

	for _, expectedCode := range []Code{T, E, S, T, S, T, R, I, N, G} {
		a.True(s.Scan())
		a.NoError(s.Err())
		a.Equal(expectedCode.String(), s.Code().String())
	}

	a.False(s.Scan())
	a.NoError(s.Err())
}

func TestScanner_ScanAudible(t *testing.T) {
	a := assert.New(t)

	codeString := "- --- / -... ."
	codeReader := ReaderFromCodeString(strings.NewReader(codeString))
	s := NewScanner(codeReader)
	s.Split(ScanAudible)

	for _, expectedCode := range []Signal{Dah, Dah, Dah, Dah, Dah, Dit, Dit, Dit, Dit} {
		a.True(s.Scan())
		a.NoError(s.Err())
		a.Equal(expectedCode.String(), s.Code().String())
	}

	a.False(s.Scan())
	a.NoError(s.Err())
}

func TestScanner_ScanSignals(t *testing.T) {
	a := assert.New(t)

	codeString := "- --- / -... ."
	code := FromCodeString(codeString)
	codeReader := ReaderFromCodeString(strings.NewReader(codeString))
	s := NewScanner(codeReader)
	s.Split(ScanSignals)

	for _, expectedSignal := range code {
		a.True(s.Scan())
		a.NoError(s.Err())
		a.Equal(expectedSignal.String(), s.Code().String())
	}

	a.False(s.Scan())
	a.NoError(s.Err())
}
