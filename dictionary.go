package morse

import (
	"unicode"
	"unicode/utf8"
)

type dictionary struct {
	runeCodeMap map[rune]Code
	codeRuneMap map[string]rune
}

var Dictionary = dictionary{
	runeCodeMap: make(map[rune]Code),
	codeRuneMap: make(map[string]rune),
}

// Add an entry for linking the rune r with the morse code c
func (d *dictionary) Add(r rune, c Code) {
	d.runeCodeMap[r] = c
	d.codeRuneMap[c.String()] = r
}

// AddCodeString is a wrapper around Add which calls FromCodeString on codeStr
func (d *dictionary) AddCodeString(r rune, codeStr string) {
	d.Add(r, FromCodeString(codeStr))
}

// FromRune returns the Morse code of the given human-readable rune, or nil if unknown
func (d *dictionary) FromRune(r rune) Code {
	c, _ := d.runeCodeMap[unicode.ToLower(r)]
	return c
}

// FromCode returns the human-readable rune of the given Morse code, or utf8.RuneError if unknown
func (d *dictionary) FromCode(c Code) rune {
	r, ok := d.codeRuneMap[c.String()]
	if !ok {
		return utf8.RuneError
	}
	return r
}
