// Package quadencoder ports Chapter13/QuadEncoder: a quadrature encoder
// reader with a sliding-window velocity estimate.
//
// Unlike the original Arduino library, callers do not need to write global
// ISR wrapper functions: Begin attaches the interrupt handlers itself using
// machine.Pin.SetInterrupt (Uno / atmega328 only).
package quadencoder

import (
	"machine"
	"time"
)

const windowSize = 80
const samplingInterval = 4000 * time.Microsecond

type QuadEncoder struct {
	pinA, pinB machine.Pin

	position     int32
	lastPosition int32
	stateA       bool
	stateB       bool

	oldestPosition        int32
	lastVelocitySampling  time.Time
	velocityCps           float32
	velocityRaw           float32
	positionWindow        [windowSize]int32
	samplingIndex         uint8
}

func New(pinA, pinB machine.Pin) *QuadEncoder {
	return &QuadEncoder{pinA: pinA, pinB: pinB}
}

func (q *QuadEncoder) PinA() machine.Pin { return q.pinA }
func (q *QuadEncoder) PinB() machine.Pin { return q.pinB }
func (q *QuadEncoder) Position() int32   { return q.position }
func (q *QuadEncoder) Velocity() float32 { return q.velocityCps }
func (q *QuadEncoder) RawVelocity() float32 { return q.velocityRaw }

// Begin configures the encoder pins and attaches the interrupt handlers.
func (q *QuadEncoder) Begin() {
	q.pinA.Configure(machine.PinConfig{Mode: machine.PinInput})
	q.pinB.Configure(machine.PinConfig{Mode: machine.PinInput})
	q.stateA = q.pinA.Get()
	q.stateB = q.pinB.Get()
	q.lastVelocitySampling = time.Now()

	q.pinA.SetInterrupt(machine.PinToggle, func(machine.Pin) {
		q.changeIsrA()
	})
	q.pinB.SetInterrupt(machine.PinToggle, func(machine.Pin) {
		q.changeIsrB()
	})
}

func (q *QuadEncoder) changeIsrA() {
	newState := q.pinA.Get()
	if newState == q.stateA {
		return
	}
	q.stateA = newState
	if newState {
		q.risingA()
	} else {
		q.fallingA()
	}
}

func (q *QuadEncoder) changeIsrB() {
	newState := q.pinB.Get()
	if newState == q.stateB {
		return
	}
	q.stateB = newState
	if newState {
		q.risingB()
	} else {
		q.fallingB()
	}
}

// Update refreshes the sliding-window velocity estimate. Call as often as
// possible from the main loop; it must never block.
func (q *QuadEncoder) Update() {
	now := time.Now()
	if now.Sub(q.lastVelocitySampling) >= samplingInterval {
		q.lastVelocitySampling = q.lastVelocitySampling.Add(samplingInterval)
		q.oldestPosition = q.positionWindow[q.samplingIndex]
		q.positionWindow[q.samplingIndex] = q.position
		q.samplingIndex = (q.samplingIndex + 1) % windowSize
		q.velocityCps = 1000000.0 * float32(q.position-q.oldestPosition) /
			float32(samplingInterval.Microseconds()*windowSize)
		q.velocityRaw = 1000000.0 * float32(q.position-q.lastPosition) /
			float32(samplingInterval.Microseconds())
		q.lastPosition = q.position
	}
}

func (q *QuadEncoder) risingA() {
	if !q.stateB {
		q.position++
	} else {
		q.position--
	}
}

func (q *QuadEncoder) fallingA() {
	if !q.stateB {
		q.position--
	} else {
		q.position++
	}
}

func (q *QuadEncoder) risingB() {
	if !q.stateA {
		q.position--
	} else {
		q.position++
	}
}

func (q *QuadEncoder) fallingB() {
	if !q.stateA {
		q.position++
	} else {
		q.position--
	}
}
