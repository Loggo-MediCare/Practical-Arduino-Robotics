// Port of Chapter02/GPIO_button/GPIO_button.ino.
package main

import (
	"machine"
	"time"

	"practicalarduinorobotics/tinygo/internal/numfmt"
)

const buttonPin = machine.D2

var numButtonEvents int32

func main() {
	machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	buttonPin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	buttonPin.SetInterrupt(machine.PinToggle, func(machine.Pin) {
		numButtonEvents++
	})

	for {
		buttonState := buttonPin.Get()
		// Invert input to turn LED ON when input is LOW.
		led.Set(!buttonState)

		if time.Now().UnixMilli()%100 == 0 {
			machine.Serial.Write([]byte(numfmt.Int(numButtonEvents) + "\r\n"))
		}
	}
}
