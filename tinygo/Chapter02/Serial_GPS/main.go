// Port of Chapter02/Serial_GPS/Serial_GPS.ino.
package main

import (
	"machine"

	"practicalarduinorobotics/tinygo/internal/softserial"
)

func main() {
	machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})
	// GPS receivers use a baud rate of 9600. Pin 10 is RX, pin 11 is TX.
	gpsSerial := softserial.New(machine.D10, machine.D11, 9600)
	gpsSerial.Configure()

	for {
		if gpsSerial.Available() {
			machine.Serial.WriteByte(gpsSerial.ReadByte())
		}
	}
}
