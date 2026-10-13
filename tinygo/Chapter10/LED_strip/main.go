// Port of Chapter10/LED_strip/LED_strip.ino.
package main

import (
	"machine"
	"time"

	"practicalarduinorobotics/tinygo/internal/dimmer"
)

type parserState uint8

const (
	waitingForData parserState = iota
	processingData
)

const rxWaitTime = 10 * time.Millisecond

func main() {
	// Transistor for red/green/blue LEDs connected to pins 2, 3, 4.
	redDimmer := dimmer.New(machine.D2, 100)
	greenDimmer := dimmer.New(machine.D3, 100)
	blueDimmer := dimmer.New(machine.D4, 100)
	redDimmer.Begin()
	greenDimmer.Begin()
	blueDimmer.Begin()

	machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})

	state := waitingForData
	var lastRxDataTimestamp time.Time

	for {
		redDimmer.Blink()
		greenDimmer.Blink()
		blueDimmer.Blink()

		switch state {
		case waitingForData:
			if machine.Serial.Buffered() > 0 {
				lastRxDataTimestamp = time.Now()
				state = processingData
			}
		case processingData:
			if time.Since(lastRxDataTimestamp) >= rxWaitTime {
				key, _ := machine.Serial.ReadByte()
				if machine.Serial.Buffered() > 0 && isDigit(peekByte()) {
					value := readInt()
					switch key {
					case 'r':
						redDimmer.SetDutyCycle(value)
					case 'g':
						greenDimmer.SetDutyCycle(value)
					case 'b':
						blueDimmer.SetDutyCycle(value)
					default:
						machine.Serial.Write([]byte("Invalid key received. Clearing all colors.\r\n"))
						redDimmer.SetDutyCycle(0)
						greenDimmer.SetDutyCycle(0)
						blueDimmer.SetDutyCycle(0)
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

// One-byte lookahead buffer, since machine.UART has no Peek method.
var lookahead byte
var haveLookahead bool

func peekByte() byte {
	if !haveLookahead {
		lookahead, _ = machine.Serial.ReadByte()
		haveLookahead = true
	}
	return lookahead
}

func readInt() int {
	value := 0
	for {
		var b byte
		if haveLookahead {
			b = lookahead
			haveLookahead = false
		} else if machine.Serial.Buffered() > 0 {
			b, _ = machine.Serial.ReadByte()
		} else {
			break
		}
		if !isDigit(b) {
			break
		}
		value = value*10 + int(b-'0')
	}
	return value
}
