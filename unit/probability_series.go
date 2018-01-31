package unit

import (
	"math"
	"math/rand"

	"buddin.us/shaden/dsp"
)

const readBit = 65535

func newProbabilitySeries(name string, _ Config) (*Unit, error) {
	io := NewIO()
	return NewUnit(io, name, &probabilitySeries{
		clock:     io.NewIn("clock", dsp.Float64(-1)),
		length:    io.NewIn("length", dsp.Float64(8)),
		lock:      io.NewIn("lock", dsp.Float64(0)),
		gates:     make([]float64, 16),
		values:    make([]float64, 16),
		gate:      io.NewOut("gate"),
		value:     io.NewOut("value"),
		lastClock: -1,
	}), nil
}

type probabilitySeries struct {
	clock, lock, length *In
	gate, value         *Out
	gates, values       []float64

	idx       int
	lastClock float64
}

func (s *probabilitySeries) ProcessSample(i int) {
	var (
		clock     = s.clock.Read(i)
		length    = dsp.Clamp(s.length.ReadSlow(i, ident), 2, 16)
		lengthInt = int(length)
		lock      = s.lock.ReadSlow(i, ident)
	)

	if isTrig(s.lastClock, clock) {
		var (
			lastGate, lastValue = s.gates[lengthInt-1], s.values[lengthInt-1]
			data                = rand.Float64()
		)
		for i := 0; i < lengthInt; i++ {
			s.gates[i], lastGate = lastGate, s.gates[i]
			s.values[i], lastValue = lastValue, s.values[i]
		}
		if lock != 1 && data > lock {
			s.gates[0] = round(data)
			s.values[0] = rand.Float64()
		}
	}

	s.gate.Write(i, s.gates[lengthInt-1])
	s.value.Write(i, s.values[lengthInt-1])

	s.lastClock = clock
}

func round(v float64) float64 {
	t := math.Trunc(v)
	if math.Abs(v-t) >= 0.5 {
		return t + math.Copysign(1, v)
	}
	return t
}
