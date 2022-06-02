package prosign

import "github.com/bhollier/morse"

var Newline = morse.FromCodeString("・－・－")

// todo add more

func init() {
	morse.Dictionary.Add('\n', Newline)
}
