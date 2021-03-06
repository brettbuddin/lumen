package midi

import "github.com/brettbuddin/shaden/unit"

const (
	sampleRate = 44100
	frameSize  = 256
)

func newUnitConfig(v map[string]interface{}) unit.Config {
	return unit.Config{
		SampleRate: sampleRate,
		FrameSize:  frameSize,
	}
}
