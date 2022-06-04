// Package prosign contains definitions for common Morse prosigns.
// Importing the package automatically adds some prosigns to morse.Dictionary
package prosign

import "github.com/bhollier/morse"

// The following prosigns can't (currently) be translated into human text so aren't in the dictionary:

var (
	// ThisIsFrom Used to precede the name or other identification of the calling station (Morse abbreviation).
	ThisIsFrom = morse.JoinLetters(morse.D, morse.E)

	// UnknownStation Used for directional signaling lights, but not in radiotelegraphy.
	UnknownStation = morse.JoinNoSpace(morse.A, morse.A)

	// NothingHeard General-purpose response to any request or inquiry for which the answer is "nothing" or "none" or
	// "not available" (Morse abbr.). Also means "I have no messages for you."
	NothingHeard = morse.JoinLetters(morse.N, morse.I, morse.L)

	// Rodger Means the last transmission has been received, but does not indicate the message was understood or will
	// be complied with.
	Rodger = morse.R

	// Over Invitation to transmit after terminating the call signal. (e.g. ・－・－・ －・－).
	Over = morse.K

	// Out End of transmission / End of message / End of telegram. (Same as EC "end copy", and character +.)
	Out = morse.JoinNoSpace(morse.A, morse.R)

	// Closing Announcing station shutdown (Morse abbr.).
	Closing = morse.JoinLetters(morse.C, morse.L)

	// Calling General call to any station (Morse abbr.).
	Calling = morse.JoinLetters(morse.C, morse.Q)

	// CallingFor General call to two or more specified stations (Morse abbr.).
	CallingFor = morse.JoinLetters(morse.C, morse.P)

	// Who What is the name or identity signal of your station? (Morse abbr.).
	Who = morse.JoinLetters(morse.C, morse.S)

	// Wait "I must pause for a few minutes." Also means
	// "I am engaged in a contact with another station [that you may not hear]; please wait quietly."
	Wait = morse.JoinNoSpace(morse.A, morse.S)

	// WaitOut I must pause for more than a few minutes.
	WaitOut = morse.JoinLetters(Wait, Out)

	// Verified Message is verified.
	Verified = morse.JoinNoSpace(morse.V, morse.E)

	// WordAfter (Morse abbr.)
	WordAfter = morse.JoinLetters(morse.W, morse.A)

	// WordBefore (Morse abbr.)
	WordBefore = morse.JoinLetters(morse.W, morse.B)

	// AllAfter The portion of the message to which I refer is all that follows the text ... (Morse abbr.)
	AllAfter = morse.JoinLetters(morse.A, morse.A)

	// AllBefore The portion of the message to which I refer is all that precedes the text ... (Morse abbr.)
	AllBefore = morse.JoinLetters(morse.A, morse.B)

	// AllBetween The portion of the message to which I refer is all that falls between ... and ... (Morse abbr.)
	AllBetween = morse.JoinLetters(morse.A, morse.B)

	// SayAgain When standing alone, a note of interrogation or request for repetition of a transmission not understood.
	// When ? is placed after a coded signal, modifies the code to be a question or request.
	SayAgain = morse.JoinNoSpace(morse.U, morse.D)

	// Interrogative Military replacement for the ? prosign; equivalent to Spanish ¿ punctuation mark. When placed
	// before a signal, modifies the signal to be a question/request.
	Interrogative = morse.JoinNoSpace(morse.I, morse.N, morse.T)

	// Correction Preceding text was in error. The following is the corrected text.
	Correction = morse.JoinNoSpace(morse.H, morse.H)

	// Correct Answer to prior question is "yes". (Morse abbr.)
	Correct = morse.C

	// Negative Answer to prior question is "no". (Morse abbr.)
	Negative = morse.N

	// Wrong Your last transmission was wrong. The correct version is ...
	Wrong = morse.JoinLetters(morse.Z, morse.W, morse.F)

	// DisregardThisTransmission The entire message just sent is in error, disregard it.
	DisregardThisTransmission = morse.JoinLetters(Correction, Out)

	// TimeIs The following is the correct UTC in HHMM 24-hour format
	TimeIs = morse.JoinLetters(morse.Q, morse.T, morse.R)

	// RequestTimeCheck Time-check request. / What is the correct time?
	RequestTimeCheck = morse.JoinLetters(TimeIs, morse.QuestionMark)

	// Break Start new section of message.
	Break = morse.JoinNoSpace(morse.B, morse.T)

	// BreakIn Signal used to interrupt a transmission already in progress (Morse abbr.)
	BreakIn = morse.JoinLetters(morse.B, morse.K)

	// Attention Message begins / Start of work / New message
	Attention = morse.JoinNoSpace(morse.K, morse.A)

	// Acknowledge Message received (Morse abbr.).
	Acknowledge = morse.JoinLetters(morse.C, morse.F, morse.M)

	// WeatherIs Weather report follows (Morse abbr.).
	WeatherIs = morse.JoinLetters(morse.W, morse.X)
)

// The following are prosigns in the dictionary:

var (
	Newline = morse.JoinNoSpace(morse.A, morse.A)
)

// todo add more

var Prosigns = []morse.Code{
	ThisIsFrom, UnknownStation, NothingHeard, Rodger, Over, Closing, Calling, CallingFor, Who, Wait, WaitOut, Verified,
	WordAfter, WordBefore, AllAfter, AllBefore, AllAfter, AllBetween, SayAgain, Interrogative, Correction, Correct,
	Negative, Wrong, DisregardThisTransmission, TimeIs, RequestTimeCheck, Break, BreakIn, Attention, Acknowledge,
	WeatherIs, Newline,
}

func init() {
	morse.Dictionary.Add('\n', Newline)
}
