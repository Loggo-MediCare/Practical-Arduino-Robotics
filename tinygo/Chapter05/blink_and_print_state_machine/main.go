// Port of Chapter05/blink_and_print_state_machine/blink_and_print_state_machine.ino.
package main

import (
	"machine"
	"time"
)

type ledState uint8

const (
	ledOff ledState = iota
	ledOn
)

var (
	blinkInterval = 200 * time.Millisecond
	printInterval = 1000 * time.Millisecond
	state         = ledOff
	led           = machine.LED
)

func main() {
	machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	lastBlinkTime := time.Now()
	lastPrintTime := time.Now()

	for {
		now := time.Now()
		if now.Sub(lastBlinkTime) >= blinkInterval {
			lastBlinkTime = lastBlinkTime.Add(blinkInterval)
			blinkTask()
		}
		if now.Sub(lastPrintTime) >= printInterval {
			lastPrintTime = lastPrintTime.Add(printInterval)
			machine.Serial.Write([]byte("Tick\r\n"))
		}
	}
}

func blinkTask() {
	switch state {
	case ledOff:
		led.High()
		state = ledOn
		blinkInterval = 50 * time.Millisecond
	case ledOn:
		led.Low()
		state = ledOff
		blinkInterval = 450 * time.Millisecond
	}
}
