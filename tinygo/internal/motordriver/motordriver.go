// Package motordriver ports Chapter13/MotorDriver: a simple driver for a
// motor controller with two direction inputs and one PWM enable input.
//
// It is used only by Chapter13/self_balancing, which targets an Arduino
// Mega2560 and drives its two motors from Timer1's two PWM channels (OC1A on
// pin 11 / PB5, OC1B on pin 12 / PB6). TinyGo's machine package does not
// implement the generic PWM abstraction for atmega2560 (only atmega328 and
// atmega1280 do), so this package configures Timer1 directly via
// device/avr, mirroring what the original sketch already did by hand for its
// Timer1 prescaler.
package motordriver

import (
	"device/avr"
	"machine"
)

// Channel selects which of Timer1's two PWM outputs to drive.
type Channel uint8

const (
	ChannelA Channel = iota // OC1A, digital pin 11 (PB5)
	ChannelB                // OC1B, digital pin 12 (PB6)
)

var timer1Configured bool

// configureTimer1 sets up Timer1 for 8 bit Phase Correct PWM (TOP = 0xFF) with
// no prescaling, giving a ~31.37kHz PWM frequency on both OC1A and OC1B. This
// matches the prescaler setting the original sketch pokes into TCCR1B by
// hand.
func configureTimer1() {
	if timer1Configured {
		return
	}
	// WGM10 = 1 (with WGM12=WGM13=0): Phase Correct PWM, 8 bit, TOP=0xFF.
	// COM1A1/COM1B1 = 1: non-inverting output on OC1A/OC1B.
	avr.TCCR1A.Set(avr.TCCR1A_WGM10 | avr.TCCR1A_COM1A1 | avr.TCCR1A_COM1B1)
	// CS10 = 1: no prescaling.
	avr.TCCR1B.Set(avr.TCCR1B_CS10)
	avr.OCR1AH.Set(0)
	avr.OCR1AL.Set(0)
	avr.OCR1BH.Set(0)
	avr.OCR1BL.Set(0)
	timer1Configured = true
}

type MotorDriver struct {
	channel  Channel
	dirPinA  machine.Pin
	dirPinB  machine.Pin
	pwm      int
	maxPwm   int
}

func New(channel Channel, dirPinA, dirPinB machine.Pin) *MotorDriver {
	return &MotorDriver{
		channel: channel,
		dirPinA: dirPinA,
		dirPinB: dirPinB,
		maxPwm:  255,
	}
}

func (m *MotorDriver) Begin() {
	configureTimer1()
	m.dirPinA.Configure(machine.PinConfig{Mode: machine.PinOutput})
	m.dirPinB.Configure(machine.PinConfig{Mode: machine.PinOutput})
	m.dirPinA.Low()
	m.dirPinB.Low()
	switch m.channel {
	case ChannelA:
		machine.PB5.Configure(machine.PinConfig{Mode: machine.PinOutput})
	case ChannelB:
		machine.PB6.Configure(machine.PinConfig{Mode: machine.PinOutput})
	}
	m.setDuty(0)
}

func constrain(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func (m *MotorDriver) SetPwm(pwm int) {
	m.pwm = constrain(pwm, -m.maxPwm, m.maxPwm)
	dir := m.pwm > 0
	m.dirPinA.Set(dir)
	m.dirPinB.Set(!dir)
	m.setDuty(abs(m.pwm))
}

func (m *MotorDriver) setDuty(duty int) {
	switch m.channel {
	case ChannelA:
		avr.OCR1AL.Set(uint8(duty))
	case ChannelB:
		avr.OCR1BL.Set(uint8(duty))
	}
}

func (m *MotorDriver) AdjustPwm(val int) {
	m.SetPwm(m.pwm + val)
}

func (m *MotorDriver) SetMaxPwm(maxPwm int) {
	m.maxPwm = maxPwm
}

func (m *MotorDriver) GetPwm() int {
	return m.pwm
}
