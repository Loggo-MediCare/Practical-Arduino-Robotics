// Port of Chapter10/Dimmer_example/Dimmer_example.ino.
package main

import (
	"machine"
	"time"

	"practicalarduinorobotics/tinygo/internal/dimmer"
)

func main() {
	myDimmer := dimmer.New(machine.LED, 500)
	myDimmer.Begin()

	for {
		// Generate saw-tooth brightness profile.
		dutyCycle := int(time.Now().UnixMilli() % 1024)
		myDimmer.SetDutyCycle(dutyCycle)
		myDimmer.Blink()
	}
}
