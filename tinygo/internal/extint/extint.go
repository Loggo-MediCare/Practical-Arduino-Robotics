// Package extint attaches interrupt handlers to the Arduino Mega2560's
// hardware external interrupt lines INT0-INT5 (digital pins 21, 20, 19, 18,
// 2, 3).
//
// TinyGo's machine package does not implement Pin.SetInterrupt for
// atmega2560 (only atmega328), so this package configures EICRA/EICRB/EIMSK
// directly, following the same pattern machine_atmega328.go uses for its own
// SetInterrupt: register a Go callback via runtime/interrupt.New, then
// unmask the interrupt at the peripheral level.
package extint

import (
	"device/avr"
	"runtime/interrupt"
	"runtime/volatile"
)

// Mode selects which edge(s) trigger the interrupt, matching Arduino's
// RISING, FALLING, and CHANGE constants.
type Mode uint8

const (
	Rising Mode = iota
	Falling
	Change
)

// Line identifies one of the Mega2560's 6 usable external interrupt lines.
type Line uint8

const (
	INT0 Line = iota // digital pin 21
	INT1              // digital pin 20
	INT2              // digital pin 19
	INT3              // digital pin 18
	INT4              // digital pin 2
	INT5              // digital pin 3
)

var callbacks [6]func()

// Attach configures line to trigger on the given edge mode and calls fn from
// interrupt context whenever it fires.
func Attach(line Line, mode Mode, fn func()) {
	callbacks[line] = fn

	var isc1, isc0 bool
	switch mode {
	case Rising:
		isc1, isc0 = true, true
	case Falling:
		isc1, isc0 = true, false
	case Change:
		isc1, isc0 = false, true
	}

	switch line {
	case INT0:
		setSenseBits(avr.EICRA, 0, isc1, isc0)
		interrupt.New(avr.IRQ_INT0, handleINT0)
		avr.EIMSK.SetBits(avr.EIMSK_INT0)
	case INT1:
		setSenseBits(avr.EICRA, 1, isc1, isc0)
		interrupt.New(avr.IRQ_INT1, handleINT1)
		avr.EIMSK.SetBits(avr.EIMSK_INT1)
	case INT2:
		setSenseBits(avr.EICRA, 2, isc1, isc0)
		interrupt.New(avr.IRQ_INT2, handleINT2)
		avr.EIMSK.SetBits(avr.EIMSK_INT2)
	case INT3:
		setSenseBits(avr.EICRA, 3, isc1, isc0)
		interrupt.New(avr.IRQ_INT3, handleINT3)
		avr.EIMSK.SetBits(avr.EIMSK_INT3)
	case INT4:
		setSenseBits(avr.EICRB, 0, isc1, isc0)
		interrupt.New(avr.IRQ_INT4, handleINT4)
		avr.EIMSK.SetBits(avr.EIMSK_INT4)
	case INT5:
		setSenseBits(avr.EICRB, 1, isc1, isc0)
		interrupt.New(avr.IRQ_INT5, handleINT5)
		avr.EIMSK.SetBits(avr.EIMSK_INT5)
	}
}

// setSenseBits sets the two ISCn1:ISCn0 sense-control bits at bit position
// index*2 within an EICRA/EICRB register.
func setSenseBits(reg *volatile.Register8, index uint8, isc1, isc0 bool) {
	mask := uint8(0b11) << (index * 2)
	reg.ClearBits(mask)
	var value uint8
	if isc1 {
		value |= 0b10 << (index * 2)
	}
	if isc0 {
		value |= 0b01 << (index * 2)
	}
	reg.SetBits(value)
}

func handleINT0(interrupt.Interrupt) {
	if callbacks[INT0] != nil {
		callbacks[INT0]()
	}
}
func handleINT1(interrupt.Interrupt) {
	if callbacks[INT1] != nil {
		callbacks[INT1]()
	}
}
func handleINT2(interrupt.Interrupt) {
	if callbacks[INT2] != nil {
		callbacks[INT2]()
	}
}
func handleINT3(interrupt.Interrupt) {
	if callbacks[INT3] != nil {
		callbacks[INT3]()
	}
}
func handleINT4(interrupt.Interrupt) {
	if callbacks[INT4] != nil {
		callbacks[INT4]()
	}
}
func handleINT5(interrupt.Interrupt) {
	if callbacks[INT5] != nil {
		callbacks[INT5]()
	}
}
