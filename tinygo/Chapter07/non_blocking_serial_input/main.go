// Port of Chapter07/non_blocking_serial_input/non_blocking_serial_input.ino.
package main

import (
	"machine"
	"time"

	"practicalarduinorobotics/tinygo/internal/blinker"
	"practicalarduinorobotics/tinygo/internal/numfmt"
)

type parserState uint8

const (
	waitingForData parserState = iota
	processingData
)

const rxWaitTime = 10 * time.Millisecond

func main() {
	machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})
	myBlinker := blinker.New(machine.LED, 50)
	myBlinker.Begin()

	state := waitingForData
	var lastRxDataTimestamp time.Time

	for {
		myBlinker.Blink()

		switch state {
		case waitingForData:
			if machine.Serial.Buffered() > 0 {
				lastRxDataTimestamp = time.Now()
				state = processingData
			}
		case processingData:
			if time.Since(lastRxDataTimestamp) >= rxWaitTime {
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
				state = waitingForData
			}
		}
	}
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

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
