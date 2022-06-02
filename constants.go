package morse

import "unicode"

const StandardWord = "PARIS"

var (
	Dit         = Signal{true, 1, "・"}
	Dah         = Signal{true, 3, "－"}
	SignalSpace = Signal{false, 1, ""}
	RuneSpace   = Signal{false, 3, " "}
	WordSpace   = Signal{false, 7, "  "}
)

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

var standardWordCode Code
var standardWordDuration uint
var standardWordFarnsworthDuration uint

func init() {
	Dictionary.Add(' ', Code{WordSpace})

	Dictionary.AddCodeString('a', "・－")
	Dictionary.AddCodeString('b', "－・・・")
	Dictionary.AddCodeString('c', "－・－・")
	Dictionary.AddCodeString('d', "－・・")
	Dictionary.AddCodeString('e', "・")
	Dictionary.AddCodeString('f', "・・－・")
	Dictionary.AddCodeString('g', "－－・")
	Dictionary.AddCodeString('h', "・・・・")
	Dictionary.AddCodeString('i', "・・")
	Dictionary.AddCodeString('j', "・－－－")
	Dictionary.AddCodeString('k', "－・－")
	Dictionary.AddCodeString('l', "・－・・")
	Dictionary.AddCodeString('m', "－－")
	Dictionary.AddCodeString('n', "－・")
	Dictionary.AddCodeString('o', "－－－")
	Dictionary.AddCodeString('p', "・－－・")
	Dictionary.AddCodeString('q', "－－・－")
	Dictionary.AddCodeString('r', "・－・")
	Dictionary.AddCodeString('s', "・・・")
	Dictionary.AddCodeString('t', "－")
	Dictionary.AddCodeString('u', "・・－")
	Dictionary.AddCodeString('v', "・・・－")
	Dictionary.AddCodeString('w', "・－－")
	Dictionary.AddCodeString('x', "－・・－")
	Dictionary.AddCodeString('y', "－・－－")
	Dictionary.AddCodeString('z', "－－・・")

	Dictionary.AddCodeString('1', "・－－－－")
	Dictionary.AddCodeString('2', "・・－－－")
	Dictionary.AddCodeString('3', "・・・－－")
	Dictionary.AddCodeString('4', "・・・・－")
	Dictionary.AddCodeString('5', "・・・・・")
	Dictionary.AddCodeString('6', "－・・・・")
	Dictionary.AddCodeString('7', "－－・・・")
	Dictionary.AddCodeString('8', "－－－・・")
	Dictionary.AddCodeString('9', "－－－－・")
	Dictionary.AddCodeString('0', "－－－－－")

	Dictionary.AddCodeString('.', "・－・－・－")
	Dictionary.AddCodeString(',', "－－・・－－")
	Dictionary.AddCodeString('?', "・・－－・・")
	Dictionary.AddCodeString('-', "－・・・・－")
	Dictionary.AddCodeString('/', "－・・－・")
	Dictionary.AddCodeString('@', "・－－・－・")
	Dictionary.AddCodeString('(', "－・－－・")
	Dictionary.AddCodeString(')', "－・－－・－")

	/* todo uncomment once jp is handled properly
	Dictionary.AddCodeString('イ',"・－")
	Dictionary.AddCodeString('ロ',"・－・－")
	Dictionary.AddCodeString('ハ',"－・・・")
	Dictionary.AddCodeString('ニ',"－・－・")
	Dictionary.AddCodeString('ホ',"－・・")
	Dictionary.AddCodeString('ヘ',"・")
	Dictionary.AddCodeString('ト',"・・－・・")
	Dictionary.AddCodeString('チ',"・・－・")
	Dictionary.AddCodeString('リ',"－－・")
	Dictionary.AddCodeString('ヌ',"・・・・")
	Dictionary.AddCodeString('ル',"－・－－・")
	Dictionary.AddCodeString('ヲ',"・－－－")
	Dictionary.AddCodeString('ワ',"－・－")
	Dictionary.AddCodeString('カ',"・－・・")
	Dictionary.AddCodeString('ヨ',"－－")
	Dictionary.AddCodeString('タ',"－・")
	Dictionary.AddCodeString('レ',"－－－")
	Dictionary.AddCodeString('ソ',"－－－・")
	Dictionary.AddCodeString('ツ',"・－－・")
	Dictionary.AddCodeString('ネ',"－－・－")
	Dictionary.AddCodeString('ナ',"・－・")
	Dictionary.AddCodeString('ラ',"・・・")
	Dictionary.AddCodeString('ム',"－")
	Dictionary.AddCodeString('ウ',"・・－")
	Dictionary.AddCodeString('ヰ',"・－・・－")
	Dictionary.AddCodeString('ノ',"・・－－")
	Dictionary.AddCodeString('オ',"・－・・・")
	Dictionary.AddCodeString('ク',"・・・－")
	Dictionary.AddCodeString('ヤ',"・－－")
	Dictionary.AddCodeString('マ',"－・・－")
	Dictionary.AddCodeString('ケ',"－・－－")
	Dictionary.AddCodeString('フ',"－－・・")
	Dictionary.AddCodeString('コ',"－－－－")
	Dictionary.AddCodeString('エ',"－・－－－")
	Dictionary.AddCodeString('テ',"・－・－－")
	Dictionary.AddCodeString('ア',"－－・－－")
	Dictionary.AddCodeString('サ',"－・－・－")
	Dictionary.AddCodeString('キ',"－・－・・")
	Dictionary.AddCodeString('ユ',"－・・－－")
	Dictionary.AddCodeString('メ',"－・・・－")
	Dictionary.AddCodeString('ミ',"・・－・－")
	Dictionary.AddCodeString('シ',"－－・－・")
	Dictionary.AddCodeString('ヱ',"・－－・・")
	Dictionary.AddCodeString('ヒ',"－－・・－")
	Dictionary.AddCodeString('モ',"－・・－・")
	Dictionary.AddCodeString('セ',"・－－－・")
	Dictionary.AddCodeString('ス',"－－－・－")
	Dictionary.AddCodeString('ン',"・－・－・")
	Dictionary.AddCodeString('゛',"・・")
	Dictionary.AddCodeString('゜',"・・－－・")
	*/

	// Generate the code for the standard word after the dictionary has been populated
	standardWordCode = append(FromText(StandardWord), WordSpace)

	// Calculate the standard word duration
	standardWordDuration = standardWordCode.duration()

	// Calculate the standard word farnsworth duration
	// (the duration of the spaces between letters and words)
	for _, s := range standardWordCode {
		if s == RuneSpace || s == WordSpace {
			standardWordFarnsworthDuration += s.duration
		}
	}
}
