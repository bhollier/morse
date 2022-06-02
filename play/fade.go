package play

import (
	"github.com/faiface/beep"
	"math"
	"time"
)

const minGain = 0.00001
const maxGain = 1.0

type fadeStreamer struct {
	beep.Streamer
	sr beep.SampleRate

	currentGain float64

	initialGain    float64
	endGain        float64
	gainGrowthRate float64
}

func newFadeStreamer(s beep.Streamer, sr beep.SampleRate) *fadeStreamer {
	return &fadeStreamer{
		Streamer:    s,
		sr:          sr,
		currentGain: minGain,
	}
}

func (s *fadeStreamer) calcGain(t int) float64 {
	return s.initialGain * math.Pow(math.E, float64(t)*s.gainGrowthRate)
}

func (s *fadeStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	n, ok = s.Streamer.Stream(samples)
	if ok && n > 0 {
		for i := range samples[:n] {
			if s.gainGrowthRate != 0 {
				s.currentGain = s.calcGain(i)
				if (s.gainGrowthRate > 0 && s.currentGain > s.endGain) ||
					(s.gainGrowthRate < 0 && s.currentGain < s.endGain) {
					s.gainGrowthRate = 0
					s.currentGain = s.endGain
				}
			}
			samples[i][0] *= s.currentGain
			samples[i][1] *= s.currentGain
		}
	}
	return
}

func (s *fadeStreamer) Fade(toGain float64, d time.Duration) {
	s.initialGain, s.endGain = s.currentGain, toGain
	s.gainGrowthRate = math.Log(toGain/s.initialGain) / float64(s.sr.N(d))
}

func (s *fadeStreamer) FadeInFor(d time.Duration) {
	s.Fade(maxGain, d)
}

func (s *fadeStreamer) FadeOutFor(d time.Duration) {
	s.Fade(minGain, d)
}
