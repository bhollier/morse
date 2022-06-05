// Package play contains methods for playing Morse Code audio
package play

import (
	"fmt"
	"github.com/bhollier/morse"
	"github.com/bhollier/morse/internal/buffer"
	"github.com/faiface/beep"
	"github.com/faiface/beep/generators"
	"io"
	"time"
)

type streamer struct {
	fadeStreamer       *fadeStreamer
	sampleRate         beep.SampleRate
	wpm, farnsworthWPM uint
	morseReader        morse.Reader
	overflow           buffer.Overflow[[2]float64]
	err                error
}

// Read the signals 1 by 1
const signalBufferSize = 1

const fadeDuration = time.Millisecond * 40

// MorseStreamer creates a beep.Streamer for streaming
// morse.Code from the given morse.Reader as audio
func MorseStreamer(sr beep.SampleRate, freq int, wpm, farnsworthWPM uint, r morse.Reader) (beep.Streamer, error) {
	if farnsworthWPM > wpm {
		return nil, fmt.Errorf("farnswordWPM (%d) > wpm (%d)", farnsworthWPM, wpm)
	}

	sinToneStreamer, err := generators.SinTone(sr, freq)
	if err != nil {
		return nil, err
	}
	return &streamer{
		fadeStreamer:  newFadeStreamer(sinToneStreamer, sr),
		sampleRate:    sr,
		wpm:           wpm,
		farnsworthWPM: farnsworthWPM,
		morseReader:   r,
	}, nil
}

func (s *streamer) Stream(samples [][2]float64) (n int, ok bool) {
	// First, try to empty the overflow from the last read
	n = s.overflow.Empty(samples)
	samples = samples[n:]

	signals := make([]morse.Signal, signalBufferSize)

	// While there is space in p and there are morse signals to stream
	for len(samples) > 0 && s.err == nil {
		signalsRead, err := s.morseReader.Read(signals)
		if err != nil {
			s.err = err
		}

		// If we got some signals from the morse reader
		if signalsRead > 0 {
			for _, signal := range signals[:signalsRead] {
				numSamples := s.sampleRate.N(signal.Duration(s.wpm, s.farnsworthWPM))
				signalSamples := make([][2]float64, numSamples)

				if signal.Audible() {
					s.fadeStreamer.FadeInFor(fadeDuration)
				} else {
					s.fadeStreamer.FadeOutFor(fadeDuration)
				}

				_, ok := s.fadeStreamer.Stream(signalSamples)
				if !ok {
					s.err = s.fadeStreamer.Err()
				}

				// Copy the signal samples into samples (with the remaining going into the buffer)
				samplesCopied := s.overflow.Copy(samples, signalSamples)
				samples = samples[samplesCopied:]
				n += samplesCopied
			}

			// If we got no signals, but there's no error
			// (this can happen e.g. if the morse reader is a
			// NonBlockingChannelReader)
		} else if s.err == nil {
			// Fade out and the rest of the signals can be silent
			s.fadeStreamer.FadeOutFor(fadeDuration)
			samplesCopied, ok := s.fadeStreamer.Stream(samples)
			if !ok {
				s.err = s.fadeStreamer.Err()
			}

			samples = samples[samplesCopied:]
			n += samplesCopied

			// If we reached EOF, we still want to add a fade,
			// otherwise the signal cuts off very messily
		} else if s.err == io.EOF {
			// Create enough samples for a final rune space
			numSamples := s.sampleRate.N(morse.RuneSpace.Duration(s.wpm, s.farnsworthWPM))
			fadeSamples := make([][2]float64, numSamples)
			s.fadeStreamer.FadeOutFor(fadeDuration)

			_, ok := s.fadeStreamer.Stream(fadeSamples)
			if !ok {
				s.err = s.fadeStreamer.Err()
			}

			// Copy the fade samples into samples (with the remaining going into the buffer)
			samplesCopied := s.overflow.Copy(samples, fadeSamples)
			samples = samples[samplesCopied:]
			n += samplesCopied
		}
	}

	return n, n > 0 && (s.err == nil || s.err == io.EOF)
}

func (s *streamer) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}
