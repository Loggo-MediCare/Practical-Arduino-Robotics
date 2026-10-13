// Port of Chapter05/quadrature_encoder/quadrature_encoder.ino.
package main

import (
	"machine"
	"time"

	"practicalarduinorobotics/tinygo/internal/numfmt"
)

const pinA = machine.D2
const pinB = machine.D3

var encoderPosition int32

func main() {
	machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})
	pinA.Configure(machine.PinConfig{Mode: machine.PinInput})
	pinB.Configure(machine.PinConfig{Mode: machine.PinInput})

	// The original sketch dynamically swaps between RISING- and
	// FALLING-only interrupt handlers to implement quadrature decoding.
	// TinyGo's Pin.SetInterrupt fires on both edges (PinToggle) and delivers
	// the pin, so we read the pin's current level to tell which edge just
	// happened and apply the same decoding table.
	pinA.SetInterrupt(machine.PinToggle, func(p machine.Pin) {
		if p.Get() {
			// Rising edge on A.
			if !pinB.Get() {
				encoderPosition++
			} else {
				encoderPosition--
			}
		} else {
			// Falling edge on A.
			if !pinB.Get() {
				encoderPosition--
			} else {
				encoderPosition++
			}
		}
	})
	pinB.SetInterrupt(machine.PinToggle, func(p machine.Pin) {
		if p.Get() {
			// Rising edge on B.
			if !pinA.Get() {
				encoderPosition--
			} else {
				encoderPosition++
			}
		} else {
			// Falling edge on B.
			if !pinA.Get() {
				encoderPosition++
			} else {
				encoderPosition--
			}
		}
	})

	for {
		machine.Serial.Write([]byte("Position:" + numfmt.Int(encoderPosition) + "\r\n"))
		time.Sleep(50 * time.Millisecond)
	}
}
