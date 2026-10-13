// Port of Chapter07/non_blocking_multiple_inputs/non_blocking_multiple_inputs.ino.
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
	blinkerA := blinker.New(machine.D11, 50)
	blinkerB := blinker.New(machine.D12, 50)
	blinkerC := blinker.New(machine.D13, 50)
	blinkerA.Begin()
	blinkerB.Begin()
	blinkerC.Begin()

	state := waitingForData
	var lastRxDataTimestamp time.Time

	for {
		blinkerA.Blink()
		blinkerB.Blink()
		blinkerC.Blink()

		switch state {
		case waitingForData:
			if machine.Serial.Buffered() > 0 {
				lastRxDataTimestamp = time.Now()
				state = processingData
			}
		case processingData:
			if time.Since(lastRxDataTimestamp) >= rxWaitTime {
				// First character is the command key.
				key, _ := machine.Serial.ReadByte()
				if machine.Serial.Buffered() > 0 {
					next, _ := machine.Serial.ReadByte()
					if isDigit(next) {
						value := readInt(next)
						switch key {
						case 'a':
							machine.Serial.Write([]byte("Setting blink interval a to "))
							machine.Serial.Write([]byte(numfmt.Int(int32(value)) + "\r\n"))
							blinkerA.SetBlinkInterval(value)
						case 'b':
							machine.Serial.Write([]byte("Setting blink interval b to "))
							machine.Serial.Write([]byte(numfmt.Int(int32(value)) + "\r\n"))
							blinkerB.SetBlinkInterval(value)
						case 'c':
							machine.Serial.Write([]byte("Setting blink interval c to "))
							machine.Serial.Write([]byte(numfmt.Int(int32(value)) + "\r\n"))
							blinkerC.SetBlinkInterval(value)
						default:
							machine.Serial.Write([]byte("Warning: Invalid key received.\r\n"))
						}
					} else {
						machine.Serial.Write([]byte("Warning: Invalid value received.\r\n"))
					}
				} else {
					machine.Serial.Write([]byte("Warning: Invalid value received.\r\n"))
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
