// Port of Chapter06/Blinker/examples/minimal_example.ino.
package main

import (
	"machine"

	"practicalarduinorobotics/tinygo/internal/blinker"
)

func main() {
	myBlinker := blinker.New(machine.LED, 50)
	myBlinker.Begin()
	for {
		myBlinker.Blink()
	}
}
