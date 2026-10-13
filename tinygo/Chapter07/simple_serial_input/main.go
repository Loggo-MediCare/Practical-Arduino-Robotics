// Port of Chapter07/simple_serial_input/simple_serial_input.ino.
package main

import (
	"machine"
	"time"

	"practicalarduinorobotics/tinygo/internal/blinker"
	"practicalarduinorobotics/tinygo/internal/numfmt"
)

func main() {
	machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})
	myBlinker := blinker.New(machine.LED, 50)
	myBlinker.Begin()

	for {
		myBlinker.Blink()
		parseInterval(myBlinker)
	}
}

// parseInterval reads Serial input and interprets it as a blink interval.
func parseInterval(myBlinker *blinker.Blinker) {
	if machine.Serial.Buffered() > 0 {
		// Wait to receive the entire message.
		time.Sleep(10 * time.Millisecond)
		b, _ := machine.Serial.ReadByte()
		if isDigit(b) {
			newInterval := readInt(b)
			machine.Serial.Write([]byte("Setting blink interval to "))
			machine.Serial.Write([]byte(numfmt.Int(int32(newInterval)) + "\r\n"))
			myBlinker.SetBlinkInterval(newInterval)
		} else {
			machine.Serial.Write([]byte("Warning: Invalid data received.\r\n"))
		}
		for machine.Serial.Buffered() > 0 {
			machine.Serial.ReadByte()
		}
	}
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

// readInt parses an integer starting with the already-read digit first,
// consuming further buffered digits from Serial.
func readInt(first byte) int {
	value := int(first - '0')
	for machine.Serial.Buffered() > 0 {
		b, _ := machine.Serial.ReadByte()
		if !isDigit(b) {
			break
		}
		value = value*10 + int(b-'0')
	}
	return value
}
