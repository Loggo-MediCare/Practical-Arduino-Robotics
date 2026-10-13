// Package dimmer ports Chapter10/Dimmer: a non-blocking software-PWM dimmer.
// It toggles an output pin with a period of ~1024us and a duty cycle given in
// the range 0..1023, mirroring the original Arduino library exactly.
package dimmer

import (
	"machine"
	"time"
)

const blinkInterval = 1024 * time.Microsecond

type Dimmer struct {
	pin        machine.Pin
	dutyCycle  time.Duration
	lastBlink  time.Time
	ledOn      bool
}

func New(pin machine.Pin, dutyCycle int) *Dimmer {
	return &Dimmer{
		pin:       pin,
		dutyCycle: time.Duration(dutyCycle) * time.Microsecond,
	}
}

func (d *Dimmer) Begin() {
	d.pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	d.lastBlink = time.Now()
}

func (d *Dimmer) Blink() {
	now := time.Now()
	if now.Sub(d.lastBlink) >= d.onOffDuration() {
		d.blinkTask()
	}
}

func (d *Dimmer) onOffDuration() time.Duration {
	if d.ledOn {
		return d.dutyCycle
	}
	return blinkInterval - d.dutyCycle
}

func (d *Dimmer) blinkTask() {
	if d.ledOn {
		d.pin.Low()
		d.ledOn = false
		d.lastBlink = d.lastBlink.Add(blinkInterval - d.dutyCycle)
	} else {
		if d.dutyCycle > 0 {
			d.pin.High()
		}
		d.ledOn = true
		d.lastBlink = d.lastBlink.Add(d.dutyCycle)
	}
}

// SetDutyCycle sets the duty cycle, constrained to 0..1023 microseconds.
func (d *Dimmer) SetDutyCycle(dutyCycle int) {
	if dutyCycle < 0 {
		dutyCycle = 0
	}
	if dutyCycle > 1023 {
		dutyCycle = 1023
	}
	d.dutyCycle = time.Duration(dutyCycle) * time.Microsecond
}
