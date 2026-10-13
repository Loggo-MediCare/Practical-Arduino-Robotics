// Port of Chapter05/blink_and_print/blink_and_print.ino.
package main

import (
	"machine"
	"time"
)

func main() {
	machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	blinkInterval := 200 * time.Millisecond
	printInterval := 1000 * time.Millisecond
	lastBlinkTime := time.Now()
	lastPrintTime := time.Now()
	ledOn := false

	for {
		now := time.Now()
		// Check if it is time to run the blink task.
		if now.Sub(lastBlinkTime) >= blinkInterval {
			lastBlinkTime = lastBlinkTime.Add(blinkInterval)
			ledOn = !ledOn
			led.Set(ledOn)
		}
		// Check if it is time to run the print task.
		if now.Sub(lastPrintTime) >= printInterval {
			lastPrintTime = lastPrintTime.Add(printInterval)
			machine.Serial.Write([]byte("Tick\r\n"))
		}
	}
}
