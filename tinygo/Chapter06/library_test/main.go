// Port of Chapter06/library_test/library_test.ino.
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
