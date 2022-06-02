package main

import (
	"fmt"
	"github.com/bhollier/morse"
	"unicode"
)

type InputMode rune

const (
	TextInputMode  = InputMode('t')
	MorseInputMode = InputMode('m')
)

func ParseInputMode(s string) (InputMode, error) {
	if s == "" {
		return 0, fmt.Errorf("missing input mode")
	}

	switch InputMode(unicode.ToLower([]rune(s)[0])) {
	case TextInputMode:
		return TextInputMode, nil
	case MorseInputMode:
		return MorseInputMode, nil
	default:
		return 0, fmt.Errorf("unknown input mode %s", s)
	}
}

func (m InputMode) ConvertInput(input string) morse.Code {
	switch m {
	case TextInputMode:
		return morse.FromText(input)
	case MorseInputMode:
		return morse.FromCodeString(input)
	default:
		panic(fmt.Errorf("unknown input mode %c", m))
	}
}
