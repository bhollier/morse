package morse

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewSignal(t *testing.T) {
	a := assert.New(t)

	// First, test the constants
	a.Equal(Dit, NewSignal(true, 1))
	a.Equal(Dah, NewSignal(true, 3))
	a.Equal(SignalSpace, NewSignal(false, 1))
	a.Equal(RuneSpace, NewSignal(false, 3))
	a.Equal(WordSpace, NewSignal(false, 7))

	a.Equal(true, Dit.Audible())
	a.Equal(true, Dah.Audible())
	a.Equal(false, SignalSpace.Audible())
	a.Equal(false, RuneSpace.Audible())
	a.Equal(false, WordSpace.Audible())

	a.Equal(uint(1), Dit.DitDuration())
	a.Equal(uint(3), Dah.DitDuration())
	a.Equal(uint(1), SignalSpace.DitDuration())
	a.Equal(uint(3), RuneSpace.DitDuration())
	a.Equal(uint(7), WordSpace.DitDuration())

	// Now check custom signals, just in case
	maxSignal := NewSignal(false, MaxSignalDuration)
	a.Equal(false, maxSignal.Audible())
	a.Equal(uint(MaxSignalDuration), maxSignal.DitDuration())

	maxSignal = NewSignal(true, MaxSignalDuration)
	a.Equal(true, maxSignal.Audible())
	a.Equal(uint(MaxSignalDuration), maxSignal.DitDuration())

	a.Panics(func() {
		NewSignal(true, MaxSignalDuration+1)
	})
	a.Panics(func() {
		NewSignal(true, 0)
	})
}

// This test should always pass (as internally []byte and []Signal
// should be the same type), but is here just in case unsafe.Pointer
// is changed somehow
func TestConvertByteArraySignalArray(t *testing.T) {
	a := assert.New(t)

	signals := []Signal{Dit, Dah, SignalSpace, RuneSpace, WordSpace}
	bytes := []byte{byte(Dit), byte(Dah), byte(SignalSpace), byte(RuneSpace), byte(WordSpace)}

	a.Equal(bytes, signalArrayToByteArray(signals))
	a.Equal(signals, byteArrayToSignalArray(bytes))
}

func TestSignal_Duration(t *testing.T) {
	a := assert.New(t)

	standardWordDuration := append(FromText(StandardWord), WordSpace).Duration(20, 0)
	a.Equal(time.Minute, (standardWordDuration * 20).Round(time.Second))

	standardWordFarnsworthDuration := append(FromText(StandardWord), WordSpace).Duration(20, 15)
	a.Equal(time.Minute, (standardWordFarnsworthDuration * 15).Round(time.Second))
}
