// Port of Chapter06/OOP_blinker/OOP_blinker.ino.
package main

import (
	"machine"
	"time"
)

type Blinker struct {
	ledPin        machine.Pin
	blinkInterval time.Duration
	lastBlinkTime time.Time
	state         bool
}

func NewBlinker(ledPin machine.Pin, blinkIntervalMs int) *Blinker {
	return &Blinker{
		ledPin:        ledPin,
		blinkInterval: time.Duration(blinkIntervalMs) * time.Millisecond,
		lastBlinkTime: time.Now(),
	}
}

func (b *Blinker) Begin() {
	b.ledPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
}

func (b *Blinker) Blink() {
	if time.Since(b.lastBlinkTime) >= b.blinkInterval {
		b.lastBlinkTime = b.lastBlinkTime.Add(b.blinkInterval)
		b.blinkTask()
	}
}

func (b *Blinker) SetBlinkInterval(blinkIntervalMs int) {
	b.blinkInterval = time.Duration(blinkIntervalMs) * time.Millisecond
}

func (b *Blinker) blinkTask() {
	b.state = !b.state
	b.ledPin.Set(b.state)
}

func main() {
	myBlinker := NewBlinker(machine.LED, 200)
	myBlinker.Begin()
	for {
		myBlinker.Blink()
	}
}
