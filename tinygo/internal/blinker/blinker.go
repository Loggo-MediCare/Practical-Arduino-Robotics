// Package blinker ports Chapter06/Blinker: a non-blocking LED blinker.
package blinker

import (
	"machine"
	"time"
)

type Blinker struct {
	pin           machine.Pin
	blinkInterval time.Duration
	lastBlink     time.Time
	state         bool
}

func New(pin machine.Pin, blinkIntervalMs int) *Blinker {
	return &Blinker{
		pin:           pin,
		blinkInterval: time.Duration(blinkIntervalMs) * time.Millisecond,
		lastBlink:     time.Now(),
	}
}

func (b *Blinker) Begin() {
	b.pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
}

func (b *Blinker) Blink() {
	now := time.Now()
	if now.Sub(b.lastBlink) >= b.blinkInterval {
		b.lastBlink = b.lastBlink.Add(b.blinkInterval)
		b.state = !b.state
		b.pin.Set(b.state)
	}
}

func (b *Blinker) SetBlinkInterval(intervalMs int) {
	b.blinkInterval = time.Duration(intervalMs) * time.Millisecond
}
