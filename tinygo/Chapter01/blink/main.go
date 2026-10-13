// Port of Chapter01/blink/blink.ino.
package main

import (
	"machine"
	"time"
)

func main() {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})
	machine.Serial.Write([]byte("Hello, world!\r\n"))

	for {
		led.High()
		machine.Serial.Write([]byte("LED on.\r\n"))
		time.Sleep(500 * time.Millisecond)
		led.Low()
		machine.Serial.Write([]byte("LED off.\r\n"))
		time.Sleep(500 * time.Millisecond)
	}
}
