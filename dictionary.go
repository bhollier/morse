package morse

import "unicode"

type dictionary struct {
	runeCodeMap map[rune]Code
	// todo codeRuneMap to go the other way
}

var Dictionary = dictionary{
	runeCodeMap: make(map[rune]Code),
}

// Add an entry for linking the rune r with the morse code c
func (d *dictionary) Add(r rune, c Code) {
	d.runeCodeMap[r] = c
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
